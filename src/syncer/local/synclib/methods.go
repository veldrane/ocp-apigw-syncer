package synclib

import (
	"context"
	"crypto/tls"
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

// Push method is just for testing data purpose
func (n *NginxInstancies) Push(ng NginxInstance, hostname string) error {

	n.Lock.Lock()
	n.Pods[hostname] = ng
	n.Lock.Unlock()

	return nil
}

func (n *NginxInstancies) Check(config *Config, p CheckPayload, ctx context.Context, logger *log.Logger) (status string) {

	var wg sync.WaitGroup
	var err error
	wgStatusDone := make(chan interface{})

	logger.Printf("[ Check ] -> Checking sync status for auth_token %s ....", *p.authToken)
	pods := n.getPods(ctx)
	httpStatus := make([]int, len(pods))

	i := 0

	for k, v := range pods {
		if k == *p.origin {
			logger.Printf("[ Check ] -> Same origin %s - skipping\n", k)
			httpStatus[i] = 200
			i++
			continue
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup, hostname string, pod NginxInstance, httpCode *int) {
			defer wg.Done()

			logger.Printf("[ Check thread ] -> checking auth_token %s on hostname %s with address %s\n", *p.authToken, hostname, pod.Address)

			*httpCode, err = getTokenStatus(ctx, p.authToken, config, &pod, logger)
			if err != nil {
				logger.Printf("[ Check thread ] -> warning check auth_token %s on pod %s failed ", *p.authToken, hostname)

			}
		}(&wg, k, v, &httpStatus[i])
		i++
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		wgStatusDone <- true
	}(&wg)

	select {
	case <-wgStatusDone:
		logger.Printf("[ Check ] -> all threads are done sucessfully for token %s, status %d \n", *p.authToken, httpStatus)
		status = evalGroup(httpStatus)
	case <-ctx.Done():
		logger.Println("[ Check ] -> warning, timeout occured for token:", *p.authToken)
		status = "Timeout"
	}

	return status
}

func evalGroup(statusCodes []int) string {

	oks, errs := 0, 0
	res := "NotSynced"

	for _, v := range statusCodes {
		switch statusCode := v; statusCode {
		case 200:
			oks++
		default:
			errs++
		}
	}

	successRate := (float32(oks) / float32(len(statusCodes))) * 100

	if successRate > 99 {
		res = "Synced"
	} else if successRate > 50 {
		res = "Partialy"
	}

	return res
}

func InitCheckPayload(authToken string, origin string) CheckPayload {

	res := CheckPayload{
		authToken: &authToken,
		origin:    &origin,
	}

	return res
}

func (n *NginxInstancies) getPods(ctx context.Context) (res map[string]NginxInstance) {

	n.Lock.RLock()
	res = n.Pods
	n.Lock.RUnlock()

	return res
}

func getTokenStatus(ctx context.Context, token *string, config *Config, pod *NginxInstance, logger *log.Logger) (res int, err error) {

	//time.Sleep(3 * time.Second)

	w, err := os.OpenFile("/tmp/sslkey.out", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		logger.Printf("failed to open file err %+v", err)
	}
	defer w.Close()

	client, err := initHttpClient(config, w, false)
	req, err := initHttpRequest(ctx, token, config, pod)
	if err != nil {
		logger.Printf("[ Get Token Status ] -> Failed create request with context with err %s\n", err)
		return 0, err
	}

	for i := 1; i <= config.Retries; i++ {

		//time.Sleep(2 * time.Second)

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
			err = fmt.Errorf("%d", statusCode)
			return res, err
		}
	}

	logger.Printf("[ Get Token Status ] -> Token %s on pod %s is not synced in the end. Exititng with 401", *token, pod.Address)
	return 401, nil
}

func initHttpClient(config *Config, w *os.File, debug bool) (*http.Client, error) {

	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost:   16,
		MaxIdleConns:          256,
		ResponseHeaderTimeout: (time.Duration(config.ConnTimeout) / 2) * time.Millisecond,
	}

	if debug {
		tr.TLSClientConfig.KeyLogWriter = w
	}

	client := http.Client{
		Transport: tr,
		Timeout:   time.Duration(config.ConnTimeout) * time.Millisecond,
	}

	return &client, nil
}

func initHttpRequest(ctx context.Context, token *string, config *Config, pod *NginxInstance) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://"+pod.Address+":"+pod.Port+config.HttpPath, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", "auth_token="+*token+"; Domain="+config.HostDomain+"; Path=/; SameSite=Strict; HttpOnly; Secure;")
	req.Host = config.HostHeader

	return req, nil
}

func IsChanged(ocpPods map[string]NginxInstance, storedPods map[string]NginxInstance, logger *log.Logger) bool {

	if numberOfOcpPods := len(ocpPods); numberOfOcpPods != len(storedPods) {
		//logger.Printf("[ Scraping thread ] -> New pods detected %d %d", len(ocpPods), len(storedPods))
		return true
	}

	for i := range ocpPods {
		if storedPods[i] == (NginxInstance{}) {
			//logger.Printf("[ Scraping thread ] -> Pods changed! New configuration will be stored!")
			return true
		}
	}

	return false
}

func (n *NginxInstancies) Update(ngs map[string]NginxInstance, logger *log.Logger) error {

	n.Lock.Lock()
	// time.Sleep(10 * time.Second) - mutex test - sharing between client and scraper thread
	for k := range n.Pods {
		delete(n.Pods, k)
	}

	for k, v := range ngs {
		if v.Address != "" {
			n.Pods[k] = v
		}
	}

	n.Lock.Unlock()

	return nil
}
