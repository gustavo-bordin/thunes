package repository

import (
	"context"

	"github.com/gustavo-bordin/thunes/internal/thunes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Insert(ctx context.Context, transaction thunes.TransactionState) error
	Find(ctx context.Context, filter primitive.M) ([]thunes.TransactionState, error)
}

type TransactionRepository struct {
	db         MongoDB
	collection string
}

func NewTransactionRepository(db MongoDB) TransactionRepository {
	return TransactionRepository{
		db:         db,
		collection: "transaction",
	}
}

func (r TransactionRepository) Insert(
	ctx context.Context,
	transaction thunes.TransactionState,
) error {
	_, err := r.db.DB.Collection(r.collection).InsertOne(ctx, transaction)
	return err
}

func (r TransactionRepository) Find(
	ctx context.Context,
	filter primitive.M,
) ([]thunes.TransactionState, error) {
	cursor, err := r.db.DB.Collection(r.collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var transactions []thunes.TransactionState

	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}
