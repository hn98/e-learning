package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

func HandleRequest(w http.ResponseWriter, r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return err
	}

	return nil
}

func HandleIDRequest(w http.ResponseWriter, r *http.Request) (primitive.ObjectID, error) {
	var request IDRequest

	err := HandleRequest(w, r, &request)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	id, err := primitive.ObjectIDFromHex(request.ID)

	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": err}
		json.NewEncoder(w).Encode(resp)
		return primitive.ObjectID{}, err
	}

	return id, nil
}

func Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func getUserToken(r *http.Request) *Token {
	tk := r.Context().Value("user")
	if tk != nil {
		return tk.(*Token)
	}
	panic("User token not found")
}

func ListStudents(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequest(w, r)
	result, _ := GetStudentList(database, id)
	Respond(w, result)
}

func ListBatches(w http.ResponseWriter, r *http.Request) {
	tk := getUserToken(r)
	log.Println("Request with token ", tk)
	id, err := primitive.ObjectIDFromHex(tk.ID)
	if err != nil {
		panic(err)
	}
	result, _ := GetBatchList(database, id)
	Respond(w, result)
}

func Timeslot(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequest(w, r)
	result, _ := GetTimeslot(database, id)
	Respond(w, result)
}

func EnrollBatch(w http.ResponseWriter, r *http.Request) {
	tk := getUserToken(r)
	log.Println("Request with token ", tk)
	batchID, _ := HandleIDRequest(w, r)
	studentID, _ := primitive.ObjectIDFromHex(tk.ID)
	err := EnrollInBatch(database, studentID, batchID)

	if err == nil {
		Respond(w, "Successfully enrolled")
	}
}

func UnenrollBatch(w http.ResponseWriter, r *http.Request) {
	tk := getUserToken(r)
	log.Println("Request with token ", tk)
	batchID, _ := HandleIDRequest(w, r)
	studentID, _ := primitive.ObjectIDFromHex(tk.ID)
	err := UnenrollFromBatch(database, studentID, batchID)

	if err == nil {
		Respond(w, "Successfully Unenrolled")
	}
}

func StudentBatchDetails(w http.ResponseWriter, r *http.Request) {
	tk := getUserToken(r)
	log.Println("Request with token ", tk)
	id, _ := primitive.ObjectIDFromHex(tk.ID)
	result, _ := GetStudentBatchDetails(database, id)
	Respond(w, result)
}

func BatchInfo(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequest(w, r)
	result, _ := GetBatchInfo(database, id)
	Respond(w, result)
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	// tk := getUserToken(r)
	// log.Println("Request with token ", tk)
	// id, _ := primitive.ObjectIDFromHex(tk.ID)
	// username := tk.Username

	username := "Rajiv.Kumar"
	id, _ := primitive.ObjectIDFromHex("5ff6eefed52c0086bc76a77a")
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
		filesDB,
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

	instructorCollection := database.Collection("Instructors")

	_, _ = instructorCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.D{
			{"$addToSet", bson.D{{"assignments", handler.Filename}}},
		},
	)

	fmt.Fprintf(w, "Successfully Uploaded File\n")
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

	if err := AllotToBatch(database, batchID, assignment); err != nil {
		json.NewEncoder(w).Encode(Exception{Message: "Could not allot assignment"})
		return
	}

	Respond(w, "Assignment alloted succcesfully")
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	fileName := "testfile"
	fsFiles := filesDB.Collection("fs.files")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var results bson.M
	err := fsFiles.FindOne(ctx, bson.M{}).Decode(&results)
	if err != nil {
		log.Fatal(err)
	}
	// you can print out the result
	fmt.Println("Results:")
	fmt.Println(results)

	bucket, _ := gridfs.NewBucket(
		filesDB,
	)
	var buf bytes.Buffer
	dStream, err := bucket.DownloadToStreamByName(fileName, &buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File size to download: %v \n", dStream)

	cd := mime.FormatMediaType("attachment", map[string]string{"filename": fileName})
	w.Header().Set("Content-Disposition", cd)
	w.Header().Set("Content-Type", "application/pdf")

	len, err := buf.WriteTo(w)
	fmt.Println(len)
}
