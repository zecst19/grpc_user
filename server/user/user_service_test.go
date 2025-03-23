package grpc_user

import (
	"context"
	"log"
	"testing"

	"github.com/IBM/sarama"
	sarama_mock "github.com/IBM/sarama/mocks"
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

func TestUserService(t *testing.T) {
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
	mock_producer := sarama_mock.NewSyncProducer(t, sarama.NewConfig())
	svc := NewUserService(user_collection, mock_producer)

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

	t.Run("Create User", func(t *testing.T) {
		in := &pb.CreateUserRequest{
			FirstName: "Mohammed",
			LastName:  "Salah",
			Nickname:  "MoSalah",
			Password:  "word5678",
			Email:     "mo@salah.com",
			Country:   "EG",
		}

		expected_response := &pb.User{
			FirstName: "Mohammed",
			LastName:  "Salah",
			Nickname:  "MoSalah",
			Email:     "mo@salah.com",
			Country:   "EG",
		}

		mock_producer.ExpectSendMessageAndSucceed()

		resp, err := svc.CreateUser(context.Background(), in)
		require.NoError(t, err)
		require.Equal(t, expected_response.FirstName, resp.FirstName)
		require.Equal(t, expected_response.LastName, resp.LastName)
		require.Equal(t, expected_response.Nickname, resp.Nickname)
		require.Equal(t, expected_response.Email, resp.Email)
		require.Equal(t, expected_response.Country, resp.Country)
	})

	t.Run("Get User", func(t *testing.T) {
		in := &pb.GetUserRequest{
			Id: "86f9f466-851a-4b93-af21-d5f52ac91006",
		}

		expected_response := &pb.User{
			Id:        "86f9f466-851a-4b93-af21-d5f52ac91006",
			FirstName: "Cristiano",
			LastName:  "Ronaldo",
			Nickname:  "CR7",
			Password:  "$2a$14$ol2a598AAUFAV3SizmaZJuvnBjfqA6CGzssNnCQ0Wn9Sr7vFxy5wy",
			Email:     "cristiano@ronaldo.com",
			Country:   "PT",
			CreatedAt: "2025-03-22T18:37:00Z",
			UpdatedAt: "2025-03-22T18:37:00Z",
		}

		mock_producer.ExpectSendMessageAndSucceed()

		resp, err := svc.GetUser(context.Background(), in)
		require.NoError(t, err)
		require.Equal(t, expected_response, resp)

	})

	t.Run("Update User", func(t *testing.T) {
		newFirstName := "Lionel"
		newLastName := "Messi"
		newNickname := "Messi"
		newCountry := "AR"

		in := &pb.UpdateUserRequest{
			Id:        "86f9f466-851a-4b93-af21-d5f52ac91006",
			FirstName: &newFirstName,
			LastName:  &newLastName,
			Nickname:  &newNickname,
			Country:   &newCountry,
		}

		expected_response := &pb.User{
			Id:        "86f9f466-851a-4b93-af21-d5f52ac91006",
			FirstName: "Lionel",
			LastName:  "Messi",
			Nickname:  "Messi",
			Email:     "cristiano@ronaldo.com",
			Country:   "AR",
			CreatedAt: "2025-03-22T18:37:00Z",
			UpdatedAt: "2025-03-22T18:37:00Z",
		}

		mock_producer.ExpectSendMessageAndSucceed()

		resp, err := svc.UpdateUser(context.Background(), in)
		require.NoError(t, err)
		require.Equal(t, expected_response.FirstName, resp.FirstName)
		require.Equal(t, expected_response.LastName, resp.LastName)
		require.Equal(t, expected_response.Nickname, resp.Nickname)
		require.Equal(t, expected_response.Email, resp.Email)
		require.Equal(t, expected_response.Country, resp.Country)
		require.Equal(t, expected_response.CreatedAt, resp.CreatedAt)
		require.NotEqual(t, expected_response.UpdatedAt, resp.UpdatedAt)

	})

	t.Run("List Users", func(t *testing.T) {
		in := &pb.ListUsersRequest{
			Page:     1,
			PageSize: 10,
		}

		expected_response := &pb.ListUsersResponse{
			Users: []*pb.User{
				{
					Id:        "86f9f466-851a-4b93-af21-d5f52ac91006",
					FirstName: "Lionel",
					LastName:  "Messi",
					Nickname:  "Messi",
					Password:  "$2a$14$ol2a598AAUFAV3SizmaZJuvnBjfqA6CGzssNnCQ0Wn9Sr7vFxy5wy",
					Email:     "cristiano@ronaldo.com",
					Country:   "AR",
					CreatedAt: "2025-03-22T18:37:00Z",
					UpdatedAt: "2025-03-22T18:37:00Z",
				},
				{
					Id:        "4bedafd6-b946-4d70-a156-d828ddcb62fb",
					FirstName: "Mohammed",
					LastName:  "Salah",
					Nickname:  "MoSalah",
					Password:  "$2a$14$SpVoNKXfhSqGP4YTkh6H6ORQCoWxN9gQt3HHreNxKl8D0Gp1bn3oa",
					Email:     "mo@salah.com",
					Country:   "EG",
					CreatedAt: "2025-03-22T18:37:00Z",
					UpdatedAt: "2025-03-22T18:37:00Z",
				},
			},
			TotalCount: 2,
		}

		mock_producer.ExpectSendMessageAndSucceed()

		resp, err := svc.ListUsers(context.Background(), in)
		require.NoError(t, err)
		require.Equal(t, len(expected_response.Users), len(resp.Users))

		require.Equal(t, expected_response.Users[0].FirstName, resp.Users[0].FirstName)
		require.Equal(t, expected_response.Users[0].LastName, resp.Users[0].LastName)
		require.Equal(t, expected_response.Users[0].Nickname, resp.Users[0].Nickname)
		require.Equal(t, expected_response.Users[0].Email, resp.Users[0].Email)
		require.Equal(t, expected_response.Users[0].Country, resp.Users[0].Country)

		require.Equal(t, expected_response.Users[1].FirstName, resp.Users[1].FirstName)
		require.Equal(t, expected_response.Users[1].LastName, resp.Users[1].LastName)
		require.Equal(t, expected_response.Users[1].Nickname, resp.Users[1].Nickname)
		require.Equal(t, expected_response.Users[1].Email, resp.Users[1].Email)
		require.Equal(t, expected_response.Users[1].Country, resp.Users[1].Country)

		require.Equal(t, expected_response.TotalCount, int32(2))
	})

	t.Run("List Users Country Filter", func(t *testing.T) {
		countryFilter := "EG"

		in := &pb.ListUsersRequest{
			Page:     1,
			PageSize: 10,
			Country:  &countryFilter,
		}

		expected_response := &pb.ListUsersResponse{
			Users: []*pb.User{
				{
					Id:        "4bedafd6-b946-4d70-a156-d828ddcb62fb",
					FirstName: "Mohammed",
					LastName:  "Salah",
					Nickname:  "MoSalah",
					Password:  "$2a$14$SpVoNKXfhSqGP4YTkh6H6ORQCoWxN9gQt3HHreNxKl8D0Gp1bn3oa",
					Email:     "mo@salah.com",
					Country:   "EG",
					CreatedAt: "2025-03-22T18:37:00Z",
					UpdatedAt: "2025-03-22T18:37:00Z",
				},
			},
			TotalCount: 2,
		}

		mock_producer.ExpectSendMessageAndSucceed()

		resp, err := svc.ListUsers(context.Background(), in)
		require.NoError(t, err)
		require.Equal(t, len(expected_response.Users), len(resp.Users))

		require.Equal(t, expected_response.Users[0].FirstName, resp.Users[0].FirstName)
		require.Equal(t, expected_response.Users[0].LastName, resp.Users[0].LastName)
		require.Equal(t, expected_response.Users[0].Nickname, resp.Users[0].Nickname)
		require.Equal(t, expected_response.Users[0].Email, resp.Users[0].Email)
		require.Equal(t, expected_response.Users[0].Country, resp.Users[0].Country)

		require.Equal(t, expected_response.TotalCount, int32(2))
	})

	t.Run("Delete User", func(t *testing.T) {
		in := &pb.DeleteUserRequest{
			Id: "86f9f466-851a-4b93-af21-d5f52ac91006",
		}

		expected_response := &pb.DeleteUserResponse{
			Success: true,
		}

		mock_producer.ExpectSendMessageAndSucceed()

		resp, err := svc.DeleteUser(context.Background(), in)
		require.NoError(t, err)
		require.Equal(t, expected_response, resp)
	})

	client.Database(dbName).Collection("users_test").Drop(ctx)
}
