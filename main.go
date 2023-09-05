package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./feedback.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS feedback (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		answer1 INTEGER,
		answer2 INTEGER,
		answer3 INTEGER
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func main() {
	http.HandleFunc("/", formPage)
	http.HandleFunc("/submit", submitForm)
	http.HandleFunc("/results", showResults)

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func formPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func submitForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	answer1 := r.FormValue("answer1")
	answer2 := r.FormValue("answer2")
	answer3 := r.FormValue("answer3")

	_, err := db.Exec("INSERT INTO feedback (answer1, answer2, answer3) VALUES (?, ?, ?)", answer1, answer2, answer3)
	if err != nil {
		http.Error(w, "Failed to insert into database", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

func showResults(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM feedback")
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]int
	for rows.Next() {
		var id, answer1, answer2, answer3 int
		if err := rows.Scan(&id, &answer1, &answer2, &answer3); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		results = append(results, map[string]int{
			"answer1": answer1,
			"answer2": answer2,
			"answer3": answer3,
		})
	}

	t, err := template.ParseFiles("results.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, results)
}
