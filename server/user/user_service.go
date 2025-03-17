package user_service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

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
	user := &pb.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Password:  req.Password, //TODO HASH
		Email:     req.Email,
		Country:   req.Country,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	res, err := svc.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create user: %v", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to convert InsertedID to ObjectID")
	}

	user.Id = oid.Hex() //TODO chcek if this is hashed
	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	oid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid ID")
	}

	var user pb.User
	err = s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get user: %v", err)
	}

	user.Id = oid.Hex()
	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	oid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid ID")
	}

	//call the get and fill these variables with it
	var user pb.User
	err = s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
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

	var updatedNickname string
	if req.Nickname != nil {
		updatedNickname = *req.Nickname
	}

	//maybe remove this from update
	var updatedPassword string
	if req.Password != nil {
		updatedPassword = *req.Password
	}

	var updatedEmail string
	if req.Email != nil {
		updatedEmail = *req.Email
	}

	var updatedCountry string
	if req.Country != nil {
		updatedCountry = *req.Country
	}

	update := bson.M{
		"$set": bson.M{
			"id":         user.Id,
			"first_name": updatedFirstName,
			"last_name":  updatedLastName,
			"nickname":   updatedNickname,
			"password":   updatedPassword,
			"email":      updatedEmail,
			"country":    updatedCountry,
			"created_at": user.CreatedAt,
			"updated_at": time.Now(),
		},
	}

	var updatedUser pb.User
	err = s.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": oid},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to update user: %v", err)
	}

	updatedUser.Id = oid.Hex()
	return &updatedUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	oid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid ID")
	}

	res, err := s.collection.DeleteOne(ctx, bson.M{"_id": oid})
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
	//missing the filters
	opts := options.Find().
		SetSkip(int64((req.Page - 1) * req.PageSize)).
		SetLimit(int64(req.PageSize))

	cursor, err := s.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list users: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user pb.User
		if err := cursor.Decode(&user); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to decode user: %v", err)
		}
		user.Id = user.Id // Convert ObjectID to string
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
