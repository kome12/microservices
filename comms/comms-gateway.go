package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"google.golang.org/grpc"

	pb "github.com/microservices/microservices/comms"
)

type server struct {
	pb.UnimplementedCommsServer
}

func (s *server) SendConfirmation(ctx context.Context, in *pb.ConfirmationRequest) (*pb.ConfirmationResponse, error) {
	log.Printf("Received: %s", in.Status)
	sendStatusEmail := sendEmail(in)
	return &pb.ConfirmationResponse{Success: sendStatusEmail}, nil
}

func sendEmail(in *pb.ConfirmationRequest) bool  {
	from := mail.NewEmail("Koichi", "kmiyatake@kiraku.io")
	// subject := "Sending with SendGrid is Fun"
	var subject string
	if in.GetStatus() == "succeeded" {
		subject = "Thank you for your purchase"
	} else {
		subject = "Error in purchase"
	}
	to := mail.NewEmail("Ko", in.Email)
	htmlContent := "Payment ID: <strong>" + in.Id + "</strong>"
	plainTextContent := "Amount: " + fmt.Sprint(in.GetAmount())
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return false
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return true
	}
}

func comms() {
	from := mail.NewEmail("Koichi", "kmiyatake@kiraku.io")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Ko", "ko.take5481@gmail.com")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

func init() {
	 err := godotenv.Load(".env")

	 if err != nil {
		 log.Fatal("Error loading .env file")
	 }
}

func main()  {
	// comms()
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
