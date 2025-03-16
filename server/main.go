package main

import (
	"context"
	"log"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	pb "github.com/zecst19/grpc-user/proto"
	userService "github.com/zecst19/grpc-user/server/user/user_service"
)

const (
	port     = ":50051"
	mongoURI = "mongodb://localhost:27017"
	dbName   = "taskdb"
)

func main() {
	// Set up a connection to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")

	// Create a new UserService instance
	user_collection := client.Database(dbName).Collection("users")
	user_service := userService.NewUserService(user_collection)

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register our service with the gRPC server
	pb.RegisterUserServiceServer(s, user_service)

	// Start listening on the specified port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server listening on %s", port)

	// Start serving
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
