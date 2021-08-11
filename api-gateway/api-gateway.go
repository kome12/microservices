package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	pbComms "github.com/microservices/microservices/comms"
	pbPayments "github.com/microservices/microservices/payments"

	// "github.com/microservices/products"
	handleapi "github.com/microservices/handleAPI"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the API Gateway!")
	fmt.Println("Endpoint hit: homepage")
}

func createCharge(ctx context.Context, cPayment pbPayments.PaymentsClient, cComms pbComms.CommsClient) bool  {
		r, err := cPayment.CreateCharge(ctx, &pbPayments.ChargeRequest{Amount: 5000})

	if err != nil {
        log.Fatalf("could not process: %v", err)
	}

	if r.GetSuccess() {
			log.Printf("Charged: %t", r.GetSuccess())
		emailRes, err := cComms.SendConfirmation(ctx, &pbComms.ConfirmationRequest{
			Status: "Successful",
		})

		log.Printf("Email success: %t", emailRes.GetSuccess())

		if err != nil {
        log.Fatalf("could not process: %v", err)
				return false
		}
		return true
	}
	return false
}

func sendEmail(ctx context.Context, cComms pbComms.CommsClient)  {
	emailRes, err := cComms.SendConfirmation(ctx, &pbComms.ConfirmationRequest{
			Status: "Successful",
		})

		log.Printf("Email success: %t", emailRes.GetSuccess())

		if err != nil {
        log.Fatalf("could not process: %v", err)
		}
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/", homePage)

	// router.HandleFunc("/createCharges", func(rw http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(rw, "Processing...")
	// 	finalResult := createCharge(ctx, cPayment, cComms)
	// 	if finalResult {
	// 		fmt.Fprintf(rw, "Successfully created charge and emailed result!")
	// 	} else {
	// 		fmt.Fprintf(rw, "Something went wrong...")
	// 	}
	// }).Methods("POST");

	// router.HandleFunc("/sendEmail", func(rw http.ResponseWriter, r *http.Request) {
	// 	sendEmail(ctx, cComms)
	// })

	// router.HandleFunc("/api/shoes", func(rw http.ResponseWriter, r *http.Request) {
	// 	EnableCors(&rw)
	// 	shoes := handleapi.GetShoes()
	// 	rw.Header().Set("Content-Type", "application/json")
  // 	json.NewEncoder(rw).Encode(shoes.Shoes)
	// }).Methods("GET")

	router.HandleFunc("/api/shoes", handleapi.GetShoes).Methods("GET")

	// router.HandleFunc("/api/shoes/{id}", func(rw http.ResponseWriter, r *http.Request) {
	// 	EnableCors(&rw)
	// 	params := mux.Vars(r)

	// 	shoe := handleapi.GetShoe(params["id"])
	// 	rw.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(rw).Encode(shoe)
	// }).Methods("GET")

	router.HandleFunc("/api/shoes/{id}", handleapi.GetShoe).Methods("GET")

	router.HandleFunc("/api/purchases", handleapi.Purchase).Methods("POST", "OPTIONS")

	// router.HandleFunc("/api/shoes", products.GetShoes).Methods("GET")
	// router.HandleFunc("/api/shoes/{id}", products.GetShoe).Methods("GET")
	
	// router.HandleFunc("/api/purchases", products.CreatePurchases).Methods("POST", "OPTIONS")

	// router.HandleFunc("/purchases", func(rw http.ResponseWriter, r *http.Request) {
	// 	products.EnableCors(&rw)
	// 	fmt.Println("came into post shoe")
	// 	rw.Header().Set("Content-Type", "application/json")
  // 	var shoeInput products.ShoeInput
  // 	_ = json.NewDecoder(r.Body).Decode(shoeInput)
	// 	fmt.Println("shoeInput:", shoeInput)

	// 	json.NewEncoder(rw).Encode(nil)
	// }).Methods("POST")

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
	// connPayment, err := grpc.Dial(os.Getenv("ADDRESS") + ":" + os.Getenv("PAYMENT_PORT"), grpc.WithInsecure(), grpc.WithBlock())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer connPayment.Close()
	// cPayment := pbPayments.NewPaymentsClient(connPayment)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	// defer cancel()

	// connComms, err := grpc.Dial(os.Getenv("ADDRESS") + ":" + os.Getenv("COMMS_PORT"), grpc.WithInsecure(), grpc.WithBlock())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer connComms.Close()
	// cComms := pbComms.NewCommsClient(connComms)

	// connProducts, err := grpc.Dial(os.Getenv("ADDRESS") + ":" + os.Getenv("PRODUCTS_PORT"), grpc.WithInsecure(), grpc.WithBlock())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer connProducts.Close()
	// cProducts := pbProducts.NewProductsClient(connProducts)



	// ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	// defer cancel()

	handleRequests()
}