package main

import (
	"context"
	"fmt"
	// "log"
	"time"
	// "errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Batch struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Timeslot    string             `bson:"timeslot,omitempty"`
	Assignments []string           `bson:"assignments,omitempty"`
}

// Instructor represents the schema for the "Instructors" collection
type Instructor struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	Name          string               `bson:"name,omitempty"`
	Email         string               `bson:"email,omitempty"`
	Qualificatiom []string             `bson:"qualifcations,omitempty"`
	Experience    []string             `bson:"experience,omitempty"`
	Fees          float64              `bson:"fees,omitempty"`
	Batches       []primitive.ObjectID `bson:"batches,omitempty"`
}

// Student represents the schema for the "Students" collection
type Student struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty"`
	Name     string               `bson:"name,omitempty"`
	Email    string               `bson:"email,omitempty"`
	Std      string               `bson:"std,omitempty"`
	Balance  float64              `bson:"balance,omitempty"`
	Location string               `bson:"location,omitempty"`
	Batches  []primitive.ObjectID `bson:"batches,omitempty"`
}

type StudentDetail struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	Name         string               `bson:"name,omitempty"`
	Email        string               `bson:"email,omitempty"`
	Std          string               `bson:"std,omitempty"`
	Balance      float64              `bson:"balance,omitempty"`
	Location     string               `bson:"location,omitempty"`
	Batches      []primitive.ObjectID `bson:"batches,omitempty"`
	BatchDetails []Batch              `bson:"batch_details,omitempty"`
}

func GetTimeslot(db *mongo.Database, id primitive.ObjectID) (string, error) {
	var batch Batch

	batchCollection := db.Collection("Batches")
	err := batchCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&batch)

	if err != nil {
		return "", err
	}

	return batch.Timeslot, nil
}

func GetBatchList(db *mongo.Database, id primitive.ObjectID) ([]primitive.ObjectID, error) {
	var instructor Instructor

	instructorCollection := db.Collection("Instructors")
	err := instructorCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&instructor)

	if err != nil {
		return nil, err
	}

	return instructor.Batches, nil
}

func UpdateInstructorInfo(db *mongo.Database, instructor Instructor) error {
	instructorCollection := db.Collection("Instructors")

	_, err := instructorCollection.ReplaceOne(
		context.Background(),
		bson.M{"_id": instructor.ID},
		instructor,
	)

	return err
}

func GetStudentList(db *mongo.Database, id primitive.ObjectID) ([]Student, error) {
	var studentList []Student
	studentCollection := db.Collection("Students")
	cursor, err := studentCollection.Find(context.Background(), bson.M{"batches": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &studentList); err != nil {
		return nil, err
	}

	return studentList, nil
}

func EnrollInBatch(db *mongo.Database, studentID primitive.ObjectID, batchID primitive.ObjectID) error {
	studentCollection := db.Collection("Students")
	_, err := studentCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": studentID},
		bson.D{
			{"$addToSet", bson.D{{"batches", batchID}}},
		},
	)
	return err
}

func UnenrollFromBatch(db *mongo.Database, studentID primitive.ObjectID, batchID primitive.ObjectID) error {
	studentCollection := db.Collection("Students")
	_, err := studentCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": studentID},
		bson.D{
			{"$pull", bson.D{{"batches", batchID}}},
		},
	)
	return err
}

func GetBatchDetails(db *mongo.Database, id primitive.ObjectID) ([]Batch, error) {
	studentCollection := db.Collection("Students")

	matchStage := bson.D{{"$match", bson.D{{"_id", id}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "Batches"}, {"localField", "batches"}, {"foreignField", "_id"}, {"as", "batch_details"}}}}

	cursor, err := studentCollection.Aggregate(context.Background(), mongo.Pipeline{matchStage, lookupStage})
	if err != nil {
		return nil, err
	}
	var showsLoaded []StudentDetail
	if err = cursor.All(context.Background(), &showsLoaded); err != nil {
		return nil, err
	}
	// TODO check len
	return showsLoaded[0].BatchDetails, nil
}

func emptyCollection(c *mongo.Collection) (int64, error) {
	deleteResult, err := c.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}

func insertSampleStudents(db *mongo.Database) (*mongo.InsertManyResult, error) {
	studentCollection := db.Collection("Students")

	batchID, err := primitive.ObjectIDFromHex("5ff37a95c8f63363476389f6")
	s1 := Student{
		Name:     "Ayush Sharma",
		Email:    "ayush.sharma@test.com",
		Std:      "Sem 1-2",
		Balance:  25000.0,
		Location: "Delhi",
		Batches:  []primitive.ObjectID{batchID},
	}
	s2 := Student{
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

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	database := client.Database("learning")

	instructorID, err := primitive.ObjectIDFromHex("5fec325018bec55548723b54")
	batchID, err := primitive.ObjectIDFromHex("5ff37a95c8f63363476389f6")
	batchID2, err := primitive.ObjectIDFromHex("5ff37a95c8f63363476389f7")
	studentID, err := primitive.ObjectIDFromHex("5ff418a8341a9f45bf4dfbe7")

	fmt.Println(UnenrollFromBatch(database, studentID, batchID))
	fmt.Println(EnrollInBatch(database, studentID, batchID))

	batchDetails, err := GetBatchDetails(database, studentID)
	fmt.Println(batchDetails)

	studentList, err := GetStudentList(database, batchID)
	fmt.Println(studentList)
	studentList, err = GetStudentList(database, batchID2)
	fmt.Println(studentList)
	// res, _ := insertSampleStudents(database)

	batches, _ := GetBatchList(database, instructorID)
	fmt.Println(batches)
}
