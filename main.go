package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"vtask/database"
	"vtask/handlers"
	"vtask/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	envFile := ".env"

	if err := godotenv.Load(envFile); err != nil {
		log.Println("Cannot load env file")
		log.Fatalf(err.Error())
	}


	
	log.Println("Successful load env file")

	database.InitDB()
	defer database.DB.Close()

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		database.DB,
		os.DirFS("migrations"),
	)

	if err != nil {
		log.Println("failed to set provider")
		log.Fatalf(err.Error())
	}

	res, err := provider.Up(context.Background())
	if err != nil{
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println(res)

	handlers.InitJWT()

	router := mux.NewRouter()

	router.Use(middleware.AuthMiddleware)

	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/ads", handlers.GetAdsFeed).Methods("GET")

	router.HandleFunc("/ads", handlers.CreateAd).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
