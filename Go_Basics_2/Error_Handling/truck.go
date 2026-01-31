package main

import (
	"errors"
	"fmt"
	"log"
)

var (
	ErrNotImplemented = errors.New("Not implemented.")
	ErrTruckNotFound = errors.New("Truck not found.")
)

type Truck struct {
	id string
	cargo int
}

func (t *Truck) LoadCargo() error {
	return ErrTruckNotFound
}

func processTruck(truck Truck) error {
	fmt.Printf("Processing Truck: %s\n", truck.id)
	if err := truck.LoadCargo(); err != nil {
		return fmt.Errorf("Error loading cargo: %w", err)
	}
	return ErrNotImplemented
}

func main(){
	trucks := []Truck{
		{id: "1"},
		{id: "2"},
		{id: "3"},
	}

	for _, truck := range trucks {
		fmt.Printf("Truck %s arrived.", truck.id)

		if err := processTruck(truck); err != nil {

			// switch err {
			// case ErrNotImplemented:
			// 	do something
			// case ErrTruckNotFound:
			// 	do something else
			// }

			log.Fatalf("Error processing truck No.%s: %v", truck.id, err)
		}
	}
}