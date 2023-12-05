package syncer

import (
	"context"
	"log"
	health "syncer/gen/health"
)

// health service example implementation.
// The example methods log the requests and return zero values.
type healthsrvc struct {
	logger *log.Logger
}

// NewHealth returns the health service implementation.
func NewHealth(logger *log.Logger) health.Service {
	return &healthsrvc{logger}
}

// Ping endpoin
func (s *healthsrvc) Get(ctx context.Context) (res *health.Health, err error) {
	res = &health.Health{Status: "Up"}
	s.logger.Print("health.get")
	return
}
