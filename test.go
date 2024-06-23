// package main

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/google/generative-ai-go/genai"
// 	"github.com/joho/godotenv"
// 	"google.golang.org/api/option"

// 	"cloud.google.com/go/firestore"
// 	firebase "firebase.google.com/go"
// )

// func main() {
// 	// Load .env file
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	ctx := context.Background()

// 	// API KEY
// 	apiKey := os.Getenv("API_KEY")
// 	if apiKey == "" {
// 		log.Fatal("API_KEY environment variable is not set in .env file")
// 	}

// 	// Initialize Firebase app
// 	app, err := firebase.NewApp(ctx, nil)
// 	if err != nil {
// 		log.Fatalf("Error initializing app: %v", err)
// 	}

// 	// Get Firestore client
// 	firestoreClient, err := app.Firestore(ctx)
// 	if err != nil {
// 		log.Fatalf("Error getting Firestore client: %v", err)
// 	}
// 	defer firestoreClient.Close()

// 	// Get Storage client
// 	storageClient, err := app.Storage(ctx)
// 	if err != nil {
// 		log.Fatalf("Error getting Storage client: %v", err)
// 	}

// 	// Upload image to Firebase Storage
// 	imageURL, err := uploadImageToStorage(ctx, storageClient, "image2.png")
// 	if err != nil {
// 		log.Fatalf("Error uploading image: %v", err)
// 	}

// 	// Fetch image data from URL
// 	imageData, err := fetchImageFromURL(imageURL)
// 	if err != nil {
// 		log.Fatalf("Error fetching image from URL: %v", err)
// 	}

// 	// Initialize Gemini client
// 	geminiClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer geminiClient.Close()

// 	// Model configuration
// 	model := geminiClient.GenerativeModel("gemini-1.5-pro")

// 	model.SetTemperature(0.9)
// 	model.SetTopK(1)

// 	model.SafetySettings = []*genai.SafetySetting{
// 		{
// 			Category:  genai.HarmCategoryHarassment,
// 			Threshold: genai.HarmBlockMediumAndAbove,
// 		},
// 		{
// 			Category:  genai.HarmCategoryHateSpeech,
// 			Threshold: genai.HarmBlockMediumAndAbove,
// 		},
// 	}

// 	// Generation type (this would come from your Flutter app)
// 	generationType := "story" // Can be "story", "poem", or "rhythm"

// 	promptText := fmt.Sprintf("You are a very creative %s writer who can see a doodle/drawing/scribble made by a child and generate a %s with good values. Remember to keep the %s simple and have a simple language and wording. The image given to you is the child's doodle/drawing/scribble.", generationType, generationType, generationType)

// 	prompt := []genai.Part{
// 		genai.ImageData("png", imageData), // Assuming PNG format, adjust if necessary
// 		genai.Text(promptText),
// 	}

// 	iter := model.GenerateContentStream(ctx, prompt...)
// 	//defer iter.Stop()

// 	var generatedText string
// 	for {
// 		resp, err := iter.Next()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		for _, candidate := range resp.Candidates {
// 			for _, part := range candidate.Content.Parts {
// 				if textPart, ok := part.(genai.Text); ok {
// 					fmt.Print(textPart) // Print without newline for continuous output
// 					generatedText += string(textPart)
// 				}
// 			}
// 		}
// 	}

// 	// Store result in Firestore
// 	err = storeInFirestore(ctx, firestoreClient, generationType, imageURL, generatedText)
// 	if err != nil {
// 		log.Printf("Error storing in Firestore: %v", err)
// 	}
// }

// func uploadImageToStorage(ctx context.Context, client *firebase.Storage, filepath string) (string, error) {
// 	bucket, err := client.Bucket("your-firebase-storage-bucket-name")
// 	if err != nil {
// 		return "", fmt.Errorf("error getting bucket: %v", err)
// 	}

// 	obj := bucket.Object(fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), filepath))

// 	// Open local file
// 	f, err := os.Open(filepath)
// 	if err != nil {
// 		return "", fmt.Errorf("os.Open: %v", err)
// 	}
// 	defer f.Close()

// 	// Upload an object with storage.Writer
// 	wc := obj.NewWriter(ctx)
// 	if _, err = io.Copy(wc, f); err != nil {
// 		return "", fmt.Errorf("io.Copy: %v", err)
// 	}
// 	if err := wc.Close(); err != nil {
// 		return "", fmt.Errorf("Writer.Close: %v", err)
// 	}

// 	// Make the object publicly accessible
// 	if err := obj.ACL().Set(ctx, "allUsers", "READER"); err != nil {
// 		return "", fmt.Errorf("ACL.Set: %v", err)
// 	}

// 	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", "your-firebase-storage-bucket-name", obj.Name), nil
// }

// func fetchImageFromURL(url string) ([]byte, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching image: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	return io.ReadAll(resp.Body)
// }

// func storeInFirestore(ctx context.Context, client *firestore.Client, generationType, imageURL, generatedText string) error {
// 	_, _, err := client.Collection("generations").Add(ctx, map[string]interface{}{
// 		"timestamp":      time.Now(),
// 		"generationType": generationType,
// 		"text":           generatedText,
// 		"imageURL":       imageURL,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("error adding document: %v", err)
// 	}

// 	log.Println("Successfully stored in Firestore")
// 	return nil
// }
