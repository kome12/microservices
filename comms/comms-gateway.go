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
	// from := mail.NewEmail("Koichi", "kmiyatake@kiraku.io")
	// // subject := "Sending with SendGrid is Fun"
	// var subject string
	// if in.GetStatus() == "succeeded" {
	// 	subject = "Thank you for your purchase"
	// } else {
	// 	subject = "Error in purchase"
	// }
	// to := mail.NewEmail("Ko", in.Email)
	// htmlContent := "Payment ID: <strong>" + in.Id + "</strong>"
	// plainTextContent := "Amount: " + fmt.Sprint(in.GetAmount())
	
	// message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	// // personalization := mail.NewPersonalization(to, from)
	// var personolization mail.Personalization
	// message.Personalizations = []*mail.Personalization{&personolization}

	// client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	// response, err := client.Send(message)
	// if err != nil {
	// 	log.Println(err)
	// 	return false
	// } else {
	// 	fmt.Println(response.StatusCode)
	// 	fmt.Println(response.Body)
	// 	fmt.Println(response.Headers)
	// 	return true
	// }

	// dynamicData := DynamicData{price: fmt.Sprint(in.GetAmount()), paymentId: in.GetId(), shoeName: in.GetProductName()}

	var byteData string
	byteData += "{"
	byteData += `"from":{"email":"kmiyatake@kiraku.io","name":"Koichi"},`
	// byteData += `"to":[{"email":"` + in.Email + `", "name":"Ko"}],`
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
	// request.Body = []byte(`{
	// 	"from": {
  //   	"email": "kmiyatake@kiraku.io", 
  //   	"name": "Koichi"
  // 	},
	// 	"personalizations": [
	// 		{
	// 			"to": [
	// 				{
	// 					"email": "ko.take5481@gmail.com"
	// 				}
	// 			],
	// 			"dynamic_template_data": {
  //        	"price": "New Value 1", 
  //        	"paymentId": "1", 
  //      	}, 
	// 		}
	// 	],
	// 	"template_id": "d-e3ac2bfba4014f6e80e5de389e7b0bf5"
	// }`)
	// "template_id": "d-e3ac2bfba4014f6e80e5de389e7b0bf5",
	// 	"personalizations": [
	// 		{
	// 			"custom_args": {
  //       	"price": "New Value 1", 
  //       	"paymentId": "1", 
  //     	}, 
	// 		}
	// 	]
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
