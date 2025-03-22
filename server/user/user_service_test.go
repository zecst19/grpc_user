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

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()
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

	// Create a new UserService instance
	user_collection := client.Database(dbName).Collection("users_test")

	svc := NewUserService(user_collection)

	// set up test cases
	tests := map[string]struct {
		in                *pb.CreateUserRequest
		expected_response *pb.User
	}{
		"Create User Success": {
			in: &pb.CreateUserRequest{
				FirstName: "Cristiano",
				LastName:  "Ronaldo",
				Nickname:  "CR7",
				Password:  "pass1234",
				Email:     "cristiano@ronaldo.com",
				Country:   "PT",
			},
			expected_response: &pb.User{
				FirstName: "Cristiano",
				LastName:  "Ronaldo",
				Nickname:  "CR7",
				Email:     "cristiano@ronaldo.com",
				Country:   "PT",
			},
		},
	}

	for scenario, tc := range tests {
		t.Run(scenario, func(t *testing.T) {
			resp, err := svc.CreateUser(context.Background(), tc.in)
			require.NoError(t, err)
			require.Equal(t, tc.expected_response.FirstName, resp.FirstName)
			require.Equal(t, tc.expected_response.LastName, resp.LastName)
			require.Equal(t, tc.expected_response.Nickname, resp.Nickname)
			require.Equal(t, tc.expected_response.Email, resp.Email)
			require.Equal(t, tc.expected_response.Country, resp.Country)
		})

	}
}

func TestUserService_GetUser(t *testing.T) {
	ctx := context.Background()
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

	// Create a new UserService instance
	user_collection := client.Database(dbName).Collection("users_test")

	svc := NewUserService(user_collection)

	user_collection.InsertOne(ctx, &pb.User{
		Id:        "86f9f466-851a-4b93-af21-d5f52ac91006",
		FirstName: "Cristiano",
		LastName:  "Ronaldo",
		Nickname:  "CR7",
		Password:  "$2a$14$ol2a598AAUFAV3SizmaZJuvnBjfqA6CGzssNnCQ0Wn9Sr7vFxy5wy",
		Email:     "cristiano@ronaldo.com",
		Country:   "PT",
		CreatedAt: "2025-03-22T18:37:00Z",
		UpdatedAt: "2025-03-22T18:37:00Z",
	},
	)
	// set up test cases
	tests := map[string]struct {
		in                *pb.GetUserRequest
		expected_response *pb.User
	}{
		"Create User Success": {
			in: &pb.GetUserRequest{
				Id: "86f9f466-851a-4b93-af21-d5f52ac91006",
			},
			expected_response: &pb.User{
				Id:        "86f9f466-851a-4b93-af21-d5f52ac91006",
				FirstName: "Cristiano",
				LastName:  "Ronaldo",
				Nickname:  "CR7",
				Password:  "$2a$14$ol2a598AAUFAV3SizmaZJuvnBjfqA6CGzssNnCQ0Wn9Sr7vFxy5wy",
				Email:     "cristiano@ronaldo.com",
				Country:   "PT",
				CreatedAt: "2025-03-22T18:37:00Z",
				UpdatedAt: "2025-03-22T18:37:00Z",
			},
		},
	}

	for scenario, tc := range tests {
		t.Run(scenario, func(t *testing.T) {
			resp, err := svc.GetUser(context.Background(), tc.in)
			require.NoError(t, err)
			require.Equal(t, tc.expected_response, resp)
		})

	}
}

func TestUserService_UpdateUser(t *testing.T) {
	ctx := context.Background()
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

	// Create a new UserService instance
	user_collection := client.Database(dbName).Collection("users_test")

	svc := NewUserService(user_collection)

	user_collection.InsertOne(ctx, &pb.User{
		Id:        "86f9f466-851a-4b93-af21-d5f52ac91006",
		FirstName: "Cristiano",
		LastName:  "Ronaldo",
		Nickname:  "CR7",
		Password:  "$2a$14$ol2a598AAUFAV3SizmaZJuvnBjfqA6CGzssNnCQ0Wn9Sr7vFxy5wy",
		Email:     "cristiano@ronaldo.com",
		Country:   "PT",
		CreatedAt: "2025-03-22T18:37:00Z",
		UpdatedAt: "2025-03-22T18:37:00Z",
	},
	)

	newFirstName := "Mohammed"
	newLastName := "Salah"
	newNickname := "MoSalah"
	newCountry := "EG"

	// set up test cases
	tests := map[string]struct {
		in                *pb.UpdateUserRequest
		expected_response *pb.User
	}{
		"Create User Success": {
			in: &pb.UpdateUserRequest{
				Id:        "86f9f466-851a-4b93-af21-d5f52ac91006",
				FirstName: &newFirstName,
				LastName:  &newLastName,
				Nickname:  &newNickname,
				Country:   &newCountry,
			},
			expected_response: &pb.User{
				Id:        "86f9f466-851a-4b93-af21-d5f52ac91006",
				FirstName: "Mohammed",
				LastName:  "Salah",
				Nickname:  "MoSalah",
				Email:     "cristiano@ronaldo.com",
				Country:   "EG",
				CreatedAt: "2025-03-22T18:37:00Z",
				UpdatedAt: "2025-03-22T18:37:00Z",
			},
		},
	}

	for scenario, tc := range tests {
		t.Run(scenario, func(t *testing.T) {
			resp, err := svc.UpdateUser(context.Background(), tc.in)
			require.NoError(t, err)
			require.Equal(t, tc.expected_response.FirstName, resp.FirstName)
			require.Equal(t, tc.expected_response.LastName, resp.LastName)
			require.Equal(t, tc.expected_response.Nickname, resp.Nickname)
			require.Equal(t, tc.expected_response.Email, resp.Email)
			require.Equal(t, tc.expected_response.Country, resp.Country)
			require.Equal(t, tc.expected_response.CreatedAt, resp.CreatedAt)
			require.NotEqual(t, tc.expected_response.UpdatedAt, resp.UpdatedAt)
		})

	}
}

func TestUserService_DeleteUser(t *testing.T) {
	ctx := context.Background()
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

	// Create a new UserService instance
	user_collection := client.Database(dbName).Collection("users_test")

	svc := NewUserService(user_collection)

	user_collection.InsertOne(ctx, &pb.User{
		Id:        "86f9f466-851a-4b93-af21-d5f52ac91006",
		FirstName: "Cristiano",
		LastName:  "Ronaldo",
		Nickname:  "CR7",
		Password:  "$2a$14$ol2a598AAUFAV3SizmaZJuvnBjfqA6CGzssNnCQ0Wn9Sr7vFxy5wy",
		Email:     "cristiano@ronaldo.com",
		Country:   "PT",
		CreatedAt: "2025-03-22T18:37:00Z",
		UpdatedAt: "2025-03-22T18:37:00Z",
	},
	)

	// set up test cases
	tests := map[string]struct {
		in                *pb.DeleteUserRequest
		expected_response *pb.DeleteUserResponse
	}{
		"Create User Success": {
			in: &pb.DeleteUserRequest{
				Id: "86f9f466-851a-4b93-af21-d5f52ac91006",
			},
			expected_response: &pb.DeleteUserResponse{
				Success: true,
			},
		},
	}

	for scenario, tc := range tests {
		t.Run(scenario, func(t *testing.T) {
			resp, err := svc.DeleteUser(context.Background(), tc.in)
			require.NoError(t, err)
			require.Equal(t, tc.expected_response, resp)
		})

	}
}

func TestUserService_ListUsers(t *testing.T) {
	ctx := context.Background()
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

	// Create a new UserService instance
	user_collection := client.Database(dbName).Collection("users_test")

	svc := NewUserService(user_collection)

	user_collection.InsertOne(ctx,
		&pb.User{
			Id:        "86f9f466-851a-4b93-af21-d5f52ac91006",
			FirstName: "Cristiano",
			LastName:  "Ronaldo",
			Nickname:  "CR7",
			Password:  "$2a$14$ol2a598AAUFAV3SizmaZJuvnBjfqA6CGzssNnCQ0Wn9Sr7vFxy5wy",
			Email:     "cristiano@ronaldo.com",
			Country:   "PT",
			CreatedAt: "2025-03-22T18:37:00Z",
			UpdatedAt: "2025-03-22T18:37:00Z",
		},
	)
	user_collection.InsertOne(ctx,
		&pb.User{
			Id:        "4bedafd6-b946-4d70-a156-d828ddcb62fb",
			FirstName: "Mohammed",
			LastName:  "Salah",
			Nickname:  "MoSalah",
			Password:  "$2a$14$ol2a598AAUFAV3SizmaZJuvnBjfqA6CGzssNnCQ0Wn9Sr7vFxy5wy",
			Email:     "mo@salah.com",
			Country:   "EG",
			CreatedAt: "2025-03-22T18:37:00Z",
			UpdatedAt: "2025-03-22T18:37:00Z",
		},
	)
	countryFilter := "PT"
	// set up test cases
	tests := map[string]struct {
		in                *pb.ListUsersRequest
		expected_response *pb.ListUsersResponse
	}{
		"List User Success": {
			in: &pb.ListUsersRequest{
				Page:     1,
				PageSize: 10,
			},
			expected_response: &pb.ListUsersResponse{
				Users: []*pb.User{
					&pb.User{
						Id:        "86f9f466-851a-4b93-af21-d5f52ac91006",
						FirstName: "Cristiano",
						LastName:  "Ronaldo",
						Nickname:  "CR7",
						Password:  "$2a$14$ol2a598AAUFAV3SizmaZJuvnBjfqA6CGzssNnCQ0Wn9Sr7vFxy5wy",
						Email:     "cristiano@ronaldo.com",
						Country:   "PT",
						CreatedAt: "2025-03-22T18:37:00Z",
						UpdatedAt: "2025-03-22T18:37:00Z",
					},
					&pb.User{
						Id:        "4bedafd6-b946-4d70-a156-d828ddcb62fb",
						FirstName: "Mohammed",
						LastName:  "Salah",
						Nickname:  "MoSalah",
						Password:  "$2a$14$ol2a598AAUFAV3SizmaZJuvnBjfqA6CGzssNnCQ0Wn9Sr7vFxy5wy",
						Email:     "mo@salah.com",
						Country:   "EG",
						CreatedAt: "2025-03-22T18:37:00Z",
						UpdatedAt: "2025-03-22T18:37:00Z",
					},
				},
				TotalCount: 2,
			},
		},
		"List User Success Filter Country": {
			in: &pb.ListUsersRequest{
				Page:     1,
				PageSize: 10,
				Country:  &countryFilter,
			},
			expected_response: &pb.ListUsersResponse{
				Users: []*pb.User{
					&pb.User{
						Id:        "86f9f466-851a-4b93-af21-d5f52ac91006",
						FirstName: "Cristiano",
						LastName:  "Ronaldo",
						Nickname:  "CR7",
						Password:  "$2a$14$ol2a598AAUFAV3SizmaZJuvnBjfqA6CGzssNnCQ0Wn9Sr7vFxy5wy",
						Email:     "cristiano@ronaldo.com",
						Country:   "PT",
						CreatedAt: "2025-03-22T18:37:00Z",
						UpdatedAt: "2025-03-22T18:37:00Z",
					},
				},
				TotalCount: 1,
			},
		},
	}

	for scenario, tc := range tests {
		t.Run(scenario, func(t *testing.T) {
			resp, err := svc.ListUsers(context.Background(), tc.in)
			require.NoError(t, err)
			require.Equal(t, tc.expected_response, resp)
		})

	}
}
