package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	p "github.com/microservices/microservices/products"
)

type server struct {
	p.UnimplementedProductsServer
}

func (s *server) GetShoes(ctx context.Context, in *p.ShoesRequest) (*p.ShoesResponse, error) {
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
	shoes, _ := readProducts()

	var shoeFound Shoe

	for _, shoe := range shoes.Shoes {
		if shoe.ID == in.GetId() {
			fmt.Println("shoe:", shoe)
			// json.NewEncoder(rw).Encode(shoe)
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

// func GetShoes(rw http.ResponseWriter, r *http.Request)  {
// 	EnableCors(&rw)

// 	shoes, _ := readProducts()
	
// 	fmt.Println("payload:", shoes)
// 	rw.Header().Set("Content-Type", "application/json")
//   json.NewEncoder(rw).Encode(shoes.Shoes)
// }

func GetShoe(rw http.ResponseWriter, r *http.Request)  {
	fmt.Println("came into getShoe")
	EnableCors(&rw)
	params := mux.Vars(r)

	shoes, _ := readProducts()
	rw.Header().Set("Content-Type", "application/json")

	var shoeFound Shoe

	for _, shoe := range shoes.Shoes {
		if shoe.ID == params["id"] {
			fmt.Println("shoe:", shoe)
			// json.NewEncoder(rw).Encode(shoe)
			shoeFound = shoe
			break
		}
	}

	if shoeFound.ID != "" {
		json.NewEncoder(rw).Encode(shoeFound)
	} else {
		json.NewEncoder(rw).Encode(nil)
	}
}

func CreatePurchases(rw http.ResponseWriter, r *http.Request)  {
		fmt.Println("came into post shoe")
		EnableCors(&rw);
		rw.Header().Set("Content-Type", "application/json")

  	var shoeInput ShoeInput
  	_ = json.NewDecoder(r.Body).Decode(&shoeInput)
		fmt.Println("shoeInput:", shoeInput)
		
		json.NewEncoder(rw).Encode(nil)
}

func readProducts() (shoes Shoes, err error)  {
	content, err := ioutil.ReadFile("data/products.json")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("content:", content)

	// var shoes Shoes
	err = json.Unmarshal(content, &shoes)
	if err != nil {
			log.Fatal("Error during Unmarshal(): ", err)
	}
	return shoes, err
}

func readShoes() (shoes p.ShoesResponse, err error) {
	content, err := ioutil.ReadFile("data/products.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("content in readShoes():", content)

	err = json.Unmarshal(content, &shoes)
	fmt.Println("unmarshal version:", shoes)
	if err != nil {
			log.Fatal("Error during Unmarshal(): ", err)
	}
	return shoes, err
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