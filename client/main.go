package main

import (
	"log"
	pb "github.com/karankumarshreds/GoProto/protofiles"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// grpc server address
const address = "localhost:8000"

func main() {
	// Set up connection with the grpc server 
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error while making connection, %v", err)
	}

	// Create a client instance 
	c := pb.NewMoneyTransactionClient(conn)
	
	// Lets invoke the remote function from client on the server 
	c.MakeTransaction(
		context.Background(), 
		&pb.TransactionRequest{
			From: "John",
			To: "Alice",
			Amount: float32(120.15),
		},
	)
}