package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	cors "github.com/PontakornDev/go-API/middleware"
	utils "github.com/PontakornDev/go-API/utils"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

const coursePath = "/courses"

type Course struct {
	CourseID   int     `json: "courseid"`
	CourseName string  `json: "coursename"`
	Price      float64 `json: "price"`
	ImageURL   string  `json: "imageurl"`
}

func handleCourses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		courseList, err := utils.GetCourseList(Db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(courseList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handlesCourses(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", coursePath))
	fmt.Println("urlPathSegments:", urlPathSegments)
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	courseID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	fmt.Println("courseID:", courseID)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		course, err := utils.GetCourseID(courseID, Db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if course == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleRequests(coursePath string) {
	coursesHandler := http.HandlerFunc(handlesCourses)
	http.Handle(fmt.Sprintf("%s/", coursePath), cors.CorsMiddleware(coursesHandler))
	courseHandler := http.HandlerFunc(handleCourses)
	http.Handle("/course", cors.CorsMiddleware(courseHandler))
}

func SetupDB() {
	var err error
	Db, err = sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/coursedb")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Database is Connected!!")
	}

	Db.SetConnMaxLifetime(time.Minute * 3)
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)
}

func main() {
	SetupDB()
	handleRequests(coursePath)
	log.Fatal(http.ListenAndServe(":80", nil))
}
