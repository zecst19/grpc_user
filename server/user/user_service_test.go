package grpc_user

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	pb "github.com/zecst19/grpc-user/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	port     = ":50051"
	mongoURI = "mongodb://localhost:27017"
	dbName   = "userTestDB"
)

func setup() *mongo.Collection {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Create a new UserService instance
	user_collection := client.Database(dbName).Collection("users")

	return user_collection
}

func TestUserService_CreateUser(t *testing.T) {
	svc := NewUserService(setup())

	// set up test cases
	tests := []struct {
		name              string
		request           pb.CreateUserRequest
		expected_response pb.User
	}{
		{
			name: "Create User Success",
			request: pb.CreateUserRequest{
				FirstName: "Cristiano",
				LastName:  "Ronaldo",
				Nickname:  "CR7",
				Password:  "pass1234",
				Email:     "cristiano@ronaldo.com",
				Country:   "PT",
			},
			expected_response: pb.User{
				Id:        "1",
				FirstName: "Cristiano",
				LastName:  "Ronaldo",
				Nickname:  "CR7",
				Password:  "pass1234",
				Email:     "cristiano@ronaldo.com",
				Country:   "PT",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := svc.CreateUser(context.Background(), &tc.request)
			require.NoError(t, err)
			require.Equal(t, tc.expected_response, resp)
		})

	}
}

func TestUserService_GetUser(t *testing.T) {

}

func TestUserService_UpdateUser(t *testing.T) {

}

func TestUserService_DeleteUser(t *testing.T) {

}

func TestUserService_ListUsers(t *testing.T) {

}
