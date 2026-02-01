package main

import (
	"errors"
	"sync"
)

var ErrTruckNotFound = errors.New("truck not found")

type FleetManager interface {
	AddTruck(id string, cargo int) error
	GetTruck(id string) (Truck, error)
	RemoveTruck(id string) error
	UpdateTruckCargo(id string, change int) error
}

type Truck struct {
	ID    string
	Cargo int
}

type truckManager struct {
	trucks map[string]*Truck
	sync.RWMutex
}

func (tm *truckManager) AddTruck(id string, cargo int) error {
	tm.Lock()
	defer tm.Unlock()
	tm.trucks[id] = &Truck{ID: id, Cargo: cargo}
	return nil
}

func (tm *truckManager) GetTruck(id string) (Truck, error) {
	tm.RLock()
	defer tm.RUnlock()
	truck, exists := tm.trucks[id]
	if exists {
		return *truck, nil
	} else {
		return Truck{}, ErrTruckNotFound
	}
}

func (tm *truckManager) RemoveTruck(id string) error {
	tm.Lock()
	defer tm.Unlock()
	delete(tm.trucks, id)
	return nil
}

func (tm *truckManager) UpdateTruckCargo(id string, change int) error {
	tm.Lock()
	defer tm.Unlock()

	truck, exists := tm.trucks[id]
	if !exists {
		return ErrTruckNotFound
	}

	truck.Cargo += change
	return nil
}

func NewTruckManager() truckManager {
	return truckManager{
		trucks: make(map[string]*Truck),
	}
}
