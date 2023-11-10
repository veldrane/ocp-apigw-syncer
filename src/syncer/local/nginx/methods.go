package nginx

import (
	"context"
	"errors"
	"log"
	"math/rand"
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

	logger.Printf("[ Check ] -> Checking sync status for auth_token %s ....", *p.token)

	pods := n.getPods(ctx)

	for k, v := range pods {
		if k == *p.origin {
			logger.Printf("[ Check ] -> Same origin %s - skipping\n", k)
			continue
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup, hostname string, pod NginxInstance) {
			defer wg.Done()
			logger.Printf("[ Check thread ] -> checking auth_token %s on hostname %s with address %s\n", *p.token, hostname, pod.Address)

			getTokenStatus(ctx, p.token, config, &pod, logger)
			time.Sleep(time.Duration(rand.Intn(3096)) * time.Millisecond)
		}(&wg, k, v)
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		status <- true
	}(&wg)

	select {
	case <-status:
		logger.Println("[ Check ] -> all threads are done sucessfully for token ", *p.token)
		err = nil
	case <-ctx.Done():
		logger.Println("[ Check ] -> warning, timeout occured for token:", *p.token)
		err = errors.New("timeout")
	}

	return err
}

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

func getTokenStatus(ctx context.Context, token *string, config *RequestConfig, pod *NginxInstance, logger *log.Logger) (err error) {

	//ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(200*time.Millisecond))
	//defer cancel()
	//req, err := http.NewRequestWithContext(ctx, "GET", "http://"+hostname+":8080", body)

	logger.Printf("[ Get Token Status ] -> trying connect to %s:%s ", pod.Address, pod.Port)

	return nil
}
