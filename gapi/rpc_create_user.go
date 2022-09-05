package gapi

import (
	"context"
	"github.com/lib/pq"
	db "github.com/rezaDastrs/banking/db/sqlc"
	"github.com/rezaDastrs/banking/pb"
	"github.com/rezaDastrs/banking/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context,req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "faild to hash password: %s",err)
	}

	arg := db.CreateUserParams{
		Username:     req.GetUsername(),
		HashPassword: hashPassword,
		FullName:     req.GetFullName(),
		Email:        req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_validation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s",err)
			}
		}
		return nil, status.Errorf(codes.Internal, "faild to create user: %s",err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp , nil
}
