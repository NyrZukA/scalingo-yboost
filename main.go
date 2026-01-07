package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// On définit le header pour dire au navigateur qu'on envoie du HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// On écrit le contenu HTML directement dans la réponse
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html lang="fr">
		<head>
			<meta charset="UTF-8">
			<title>Ma page Go</title>
		</head>
		<body>
			<h1>Hello World</h1>
			<p>Ceci est renvoyé par un serveur en Go !</p>
		</body>
		</html>
	`)
}

func main() {
	// On associe la route "/" à notre fonction handler
	http.HandleFunc("/", handler)

	fmt.Println("Serveur lancé sur http://localhost:8080")

	// On lance le serveur sur le port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Erreur lors du lancement du serveur : %s\n", err)
	}
}
