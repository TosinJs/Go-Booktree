package main

import (
	"context"
	"net/http"
	"tosinjs/go-booktree/pkg/models"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *application) getBooksHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	books, err := app.bookCollection.Find(context.TODO(), bson.M{})
	gotBooks := []models.Book{}
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	for books.Next(context.TODO()) {
		var book models.Book
		books.Decode(&book)
		gotBooks = append(gotBooks, book)
	}
	err = app.writeJSON(w, http.StatusOK, gotBooks, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input struct {
		Name        string   `json:"name" validate:"required"`
		Author      string   `json:"author" validate:"required"`
		Description string   `json:"description" validate:"required"`
		Genre       []string `json:"genre" validate:"required"`
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
	book := models.Book{
		Id:          primitive.NewObjectID(),
		Name:        input.Name,
		Author:      input.Author,
		Description: input.Description,
		Genre:       input.Genre,
	}
	_, err = app.bookCollection.InsertOne(context.TODO(), book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, book, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	idHex, _ := primitive.ObjectIDFromHex(id)
	var book models.Book
	err := app.bookCollection.FindOne(context.TODO(), bson.M{"id": idHex}).Decode(&book)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.writeJSON(w, http.StatusFound, book, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// id := ps.ByName("id")
	// newData := bson.M {
	// 	"$set": bson.M{
	// 		""
	// 	}
	// }
	// _,err := app.bookCollection.FindOneAndUpdate(context.TODO(), bson.M{"id": id}, )
}

func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	idHex, _ := primitive.ObjectIDFromHex(id)
	_, err := app.bookCollection.DeleteOne(context.TODO(), bson.M{"id": idHex})
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.writeJSON(w, http.StatusOK, map[string]string{"message": "Successfully Deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
