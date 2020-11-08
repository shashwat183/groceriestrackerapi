package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Printf("Starting API")
	initialisedb()
	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)
	router.HandleFunc("/", home)
	router.HandleFunc("/grocery", createGrocery).Methods("POST")
	router.HandleFunc("/grocery", getGroceries).Methods("GET")
	router.HandleFunc("/grocery/{name}", getGrocery).Methods("GET")
	router.HandleFunc("/grocery/{name}", updateGrocery).Methods("PATCH")
	router.HandleFunc("/grocery/{name}", deleteGrocery).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
