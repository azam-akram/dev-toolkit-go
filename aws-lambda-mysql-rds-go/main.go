package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Student struct
type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// Initialize database connection
func initDB() {
	var err error

	// Read database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	// Ensure all environment variables are set
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbName == "" {
		log.Fatal("Missing required environment variables for database connection")
	}

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
}

// Lambda handler function
func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	initDB()
	defer db.Close()

	switch req.HTTPMethod {
	case "POST":
		return createStudent(req)
	case "GET":
		return getStudents()
	default:
		return events.APIGatewayProxyResponse{StatusCode: 405, Body: "Method Not Allowed"}, nil
	}
}

// Create student (POST)
func createStudent(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var student Student
	err := json.Unmarshal([]byte(req.Body), &student)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid Request"}, nil
	}

	fmt.Printf("Student Name: %s, Student Email = %s", student.Name, student.Email)

	query := "INSERT INTO Student (name, email, age) VALUES (?, ?, ?)"
	_, err = db.Exec(query, student.Name, student.Email, student.Age)
	if err != nil {
		log.Printf("DB Error: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to insert student"}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 201, Body: "Student created"}, nil
}

// Get all students (GET)
func getStudents() (events.APIGatewayProxyResponse, error) {
	rows, err := db.Query("SELECT id, name, email, age FROM Student")
	if err != nil {
		log.Printf("DB Error: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to fetch students"}, nil
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.Age); err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}
		students = append(students, s)
	}

	resp, _ := json.Marshal(students)
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(resp)}, nil
}

func main() {
	lambda.Start(handler)
}
