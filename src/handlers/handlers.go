package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/michaelNuel/movies-api/src/config"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var db *sql.DB

func init() {
	config.LoadConfig()
	var err error
	db, err = sql.Open("postgres", config.DB_URL)
	if err != nil {
		panic("Failed to connect to the database:" + err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic("Database ping failed: " + err.Error())
	}
}

func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT id, isbn, title, director_firstname, director_lastname FROM movies")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		var df, dl sql.NullString // For nullable director fields
		err := rows.Scan(&m.ID, &m.Isbn, &m.Title, &df, &dl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		m.Director = &Director{Firstname: df.String, Lastname: dl.String}
		movies = append(movies, m)
	}
	json.NewEncoder(w).Encode(movies)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	params := mux.Vars(r)
	_, err := db.Exec("DELETE FROM movies WHERE id = $1", params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json.NewEncoder(w).Encode(GetMovies(w, r))
	GetMovies(w, r)
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	params := mux.Vars(r)
	var m Movie
	var df, dl sql.NullString
	err := db.QueryRow("SELECT id, isbn, title, director_firstname, director_lastname FROM movies WHERE id = $1", params["id"]).Scan(&m.ID, &m.Isbn, &m.Title, &df, &dl)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m.Director = &Director{
		Firstname: df.String,
		Lastname:  dl.String,
	}
	json.NewEncoder(w).Encode(m)
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
	err := db.QueryRow("INSERT INTO movies (isbn, title, director_firstname, director_lastname) VALUES ($1, $2, $3, $4) RETURNING id",
		movie.Isbn, movie.Title, movie.Director.Firstname, movie.Director.Lastname).Scan(&movie.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(movie)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	params := mux.Vars(r)
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid input :" +err.Error(), http.StatusBadRequest)
		return
	}
   if movie.Title == "" || movie.Isbn == "" {
	http.Error(w, "Title and ISBN are required", http.StatusBadRequest)
	return
   }

   _, err := db.Exec("UPDATE movies SET isbn = $1, title = $2, director_firstname = $3, director_lastname = $4 WHERE id = $5",
	movie.Isbn, movie.Title, movie.Director.Firstname, movie.Director.Lastname, params["id"])
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	movie.ID = params["id"]
	json.NewEncoder(w).Encode(movie)


}
