package main

import (
    "fmt"
    "log"
    "net/http"
    "daily_task/Internal/application/tasks"
	"daily_task/package/database"
    "github.com/gorilla/mux"
)

func main() {

	dsn := "root:123@tcp(localhost:3306)/"
    db, err := database.NewMySQLConnection(dsn)


    if err != nil {
        log.Fatal(err)
    }

	defer db.Close()

	

    taskRepo := database.NewMySQLTaskRepository(db)
	taskRepo.SetupDatabase()
    taskService := tasks.NewTaskService(taskRepo)
    taskHandler := tasks.NewTaskHandler(taskService)


    r := mux.NewRouter()
    r.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks", taskHandler.ShowTasks).Methods("GET") 
	r.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")
    
    
    

    http.Handle("/", r)
    fmt.Println("Server is running, on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

