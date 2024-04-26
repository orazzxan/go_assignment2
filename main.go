package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123123"
	dbname   = "database"
)

type Task struct {
	ID        int
	Name      string
	Completed bool
}

func main() {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT FALSE
	)`)
	if err != nil {
		fmt.Println(err)
		return
	}

	task := Task{Name: "Task1"}
	taskID, err := createTask(db, task)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Task1 ID:", taskID)

	tasks, err := fetchTasks(db)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Tasks:")
	for _, t := range tasks {
		fmt.Printf("- %s (completed: %t)\n", t.Name, t.Completed)
	}
}

func createTask(db *sql.DB, task Task) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO tasks(name) VALUES($1) RETURNING id", task.Name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func fetchTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, name, completed FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Name, &task.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
