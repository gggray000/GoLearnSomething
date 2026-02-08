package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	pb_d "ride-sharing/shared/proto/driver"
)

type inmemRepository struct {
	trips     map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

func NewInmemRepository() *inmemRepository {
	return &inmemRepository{
		trips:     make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (r *inmemRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

func (r *inmemRepository) SaveRideFare(ctx context.Context, f *domain.RideFareModel) error {
	r.rideFares[f.ID.Hex()] = f
	return nil
}

func (r *inmemRepository) GetRideFareByID(ctx context.Context, id string) (*domain.RideFareModel, error) {
	fare, exist := r.rideFares[id]
	if !exist {
		return nil, fmt.Errorf("Fare does not exist with ID: %s", id)
	}
	return fare, nil
}

func (r *inmemRepository) GetTripByID(ctx context.Context, id string)(*domain.TripModel, error) {
	trip, exist := r.trips[id]
	if !exist {
		return nil, fmt.Errorf("Trip does not exist with ID: %s", id)
	}
	return trip, nil
}

func (r *inmemRepository) UpdateTrip(ctx context.Context, tripID string, status string, driver *pb_d.Driver) error {
	trip, exist := r.trips[tripID]
	if !exist {
		return fmt.Errorf("Trip does not exist with ID: %s", tripID)
	}
	trip.Status = status
	if driver != nil{
		trip.Driver = &pb.TripDriver{
			Id: driver.Id,
			Name: driver.Name,
			CarPlate: driver.CarPlate,
			ProfilePicture: driver.ProfilePicture,
		}
	}
	return nil
}
