package main

import (
	"log"
	"net"

	"github.com/iKayrat/auth-grpc/internal/config"
	storage "github.com/iKayrat/auth-grpc/internal/inmemory"
	gapi "github.com/iKayrat/auth-grpc/internal/servers"
	"github.com/iKayrat/auth-grpc/internal/services/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	Userstorage := storage.InitStorage()
	// var store storage.Store = Userstorage
	store := storage.New(Userstorage)

	server, err := gapi.NewGrpcServer(&store)
	if err != nil {
		log.Fatal("Error creating server", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(server.AuthInterceptor),
	)

	pb.RegisterProfileServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", conf.ServerAddress)
	if err != nil {
		log.Fatal("Error creating server", err)
	}

	log.Printf("gRPC server listening on %s", conf.ServerAddress)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Error serving gRPC server", err)
	}
}
