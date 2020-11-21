imdbMovieAPI in Golang

type Movie struct { Title string
ReleasedYear int
Rating float64 Id string
Genres []string }

API:
1. Find movie by title by exact value that’s passed in the API. 
Notes: If there is no match in local database, it will use imdb-api package for the search. If that returns result(s), then it will store the result in database and return first value.
2. Allows updates to genres and ratings of the movie. 
Implemented following search APIs against local database only:
3. Search by Id
4. Search by particular year or year range
5. Search movies with rating higher or lower than passed in value.
6. Search movies with passed in genres value

NOTE: Change to your API key for IMDB in gomdb.go file.

API Routes examples(Postman): 
id: localhost:80/api/movies/id/tt0848228 
rating less than: localhost:80/api/movies/rating-less/7.4 
rating higher than: localhost:80/api/movies/rating-high/8 
year: localhost:80/api/movies/year/2012 
year range: localhost:80/api/movies/year/2012-2017 
genres: localhost:80/api/movies/genres/Action 
Update (PUT): 
  localhost:80/api/movies/update/tt0848228 
  body: { "rating": 5.7, "genres": [ "Crime", "Drama", "Mystery", "Thriller"] }

Contact me at patel1gv@uwindsor.ca, if you have any questions.

Thank you!
