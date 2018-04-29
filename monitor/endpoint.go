package monitor

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makeStatusRequestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		status := s.StatusResults()

		return status, nil
	}
}
