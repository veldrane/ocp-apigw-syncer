package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	syncer "syncer"
	checker "syncer/gen/checker"
	root "syncer/gen/root"
	"syscall"

	nginx "github.com/nginx"
)

func main() {
	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		hostF     = flag.String("host", "", "Server host (valid values: )")
		domainF   = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		httpPortF = flag.String("http-port", "", "HTTP port (overrides host HTTP port specified in service design)")
		secureF   = flag.Bool("secure", false, "Use secure scheme (https or grpcs)")
		dbgF      = flag.Bool("debug", false, "Log request and response bodies")
	)
	flag.Parse()

	// Setup logger. Replace logger with your own log package of choice.
	var (
		logger *log.Logger
	)
	{
		logger = log.New(os.Stderr, "[syncer] ", log.Ltime)
	}

	var (
		nginxs        nginx.NginxInstancies
		requestConfig nginx.RequestConfig
	)
	{
		requestConfig = nginx.RequestConfig{HostHeader: "api-apigwp-cz.t.dc1.cz.ipa.ifortuna.cz", Retries: 5, SyncTimeout: 20}

		nginxs = nginx.New()
		nginxs.Push(nginx.NginxInstance{Address: "127.0.0.1", Port: "8080"}, "ng-plus-apigw-6cc76b4d5-vxtvg")
		nginxs.Push(nginx.NginxInstance{Address: "127.0.0.11", Port: "8080"}, "ng-plus-apigw-6cc76b4d5-asdvg")
		nginxs.Push(nginx.NginxInstance{Address: "127.0.0.12", Port: "8080"}, "ng-plus-apigw-6cc76b4d5-rtypb")
		nginxs.Push(nginx.NginxInstance{Address: "127.0.0.13", Port: "8080"}, "ng-plus-apigw-6cc76b4d5-adfse")
	}

	// Initialize the services.
	var (
		checkerSvc checker.Service
		rootSvc    root.Service
	)
	{
		checkerSvc = syncer.NewChecker(&requestConfig, &nginxs, logger)
		rootSvc = syncer.NewRoot(logger)
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		checkerEndpoints *checker.Endpoints
		rootEndpoints    *root.Endpoints
	)
	{
		checkerEndpoints = checker.NewEndpoints(checkerSvc)
		rootEndpoints = root.NewEndpoints(rootSvc)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc1 := make(chan error)
	errc2 := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c1 := make(chan os.Signal, 1)
		signal.Notify(c1, syscall.SIGINT, syscall.SIGTERM)
		errc1 <- fmt.Errorf("%s", <-c1)
	}()

	// OCP Gatherer interupt handler
	go func() {
		c2 := make(chan os.Signal, 1)
		signal.Notify(c2, syscall.SIGINT, syscall.SIGTERM)
		errc2 <- fmt.Errorf("%s", <-c2)
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Start the servers and send errors (if any) to the error channel.
	switch *hostF {
	case "":
		{
			addr := "http://localhost:8080"
			u, err := url.Parse(addr)
			if err != nil {
				logger.Fatalf("invalid URL %#v: %s\n", addr, err)
			}
			if *secureF {
				u.Scheme = "https"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *httpPortF != "" {
				h, _, err := net.SplitHostPort(u.Host)
				if err != nil {
					logger.Fatalf("invalid URL %#v: %s\n", u.Host, err)
				}
				u.Host = net.JoinHostPort(h, *httpPortF)
			} else if u.Port() == "" {
				u.Host = net.JoinHostPort(u.Host, "80")
			}
			handleHTTPServer(ctx, u, checkerEndpoints, rootEndpoints, &wg, errc1, logger, *dbgF)
		}

	default:
		logger.Fatalf("invalid host argument: %q (valid hosts: )\n", *hostF)
	}

	handleBackgroundGatherer(ctx, logger, errc2)

	// Wait for signal.
	logger.Printf("Main (%v)", <-errc1)
	logger.Printf("Main (%v)", <-errc2)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	logger.Println("exited")
}
