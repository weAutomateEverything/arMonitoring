package monitor


import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

func makeStatusRequestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		status := s.

		return status, nil
	}
}
