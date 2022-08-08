package handlers

import (
	"context"
	"net/http"
	"tosinjs/gin-booktree/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BooksHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewBooksHandler(collection *mongo.Collection, ctx context.Context) *BooksHandler {
	return &BooksHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (handler *BooksHandler) GetBooksHandler(c *gin.Context) {
	booksCursor, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	defer booksCursor.Close(handler.ctx)
	books := make([]models.Book, 0)
	for booksCursor.Next(handler.ctx) {
		var book models.Book
		err = booksCursor.Decode(&book)
		books = append(books, book)
	}
	c.JSON(http.StatusFound, books)
}

func (handler *BooksHandler) GetBookHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid book id",
		})
		return
	}
	var book models.Book
	err = handler.collection.FindOne(handler.ctx, bson.M{"id": idHex}).Decode(&book)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "item not found",
		})
		return
	}
	c.JSON(http.StatusFound, book)
}

func (handler *BooksHandler) CreateBookHandler(c *gin.Context) {
	var input struct {
		Name        string   `json:"name" validate:"required"`
		Author      string   `json:"author" validate:"required"`
		Description string   `json:"description" validate:"required"`
		Genre       []string `json:"genre" validate:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	book := models.Book{
		Id:          primitive.NewObjectID(),
		Name:        input.Name,
		Author:      input.Author,
		Description: input.Description,
		Genre:       input.Genre,
	}
	_, err = handler.collection.InsertOne(handler.ctx, book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while inserting a new recipe",
		})
		return
	}
	c.JSON(http.StatusAccepted, book)
}

func (handler *BooksHandler) DeleteBookHandler(c *gin.Context) {
	id := c.Param("id")
	idHex, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"id": idHex,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Deleted Successfully",
	})
}

func (handler *BooksHandler) UpdateBookHandler(c *gin.Context) {

}
