package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dghubble/oauth1"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	consumerKey    = os.Getenv("CONSUMER_KEY")
	consumerSecret = os.Getenv("CONSUMER_SECRET")
)


func main() {
	//Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		fmt.Println("Error loading .env file")
	}

	fmt.Println("Starting Server")

	m := mux.NewRouter()
	m.HandleFunc("/", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(200)
		fmt.Fprintf(writer, "Server is up and running")
	})
	m.HandleFunc("/tweet", tweet).Methods("POST")

	server := &http.Server{
		Handler: m,
	}
	server.Addr = ":8080"

	server.ListenAndServe()
}

func tweet(writer http.ResponseWriter, request *http.Request) {
	requestTokenURL := "https://api.twitter.com/oauth/request_token?oauth_callback=oob&x_auth_access_type=write"
	accessTokenURL := "https://api.twitter.com/oauth/access_token"
	tweetURL := "https://api.twitter.com/2/tweets"

	// Get request token
	oauthConfig := &oauth1.Config{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		CallbackURL:    "oob",
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: requestTokenURL,
			AuthorizeURL:    "",
			AccessTokenURL:  accessTokenURL,
		},
	}
	requestToken, _, err := oauthConfig.RequestToken()
	if err != nil {
		log.Fatal(err)
	}

	authURL, err := oauthConfig.AuthorizationURL(requestToken)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Please go here and authorize:", authURL)

	// Get verifier
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Paste the PIN here: ")
	verifier, _ := reader.ReadString('\n')
	verifier = strings.TrimSpace(verifier)

	// Get access token
	accessToken, accessSecret, err := oauthConfig.AccessToken(requestToken, "", verifier)
	if err != nil {
		log.Fatal(err)
	}

	// Make tweet request
	token := oauth1.NewToken(accessToken, accessSecret)
	client := oauthConfig.Client(oauth1.NoContext, token)
	req, err := http.NewRequest("POST", tweetURL, strings.NewReader(`{"text": "Hello world!"}`))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(writer, "%s\n", body)
}



// func WebhookHandler(writer http.ResponseWriter, request *http.Request) {
// 	fmt.Println("Handler called")
// 	//Read the body of the tweet
// 	body, _ := ioutil.ReadAll(request.Body)
// 	//Initialize a webhok load obhject for json decoding
// 	var load WebhookLoad
// 	err := json.Unmarshal(body, &load)
// 	if err != nil {
// 		fmt.Println("An error occured: " + err.Error())
// 	}
// 	//Check if it was a tweet_create_event and tweet was in the payload and it was not tweeted by the bot
// 	if len(load.TweetCreateEvent) < 1 || load.UserId == load.TweetCreateEvent[0].User.IdStr {
// 		return
// 	}
// 	//Send Hello world as a reply to the tweet, replies need to begin with the handles
// 	//of accounts they are replying to
// 	_, err = SendTweet("@"+load.TweetCreateEvent[0].User.Handle+" Hello World", load.TweetCreateEvent[0].IdStr)
// 	if err != nil {
// 		fmt.Println("An error occured:")
// 		fmt.Println(err.Error())
// 	} else {
// 		fmt.Println("Tweet sent successfully")
// 	}
// }

// func CrcCheck(writer http.ResponseWriter, request *http.Request) {
// 	//Set response header to json type
// 	writer.Header().Set("Content-Type", "application/json")
// 	//Get crc token in parameter
// 	token := request.URL.Query()["crc_token"]
// 	if len(token) < 1 {
// 		fmt.Fprintf(writer, "No crc_token given")
// 		return
// 	}

// 	//Encrypt and encode in base 64 then return
// 	h := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
// 	h.Write([]byte(token[0]))
// 	encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))
// 	//Generate response string map
// 	response := make(map[string]string)
// 	response["response_token"] = "sha256=" + encoded
// 	//Turn response map to json and send it to the writer
// 	responseJson, _ := json.Marshal(response)
// 	fmt.Fprintf(writer, string(responseJson))
// }