package tracker

import (
	"fmt"
	"time"
)

type Tracker struct {
	entries []entry
	stack   []uint64
	// starts from 1 because 0 means that there is no parent entry
	nextEntryID uint64
}

type entry struct {
	id       uint64
	name     string
	start    time.Time
	end      time.Time
	parentID uint64
	tracker  *Tracker
}

// Creates a new Tracker and sets inital capacity of entries array.
// We do it for performance reasons
func New(capacity int) *Tracker {
	tracker := &Tracker{
		entries:     make([]entry, 0, capacity),
		stack:       make([]uint64, 0),
		nextEntryID: 1,
	}
	return tracker
}

func (t *Tracker) Step(name string, fn func()) {
	entry := t.Start(name)
	defer entry.Stop()
	fn()
}

func (t *Tracker) Start(name string) *entry {
	id := t.nextEntryID
	t.nextEntryID += 1

	startEntry := entry{
		id:      id,
		name:    name,
		start:   time.Now(),
		tracker: t,
	}

	if len(t.stack) > 0 {
		parent := t.stack[len(t.stack)-1]
		startEntry.parentID = parent
	}

	t.entries = append(t.entries, startEntry)
	t.stack = append(t.stack, startEntry.id)

	return &t.entries[len(t.entries)-1]
}

func (e *entry) Stop() {
	e.end = time.Now()
	lastID := e.tracker.stack[len(e.tracker.stack)-1]

	if e.id != lastID {
		panic(fmt.Sprintf("tried to close %d entry but the last entry is %d", e.id, lastID))
	}

	e.tracker.stack = e.tracker.stack[:len(e.tracker.stack)-1]
}

type TrackerResult struct {
	Entries []ResultEntry
}

type ResultEntry struct {
	ID       uint64    `json:"id"`
	ParentID uint64    `json:"parentID"`
	Name     string    `json:"name"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}

func (t *Tracker) Result() TrackerResult {
	if len(t.stack) != 0 {
		panic(fmt.Sprintf("cannot create results, stack len is %d", len(t.stack)))
	}

	res := TrackerResult{
		Entries: make([]ResultEntry, len(t.entries)),
	}

	for _, e := range t.entries {
		res.Entries = append(res.Entries, ResultEntry{
			ID:       e.id,
			Name:     e.name,
			ParentID: e.parentID,
			Start:    e.start,
			End:      e.end,
		})
	}

	return res
}

// Event: Dusan
// ├── processing (1.063554ms)
// │   ├── println (2.81µs)
// │   └── something (1.060214ms)
// └── time sleep (3.17554ms)
//     └── 123 (1.058523ms)
// Total Duration: 4.239354ms
//
// type StepTiming struct {
// 	fullName string
// 	duration time.Duration
// 	level    int
// }
//
// type StepNode struct {
// 	name     string
// 	duration time.Duration
// 	children []*StepNode
// 	level    int
// }
//
// func (t *Tracker) PrintReportTree() {
// 	fmt.Println("=== Tracker Report  ===")
// 	eventEntries := make(map[int][]timeEntry)
//
// 	for _, entry := range t.entries {
// 		eventEntries[entry.index] = append(eventEntries[entry.index], entry)
// 	}
//
// 	for eventIdx := 0; eventIdx <= t.currentEvent; eventIdx++ {
// 		entries := eventEntries[eventIdx]
// 		if len(entries) == 0 {
// 			continue
// 		}
//
// 		eventName, exists := t.eventNames[eventIdx]
// 		if !exists {
// 			eventName = "Unknown"
// 		}
//
// 		fmt.Printf("Event: %s\n", eventName)
//
// 		var eventStart, eventEnd time.Time
// 		stepTimings := make(map[string]time.Duration)
// 		stepStarts := make(map[string]time.Time)
//
// 		for _, entry := range entries {
// 			if entry.name == "" {
// 				if entry.isStart {
// 					eventStart = entry.time
// 				} else {
// 					eventEnd = entry.time
// 				}
// 			} else {
// 				if entry.isStart {
// 					stepStarts[entry.name] = entry.time
// 				} else {
// 					startTime, exists := stepStarts[entry.name]
// 					if exists {
// 						duration := entry.time.Sub(startTime)
// 						stepTimings[entry.name] = duration
// 						delete(stepStarts, entry.name)
// 					}
// 				}
// 			}
// 		}
//
// 		root := t.buildStepTree(stepTimings)
//
// 		t.printStepTree(root, "", true, true)
//
// 		if !eventStart.IsZero() && !eventEnd.IsZero() {
// 			eventDuration := eventEnd.Sub(eventStart)
// 			fmt.Printf("Total Duration: %v\n", eventDuration)
// 		}
// 		fmt.Println()
// 	}
// }
//
// func (t *Tracker) buildStepTree(timings map[string]time.Duration) *StepNode {
// 	root := &StepNode{name: "root", children: make([]*StepNode, 0)}
// 	nodeMap := make(map[string]*StepNode)
// 	nodeMap[""] = root
//
// 	for fullPath, duration := range timings {
// 		parts := strings.Split(fullPath, "/")
//
// 		currentPath := ""
// 		var parent *StepNode = root
//
// 		for i, part := range parts {
// 			if i > 0 {
// 				currentPath += "/"
// 			}
// 			currentPath += part
//
// 			node, exists := nodeMap[currentPath]
// 			if !exists {
// 				node = &StepNode{
// 					name:     part,
// 					duration: duration,
// 					children: make([]*StepNode, 0),
// 					level:    i,
// 				}
// 				nodeMap[currentPath] = node
// 				parent.children = append(parent.children, node)
// 			} else if i == len(parts)-1 {
// 				node.duration = duration
// 			}
// 			parent = node
// 		}
// 	}
//
// 	return root
// }
//
// func (t *Tracker) printStepTree(node *StepNode, prefix string, isLast bool, isRoot bool) {
// 	if !isRoot {
// 		connector := "├── "
// 		if isLast {
// 			connector = "└── "
// 		}
//
// 		durationStr := fmt.Sprintf("(%v)", node.duration)
// 		fmt.Printf("%s%s%s %s\n", prefix, connector, node.name, durationStr)
// 	}
//
// 	for i, child := range node.children {
// 		childPrefix := prefix
// 		if !isRoot {
// 			if isLast {
// 				childPrefix += "    "
// 			} else {
// 				childPrefix += "│   "
// 			}
// 		}
// 		t.printStepTree(child, childPrefix, i == len(node.children)-1, false)
// 	}
// }
