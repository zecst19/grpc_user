package main

import (
	"context"
	"log"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	pb "github.com/zecst19/grpc-user/proto"
	userService "github.com/zecst19/grpc-user/server/user"
)

const (
	port     = ":50051"
	mongoURI = "mongodb://localhost:27017"
	dbName   = "userDB"
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
	server := grpc.NewServer()

	// Server Health Check
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(server, healthcheck)

	go func() {
		// asynchronously inspect dependencies and toggle serving status as needed
		next := healthpb.HealthCheckResponse_SERVING

		for {
			healthcheck.SetServingStatus("", next)

			if next == healthpb.HealthCheckResponse_SERVING {
				next = healthpb.HealthCheckResponse_NOT_SERVING
			} else {
				next = healthpb.HealthCheckResponse_SERVING
			}

			log.Print("Health:", next.Descriptor().Name())
			time.Sleep(time.Second * 5)
		}
	}()

	// Register our service with the gRPC server
	pb.RegisterUserServiceServer(server, user_service)

	// Start listening on the specified port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server listening on %s", port)

	// Start serving
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
