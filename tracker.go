package mint

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync/atomic"
	"time"
)

type contextKey int

const currentStepIDKey contextKey = 0

var activeStepContext = context.Background()

type step struct {
	id        int
	parentID  int
	eventID   int
	name      string
	startTime time.Time
	duration  time.Duration
	children  []*step
}

type TrackerEventHandle struct {
	id      int
	name    string
	tracker *Tracker
}

type TrackerStepHandle struct {
	id      int
	name    string
	eventID int
	tracker *Tracker
}

type Tracker struct {
	steps  map[int]*step
	events map[int]string
	nextID int64
}

type StepResult struct {
	ID        int           `json:"id"`
	ParentID  int           `json:"parentId"`
	Name      string        `json:"name"`
	Duration  time.Duration `json:"duration"`
	Children  []*StepResult `json:"children"`
	StartTime time.Time     `json:"-"`
}

type EventReport struct {
	EventID       int           `json:"eventId"`
	Name          string        `json:"eventName"`
	TotalDuration time.Duration `json:"totalDuration"`
	RootSteps     []*StepResult `json:"steps"`
	StartTime     time.Time     `json:"-"`
}

type Report struct {
	Events []*EventReport `json:"events"`
}

func NewTracker() *Tracker {
	return &Tracker{
		steps:  make(map[int]*step),
		events: make(map[int]string),
		nextID: 1,
	}
}

func (t *Tracker) EventStart(name string) *TrackerEventHandle {
	eventID := int(atomic.AddInt64(&t.nextID, 1))
	t.events[eventID] = name
	return &TrackerEventHandle{id: eventID, name: name, tracker: t}
}

func (e *TrackerEventHandle) StepStart(stepName string) *TrackerStepHandle {
	parentID := -1
	if pID, ok := activeStepContext.Value(currentStepIDKey).(int); ok {
		parentID = pID
	}
	return e.tracker.stepStart(stepName, e.id, parentID)
}

func (e *TrackerEventHandle) Step(stepName string, fn func()) {
	step := e.StepStart(stepName)
	defer step.Stop()
	e.tracker.withStepContext(step.id, fn)
}

func (e *TrackerEventHandle) StepWithHandle(stepName string, fn func(*TrackerStepHandle)) {
	step := e.StepStart(stepName)
	defer step.Stop()
	e.tracker.withStepContext(step.id, func() { fn(step) })
}

func (s *TrackerStepHandle) Stop() {
	s.tracker.stepStop(s.id)
}

func (s *TrackerStepHandle) StepStart(stepName string) *TrackerStepHandle {
	return s.tracker.stepStart(stepName, s.eventID, s.id)
}

func (s *TrackerStepHandle) Step(stepName string, fn func()) {
	step := s.StepStart(stepName)
	defer step.Stop()
	s.tracker.withStepContext(step.id, fn)
}

func (s *TrackerStepHandle) StepWithHandle(stepName string, fn func(*TrackerStepHandle)) {
	step := s.StepStart(stepName)
	defer step.Stop()
	s.tracker.withStepContext(step.id, func() { fn(step) })
}

func (t *Tracker) withStepContext(stepID int, fn func()) {
	oldCtx := activeStepContext
	activeStepContext = context.WithValue(oldCtx, currentStepIDKey, stepID)
	defer func() { activeStepContext = oldCtx }()
	fn()
}

func (t *Tracker) stepStart(stepName string, eventID, parentID int) *TrackerStepHandle {
	stepID := int(atomic.AddInt64(&t.nextID, 1))
	s := &step{
		id:        stepID,
		parentID:  parentID,
		eventID:   eventID,
		name:      stepName,
		startTime: time.Now(),
	}
	t.steps[stepID] = s
	return &TrackerStepHandle{id: stepID, name: stepName, eventID: eventID, tracker: t}
}

func (t *Tracker) stepStop(stepID int) {
	if s := t.steps[stepID]; s != nil {
		s.duration = time.Since(s.startTime)
	}
}

func (t *Tracker) Reset() {
	t.steps = make(map[int]*step)
	t.events = make(map[int]string)
	atomic.StoreInt64(&t.nextID, 1)
}

func (t *Tracker) GenerateReport(eventIDs ...int) *Report {
	eventSteps := make(map[int][]*step)

	// Group steps by event
	for _, s := range t.steps {
		eventSteps[s.eventID] = append(eventSteps[s.eventID], s)
	}

	var idsToProcess []int
	if len(eventIDs) > 0 {
		for _, id := range eventIDs {
			if _, exists := eventSteps[id]; exists {
				idsToProcess = append(idsToProcess, id)
			}
		}
	} else {
		for id := range eventSteps {
			idsToProcess = append(idsToProcess, id)
		}
	}
	sort.Ints(idsToProcess)

	report := &Report{Events: make([]*EventReport, 0, len(idsToProcess))}

	for _, eventID := range idsToProcess {
		steps := eventSteps[eventID]
		eventName := t.events[eventID]
		if eventName == "" {
			eventName = "Unknown"
		}

		// Calculate event bounds
		var eventStart, eventEnd time.Time
		for _, s := range steps {
			if eventStart.IsZero() || s.startTime.Before(eventStart) {
				eventStart = s.startTime
			}
			stepEnd := s.startTime.Add(s.duration)
			if eventEnd.IsZero() || stepEnd.After(eventEnd) {
				eventEnd = stepEnd
			}
		}

		var totalDuration time.Duration
		if !eventStart.IsZero() && !eventEnd.IsZero() {
			totalDuration = eventEnd.Sub(eventStart)
		}

		report.Events = append(report.Events, &EventReport{
			EventID:       eventID,
			Name:          eventName,
			TotalDuration: totalDuration,
			RootSteps:     t.buildStepHierarchy(steps),
			StartTime:     eventStart,
		})
	}

	return report
}

func (t *Tracker) buildStepHierarchy(steps []*step) []*StepResult {
	stepMap := make(map[int]*StepResult, len(steps))

	for _, s := range steps {
		stepMap[s.id] = &StepResult{
			ID:        s.id,
			ParentID:  s.parentID,
			Name:      s.name,
			Duration:  s.duration,
			StartTime: s.startTime,
			Children:  make([]*StepResult, 0),
		}
	}

	var roots []*StepResult
	for _, result := range stepMap {
		if result.ParentID == -1 {
			roots = append(roots, result)
		} else if parent := stepMap[result.ParentID]; parent != nil {
			parent.Children = append(parent.Children, result)
		}
	}

	var sortSteps func([]*StepResult)
	sortSteps = func(stepList []*StepResult) {
		sort.Slice(stepList, func(i, j int) bool {
			return stepList[i].StartTime.Before(stepList[j].StartTime)
		})
		for _, step := range stepList {
			sortSteps(step.Children)
		}
	}
	sortSteps(roots)

	return roots
}

func (r *Report) ToTree() {
	for _, event := range r.Events {
		fmt.Printf("Event: %s\n", event.Name)
		for i, step := range event.RootSteps {
			r.printStep(step, "", i == len(event.RootSteps)-1)
		}
		fmt.Printf("Total Duration: %s\n\n", formatDuration(event.TotalDuration))
	}
}

func formatDuration(d time.Duration) string {
	ns := d.Nanoseconds()

	if ns < 1000 {
		return fmt.Sprintf("%dns", ns)
	} else if ns < 1000000 {
		us := float64(ns) / 1000.0
		return fmt.Sprintf("%.3fμs", us)
	} else if ns < 1000000000 {
		ms := float64(ns) / 1000000.0
		return fmt.Sprintf("%.3fms", ms)
	} else {
		s := float64(ns) / 1000000000.0
		return fmt.Sprintf("%.3fs", s)
	}
}

func (r *Report) printStep(step *StepResult, prefix string, isLast bool) {
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	fmt.Printf("%s%s%s (%s)\n", prefix, connector, step.Name, formatDuration(step.Duration))

	childPrefix := prefix + "│   "
	if isLast {
		childPrefix = prefix + "    "
	}

	for i, child := range step.Children {
		r.printStep(child, childPrefix, i == len(step.Children)-1)
	}
}

func (r *Report) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}
