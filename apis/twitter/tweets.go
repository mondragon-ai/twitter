package twitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
	// signature := generateSignature(baseString, consumerSecret, accessTokenSecret)
    signature := auth.PrepareOAuthSignature(consumerKey, accessTokenKey, consumerSecret, accessTokenSecret)

	// Prepare tweet data
	tweetData := map[string]string{
		"text": "Hello World!",
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
	if err != nil {
		log.Fatal("Error sending request: ", err)
	}
	defer resp.Body.Close()

	// Print response
	fmt.Println("Response Status:", resp.Status)
}
