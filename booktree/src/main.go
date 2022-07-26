package main

import (
	"context"
	"net/http"
	"os"
	"tosinjs/go-booktree/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type application struct {
	bookCollection *mongo.Collection
	logger         *logger.Logger
}

func main() {
	logger := logger.New(os.Stdout, logger.LevelInfo)
	client := connectDB()
	defer client.Disconnect(context.TODO())

	bc := client.Database("booktree").Collection("books")

	app := &application{
		bookCollection: bc,
		logger:         logger,
	}

	if err := http.ListenAndServe(":3000", app.router()); err != nil {
		logger.PrintFatal(err, nil)
	}
}

func connectDB() *mongo.Client {
	logger := logger.New(os.Stdout, logger.LevelInfo)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("lmao"))
	if err != nil {
		logger.PrintFatal(err, nil)
		panic(err)
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		logger.PrintFatal(err, nil)
		panic(err)
	}
	logger.PrintInfo("Successfully Connected to the Database", nil)
	return client
}
