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

func makeSetGlobalStatusRequestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		s.storeGlobalStateDaily()

		return nil, nil
	}
}

func makeGetDatedGlobalStateRequestEndpoint(s Service) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(*datedGlobalStateRequest)

		resp, err := s.getDatedGlobalStateDaily(req.Date)

		return resp, nil
	}
}

type datedGlobalStateRequest struct {
	Date string
}
