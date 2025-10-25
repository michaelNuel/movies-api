package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)


type Movie struct {
	ID       string   `json:"id"`
	Isbn     string   `json:"isbn"`
	Title    string   `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}


var Movies []Movie 

func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(Movies)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	params := mux.Vars(r)
	for index, item := range Movies {
		if item.ID == params["id"] {
			Movies = append(Movies[:index], Movies[index+1:]...)
			break 
		}
	}
	json.NewEncoder(w).Encode(Movies)
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	params := mux.Vars(r)
	for _, item := range Movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	 var movie Movie 
	_ = json.NewDecoder(r.Body).Decode(&movie) 
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	Movies = append(Movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	params := mux.Vars(r)

	for index, item := range Movies {
		if item.ID == params["id"] {
			Movies = append(Movies[:index], Movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID  = params["id"]
			Movies = append(Movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

}