// Some CRUD stuff.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Unstructured Database
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

var movies []*Movie

func main() {
	r := mux.NewRouter()

	movies = append(movies, &Movie{
		ID:    "1",
		Isbn:  "438227",
		Title: "Batman",
		Director: &Director{
			Firstname: "Christopher",
			Lastname:  "Nolan",
		},
	})
	movies = append(movies, &Movie{
		ID:    "2",
		Isbn:  "438228",
		Title: "Batman 2",
		Director: &Director{
			Firstname: "Christopher",
			Lastname:  "Nolan",
		},
	})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}

func getMovies(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(movies)
}

func deleteMovie(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, m := range movies {
		if m.ID == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			break
		}
	}
	json.NewEncoder(rw).Encode(movies)
}

func getMovie(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, m := range movies {
		if m.ID == params["id"] {
			json.NewEncoder(rw).Encode(m)
			break
		}
	}
}

// $postParams = @{id='4';isbn='438227';title='Batman';director=@{firstname='Christopher';lastname='Nolan'}} | ConvertTo-Json
// Invoke-WebRequest -UseBasicParsing  http://localhost:8080 -ContentType "application/json" -Method POST -Body $postParams
func createMovie(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		json.NewEncoder(rw).Encode(http.ErrBodyNotAllowed)
	}
	fmt.Println(movie)
	movie.ID = strconv.Itoa(rand.Intn(100))
	movies = append(movies, &movie)
	json.NewEncoder(rw).Encode(movie)
}

func updateMovie(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, m := range movies {
		// Not ideal when working with DBs, but for practice is OK.
		// Delete the old movie
		if m.ID == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			// Create "new/updated" movie. Same ID
			var movie Movie
			if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
				json.NewEncoder(rw).Encode(http.ErrBodyNotAllowed)
			}
			movie.ID = m.ID
			movies = append(movies, &movie)
			json.NewEncoder(rw).Encode(movie)
			break
		}
	}
}
