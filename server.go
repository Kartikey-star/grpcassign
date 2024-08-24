package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/kartikey-star/grpcassign/user"
	"google.golang.org/grpc"
)

type UserServer struct {
	user.UnimplementedUserServiceServer
	users map[string]*user.UserReceipt
}

func (u *UserServer) Create(ctx context.Context, req *user.UserRequest) (*user.UserReceipt, error) {
	receipt := user.UserReceipt{
		User:    req.User,
		From:    req.From,
		To:      req.To,
		Section: req.Section,
		Price:   20,
		Status:  true,
		Seat:    req.Seat,
	}
	u.users[req.GetUser().Emailaddress] = &receipt
	return &receipt, nil
}
func (u *UserServer) Get(ctx context.Context, req *user.ReadUserRequest) (*user.UserReceipt, error) {
	if val, ok := u.users[req.GetUser().Emailaddress]; ok {
		return val, nil
	}
	return &user.UserReceipt{}, fmt.Errorf("User not found")
}

func (u *UserServer) Update(ctx context.Context, req *user.UpdateUserRequest) (*user.UserReceipt, error) {
	if _, ok := u.users[req.GetUser().Emailaddress]; ok {
		updateUserReceipt := user.UserReceipt{
			User:    req.User,
			From:    req.From,
			To:      req.To,
			Section: req.Section,
			Price:   20,
			Status:  true,
			Seat:    req.Seat,
		}
		u.users[req.GetUser().Emailaddress] = &updateUserReceipt
		return &updateUserReceipt, nil
	}
	return &user.UserReceipt{}, fmt.Errorf("User not found")
}

func (u *UserServer) Delete(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserReceipt, error) {
	if _, ok := u.users[req.GetUser().Emailaddress]; ok {
		delete(u.users, req.User.Emailaddress)
		return &user.DeleteUserReceipt{User: req.User, Status: true}, nil
	}
	return &user.DeleteUserReceipt{}, fmt.Errorf("User not found")
}

func (u *UserServer) List(ctx context.Context, req *user.ListUserRequest) (*user.ListUserResponse, error) {
	var users []*user.UserReceipt
	for _, val := range u.users {
		if val.Section == req.Section {
			users = append(users, val)
		}
	}
	listResponse := user.ListUserResponse{
		UserReceiptlist: users,
	}
	return &listResponse, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service := &UserServer{}
	service.users = make(map[string]*user.UserReceipt, 0)
	user.RegisterUserServiceServer(s, service)
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
	}
}
