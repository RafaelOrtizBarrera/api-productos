package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type producto struct {
	Marca            string `json:"marca"`
	Modelo           string `json:"modelo"`
	Precio           string `json:"precio"`
	Tipo             string `json:"tipo"`
	DescripcionCorta string `json:"descripcionCorta"`
	NombreImg        string `json:"nombreImg"`
	URLImg           string `json:"urlImg"`
}

func crearProducto(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	/* rescatando otros datos*/
	marca := r.FormValue("marca")
	precio := r.FormValue("precio")
	modelo := r.FormValue("modelo")
	descripcionCorta := r.FormValue("descripcionCorta")
	fmt.Printf("marca %s\n", marca)
	fmt.Printf("precio %s\n", precio)
	fmt.Printf("modelo %s\n", modelo)
	fmt.Printf("descripcionCorta %s\n", descripcionCorta)

	/* procesando archivos */
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("imagen")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	fileName := header.Filename
	uploadFileMongo(file, fileName)
	urlImg := uploadFileAWS(header)
	fmt.Fprintf(w, "%v", header.Header)
	fmt.Printf("File name %s\n", fileName)
	nuevaBicicleta := producto{
		Marca:            marca,
		Modelo:           modelo,
		Precio:           precio,
		DescripcionCorta: descripcionCorta,
		NombreImg:        fileName,
		URLImg:           urlImg,
	}

	//json.Unmarshal(reqBody, &nuevoProducto)
	res2B, _ := json.Marshal(&nuevaBicicleta)
	fmt.Println(string(res2B))
	client := conexionMongoDB()
	collection := client.Database("bike-store").Collection("productos")
	// Insert a single document
	insertResult, err := collection.InsertOne(context.TODO(), &nuevaBicicleta)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&nuevaBicicleta)
}

func obtenerTodosProductos(w http.ResponseWriter, r *http.Request) {
	findOptions := options.Find()
	findOptions.SetLimit(20)
	var results []*producto
	client := conexionMongoDB()
	collection := client.Database("bike-store").Collection("productos")
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Iterate through the cursor
	for cur.Next(context.TODO()) {
		var elem producto
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
