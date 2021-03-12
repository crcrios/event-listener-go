package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleRequests() {
	router := mux.NewRouter()

	router.HandleFunc("/odts", getPosts).Methods("GET")

	http.ListenAndServe(":8000", router)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("{\"saludo\": \"Hola\"}")
}
