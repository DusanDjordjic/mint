package mint

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Config struct {
	Events []Event `json:"events"`
	// If we are doing fuzz than we need events table
	// same for the load tests
	// EventsTable eventsTable `json:"eventsTable"`
	// if we are doing normal tests then we need tests cases
	// PreTest  []Event
	// Tests    [][]Event
	// PostTest []Event
}

func ParseConfig(reader io.Reader) (Config, error) {
	var cfg Config
	err := json.NewDecoder(reader).Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse json, %w", err)
	}

	m := make(map[Event]int, len(cfg.Events))
	for _, e := range cfg.Events {
		x, ok := m[e]
		if ok {
			x++
			m[e] = x
		} else {
			m[e] = 1
		}
	}

	b := strings.Builder{}
	for event, count := range m {
		if count > 1 {
			b.WriteString(fmt.Sprintf("event %s occurs %d times", event, count))
		}
	}

	if b.Len() > 0 {
		return Config{}, fmt.Errorf("events have duplicates, %s", b.String())
	}

	// err = cfg.EventsTable.Validate(len(cfg.Events))
	// if err != nil {
	// 	return Config{}, err
	// }

	return cfg, nil
}
