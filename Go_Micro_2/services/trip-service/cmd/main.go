package main

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"time"
)

func main() {
	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)
	ctx := context.Background()
	fare := &domain.RideFareModel{
		UserID: "42",
	}
	t, err := svc.CreateTrip(ctx, fare)
	if err != nil {
		log.Println(err)
	}

	log.Println(t)

	// keep the program alive
	for {
		time.Sleep(time.Second)
	}
}
