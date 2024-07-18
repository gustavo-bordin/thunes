package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gustavo-bordin/thunes/internal/thunes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockRepository struct {
	insertCalled bool
}

func (m *mockRepository) Find(
	ctx context.Context,
	filter primitive.M,
) ([]thunes.TransactionState, error) {
	return nil, nil
}

func (m *mockRepository) Insert(
	ctx context.Context,
	transaction thunes.TransactionState,
) error {
	m.insertCalled = true
	return nil
}

func TestHandleTransactionState(t *testing.T) {
	mockRepo := &mockRepository{}
	handler := NewHandler(mockRepo)

	transaction := thunes.TransactionState{
		ID:     1,
		Status: 1,
	}

	payload, err := json.Marshal(transaction)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest("POST", "/states", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.handleTransactionState(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rec.Code)
	}

	if !mockRepo.insertCalled {
		t.Error("Expected Insert method to be called on repository")
	}
}
