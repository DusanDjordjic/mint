package tracker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tracker := New(1000)
	assert.EqualValues(t, 1, tracker.nextEntryID)
	assert.Zero(t, len(tracker.entries))
	assert.Equal(t, cap(tracker.entries), 1000)

	assert.Zero(t, len(tracker.stack))
	assert.Zero(t, cap(tracker.stack))
}

func TestEntryStart(t *testing.T) {
	tracker := New(1000)

	entryName := "test-entry"
	entry := tracker.Start(entryName)

	assert.EqualValues(t, 1, len(tracker.entries))
	assert.EqualValues(t, 1, len(tracker.stack))
	assert.EqualValues(t, 1, tracker.stack[0])

	assert.Equal(t, entry.name, entryName)
	assert.EqualValues(t, entry.id, 1)
	assert.EqualValues(t, entry.parentID, 0)
	assert.False(t, entry.start.IsZero())
	assert.True(t, entry.stop.IsZero())
	assert.Equal(t, entry.tracker, tracker)
}

func TestEntryStartStop(t *testing.T) {
	tracker := New(1000)

	entryName := "test-entry"
	entry := tracker.Start(entryName)
	assert.EqualValues(t, 1, len(tracker.stack))
	assert.EqualValues(t, 1, tracker.stack[0])
	entry.Stop()

	assert.EqualValues(t, 1, len(tracker.entries))
	assert.EqualValues(t, 0, len(tracker.stack))

	assert.Equal(t, entry.name, entryName)
	assert.EqualValues(t, entry.id, 1)
	assert.EqualValues(t, entry.parentID, 0)
	assert.False(t, entry.start.IsZero())
	assert.False(t, entry.stop.IsZero())
	assert.Equal(t, entry.tracker, tracker)

	dur := entry.stop.Sub(entry.start)
	assert.GreaterOrEqual(t, int64(dur), int64(0))
}

func TestSubEntries(t *testing.T) {
	tracker := New(1000)

	entryName1 := "test-entry1"
	entryName2 := "test-entry2"
	entryName3 := "test-entry3"
	entry1 := tracker.Start(entryName1)
	entry2 := tracker.Start(entryName2)
	entry3 := tracker.Start(entryName3)

	assert.EqualValues(t, 3, len(tracker.stack))
	assert.EqualValues(t, 3, tracker.stack[2])

	entry3.Stop()

	assert.EqualValues(t, 2, len(tracker.stack))
	assert.EqualValues(t, 2, tracker.stack[1])

	entry2.Stop()

	assert.EqualValues(t, 1, len(tracker.stack))
	assert.EqualValues(t, 1, tracker.stack[0])

	entry1.Stop()

	assert.EqualValues(t, 3, len(tracker.entries))
	assert.EqualValues(t, 0, len(tracker.stack))

	assert.Equal(t, entry1.name, entryName1)
	assert.Equal(t, entry2.name, entryName2)
	assert.Equal(t, entry3.name, entryName3)

	assert.EqualValues(t, entry1.parentID, 0)
	assert.EqualValues(t, entry2.parentID, entry1.id)
	assert.EqualValues(t, entry3.parentID, entry2.id)
}

func TestEntryStopInvalidOrder(t *testing.T) {
	assert.Panics(t, func() {
		tracker := New(1000)

		entryName1 := "test-entry1"
		entryName2 := "test-entry2"
		entry1 := tracker.Start(entryName1)
		entry2 := tracker.Start(entryName2)

		// entry2 should be stopped before entry1
		entry1.Stop()
		entry2.Stop()
	})
}

func TestResultMissingEntryStop(t *testing.T) {
	assert.Panics(t, func() {
		tracker := New(1000)

		entryName := "test-entry"
		tracker.Start(entryName)
		tracker.Result()
	})
}

func TestResultOnlyOneEntry(t *testing.T) {
	tracker := New(1000)

	entryName := "test-entry"
	expected := tracker.Start(entryName)
	expected.Stop()
	res := tracker.Result()
	entry := res.Entries[0]
	assert.Equal(t, expected.id, entry.ID)
}
