This was created during my time as a student at [Code Chrysalis](https://www.codechrysalis.io/)

# Microservices in Go using gRPC

This is a small group of microservices written in Go.

It has 4 parts

1. API Gateway
2. Payment Gateway (integrated with Stripe API)
3. Communications Gateway (integrated with SendGrid API)
4. Products (reading data from JSON file)

## .env Parameters

```
PORT
PAYMENT_PORT
COMMS_PORT
PRODUCTS_PORT
ADDRESS
STRIPE_KEY
SENDGRID_API_KEY
```

PORT = Port for the API Gateway

PAYMENT_PORT = Port for the Payment Gateway

COMMS_PORT = Port for the Communications Gateway

PRODUCTS_PORT = Port for the Products

ADDRESS = Address for each of the ports, ie. localhost

STRIPE_KEY = Stripe API key

SENDGRID_API_KEY = SendGrid API key

## How to run the microservices locally

1. Open 4 terminals windows. Please run the below steps in separate windows.
2. From the root of the project, run `go run api-gateway/api-gateway.go` to start the API Gateway
3. From the root of the project, run `go run payment-gateway/payment-gateway.go` to start the Payment Gateway
4. From the root of the project, run `go run comms-gateway/comms-gateway.go` to start the Communications Gateway
5. From the root of the project, run `go run products/products.go` to start the Products

## Tech Stack

- [Go](https://golang.org/)
- [gRPC](https://grpc.io/)
- [Stripe](https://stripe.com/)
- [SendGrid](https://sendgrid.com/)
