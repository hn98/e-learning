package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Assignment struct {
	Name     string `bson:"name,omitempty"`
	Filename string `bson:"filename,omitempty"`
	Deadline string `bson:"deadline,omitempty"`
}

type Batch struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Timeslot    string             `bson:"timeslot,omitempty"`
	Assignments []Assignment       `bson:"assignments,omitempty"`
}

// Instructor represents the schema for the "Instructors" collection
type Instructor struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	Username      string               `bson:"username,omitempty"`
	Password      string               `bson:"password,omitempty"`
	Name          string               `bson:"name,omitempty"`
	Email         string               `bson:"email,omitempty"`
	Qualificatiom []string             `bson:"qualifcations,omitempty"`
	Experience    []string             `bson:"experience,omitempty"`
	Fees          float64              `bson:"fees,omitempty"`
	Assignments   []string             `bson:"assignments,omitempty"`
	Batches       []primitive.ObjectID `bson:"batches,omitempty"`
}

// Student represents the schema for the "Students" collection
type Student struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty"`
	Username string               `bson:"username,omitempty"`
	Password string               `bson:"password,omitempty"`
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

type IDRequest struct {
	ID string
}

type AssignmentRequest struct {
	Name     string
	Filename string
	Deadline string
	BatchID  string
}

type DownloadRequest struct {
	Filename string
}

type LoginRequest struct {
	Username string
	Password string
}

type Token struct {
	ID       string
	Username string
	Role     string
	*jwt.StandardClaims
}

type Exception struct {
	Message string `json:"message"`
}
