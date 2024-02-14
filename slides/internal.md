---
title: Syncer internal
separator: <!--s-->
verticalSeparator: <!--v-->
revealOptions:
transition: 'none'
---


<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=left2-small>

#### General

- main two parts 

        - main part is writen in Go (mainly this presentation)
        - second part is integrated in api gateway

- skeleton is writen in Goa framework

        - http server
        - logging
        - root contexts

</div>

<!--s-->

<!-- .slide: data-background="images/syncer-src-description.png" data-background-size="1920px" -->

<!--s-->


<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=left2-small>

#### Design goa files

</div>

<!--s-->

<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=left2-small>

#### Design goa files

- describes api
- contains http, skelletons of handlers
- defines input and output header
- important file checker.go
- other files contain link to swagger, some definition etc

<BR>
<BR>

<div id=resources>

[https://github.com/goadesign/goa](https://github.com/goadesign/goa)<BR>
[https://pkg.go.dev/goa.design/goa/v3/dsl](https://pkg.go.dev/goa.design/goa/v3/dsl)</BR>

</div>

</div>

<div id=right-small>

```bash
mkdir -p syncer/design
cp ocp-apigwp-syncer/src/syncer/design syncer/design
cd syncer
go mod init syncer
goa gen syncer/design
goa example syncer/design
go get goa.design/goa/v3/http@v3.14.1
go build cmd/syncer/* -o syncer
./syncer
```

</div>

<!--s-->

<!-- .slide: data-background="images/check-service-design-file.png" data-background-size="1920px" -->

<!--s-->

<!-- .slide: data-background="images/syncer-architecture-background-black-text-1.png" data-background-size="1920px" -->

#### General Syncer architecture

<!--s-->

<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=two-columns-black>

#### General Syncer architecture

</div>

<BR>

<div id=left2-small>

<BR>

- the scraper

    - communicate with k8s/openshift api
    - update NginxInstancies based on the running pod
    - uses internal library synclib (src/syncer/local/synclib)
    - entry point is called from src/syncer/cmd/syncer/main.go
    - main code is writen in src/syncer/cmd/syncer/background.go and function handleBackgroundGatherer

</div>

<div id=right2-small>

<BR>

- http handler

    - entrypoint is a product of goa framework 
        - (src/syncer/checker.go) function "Get"
    - <div id=important>main code and most important function is "Check" </div>
        - library synclib (src/syncer/local/synclib) methods.go
    - other code mostly in synclib library

</div>


<!--s-->

<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->


<div id=two-columns-black>

#### Scraper process

</div>

<BR>

<div id=left2-small>

- separate gorutine periodically checks the ocp api

- if detects discrepancy between ocp and []nginxInstancies then:
    - locks the instance of the nginx instancies
    - rewrite nginx instancies based on the oc api

<BR>
<BR>

- init of the scraper process is done by main function in src/syncer/cmd/syncer/main.go

<BR>
<BR>

```go
# src/syncer/cmd/main.go
func main() {
.
.
	var (
		nginxs l.NginxInstancies
		config l.Config
	)
	{
		config = l.GetConfig() #parse config files and fill config structures
		nginxs = l.New() #Initialize emtpy NginxInstancies
	}
.
	handleBackgroundGatherer(ctx, &nginxs, &config, logger, errc2)
.
}
```

</div>

<div id=right2-small>

<BR>
<BR>
<BR>
<BR>
<BR>
<BR>
<BR>
<BR>
<BR>
<BR>
<BR>

```go
# src/syncer/cmd/syncer/background.go
func handleBackgroundGatherer(ctx context.Context, pods *l.NginxInstancies, config *l.Config, logger *log.Logger, errc chan error) {
	go func() {
		logger.Printf("[ Scraping thread ] -> Started sucessfully with period %s seconds", strconv.Itoa(10))
		ocpSession := ocp4cli.Session()
		ctx := context.Background()
		go func() {
			for {
				runningPods, _ := ocpSession.GetPods(&ctx, config)
				//logger.Printf("[ Scraping thread ] -> Waking up, checking ocp configuration....")
				if l.IsChanged(runningPods, pods.Pods, logger) {
					pods.Update(runningPods, logger)
					var pl string
					for k := range runningPods {
						pl = pl + k + " "
					}
					logger.Printf("[ Scraping thread ] -> Pods updated: %s", pl)
				}
				time.Sleep(time.Duration(10) * time.Second)
			}
		}()
		errc <- fmt.Errorf("%s", "[ Scraping thread ] -> scraping thread is dead baby")
	}()
}
```

</div>

<!--s-->

<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=two-columns-black>

#### Scraper process - how scraper looks for the right pods ?

</div>

<BR>

<div id=left2-small>

- scraper gets the latest revision of the apigw deployment

- go through all replication set and found one with this revision number

- from replicationset name gets the pod hash

- then is able to get from ocp api the list of current running pods

</div>

<div id=right2-small>

```bash
$ oc get deployment ng-plus-apigw -o yaml | grep revision
    deployment.kubernetes.io/revision: "73"
$ oc get rs ng-plus-apigw-5d5f657889 -o yaml | grep revision
    deployment.kubernetes.io/revision: "73"
$ oc get pods
NAME                                     READY   STATUS    RESTARTS        AGE
ng-plus-apigw-5d5f657889-klc9r           2/2     Running   0               55m
```
<BR>

```go
# src/syncer/local/ocp4cli/public.go
func (session *SessionT) GetPods(ctx *context.Context, config *l.Config) (map[string]l.NginxInstance, error) {

	podList := l.New()
	var pod l.NginxInstance

	revision, err := session.getDeploymentRevision(ctx, &config.Deployment, &config.Namespace)
	if err != nil {
		panic(err)
	}

	replicaSet, _ := session.getRsBasedOnRevision(ctx, &revision, &config.Namespace, &config.Deployment)
	if err != nil {
		panic(err)
	}

	rse := strings.Split(replicaSet, "-")
	podsHash := rse[len(rse)-1]
```

</div>

<!--s-->


<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=two-columns-black>

#### Scraper process - check discrepancy and update NginxInstancies

</div>

<BR>

<div id=left2-small>

- function isChanged():
        - compare number of pods
        - compare if the name of the pods are different
        - if any of that is different the return true
<BR>

- function Update():
        - lock the NginxInstancies for read/write. Any http request must wait when update is finished.
        - delete the NginxInstancies
        - write new pods
        - unlock NginxInstancies read/write

</div>

<BR>

<div id=right2-small>

```go
# src/syncer/local/synclib/methods.go
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
```
<BR>

```go
# src/syncer/local/synclib/methods.go
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
```

</div>

<!--s-->


<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=left2-small>

#### http handler part

</div>

<!--s-->


<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=two-columns-black>

#### http handler part - checker.go

</div>

<BR>

<div id=left2-small>

- skell is the product of the goa framework

- service has to be inicialized - function NewChecksrv and type checksrvc
	- added pointer to NginxInstancies function
	- added pointer to the configguration type

- logger is pointed by default by goa
- all services defined by goa design files are initialized by src/syncer/cmd/syncer/main.go

</div>

<div id=right2-small>

```go
# src/syncer/checker.go
// checker service example implementation.
// The example methods log the requests and return zero values.
type checkersrvc struct {
	requestConfig	*l.Config
	nginxs	*l.NginxInstancies
	logger	*log.Logger
}

// NewChecker returns the checker service implementation.
// Function return basic logger interface, list of the nginx pods and structure
// describes global config
func NewChecker(requestConfig *l.Config, nginxs *l.NginxInstancies, logger *log.Logger) checker.Service {
	return &checkersrvc{requestConfig, nginxs, logger}
}
```

</div>

<!--s-->

<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=two-columns-black>

#### http handler part - checker.go, function Get()

</div>

<div id=left2-small>

- this function is called when apigw calls /v1/synced

- checks if apigw uses replicas, if not return status "synced" and finishes

- take the payload(from http headers) and wrape arround the special structure
	- ugly, blee (i expect more parameters, but its better to pass one structure than multiple variables)

- call the Check() in src/syncer/local/methods.go and store result in status variable

</div>

<div id=right2-small>

```go
# src/syncer/checker.go
// Get last full report Main handler of the endpoint /v1/synced
// Endpoint return one of the four statuses: Synced, Partialy, NotSynced, Timeout
func (s *checkersrvc) Get(ctx context.Context, p *checker.GetPayload) (res *checker.Sync, err error) {

	if numPods := len(s.nginxs.Pods); numPods < 2 {
		s.logger.Printf("[ Get ] -> Not sync check required for token %s - nginx pods doesnt have replicas", p.AuthToken)
		return &checker.Sync{Status: "Synced"}, nil
	}

	cp := l.InitCheckPayload(p.AuthToken, p.Origin)

	ctxCheck, cancel := context.WithTimeout(ctx, time.Duration(s.requestConfig.Deadline*int(time.Millisecond)))
	defer cancel()

	status := s.nginxs.Check(s.requestConfig, cp, ctxCheck, s.logger)

	return &checker.Sync{Status: status}, nil
}
```

</div>

<!--s-->

<!-- .slide: data-background="images/syncer-check-function-2.png" data-background-size="1920px" -->

<!--s-->

<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=two-columns-black>

#### HTTP Handler - function getTokenStatus()

</div>

<div id=left2-small>

- called inside the check gorutine

- first init http client then request itself inside loop

- support for ssl dump keys in case of the debug

</div>

<div id=right2-small>

```go 
# src/syncer/local/synclib/methods.go
func getTokenStatus(ctx context.Context, token *string, config *Config, pod *NginxInstance, logger *log.Logger) (res int, err error) {
.
.
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
			err = fmt.Errorf("%d", statusCode)
			return res, err
		}
	}
```
</div>

<!--s-->

<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->



<div id=left2-small>

#### Used golang design paterns and features

</div>

<!--s-->

<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=left2-small>

#### Gorutines - paralelism in golang

- key word "go" when calling async piece of code

- part of the golang language

- own golang async runtime! Run in userspace (not context switches!)

- fast but small memory footprint

- for communication between multiple gorutines channels are used in general

</div>

<div id=right2-small>

```go
# src/syncer/local/synclib/methods.go
.
.
go func(wg *sync.WaitGroup, hostname string, pod NginxInstance, httpCode *int) {
    defer wg.Done()

    logger.Printf("[ Check thread ] -> checking auth_token %s on hostname %s with address %s\n", *p.authToken, hostname, pod.Address)

    *httpCode, err = getTokenStatus(ctx, p.authToken, config, &pod, logger)
    if err != nil {
        logger.Printf("[ Check thread ] -> warning check auth_token %s on pod %s failed ", *p.authToken, hostname)

    }
}(&wg, k, v, &httpStatus[i])
.
.
```

</div>

<!--s-->

<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=left2-small>

#### Waitgroups

- used for synchronization multiple gorutines

- common patern: 
    - before start element to waitgroup is added

    - "defer" function ensures remove element from waitgroup when gorutine is done

    - wg.wait waits until all gorutines are done

</div>

<div id=right2-small>

```go
# src/syncer/local/synclib/methods.go
var wg sync.WaitGroup
wg.Add(1)
go func(wg *sync.WaitGroup, hostname string, pod NginxInstance, httpCode *int) {
    defer wg.Done() // when function finishes calls wg.Done

    logger.Printf("[ Check thread ] -> checking auth_token %s on hostname %s with address %s\n", *p.authToken, hostname, pod.Address)

    *httpCode, err = getTokenStatus(ctx, p.authToken, config, &pod, logger)
    if err != nil {
        logger.Printf("[ Check thread ] -> warning check auth_token %s on pod %s failed ", *p.authToken, hostname)

    }
}(&wg, k, v, &httpStatus[i]) // we have to pass reference to waiting group for calling by defer

go func(wg *sync.WaitGroup) {
    wg.Wait()   // stop code in the gorutine and wait until whole waitgroup is done, then continue 
    wgStatusDone <- true
}(&wg)

```

</div

<!--s-->

<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=left2-small>

#### Contexts

- uses for distuinguish between instancies of the same objects

- typicaly we have http request and we need to point to specific attributes of that (headers, parameters, body etc)

- context can be chained, so (for example) we can propagate to specific http request information from http server instance who has a root context

- has a ability to cancel action and clean all objects in chain

- for example [server] -> [http request] -> [spread gorutines] -> [http client request]

- if context on http request is canceled, all following requested (gorutines, client) are clean canceled as well. 

- syncer uses contexts for canceling long http client requests

</div>

<div id=right2-small>

```go
# src/syncer/check.go
func (s *checkersrvc) Get(ctx context.Context, p *checker.GetPayload) (res *checker.Sync, err error) {
    .
	ctxCheck, cancel := context.WithTimeout(ctx, time.Duration(s.requestConfig.Deadline*int(time.Millisecond)))
	defer cancel()
    .
```

```go
# src/syncer/local/synclib/methods.go
func (n *NginxInstancies) Check(config *Config, p CheckPayload, ctx context.Context, logger *log.Logger) (status string) {
    .
    .
	logger.Printf("[ Check ] -> Checking sync status for auth_token %s ....", *p.authToken)
	pods := n.getPods(ctx)
    .
	for k, v := range pods {
        .
		go func(wg *sync.WaitGroup, hostname string, pod NginxInstance, httpCode *int) {
			defer wg.Done()

			logger.Printf("[ Check thread ] -> checking auth_token %s on hostname %s with address %s\n", *p.authToken, hostname, pod.Address)

			*httpCode, err = getTokenStatus(ctx, p.authToken, config, &pod, logger)
			if err != nil {
				logger.Printf("[ Check thread ] -> warning check auth_token %s on pod %s failed ", *p.authToken, hostname)

			}
		}(&wg, k, v, &httpStatus[i])

	}
    .
	select {
	case <-wgStatusDone:
		logger.Printf("[ Check ] -> all threads are done sucessfully for token %s, status %d \n", *p.authToken, httpStatus)
		status = evalGroup(httpStatus)
	case <-ctx.Done():
		logger.Println("[ Check ] -> warning, timeout occured for token:", *p.authToken)
		status = "Timeout"
	}
    .   
	return status
}
```

</div>

<!--s-->
