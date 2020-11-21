package localDbSearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rajpatel/movieAPI/gomdb"
	imdb "github.com/rajpatel/movieAPI/gomdb"
	"github.com/rajpatel/movieAPI/helper"
	"github.com/rajpatel/movieAPI/models"
	"go.mongodb.org/mongo-driver/bson"
)

var collection = helper.ConnectToDB() //calling helper class ConnectToDB function
/*
func GetAllMoviesByFilter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movies []models.Movie
	var params = mux.Vars(r)
	filter := bson.M{}

	if params["searchby"] == "year" {
		var yearRange = strings.Split(params["year"], "-") //split by -
		//var yearRange = strings.Fields(params["year"]) //split by space
		filter := bson.M{"releasedYear": yearRange[0]}
		_ = filter
		if len(yearRange) > 1 {
			if yearRange[0] >= yearRange[1] {
				filter := bson.M{"releasedYear": bson.M{"lte": yearRange[0], "gte": yearRange[1]}}
				_ = filter
			} else {
				filter := bson.M{"releasedYear": bson.M{"lte": yearRange[1], "gte": yearRange[0]}}
				_ = filter
			}
		}
	} else if params["searchby"] == "rating-less" {
		filter := bson.M{"releasedYear": bson.M{"lt": params["rating"]}}
		_ = filter
	} else if params["searchby"] == "rating-high" {
		filter := bson.M{"releasedYear": bson.M{"gt": params["rating"]}}
		_ = filter
	} else if params["searchby"] == "genres" {
		filter := bson.M{"genres": params["genres"]}
		_ = filter
	}
	//cur, err := collection.Find(context.TODO(), bson.M{})
	cur, err := collection.Find(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var movie models.Movie

		err := cur.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}

		movies = append(movies, movie)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(movies)
}
*/
func SearchMovieByTitle(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var movie models.Movie
	// we get params with mux.
	var params = mux.Vars(r)

	filter := bson.M{"title": params["title"]}
	err := collection.FindOne(context.TODO(), filter).Decode(&movie)

	//if it doesn't exist in local database
	if err != nil {
		//fmt.Println("Doesn't exist in local database")

		//helper.GetError(err, w)
		//searching in imdb database
		cl := imdb.Init()
		query := &imdb.QueryData{Title: params["title"], SearchType: gomdb.MovieSearch}
		res2, err2 := cl.MovieByTitle(query)

		if err2 != nil {
			http.Error(w, "Doesn't exist in imdb database", http.StatusBadRequest)
			return
		}

		releasedYearInt, err := strconv.Atoi(res2.Year)
		_ = err
		RatingFloat, err := strconv.ParseFloat(res2.ImdbRating, 64)
		_ = err
		var genresArray = strings.Split(res2.Genre, ", ") //split by , and space

		mv := &models.Movie{
			Title:        res2.Title,
			ReleasedYear: releasedYearInt,
			Rating:       RatingFloat,
			Id:           res2.ImdbID,
			Genres:       genresArray,
		}
		CreateMovie(mv, w, r)
		return
	}

	json.NewEncoder(w).Encode(movie) //sending the result
}

func SearchMovieById(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var movie models.Movie
	// we get params with mux.
	var params = mux.Vars(r)

	filter := bson.M{"id": params["id"]}

	err := collection.FindOne(context.TODO(), filter).Decode(&movie)

	//if it doesn't exist in local database
	if err != nil {
		http.Error(w, "No documents found", http.StatusBadRequest)
		//fmt.Println("Doesn't exist in local database")
		return
	}

	json.NewEncoder(w).Encode(movie) //sending the result
}

func SearchMoviesByYear(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")
	var movies []models.Movie
	// we get params with mux.
	var params = mux.Vars(r)
	var yearRange = strings.Split(params["year"], "- ") //split by -

	y0, err := strconv.Atoi(yearRange[0])
	_ = err

	if len(yearRange) > 1 {
		y1, err := strconv.Atoi(yearRange[1])
		_ = err
		if y0 >= y1 {
			filter := bson.M{"releasedYear": bson.M{"$lte": y0, "$gte": y1}}
			cur, err := collection.Find(context.TODO(), filter)

			if err != nil {
				http.Error(w, "No documents found", http.StatusBadRequest)
				return
			}

			defer cur.Close(context.TODO())

			for cur.Next(context.TODO()) {

				var movie models.Movie

				err := cur.Decode(&movie)
				if err != nil {
					log.Fatal(err)
				}

				movies = append(movies, movie)
			}
			if err := cur.Err(); err != nil {
				log.Fatal(err)
			}

			json.NewEncoder(w).Encode(movies)
		} else {
			filter := bson.M{"releasedYear": bson.M{"$lte": y1, "$gte": y0}}
			cur, err := collection.Find(context.TODO(), filter)

			if err != nil {
				http.Error(w, "No documents found", http.StatusBadRequest)
				return
			}

			defer cur.Close(context.TODO())

			for cur.Next(context.TODO()) {

				var movie models.Movie

				err := cur.Decode(&movie)
				if err != nil {
					log.Fatal(err)
				}

				movies = append(movies, movie)
			}
			if err := cur.Err(); err != nil {
				log.Fatal(err)
			}

			json.NewEncoder(w).Encode(movies)
		}
	} else {
		filter := bson.M{"releasedYear": y0}
		cur, err := collection.Find(context.TODO(), filter)

		if err != nil {
			http.Error(w, "No documents found", http.StatusBadRequest)
			return
		}

		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {

			var movie models.Movie

			err := cur.Decode(&movie)
			if err != nil {
				log.Fatal(err)
			}

			movies = append(movies, movie)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(movies)
	}
}

func SearchMoviesByRatingLower(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var movies []models.Movie
	// we get params with mux.
	var params = mux.Vars(r)
	rt, err := strconv.ParseFloat(params["rating"], 64)

	filter := bson.M{"rating": bson.M{"$lt": rt}}
	//if it doesn't exist in local database
	cur, err := collection.Find(context.TODO(), filter)

	if err != nil {
		http.Error(w, "No documents found", http.StatusBadRequest)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var movie models.Movie

		err := cur.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(movies)
}
func SearchMoviesByRatingHigher(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var movies []models.Movie
	// we get params with mux.
	var params = mux.Vars(r)
	rt, err := strconv.ParseFloat(params["rating"], 64)

	filter := bson.M{"rating": bson.M{"$gt": rt}}

	cur, err := collection.Find(context.TODO(), filter)

	if err != nil {
		http.Error(w, "No documents found", http.StatusBadRequest)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var movie models.Movie

		err := cur.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}

		movies = append(movies, movie)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(movies)
}
func SearchMoviesByGenres(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var movies []models.Movie
	// we get params with mux.
	var params = mux.Vars(r)

	filter := bson.M{"genres": params["genres"]}

	//arrayGenres := [1]string{params["genres"]}

	//filter := bson.M{"genres": bson.M{"$in": arrayGenres}}

	//filter := bson.M{"genres": { "$all": }

	cur, err := collection.Find(context.TODO(), filter)

	if err != nil {
		http.Error(w, "No documents found", http.StatusBadRequest)
		return
	}
	if cur == nil {
		fmt.Println("empty")
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var movie models.Movie

		err := cur.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}

		movies = append(movies, movie)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(movies)
}

func CreateMovie(movie *models.Movie, w http.ResponseWriter, r *http.Request) {
	filter := bson.M{"id": movie.Id}

	err := collection.FindOne(context.TODO(), filter).Decode(&movie)
	if err != nil {
		_, err1 := collection.InsertOne(context.TODO(), movie)

		if err1 != nil {
			http.Error(w, "Error adding in local database", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(movie) //sending the result
	}
	//if it exists in the local database, then just return it without inserting
	json.NewEncoder(w).Encode(movie) //sending the result
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//id, _ := primitive.ObjectIDFromHex(params["id"])
	imdbId := params["id"]

	var movie models.Movie

	// Creating filter with passed id
	filter := bson.M{"id": imdbId}

	err1 := json.NewDecoder(r.Body).Decode(&movie)
	if err1 != nil {
		http.Error(w, "Error in parameters", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"rating": movie.Rating,
			"genres": movie.Genres,
		},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&movie)

	if err != nil {
		http.Error(w, "Error while updating local database", http.StatusBadRequest)
		return
	}

	//movie.Id = imdbId
	json.NewEncoder(w).Encode("Updated, check the database!")
}
