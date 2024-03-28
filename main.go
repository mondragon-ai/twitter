package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
    consumerKey    = os.Getenv("CONSUMER_KEY")
    consumerSecret = os.Getenv("CONSUMER_SECRET")
    accessTokenKey = os.Getenv("ACCESS_TOKEN_KEY")
    accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
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

	//Start Server
	server := &http.Server{
		Handler: m,
	}
	server.Addr = ":8080"

	server.ListenAndServe()
}

// func generateNonce() string {
//     const nonceLength = 32
//     encodedLength := base64.StdEncoding.EncodedLen(nonceLength)
//     b := make([]byte, encodedLength)
//     _, err := rand.Read(b)
//     if err != nil {
//         fmt.Println("Error generating nonce:", err)
//         return ""
//     }
//     return base64.StdEncoding.EncodeToString(b)[:nonceLength]
// }

func generateNonce() string {
	const nonceLength = 32
	b := make([]byte, nonceLength)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error generating nonce:", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

// func tweet(writer http.ResponseWriter, request *http.Request) {
// 	// Generate oauth_nonce
//     oauthNonce := generateNonce()

//     // Generate oauth_timestamp
//     oauthTimestamp := strconv.FormatInt(time.Now().Unix(), 10)

//     // Twitter API endpoint
//     url := "https://api.twitter.com/2/tweets"

//     // Construct OAuth parameters
//     parameters := map[string]string{
//         "oauth_consumer_key":     consumerKey,
//         "oauth_nonce":            oauthNonce,
//         "oauth_signature_method": "HMAC-SHA1",
//         "oauth_timestamp":        oauthTimestamp,
//         "oauth_token":            accessTokenKey,
//         "oauth_version":          "1.0",
//     }

//     // Add query parameters to the parameters map
//     queryValues := request.URL.Query()
//     for key, values := range queryValues {
//         parameters[key] = values[0] // Only consider the first value if there are multiple
//     }

//     // Generate signature base string
//     signatureBaseString := fmt.Sprintf("%s&%s&", "POST", urlEncode(url))
//     var parameterStrings []string
//     for key, value := range parameters {
//         parameterStrings = append(parameterStrings, fmt.Sprintf("%s=%s", urlEncode(key), urlEncode(value)))
//     }
//     signatureBaseString += urlEncode(strings.Join(parameterStrings, "&"))

//     // Generate signing key
//     signingKey := urlEncode(consumerSecret) + "&" + urlEncode(accessTokenSecret)

//     // Generate signature
//     signature := generateSignature(signatureBaseString, signingKey)

//     // Construct OAuth header
//     oauthHeader := fmt.Sprintf(`OAuth oauth_consumer_key="%s", oauth_nonce="%s", oauth_signature="%s", oauth_signature_method="HMAC-SHA1", oauth_timestamp="%s", oauth_token="%s", oauth_version="1.0"`,
//         consumerKey, oauthNonce, urlEncode(signature), oauthTimestamp, accessTokenKey)

//     tweetText := "TEST BODY"

//     // Construct the request body
//     requestBody := []byte(fmt.Sprintf(`{"text": "%s"}`, tweetText))

//     // Create HTTP request
//     req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
//     if err != nil {
//         fmt.Println("Error creating request:", err)
//         return
//     }
//     // Set content type header
//     req.Header.Set("Content-Type", "application/json")

//     // Set authorization header
//     req.Header.Set("Authorization", oauthHeader)

//     // Make HTTP request
//     client := &http.Client{}
//     resp, err := client.Do(req)
//     if err != nil {
//         fmt.Println("Error making request:", err)
//         return
//     }
//     defer resp.Body.Close()

//     // Print response status
//     fmt.Println("Response Status:", resp.Status)
// }

// func generateSignature(baseString, signingKey string) string {
// 	hash := hmac.New(sha1.New, []byte(signingKey))
// 	hash.Write([]byte(baseString))
// 	signature := hash.Sum(nil)
// 	return base64.StdEncoding.EncodeToString(signature)
// }

// func urlEncode(s string) string {
// 	return strings.ReplaceAll(base64.URLEncoding.EncodeToString([]byte(s)), "=", "%3D")
// }

func tweet(writer http.ResponseWriter, request *http.Request) {
	// Construct OAuth parameters
	oauthNonce := generateNonce()
	oauthTimestamp := strconv.FormatInt(time.Now().Unix(), 10)

	oauthParams := map[string]string{
		"oauth_consumer_key":     consumerKey,
		"oauth_token":            accessTokenKey,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        oauthTimestamp,
		"oauth_nonce":            oauthNonce,
		"oauth_version":          "1.0",
	}

	// Generate base string
	baseString := generateBaseString(request.Method, request.URL.String(), oauthParams)

	// Generate signature
	signature := generateSignature(baseString, consumerSecret, accessTokenSecret)

	// Add signature to OAuth parameters
	oauthParams["oauth_signature"] = signature

	fmt.Println(oauthParams)

	// Construct Authorization header
	authHeader := "OAuth "
	for key, value := range oauthParams {
		authHeader += fmt.Sprintf("%s=\"%s\", ", key, value)
	}
	authHeader = strings.TrimSuffix(authHeader, ", ")
	fmt.Println(authHeader)

	// Prepare tweet data
	tweetData := map[string]string{
		"text": "Hello World!",
	}
	fmt.Println(tweetData)
	tweetJSON, _ := json.Marshal(tweetData)
	fmt.Println(tweetJSON)

	// Send tweet request
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.twitter.com/2/tweets", bytes.NewBuffer(tweetJSON))
	if err != nil {
		log.Fatal("Error creating request: ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request: ", err)
	}
	defer resp.Body.Close()

	// Print response
	fmt.Println("Response Status:", resp.Status)
}

// func generateNonce() string {
// 	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	nonce := make([]byte, 32)
// 	for i := range nonce {
// 		nonce[i] = charset[rand.Intn(len(charset))]
// 	}
// 	return string(nonce)
// }

func generateBaseString(method, urlStr string, params map[string]string) string {
	var paramString string
	for key, value := range params {
		paramString += key + "=" + value + "&"
	}
	paramString = strings.TrimSuffix(paramString, "&")

	return fmt.Sprintf("%s&%s&%s", method, urlStr, strings.ReplaceAll(paramString, "=", "%3D"))
}

func generateSignature(baseString, consumerSecret, tokenSecret string) string {
	signingKey := consumerSecret + "&" + tokenSecret
	h := hmac.New(sha1.New, []byte(signingKey))
	h.Write([]byte(baseString))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature
}
// func tweet(writer http.ResponseWriter, request *http.Request) {
//     // Prepare the request body
//     requestBody := strings.NewReader(`{"text": "Hello World!"}`)

//     // Create OAuth1 configuration
//     config := oauth1.NewConfig(consumerKey, consumerSecret)
//     token := oauth1.NewToken(accessTokenKey, accessTokenSecret)
//     httpClient := config.Client(oauth1.NoContext, token)

//     // Prepare the HTTP request
//     req, err := http.NewRequest("POST", "https://api.twitter.com/2/tweets", requestBody)
//     if err != nil {
//         http.Error(writer, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     // Set request headers
//     req.Header.Set("Content-Type", "application/json")

//     // Send the HTTP request
//     response, err := httpClient.Do(req)
//     if err != nil {
//         http.Error(writer, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     defer response.Body.Close()

//     // Read response body
//     responseData, err := ioutil.ReadAll(response.Body)
//     if err != nil {
//         http.Error(writer, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     // Write response back to the client
//     writer.WriteHeader(response.StatusCode)
//     writer.Write(responseData)
// }

// func tweet(writer http.ResponseWriter, request *http.Request) {
// 	requestTokenURL := "https://api.twitter.com/oauth/request_token?oauth_callback=oob&x_auth_access_type=write"
// 	accessTokenURL := "https://api.twitter.com/oauth/access_token"
// 	tweetURL := "https://api.twitter.com/2/tweets"

// 	// Get request token
// 	oauthConfig := &oauth1.Config{
// 		ConsumerKey:    consumerKey,
// 		ConsumerSecret: consumerSecret,
// 		CallbackURL:    "oob",
// 		Endpoint: oauth1.Endpoint{
// 			RequestTokenURL: requestTokenURL,
// 			AuthorizeURL:    "",
// 			AccessTokenURL:  accessTokenURL,
// 		},
// 	}
// 	requestToken, _, err := oauthConfig.RequestToken()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	authURL, err := oauthConfig.AuthorizationURL(requestToken)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Please go here and authorize:", authURL)

// 	// Get verifier
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Paste the PIN here: ")
// 	verifier, _ := reader.ReadString('\n')
// 	verifier = strings.TrimSpace(verifier)

// 	// Get access token
// 	accessToken, accessSecret, err := oauthConfig.AccessToken(requestToken, "", verifier)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Make tweet request
// 	token := oauth1.NewToken(accessToken, accessSecret)
// 	client := oauthConfig.Client(oauth1.NoContext, token)
// 	req, err := http.NewRequest("POST", tweetURL, strings.NewReader(`{"text": "Hello world!"}`))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Fprintf(writer, "%s\n", body)
// }



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