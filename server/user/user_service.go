package grpc_user

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	pb "github.com/zecst19/grpc-user/proto"
)

var (
	topic = "user-topic"
)

type Message struct {
	Event string `json:"event"`
	Value any    `json:"value"`
}

type UserService struct {
	pb.UnimplementedUserServiceServer
	collection    *mongo.Collection
	kafkaProducer sarama.SyncProducer
}

func NewUserService(collection *mongo.Collection, kafkaProducer sarama.SyncProducer) *UserService {
	return &UserService{collection: collection, kafkaProducer: kafkaProducer}
}

func (svc *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	//TODO add checks
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password: %v", err)
	}

	user := &pb.User{
		Id:        uuid.New().String(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Password:  string(hashedPassword),
		Email:     req.Email,
		Country:   req.Country,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	_, err = svc.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create user: %v", err)
	}

	// Define the message to be sent
	message := Message{
		Event: "user.created",
		Value: user,
	}

	serializedMessage, err := json.Marshal(message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to serialize Producer message: %v", err)
	}

	err = svc.produceMessage(serializedMessage)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to send Producer message: %v", err)
	}

	return user, nil
}

func (svc *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	var user pb.User
	err := svc.collection.FindOne(ctx, bson.M{"id": req.Id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get user: %v", err)
	}

	// Define the message to be sent
	message := Message{
		Event: "user.get",
		Value: &user,
	}

	serializedMessage, err := json.Marshal(message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to serialize Producer message: %v", err)
	}

	err = svc.produceMessage(serializedMessage)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to send Producer message: %v", err)
	}

	return &user, nil
}

func (svc *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	//call the get and fill these variables with it
	var user pb.User
	err := svc.collection.FindOne(ctx, bson.M{"id": req.Id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get user: %v", err)
	}

	updatedFirstName := user.FirstName
	if req.FirstName != nil {
		updatedFirstName = *req.FirstName
	}

	updatedLastName := user.LastName
	if req.LastName != nil {
		updatedLastName = *req.LastName
	}

	updatedNickname := user.Nickname
	if req.Nickname != nil {
		updatedNickname = *req.Nickname
	}

	updatedEmail := user.Email
	if req.Email != nil {
		updatedEmail = *req.Email
	}

	updatedCountry := user.Country
	if req.Country != nil {
		updatedCountry = *req.Country
	}

	update := bson.M{
		"$set": &pb.User{
			Id:        user.Id,
			FirstName: updatedFirstName,
			LastName:  updatedLastName,
			Nickname:  updatedNickname,
			Password:  user.Password,
			Email:     updatedEmail,
			Country:   updatedCountry,
			CreatedAt: user.CreatedAt,
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
	}

	var updatedUser pb.User
	err = svc.collection.FindOneAndUpdate(
		ctx,
		bson.M{"id": req.Id},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to update user: %v", err)
	}

	// Define the message to be sent
	message := Message{
		Event: "user.update",
		Value: &updatedUser,
	}

	serializedMessage, err := json.Marshal(message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to serialize Producer message: %v", err)
	}

	err = svc.produceMessage(serializedMessage)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to send Producer message: %v", err)
	}

	return &updatedUser, nil
}

func (svc *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	res, err := svc.collection.DeleteOne(ctx, bson.M{"id": req.Id})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete user: %v", err)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	// Define the message to be sent
	message := Message{
		Event: "user.delete",
		Value: req.Id,
	}

	serializedMessage, err := json.Marshal(message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to serialize Producer message: %v", err)
	}

	err = svc.produceMessage(serializedMessage)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to send Producer message: %v", err)
	}

	return &pb.DeleteUserResponse{Success: true}, nil
}

func (svc *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	var users []*pb.User

	opts := options.Find().
		SetSkip(int64((req.Page - 1) * req.PageSize)).
		SetLimit(int64(req.PageSize))

	filter := bson.M{}
	if req.Country != nil {
		filter["country"] = req.Country
	}

	if req.LastName != nil {
		filter["last_name"] = req.LastName
	}

	//if req.FromDate != nil {
	//	filter["created_at"] = "$gte: " + *req.FromDate
	//}

	//if req.ToDate != nil {
	//	filter["created_at"] = "$lte: " + *req.ToDate
	//}

	cursor, err := svc.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list users: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user pb.User
		if err := cursor.Decode(&user); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to decode user: %v", err)
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Cursor error: %v", err)
	}

	totalCount, err := svc.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to count users: %v", err)
	}

	// Define the message to be sent
	message := Message{
		Event: "user.list",
		Value: users,
	}

	serializedMessage, err := json.Marshal(message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to serialize Producer message: %v", err)
	}

	err = svc.produceMessage(serializedMessage)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to send Producer message: %v", err)
	}

	return &pb.ListUsersResponse{
		Users:      users,
		TotalCount: int32(totalCount),
	}, nil
}

func (svc *UserService) produceMessage(message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := svc.kafkaProducer.SendMessage(msg)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to produce message: %v", err)
	}
	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

	return nil
}
