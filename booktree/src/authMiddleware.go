package main

import (
	"net/http"
	"strings"
	"tosinjs/go-booktree/pkg/jwt_auth"

	"github.com/julienschmidt/httprouter"
)

func (app *application) requiresAuthentication(next httprouter.Handle) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Add("vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			app.authenticationRequiredResponse(w, r)
			return
		}
		authArr := strings.Split(authorizationHeader, " ")
		if len(authArr) != 2 {
			app.invalidAuthResponse(w, r)
			return
		}
		isValid, err := jwt_auth.VerifyJWTToken(authArr[1])
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		if !isValid {
			app.invalidAuthResponse(w, r)
			return
		}
		next(w, r, ps)
	}
	return fn
}
