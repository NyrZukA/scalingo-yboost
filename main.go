package main

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed templates static
var content embed.FS

type Todo struct {
	ID      int
	Content string
	Done    bool
}

var db *sql.DB

func main() {
	rawURL := os.Getenv("SCALINGO_MYSQL_URL")
	if rawURL == "" {
		rawURL = os.Getenv("DATABASE_URL")
	}

	if rawURL != "" {
		dsn, err := parseMySQLURL(rawURL)
		if err != nil {
			log.Fatal(err)
		}

		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}
		createTable()
	}

	staticFS, _ := fs.Sub(content, "static")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/toggle", toggleHandler)
	http.HandleFunc("/delete", deleteHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func parseMySQLURL(rawURL string) (string, error) {
	if !strings.HasPrefix(rawURL, "mysql://") && strings.Contains(rawURL, "@") {
		rawURL = "mysql://" + rawURL
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	password, _ := u.User.Password()
	return fmt.Sprintf("%s:%s@tcp(%s)%s?parseTime=true",
		u.User.Username(),
		password,
		u.Host,
		u.Path,
	), nil
}

func createTable() {
	query := `CREATE TABLE IF NOT EXISTS todos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		content TEXT NOT NULL,
		done BOOLEAN DEFAULT FALSE
	);`
	if db != nil {
		db.Exec(query)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var todos []Todo

	if db != nil {
		rows, err := db.Query("SELECT id, content, done FROM todos ORDER BY id DESC")
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var t Todo
				if err := rows.Scan(&t.ID, &t.Content, &t.Done); err == nil {
					todos = append(todos, t)
				}
			}
		}
	}

	tmpl, err := template.ParseFS(content, "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	tmpl.Execute(w, todos)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && db != nil {
		content := r.FormValue("content")
		if content != "" {
			db.Exec("INSERT INTO todos (content) VALUES (?)", content)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && db != nil {
		id := r.FormValue("id")
		db.Exec("UPDATE todos SET done = !done WHERE id = ?", id)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && db != nil {
		id := r.FormValue("id")
		db.Exec("DELETE FROM todos WHERE id = ?", id)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
