package mint

import (
	"fmt"
)

type eventsTable struct {
	Probabilities []uint16 `json:"probabilities"`
}

func (table *eventsTable) Validate(numberOfEvents int) error {
	if len(table.Probabilities) != numberOfEvents*numberOfEvents {
		return fmt.Errorf("number of probabilities must be number of events ^2, but we got %d events and %d probabilities", numberOfEvents, len(table.Probabilities))
	}

	for row := range numberOfEvents {
		var sum uint16 = 0
		for col := range numberOfEvents {

			sum += table.Probabilities[row*numberOfEvents+col]
		}

		if sum != 100 {
			return fmt.Errorf("sum of probabilities in %d. row is not 100 but %d", row+1, sum)
		}
	}

	return nil
}

type Manager struct {
	config config
}

func (m *Manager) Init(config config) {
	m.config = config
}

func (m *Manager) Start() {}

func New() Manager {
	return Manager{}
}
