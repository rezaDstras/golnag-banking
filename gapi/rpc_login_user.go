package gapi

import (
	"context"
	"database/sql"
	db "github.com/rezaDastrs/banking/db/sqlc"
	"github.com/rezaDastrs/banking/pb"
	"github.com/rezaDastrs/banking/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser( ctx context.Context,req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	err = util.ChaeckPassword(req.GetPassword(), user.HashPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	accessToken,accessPayload , err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token ")
	}

	refreshToken,refreshPayload , err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refrsh token")
	}

	session , err :=server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpireAt:   refreshPayload.ExpireAt  ,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	rep := &pb.LoginUserResponse{
		SessionId: session.ID.String(),
		AccessToken: accessToken,
		AccessTokenExpireAt: timestamppb.New(accessPayload.ExpireAt),
		RefreshToken: refreshToken,
		RefreshTokenExpireAt: timestamppb.New(refreshPayload.ExpireAt),
		User:        convertUser(user),
	}

	return rep , nil
}
