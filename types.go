package main

import "go.mongodb.org/mongo-driver/bson/primitive"

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

// StudentDetail represents schema for joining "Students" with "Batches" collection
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
