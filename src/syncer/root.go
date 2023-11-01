package syncer

import (
	"context"
	"log"
	root "syncer/gen/root"
)

// root service example implementation.
// The example methods log the requests and return zero values.
type rootsrvc struct {
	logger *log.Logger
}

// NewRoot returns the root service implementation.
func NewRoot(logger *log.Logger) root.Service {
	return &rootsrvc{logger}
}

// Return default redirect
func (s *rootsrvc) Default(ctx context.Context) (res *root.DefaultResult, err error) {
	res = &root.DefaultResult{}
	s.logger.Print("root.default")
	return
}
