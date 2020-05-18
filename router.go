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
	fmt.Fprintf(w, "Bienvenidos a la version %s", viper.GetString("api.version"))
}

func main() {
	viperEnvVariable()
	log.Printf("api version %s \n", viper.GetString("api.version"))
	log.Printf("api port %s \n", viper.GetString("port"))

	router := mux.NewRouter().StrictSlash(true)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/", homeLink)
	api.HandleFunc("/productos", crearProducto).Methods(http.MethodPost)
	api.HandleFunc("/productos", obtenerTodosProductos).Methods(http.MethodGet)
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(viper.GetString("port"), handler))
}
