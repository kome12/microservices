package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	pb "github.com/microservices/microservices/payments"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPaymentsServer
}

func (s *server) CreateCharge(ctx context.Context, in *pb.ChargeRequest) (*pb.ChargeResponse, error) {
		log.Printf("Received: %v", in.GetAmount())
		stripeChargeResponse := createCharge(in.GetAmount())
		return &pb.ChargeResponse{
			Success: stripeChargeResponse.Status == "succeeded",
			Id:      stripeChargeResponse.ID,
			Amount:  stripeChargeResponse.Amount,
			Status: stripeChargeResponse.Status,
		}, nil
}

func createCharge(amount int64) *stripe.Charge {
	fmt.Println("preparing new charge to send to Stripe...")
	stripe.Key = os.Getenv("STRIPE_KEY")

	// `source` is obtained with Stripe.js; see https://stripe.com/docs/payments/accept-a-payment-charges#web-create-token
	params := &stripe.ChargeParams{
		Amount: stripe.Int64(amount),
		Currency: stripe.String(string(stripe.CurrencyJPY)),
		Description: stripe.String("My First Test Charge (created for API docs)"),
		Source: &stripe.SourceParams{Token: stripe.String("tok_mastercard")},
	}
	c, _ := charge.New(params)

	return c
}

func init() {
	 err := godotenv.Load(".env")

	 if err != nil {
		 log.Fatal("Error loading .env file")
	 }
}

func main() {
	lis, err := net.Listen("tcp", ":" + os.Getenv("PAYMENT_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPaymentsServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
