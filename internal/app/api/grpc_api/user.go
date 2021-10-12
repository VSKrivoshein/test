package grpc_api

import (
	context "context"
	api "github.com/VSKrivoshein/test/internal/app/api/grpc_api/proto"
	"github.com/VSKrivoshein/test/internal/app/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	Service service.User
	api.UnimplementedUserServer
}

func NewGRPCServer(service service.User) api.UserServer {
	return &GRPCServer{Service: service}
}

func (s *GRPCServer) Create(ctx context.Context, request *api.CreateRequest) (*emptypb.Empty, error) {
	if err := s.Service.Create(request.GetEmail(), request.GetPassword()); err != nil {
		return nil, err
	}
	return new(emptypb.Empty), nil
}

func (s *GRPCServer) Delete(ctx context.Context, request *api.DeleteRequest) (*emptypb.Empty, error) {
	if err := s.Service.Delete(request.GetEmail(), request.GetPassword()); err != nil {
		return nil, err
	}
	return new(emptypb.Empty), nil
}

func (s *GRPCServer) GetAll(ctx context.Context, empty *emptypb.Empty) (*api.GetAllResponse, error) {
	userList, err := s.Service.GetAll()
	if err != nil {
		return nil, err
	}
	return &api.GetAllResponse{UserList: userList}, nil
}
