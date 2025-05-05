package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/fibergun/BoysApp/utils"
)

var (
	sentences []string
	mu        sync.Mutex
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Load sentences from the Excel file
	loadedSentences, err := utils.LoadSentences()
	if err != nil {
		fmt.Println("Error loading sentences from Excel:", err)
	} else {
		mu.Lock()
		sentences = append(sentences, loadedSentences...)
		mu.Unlock()
		fmt.Println("Loaded sentences from Excel:", sentences)
	}

	// Register API routes
	registerRoutes()

	// Start the server
	http.ListenAndServe(":8080", nil)
}
