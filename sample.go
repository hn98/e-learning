package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func emptyCollection(c *mongo.Collection) (int64, error) {
	deleteResult, err := c.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}

func insertSampleStudents(db *mongo.Database) (*mongo.InsertManyResult, error) {
	studentCollection := db.Collection("Students")
	emptyCollection(studentCollection)

	pass, err := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	batchID, err := primitive.ObjectIDFromHex("5ff37a95c8f63363476389f6")

	s1 := Student{
		Username: "Ayush.Sharma",
		Password: string(pass),
		Name:     "Ayush Sharma",
		Email:    "ayush.sharma@test.com",
		Std:      "Sem 1-2",
		Balance:  25000.0,
		Location: "Delhi",
		Batches:  []primitive.ObjectID{batchID},
	}
	s2 := Student{
		Username: "Rahul.Singh",
		Password: string(pass),
		Name:     "Rahul Singh",
		Email:    "rahul.singh@test.com",
		Std:      "Sem 1-2",
		Balance:  22000.0,
		Location: "Jaipur",
		Batches:  []primitive.ObjectID{batchID},
	}

	studentList := []interface{}{s1, s2}

	result, err := studentCollection.InsertMany(context.Background(), studentList)
	return result, err
}

func insertSampleInstructors(db *mongo.Database) (*mongo.InsertOneResult, error) {
	instructorCollection := db.Collection("Instructors")
	emptyCollection(instructorCollection)

	pass, err := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	batchID1, err := primitive.ObjectIDFromHex("5ff37a95c8f63363476389f6")
	batchID2, err := primitive.ObjectIDFromHex("5ff37a95c8f63363476389f7")

	batch_ids := []primitive.ObjectID {batchID1, batchID2}

	i1 := Instructor{
		Username: "Rajiv.Kumar",
		Password: string(pass),
		Name:     "Rajiv Kumar",
		Email:    "rajiv.kumar@test.com",
		Fees: 1210.0,
		Batches: batch_ids,
	}

	result, err := instructorCollection.InsertOne(context.Background(), i1)
	return result, err
}