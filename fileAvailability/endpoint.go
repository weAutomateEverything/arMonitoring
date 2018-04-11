package fileAvailability

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makeStatusRequestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		status := s.CreateJSONResponse()

		return status, nil
	}
}
