package mint

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type config struct {
	Events      []Event     `json:"events"`
	EventsTable eventsTable `json:"eventsTable"`
}

func ParseConfig(reader io.Reader) (config, error) {
	var cfg config
	err := json.NewDecoder(reader).Decode(&cfg)
	if err != nil {
		return config{}, fmt.Errorf("failed to parse json, %w", err)
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
		return config{}, fmt.Errorf("events have duplicates, %s", b.String())
	}

	err = cfg.EventsTable.Validate(len(cfg.Events))
	if err != nil {
		return config{}, err
	}

	return cfg, nil
}
