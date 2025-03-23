package grpc_user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	pb "github.com/zecst19/grpc-user/proto"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	collection *mongo.Collection
}

func NewUserService(collection *mongo.Collection) *UserService {
	return &UserService{collection: collection}
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

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	var user pb.User
	err := s.collection.FindOne(ctx, bson.M{"id": req.Id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get user: %v", err)
	}

	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	//call the get and fill these variables with it
	var user pb.User
	err := s.collection.FindOne(ctx, bson.M{"id": req.Id}).Decode(&user)
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
	err = s.collection.FindOneAndUpdate(
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

	return &updatedUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	res, err := s.collection.DeleteOne(ctx, bson.M{"id": req.Id})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete user: %v", err)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	return &pb.DeleteUserResponse{Success: true}, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
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

	cursor, err := s.collection.Find(ctx, filter, opts)
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

	totalCount, err := s.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to count users: %v", err)
	}

	return &pb.ListUsersResponse{
		Users:      users,
		TotalCount: int32(totalCount),
	}, nil
}
