package gapi

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/iKayrat/auth-grpc/internal/model"
	"github.com/iKayrat/auth-grpc/internal/services/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.User, error) {

	var u model.User

	u.Username = req.Username
	u.Email = req.Email
	u.Password = req.Password
	u.Admin = req.Admin

	user, err := server.Store.Create(ctx, u)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	idd := user.ID.String()
	fmt.Println(idd)
	resp := &pb.User{
		Id:       idd,
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
		Admin:    user.Admin,
	}

	// resp.Id = idd
	// resp.Username = user.Username
	// resp.Email = user.Email
	// resp.Password = user.Password
	// resp.Admin = user.Admin

	return resp, nil

}

func (server *Server) GetUserById(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	user, exists := server.Store.GetById(ctx, id)
	if !exists {
		return nil, errors.New("user not exists")
	}

	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, err.Error())
	// }

	idd := user.ID.String()

	resp := &pb.User{
		Id:       idd,
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
		Admin:    user.Admin,
	}

	return resp, nil

}

func (server *Server) GetUsers(ctx context.Context, req *pb.Empty) (*pb.GetUsersResp, error) {

	users := server.Store.GetAll(ctx)

	// resp := pb.GetUsersResp{
	// 	Users: make([]*pb.User, len(users)),
	// }
	resp := new(pb.GetUsersResp)

	log.Println("resp:", resp)
	log.Println("users:", len(users))
	log.Println(len(resp.Users))

	pbUsers := make([]*pb.User, len(users))
	log.Printf("resp:%#v\n", pbUsers)
	log.Printf("resp:%#v\n", users)

	for i := 0; i < len(users); i++ {

		idd := users[i].ID.String()

		pbUsers[i].Id = idd
		pbUsers[i].Username = users[i].Username
		pbUsers[i].Email = users[i].Email
		pbUsers[i].Password = users[i].Password
		pbUsers[i].Admin = users[i].Admin
	}
	log.Println("resp:", resp)

	resp.Users = append(resp.Users, pbUsers...)

	// for i := 0; i < len(users); i++ {

	// 	resp.Users[i].Id = pbUsers[i].Id
	// 	resp.Users[i].Username = pbUsers[i].Username
	// 	resp.Users[i].Email = pbUsers[i].Email
	// 	resp.Users[i].Password = pbUsers[i].Password
	// 	resp.Users[i].Admin = pbUsers[i].Admin
	// }

	return resp, nil
}
