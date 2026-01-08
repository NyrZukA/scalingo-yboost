package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	// 1. Gestion des fichiers statiques (CSS, JS, Images)
	// On dit au serveur : "Quand une URL commence par /assets/, va chercher dans le dossier 'assets'"
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// 2. Gestion de la page d'accueil
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// On va chercher le fichier dans le dossier templates
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Erreur interne : impossible de charger le template", http.StatusInternalServerError)
			fmt.Println("Erreur template:", err)
			return
		}
		tmpl.Execute(w, nil)
	})

	// 3. Lancement du serveur
	fmt.Println("🚀 Serveur prêt sur : http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erreur serveur:", err)
	}
}
