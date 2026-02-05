package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	trip_types "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type TripService interface {
// 	CreateTrip(ctx context.Context, fare RideFareModel) (*TripModel, error)
// }

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{repo: repo}
}

func (s *service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	t := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
		Driver:   &trip.TripDriver{},
	}
	return s.repo.CreateTrip(ctx, t)
}

func (s *service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*trip_types.OsrmApiResponse, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude,
		pickup.Latitude,
		destination.Longitude,
		destination.Latitude,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch route from OSRM API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response from OSRM: %v", err)
	}

	var routeResp trip_types.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("Failed to parse response from OSRM: %v", err)
	}

	return &routeResp, nil
}

func (s *service) StartTrip() {

}

func (s *service) EstimatePackagesPriceWithRoute(route *trip_types.OsrmApiResponse) []*domain.RideFareModel {
	baseFares := getBaseFares()

	estimatedFare := make([]*domain.RideFareModel, len(baseFares))

	for i, f := range baseFares {
		estimatedFare[i] = estimateFareRoute(f, route)
	}

	return estimatedFare
}

func (s *service) GenerateTripFares(ctx context.Context, rideFares []*domain.RideFareModel, userID string, route *trip_types.OsrmApiResponse) ([]*domain.RideFareModel, error) {
	fares := make([]*domain.RideFareModel, len(rideFares))

	for i, f := range rideFares {
		id := primitive.NewObjectID()
		fare := &domain.RideFareModel{
			ID:                id,
			UserID:            userID,
			TotalPriceInCents: f.TotalPriceInCents,
			PackageSlug:       f.PackageSlug,
			Route:             *route,
		}

		if err := s.repo.SaveRideFare(ctx, fare); err != nil {
			return nil, fmt.Errorf("Failed to save trip fare %w", &err)
		}

		fares[i] = fare
	}

	return fares, nil
}

func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 500,
		},
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 500,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 800,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 1000,
		},
	}
}

func estimateFareRoute(f *domain.RideFareModel, route *trip_types.OsrmApiResponse) *domain.RideFareModel {
	pricingCfg := trip_types.DefaultPricingConfig()
	carPackagePrice := f.TotalPriceInCents
	distanceInKm := route.Routes[0].Distance / 1000
	durationInMin := route.Routes[0].Duration / 60
	log.Println("Duration: %w", durationInMin)

	distanceFare := distanceInKm * pricingCfg.PricePerKm
	timeFare := durationInMin * pricingCfg.PricingPerMinute
	totalPrice := carPackagePrice + distanceFare + timeFare

	return &domain.RideFareModel{
		TotalPriceInCents: totalPrice,
		PackageSlug:       f.PackageSlug,
	}
}

func (s *service) GetAndValidateFare(ctx context.Context, fareID, userID string) (*domain.RideFareModel, error) {
	fare, err := s.repo.GetRideFareByID(ctx, fareID)
	if err != nil { 
		return nil, fmt.Errorf("Failed to get trip fare: %w", err)
	}
	if fare == nil {
		return nil, fmt.Errorf("Fare does not exist")
	}
	if userID != fare.UserID {
		return nil, fmt.Errorf("Fare does not belong to user")
	}
	return fare, nil
}
