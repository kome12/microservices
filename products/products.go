package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	p "github.com/microservices/microservices/products"
)

type server struct {
	p.UnimplementedProductsServer
}

func (s *server) GetShoes(ctx context.Context, in *p.ShoesRequest) (*p.ShoesResponse, error) {
	fmt.Println("came into getShoes...")
	shoes, _ := readProducts()
	var shoesRequest p.ShoesResponse
	for _, shoe := range shoes.Shoes {
		var shoeRequest p.ShoeResponse
		shoeRequest.Id = shoe.ID
		shoeRequest.Name = shoe.Name
		shoeRequest.Price = shoe.Price
		shoesRequest.Shoes = append(shoesRequest.Shoes, &shoeRequest)
	}
	return &p.ShoesResponse{Shoes: shoesRequest.Shoes}, nil
}

func (s *server) GetShoe(ctx context.Context, in *p.ShoeRequest) (*p.ShoeResponse, error) {
	fmt.Println("came into getShoe (just 1)...")
	shoes, _ := readProducts()

	var shoeFound Shoe

	for _, shoe := range shoes.Shoes {
		if shoe.ID == in.GetId() {
			shoeFound = shoe
			break
		}
	}

	if shoeFound.ID != "" {
		return &p.ShoeResponse{Id: shoeFound.ID, Name: shoeFound.Name, Price: shoeFound.Price}, nil
	} else {
		return &p.ShoeResponse{}, nil
	}
}

type ShoeInput struct {
	ShoeID string `json:"shoeId"`
	Price int64 `json:"price"`
	Email string `json:"email"`
}

type Shoe struct {
	ID string `json:"id"`
  Name string `json:"name"`
  Price int64 `json:"price"`
}

type Shoes struct {
	Shoes []Shoe `json:"shoes"`
}

func readProducts() (shoes Shoes, err error)  {
	content, err := ioutil.ReadFile("data/products.json")

	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(content, &shoes)
	if err != nil {
			log.Fatal("Error during Unmarshal(): ", err)
	}
	return shoes, err
}

func init() {
	 err := godotenv.Load(".env")

	 if err != nil {
		 log.Fatal("Error loading .env file")
	 }
}

func main()  {
	lis, err := net.Listen("tcp", ":" + os.Getenv("PRODUCTS_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	p.RegisterProductsServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}