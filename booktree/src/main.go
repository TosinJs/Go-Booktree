package main

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type application struct {
	bookCollection *mongo.Collection
}

func main() {

	client := connectDB()
	defer client.Disconnect(context.TODO())

	bc := client.Database("booktree").Collection("books")

	app := &application{
		bookCollection: bc,
	}

	if err := http.ListenAndServe(":3000", app.router()); err != nil {
		log.Fatal(err)
	}
}

func connectDB() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("lmaofraudman"))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully Connected to the Database")
	return client
}
