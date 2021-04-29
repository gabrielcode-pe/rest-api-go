package main

//Imports
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Task struct
type task struct {
	Id      int    `json:id`
	Name    string `json:name`
	Content string `json:content`
}

// slice all tasks
type allTasks = []task

// tasks
var tasks = allTasks{
	{
		Id:      1,
		Name:    "One Task",
		Content: "Some Content",
	},
}

// Hanndler functions
func homeRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to rest API with GO")
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Insert a valid task")
	}

	json.Unmarshal(reqBody, &newTask)

	newTask.Id = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskId, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Task ID invalid")
		return
	}

	for _, task := range tasks {
		if task.Id == taskId {

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}

}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskId, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Task ID invalid")
		return
	}

	for i, task := range tasks {
		if task.Id == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...)

			fmt.Fprintf(w, "Task with ID %v has been deleted!", taskId)
		}
	}

}

// Main function
func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homeRoute)

	router.HandleFunc("/tasks", getTasks).Methods("GET")

	router.HandleFunc("/tasks", createTask).Methods("POST")

	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")

	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":5000", router))
}
