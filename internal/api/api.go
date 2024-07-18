package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gustavo-bordin/thunes/config"
)

type stateApi struct {
	handler handler
	port    int16
}

func NewStateApi(handler handler, cfg config.ApiConfig) stateApi {
	return stateApi{handler: handler, port: cfg.Port}
}

func (s stateApi) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handler.handleTransactionState)

	log.Printf("Server is starting in %d", s.port)
	port := fmt.Sprintf(":%d", s.port)

	if err := http.ListenAndServe(port, mux); err != nil {
		panic(err)
	}
}
