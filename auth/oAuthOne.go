package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

type TwitterAuth interface {
    PrepareOAuthSignature(oauthConsumerKey, oauthToken, consumerSecret, tokenSecret string) string
}


func PrepareOAuthSignature(oauthConsumerKey, oauthToken, consumerSecret, tokenSecret string) string {
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
	

    // Construct signature base string
    httpMethod := "POST"
    baseURL := "https://api.twitter.com/2/tweets"
    baseStr := fmt.Sprintf("%s&%s&%s", httpMethod, url.QueryEscape(baseURL), url.QueryEscape(paramStr.String()))

    // Construct signing key
    signingKey := fmt.Sprintf("%s&%s", url.QueryEscape(consumerSecret), url.QueryEscape(tokenSecret))

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