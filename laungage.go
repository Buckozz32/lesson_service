package handlers

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/sessions"

	"path/to/project/models"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("laungage")))

// Render the language selection page
func RenderLanguagePage(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authenticated
	session, _ := sessions.Store.Get(r, "laungage")
	userID, ok := session.Values["userID"].(int)

	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get the user model
	user, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get all languages from the database
	languages, err := models.GetAllLanguages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		User      *models.User
		Languages []models.Language
	}{
		User:      user,
		Languages: languages,
	}

	tmpl, err := template.ParseFiles("views/language.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handle language selection form submission
func SelectLanguage(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get form values
	languageName := r.FormValue("language")

	// Get the language model
	language, err := models.GetLanguageByName(languageName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the selected language for the user
	session, _ := sessions.Store.Get(r, "laungage")
	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	err = models.SetLanguageForUser(userID, language.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect the user to the first lesson for the selected language
	http.Redirect(w, r, "/lesson/1", http.StatusSeeOther)
}
