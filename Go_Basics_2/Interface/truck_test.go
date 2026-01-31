package main

import "testing"

func TestMain(t *testing.T) {
	
	t.Run("processTruck", func(t *testing.T){
		t.Run("should load and unload a truck cargo", func(t *testing.T){
			normal := &NormalTruck{id: "1", cargo: 5}
			electric := &ElectricTruck{id: "2"} 

		err := processTruck(normal); 
		if err != nil {
			t.Fatalf("Error processing truck: %s", err)
		}

		err = processTruck(electric); 
		if err != nil {
			t.Fatalf("Error processing truck: %s", err)
		}

		if normal.cargo != 0 {
			t.Fatal("Normal truck cargo should be 0")
		}

		if electric.battery != -2 {
			t.Fatal("Electric truck battery should be -2")
		}

		})
	})
}