package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {
	viperEnvVariable()
	router := mux.NewRouter().StrictSlash(true)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/", homeLink)
	api.HandleFunc("/productos", crearProducto).Methods(http.MethodPost)
	api.HandleFunc("/productos", obtenerTodosProductos).Methods(http.MethodGet)
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(viper.GetString("port"), handler))
}
