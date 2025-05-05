package main

import (
	"encoding/json"
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

	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/add", addSentence)
	http.HandleFunc("/random", getRandomSentence)
	http.HandleFunc("/save", saveSentencesToExcel) // New route for saving sentences to Excel

	http.ListenAndServe(":8080", nil)
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Sentence App</title>
    </head>
    <body>
        <h1>Sentence App</h1>
        <form id="sentenceForm">
            <input type="text" id="sentenceInput" placeholder="Enter a sentence" required>
            <button type="button" onclick="addSentence()">Add Sentence</button>
        </form>
        <br>
        <button onclick="getRandomSentence()">Get Random Sentence</button>
        <button onclick="saveSentences()">Save Sentences to Excel</button> <!-- New button for saving sentences -->
        <p id="randomSentence"></p>

        <script>
            function addSentence() {
                const input = document.getElementById('sentenceInput').value;
                fetch('/add', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ sentence: input })
                }).then(() => {
                    alert('Sentence added!');
                    document.getElementById('sentenceInput').value = '';
                });
            }

            function getRandomSentence() {
                fetch('/random')
                    .then(response => response.json())
                    .then(data => {
                        document.getElementById('randomSentence').innerText = data.sentence || 'No sentences available.';
                    });
            }

            function saveSentences() {
                fetch('/save', {
                    method: 'POST'
                }).then(response => {
                    if (response.ok) {
                        alert('Sentences saved to Excel!');
                    } else {
                        alert('Failed to save sentences.');
                    }
                });
            }
        </script>
    </body>
    </html>
    `
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
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
