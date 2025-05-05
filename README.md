# BoysApp

BoysApp is a simple web application that allows users to add sentences and retrieve a random sentence. The application also provides functionality to save all added sentences into an Excel file.

## Project Structure

```
BoysApp
├── main
│   ├── main.go          # Entry point of the application
│   └── utils
│       └── excel.go     # Functions to handle Excel file operations
├── go.mod               # Module definition file
├── go.sum               # Checksums for module dependencies
└── README.md            # Documentation for the project
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd BoysApp
   ```

2. **Install dependencies:**
   Ensure you have Go installed on your machine. Run the following command to install the necessary dependencies:
   ```
   go mod tidy
   ```

3. **Run the application:**
   Start the server by executing:
   ```
   go run main/main.go
   ```
   The application will be available at `http://localhost:8080`.

## Usage

- **Add a Sentence:**
  Use the input field on the main page to enter a sentence and click "Add Sentence". This will store the sentence in memory.

- **Get a Random Sentence:**
  Click the "Get Random Sentence" button to retrieve a random sentence from the stored sentences.

- **Save Sentences to Excel:**
  (To be implemented) A feature will be added to save all sentences into an Excel file.

## Dependencies

This project uses the following Go packages:
- `github.com/xuri/excelize/v2` for handling Excel file operations.

## License

This project is licensed under the MIT License.