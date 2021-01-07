package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func GetStudentBatchDetails(db *mongo.Database, id primitive.ObjectID) ([]Batch, error) {
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

func GetBatchInfo(db *mongo.Database, id primitive.ObjectID) (Batch, error) {
	var batch Batch

	batchCollection := db.Collection("Batches")
	err := batchCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&batch)

	return batch, err
}


