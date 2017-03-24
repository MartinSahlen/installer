package main

import (
	"log"
	"net/http"
	"time"

	"github.com/MartinSahlen/installer/brew"
	"github.com/gorilla/mux"
)

func main() {

	db, err := brew.NewDB()

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	//add some more routes
	r.HandleFunc("/{id}", InstallAppHandler(db))
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
