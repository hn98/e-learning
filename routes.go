package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var instructorRoutes = Routes{
	Route{
		"ListBatches",
		"POST",
		"/listBatches",
		ListBatches,
	}, Route{
		"ListStudents",
		"POST",
		"/listStudents",
		ListStudents,
	},
	Route{
		"Timeslot",
		"POST",
		"/timeslot",
		Timeslot,
	},
	Route{
		"UploadFile",
		"POST",
		"/upload",
		UploadFile,
	},
}
var studentRoutes = Routes{
	Route{
		"EnrollBatch",
		"POST",
		"/enrollBatch",
		EnrollBatch,
	}, Route{
		"UnenrollBatch",
		"POST",
		"/unenrollBatch",
		UnenrollBatch,
	},
	Route{
		"StudentBatchDetails",
		"POST",
		"/studentBatchDetails",
		StudentBatchDetails,
	},
	Route{
		"BatchInfo",
		"POST",
		"/batchInfo",
		BatchInfo,
	},
	Route{
		"DownloadFile",
		"GET",
		"/download",
		DownloadFile,
	},
}
var loginRoutes = Routes{
	Route{
		"InstructorLogin",
		"POST",
		"/login/instructor",
		InstructorLogin,
	},
	Route{
		"StudentLogin",
		"POST",
		"/login/student",
		StudentLogin,
	},
}
