package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

func HandleIDRequests(w http.ResponseWriter, r *http.Request) (primitive.ObjectID, error) {
	var request IDRequest
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &request); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return primitive.ObjectID{}, err
	}

	ID, err := primitive.ObjectIDFromHex(request.ID)

	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": err}
		json.NewEncoder(w).Encode(resp)
		return primitive.ObjectID{}, err
	}

	return ID, nil
}

func HandleEnrollmentRequests(w http.ResponseWriter, r *http.Request) (primitive.ObjectID, primitive.ObjectID, error) {
	var request EnrollmentRequest
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &request); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return primitive.ObjectID{}, primitive.ObjectID{}, err
	}

	studentID, err := primitive.ObjectIDFromHex(request.studentID)
	batchID, err := primitive.ObjectIDFromHex(request.batchID)
	// TODO: Handle error

	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": err}
		json.NewEncoder(w).Encode(resp)
		return primitive.ObjectID{}, primitive.ObjectID{}, err
	}

	return studentID, batchID, nil
}

func Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func ListStudents(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequests(w, r)
	result, _ := GetStudentList(database, id)
	Respond(w, result)
}

func ListBatches(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequests(w, r)
	result, _ := GetBatchList(database, id)
	Respond(w, result)
}

func Timeslot(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequests(w, r)
	result, _ := GetTimeslot(database, id)
	Respond(w, result)
}

func EnrollBatch(w http.ResponseWriter, r *http.Request) {
	studentID, batchID, _ := HandleEnrollmentRequests(w, r)
	err := EnrollInBatch(database, studentID, batchID)

	if err == nil {
		Respond(w, "Successfully enrolled")
	}
}

func UnenrollBatch(w http.ResponseWriter, r *http.Request) {
	studentID, batchID, _ := HandleEnrollmentRequests(w, r)
	err := UnenrollFromBatch(database, studentID, batchID)

	if err == nil {
		Respond(w, "Successfully Unenrolled")
	}
}

func StudentBatchDetails(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequests(w, r)
	result, _ := GetStudentBatchDetails(database, id)
	Respond(w, result)
}

func BatchInfo(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequests(w, r)
	result, _ := GetBatchInfo(database, id)
	Respond(w, result)
}
