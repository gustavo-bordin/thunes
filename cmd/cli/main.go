package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gustavo-bordin/thunes/config"
	"github.com/gustavo-bordin/thunes/internal/cli"
	"github.com/gustavo-bordin/thunes/internal/ngrok"
	"github.com/gustavo-bordin/thunes/internal/repository"
	"github.com/gustavo-bordin/thunes/internal/thunes"
)

func main() {
	cfg := config.NewCliConfig()
	cfg.Load()

	db := repository.NewMongoDB(cfg.Mongo)
	transactionRepo := repository.NewTransactionRepository(db)

	ngrok := ngrok.NewNgrok(cfg)
	ngrokUrl, err := ngrok.GetNgrokUrl()
	if err != nil {
		panic(err)
	}

	tc := thunes.NewClient(cfg)
	cli := cli.NewRootScreen(tc, transactionRepo, *ngrokUrl)

	p := tea.NewProgram(cli, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("error on executing program %s", err)
	}
}
