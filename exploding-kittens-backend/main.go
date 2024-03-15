package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "exploding-kittens-backend/handlers"
)

func main() {
    r := mux.NewRouter()

    // Define API endpoints
    r.HandleFunc("/api/drawCard", handlers.DrawCard).Methods("POST")
    r.HandleFunc("/api/shuffleDeck", handlers.ShuffleDeck).Methods("POST")
    r.HandleFunc("/api/defuseBomb", handlers.DefuseBomb).Methods("POST")
    r.HandleFunc("/api/leaderboard", handlers.Leaderboard).Methods("GET")
    r.HandleFunc("/api/getScore", handlers.GetScore).Methods("POST")
    r.HandleFunc("/api/updateScore", handlers.UpdateScore).Methods("POST")

    // Start the server
    log.Println("Server is running on port 5173")
    log.Fatal(http.ListenAndServe(":5173", r))
}
