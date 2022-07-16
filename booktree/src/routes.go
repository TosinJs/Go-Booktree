package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) router() http.Handler {
	router := httprouter.New()

	router.GET("/books", app.getBooksHandler)
	router.POST("/books", app.createBookHandler)
	router.GET("/books/:id", app.getBookHandler)
	router.PUT("/books/:id", app.updateBookHandler)
	router.DELETE("/books/:id", app.deleteBookHandler)

	return router
}
