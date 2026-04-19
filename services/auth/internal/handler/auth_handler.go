package handler

import (
	"context"

	"github.com/jram17/second-brain/services/auth/internal/model"
	"github.com/jram17/second-brain/services/auth/pkg/jwt"
	pb "github.com/jram17/second-brain/services/auth/pkg/pb"
)

type AuthHandler struct{
	pb.UnimplementedAuthServiceServer
	store *model.Store
	jwtMaker *jwt.Maker
}

//make the constructor
func NewAuthHandler(store *model.Store, jwtMaker *jwt.Maker) *AuthHandler {
	return &AuthHandler{
		store: store,
		jwtMaker: jwtMaker,
	}
}



