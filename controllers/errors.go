package controllers

import (
	"errors"
	"net/http"
)

func Unknown(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	RenderHTML(w, r, "dashboard", "errors/unknown", Data{Error: errors.New("unknown")})
}

/*
func Unauthorized(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(401)
	return RenderStatic(w, r, "errors/unauthorized")
}
*/

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	RenderHTML(w, r, "dashboard", "errors/notfound", Data{})
	//http.Redirect(w, r, "/notfound", http.StatusSeeOther)
}
