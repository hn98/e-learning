package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTimeslot(id primitive.ObjectID) (string, error) {
	var batch Batch

	batchCollection := dbClient.Database("learning").Collection("Batches")
	err := batchCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&batch)

	if err != nil {
		return "", err
	}

	return batch.Timeslot, nil
}

func getInstructor(id primitive.ObjectID) (Instructor, error) {
	var instructor Instructor

	instructorCollection := dbClient.Database("learning").Collection("Instructors")
	err := instructorCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&instructor)

	return instructor, err
}

func GetBatchList(id primitive.ObjectID) ([]primitive.ObjectID, error) {
	var instructor Instructor

	instructorCollection := dbClient.Database("learning").Collection("Instructors")
	err := instructorCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&instructor)

	if err != nil {
		return nil, err
	}

	return instructor.Batches, nil
}

func UpdateInstructorInfo(id primitive.ObjectID, req InstructorProfileRequest) error {
	instructorCollection := dbClient.Database("learning").Collection("Instructors")

	_, err := instructorCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{
				{"name", req.Name},
				{"email", req.Email},
				{"qualification", req.Qualification},
				{"experience", req.Experience},
				{"fees", req.Fees},
			}},
		},
	)

	return err
}

func GetStudentList(id primitive.ObjectID) ([]Student, error) {
	var studentList []Student
	studentCollection := dbClient.Database("learning").Collection("Students")
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

func UpdateStudentInfo(id primitive.ObjectID, req StudentProfileRequest) error {
	studentCollection := dbClient.Database("learning").Collection("Students")

	_, err := studentCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{
				{"name", req.Name},
				{"email", req.Email},
				{"location", req.Location},
			}},
		},
	)

	return err
}

func AllotToBatch(batchID primitive.ObjectID, assignment Assignment) error {
	batchCollection := dbClient.Database("learning").Collection("Batches")
	_, err := batchCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": batchID},
		bson.D{
			{"$addToSet", bson.D{{"assignments", assignment}}},
		},
	)
	return err
}

func EnrollInBatch(studentID primitive.ObjectID, batchID primitive.ObjectID) error {
	studentCollection := dbClient.Database("learning").Collection("Students")
	_, err := studentCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": studentID},
		bson.D{
			{"$addToSet", bson.D{{"batches", batchID}}},
		},
	)
	return err
}

func UnenrollFromBatch(studentID primitive.ObjectID, batchID primitive.ObjectID) error {
	studentCollection := dbClient.Database("learning").Collection("Students")
	_, err := studentCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": studentID},
		bson.D{
			{"$pull", bson.D{{"batches", batchID}}},
		},
	)
	return err
}

func GetStudentBatchDetails(id primitive.ObjectID) ([]Batch, error) {
	studentCollection := dbClient.Database("learning").Collection("Students")

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

func GetBatchInfo(id primitive.ObjectID) (Batch, error) {
	var batch Batch

	batchCollection := dbClient.Database("learning").Collection("Batches")
	err := batchCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&batch)

	return batch, err
}
