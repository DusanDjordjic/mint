package mint

import (
	"fmt"
	"strings"
	"time"
)

type timeEntry struct {
	eventIndex int
	eventName  string
	stepName   string
	time       time.Time
	isStart    bool
}

type Tracker struct {
	entries      []timeEntry
	currentEvent int
	eventNames   map[int]string
	stepStack    []string
}

type StepTiming struct {
	fullName string
	duration time.Duration
	level    int
}

type StepNode struct {
	name     string
	duration time.Duration
	children []*StepNode
	level    int
}

func NewTracker() *Tracker {
	return &Tracker{
		entries:      make([]timeEntry, 0, 1000),
		currentEvent: -1,
		eventNames:   make(map[int]string),
		stepStack:    make([]string, 0),
	}
}

func (t *Tracker) EventStart(name string) {
	t.currentEvent++
	t.eventNames[t.currentEvent] = name
	t.entries = append(t.entries, timeEntry{
		eventIndex: t.currentEvent,
		eventName:  name,
		stepName:   "",
		time:       time.Now(),
		isStart:    true,
	})
}

func (t *Tracker) EventStop() {
	if t.currentEvent >= 0 {
		t.entries = append(t.entries, timeEntry{
			eventIndex: t.currentEvent,
			eventName:  t.eventNames[t.currentEvent],
			stepName:   "",
			time:       time.Now(),
			isStart:    false,
		})
	}
}

func (t *Tracker) Reset() {
	t.entries = t.entries[:0]
	t.currentEvent = -1
	t.eventNames = make(map[int]string)
	t.stepStack = t.stepStack[:0]
}

func (t *Tracker) Step(name string, fn func()) {
	t.StepStart(name)
	defer t.StepStop(name)
	fn()
}

func (t *Tracker) StepStart(stepName string) {
	if t.currentEvent >= 0 {
		fullStepName := stepName
		if len(t.stepStack) > 0 {
			fullStepName = strings.Join(t.stepStack, "/") + "/" + stepName
		}

		t.stepStack = append(t.stepStack, stepName)

		t.entries = append(t.entries, timeEntry{
			eventIndex: t.currentEvent,
			eventName:  t.eventNames[t.currentEvent],
			stepName:   fullStepName,
			time:       time.Now(),
			isStart:    true,
		})
	}
}

func (t *Tracker) StepStop(stepName string) {
	if t.currentEvent >= 0 && len(t.stepStack) > 0 {
		fullStepName := stepName
		if len(t.stepStack) > 1 {
			fullStepName = strings.Join(t.stepStack[:len(t.stepStack)-1], "/") + "/" + stepName
		}

		t.entries = append(t.entries, timeEntry{
			eventIndex: t.currentEvent,
			eventName:  t.eventNames[t.currentEvent],
			stepName:   fullStepName,
			time:       time.Now(),
			isStart:    false,
		})

		if len(t.stepStack) > 0 {
			t.stepStack = t.stepStack[:len(t.stepStack)-1]
		}
	}
}

// ukrao. builduje treeview
// Event: Dusan
// ├── processing (1.063554ms)
// │   ├── println (2.81µs)
// │   └── something (1.060214ms)
// └── time sleep (3.17554ms)
//     └── 123 (1.058523ms)
// Total Duration: 4.239354ms

func (t *Tracker) PrintReportTree() {
	fmt.Println("=== Tracker Report  ===")
	eventEntries := make(map[int][]timeEntry)

	for _, entry := range t.entries {
		eventEntries[entry.eventIndex] = append(eventEntries[entry.eventIndex], entry)
	}

	for eventIdx := 0; eventIdx <= t.currentEvent; eventIdx++ {
		entries := eventEntries[eventIdx]
		if len(entries) == 0 {
			continue
		}

		eventName, exists := t.eventNames[eventIdx]
		if !exists {
			eventName = "Unknown"
		}

		fmt.Printf("Event: %s\n", eventName)

		var eventStart, eventEnd time.Time
		stepTimings := make(map[string]time.Duration)
		stepStarts := make(map[string]time.Time)

		for _, entry := range entries {
			if entry.stepName == "" {
				if entry.isStart {
					eventStart = entry.time
				} else {
					eventEnd = entry.time
				}
			} else {
				if entry.isStart {
					stepStarts[entry.stepName] = entry.time
				} else {
					startTime, exists := stepStarts[entry.stepName]
					if exists {
						duration := entry.time.Sub(startTime)
						stepTimings[entry.stepName] = duration
						delete(stepStarts, entry.stepName)
					}
				}
			}
		}

		root := t.buildStepTree(stepTimings)

		t.printStepTree(root, "", true, true)

		if !eventStart.IsZero() && !eventEnd.IsZero() {
			eventDuration := eventEnd.Sub(eventStart)
			fmt.Printf("Total Duration: %v\n", eventDuration)
		}
		fmt.Println()
	}
}

func (t *Tracker) buildStepTree(timings map[string]time.Duration) *StepNode {
	root := &StepNode{name: "root", children: make([]*StepNode, 0)}
	nodeMap := make(map[string]*StepNode)
	nodeMap[""] = root

	for fullPath, duration := range timings {
		parts := strings.Split(fullPath, "/")

		currentPath := ""
		var parent *StepNode = root

		for i, part := range parts {
			if i > 0 {
				currentPath += "/"
			}
			currentPath += part

			node, exists := nodeMap[currentPath]
			if !exists {
				node = &StepNode{
					name:     part,
					duration: duration,
					children: make([]*StepNode, 0),
					level:    i,
				}
				nodeMap[currentPath] = node
				parent.children = append(parent.children, node)
			} else if i == len(parts)-1 {
				node.duration = duration
			}
			parent = node
		}
	}

	return root
}

func (t *Tracker) printStepTree(node *StepNode, prefix string, isLast bool, isRoot bool) {
	if !isRoot {
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		durationStr := fmt.Sprintf("(%v)", node.duration)
		fmt.Printf("%s%s%s %s\n", prefix, connector, node.name, durationStr)
	}

	for i, child := range node.children {
		childPrefix := prefix
		if !isRoot {
			if isLast {
				childPrefix += "    "
			} else {
				childPrefix += "│   "
			}
		}
		t.printStepTree(child, childPrefix, i == len(node.children)-1, false)
	}
}
