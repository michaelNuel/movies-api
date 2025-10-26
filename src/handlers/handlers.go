package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/michaelNuel/movies-api/src/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Movie struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Isbn      string    `json:"isbn" gorm:"not null"`
	Title     string    `json:"title" gorm:"not null"`
	Director  Director  `json:"director" gorm:"embedded;embeddedPrefix:director_"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var db *gorm.DB

func init() {
	config.LoadConfig()
	var err error
	db, err = gorm.Open(postgres.Open(config.DB_URL), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}
	
	// Auto-migrate the schema
	err = db.AutoMigrate(&Movie{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
}

func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var movies []Movie
	if err := db.Find(&movies).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(movies)
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	
	var movie Movie
	if err := db.First(&movie, params["id"]).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(movie)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	if movie.Title == "" || movie.Isbn == "" {
		http.Error(w, "Title and ISBN are required", http.StatusBadRequest)
		return
	}
	
	if err := db.Create(&movie).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(movie)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	if movie.Title == "" || movie.Isbn == "" {
		http.Error(w, "Title and ISBN are required", http.StatusBadRequest)
		return
	}
	
	// Check if movie exists
	var existing Movie
	if err := db.First(&existing, params["id"]).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Update the movie
	movie.ID = existing.ID
	if err := db.Save(&movie).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(movie)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	
	if err := db.Delete(&Movie{}, params["id"]).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	GetMovies(w, r)
}