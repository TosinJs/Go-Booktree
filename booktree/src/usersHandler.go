package main

import (
	"context"
	"net/http"
	"tosinjs/go-booktree/pkg/jwt_auth"
	"tosinjs/go-booktree/pkg/models"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type resp struct {
	Username string
	Token    string
}

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := models.User{
		Id:       primitive.NewObjectID(),
		Username: input.Username,
		Password: input.Password,
	}
	var mockUser models.User
	err = app.userCollection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&mockUser)
	if err != nil && err.Error() != "mongo: no documents in result" {
		app.serverErrorResponse(w, r, err)
		return
	}
	if mockUser.Username == user.Username {
		app.duplicateUsernameErrorResponse(w, r)
		return
	}
	token, err := jwt_auth.GenerateJWTToken(10, user.Username)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	_, err = app.userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, resp{Username: user.Username, Token: token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := models.User{
		Id:       primitive.NewObjectID(),
		Username: input.Username,
		Password: input.Password,
	}
	var mockUser models.User
	err = app.userCollection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&mockUser)
	if err != nil && err.Error() != "mongo: no documents in result" {
		app.serverErrorResponse(w, r, err)
		return
	}
	if err != nil && err.Error() == "mongo: no documents in result" {
		app.invalidUserErrorResponse(w, r)
		return
	}
	if mockUser.Password != user.Password {
		app.invalidUserErrorResponse(w, r)
		return
	}
	token, err := jwt_auth.GenerateJWTToken(10, user.Username)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, resp{Username: user.Username, Token: token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
