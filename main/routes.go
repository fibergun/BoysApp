package main

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/fibergun/BoysApp/utils"
)

func registerRoutes() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Register API routes
	http.HandleFunc("/add", addSentence)
	http.HandleFunc("/random", getRandomSentence)
	http.HandleFunc("/save", saveSentencesToExcel)
}

func addSentence(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Sentence string `json:"sentence"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	mu.Lock()
	sentences = append(sentences, input.Sentence)
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func getRandomSentence(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if len(sentences) == 0 {
		json.NewEncoder(w).Encode(map[string]string{"sentence": ""})
		return
	}

	randomIndex := rand.Intn(len(sentences))
	json.NewEncoder(w).Encode(map[string]string{"sentence": sentences[randomIndex]})
}

func saveSentencesToExcel(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if err := utils.SaveSentences(sentences); err != nil {
		http.Error(w, "Failed to save sentences", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
