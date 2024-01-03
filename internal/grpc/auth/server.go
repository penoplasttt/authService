package auth

import (
	"context"

	pb "github.com/penoplasttt/authService/protos/gen/go"
	"google.golang.org/grpc"
)

type serverAPI struct {
	pb.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	pb.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login (ctx context.Context, req *pb.LoginRequest)(*pb.LoginResponse,error){
	return &pb.LoginResponse{Token: "123"}, nil
	
}

func (s *serverAPI) Register (ctx context.Context, req *pb.RegisterRequest)(*pb.RegisterResponse, error){
	panic("implement me")
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *pb.IsAdminRequest)(*pb.IsAdminResponse, error){
	panic("implement me")
}