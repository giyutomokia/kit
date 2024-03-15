package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

// Define a global Redis client
var rdb *redis.Client

// Initialize the Redis client
func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})
}

// DrawCard handles drawing a card from the deck
func DrawCard(w http.ResponseWriter, r *http.Request) {
	// Implement logic to draw a card
	deck, err := rdb.LRange(r.Context(), "deck", 0, -1).Result()
	if err != nil {
		http.Error(w, "Failed to retrieve deck from Redis", http.StatusInternalServerError)
		return
	}

	if len(deck) == 0 {
		http.Error(w, "Deck is empty", http.StatusInternalServerError)
		return
	}

	// Draw a random card from the deck
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(deck))
	drawnCard := deck[index]

	// Remove the drawn card from the deck
	_, err = rdb.LRem(r.Context(), "deck", 1, drawnCard).Result()
	if err != nil {
		http.Error(w, "Failed to remove drawn card from deck", http.StatusInternalServerError)
		return
	}

	// Return the drawn card as JSON response
	resp := map[string]string{"drawnCard": drawnCard}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// ShuffleDeck handles shuffling the deck
func ShuffleDeck(w http.ResponseWriter, r *http.Request) {
	// Implement logic to shuffle the deck
	deck := []string{"😼", "🙅‍♂️", "🔀", "💣"}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	// Store the shuffled deck in Redis
	_, err := rdb.Del(r.Context(), "deck").Result()
	if err != nil {
		http.Error(w, "Failed to clear deck in Redis", http.StatusInternalServerError)
		return
	}

	for _, card := range deck {
		_, err := rdb.LPush(r.Context(), "deck", card).Result()
		if err != nil {
			http.Error(w, "Failed to push card to deck in Redis", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// DefuseBomb handles defusing the bomb
func DefuseBomb(w http.ResponseWriter, r *http.Request) {
	// Implement logic to defuse the bomb
	// In this example, we're simulating defusing the bomb by returning a success message
	resp := map[string]string{"message": "Bomb defused successfully!"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// Leaderboard handles getting the leaderboard
func Leaderboard(w http.ResponseWriter, r *http.Request) {
	// Implement logic to get the leaderboard
	// In this example, we're returning a dummy leaderboard
	leaderboard := []map[string]interface{}{
		{"username": "user1", "score": 10},
		{"username": "user2", "score": 8},
		{"username": "user3", "score": 5},
	}

	jsonResp, err := json.Marshal(leaderboard)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// GetScore handles getting the score for a user
func GetScore(w http.ResponseWriter, r *http.Request) {
	// Implement logic to get the score for a user
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	score, err := rdb.Get(r.Context(), fmt.Sprintf("score:%s", username)).Int()
	if err != nil {
		if err == redis.Nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve score", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{"username": username, "score": score}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// UpdateScore handles updating the score for a user
func UpdateScore(w http.ResponseWriter, r *http.Request) {
	// Implement logic to update the score for a user
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Increment the user's score by 1
	_, err := rdb.Incr(r.Context(), fmt.Sprintf("score:%s", username)).Result()
	if err != nil {
		http.Error(w, "Failed to update score", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"message": "Score updated successfully"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
