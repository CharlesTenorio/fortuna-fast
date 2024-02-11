package handler

import (
	"net/http"
)

func DefaultMethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErroHttpMsgMethodNotAllowed.Write(w)
	})
}

func DefaultNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErroHttpMsgPageNotFound.Write(w)
	})
}
