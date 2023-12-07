///////////////////////////////////////////////////////////////////////
/////////////// Automatic Factory Floor using Robots //////////////////
///////////////////////////////////////////////////////////////////////

// Summary:

// Factories are operated through heavily automated orchestrations of
// robots and assembly stations. Factory Floor is made up of a collection
// of specialized robots (e.g. transportation robots, painting robots,
// welding robots and assembly robots) which serve different tasks
// (e.g. transporting components, welding, paiting, etc.) and stations
// (e.g. pickup stations, assembly stations, welding stations, painting
// stations, dropoff stations), where certain tasks must be carried out.
// The factory also contains a control system that coordinates the robots
// on the factory floor, dispatching them to perform the appropriate task
// to the appropriate station.

///////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"log"
	"time"
)

//////////////////// Definition of Data Structures ////////////////////

/////// robots and their tasks ///////

// one worker/robot of a certain specialization
type Worker struct {
	id             int
	specialization *WorkerSet
	inbox          chan TaskSet
	next_facility  chan *Facility
	task_completed chan bool
}

// list of workers of a certain specialization
type WorkerSet struct {
	workers        []*Worker
	specialization string
	freeWorkers    chan *Worker
}

// robot's task
type Task struct {
	FacilityType    *FacilitySet
	Facility        *Facility // null-pointer if facility not assigned yet
	Transporter     *Worker   // null-pointer if transporter not assigned yet
	description     string    // like paint in blue
	assignedWorkers []*Worker
	completed       bool
	tasksetID       int
}

// list of tasks
type TaskSet struct {
	id    int
	tasks []*Task
}

/////////// facitilies ///////////

// one facility of a certain type
type Facility struct {
	id             int
	facilityType   string
	workerArrival  chan *Worker
	taskAssignment chan *Task
}

// list of facilities of a certain type
type FacilitySet struct {
	facilities     []*Facility
	facilityType   string
	freeFacilities chan *Facility
	taskAssignment chan *Task
}

// //////// control center //////////
type ControlCenter struct {

	// parts of the factory
	// each facility set contains a list of facilities of a certain type
	// and a channel to communicate with the control center
	PickupStations   *FacilitySet
	AssemblyStations *FacilitySet
	WeldingStations  *FacilitySet
	PaintingStations *FacilitySet
	DropoffStations  *FacilitySet

	// each worker set contains a list of workers of a certain specialization
	// worker of certain specialization use channels to communicate with the control center
	AssemblyWorkers  *WorkerSet
	WeldingWorkers   *WorkerSet
	PaintingWorkers  *WorkerSet
	TransportWorkers *WorkerSet

	// inbox when new components arrive
	request chan *TaskSet // when a truck comes in with a component it sends a request to the control center

	// channels to communicate
	workerArrival   chan *Worker
	taskSetFinished chan *TaskSet

	// program time
	ProgramTime *ProgramTime

	// counter for completed task sets
	CompletedTaskSets int
}

// ///// time ///////
type ProgramTime struct {
	time int
}

// //////////////////// Time functionality //////////////////////
// simple timer for the program
// not necessary here but generally useful

// get program time
func (pt *ProgramTime) GetCurrentTime() int {
	return pt.time
}

// start program time
func StartProgramTime() *ProgramTime {
	// receives current time
	ticker := time.NewTicker(time.Second)

	// store program time
	var programTime ProgramTime

	// increment program time
	go func() {
		for range ticker.C {
			programTime.time++
		}
	}()

	return &programTime
}

// //////////////////// Boot Facility //////////////////////
func (controlCenter *ControlCenter) Boot() {

	///// Start stations /////

	// start pickup stations
	// every pickup station is free at the beginning
	// every pickup station is started in a separate go routine
	for _, pickupStation := range controlCenter.PickupStations.facilities {
		controlCenter.PickupStations.freeFacilities <- pickupStation
		go pickupStation.RunPickupStation(controlCenter.PickupStations.freeFacilities)
	}

	// start dropoff stations
	// every dropoff station is free at the beginning
	// every dropoff station is started in a separate go routine
	for _, dropoffStation := range controlCenter.DropoffStations.facilities {
		controlCenter.DropoffStations.freeFacilities <- dropoffStation
		go dropoffStation.RunDropoffStation(controlCenter.DropoffStations.freeFacilities)
	}

	// start all assembly stations
	// every assembly station is free at the beginning
	// every assembly station is started in a separate go routine
	for _, assemblyStation := range controlCenter.AssemblyStations.facilities {
		controlCenter.AssemblyStations.freeFacilities <- assemblyStation
		go assemblyStation.RunAssemblyStation(controlCenter.AssemblyStations.freeFacilities)
	}

	// start all welding stations
	// every welding station is free at the beginning
	// every welding station is started in a separate go routine
	for _, weldingStation := range controlCenter.WeldingStations.facilities {
		controlCenter.WeldingStations.freeFacilities <- weldingStation
		go weldingStation.RunWeldingStation(controlCenter.WeldingStations.freeFacilities)
	}

	// start all painting stations
	// every painting station is free at the beginning
	// every painting station is started in a separate go routine
	for _, paintingStation := range controlCenter.PaintingStations.facilities {
		controlCenter.PaintingStations.freeFacilities <- paintingStation
		go paintingStation.RunPaintingStation(controlCenter.PaintingStations.freeFacilities)
	}

	///// Start workers /////

	// start all transport workers
	// every transport worker is free at the beginning
	// every transport worker is started in a separate go routine
	for _, transportWorker := range controlCenter.TransportWorkers.workers {
		transportWorker.specialization = controlCenter.TransportWorkers
		controlCenter.TransportWorkers.freeWorkers <- transportWorker
		go transportWorker.RunTransportWorker(controlCenter, controlCenter.TransportWorkers.freeWorkers)
	}

	// start all assembly workers
	// every assembly worker is free at the beginning
	// every assembly worker is started in a separate go routine
	for _, assemblyWorker := range controlCenter.AssemblyWorkers.workers {
		assemblyWorker.specialization = controlCenter.AssemblyWorkers
		controlCenter.AssemblyWorkers.freeWorkers <- assemblyWorker
		go assemblyWorker.RunAssemblyWorker(controlCenter.AssemblyWorkers.freeWorkers)
	}

	// start all welding workers
	// every welding worker is free at the beginning
	// every welding worker is started in a separate go routine
	for _, weldingWorker := range controlCenter.WeldingWorkers.workers {
		weldingWorker.specialization = controlCenter.WeldingWorkers
		controlCenter.WeldingWorkers.freeWorkers <- weldingWorker
		go weldingWorker.RunWeldingWorker(controlCenter.WeldingWorkers.freeWorkers)
	}

	// start all painting workers
	// every painting worker is free at the beginning
	// every painting worker is started in a separate go routine
	for _, paintingWorker := range controlCenter.PaintingWorkers.workers {
		paintingWorker.specialization = controlCenter.PaintingWorkers
		controlCenter.PaintingWorkers.freeWorkers <- paintingWorker
		go paintingWorker.RunPaintingWorker(controlCenter.PaintingWorkers.freeWorkers)
	}

	// start control center
	go controlCenter.RunControlCenter()

	// Print factory status
	fmt.Println("🌱\nbooted factory with", len(controlCenter.PickupStations.facilities), "pick-up stations,", len(controlCenter.AssemblyStations.facilities), "assembly stations,", len(controlCenter.WeldingStations.facilities), "welding stations,", len(controlCenter.PaintingStations.facilities), "painting stations,", len(controlCenter.DropoffStations.facilities), "drop-off stations,", len(controlCenter.AssemblyWorkers.workers), "assembly workers,", len(controlCenter.WeldingWorkers.workers), "welding workers,", len(controlCenter.PaintingWorkers.workers), "painting workers and", len(controlCenter.TransportWorkers.workers), "transport workers\n🌱")

}

// /////////////// Run Control center /////////////////

func (controlCenter *ControlCenter) RunControlCenter() {
	//// Request handling ////

	// on incoming requests:1. assign free pickup station
	// 						2. notify assigned pickup station
	// 						3. assign free transportations worker
	go controlCenter.HandleRequests()

	//// Task specification ////

	// The transportation robot has the task on its list, e. g. weld steel bar
	// on incoming task:  1. assign free facility of the specific type
	// 					  2. notify transportation worker
	// 					  3. assign free workers of the specific type
	go controlCenter.HandleWeldingAssignments()
	go controlCenter.HandleAssemblyAssignments()
	go controlCenter.HandlePaintingAssignments()
	go controlCenter.HandleDropoffAssignments()
	go controlCenter.TaskFinishedInbox()
}

func (controlCenter *ControlCenter) TaskFinishedInbox() {
	for {
		taskset := <-controlCenter.taskSetFinished
		fmt.Println("\n✅ taskset", taskset.id, "was completed ✅\n ")
	}
}

// We use central controls for each station type
// thus requests from transporters to welding stations
// and e. g. painting stations can be handled in parallel
// (assignment of station and workers and respective notifications)
// but when e. g. 1000 transporter get to painting stations
// this is handled by the control one by one

// Higher parallelization and decentralization could be obtained
// by moving the code now in the handlers to the transporters
// so that the transporters directly race on the channels containing
// the free facilities and workers. To prevent deadlocks in the case
// of the welding station (now avoided) where one transporter gets
// hold of one worker and the other of the other and then forever wait
// some kind of lock mechanism would have to be used.

// get initial request from the trucks and get things going
func (controlCenter *ControlCenter) HandleRequests() {
	for {
		// wait for request to arrive
		request := <-controlCenter.request
		fmt.Println("\n📨: taskset", request.id, "received.\n ")
		// assign facility
		pickupStation := <-controlCenter.PickupStations.freeFacilities
		request.tasks[0].Facility = pickupStation
		pickupStation.taskAssignment <- request.tasks[0]
		// assign free transportation worker
		transportWorker := <-controlCenter.TransportWorkers.freeWorkers
		// for all tasks, assign free transportation worker
		for _, task := range request.tasks {
			task.Transporter = transportWorker
		}
		transportWorker.inbox <- *request
	}
}

// handles welding assignments
// assigns free welding facility and welding workers and notifies transportation worker
func (controlCenter *ControlCenter) HandleWeldingAssignments() {
	for {
		task := <-controlCenter.WeldingStations.taskAssignment
		// assign facility
		facility := <-task.FacilityType.freeFacilities
		task.Facility = facility
		facility.taskAssignment <- task
		// notify transportation worker
		task.Transporter.next_facility <- facility
		// assign workers
		worker1 := <-controlCenter.WeldingWorkers.freeWorkers
		worker2 := <-controlCenter.WeldingWorkers.freeWorkers
		task.assignedWorkers = []*Worker{worker1, worker2}
		worker1.inbox <- TaskSet{99, []*Task{task}} // 99 is default id for trivial tasks
		worker2.inbox <- TaskSet{99, []*Task{task}}
	}
}

// handles assembly assignments
// assigns free assembly facility and assembly workers and notifies transportation worker
func (controlCenter *ControlCenter) HandleAssemblyAssignments() {
	for {
		task := <-controlCenter.AssemblyStations.taskAssignment
		// assign facility
		facility := <-task.FacilityType.freeFacilities
		task.Facility = facility
		facility.taskAssignment <- task
		// notify transportation worker
		task.Transporter.next_facility <- facility
		// assign workers
		worker := <-controlCenter.AssemblyWorkers.freeWorkers
		task.assignedWorkers = []*Worker{worker}
		worker.inbox <- TaskSet{99, []*Task{task}} // 99 is default id for trivial tasks
	}
}

// handles painting assignments
// assigns free painting facility and painting workers and notifies transportation worker
func (controlCenter *ControlCenter) HandlePaintingAssignments() {
	for {
		task := <-controlCenter.PaintingStations.taskAssignment
		// assign facility
		facility := <-task.FacilityType.freeFacilities
		task.Facility = facility
		facility.taskAssignment <- task
		// notify transportation worker
		task.Transporter.next_facility <- facility
		// assign workers
		worker := <-controlCenter.PaintingWorkers.freeWorkers
		task.assignedWorkers = []*Worker{worker}
		worker.inbox <- TaskSet{99, []*Task{task}} // 99 is default id for trivial tasks
	}
}

// handles dropoff assignments
// assigns free dropoff facility and notifies transportation worker
func (controlCenter *ControlCenter) HandleDropoffAssignments() {
	for {
		task := <-controlCenter.DropoffStations.taskAssignment
		// assign facility
		facility := <-task.FacilityType.freeFacilities
		task.Facility = facility
		facility.taskAssignment <- task
		// notify transportation worker
		task.Transporter.next_facility <- facility
	}
}

// //////////////////// Run Facilities //////////////////////

func (facility *Facility) work() {
	// dummy function that simulates some predetermined
	// amount of time for the task to be completed
	time.Sleep(1 * time.Second)
}

// pickup station (only 1 transportation worker per 1 pickup station)
func (pickupStation *Facility) RunPickupStation(freeFacilities chan *Facility) {
	for {
		// wait for task to arrive
		task := <-pickupStation.taskAssignment
		// print task X arrived at pickup station Y
		fmt.Println("[", task.tasksetID, "]", "📝 ➢ 📤: task", task.description, "arrived at pickup station", pickupStation.id)
		// wait for transportation worker to arrive
		transportWorker := <-pickupStation.workerArrival
		// print transportation worker Z arrived at pickup station Y
		fmt.Println("[", task.tasksetID, "]", "🚚 ➢ 📤: transportation worker", transportWorker.id, "arrived at pickup station", pickupStation.id)
		// check that the correct worker type arrived, not strictly necessary as guarantueed by how the
		// control center operatores
		if transportWorker.specialization.specialization != "transport" {
			log.Fatal("Wrong worker arrived at Pickup-Station!")
		}
		// do pickup (sleep)
		pickupStation.work()
		// notify transportation worker that task is completed
		transportWorker.task_completed <- true
		// free facility
		freeFacilities <- pickupStation
		// print pickup station Y is free again
		fmt.Println("[", task.tasksetID, "]", "🕊️ : Pickup station", pickupStation.id, "is free again")
	}
}

// assembly station (1 assembly worker, 1 transportation worker per 1 assembly station)
func (assemblyStation *Facility) RunAssemblyStation(freeFacilities chan *Facility) {
	for {
		// wait for task to arrive
		task := <-assemblyStation.taskAssignment
		// print task X arrived at assembly station Y
		fmt.Println("[", task.tasksetID, "]", "📝 ➢ 🦾: task", task.description, "arrived at assembly station", assemblyStation.id)
		// wait for first worker to arrive
		worker1 := <-assemblyStation.workerArrival
		// print first worker Z arrived at assembly station Y
		fmt.Println("[", task.tasksetID, "]", worker1.to_emoji(), " ➢ 🦾: ", worker1.specialization.specialization, " worker", worker1.id, "arrived at assembly station", assemblyStation.id)
		// wait for second worker to arrive
		worker2 := <-assemblyStation.workerArrival
		// print second worker Z arrived at assembly station Y
		fmt.Println("[", task.tasksetID, "]", worker2.to_emoji(), " ➢ 🦾: ", worker2.specialization.specialization, " worker", worker2.id, "arrived at assembly station", assemblyStation.id)
		// check correct workers arrived
		if !((worker1.specialization.specialization == "transport") && (worker2.specialization.specialization == "assembly") || (worker2.specialization.specialization == "transport") && (worker1.specialization.specialization == "assembly")) {
			log.Fatal("Wrong workers arrived at Assembly-Station!")
		}
		// do assembly (sleep)
		assemblyStation.work()
		fmt.Println("[", task.tasksetID, "]", "🦾 ➢ ✅: assembly task finished")
		// notify all assigned workers that task is completed
		worker1.task_completed <- true
		worker2.task_completed <- true
		// free facility
		freeFacilities <- assemblyStation
		// print type of facility Y is free again
		fmt.Println("[", task.tasksetID, "]", "🕊️ : Assembly station", assemblyStation.id, "is free again")
	}
}

// welding station (2 welders, 1 transportation worker per 1 welding station)
func (weldingStation *Facility) RunWeldingStation(freeFacilities chan *Facility) {
	for {
		// wait for task to arrive
		task := <-weldingStation.taskAssignment
		// print task X arrived at welding station Y
		fmt.Println("[", task.tasksetID, "]", "📝 ➢ 🔨: task", task.description, "arrived at welding station", weldingStation.id)
		// wait for first worker to arrive
		worker1 := <-weldingStation.workerArrival
		// print first worker Z arrived at welding station Y
		fmt.Println("[", task.tasksetID, "]", worker1.to_emoji(), " ➢ 🔨: ", worker1.specialization.specialization, " worker", worker1.id, "arrived at welding station", weldingStation.id)
		// wait for second worker to arrive
		worker2 := <-weldingStation.workerArrival
		// print second worker Z arrived at welding station Y
		fmt.Println("[", task.tasksetID, "]", worker2.to_emoji(), " ➢ 🔨: ", worker2.specialization.specialization, " worker", worker2.id, "arrived at welding station", weldingStation.id)
		// wait for third worker to arrive
		worker3 := <-weldingStation.workerArrival
		// print third worker Z arrived at welding station Y
		fmt.Println("[", task.tasksetID, "]", worker3.to_emoji(), " ➢ 🔨: ", worker3.specialization.specialization, " worker", worker3.id, "arrived at welding station", weldingStation.id)
		// check correct workers arrived
		weld_count := 0
		tranp_count := 0
		for _, worker := range []*Worker{worker1, worker2, worker3} {
			if worker.specialization.specialization == "welding" {
				weld_count++
			} else if worker.specialization.specialization == "transport" {
				tranp_count++
			}
		}
		if !((weld_count == 2) && (tranp_count == 1)) {
			log.Fatal("Wrong workers arrived at Welding-Station!")
		}
		// do welding (sleep)
		weldingStation.work()
		fmt.Println("[", task.tasksetID, "]", "🔨 ➢ ✅: welding task finished")
		// notify all assigned workers that task is completed
		worker1.task_completed <- true
		worker2.task_completed <- true
		worker3.task_completed <- true
		// free facility
		freeFacilities <- weldingStation
		// print type of facility Y is free again
		fmt.Println("[", task.tasksetID, "]", "🕊️ : Welding station", weldingStation.id, "is free again")
	}
}

// painting station (1 painter, 1 transportation worker per 1 painting station)
func (paintingStation *Facility) RunPaintingStation(freeFacilities chan *Facility) {
	for {
		// wait for task to arrive
		task := <-paintingStation.taskAssignment
		// print task X arrived at painting station Y
		fmt.Println("[", task.tasksetID, "]", "📝 ➢ 🎨: task", task.description, "arrived at painting station", paintingStation.id)
		// wait for first worker to arrive
		worker1 := <-paintingStation.workerArrival
		// print first worker Z arrived at painting station Y
		fmt.Println("[", task.tasksetID, "]", worker1.to_emoji(), " ➢ 🎨: ", worker1.specialization.specialization, " worker", worker1.id, "arrived at painting station", paintingStation.id)
		// wait second worker to arrive
		worker2 := <-paintingStation.workerArrival
		// print second worker Z arrived at painting station Y
		fmt.Println("[", task.tasksetID, "]", worker2.to_emoji(), " ➢ 🎨: ", worker2.specialization.specialization, " worker", worker2.id, "arrived at painting station", paintingStation.id)
		// assert correct workers arrived
		if !((worker1.specialization.specialization == "transport") && (worker2.specialization.specialization == "painting") || (worker2.specialization.specialization == "transport") && (worker1.specialization.specialization == "painting")) {
			log.Fatal("Wrong workers arrived at Assembly-Station!")
		}
		// do painting (sleep)
		paintingStation.work()
		fmt.Println("[", task.tasksetID, "]", "🎨 ➢ ✅: painting task finished")
		// notify all assigned workers that task is completed
		worker1.task_completed <- true
		worker2.task_completed <- true
		// free facility
		freeFacilities <- paintingStation
		// print type of facility Y is free again
		fmt.Println("[", task.tasksetID, "]", "🕊️ : Painting station", paintingStation.id, "is free again")
	}
}

// dropoff station (1 transportation worker per 1 dropoff station)
func (dropoffStation *Facility) RunDropoffStation(freeFacilities chan *Facility) {
	for {
		// wait for task to arrive
		task := <-dropoffStation.taskAssignment
		// print task X arrived at dropoff station Y
		fmt.Println("[", task.tasksetID, "]", "📝 ➢ ✈: task", task.description, "arrived at dropoff station", dropoffStation.id)
		// wait for transportation worker to arrive
		transportWorker := <-dropoffStation.workerArrival
		// print transportation worker Z arrived at dropoff station Y
		fmt.Println("[", task.tasksetID, "]", "🚚 ➢ ✈: transportation worker", transportWorker.id, "arrived at dropoff station", dropoffStation.id)
		// check correct workers arrived
		if transportWorker.specialization.specialization != "transport" {
			log.Fatal("Wrong worker arrived at Dropoff-Station!")
		}
		// do dropoff (sleep)
		dropoffStation.work()
		// notify transportation worker that task is completed
		fmt.Println("[", task.tasksetID, "]", "✈ ➢ ✅: dropoff task finished")
		transportWorker.task_completed <- true
		// free facility
		freeFacilities <- dropoffStation
		// print type of facility Y is free again
		fmt.Println("[", task.tasksetID, "]", "🕊️ : Dropoff station", dropoffStation.id, "is free again")
	}
}

// //////////////////// Run Workers //////////////////////

// helper function for printing emojis
func (worker *Worker) to_emoji() string {
	switch worker.specialization.specialization {
	case "assembly":
		return "👷"
	case "welding":
		return "🧑‍"
	case "painting":
		return "🧑‍"
	case "transport":
		return "🚚"
	default:
		return "👷"
	}
}

func (worker *Worker) commute() {
	// dummy function that simulates some predetermined
	// amount of time for traveling between facilities
	time.Sleep(1 * time.Second)
}

// transportation worker
func (transportWorker *Worker) RunTransportWorker(controlCenter *ControlCenter, backToControl chan *Worker) {
	for {
		// wait for task to arrive
		taskset := <-transportWorker.inbox
		// print task X arrived at transportation worker Y
		fmt.Println("📝 ➢➢ 🚚: taskset", taskset.id, "arrived at transportation worker", transportWorker.id)
		// go to pickup station, commute (sleep)
		transportWorker.commute()
		// notify assigned pickup station
		taskset.tasks[0].Facility.workerArrival <- transportWorker
		// wait for task to be completed
		<-transportWorker.task_completed
		// go through all other tasks
		for _, task := range taskset.tasks[1:] {
			// send handling request to control center
			task.FacilityType.taskAssignment <- task
			// wait for next facility
			next_facility := <-transportWorker.next_facility
			// print next facility
			fmt.Println("[", taskset.id, "]", "🚚: next facility of transportation worker", transportWorker.id, "is", next_facility.facilityType, "number", next_facility.id)
			// transport, commute (sleep)
			transportWorker.commute()
			// notify next assigned facility
			next_facility.workerArrival <- transportWorker
			// wait for task to be completed
			<-transportWorker.task_completed
			// set task as completed
			task.completed = true
		}
		controlCenter.taskSetFinished <- &taskset
		controlCenter.CompletedTaskSets++
		// go back to control center, commute (sleep)
		transportWorker.commute()
		// worker notifies control center
		// print transportation worker Y arrived at control center
		fmt.Println("[", taskset.id, "]", "🏠:", "transport worker", transportWorker.id, "arrived at control center")
		// the worker is taking the specific "entrance" for workers of his specialization
		// think of a control center with a room for the transporters, welders, ...
		backToControl <- transportWorker
	}
}

// assembly worker
func (assemblyWorker *Worker) RunAssemblyWorker(backToControl chan *Worker) {
	for {
		// wait for task to arrive
		taskset := <-assemblyWorker.inbox
		task := taskset.tasks[0]
		// print task X arrived at assembly worker Y
		fmt.Println("[", task.tasksetID, "]", "📝 ➢ 👷: task", task.description, "arrived at assembly worker", assemblyWorker.id)
		// go to assembly station, commute (sleep)
		assemblyWorker.commute()
		// notify assigned assembly station
		task.Facility.workerArrival <- assemblyWorker
		// wait for task to be completed
		<-assemblyWorker.task_completed
		// go back to control center, commute (sleep)
		assemblyWorker.commute()
		// print assembly worker Y arrived at control center
		fmt.Println("[", task.tasksetID, "]", "🏠:", "assembly worker", assemblyWorker.id, "arrived at control center")
		// notify control center
		backToControl <- assemblyWorker
	}
}

// welding worker
func (weldingWorker *Worker) RunWeldingWorker(backToControl chan *Worker) {
	for {
		// wait for task to arrive
		taskset := <-weldingWorker.inbox
		task := taskset.tasks[0]
		// print task X arrived at welding worker Y
		fmt.Println("[", task.tasksetID, "]", "📝 ➢ 🧑‍: task", task.description, "arrived at welding worker", weldingWorker.id)
		// go to welding station, commute (sleep)
		weldingWorker.commute()
		// notify assigned welding station
		task.Facility.workerArrival <- weldingWorker
		// wait for task to be completed
		<-weldingWorker.task_completed
		// go back to control center, commute (sleep)
		weldingWorker.commute()
		// print welding worker Y arrived at control center
		fmt.Println("[", task.tasksetID, "]", "🏠:", "welding worker", weldingWorker.id, "arrived at control center")
		// notify control center
		backToControl <- weldingWorker
	}
}

// painting worker
func (paintingWorker *Worker) RunPaintingWorker(backToControl chan *Worker) {
	for {
		// wait for task to arrive
		taskset := <-paintingWorker.inbox
		task := taskset.tasks[0]
		// print task X arrived at painting worker Y
		fmt.Println("[", task.tasksetID, "]", "📝 ➢ 🧑‍: task", task.description, "arrived at painting worker", paintingWorker.id)
		// go to painting station, commute (sleep)
		paintingWorker.commute()
		// notify assigned painting station
		task.Facility.workerArrival <- paintingWorker
		// wait for task to be completed
		<-paintingWorker.task_completed
		// go back to control center, commute (sleep)
		paintingWorker.commute()
		// print painting worker Y arrived at control center
		fmt.Println("[", task.tasksetID, "]", "🏠:", "painting worker", paintingWorker.id, "arrived at control center")
		// notify control center
		backToControl <- paintingWorker
	}
}

// ///////// Build factory ///////////
// construct the factory with the control center struct
// factory has:
//
//	N robots of each kind
//	W facilities on welding stations
//	P facilities on painting stations
//	A facilities on assembly stations
//	I facilities on pick-up stations
//	D facilities on drop-off stations
func BuildFactory(pickupStations int, assemblyStations int, weldingStations int, paintingStations int, dropoffStations int, assemblyWorkers int, weldingWorkers int, paintingWorkers int, transportWorkers int, program_time *ProgramTime) ControlCenter {

	// Start by creating the facility sets

	// Generate the pickup station set with I pickup stations
	pickups := FacilitySet{make([]*Facility, pickupStations), "pickup", make(chan *Facility, pickupStations), make(chan *Task)}
	for i := 0; i < pickupStations; i++ {
		pickups.facilities[i] = &Facility{i, "pickup", make(chan *Worker), make(chan *Task)}
	}

	// Generate the assembly station set with A assembly stations
	assemblies := FacilitySet{make([]*Facility, assemblyStations), "assembly", make(chan *Facility, assemblyStations), make(chan *Task)}
	for i := 0; i < assemblyStations; i++ {
		assemblies.facilities[i] = &Facility{i, "assembly", make(chan *Worker), make(chan *Task)}
	}

	// Generate the welding station set with W welding stations
	weldings := FacilitySet{make([]*Facility, weldingStations), "welding", make(chan *Facility, weldingStations), make(chan *Task)}
	for i := 0; i < weldingStations; i++ {
		weldings.facilities[i] = &Facility{i, "welding", make(chan *Worker), make(chan *Task)}
	}

	// Generate the painting station set with P painting stations
	paintings := FacilitySet{make([]*Facility, paintingStations), "painting", make(chan *Facility, paintingStations), make(chan *Task)}
	for i := 0; i < paintingStations; i++ {
		paintings.facilities[i] = &Facility{i, "painting", make(chan *Worker), make(chan *Task)}
	}

	// Generate the dropoff station set with D dropoff stations
	dropoffs := FacilitySet{make([]*Facility, dropoffStations), "dropoff", make(chan *Facility, dropoffStations), make(chan *Task)}
	for i := 0; i < dropoffStations; i++ {
		dropoffs.facilities[i] = &Facility{i, "dropoff", make(chan *Worker), make(chan *Task)}
	}

	// Generate the worker sets

	// Generate the assembly worker set with N assembly workers
	assemblers := WorkerSet{make([]*Worker, assemblyWorkers), "assembly", make(chan *Worker, assemblyWorkers)}
	for i := 0; i < assemblyWorkers; i++ {
		assemblers.workers[i] = &Worker{i, nil, make(chan TaskSet), make(chan *Facility), make(chan bool)}
	}

	// Generate the welding worker set with N welding workers
	welders := WorkerSet{make([]*Worker, weldingWorkers), "welding", make(chan *Worker, weldingWorkers)}
	for i := 0; i < weldingWorkers; i++ {
		welders.workers[i] = &Worker{i, nil, make(chan TaskSet), make(chan *Facility), make(chan bool)}
	}

	// Generate the painting worker set with N painting workers
	painters := WorkerSet{make([]*Worker, paintingWorkers), "painting", make(chan *Worker, paintingWorkers)}
	for i := 0; i < paintingWorkers; i++ {
		painters.workers[i] = &Worker{i, nil, make(chan TaskSet), make(chan *Facility), make(chan bool)}
	}

	// Generate the transportation worker set with N transportation workers
	transporters := WorkerSet{make([]*Worker, transportWorkers), "transport", make(chan *Worker, transportWorkers)}
	for i := 0; i < transportWorkers; i++ {
		transporters.workers[i] = &Worker{i, nil, make(chan TaskSet), make(chan *Facility), make(chan bool)}
	}

	// Create the control center
	controlCenter := ControlCenter{&pickups, &assemblies, &weldings, &paintings, &dropoffs, &assemblers, &welders, &painters, &transporters, make(chan *TaskSet), make(chan *Worker), make(chan *TaskSet), program_time, 0}
	return controlCenter
}

// //////// Simple Task Set generator ///////////

// generates a task set with the specified id, stations and tasks
// stations and tasks must be of the same length
func gen_task_set(control_center *ControlCenter, id int, stations []string, tasks []string) TaskSet {
	taskset := TaskSet{id, make([]*Task, len(tasks))}
	for i, station := range stations {
		switch station {
		case "pickup":
			taskset.tasks[i] = &Task{control_center.PickupStations, nil, nil, tasks[i], nil, false, id}
		case "welding":
			taskset.tasks[i] = &Task{control_center.WeldingStations, nil, nil, tasks[i], nil, false, id}
		case "assembly":
			taskset.tasks[i] = &Task{control_center.AssemblyStations, nil, nil, tasks[i], nil, false, id}
		case "painting":
			taskset.tasks[i] = &Task{control_center.PaintingStations, nil, nil, tasks[i], nil, false, id}
		case "dropoff":
			taskset.tasks[i] = &Task{control_center.DropoffStations, nil, nil, tasks[i], nil, false, id}
		default:
			fmt.Println("Error: task", station, "not recognized")
		}
	}
	return taskset
}

// //////////////////// Main //////////////////////
func main() {
	// Start program time
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

	/////////////////////// Simple Test ///////////////////////
	tasksetA := gen_task_set(&controlCenter, 1, []string{"pickup", "welding", "assembly", "painting", "dropoff"}, []string{"pickup steel bar", "weld steel bar", "assemble steel bar", "paint steel bar in blue", "dropoff steel bar"})
	controlCenter.request <- &tasksetA

	tasksetB := gen_task_set(&controlCenter, 2, []string{"pickup", "welding", "assembly", "painting", "dropoff"}, []string{"pickup steel wool", "weld steel wool", "assemble steel wool", "paint steel wool in red", "dropoff steel wool"})
	controlCenter.request <- &tasksetB

	tasksetC := gen_task_set(&controlCenter, 3, []string{"pickup", "welding", "assembly", "painting", "dropoff"}, []string{"pickup steel pot", "weld steel pot", "assemble steel pot", "paint steel pot in green", "dropoff steel pot"})
	controlCenter.request <- &tasksetC

	// dummy "keep-alive-system"
	// in the real world, the factory work "forever"
	// and process requests as they come in
	for programTime.GetCurrentTime() < 100 {
		time.Sleep(1 * time.Second)
		// if all task sets are completed wait 2 seconds
		// for transport workers to return to control center
		if controlCenter.CompletedTaskSets == 3 {
			time.Sleep(2 * time.Second)
			break
		}
	}

	// program terminates
}

// //////////////////// TODO //////////////////////
// - compacify code, many things are repeated
