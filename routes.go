package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"ListBatches",
		"POST",
		"/instructor/listBatches",
		ListBatches,
	}, Route{
		"ListStudents",
		"POST",
		"/instructor/listStudents",
		ListStudents,
	},
	Route{
		"Timeslot",
		"POST",
		"/instructor/timeslot",
		Timeslot,
	},
	Route{
		"EnrollBatch",
		"POST",
		"/student/enrollBatch",
		EnrollBatch,
	}, Route{
		"UnenrollBatch",
		"POST",
		"/student/unenrollBatch",
		UnenrollBatch,
	},
	Route{
		"StudentBatchDetails",
		"POST",
		"/student/studentBatchDetails",
		StudentBatchDetails,
	},
	Route{
		"BatchInfo",
		"POST",
		"/student/batchInfo",
		BatchInfo,
	},
	Route{
		"Login",
		"POST",
		"/login",
		Login,
	},
}
