package main

import (
	"encoding/json"
	"html/template"
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
	offset := "0"
	limit := "20"
	order_direction := "desc"
	owner := ""

	if x, ok := q["owner"]; ok {
		owner = x[0]
	} else {
		http.ServeFile(w, r, "./assets/index.html")
		return
	}

	if x, ok := q["order_direction"]; ok {
		order_direction = x[0]
	}

	if x, ok := q["offset"]; ok {
		offset = x[0]
	}

	if x, ok := q["limit"]; ok {
		limit = x[0]
	}

	resp, err := http.Get(opensea_url + "&owner=" + owner + "&order_direction=" + order_direction + "&offset=" + offset + "&limit=" + limit)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	m := map[string]interface{}{}
	t := template.New("nfts.html")              // Create a template.
	t, err = t.ParseFiles("./assets/nfts.html") // Parse template file.

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal([]byte(body), &m); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, m); err == nil {
		log.Println("success ", resp)
	} else {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

// Route declaration
func router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")
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
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
