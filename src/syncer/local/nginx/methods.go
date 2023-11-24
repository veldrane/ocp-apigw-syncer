package nginx

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func New() NginxInstancies {

	return NginxInstancies{Pods: make(map[string]NginxInstance)}
}

func (n *NginxInstancies) Push(ng NginxInstance, hostname string) error {

	n.Lock.Lock()
	n.Pods[hostname] = ng
	n.Lock.Unlock()

	return nil
}

func (n *NginxInstancies) Delete(hostname string) error {

	n.Lock.Lock()
	delete(n.Pods, hostname)
	n.Lock.Unlock()

	return nil
}

func (n *NginxInstancies) Check(config *RequestConfig, p CheckPayload, ctx context.Context, logger *log.Logger) (err error) {

	var wg sync.WaitGroup
	status := make(chan interface{})
	defer close(status)

	logger.Printf("[ Check ] -> Checking sync status for auth_token %s ....", *p.token)
	pods := n.getPods(ctx)
	httpStatus := make([]int, len(pods))

	i := 0

	for k, v := range pods {
		if k == *p.origin {
			logger.Printf("[ Check ] -> Same origin %s - skipping\n", k)
			continue
		}

		wg.Add(1)

		go func(wg *sync.WaitGroup, hostname string, pod NginxInstance, httpCode *int) {
			defer wg.Done()
			logger.Printf("[ Check thread ] -> checking auth_token %s on hostname %s with address %s\n", *p.token, hostname, pod.Address)

			*httpCode, err = getTokenStatus(ctx, p.token, config, &pod, logger)

			if err != nil {
				logger.Printf("[ Check thread ] -> warning check auth_token %s on pod %s failed ", *p.token, hostname)

			}
		}(&wg, k, v, &httpStatus[i])
		i++
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		status <- true
	}(&wg)

	select {
	case <-status:
		logger.Printf("[ Check ] -> all threads are done sucessfully for token %s, status %d \n", *p.token, httpStatus)
		err = nil
	case <-ctx.Done():
		logger.Println("[ Check ] -> warning, timeout occured for token:", *p.token)
		err = errors.New("timeout")
	}

	return err
}

//func evalGroup(statusCodes []*int) error {

//	for _, v := range statusCodes {

//	}

//	return nil
//}

func InitCheckPayload(token string, origin string) CheckPayload {

	res := CheckPayload{
		token:  &token,
		origin: &origin,
	}

	return res
}

func (n *NginxInstancies) getPods(ctx context.Context) (res map[string]NginxInstance) {

	n.Lock.RLock()
	res = n.Pods
	n.Lock.RUnlock()

	return res
}

func getTokenStatus(ctx context.Context, token *string, config *RequestConfig, pod *NginxInstance, logger *log.Logger) (res int, err error) {

	w, err := os.OpenFile("/tmp/sslkey.out", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		logger.Printf("failed to open file err %+v", err)
	}
	defer w.Close()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://"+pod.Address+":"+pod.Port+config.HttpPath, nil)
	if err != nil {
		logger.Printf("[ Get Token Status ] -> Failed new request wirh context with error %s\n", err)
		return 0, err
	}

	req.Header.Set("Cookie", "auth_token="+*token+"; Domain="+config.HostDomain+"; Path=/; SameSite=Strict; HttpOnly; Secure;")
	req.Host = config.HostHeader

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true, KeyLogWriter: w},
	}

	client := http.Client{
		Transport: tr,
		Timeout:   500 * time.Millisecond,
	}

	for i := 1; i <= config.Retries; i++ {

		resp, err := client.Do(req)
		if err != nil {
			logger.Printf("[ Get Token Status ] -> Failed client.Do with err %s\n", err)
			return 0, err
		}

		switch statusCode := resp.StatusCode; statusCode {
		case 200:
			return resp.StatusCode, nil
		case 401:
			logger.Printf("[ Get Token Status ] -> Token %s on pod %s not sync... retry\n", *token, pod.Address)
			time.Sleep(time.Duration(config.SyncTimeout * int(time.Millisecond)))
			continue
		default:
			res = 0
			err = errors.New(fmt.Sprint(statusCode))
			return res, err
		}

	}

	logger.Printf("[ Get Token Status ] -> Token %s on pod %s is not synced finaly. Exititng with 401", *token, pod.Address)

	return 401, nil
}
