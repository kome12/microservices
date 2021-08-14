package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	// "github.com/microservices/products"
	handleapi "github.com/microservices/handleAPI"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the API Gateway!")
	fmt.Println("Endpoint hit: homepage")
}

func handleRequests() {
	router := mux.NewRouter()

	router.HandleFunc("/", homePage)
	router.HandleFunc("/api/shoes", handleapi.GetShoes).Methods("GET")
	router.HandleFunc("/api/shoes/{id}", handleapi.GetShoe).Methods("GET")
	router.HandleFunc("/api/purchases", handleapi.Purchase).Methods("POST", "OPTIONS")

	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":" + port, router))
}

func EnableCors(w *http.ResponseWriter) {
	header := (*w).Header()
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

func init() {
	 err := godotenv.Load(".env")

	 if err != nil {
		 log.Fatal("Error loading .env file")
	 }
}

func main() {
	handleRequests()
}