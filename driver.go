package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Student represents a student entity
type Student struct {
	ID        gocql.UUID `json:"id"`
	Firstname string     `json:"firstname"`
	Lastname  string     `json:"lastname"`
	Age       int        `json:"age"`
}

var Session *gocql.Session

func init() {

	cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster.Keyspace = "university"
	Session, _ = cluster.CreateSession()
}

func main() {
	// Create the keyspace and table if they don't exist
	err := createKeyspaceAndTable()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize routes
	r := mux.NewRouter()
	r.HandleFunc("/students", GetStudents).Methods("GET")
	r.HandleFunc("/students/{id}", GetStudentById).Methods("GET")
	r.HandleFunc("/students", CreateStudent).Methods("POST")
	r.HandleFunc("/students/{id}", UpdateStudentById).Methods("PUT")
	r.HandleFunc("/students/{id}", DeleteStudentById).Methods("DELETE")

	// Start the HTTP server
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createKeyspaceAndTable() error {
	err := Session.Query(`
		CREATE KEYSPACE IF NOT EXISTS university 
		WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1};
	`).Exec()
	if err != nil {
		return err
	}

	err = Session.Query(`
		CREATE TABLE IF NOT EXISTS university.students (
			id UUID PRIMARY KEY,
			firstname TEXT,
			lastname TEXT,
			age INT
		);
	`).Exec()
	return err
}

// GetStudents returns all students
func GetStudents(w http.ResponseWriter, r *http.Request) {
	var students []Student
	iter := Session.Query("SELECT id, firstname, lastname, age FROM students").Iter()
	for {
		var student Student
		if !iter.Scan(&student.ID, &student.Firstname, &student.Lastname, &student.Age) {
			break
		}
		students = append(students, student)
	}

	respondJSON(w, http.StatusOK, students)
}

// GetStudent returns a single student by ID
func GetStudentById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var student Student
	query := "SELECT id, firstname, lastname, age FROM students WHERE id = ?"
	if err := Session.Query(query, params["id"]).Scan(
		&student.ID, &student.Firstname, &student.Lastname, &student.Age); err != nil {
		respondError(w, http.StatusNotFound, "Student not found")
		return
	}

	respondJSON(w, http.StatusOK, student)
}

// CreateStudent creates a new student
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var newStudent Student
	if err := json.NewDecoder(r.Body).Decode(&newStudent); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	newStudent.ID = gocql.TimeUUID()
	if err := Session.Query(`
		INSERT INTO students (id, firstname, lastname, age) VALUES (?, ?, ?, ?)`,
		newStudent.ID, newStudent.Firstname, newStudent.Lastname, newStudent.Age).Exec(); err != nil {
		respondError(w, http.StatusInternalServerError, "Error inserting student")
		return
	}

	respondJSON(w, http.StatusCreated, newStudent)
}

// UpdateStudent updates a student by ID
func UpdateStudentById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedStudent Student
	if err := json.NewDecoder(r.Body).Decode(&updatedStudent); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := Session.Query(`
		UPDATE students SET firstname = ?, lastname = ?, age = ? WHERE id = ?`,
		updatedStudent.Firstname, updatedStudent.Lastname, updatedStudent.Age, params["id"]).Exec(); err != nil {
		respondError(w, http.StatusInternalServerError, "Error updating student")
		return
	}

	respondJSON(w, http.StatusOK, updatedStudent)
}

// DeleteStudent deletes a student by ID
func DeleteStudentById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := "DELETE FROM students WHERE id = ?"
	if err := Session.Query(query, params["id"]).Exec(); err != nil {
		respondError(w, http.StatusInternalServerError, "Error deleting student")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Student deleted"})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
