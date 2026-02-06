package main

import (
	"context"
	pb "ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpc_handler struct {
	pb.UnimplementedDriverServiceServer
	Service *Service
}

func NewGRPCHandler(s *grpc.Server, service *Service) {
	handler := &grpc_handler{
		Service: service,
	}
	pb.RegisterDriverServiceServer(s, handler)
}

func (h *grpc_handler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	driver, error := h.Service.RegisterDriver(req.GetDriverID(), req.GetPackageSlug())
	if error != nil {
		return nil, status.Errorf(codes.Internal, "Unable to register driver: %v", error)
	}
	return &pb.RegisterDriverResponse{
		Driver: driver,
	}, nil
}
func (h *grpc_handler) UnregisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	h.Service.UnregisterDriver(req.GetDriverID())

	return &pb.RegisterDriverResponse{
		Driver: &pb.Driver{
			Id: req.GetDriverID(),
		},
	}, nil
}
