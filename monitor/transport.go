package monitor

import (
	"context"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"
	"net/http"
)

//MakeHandler returns a restful http handler for the file status service
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	fileStatusEndpoint := kithttp.NewServer(makeStatusRequestEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	setDatedGlobalState := kithttp.NewServer(makeSetGlobalStatusRequestEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	getDatedGlobalState := kithttp.NewServer(makeGetDatedGlobalStateRequestEndpoint(service), handleGetDatedGlobalStateRequest, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/fileStatus", fileStatusEndpoint).Methods("GET")
	r.Handle("/setGlobalState", setDatedGlobalState).Methods("GET")
	r.Handle("/backdated", getDatedGlobalState).Queries("date", "{date}").Methods("GET")

	return r

}

func handleGetDatedGlobalStateRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	var q = &datedGlobalStateRequest{Date: r.FormValue("date")}

	return q, nil

}
