package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var tasks []Task

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tasks", addTask).Methods("POST")
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	http.ListenAndServe(":8080", router)

}

func addTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newTask.ID = len(tasks) + 1
	fmt.Println(newTask.ID)
	tasks = append(tasks, newTask)
	json.NewEncoder(w).Encode(newTask)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	taskID, _ := strconv.Atoi(vars["id"])
	index := -1
	for i, task := range tasks {
		if task.ID == taskID {
			index = i
			break
		}
	}
	if index == -1 {
		http.Error(w, "Задач не найдена", http.StatusNotFound)
		return
	}
	tasks = append(tasks[:index], tasks[index+1:]...)
	json.NewEncoder(w).Encode(map[string]string{"message": "Задача успешно удалена"})
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	taskID, _ := strconv.Atoi(vars["id"])
	index := -1
	for i, task := range tasks {
		if task.ID == taskID {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "Задча не найдена", http.StatusNotFound)
		return
	}

	var updateStatus struct {
		Status string `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&updateStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[index].Status = updateStatus.Status

	json.NewEncoder(w).Encode(tasks[index])
}
