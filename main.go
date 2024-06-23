package main

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
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

	pathToImage := "image2.png"

	imgData, err := os.ReadFile(pathToImage)
	if err != nil {
		log.Fatal(err)
	}

	prompt := []genai.Part{
		genai.ImageData("png", imgData),
		// genai.ImageData("jpeg", imgData2),
		genai.Text("You are a very creative story writer who can see a doodle/drawing/scribble made by a child and generate a story with good values. Remember the keep the story simple and have a simple language and wording. The image given to you is the child's doodle/drawing/scribble."),
	}

	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		log.Fatal(err)
	}

	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			if textPart, ok := part.(genai.Text); ok {
				log.Println(textPart)
			}
		}
	}
}
