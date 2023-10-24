package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// User represents a user in the system.
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
	Role     string
}

var db *gorm.DB
var store = sessions.NewCookieStore([]byte("your-secret-key"))

func main() {
	// Open a SQLite database.
	database, err := gorm.Open(sqlite.Open("security_app.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}
	db = database

	// Automatically create the user table.
	db.AutoMigrate(&User{})

	// Create some example users.
	createUser("admin", "adminpassword", "admin")
	createUser("user1", "user1password", "user")
	createUser("user2", "user2password", "user")

	// Set up the web server.
	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/logout", logoutHandler)
	r.HandleFunc("/register", registerHandler).Methods("POST")
	r.HandleFunc("/reset-password", resetPasswordHandler).Methods("POST")
	r.HandleFunc("/profile", profileHandler)

	http.Handle("/", r)
	fmt.Println("Listening on :8080...")
	http.ListenAndServe(":8080", nil)
}

func createUser(username, password, role string) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}

	user := User{
		Username: username,
		Password: string(hashedPassword),
		Role:     role,
	}

	result := db.Create(&user)
	if result.Error != nil {
		fmt.Println("Error creating user:", result.Error)
	} else {
		fmt.Println("User created:", username)
	}
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var user User
	result := db.Where("username = ?", username).First(&user)

	if result.Error != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, "session-name")
	session.Values["user"] = user.Username
	session.Save(r, w)

	fmt.Fprintln(w, "Login successful")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	delete(session.Values, "user")
	session.Save(r, w)
	fmt.Fprintln(w, "Logged out")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Check if the username is already taken.
	var existingUser User
	result := db.Where("username = ?", username).First(&existingUser)
	if result.Error == nil {
		http.Error(w, "Username is already taken", http.StatusConflict)
		return
	}

	createUser(username, password, "user")
	fmt.Fprintln(w, "Registration successful")
}

func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	newPassword := r.FormValue("newPassword")

	var user User
	result := db.Where("username = ?", username).First(&user)

	if result.Error != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)
	result = db.Save(&user)

	if result.Error != nil {
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Password reset successful")
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["user"].(string)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	fmt.Fprintln(w, "Welcome, "+username)
	fmt.Fprintln(w, "Role: "+user.Role)
}
