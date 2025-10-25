package main

import (
	"fmt"
	"log"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/michaelNuel/movies-api/src/handlers"
)

func main() {
	r := mux.NewRouter()

	handlers.Movies = append(handlers.Movies, handlers.Movie{ID: "1", Isbn: "438227", Title: "Movie One", Director: &handlers.Director{Firstname: "John", Lastname: "Doe"}})
	handlers.Movies = append(handlers.Movies, handlers.Movie{ID: "2", Isbn: "454555", Title: "Movie Two", Director: &handlers.Director{Firstname: "Steve", Lastname: "Smith"}})
	handlers.Movies = append(handlers.Movies, handlers.Movie{ID: "3", Isbn: "123456", Title: "Movie Three", Director: &handlers.Director{Firstname: "Mary", Lastname: "Jane"}})

	r.HandleFunc("/movies", handlers.GetMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", handlers.GetMovie).Methods("GET")
	r.HandleFunc("/movies", handlers.CreateMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", handlers.UpdateMovie).Methods("PUt")
	r.HandleFunc("/movies/{id}", handlers.DeleteMovie).Methods("DELETE")

	fmt.Println("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))

}
