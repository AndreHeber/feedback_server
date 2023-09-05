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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS questions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        question TEXT,
        type TEXT
    );`)
	if err != nil {
		log.Fatalf("Failed to create questions table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS feedback (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        question_id INTEGER,
        numeric_answer INTEGER,
        text_answer TEXT,
        FOREIGN KEY(question_id) REFERENCES questions(id)
    );`)
	if err != nil {
		log.Fatalf("Failed to create feedback table: %v", err)
	}

	// Populate questions, handle errors gracefully.
	_, err = db.Exec(`INSERT OR IGNORE INTO questions (id, question, type) VALUES 
                     (1, 'How well do I listen to others?', 'numeric'), 
                     (2, 'How clear is my communication?', 'numeric'), 
                     (3, 'Any additional comments?', 'text')`)
	if err != nil {
		log.Printf("Failed to insert questions: %v", err)
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
	rows, err := db.Query("SELECT id, question, type FROM questions")
	if err != nil {
		log.Printf("Error querying questions: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var questions []map[string]interface{}
	for rows.Next() {
		var id int
		var question, qType string
		if err := rows.Scan(&id, &question, &qType); err != nil {
			log.Printf("Error scanning rows: %v", err)
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		questions = append(questions, map[string]interface{}{
			"id":       id,
			"question": question,
			"type":     qType,
		})
	}

	t, err := template.ParseFiles("form.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, questions)
}

func submitForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	for key, values := range r.PostForm {
		var questionID int
		fmt.Sscanf(key, "q%d", &questionID)
		if questionID == 0 {
			continue
		}

		value := values[0]
		row := db.QueryRow("SELECT type FROM questions WHERE id = ?", questionID)
		var qType string
		if err := row.Scan(&qType); err != nil {
			log.Printf("Error scanning question type: %v", err)
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}

		if qType == "numeric" {
			_, err := db.Exec("INSERT INTO feedback (question_id, numeric_answer) VALUES (?, ?)", questionID, value)
			if err != nil {
				log.Printf("Error inserting numeric answer: %v", err)
				http.Error(w, "Failed to insert into database", http.StatusInternalServerError)
				return
			}
		} else if qType == "text" {
			_, err := db.Exec("INSERT INTO feedback (question_id, text_answer) VALUES (?, ?)", questionID, value)
			if err != nil {
				log.Printf("Error inserting text answer: %v", err)
				http.Error(w, "Failed to insert into database", http.StatusInternalServerError)
				return
			}
		}
	}

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

func showResults(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT q.question, q.type, f.numeric_answer, f.text_answer 
                        FROM feedback f 
                        INNER JOIN questions q ON f.question_id = q.id`)
	if err != nil {
		log.Printf("Error querying results: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var question, qType string
		var numericAnswer sql.NullInt64
		var textAnswer sql.NullString
		if err := rows.Scan(&question, &qType, &numericAnswer, &textAnswer); err != nil {
			log.Printf("Error scanning results: %v", err)
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		results = append(results, map[string]interface{}{
			"question":      question,
			"type":          qType,
			"numericAnswer": numericAnswer,
			"textAnswer":    textAnswer,
		})
	}

	t, err := template.ParseFiles("results.html")
	if err !=

 nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, results)
}
