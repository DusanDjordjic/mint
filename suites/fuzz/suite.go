package fuzz

import (
	"fmt"
	"math/rand/v2"

	"github.com/DusanDjordjic/mint"
	"github.com/DusanDjordjic/mint/tracker"
)

type FuzzSuite struct {
	events      mint.Events
	eventsTable EventsTable

	// Current event in the fuzz test used to generate a new one
	currentEvent mint.Event

	// Min and max number of events that are generated for a test
	min uint32
	max uint32

	// RNG used while fuzz testing, created from the seed
	r *rand.Rand

	executor mint.Executor
}

func New(
	seed uint64,
	events mint.Events,
	table EventsTable,
	start mint.Event,
	min uint32,
	max uint32,
	executor mint.Executor,
) FuzzSuite {
	if max < min {
		panic(fmt.Sprintf("max %d < min %d", max, min))
	}

	return FuzzSuite{
		events:       events,
		eventsTable:  table,
		currentEvent: start,
		min:          min,
		max:          max,
		executor:     executor,
		// TODO check how does PCG work and if the second seed is a good one
		r: rand.New(rand.NewPCG(seed, (^seed<<1)|1)),
	}
}

// returts random number of steps to run
func (suite *FuzzSuite) nSteps() uint32 {
	x := suite.max - suite.min
	return suite.r.Uint32N(x+1) + suite.min
}

func (suite *FuzzSuite) Run() {
	count := suite.nSteps()
	test := make([]mint.Event, 1+count)
	test[0] = suite.currentEvent

	var i uint32
	for i = 1; i <= count; i++ {
		x := int16(suite.r.IntN(101))
		nextIndex := suite.eventsTable.Next(x, suite.currentEvent.Index)
		nextEvent := suite.events[nextIndex]
		test[i] = nextEvent
		suite.currentEvent = nextEvent
	}

	tracker := tracker.New()
	suite.executor.Init(tracker)
	for _, event := range test {
		suite.executor.Execute(event, tracker)
	}
	suite.executor.Deinit(tracker)

}

type EventsTable struct {
	Probabilities []int16 `json:"probabilities"`
	nEvents       int
}

// based on the current event and the table of probabilities
// generates a new event. X [0, 100] is used to get the next event
func (table *EventsTable) Next(x int16, event int) int {
	// row of probabilites for the passed event
	row := table.Probabilities[event*table.nEvents : (event+1)*table.nEvents]
	for i, y := range row {
		x = x - y
		if x <= 0 {
			return i
		}
	}

	panic(fmt.Sprintf("unreachable: failed to generate next event for event: %d, x: %d, %v\n",
		event, x, row,
	))
}

func (table *EventsTable) Validate(numberOfEvents int) error {
	if len(table.Probabilities) != numberOfEvents*numberOfEvents {
		return fmt.Errorf("number of probabilities must be number of events ^2, but we got %d events and %d probabilities", numberOfEvents, len(table.Probabilities))
	}

	for row := range numberOfEvents {
		var sum int16 = 0
		for col := range numberOfEvents {

			sum += table.Probabilities[row*numberOfEvents+col]
		}

		if sum != 100 {
			return fmt.Errorf("sum of probabilities in %d. row is not 100 but %d", row+1, sum)
		}
	}

	return nil
}
