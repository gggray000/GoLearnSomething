package domain

import (
	"context"
	trip_types "ride-sharing/services/trip-service/pkg/types"
	pb "ride-sharing/shared/proto/trip"
	pb_d "ride-sharing/shared/proto/driver"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID       primitive.ObjectID
	UserID   string
	Status   string
	RideFare *RideFareModel
	Driver   *pb.TripDriver
}

func (t *TripModel) ToProto() *pb.Trip {
	return &pb.Trip{
		Id: t.ID.Hex(),
		UserID: t.UserID,
		SelectedFare: t.RideFare.toProto(),
		Status: t.Status,
		Driver: t.Driver,
		Route: t.RideFare.Route.ToProto(),
	}
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(ctx context.Context, f *RideFareModel) error
	GetRideFareByID(ctx context.Context, id string) (*RideFareModel, error)
	GetTripByID(ctx context.Context, id string)(*TripModel, error)
	UpdateTrip(ctx context.Context, tripID string, status string, driver *pb_d.Driver) error
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*trip_types.OsrmApiResponse, error)
	EstimatePackagesPriceWithRoute(route *trip_types.OsrmApiResponse) []*RideFareModel
	GenerateTripFares(ctx context.Context, fares []*RideFareModel, userID string, route *trip_types.OsrmApiResponse) ([]*RideFareModel, error)
	GetAndValidateFare(ctx context.Context, fareID, userID string) (*RideFareModel, error)
	GetTripByID(ctx context.Context, id string)(*TripModel, error)
	UpdateTrip(ctx context.Context, tripID string, status string, driver *pb_d.Driver) error
}
