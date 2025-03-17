package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/zecst19/grpc-user/proto"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a new user
	createResp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		FirstName: "John",
		LastName:  "Cena",
		Nickname:  "BuenosDiasMatosinhos",
		Password:  "password1234", //TODO HASH
		Email:     "john.cena@wwe.com",
		Country:   "US",
	})
	if err != nil {
		log.Fatalf("Could not create user: %v", err)
	}
	fmt.Printf("Created user: %v\n", createResp)

	// Get the created user
	getResp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: createResp.Id})
	if err != nil {
		log.Fatalf("Could not get user: %v", err)
	}
	fmt.Printf("Got user: %v\n", getResp)

	newFirstName := "Cristiano"
	newLastName := "Ronaldo"
	newCountry := "PT"
	// Update the user
	updateResp, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:        getResp.Id,
		FirstName: &newFirstName,
		LastName:  &newLastName,
		Country:   &newCountry,
	})
	if err != nil {
		log.Fatalf("Could not update user: %v", err)
	}
	fmt.Printf("Updated user: %v\n", updateResp)

	// Get the udpated user
	getUpdatedResp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: createResp.Id})
	if err != nil {
		log.Fatalf("Could not get user: %v", err)
	}
	fmt.Printf("Got updated user: %v\n", getUpdatedResp)

	// List users
	listResp, err := client.ListUsers(ctx, &pb.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Fatalf("Could not list users: %v", err)
	}
	fmt.Printf("Listed %d users, total count: %d\n", len(listResp.Users), listResp.TotalCount)
	for _, user := range listResp.Users {
		fmt.Printf("- %s: %s\n", user.Id, user.Nickname)
	}

	// Delete the user
	deleteResp, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: createResp.Id})
	if err != nil {
		log.Fatalf("Could not delete user: %v", err)
	}
	fmt.Printf("Deleted user, success: %v\n", deleteResp.Success)
}
