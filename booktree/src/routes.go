package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) router() http.Handler {
	router := httprouter.New()

	app.bookRouter(router)
	app.userRouter(router)

	return app.logRequests(router)
}

func (app *application) bookRouter(router *httprouter.Router) {
	router.GET("/books", app.getBooksHandler)
	router.GET("/books/:id", app.getBookHandler)
	router.POST("/books", app.requiresAuthentication(app.createBookHandler))
	router.PUT("/books/:id", app.requiresAuthentication(app.updateBookHandler))
	router.DELETE("/books/:id", app.requiresAuthentication(app.deleteBookHandler))
}

func (app *application) userRouter(router *httprouter.Router) {
	router.POST("/users/login", app.loginHandler)
	router.POST("/users/signup", app.signupHandler)
}
