package repository

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gustavo-bordin/thunes/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var once sync.Once
var db *mongo.Database

type MongoDB struct {
	DB *mongo.Database
}

func NewMongoDB(cfg config.MongoConfig) MongoDB {
	once.Do(func() {
		timeout := 10 * time.Second
		ctx, cancel := context.WithTimeout(context.TODO(), timeout)
		defer cancel()

		options := options.Client().ApplyURI(cfg.Url)
		client, err := mongo.Connect(ctx, options)
		if err != nil {
			log.Fatal(err)
		}

		db = client.Database(cfg.DbName)
	})

	return MongoDB{DB: db}
}
