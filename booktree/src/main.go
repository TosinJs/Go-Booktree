package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"tosinjs/go-booktree/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoURI string
var port string

type application struct {
	bookCollection *mongo.Collection
	userCollection *mongo.Collection
	logger         *logger.Logger
}

func main() {
	flag.StringVar(&mongoURI, "mongoURI", "", "Enter a Valid MongoURI")
	flag.StringVar(&port, "port", "8080", "HTTP Port to Start the Server")
	flag.Usage()
	flag.Parse()
	logger := logger.New(os.Stdout, logger.LevelInfo)
	client := connectDB()
	defer client.Disconnect(context.TODO())

	bc := client.Database("booktree").Collection("books")
	uc := client.Database("booktree").Collection("users")

	app := &application{
		bookCollection: bc,
		userCollection: uc,
		logger:         logger,
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), app.router()); err != nil {
		logger.PrintFatal(err, nil)
	}
}

func connectDB() *mongo.Client {
	logger := logger.New(os.Stdout, logger.LevelInfo)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
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
