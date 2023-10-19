package handlers
 
import (
    "html/template"
    "net/http"
)
 
// Render the registration page
func RenderRegistrationPage(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("views/register.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
 
// Handle registration form submission
func RegisterUser(w http.ResponseWriter, r *http.Request) {
    // Parse form data
    err := r.ParseForm()
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
 
    // Get form values
    username := r.FormValue("username")
    password := r.FormValue("password")
 
    // Validate form data
    if username == "" || password == "" {
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        return
    }
 
    // Create a new user model
    user := models.User{
        Username: username,
        Password: password,
    }
 
    // Save the new user in the database
    err = user.Save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
 
    // Redirect the user to the login page
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Render the login page
func RenderLoginPage(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("views/login.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
 
// Handle login form submission
func LoginUser(w http.ResponseWriter, r *http.Request) {
    // Parse form data
    err := r.ParseForm()
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
 
    // Get form values
    username := r.FormValue("username")
    password := r.FormValue("password")
 
    // Validate form data
    if username == "" || password == "" {
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        return
    }
 
    // Check if the user exists in the database
    user, err := models.GetUserByUsername(username)
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        return
    }
 
    // Check if the password is correct
    if user.Password != password {
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        return
    }
 
    // Set a session cookie for the user
    session, err := user.CreateSession()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    cookie := &http.Cookie{
        Name:     "session",
        Value:    session.ID,
        HttpOnly: true,
    }
    http.SetCookie(w, cookie)
 
    // Redirect the user to the language page
    http.Redirect(w, r, "/language", http.StatusSeeOther)
}
