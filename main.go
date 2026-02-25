package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Todo struct {
	ID      int
	Content string
	Done    bool // MySQL stocke ça en TINYINT(1), le driver le convertit en bool
}

var db *sql.DB

func main() {
	// 1. Récupération de l'URL MySQL de Scalingo
	// Scalingo fournit souvent SCALINGO_MYSQL_URL
	rawURL := os.Getenv("SCALINGO_MYSQL_URL")
	if rawURL == "" {
		// Fallback si la variable s'appelle simplement DATABASE_URL
		rawURL = os.Getenv("DATABASE_URL")
	}

	if rawURL == "" {
		log.Println("ATTENTION: Aucune URL de base de données trouvée (SCALINGO_MYSQL_URL).")
	} else {
		// 2. Conversion de l'URL (mysql://...) vers le format DSN (user:pass@tcp...)
		dsn, err := parseMySQLURL(rawURL)
		if err != nil {
			log.Fatal("Erreur parsing URL:", err)
		}

		// 3. Connexion
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Test de connexion
		if err := db.Ping(); err != nil {
			log.Fatal("Impossible de se connecter à MySQL:", err)
		}
		
		createTable()
	}

	// Servir le CSS
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/toggle", toggleHandler)
	http.HandleFunc("/delete", deleteHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Serveur démarré sur le port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Fonction utilitaire pour convertir l'URL Scalingo en DSN Go
func parseMySQLURL(rawURL string) (string, error) {
	// Enlever le préfixe mysql:// si présent car url.Parse peut être capricieux
	if !strings.HasPrefix(rawURL, "mysql://") && strings.Contains(rawURL, "@") {
		rawURL = "mysql://" + rawURL
	}
	
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	password, _ := u.User.Password()
	// Format attendu: user:password@tcp(host:port)/dbname
	dsn := fmt.Sprintf("%s:%s@tcp(%s)%s?parseTime=true", 
		u.User.Username(), 
		password, 
		u.Host, 
		u.Path,
	)
	return dsn, nil
}

func createTable() {
	// Syntaxe MySQL : AUTO_INCREMENT au lieu de SERIAL
	query := `CREATE TABLE IF NOT EXISTS todos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		content TEXT NOT NULL,
		done BOOLEAN DEFAULT FALSE
	);`
	if db != nil {
		_, err := db.Exec(query)
		if err != nil {
			log.Println("Erreur création table:", err)
		}
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
				// Le driver MySQL gère la conversion TINYINT -> bool automatiquement
				if err := rows.Scan(&t.ID, &t.Content, &t.Done); err == nil {
					todos = append(todos, t)
				}
			}
		} else {
			log.Println("Erreur Query:", err)
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, todos)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && db != nil {
		content := r.FormValue("content")
		if content != "" {
			// MySQL utilise '?' comme placeholder
			db.Exec("INSERT INTO todos (content) VALUES (?)", content)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	if db != nil {
		id := r.FormValue("id")
		// Syntaxe MySQL standard
		db.Exec("UPDATE todos SET done = !done WHERE id = ?", id)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if db != nil {
		id := r.FormValue("id")
		db.Exec("DELETE FROM todos WHERE id = ?", id)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}