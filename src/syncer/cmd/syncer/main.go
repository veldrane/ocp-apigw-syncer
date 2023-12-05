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
	health "syncer/gen/health"
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
		nginxs nginx.NginxInstancies
		config nginx.Config
	)
	{
		config = nginx.Config{
			HostHeader:  "api-apigwp-cz.t.dc1.cz.ipa.ifortuna.cz",
			HttpPath:    "/check",
			HostDomain:  "ifortuna.cz",
			HttpsPort:   "8443",
			Deployment:  "ng-plus-apigw",
			Namespace:   "apigwp-cz",
			Retries:     5,
			SyncTimeout: 100}

		nginxs = nginx.New()
	}

	// Initialize the services.
	var (
		checkerSvc checker.Service
		healthSvc  health.Service
		rootSvc    root.Service
	)
	{
		checkerSvc = syncer.NewChecker(&config, &nginxs, logger)
		rootSvc = syncer.NewRoot(logger)
		healthSvc = syncer.NewHealth(logger)
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		checkerEndpoints *checker.Endpoints
		rootEndpoints    *root.Endpoints
		healthEndpoints  *health.Endpoints
	)
	{
		checkerEndpoints = checker.NewEndpoints(checkerSvc)
		rootEndpoints = root.NewEndpoints(rootSvc)
		healthEndpoints = health.NewEndpoints(healthSvc)
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
			addr := "http://0.0.0.0:8080"
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
			handleHTTPServer(ctx, u, healthEndpoints, checkerEndpoints, rootEndpoints, &wg, errc1, logger, *dbgF)
		}

	default:
		logger.Fatalf("invalid host argument: %q (valid hosts: )\n", *hostF)
	}

	handleBackgroundGatherer(ctx, &nginxs, &config, logger, errc2)

	// Wait for signal.
	logger.Printf("Main (%v)", <-errc1)
	logger.Printf("Main (%v)", <-errc2)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	logger.Println("exited")
}
