package gapi

import (
	storage "github.com/iKayrat/auth-grpc/internal/inmemory"
	"github.com/iKayrat/auth-grpc/internal/services/pb"
)

type Server struct {
	pb.ProfileServiceServer
	Store storage.Store
}

func NewGrpcServer(store *storage.Store) (*Server, error) {
	server := &Server{
		Store: *store,
	}

	return server, nil
}
