package handleapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	pbComms "github.com/microservices/microservices/comms"
	pbPayments "github.com/microservices/microservices/payments"
	pbProducts "github.com/microservices/microservices/products"
	"google.golang.org/grpc"
)

type ShoeInput struct {
	ShoeID string `json:"shoeId"`
	Price int64 `json:"price"`
	Email string `json:"email"`
	ProductName string `json:"productName"`
}

type PurchaseStatus struct {
	PaymentSuccessful bool
	EmailSuccessful bool
	Message string
}

func handleError(err error) {
	if err != nil {
    log.Fatalf("could not process: %v", err)
	}
}

func handleGRPCConnectionError(err error) {
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
}

func openPaymentGRPC() *grpc.ClientConn {
	connPayment, err := grpc.Dial(os.Getenv("ADDRESS") + ":" + os.Getenv("PAYMENT_PORT"), grpc.WithInsecure(), grpc.WithBlock())
	handleGRPCConnectionError(err)
	return connPayment
}

func openCommsGRPC() *grpc.ClientConn {
	connComms, err := grpc.Dial(os.Getenv("ADDRESS") + ":" + os.Getenv("COMMS_PORT"), grpc.WithInsecure(), grpc.WithBlock())
	handleGRPCConnectionError(err)
	return connComms
}

func openProductsGRPC() *grpc.ClientConn {
	connProducts, err := grpc.Dial(os.Getenv("ADDRESS") + ":" + os.Getenv("PRODUCTS_PORT"), grpc.WithInsecure(), grpc.WithBlock())
	handleGRPCConnectionError(err)
	return connProducts
}

func createContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	return ctx, cancel
}

func GetShoes(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("api processing getShoes...")
	EnableCors(&rw)
	connProducts := openProductsGRPC()
	cProducts := pbProducts.NewProductsClient(connProducts)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	res, err := cProducts.GetShoes(ctx, &pbProducts.ShoesRequest{});

	handleError(err)

	defer connProducts.Close()
	// return r;
	rw.Header().Set("Content-Type", "application/json")
  json.NewEncoder(rw).Encode(res.Shoes)
}

func GetShoe(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("api processing getShoe (just 1)...")
	EnableCors(&rw)
	connProducts := openProductsGRPC()
	cProducts := pbProducts.NewProductsClient(connProducts)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	params := mux.Vars(r)
	res, err := cProducts.GetShoe(ctx, &pbProducts.ShoeRequest{Id: params["id"]})

	handleError(err)

	defer connProducts.Close()
	// return r
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(res)
}

func Purchase(rw http.ResponseWriter, r *http.Request) {
	EnableCors(&rw)
	if r.Method == "POST" {
		fmt.Println("api processing purchase...")
		connPayment := openPaymentGRPC()
		cPayment := pbPayments.NewPaymentsClient(connPayment)
		
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		defer connPayment.Close()
	
		var shoeInput ShoeInput
		_ = json.NewDecoder(r.Body).Decode(&shoeInput)
	
		if shoeInput.Email != "" {
			fmt.Println("api processing purchase (stripe)...")
			res, err := cPayment.CreateCharge(ctx, &pbPayments.ChargeRequest{Amount: shoeInput.Price})
	
			handleError(err)
	
			rw.Header().Set("Content-Type", "application/json")
	
			var status PurchaseStatus
			status.PaymentSuccessful = res.GetSuccess()
	
			if res.GetSuccess() {
				fmt.Println("api processing purchase (sendgrid)...")
				resComms := sendEmail(shoeInput, res)
				status.EmailSuccessful = resComms
		
				if resComms {
					log.Printf("Charged: %t", res.GetSuccess())
					status.Message = "Successful!"
					json.NewEncoder(rw).Encode(status)
				} else {
					status.Message = "Unsuccessful"
					json.NewEncoder(rw).Encode(status)
				}
			} else {
				status.Message = "Purchase unsuccessful"
				json.NewEncoder(rw).Encode(status)
			}
		}
	}
}

func sendEmail(shoeInput ShoeInput, paymentRes *pbPayments.ChargeResponse) bool {
	connComms := openCommsGRPC()
	cComms := pbComms.NewCommsClient(connComms)
	ctx, cancel := createContext()
	defer cancel()
	defer connComms.Close()

	res, err := cComms.SendConfirmation(ctx, &pbComms.ConfirmationRequest{
		Email:  shoeInput.Email,
		Id:     paymentRes.Id,
		Status: paymentRes.GetStatus(),
		Amount: paymentRes.Amount,
		ProductName: shoeInput.ProductName,
	})

	handleError(err)

	return res.GetSuccess()
}

func init() {
	 err := godotenv.Load(".env")

	 if err != nil {
		 log.Fatal("Error loading .env file")
	 }
}

func EnableCors(w *http.ResponseWriter) {
	header := (*w).Header()
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}