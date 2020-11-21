package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rajpatel/movieAPI/helper"
	localDbSearch "github.com/rajpatel/movieAPI/localDbSearch"
)

func main() {
	//Init Router
	router := mux.NewRouter()

	router.HandleFunc("/api/movies/{title}", localDbSearch.SearchMovieByTitle).Methods("GET")
	router.HandleFunc("/api/movies/update/{id}", localDbSearch.UpdateMovie).Methods("PUT")
	router.HandleFunc("/api/movies/id/{id}", localDbSearch.SearchMovieById).Methods("GET")
	router.HandleFunc("/api/movies/year/{year}", localDbSearch.SearchMoviesByYear).Methods("GET")
	router.HandleFunc("/api/movies/rating-less/{rating}", localDbSearch.SearchMoviesByRatingLower).Methods("GET")
	router.HandleFunc("/api/movies/rating-high/{rating}", localDbSearch.SearchMoviesByRatingHigher).Methods("GET")
	router.HandleFunc("/api/movies/genres/{genres}", localDbSearch.SearchMoviesByGenres).Methods("GET")

	/*
		//Possible searchby values: "year", "rating-less", "rating-high", "genres"
		router.HandleFunc("/api/movies/{searchby}/{year}", localDbSearch.GetAllMoviesByFilter).Methods("GET")
		router.HandleFunc("/api/movies/{searchby}/{rating}", localDbSearch.GetAllMoviesByFilter).Methods("GET")
		router.HandleFunc("/api/movies/{searchby}/{rating}", localDbSearch.GetAllMoviesByFilter).Methods("GET")
		router.HandleFunc("/api/movies/{searchby}/{genres}", localDbSearch.GetAllMoviesByFilter).Methods("GET")
	*/
	config := helper.GetConfiguration()
	log.Fatal(http.ListenAndServe(config.Port, router))
}
