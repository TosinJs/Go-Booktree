package handlers

import (
	"context"
	"net/http"
	"tosinjs/gin-booktree/pkg/auth"
	"tosinjs/gin-booktree/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUsersHandler(collection *mongo.Collection, ctx context.Context) *UsersHandler {
	return &UsersHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (handler *UsersHandler) SignupUserHandler(c *gin.Context) {
	var input struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request parameters",
		})
		return
	}
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request parameters",
		})
		return
	}
	user := models.User{
		Id:       primitive.NewObjectID(),
		Username: input.Username,
		Password: input.Password,
	}
	var mockUser models.User
	err = handler.collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&mockUser)
	if err != nil && err.Error() != "mongo: no documents in result" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	if mockUser.Username == user.Username {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "duplicate username",
		})
		return
	}

	token, err := auth.GenerateJWTToken(10, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	_, err = handler.collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	c.JSON(http.StatusAccepted, token)
}

func (handler *UsersHandler) LoginUserHandler(c *gin.Context) {
	var input struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request parameters",
		})
		return
	}
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request parameters",
		})
		return
	}
	user := models.User{
		Username: input.Username,
		Password: input.Password,
	}
	var mockUser models.User
	err = handler.collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&mockUser)
	if err != nil && err.Error() != "mongo: no documents in result" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	if err != nil && err.Error() == "mongo: no documents in result" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "invalid user credentials",
		})
		return
	}
	if mockUser.Password != user.Password {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "invalid user credentials",
		})
		return
	}
	token, err := auth.GenerateJWTToken(10, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	c.JSON(http.StatusAccepted, token)
}
