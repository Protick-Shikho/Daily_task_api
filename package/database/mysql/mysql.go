package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"daily_task/Internal/application/tasks"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLTaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *MySQLTaskRepository {
	return &MySQLTaskRepository{db: db}
}

func NewMySQLConnection(dsn string) (*sql.DB, error) {
	
	// Open a new connection to the MySQL database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open database: %v", err)
	}

	// Ping the database to ensure the connection is established
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("could not connect to database: %v", err)
	}

	log.Println("Successfully connected to MySQL")
	return db, nil
}

func (m *MySQLTaskRepository) SetupDatabase() {
	if _, err := m.db.Exec("CREATE DATABASE IF NOT EXISTS Daily_task"); err != nil {
		log.Fatalf("Error creating database: %v", err)
	}

	if _, err := m.db.Exec("USE Daily_task"); err != nil {
		log.Fatalf("Error selecting database: %v", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status ENUM('pending', 'completed') NOT NULL,  -- Corrected ENUM definition
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

	if _, err := m.db.Exec(createTableSQL); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

func (r *MySQLTaskRepository) Create(task *tasks.Task) (*tasks.Task, error) {

	if task.Status == "" {
		task.Status = "pending"
	}
	
	query := "INSERT INTO tasks (title, description, status) VALUES (?, ?, ?)"
	result, err := r.db.Exec(query, task.Title, task.Description, task.Status)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	task.ID = int(id)
	return task, nil
}

func (m *MySQLTaskRepository) ShowTasks() ([]tasks.Task, error) {
	// Query to fetch all tasks from the 'tasks' table
	query := "SELECT * FROM tasks"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ShowTasks: %v", err)
	}
	defer rows.Close()

	var tasksList []tasks.Task
	for rows.Next() {
		var task tasks.Task
		var createdAtBytes []byte

		// Scan the row into the task struct
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &createdAtBytes); err != nil {
			return nil, fmt.Errorf("ShowTasks (scanning): %v", err)
		}

		dateOnly := string(createdAtBytes)[:10]
		createdAt, err := time.Parse("2006-01-02", dateOnly)
		if err != nil {
			return nil, fmt.Errorf("ShowTasks (parsing created_at): %v", err)
		}
		task.CreatedAt = createdAt

		tasksList = append(tasksList, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ShowTasks (rows error): %v", err)
	}

	return tasksList, nil
}

func (r *MySQLTaskRepository) GetTaskByID(id int) (*tasks.Task, error) {
    query := "SELECT id, title, description, status, created_at FROM tasks WHERE id = ?"
    row := r.db.QueryRow(query, id)

    var task tasks.Task
    var createdAt string
    if err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &createdAt); err != nil {
        return nil, fmt.Errorf("GetTaskByID: %v", err)
    }

    parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAt)
    if err != nil {
        return nil, fmt.Errorf("time parsing error: %v", err)
    }

    task.CreatedAt = parsedTime

    return &task, nil
}




func (r *MySQLTaskRepository) UpdateTask(id int) (*tasks.Task, error) {
    
    query := "UPDATE tasks SET status = 'completed' WHERE id = ?"
    _, err := r.db.Exec(query, id)
    if err != nil {
        return nil, fmt.Errorf("UpdateTask: %v", err)
    }

    task, err := r.GetTaskByID(id)
    if err != nil {
        return nil, fmt.Errorf("UpdateTask (fetching updated task): %v", err)
    }

    return task, nil
}


func (m *MySQLTaskRepository) DeleteTask(id int) (*tasks.Task, error) {
	task, err := m.GetTaskByID(id)
    if err != nil {
        return nil, fmt.Errorf("DeleteTask: Error fetching task before deletion: %v", err)
    }

    query := "DELETE FROM tasks WHERE id = ?"
    result, err := m.db.Exec(query, id)
    if err != nil {
        return nil, fmt.Errorf("DeleteTask: Error executing query: %v", err)
    }

	rowsAffected, err := result.RowsAffected()
    if err != nil {
        return nil, fmt.Errorf("DeleteTask: Error getting rows affected: %v", err)
    }

    if rowsAffected == 0 {
        log.Printf("DeleteTask: No task found with ID %d", id)
        return nil, fmt.Errorf("DeleteTask: No task found with ID %d", id)
    }

    log.Printf("DeleteTask: Task with ID %d deleted successfully", id)

    tasks, err := m.ShowTasks()
    if err != nil {
        return nil, fmt.Errorf("DeleteTask: Error fetching all tasks: %v", err)
    }

    // Log the remaining tasks
    log.Println("Remaining tasks after deletion:")
    for _, t := range tasks {
        log.Printf("Task ID: %d, Title: %s, Status: %s", t.ID, t.Title, t.Status)
    }


	return task, nil
}


func (m *MySQLTaskRepository) Close() error {
	if m.db != nil {
		err := m.db.Close()
		if err != nil {
			return fmt.Errorf("failed to close the database: %w", err)
		}
	}
	return nil
}
