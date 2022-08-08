package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"tosinjs/gin-booktree/src/handlers"
	"tosinjs/gin-booktree/src/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var booksHandler *handlers.BooksHandler
var usersHandler *handlers.UsersHandler
var mongoURI string
var port string

func init() {
	flag.StringVar(&mongoURI, "mongoURI", "", "Enter a Valid MongoURI")
	flag.StringVar(&port, "port", "8080", "HTTP Port to Start the Server")
	flag.Usage()
	flag.Parse()
	ctx := context.Background()
	fmt.Println(mongoURI)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	database := client.Database("booktree")
	booksHandler = handlers.NewBooksHandler(database.Collection("books"), ctx)
	usersHandler = handlers.NewUsersHandler(database.Collection("users"), ctx)
}

func main() {
	router := gin.Default()
	requireTokenRoutes := router.Group("/")
	router.GET("/books", booksHandler.GetBooksHandler)
	router.GET("/books/:id", booksHandler.GetBookHandler)

	requireTokenRoutes.Use(middleware.RequiresAuthToken())
	requireTokenRoutes.POST("/books", booksHandler.CreateBookHandler)
	requireTokenRoutes.DELETE("/books/:id", booksHandler.DeleteBookHandler)
	requireTokenRoutes.PATCH("/books/:id", booksHandler.UpdateBookHandler)

	router.POST("/users/signup", usersHandler.SignupUserHandler)
	router.POST("/users/login", usersHandler.LoginUserHandler)

	router.Run(fmt.Sprintf(":%s", port))
}
