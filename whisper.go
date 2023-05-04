package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Data struct {
	ID          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Price       string `json:"Price"`
	Allergies   string `json:"Allergies"`
}

type WhisperResponse struct {
	Text string `json:"text"`
	Data []Data `json:"data"`
}

func main() {
	http.HandleFunc("/upload-audio", handleUploadAudio)
	http.ListenAndServe(":8090", nil)
}

func handleUploadAudio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

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
	// Define the URL and payload
	url := "http://127.0.0.1:7860/run/predict"
	payload := map[string]interface{}{
		"data": []string{resp.Text},
	}

	// Convert the payload to JSON format
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Create a new HTTP POST request
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}

	// Set the content type header
	request.Header.Set("Content-Type", "application/json")

	// Send the request and get the response
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	// Read the response body
	defer response.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		panic(err)
	}

	// Print the response
	//fmt.Println(result["data"])
	// it will return some output
	//fmt.Println("result-> ", result)
	var datas []Data
	for _, v := range result["data"].([]interface{}) {
		str := v.(string)
		str = strings.ReplaceAll(str, "'", "\"")
		//str, _ = strconv.Unquote(str)
		fmt.Println("str-> ", str)
		//fmt.Println("str->", str[1])
		err = json.Unmarshal([]byte(str), &datas)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println("datas-> ", datas)
	}

	//var war WhisperApiResponse
	//as := result["data"].([]string)
	//b, _ := io.ReadAll(response.Body)
	//fmt.Println(string(b))
	//b, _ := fmt.Fprintf(result["data"][0])
	//fmt.Println(as[0])

	fresp := WhisperResponse{
		Text: resp.Text,
		Data: datas,
	}
	// Serialize the response as JSON
	jsonResponse, err := json.Marshal(fresp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(fresp)
	// Set the Content-Type header to indicate that the response is in JSON format
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the HTTP response body
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}
}
