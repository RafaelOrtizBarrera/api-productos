package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func conexionMongoDB() *mongo.Client {
	var err error
	var client *mongo.Client
	uri := viper.GetString("mongo.url")
	opts := options.Client()
	opts.ApplyURI(uri)
	opts.SetMaxPoolSize(5)
	if client, err = mongo.Connect(context.Background(), opts); err != nil {
		fmt.Println(err.Error())
	}
	return client
}

func uploadFileMongo(file multipart.File, filename string) {

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Fatal(err)
	}

	conn := conexionMongoDB()
	bucketProducto := viper.GetString("mongo.bucket")
	bucket, err := gridfs.NewBucket(
		conn.Database(viper.GetString("mongo.dataBase")),
		&options.BucketOptions{Name: &bucketProducto},
	)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	uploadStream, err := bucket.OpenUploadStream(
		filename,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Printf("Write file to DB was successful. File size: %d\n", fileSize)
}
