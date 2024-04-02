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

	completion, err := openai.OpenAIChatCompletion("##Twitter Bot Prompt\n\n### Description: This Twitter bot generates tech and investment-oriented tweets with a focus on wit and satire. The tweets should be insightful yet predominantly parody or witty. \n\n### Bio:Tech enthusiast whos just as likely to be lost in a foreign city as debugging lines of code. Foodie pro navigating the world one misadventure at a time 🧳🏋🏼🍲\n\n### Prompt: Generate a tweet related to technology or investment, infused with wit or satire, within 160 characters. \n\nTopic: Artificial Intelligence replacing devs. \n\n Ommit hashtags and do not use quotations around the tweet generatred. \n\nHere are examples of popular tweets you can use to get inspiration. do not include any of the users:  Tech and marketing co founders are the deadliest combination for any startup. - @iuditg \n If your family doesnt think youre crazy, are you even a startup founder? - @dagorenouf \nIm a Pull Stack Developer.I just pull things off the Internet and put it into my code. - @TheJackForge \nDebugging is like an onion. There are multiple layers to it, and the more you peel them back, the more likely youre going to start crying at inappropriate times. - @iamdevloper\nif youre not happy single, you wont be happy in a relationship.true happiness comes from closing 100 chrome tabs after solving an obscure programming bug, not from someone else -@cszhu \n No meetings and uninterrupted coding is the best thing that can happen to a developer. -@ vittoStack \n Programming is 60% thinking, 5% coding, and 45% yelling the F*** word at your laptop.  -@ vittoStack\n Every startup has acquisition as a backup plan, except for the ones that win. -@naval \n A car with a driver is the new horse and buggy. -@naval \nThe narrative for this cycle is “this is the last cycle.” -@naval.")
	if err != nil {
		// log.Fatal("Error creating request: ", err)
		http.Error(writer, "Error generating completion from OpenAI: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare tweet data
	tweetData := map[string]string{
		"text": completion
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
