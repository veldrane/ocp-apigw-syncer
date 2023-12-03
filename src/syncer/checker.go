package syncer

import (
	"context"
	"log"
	checker "syncer/gen/checker"
	"time"

	nginx "github.com/nginx"
)

// checker service example implementation.
// The example methods log the requests and return zero values.
type checkersrvc struct {
	requestConfig *nginx.Config
	nginxs        *nginx.NginxInstancies
	logger        *log.Logger
}

// NewChecker returns the checker service implementation.
func NewChecker(requestConfig *nginx.Config, nginxs *nginx.NginxInstancies, logger *log.Logger) checker.Service {
	return &checkersrvc{requestConfig, nginxs, logger}
}

// Get last full report
func (s *checkersrvc) Get(ctx context.Context, p *checker.GetPayload) (res *checker.Sync, err error) {

	status := "synced"

	if len(s.nginxs.Pods) < 2 {
		s.logger.Printf("Not sync check required for token %s - nginx pods doesnt have replicas", p.AuthToken)
		return &checker.Sync{Status: status}, nil
	}

	cp := nginx.InitCheckPayload(p.AuthToken, p.Origin)

	ctxCheck, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err = s.nginxs.Check(s.requestConfig, cp, ctxCheck, s.logger)

	//s.logger.Printf("Context found %s\n", ctx)

	//if err != nil {
	//	status = "not_synced"
	//}

	return &checker.Sync{Status: err.Error()}, nil
}
