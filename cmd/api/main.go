package main

import (
	"github.com/gustavo-bordin/thunes/config"
	"github.com/gustavo-bordin/thunes/internal/api"
	"github.com/gustavo-bordin/thunes/internal/repository"
)

func main() {
	cfg := config.NewApiConfig()
	cfg.Load()

	db := repository.NewMongoDB(cfg.Mongo)
	transactionRepo := repository.NewTransactionRepository(db)
	handler := api.NewHandler(transactionRepo)
	stateApi := api.NewStateApi(handler, cfg)
	stateApi.Start()
}
