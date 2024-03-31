package twitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/twitter/apis/openai"
	"github.com/twitter/auth"
)


var (
	consumerKey = os.Getenv("CONSUMER_KEY")
	consumerSecret = os.Getenv("CONSUMER_SECRET")
	accessTokenKey = os.Getenv("ACCESS_TOKEN_KEY")
	accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
)

func Tweet(writer http.ResponseWriter, request *http.Request) {

	// Retrieve environment variables
	consumerKey = os.Getenv("CONSUMER_KEY")
	consumerSecret = os.Getenv("CONSUMER_SECRET")
	accessTokenKey = os.Getenv("ACCESS_TOKEN_KEY")
	accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")

	// Generate signature
    signature := auth.PrepareOAuthSignature(consumerKey, accessTokenKey, consumerSecret, accessTokenSecret)

	completion, err := openai.OpenAIChatCompletion("Imagine You are a satiracal influencer that focuses on investing for tech, and a software engineer. Make witty or isnightful remarks that will be short tweets about the current news of tech. choose from topics from openAI, LLMs, valuations, technologies, languages, etc. Make the tweets to the bet of your ability sharable and inspure engagement.")
	if err != nil {
		// log.Fatal("Error creating request: ", err)
		http.Error(writer, "Error generating completion from OpenAI: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("[completion]")
	fmt.Println(completion)

	// Prepare tweet data
	tweetData := map[string]string{
		"text": completion,
	}
	tweetJSON, _ := json.Marshal(tweetData)

	// Send tweet request
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.twitter.com/2/tweets", bytes.NewBuffer(tweetJSON))
	if err != nil {
		log.Fatal("Error creating request: ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", signature)

	resp, err := client.Do(req)
	r := json.NewDecoder(resp.Body)
	if err != nil {
		log.Fatal("Error sending request: ", err)
	}
	defer resp.Body.Close()

	// Print response
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Status:", r)
}
