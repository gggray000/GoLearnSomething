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

type Truck interface {
	LoadCargo() error
	UnloadCargo() error
}

type NormalTruck struct {
	id string
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
	id string
	cargo int
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

func processTruck(truck Truck) error {
	fmt.Printf("Processing truck %+v\n", truck)
	if err := truck.LoadCargo(); err != nil {
		return fmt.Errorf("Error loading cargo: %w", err)
	}
	if err := truck.UnloadCargo(); err != nil {
		return fmt.Errorf("Error unloading cargo: %w", err)
	}
	return nil
}

func main(){

	person := make(map[string]any, 0) //or map[string]interface{}
	person["name"] = "Tiago"
	person["age"] = 42

	age, exists := person["age"].(int)
	if !exists {
		log.Fatal("age doesn't exist")
		return
	}
	log.Println(age)
}
