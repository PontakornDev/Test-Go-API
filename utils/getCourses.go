package utils

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Course struct {
	CourseID   int     `json: "courseid"`
	CourseName string  `json: "coursename"`
	Price      float64 `json: "price"`
	ImageURL   string  `json: "imageurl"`
}

func GetCourseList(Db *sql.DB) ([]Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := Db.QueryContext(ctx, `SELECT
	courseid,
	coursename,
	price,
	image_url
	FROM courseonline
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer result.Close()
	courses := make([]Course, 0)
	for result.Next() {
		var course Course
		result.Scan(&course.CourseID,
			&course.CourseName,
			&course.Price,
			&course.ImageURL)

		courses = append(courses, course)
	}
	return courses, nil
}

func GetCourseID(courseid int, Db *sql.DB) (*Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := Db.QueryRowContext(ctx, `SELECT
	courseid,
	coursename,
	price,
	image_url
	FROM courseonline
	WHERE courseid = ?`, courseid)

	course := &Course{}
	err := row.Scan(
		&course.CourseID,
		&course.CourseName,
		&course.Price,
		&course.ImageURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return course, nil
}
