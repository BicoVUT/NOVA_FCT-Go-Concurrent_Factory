///////////////////////////////////////////////////////////////////////
/////////////// Automatic Factory Floor using Robots //////////////////
///////////////////////////////////////////////////////////////////////

// This file contains the test cases for the factory floor

package main

import (
	"testing"
	"time"
)

/////////////////////// Helper Functions ///////////////////////

// return a number of current free workers
func FreeWorkersCount(controlCenter *ControlCenter) int {
	var workers int = 0
	workers += len(controlCenter.TransportWorkers.freeWorkers)
	workers += len(controlCenter.AssemblyWorkers.freeWorkers)
	workers += len(controlCenter.PaintingWorkers.freeWorkers)
	workers += len(controlCenter.WeldingWorkers.freeWorkers)
	return workers
}

// return a number of current free facilities
func FreeFacilitiesCount(controlCenter *ControlCenter) int {
	var facilities int = 0
	facilities += len(controlCenter.PickupStations.freeFacilities)
	facilities += len(controlCenter.DropoffStations.freeFacilities)
	facilities += len(controlCenter.AssemblyStations.freeFacilities)
	facilities += len(controlCenter.PaintingStations.freeFacilities)
	facilities += len(controlCenter.WeldingStations.freeFacilities)
	return facilities
}

////////////////////////// Test Cases //////////////////////////

// Test for non-existive pickup station
func TestNoPickup(t *testing.T) {
	programTime := StartProgramTime()

	N := 1 // N robots of each kind
	W := 0 // W facilities on welding stations
	P := 0 // P facilities on painting stations
	A := 0 // A facilities on assembly stations
	I := 0 // I facilities on pick-up stations
	D := 1 // D facilities on drop-off stations

	// Build the factory with specified number of facilities and workers
	controlCenter := BuildFactory(I, A, W, P, D, N, N, N, N, programTime)

	// Boot the control center
	go controlCenter.Boot()

	// Create an impossible pickup and dropoff task (pick-up station is missing)
	impossibleTask := gen_task_set(&controlCenter, 1, []string{"pickup", "dropoff"}, []string{"pickup component", "dropoff component"})

	// Send the task to the control center
	controlCenter.request <- &impossibleTask

	// Simulate the passage of time for 2 seconds
	for programTime.GetCurrentTime() < 2 {
		time.Sleep(1 * time.Second)
	}

	// task should not be completed
	if controlCenter.CompletedTaskSets != 0 {
		t.Errorf("Pickup station is missing, number of completed task should be 0, not %d", controlCenter.CompletedTaskSets)
	}
}

// Test for non-existive dropoff station
func TestNoDropoff(t *testing.T) {
	programTime := StartProgramTime()

	N := 1 // N robots of each kind
	W := 0 // W facilities on welding stations
	P := 0 // P facilities on painting stations
	A := 0 // A facilities on assembly stations
	I := 1 // I facilities on pick-up stations
	D := 0 // D facilities on drop-off stations

	// Build the factory with specified number of facilities and workers
	controlCenter := BuildFactory(I, A, W, P, D, N, N, N, N, programTime)

	// Boot the control center
	go controlCenter.Boot()

	// Create an impossible pickup and dropoff task (dropoff station is missing)
	impossibleTask := gen_task_set(&controlCenter, 1, []string{"pickup", "dropoff"}, []string{"pickup component", "dropoff component"})

	// Send the task to the control center
	controlCenter.request <- &impossibleTask

	// Simulate the passage of time for 3 seconds
	for programTime.GetCurrentTime() < 3 {
		time.Sleep(1 * time.Second)
	}

	// task should not be completed
	if controlCenter.CompletedTaskSets != 0 {
		t.Errorf("Dropoff station is missing, number of completed task should be 0, not %d", controlCenter.CompletedTaskSets)
	}
}

// Test for functionality of pickup and dropoff station
func TestPickupAndDropoffStation(t *testing.T) {
	programTime := StartProgramTime()

	N := 1 // N robots of each kind
	W := 0 // W facilities on welding stations
	P := 0 // P facilities on painting stations
	A := 0 // A facilities on assembly stations
	I := 1 // I facilities on pick-up stations
	D := 1 // D facilities on drop-off stations

	// Build the factory with specified number of facilities and workers
	controlCenter := BuildFactory(I, A, W, P, D, N, N, N, N, programTime)

	// Boot the control center
	go controlCenter.Boot()

	// Create a task set for pickup and dropoff station
	pickupAndDropoffTask := gen_task_set(&controlCenter, 1, []string{"pickup", "dropoff"}, []string{"pickup steel bar", "dropoff steel bar"})

	// Send the pickup and dropoff task set to the control center
	controlCenter.request <- &pickupAndDropoffTask

	// Simulate the passage of time for 10 seconds
	for programTime.GetCurrentTime() < 10 {
		time.Sleep(1 * time.Second)
		if controlCenter.CompletedTaskSets == 1 {
			time.Sleep(2 * time.Second)
			break
		}
	}

	// Check if the pickup and dropoff task is complete
	if controlCenter.CompletedTaskSets != 1 {
		t.Errorf("Total number of tasks done in Factory is %d, want 1", controlCenter.CompletedTaskSets)
	}
}

// Test for functionality of welding station
func TestWeldingStation(t *testing.T) {
	programTime := StartProgramTime()

	N := 2 // N robots of each kind
	W := 1 // W facilities on welding stations
	P := 0 // P facilities on painting stations
	A := 0 // A facilities on assembly stations
	I := 1 // I facilities on pick-up stations
	D := 1 // D facilities on drop-off stations

	// Build the factory with specified number of facilities and workers
	controlCenter := BuildFactory(I, A, W, P, D, N, N, N, N, programTime)

	// Boot the control center
	go controlCenter.Boot()

	// Create a task set for welding station
	weldingTask := gen_task_set(&controlCenter, 1, []string{"pickup", "welding", "dropoff"}, []string{"pickup steel bar", "weld steel bar", "dropoff steel bar"})

	// Send the welding task set to the control center
	controlCenter.request <- &weldingTask

	for programTime.GetCurrentTime() < 10 {
		time.Sleep(1 * time.Second)
		if controlCenter.CompletedTaskSets == 1 {
			time.Sleep(2 * time.Second)
			break
		}
	}

	// Check if the welding task is complete
	if controlCenter.CompletedTaskSets != 1 {
		t.Errorf("Total number of tasks done in Factory is %d, want 1", controlCenter.CompletedTaskSets)
	}
}

// Test for functionality of painting station
func TestPaintingStation(t *testing.T) {
	programTime := StartProgramTime()

	N := 1 // N robots of each kind
	W := 0 // W facilities on welding stations
	P := 1 // P facilities on painting stations
	A := 0 // A facilities on assembly stations
	I := 1 // I facilities on pick-up stations
	D := 1 // D facilities on drop-off stations

	// Build the factory with specified number of facilities and workers
	controlCenter := BuildFactory(I, A, W, P, D, N, N, N, N, programTime)

	// Boot the control center
	go controlCenter.Boot()

	// Create a task set for painting station
	paintingTask := gen_task_set(&controlCenter, 1, []string{"pickup", "painting", "dropoff"}, []string{"pickup component", "paint component", "dropoff component"})

	// Send the painting task set to the control center
	controlCenter.request <- &paintingTask

	for programTime.GetCurrentTime() < 10 {
		time.Sleep(1 * time.Second)
		if controlCenter.CompletedTaskSets == 1 {
			time.Sleep(2 * time.Second)
			break
		}
	}

	// Check if the painting task is complete
	if controlCenter.CompletedTaskSets != 1 {
		t.Errorf("Total number of tasks done in Factory is %d, want 1", controlCenter.CompletedTaskSets)
	}
}

// Test for functionality of assembly station
func TestAssemblyStation(t *testing.T) {
	programTime := StartProgramTime()

	N := 1 // N robots of each kind
	W := 0 // W facilities on welding stations
	P := 0 // P facilities on painting stations
	A := 1 // A facilities on assembly stations
	I := 1 // I facilities on pick-up stations
	D := 1 // D facilities on drop-off stations

	// Build the factory with specified number of facilities and workers
	controlCenter := BuildFactory(I, A, W, P, D, N, N, N, N, programTime)

	// Boot the control center
	go controlCenter.Boot()

	// Create a task set for assembly station
	assemblyTask := gen_task_set(&controlCenter, 1, []string{"pickup", "assembly", "dropoff"}, []string{"pickup component", "assemble component", "dropoff component"})

	// Send the assembly task set to the control center
	controlCenter.request <- &assemblyTask

	for programTime.GetCurrentTime() < 10 {
		time.Sleep(1 * time.Second)
		if controlCenter.CompletedTaskSets == 1 {
			time.Sleep(2 * time.Second)
			break
		}
	}

	// Check if the assembly task is complete
	if controlCenter.CompletedTaskSets != 1 {
		t.Errorf("Total number of tasks done in Factory is %d, want 1", controlCenter.CompletedTaskSets)
	}
}

// Test if we properly free workers and facilities
func TestFree(t *testing.T) {
	programTime := StartProgramTime()

	N := 2 // N robots of each kind
	W := 2 // W facilities on welding stations
	P := 2 // P facilities on painting stations
	A := 2 // A facilities on assembly stations
	I := 2 // I facilities on pick-up stations
	D := 2 // D facilities on drop-off stations

	// Build the factory with specified number of facilities and workers
	controlCenter := BuildFactory(I, A, W, P, D, N, N, N, N, programTime)

	// Boot the control center
	go controlCenter.Boot()

	// generate task set
	taskset := gen_task_set(&controlCenter, 1, []string{"pickup", "welding", "assembly", "painting", "dropoff"}, []string{"pickup steel bar", "weld steel bar", "assemble steel bar", "paint steel bar in blue", "dropoff steel bar"})

	// send task set to control center
	controlCenter.request <- &taskset

	for programTime.GetCurrentTime() < 20 {
		time.Sleep(1 * time.Second)
		if controlCenter.CompletedTaskSets == 1 {
			time.Sleep(2 * time.Second)
			break
		}
	}

	// check if all tasks are done
	if controlCenter.CompletedTaskSets != 1 {
		t.Errorf("Total number of tasks done in Factory is %d, want 1", controlCenter.CompletedTaskSets)
	}

	// check if all workers are free
	if FreeWorkersCount(&controlCenter) != 4*N {
		t.Errorf("Total number of free workers in Factory is %d, want %d", FreeWorkersCount(&controlCenter), 4*N)
	}

	// check if all facilities are free
	if FreeFacilitiesCount(&controlCenter) != I+D+A+W+P {
		t.Errorf("Total number of free facilities in Factory is %d, want %d", FreeFacilitiesCount(&controlCenter), I+D+A+W+P)
	}

}

// Test if all requested tasks are completed and all workers are free
func TestAllCompletedTaskSets(t *testing.T) {
	programTime := StartProgramTime()

	N := 2 // N robots of each kinds
	W := 2 // W facilities on welding stations
	P := 2 // P facilities on painting stations
	A := 2 // A facilities on assembly stations
	I := 2 // I facilities on pick-up stations
	D := 2 // D facilities on drop-off stations

	// Build the factory with specified number of facilities and workers
	controlCenter := BuildFactory(I, A, W, P, D, N, N, N, N, programTime)

	// Boot the control center
	go controlCenter.Boot()

	// generate task sets
	tasksetA := gen_task_set(&controlCenter, 1, []string{"pickup", "welding", "assembly", "painting", "dropoff"}, []string{"pickup steel bar", "weld steel bar", "assemble steel bar", "paint steel bar in blue", "dropoff steel bar"})
	controlCenter.request <- &tasksetA

	tasksetB := gen_task_set(&controlCenter, 2, []string{"pickup", "welding", "assembly", "painting", "dropoff"}, []string{"pickup steel wool", "weld steel wool", "assemble steel wool", "paint steel wool in red", "dropoff steel wool"})
	controlCenter.request <- &tasksetB

	tasksetC := gen_task_set(&controlCenter, 3, []string{"pickup", "welding", "assembly", "painting", "dropoff"}, []string{"pickup steel pot", "weld steel pot", "assemble steel pot", "paint steel pot in green", "dropoff steel pot"})
	controlCenter.request <- &tasksetC

	for programTime.GetCurrentTime() < 30 {
		time.Sleep(1 * time.Second)
		if controlCenter.CompletedTaskSets == 3 {
			time.Sleep(2 * time.Second)
			break
		}
	}

	// check if all tasks are done
	if controlCenter.CompletedTaskSets != 3 {
		t.Errorf("Total number of tasks done in Factory is %d, want 3", controlCenter.CompletedTaskSets)
		t.Errorf("If failed, check if the program time is correct")
	}

	// check if all workers are free
	if FreeWorkersCount(&controlCenter) != 4*N {
		t.Errorf("Total number of free workers in Factory is %d, want %d", FreeWorkersCount(&controlCenter), 4*N)
		t.Errorf("If failed, check if the program time is correct")
	}

	// check if all facilities are free
	if FreeFacilitiesCount(&controlCenter) != I+D+A+W+P {
		t.Errorf("Total number of free facilities in Factory is %d, want %d", FreeFacilitiesCount(&controlCenter), I+D+A+W+P)
		t.Errorf("If failed, check if the program time is correct")
	}

}

// Test if all requested tasks are completed and all workers are free (more complex)
func TestComplexTasks(t *testing.T) {
	programTime := StartProgramTime()

	N := 2 // N robots of each kinds
	W := 1 // W facilities on welding stations
	P := 1 // P facilities on painting stations
	A := 1 // A facilities on assembly stations
	I := 2 // I facilities on pick-up stations
	D := 2 // D facilities on drop-off stations

	// Build the factory with specified number of facilities and workers
	controlCenter := BuildFactory(I, A, W, P, D, N, N, N, N, programTime)

	// Boot the control center
	go controlCenter.Boot()

	// task set A
	tasksetA := gen_task_set(&controlCenter, 1, []string{"pickup", "welding", "assembly", "painting", "dropoff"}, []string{"pickup steel bar", "weld steel bar", "assemble steel bar", "paint steel bar in blue", "dropoff steel bar"})
	controlCenter.request <- &tasksetA

	// task set B
	tasksetB := gen_task_set(&controlCenter, 2, []string{"pickup", "welding", "painting", "assembly", "dropoff"}, []string{"pickup steel wool", "weld steel wool", "paint steel wool in red", "assemble steel wool", "dropoff steel wool"})
	controlCenter.request <- &tasksetB

	// task set C
	tasksetC := gen_task_set(&controlCenter, 3, []string{"pickup", "painting", "welding", "assembly", "dropoff"}, []string{"pickup steel pot", "paint steel pot in green", "weld steel pot", "assemble steel pot", "dropoff steel pot"})
	controlCenter.request <- &tasksetC

	// task set D
	tasksetD := gen_task_set(&controlCenter, 4, []string{"pickup", "welding", "assembly", "painting", "dropoff"}, []string{"pickup steel bar", "weld steel bar", "assemble steel bar", "paint steel bar in blue", "dropoff steel bar"})
	controlCenter.request <- &tasksetD

	// task set E
	tasksetE := gen_task_set(&controlCenter, 5, []string{"pickup", "welding", "painting", "assembly", "dropoff"}, []string{"pickup steel wool", "weld steel wool", "paint steel wool in red", "assemble steel wool", "dropoff steel wool"})
	controlCenter.request <- &tasksetE

	for programTime.GetCurrentTime() < 50 {
		time.Sleep(1 * time.Second)
		if controlCenter.CompletedTaskSets == 5 {
			time.Sleep(2 * time.Second)
			break
		}
	}

	// check if all tasks are done
	if controlCenter.CompletedTaskSets != 5 {
		t.Errorf("Total number of tasks done in Factory is %d, want 5", controlCenter.CompletedTaskSets)
		t.Errorf("If failed, check if the program time is correct")
	}

	// check if all workers are free
	if FreeWorkersCount(&controlCenter) != 4*N {
		t.Errorf("Total number of free workers in Factory is %d, want %d", FreeWorkersCount(&controlCenter), 4*N)
		t.Errorf("If failed, check if the program time is correct")
	}

	// check if all facilities are free
	if FreeFacilitiesCount(&controlCenter) != I+D+A+W+P {
		t.Errorf("Total number of free facilities in Factory is %d, want %d", FreeFacilitiesCount(&controlCenter), I+D+A+W+P)
		t.Errorf("If failed, check if the program time is correct")
	}

}
