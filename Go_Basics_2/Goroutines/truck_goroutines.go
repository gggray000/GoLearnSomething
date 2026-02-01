package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	ErrNotImplemented = errors.New("Not implemented.")
	ErrTruckNotFound  = errors.New("Truck not found.")
)

type contextKey string
var usersIDKey contextKey = "userID"

type Truck interface {
	LoadCargo() error
	UnloadCargo() error
}

type NormalTruck struct {
	id    string
	cargo int
}

func (t *NormalTruck) LoadCargo() error {
	t.cargo += 2
	return nil
}

func (t *NormalTruck) UnloadCargo() error {
	t.cargo = 0
	return nil
}

type ElectricTruck struct {
	id      string
	cargo   int
	battery float64
}

func (t *ElectricTruck) LoadCargo() error {
	t.cargo += 1
	t.battery -= 1
	return nil
}

func (t *ElectricTruck) UnloadCargo() error {
	t.cargo = 0
	t.battery -= 1
	return nil
}

func processTruck(ctx context.Context, truck Truck) error {
	fmt.Printf("Processing truck %+v\n", truck)

	// access user id from ctx
	userID := ctx.Value(usersIDKey)
	log.Println(userID)

	ctx, cancel := context.WithTimeout(ctx, time.Second*2) 
	defer cancel()

	delay := time.Second * 3
	select {
	case <- ctx.Done():
		return ctx.Err()
	case <- time.After(delay):
		break
	}

	if err := truck.LoadCargo(); err != nil {
		return fmt.Errorf("Error loading cargo: %w", err)
	}
	if err := truck.UnloadCargo(); err != nil {
		return fmt.Errorf("Error unloading cargo: %w", err)
	}
	fmt.Printf("Finished processing truck %+v\n", truck)
	return nil
}

func processFleet(ctx context.Context, trucks []Truck) error {
	var waitGroup sync.WaitGroup
	errorsChan := make(chan error, len(trucks))

	for _, t := range trucks {
		waitGroup.Add(1)
		go func(t Truck) {
			if err := processTruck(ctx, t); err != nil {
				log.Println(err)
				errorsChan <- err
			}
			waitGroup.Done()
		}(t)
	}
	waitGroup.Wait()
	close(errorsChan)
	
	var errs []error
	for err := range errorsChan {
		log.Printf("Error processing truck: %v\n", err)
		errs = append(errs, err)
	}

	if len(errs)>0{
		return fmt.Errorf("Fleet processing had %d errors", len(errs))
	}
	return nil
}

func main() {

	ctx := context.Background()
	ctx = context.WithValue(ctx, usersIDKey, 42)

	fleet := []Truck{
		&NormalTruck{id: "1", cargo: 0},
		&ElectricTruck{id: "2", cargo: 0, battery: 100},
		&NormalTruck{id: "3", cargo: 10},
		&ElectricTruck{id: "4", cargo: 0, battery: 50},
	}

	if err := processFleet(ctx, fleet); err != nil {
		fmt.Printf("Error processing fleet: %v\n", err)
	} else {
		fmt.Println("All trucks are processed sucessfully")
	}

	m := make(map[string]int)
	var wg sync.WaitGroup

	for i:=0; i<100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(1 * time.Second)

			m[fmt.Sprintf("key-%d", i)] = i // race condition
		}(i)
	}
	wg.Wait()
	fmt.Println("Map:", m)

}
