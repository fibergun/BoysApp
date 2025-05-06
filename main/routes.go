package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	http.HandleFunc("/remove/", removeQuest)
	http.HandleFunc("/clear", clearAllClaims)
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
	id := len(sentences) + 1 // ID is the line number in the Excel file
	sentences = append(sentences, input.Sentence)

	// Save updated sentences to the Excel file
	if err := utils.SaveSentencesWithPlayers(sentences, nil); err != nil {
		mu.Unlock()
		http.Error(w, "Failed to save sentences to Excel", http.StatusInternalServerError)
		return
	}
	mu.Unlock()

	// Return the ID and sentence as a JSON response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       id,
		"sentence": input.Sentence,
	})
}

func getRandomSentence(w http.ResponseWriter, r *http.Request) {
	playerName := r.URL.Query().Get("player")
	if playerName == "" {
		http.Error(w, "Player name is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Find a quest without a player name
	var availableQuestIndex int = -1
	for i, sentence := range sentences {
		if !strings.Contains(sentence, " - ") { // Check if the sentence has no player name
			availableQuestIndex = i
			break
		}
	}

	if availableQuestIndex == -1 {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Sorry, there are currently no more quests.",
		})
		return
	}

	// Assign the quest to the player
	quest := sentences[availableQuestIndex]
	sentences[availableQuestIndex] = quest + " - " + playerName

	// Save updated sentences to the Excel file
	if err := utils.SaveSentencesWithPlayers(sentences, nil); err != nil {
		http.Error(w, "Failed to save sentences to Excel", http.StatusInternalServerError)
		return
	}

	// Return the quest
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       availableQuestIndex + 1,
		"sentence": quest,
	})
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

func removeQuest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/remove/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 || id > len(sentences) {
		http.Error(w, "Invalid quest ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	sentences = append(sentences[:id-1], sentences[id:]...) // Remove the quest

	// Save updated sentences to the Excel file
	if err := utils.SaveSentencesWithPlayers(sentences, nil); err != nil {
		mu.Unlock()
		http.Error(w, "Failed to save sentences to Excel", http.StatusInternalServerError)
		return
	}
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func clearAllClaims(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Remove all player names from the sentences
	for i, sentence := range sentences {
		if strings.Contains(sentence, " - ") {
			sentences[i] = strings.Split(sentence, " - ")[0] // Keep only the quest text
		}
	}

	// Debug: Print the updated sentences
	fmt.Println("Updated sentences:", sentences)

	// Save updated sentences to the Excel file
	if err := utils.SaveSentencesWithPlayers(sentences, nil); err != nil {
		fmt.Println("Error saving sentences:", err)
		http.Error(w, "Failed to save sentences to Excel", http.StatusInternalServerError)
		return
	}

	fmt.Println("All claims removed successfully.")
	w.WriteHeader(http.StatusOK)
}
