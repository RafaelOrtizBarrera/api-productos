package main

import (
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func contieneHeader(r *http.Request, key string) bool {
	contentType := r.Header.Get(key)
	return contentType == "application/json"
}
