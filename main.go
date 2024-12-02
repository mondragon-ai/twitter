package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/twitter/config"
	"github.com/twitter/controller"
	"github.com/twitter/helper"
	"github.com/twitter/mentions"
	"github.com/twitter/router"
	"github.com/twitter/service"
)

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

    server := http.Server{Addr: ":" + port, Handler: routes}

    log.Printf("Starting server on port %s\n", port)
    err = server.ListenAndServe()
    helper.PanicIfError(err)

}


// package main

// import (
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/gorilla/mux"
// 	"github.com/robfig/cron"
// 	"github.com/twitter/apis/twitter"
// )


// func main() {
// 	//Load env
// 	// err := godotenv.Load()
// 	// if err != nil {
// 	// 	log.Fatal("Error loading .env file")
// 	// 	fmt.Println("Error loading .env file")
// 	// }

// 	fmt.Println("Starting Server")

// 	m := mux.NewRouter()
// 	m.HandleFunc("/", func(writer http.ResponseWriter, _ *http.Request) {
// 		writer.WriteHeader(200)
// 		fmt.Fprintf(writer, "Server is up and running")
// 	})
// 	m.HandleFunc("/tweet", twitter.Tweet).Methods("POST")

// 	//Start Server
// 	server := &http.Server{
// 		Handler: m,
// 	}

//     // Retrieve the port from the environment variable or default to 8080
//     port := os.Getenv("PORT")
//     if port == "" {
//         port = "8080"
//     }

//     server.Addr = ":" + port
// 	// server.Addr = ":8080"

// 	// Setup Cron job
// 	c := cron.New()
// 	c.AddFunc("0 30 * * * *", func() {
// 		fmt.Println("Ready to Tweet...")
// 		if shouldPost() && isWithinAllowedTimezone() {
// 			fmt.Println("[TWEETED]")
// 			twitter.Post()
// 		}
// 	})
// 	c.Start()

// 	// Start HTTP server
// 	go func() {
// 		if err := server.ListenAndServe(); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	// Wait for the Cron job to run
// 	// time.Sleep(5 * time.Minute)


//     // Wait for a signal to gracefully exit the program
//     sig := make(chan os.Signal, 1)
//     signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
//     <-sig


// 	// Stop the Cron job scheduler
// 	c.Stop()
// }


// func shouldPost() bool {
//     // You can adjust the weights as needed.
//     weights := map[bool]int{
//         true:  2,  
//         false: 8, 
//     }
//     totalWeight := 0
//     for _, weight := range weights {
//         totalWeight += weight
//     }
//     num := rand.Intn(totalWeight)
//     for val, weight := range weights {
//         if num < weight {
//             return val
//         }
//         num -= weight
//     }
//     return false
// }


// func isWithinAllowedTimezone() bool {
//     loc, err := time.LoadLocation("America/Chicago")
//     if err != nil {
//         return false
//     }
//     current := time.Now().In(loc)
//     hour := current.Hour()
//     return hour >= 5 && hour < 23 
// }

