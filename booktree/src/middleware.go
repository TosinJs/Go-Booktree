package main

import "net/http"

func (app *application) logRequests(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		app.logger.PrintInfo(r.Method, map[string]string{
			"route": r.URL.String(),
		})
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
