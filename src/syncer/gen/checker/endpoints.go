// Code generated by goa v3.13.2, DO NOT EDIT.
//
// checker endpoints
//
// Command:
// $ goa gen syncer/design

package checker

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "checker" service endpoints.
type Endpoints struct {
	Get goa.Endpoint
}

// NewEndpoints wraps the methods of the "checker" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Get: NewGetEndpoint(s),
	}
}

// Use applies the given middleware to all the "checker" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Get = m(e.Get)
}

// NewGetEndpoint returns an endpoint function that calls the method "get" of
// service "checker".
func NewGetEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		res, err := s.Get(ctx)
		if err != nil {
			return nil, err
		}
		vres := NewViewedSync(res, "default")
		return vres, nil
	}
}
