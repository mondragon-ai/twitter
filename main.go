package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/robfig/cron"
	"github.com/twitter/apis/twitter"
)


func main() {
	//Load env
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// 	fmt.Println("Error loading .env file")
	// }

	fmt.Println("Starting Server")

	m := mux.NewRouter()
	m.HandleFunc("/", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(200)
		fmt.Fprintf(writer, "Server is up and running")
	})
	m.HandleFunc("/tweet", twitter.Tweet).Methods("POST")

	//Start Server
	server := &http.Server{
		Handler: m,
	}
	server.Addr = ":8080"

	// Setup Cron job
	c := cron.New()
	c.AddFunc("0 30 * * * *", func() {
		fmt.Println("Ready to Tweet...")
		if shouldPost() && isWithinAllowedTimezone() {
			fmt.Println("[TWEETED]")
			twitter.Post()
		}
	})
	c.Start()

	// Start HTTP server
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for the Cron job to run
	// time.Sleep(5 * time.Minute)


    // Wait for a signal to gracefully exit the program
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig


	// Stop the Cron job scheduler
	c.Stop()
}


func shouldPost() bool {
    // You can adjust the weights as needed.
    weights := map[bool]int{
        true:  7,  
        false: 3, 
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