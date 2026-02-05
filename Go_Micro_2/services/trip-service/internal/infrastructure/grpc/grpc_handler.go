package grpc_handler

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	service domain.TripService
}

func NewGRPCHandler(server *grpc.Server, service domain.TripService) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}
	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetStartLocation()
	destination := req.GetEndLocation()

	pickupCoord := &types.Coordinate{
		Latitude:  pickup.Latitude,
		Longitude: pickup.Longitude,
	}

	destinationCoord := &types.Coordinate{
		Latitude:  destination.Latitude,
		Longitude: destination.Longitude,
	}

	route, err := h.service.GetRoute(ctx, pickupCoord, destinationCoord)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "Fail to get route: %v", err)
	}

	estimatedFares := h.service.EstimatePackagesPriceWithRoute(route)
	fares, err := h.service.GenerateTripFares(ctx, estimatedFares, req.UserID, route)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Fail to generate ride fares: %v", err)
	}

	return &pb.PreviewTripResponse{
		Route:     route.ToProto(),
		RideFares: domain.ToRideFaresProto(fares),
	}, nil
}

func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest)(*pb.CreateTripResponse, error) {
	rideFare, err := h.service.GetAndValidateFare(ctx, req.GetRideFareID(), req.GetUserID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Invalid fare: %v", err)
	}
	trip, err := h.service.CreateTrip(ctx, rideFare)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create trip: %v", err)
	}
	return &pb.CreateTripResponse{
		TripID: trip.ID.Hex(),
	}, nil
}
