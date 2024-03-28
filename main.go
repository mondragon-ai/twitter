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
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	consumerKey = os.Getenv("CONSUMER_KEY")
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

func GenerateOAuthNonce(length int) (string, error) {
    // Generate random bytes
    randomBytes := make([]byte, length)
    _, err := rand.Read(randomBytes)
    if err != nil {
        return "", err
    }

    // Base64 encode the random bytes
    base64Str := base64.StdEncoding.EncodeToString(randomBytes)

    // Strip non-word characters
    var result strings.Builder
    for _, char := range base64Str {
        if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
            result.WriteRune(char)
        }
    }

    return result.String(), nil
}


func tweet(writer http.ResponseWriter, request *http.Request) {

	// Retrieve environment variables
	consumerKey = os.Getenv("CONSUMER_KEY")
	consumerSecret = os.Getenv("CONSUMER_SECRET")
	accessTokenKey = os.Getenv("ACCESS_TOKEN_KEY")
	accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")

	// Generate signature
	// signature := generateSignature(baseString, consumerSecret, accessTokenSecret)
    signature := prepareOAuthSignature(consumerKey, accessTokenKey, consumerSecret, accessTokenSecret)
	fmt.Println("[signature]")
	fmt.Println(signature)

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
func generateSignature(method, baseURL string, consumerSecret, tokenSecret string) (string, error) {
	// Step 1: Collect parameters and encode them
    // Step 4: Create base string
    baseString := strings.ToUpper(method) + "&" + url.QueryEscape(baseURL)


    // Step 5: Create signing key
    signingKey := url.QueryEscape(consumerSecret) + "&" + url.QueryEscape(tokenSecret)

    // Step 6: Calculate signature
    h := hmac.New(sha1.New, []byte(signingKey))
    h.Write([]byte(baseString))
    signature := h.Sum(nil)

    // Step 7: Base64 encode the signature
    encodedSignature := base64.StdEncoding.EncodeToString(signature)

    return encodedSignature, nil
}


func prepareOAuthSignature(oauthConsumerKey, oauthToken, consumerSecret, tokenSecret string) string {
    // Constants
    oauthSignatureMethod := "HMAC-SHA1"
    oauthVersion := "1.0"
    oauthNonce := "qac8abeMCg8" // Assuming you have your way of generating nonce
    oauthTimestamp := fmt.Sprintf("%d", time.Now().Unix())

    // Collect parameters
    params := map[string]string{
        "oauth_consumer_key":     oauthConsumerKey,
        "oauth_token":            oauthToken,
        "oauth_signature_method": oauthSignatureMethod,
        "oauth_timestamp":        oauthTimestamp,
        "oauth_nonce":            oauthNonce,
        "oauth_version":          oauthVersion,
    }

    // Sort parameters by key
    var keys []string
    for k := range params {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    // Construct parameter string
    var paramStr strings.Builder
    for _, k := range keys {
        if paramStr.Len() > 0 {
            paramStr.WriteString("&")
        }
        paramStr.WriteString(url.QueryEscape(k))
        paramStr.WriteString("=")
        paramStr.WriteString(url.QueryEscape(params[k]))
    }
	fmt.Println("[paramStr]")
	fmt.Println(paramStr)
	

    // Construct signature base string
    httpMethod := "POST"
    baseURL := "https://api.twitter.com/2/tweets"
    baseStr := fmt.Sprintf("%s&%s&%s", httpMethod, url.QueryEscape(baseURL), url.QueryEscape(paramStr.String()))
	fmt.Println("[baseStr]")
	fmt.Println(baseStr)

    // Construct signing key
    signingKey := fmt.Sprintf("%s&%s", url.QueryEscape(consumerSecret), url.QueryEscape(tokenSecret))
	fmt.Println("[signingKey]")
	fmt.Println(signingKey)

    // Calculate signature
    h := hmac.New(sha1.New, []byte(signingKey))
    h.Write([]byte(baseStr))
    signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

    // Construct authorization header
    var authHeader strings.Builder
    authHeader.WriteString(`OAuth oauth_consumer_key="`)
    authHeader.WriteString(url.QueryEscape(oauthConsumerKey))
    authHeader.WriteString(`", oauth_token="`)
    authHeader.WriteString(url.QueryEscape(oauthToken))
    authHeader.WriteString(`", oauth_signature_method="`)
    authHeader.WriteString(url.QueryEscape(oauthSignatureMethod))
    authHeader.WriteString(`", oauth_timestamp="`)
    authHeader.WriteString(url.QueryEscape(oauthTimestamp))
    authHeader.WriteString(`", oauth_nonce="`)
    authHeader.WriteString(url.QueryEscape(oauthNonce))
    authHeader.WriteString(`", oauth_version="`)
    authHeader.WriteString(url.QueryEscape(oauthVersion))
    authHeader.WriteString(`", oauth_signature="`)
    authHeader.WriteString(url.QueryEscape(signature))
    authHeader.WriteString(`"`)

    return authHeader.String()
}

func percentEncode(src string) string {
    var encoded strings.Builder

    for _, b := range src {
        switch {
        case (b >= '0' && b <= '9') || (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || b == '-' || b == '.' || b == '_' || b == '~':
            encoded.WriteRune(b)
        default:
            encoded.WriteString(fmt.Sprintf("%%%X", b))
        }
    }

    return encoded.String()
}