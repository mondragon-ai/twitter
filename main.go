package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"

	"github.com/twitter/config"
	"github.com/twitter/controller"
	"github.com/twitter/helper"
	"github.com/twitter/mentions"
	"github.com/twitter/router"
	"github.com/twitter/service"
)

// Middleware to enable CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	godotenv.Load(".env")
	fmt.Printf("Start server")

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// database
	db, err := config.DatabaseConnection(connString)
	helper.PanicIfError(err)

	// Ensure the mention table exists
	// err = config.ResetDB(db)
	// if err != nil {
	// 	log.Fatal("Could not delete tables: ", err)
	// }

	err = config.CreateDB(db)
	if err != nil {
		log.Fatal("Could not create tables: ", err)
	}

	// repository
	mentionRepository := mentions.MentionCrud(db)

	// service
	mentionService := service.NewMentionServiceImpl(mentionRepository)
	twitterService := service.NewTwitterServiceImpl(mentionRepository)

	// Controllers
	mentionsController := controller.NewMentionsController(mentionService)
	twitterController := controller.NewTwitterController(twitterService)

	// router
	routes := router.MentionsRouter(mentionsController, twitterController)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    server := http.Server{
		Addr: ":" + port, 
		Handler: enableCORS(routes),
	}

	// Setup Cron job
	c := cron.New()
	c.AddFunc("0 30 * * * *", func() {
		log.Println("Ready to Tweet...")
		if shouldPost() && isWithinAllowedTimezone() {
			log.Println("[TWEETED]")
			InternalCall()
		}
	})
	c.Start()

	// Channel to receive termination signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Run the server in a goroutine
	go func() {
		log.Printf("Starting server on port %s\n", port)
		err := server.ListenAndServe()
		helper.PanicIfError(err)
	}()

	// Wait for a termination signal
	<-sig
	log.Println("Shutting down gracefully...")

	// Stop the Cron job scheduler
	c.Stop()
	log.Println("Cron job scheduler stopped")
}

func InternalCall() {

	tweetType := postType()
	body := map[string]interface{}{
		"type": tweetType,
	}

    // Prepare request body if applicable
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Printf("Failed to marshal request body: %v", err)
	}
	requestBody := bytes.NewBuffer(bodyJSON)

	url := "http://localhost:8080/api/twitter/tweet"
    req, err := http.NewRequest(http.MethodPost, url, requestBody)
    if err != nil {
        log.Fatal(err)
    }
    req.Header.Add("Accept", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Fatal(err)
    }

    if resp == nil || resp.Body == nil {
        log.Fatal("Received nil response or response body from Twitter API")
    }

    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := io.ReadAll(resp.Body)
        log.Printf("Twitter API error: %s", string(bodyBytes))
    }

    defer resp.Body.Close()
    if err != nil {
        log.Fatal(err)
    }
}

func postType() string {
    weights := map[string]int{
        "clone":  15,
        "create": 60,
        "article": 15,
        "thread": 10,
    }
    totalWeight := 0
    for _, weight := range weights {
        totalWeight += weight
    }

    num := rand.Intn(totalWeight)
    for val, weight := range weights {
        if num < weight {
            return val
        }
        num -= weight
    }
    return "create"
}

func shouldPost() bool {
    // You can adjust the weights as needed.
    weights := map[bool]int{
        true:  6,  
        false: 4, 
    }
    totalWeight := 0
    for _, weight := range weights {
        totalWeight += weight
    }
    num := rand.Intn(totalWeight)
    for val, weight := range weights {
        if num < weight {
            return val
        }
        num -= weight
    }
    return false
}

func isWithinAllowedTimezone() bool {
    loc, err := time.LoadLocation("America/Chicago")
    if err != nil {
        return false
    }
    current := time.Now().In(loc)
    hour := current.Hour()
    return hour >= 5 && hour < 23 
}

