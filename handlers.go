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

func HandleRequest(w http.ResponseWriter, r *http.Request, v interface{}) (error) {
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
		return  err
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

func HandleEnrollmentRequest(w http.ResponseWriter, r *http.Request) (primitive.ObjectID, primitive.ObjectID, error) {
	var request EnrollmentRequest

	err := HandleRequest(w, r, &request)
	if err != nil {
		return primitive.ObjectID{}, primitive.ObjectID{}, err
	}

	studentID, err := primitive.ObjectIDFromHex(request.studentID)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": err}
		json.NewEncoder(w).Encode(resp)
		return primitive.ObjectID{}, primitive.ObjectID{}, err
	}
	batchID, err := primitive.ObjectIDFromHex(request.batchID)
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
	id, _ := HandleIDRequest(w, r)
	result, _ := GetStudentList(database, id)
	Respond(w, result)
}

func ListBatches(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequest(w, r)
	result, _ := GetBatchList(database, id)
	Respond(w, result)
}

func Timeslot(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequest(w, r)
	result, _ := GetTimeslot(database, id)
	Respond(w, result)
}

func EnrollBatch(w http.ResponseWriter, r *http.Request) {
	studentID, batchID, _ := HandleEnrollmentRequest(w, r)
	err := EnrollInBatch(database, studentID, batchID)

	if err == nil {
		Respond(w, "Successfully enrolled")
	}
}

func UnenrollBatch(w http.ResponseWriter, r *http.Request) {
	studentID, batchID, _ := HandleEnrollmentRequest(w, r)
	err := UnenrollFromBatch(database, studentID, batchID)

	if err == nil {
		Respond(w, "Successfully Unenrolled")
	}
}

func StudentBatchDetails(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequest(w, r)
	result, _ := GetStudentBatchDetails(database, id)
	Respond(w, result)
}

func BatchInfo(w http.ResponseWriter, r *http.Request) {
	id, _ := HandleIDRequest(w, r)
	result, _ := GetBatchInfo(database, id)
	Respond(w, result)
}


func UploadFile(w http.ResponseWriter, r *http.Request) {
    fmt.Println("File Upload Endpoint Hit")

    // Max upload size of 10 MB files.
    r.ParseMultipartForm(10 << 20)

    file, handler, err := r.FormFile("myFile")
    if err != nil {
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
    defer file.Close()
    fmt.Printf("Uploaded File: %+v\n", handler.Filename)
    fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	fmt.Println(handler.Header.Get("Content-Type"))

    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Println(err)
	}
	bucket, err := gridfs.NewBucket(
		filesDB,
	)
	if err != nil {
		panic(err)
	}

	uploadStream, err := bucket.OpenUploadStream(
		handler.Filename, // this is the name of the file which will be saved in the database
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

    fmt.Fprintf(w, "Successfully Uploaded File\n")
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
