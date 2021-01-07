package main

import (
	"context"
	// "fmt"
	"log"
	"net/http"
	"time"
	// "errors"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database
var filesDB *mongo.Database

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	database = client.Database("learning")
	filesDB = client.Database("myFiles")

	// Reset stduent and instructor collection
	// insertSampleStudents(database)
	// insertSampleInstructors(database)

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
