package monitor

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makeStatusRequestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		status := "poop"

		return status, nil
	}
}
