package handler

import (
	"context"
	"time"

	"github.com/jram17/second-brain/services/auth/internal/model"
	"github.com/jram17/second-brain/services/auth/pkg/jwt"
	pb "github.com/jram17/second-brain/services/auth/pkg/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	store    *model.Store
	jwtMaker *jwt.Maker
}

// make the constructor
func NewAuthHandler(store *model.Store, jwtMaker *jwt.Maker) *AuthHandler {
	return &AuthHandler{
		store:    store,
		jwtMaker: jwtMaker,
	}
}

// make the signup
func (a *AuthHandler) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	//check if the user exists
	_, err := a.store.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}
	//hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}
	//store the user now
	user, err := a.store.CreateUser(ctx, req.Email, req.Username, string(hashed))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}
	//make the tokens
	refresh, err := a.jwtMaker.GenerateToken(user.ID.Hex(), req.Email, time.Hour*7*24)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}
	access, err := a.jwtMaker.GenerateToken(user.ID.Hex(), req.Email, time.Minute*15)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}
	return &pb.SignupResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// make the login
func (a *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// check if the user exist!!
	user, err := a.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	//now that user exists check the password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password")
	}
	//no err then user is authenticated!!!
	refresh, err := a.jwtMaker.GenerateToken(user.ID.Hex(), req.Email, time.Hour*7*24)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}
	access, err := a.jwtMaker.GenerateToken(user.ID.Hex(), req.Email, time.Minute*15)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}
	return &pb.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// make the validatetoken
func (a *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := a.jwtMaker.ValidateToken(req.AccessToken)
	if err != nil {
		return &pb.ValidateResponse{Valid: false}, nil
	}
	return &pb.ValidateResponse{
		Userid: claims.UserId,
		Valid:  true,
	}, nil
}

// make the refreshToken
func (a *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	claims, err := a.jwtMaker.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token")
	}
	access, err := a.jwtMaker.GenerateToken(claims.UserId, claims.Email, time.Minute*15)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}
	return &pb.RefreshTokenResponse{
		AccessToken:  access,
		RefreshToken: req.RefreshToken,
	}, nil
}
