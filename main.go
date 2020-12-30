package main

import (
	"context"
	"fmt"
	"time"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Batch struct {
	ID		primitive.ObjectID	`bson:"_id,omitempty"`
	Name	string				`bson:"name,omitempty"`
	Timeslot	string			`bson:"timeslot,omitempty"`
	Assignments []string		`bson:"assignments,omitempty"`
}

// Instructor represents the schema for the "Instructors" collection
type Instructor struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name,omitempty"`
	Email string             `bson:"email,omitempty"`
	Qualificatiom   []string `bson:"qualifcations,omitempty"`
	Experience	[]string	 `bson:"experience,omitempty"`	
	Fees	float64			 `bson:"fees,omitempty"`
	Batches []Batch `bson:"batches,omitempty"`
}

// Student represents the schema for the "Students" collection
type Student struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name  string              `bson:"name,omitempty"`
	Email string              `bson:"email,omitempty"`
	Std   string 			  `bson:"std,omitempty"`
	Balance	float64			  `bson:"balance,omitempty"`
	Location string 		  `bson:"location,omitempty"`
	Batches []primitive.ObjectID `bson:"batches,omitempty"`
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	database := client.Database("learning")
	instructors := database.Collection("Instructors")

	i1 := Instructor{
		Name: "Rajiv Kumar",
		Email: "rajiv.kumar@test.com",
		Fees: 1200.0,
		Batches: []Batch{
			{Name: "Math 101", Timeslot: "MWF 2"},
			{Name: "Math 102", Timeslot: "TTS 3"},
		},
	}

	result, err := instructors.InsertOne(ctx, i1)

	if err != nil {
		panic(err)
	}
	fmt.Println(result.InsertedID)
	fmt.Println(i1)

}