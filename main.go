package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

const (
	opensea_url = "https://api.opensea.io/api/v1/assets?order_direction=desc&offset=0&limit=20"
)

func handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userUUID := mux.Vars(r)["uuid"]
	offset := "0"
	limit := "20"
	order_direction := "desc"

	if x, ok := q["order_direction"]; ok {
		order_direction = x[0]
	}

	if x, ok := q["offset"]; ok {
		offset = x[0]
	}

	if x, ok := q["limit"]; ok {
		limit = x[0]
	}

	resp, err := http.Get(opensea_url + "&owner=" + userUUID + "&order_direction=" + order_direction + "&offset=" + offset + "&limit=" + limit)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("sucess ", resp)
	w.Write(body)
	w.WriteHeader(http.StatusOK)
}

// Route declaration
func router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users/{uuid}/nft", handler).Methods("GET")
	return r
}

// Initiate web server
func main() {
	log.Print("starting server...")

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	router := router()
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
