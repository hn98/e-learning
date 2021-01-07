package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func UpdateInstructorProfile(w http.ResponseWriter, r *http.Request) {
	tk := getUserToken(r)
	log.Println("Request with token ", tk)
	id, _ := primitive.ObjectIDFromHex(tk.ID)

	var req InstructorProfileRequest

	if err := HandleRequest(w, r, &req); err != nil {
		panic(err)
	}

	if err := UpdateInstructorInfo(id, req); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "Profile update failed"})
		return
	}
	Respond(w, "Profile updated successfully")
}

func ListStudents(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequest(w, r)
	result, _ := GetStudentList(id)
	Respond(w, result)
}

func ListBatches(w http.ResponseWriter, r *http.Request) {
	tk := getUserToken(r)
	log.Println("Request with token ", tk)
	id, err := primitive.ObjectIDFromHex(tk.ID)
	if err != nil {
		panic(err)
	}
	result, _ := GetBatchList(id)
	Respond(w, result)
}

func Timeslot(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequest(w, r)
	result, _ := GetTimeslot(id)
	Respond(w, result)
}

func UpdateStudentProfile(w http.ResponseWriter, r *http.Request) {
	tk := getUserToken(r)
	log.Println("Request with token ", tk)
	id, _ := primitive.ObjectIDFromHex(tk.ID)

	var req StudentProfileRequest

	if err := HandleRequest(w, r, &req); err != nil {
		panic(err)
	}

	if err := UpdateStudentInfo(id, req); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "Profile update failed"})
		return
	}
	Respond(w, "Profile updated successfully")
}

func EnrollBatch(w http.ResponseWriter, r *http.Request) {
	tk := getUserToken(r)
	log.Println("Request with token ", tk)
	batchID, _ := HandleIDRequest(w, r)
	studentID, _ := primitive.ObjectIDFromHex(tk.ID)
	err := EnrollInBatch(studentID, batchID)

	if err == nil {
		Respond(w, "Successfully enrolled")
	}
}

func UnenrollBatch(w http.ResponseWriter, r *http.Request) {
	tk := getUserToken(r)
	log.Println("Request with token ", tk)
	batchID, _ := HandleIDRequest(w, r)
	studentID, _ := primitive.ObjectIDFromHex(tk.ID)
	err := UnenrollFromBatch(studentID, batchID)

	if err == nil {
		Respond(w, "Successfully Unenrolled")
	}
}

func StudentBatchDetails(w http.ResponseWriter, r *http.Request) {
	tk := getUserToken(r)
	log.Println("Request with token ", tk)
	id, _ := primitive.ObjectIDFromHex(tk.ID)
	result, _ := GetStudentBatchDetails(id)
	Respond(w, result)
}

func BatchInfo(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequest(w, r)
	result, _ := GetBatchInfo(id)
	Respond(w, result)
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	tk := getUserToken(r)
	log.Println("Request with token ", tk)
	id, _ := primitive.ObjectIDFromHex(tk.ID)
	username := tk.Username

	// TEST
	// username := "Rajiv.Kumar"
	// id, _ := primitive.ObjectIDFromHex("5ff6eefed52c0086bc76a77a")
	// Max upload size of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}
	defer file.Close()
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	bucket, err := gridfs.NewBucket(
		dbClient.Database("myFiles"),
	)
	if err != nil {
		panic(err)
	}

	uploadStream, err := bucket.OpenUploadStream(
		username + handler.Filename, // this is the name of the file which will be saved in the database
	)
	if err != nil {
		panic(err)
	}
	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(fileBytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Write file to DB was successful. File size: %d \n", fileSize)

	instructorCollection := dbClient.Database("learning").Collection("Instructors")

	_, _ = instructorCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.D{
			{"$addToSet", bson.D{{"assignments", handler.Filename}}},
		},
	)
}

func AllotAssignment(w http.ResponseWriter, r *http.Request) {
	tk := getUserToken(r)
	instructorID, _ := primitive.ObjectIDFromHex(tk.ID)
	log.Println("Request with token ", tk)

	var req AssignmentRequest
	err := HandleRequest(w, r, &req)
	if err != nil {
		panic(err)
	}

	batchID, _ := primitive.ObjectIDFromHex(req.BatchID)
	instructor, _ := getInstructor(instructorID)

	if !(findString(instructor.Assignments, req.Filename)) {
		json.NewEncoder(w).Encode(Exception{Message: "Can not find file with given name"})
		return
	}

	assignment := Assignment{
		Name:     req.Name,
		Filename: tk.Username + req.Filename,
		Deadline: req.Deadline,
	}

	if err := AllotToBatch(batchID, assignment); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "Could not allot assignment"})
		return
	}

	Respond(w, "Assignment alloted succcesfully")
}

func FindExamDetails(w http.ResponseWriter, r *http.Request) {
	tk := getUserToken(r)
	log.Println("Request with token ", tk)
	id, _ := primitive.ObjectIDFromHex(tk.ID)
	batches, _ := GetStudentBatchDetails(id)

	var assignments []Assignment

	for _, batch := range batches {
		assignments = append(assignments, batch.Assignments...)
	}
	Respond(w, assignments)
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	var req DownloadRequest
	if err := HandleRequest(w, r, &req); err != nil {
		panic(err)
	}

	filesDB := dbClient.Database("myFiles")
	fsFiles := filesDB.Collection("fs.files")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var results bson.M
	if err := fsFiles.FindOne(ctx, bson.M{"filename": req.Filename}).Decode(&results); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "Can not find file with given name"})
		return
	}

	log.Println("Result ", results)

	bucket, _ := gridfs.NewBucket(
		filesDB,
	)
	var buf bytes.Buffer
	_, err := bucket.DownloadToStreamByName(req.Filename, &buf)
	if err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "Download failed"})
		return
	}

	cd := mime.FormatMediaType("attachment", map[string]string{"filename": req.Filename})
	w.Header().Set("Content-Disposition", cd)
	w.Header().Set("Content-Type", "application/pdf")

	len, err := buf.WriteTo(w)
	fmt.Println(len)
}
