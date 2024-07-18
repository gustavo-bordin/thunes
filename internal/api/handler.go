package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gustavo-bordin/thunes/internal/repository"
	"github.com/gustavo-bordin/thunes/internal/thunes"
)

type handler struct {
	transactionRepo repository.Repository
}

func NewHandler(
	transactionRepo repository.Repository,
) handler {
	return handler{transactionRepo: transactionRepo}
}

func (h handler) handleTransactionState(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling transaction state")

	var transaction thunes.TransactionState
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.transactionRepo.Insert(ctx, transaction)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Saved transaction state")
	w.WriteHeader(http.StatusCreated)
}
