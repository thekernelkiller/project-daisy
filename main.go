package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"

	firebase "firebase.google.com/go"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	// API KEY
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is not set in .env file")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Model configuration
	model := client.GenerativeModel("gemini-1.5-flash")

	model.SetTemperature(0.9)
	model.SetTopK(1)

	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
	}

	generationType := "rhythm" // Generation type (this would come from your Flutter app)

	pathToImage := "image1.png"

	imgData, err := os.ReadFile(pathToImage)
	if err != nil {
		log.Fatal(err)
	}

	promptText := fmt.Sprintf("You are a very creative %s writer who can see a doodle/drawing/scribble made by a child and generate a %s with good values. Remember to keep the %s simple and have a simple language and wording. The image given to you is the child's doodle/drawing/scribble.", generationType, generationType, generationType)

	prompt := []genai.Part{
		genai.ImageData("png", imgData),
		genai.Text(promptText),
	}

	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		log.Fatal(err)
	}

	var generatedText string
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			if textPart, ok := part.(genai.Text); ok {
				generatedText = string(textPart)
				log.Println(generatedText)
			}
		}
	}

	// Firebase integration
	err = storeInFirebase(ctx, generationType, imgData, generatedText)
	if err != nil {
		log.Printf("Error storing in Firebase: %v", err)
	}
}

func storeInFirebase(ctx context.Context, generationType string, imgData []byte, generatedText string) error {
	// Initialize Firebase app
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	// Get Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("error getting Firestore client: %v", err)
	}
	defer client.Close()

	// Store data in Firestore
	_, _, err = client.Collection("generations").Add(ctx, map[string]interface{}{
		"timestamp":      time.Now(),
		"generationType": generationType,
		"text":           generatedText,
		"image":          imgData, // Note: storing image directly in Firestore is not ideal for large files
	})
	if err != nil {
		return fmt.Errorf("error adding document: %v", err)
	}

	log.Println("Successfully stored in Firebase")
	return nil
}
