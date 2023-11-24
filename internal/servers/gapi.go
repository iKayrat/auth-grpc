package gapi

import (
	storage "github.com/iKayrat/auth-grpc/internal/inmemory"
	"github.com/iKayrat/auth-grpc/internal/services/pb"
)

type Server struct {
	pb.ProfileServiceServer
	Store storage.Store

	// authService services.AuthService
	// userService services.UserService
}

func NewGrpcServer(store storage.Store) (*Server, error) {
	server := &Server{
		Store: store,
	}

	return server, nil
}
