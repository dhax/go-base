package main

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	http.HandleFunc("/upload-audio", handleUploadAudio)
	http.ListenAndServe(":8090", nil)
}

func handleUploadAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	// Get audio file
	file, handler, err := r.FormFile("audio")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	defer file.Close()

	// Read audio data
	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	// Write audio data to file
	fileName := handler.Filename
	fileExt := filepath.Ext(fileName)
	fileName = fileName[0 : len(fileName)-len(fileExt)]
	fileName = fileName + ".webm"
	filePath := filepath.Join("upload", fileName)
	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	fmt.Printf("Successfully uploaded audio to %s\n", filePath)

	// call to whisper API
	authToken := os.Getenv("API_KEY")
	c := openai.NewClient(authToken)
	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: "./upload/blob.webm",
	}

	resp, err := c.CreateTranslation(ctx, req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return
	}
	fmt.Println(resp.Text)

	// call GPT for suggestion
	// it will return some output
	fmt.Fprintf(w, "%s", resp.Text)
}
