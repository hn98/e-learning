package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	// "errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// instructorID, err := primitive.ObjectIDFromHex("5fec325018bec55548723b54")
	batchID, err := primitive.ObjectIDFromHex("5ff37a95c8f63363476389f6")
	batchID2, err := primitive.ObjectIDFromHex("5ff37a95c8f63363476389f7")

	// Reset stduent and instructor collection
	insertSampleStudents(database)
	insertSampleInstructors(database)

	studentList, err := GetStudentList(database, batchID)
	fmt.Println(studentList)
	studentList, err = GetStudentList(database, batchID2)
	fmt.Println(studentList)
	// res, _ := insertSampleStudents(database)

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
