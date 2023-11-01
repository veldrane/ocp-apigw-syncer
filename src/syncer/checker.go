package syncer

import (
	"context"
	"log"
	checker "syncer/gen/checker"
)

// checker service example implementation.
// The example methods log the requests and return zero values.
type checkersrvc struct {
	pods   *[]NginxInstance
	logger *log.Logger
}

// NewChecker returns the checker service implementation.
func NewChecker(pods *[]NginxInstance, logger *log.Logger) checker.Service {
	return &checkersrvc{pods, logger}
}

// Get last full report
func (s *checkersrvc) Get(ctx context.Context) (res *checker.Sync, err error) {
	res = &checker.Sync{}
	s.logger.Printf("checker.get %s\n", s.pods)
	return
}
