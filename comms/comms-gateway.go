package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"google.golang.org/grpc"

	pb "github.com/microservices/microservices/comms"
)

type server struct {
	pb.UnimplementedCommsServer
}

type DynamicData struct {
	price string
	paymentId string
	shoeName string
}

func (s *server) SendConfirmation(ctx context.Context, in *pb.ConfirmationRequest) (*pb.ConfirmationResponse, error) {
	log.Printf("Received: %s", in.Status)
	sendStatusEmail := sendEmail(in)
	return &pb.ConfirmationResponse{Success: sendStatusEmail}, nil
}

func sendEmail(in *pb.ConfirmationRequest) bool  {
	fmt.Println("preparing new email to send via SendGrid...")

	// TODO: Refactor below so it's cleaner and more dynamic
	var byteData string
	byteData += "{"
	byteData += `"from":{"email":"kmiyatake@kiraku.io","name":"Koichi"},`
	byteData += `
		"personalizations":[
			{
				"to": [
					{
						"email":"` + in.Email + `"},
				],
				"dynamic_template_data": {
					"price":"` + fmt.Sprint(in.GetAmount()) + `",
					"paymentId":"` + in.GetId() + `",
					"shoeName":"` + in.GetProductName() + `",
				}
			}
		],`
	byteData += `"template_id":"d-e3ac2bfba4014f6e80e5de389e7b0bf5"`
	byteData += "}"

	apiKey := os.Getenv("SENDGRID_API_KEY")
	host := "https://api.sendgrid.com"
	request := sendgrid.GetRequest(apiKey, "/v3/mail/send", host)
	request.Method = "POST"
	request.Body = []byte(byteData)
	response, err := sendgrid.API(request)

	if err != nil {
		log.Println("err in email:", err)
		return false
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return response.StatusCode == 200 || response.StatusCode == 202
	}
}

func init() {
	 err := godotenv.Load(".env")

	 if err != nil {
		 log.Fatal("Error loading .env file")
	 }
}

func main()  {
	lis, err := net.Listen("tcp", ":" + os.Getenv("COMMS_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCommsServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
