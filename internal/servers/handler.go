package gapi

import (
	"context"
	"errors"
	"fmt"

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

	users, err := server.Store.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := new(pb.GetUsersResp)

	for i := 0; i < len(users); i++ {

		idd := users[i].ID.String()

		resp.Users = append(resp.Users, &pb.User{
			Id:       idd,
			Email:    users[i].Email,
			Username: users[i].Username,
			Password: users[i].Password,
		})
	}

	return resp, nil
}

func (server *Server) UpdateUser(ctx context.Context, in *pb.GetByIdRequest) (*pb.User, error) {

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, err
	}
	var u model.User = model.User{
		ID:       id,
		Email:    in.Email,
		Username: in.Username,
		Password: in.Password,
		Admin:    in.Admin,
	}

	updatedUser, err := server.Store.Update(ctx, u)
	if err != nil {
		return nil, errors.New("user not exists")
	}

	resp := &pb.User{
		Id:       updatedUser.ID.String(),
		Email:    updatedUser.Email,
		Username: updatedUser.Username,
		Password: updatedUser.Password,
		Admin:    updatedUser.Admin,
	}

	return resp, nil

}
func (server *Server) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.ResponseMsg, error) {
	err := server.Store.Delete(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("User with (%s) ID deleted", in.Id)

	resp := new(pb.ResponseMsg)
	resp.Msg = msg

	return resp, nil
}
