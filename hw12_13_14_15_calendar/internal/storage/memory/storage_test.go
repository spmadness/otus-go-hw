package memorystorage

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("create success", func(t *testing.T) {
		s := New()

		e := storage.Event{
			DateStart: "2022-10-10 00:02:15",
		}

		err := s.CreateEvent(e)
		require.NoError(t, err)
	})

	t.Run("create fail: date busy", func(t *testing.T) {
		id := "1"

		s := &Storage{
			eventsByID: map[string]*storage.Event{
				id: {ID: id, DateStart: "2022-10-10 00:02:15"},
			},
			eventsByDateStart: map[string]*storage.Event{
				"2022-10-10 00:02:15": {ID: id, DateStart: "2022-10-10 00:02:15"},
			},
			eventsByDay: map[int64][]storage.Event{
				1665360000: {
					{ID: id, DateStart: "2022-10-10 00:02:15"},
				},
			},
		}

		e := storage.Event{
			DateStart: "2022-10-10 00:02:15",
		}

		err := s.CreateEvent(e)
		require.ErrorIs(t, err, ErrDateBusy)
	})

	t.Run("update success", func(t *testing.T) {
		id := "1"

		s := &Storage{
			eventsByID: map[string]*storage.Event{
				id: {ID: id, DateStart: "2022-10-10 00:02:15"},
			},
			eventsByDateStart: map[string]*storage.Event{
				"2022-10-10 00:02:15": {ID: id, DateStart: "2022-10-10 00:02:15"},
			},
			eventsByDay: map[int64][]storage.Event{
				1665360000: {
					{ID: id, DateStart: "2022-10-10 00:02:15"},
				},
			},
		}

		datetime := "2023-01-01 10:00:00"
		updateEvent := storage.Event{DateStart: datetime}

		err := s.UpdateEvent(id, updateEvent)
		require.NoError(t, err)

		val, err := s.Event(id)
		require.NoError(t, err)

		require.Equal(t, datetime, val.DateStart)
	})

	t.Run("update fail: date busy", func(t *testing.T) {
		id := "1"

		s := &Storage{
			eventsByID: map[string]*storage.Event{
				id: {ID: id, DateStart: "2022-10-10 00:02:15"},
			},
			eventsByDateStart: map[string]*storage.Event{
				"2022-10-10 00:02:15": {ID: id, DateStart: "2022-10-10 00:02:15"},
			},
			eventsByDay: map[int64][]storage.Event{
				1665360000: {
					{ID: id, DateStart: "2022-10-10 00:02:15"},
				},
			},
		}

		updateEvent := storage.Event{DateStart: "2022-10-10 00:02:15"}
		err := s.UpdateEvent("1", updateEvent)
		require.Error(t, err, ErrDateBusy)
	})

	t.Run("update fail: no event", func(t *testing.T) {
		id := "1"

		s := &Storage{
			eventsByID: map[string]*storage.Event{
				id: {ID: id, DateStart: "2022-10-10 00:02:15"},
			},
		}

		updateEvent := storage.Event{DateStart: "2023-01-01 10:00:00"}
		err := s.UpdateEvent("2", updateEvent)
		require.Error(t, err, ErrEventNotExist)
	})

	t.Run("delete success", func(t *testing.T) {
		id := "1"
		dayTimeUnix := int64(1665360000)

		s := &Storage{
			eventsByID: map[string]*storage.Event{
				id: {ID: id, DateStart: "2022-10-10 00:02:15"},
			},
			eventsByDateStart: map[string]*storage.Event{
				"2022-10-10 00:02:15": {ID: id, DateStart: "2022-10-10 00:02:15"},
			},
			eventsByDay: map[int64][]storage.Event{
				dayTimeUnix: {
					{ID: "1", DateStart: "2022-10-10 00:02:15"},
				},
			},
		}
		err := s.DeleteEvent(id)
		require.NoError(t, err)

		_, err = s.Event(id)
		require.Error(t, err, ErrEventNotExist)
		require.True(t, len(s.eventsByID) == 0)
		require.True(t, len(s.eventsByDateStart) == 0)
		require.True(t, len(s.eventsByDay[dayTimeUnix]) == 0)
	})

	t.Run("delete fail: no event", func(t *testing.T) {
		id := "1"
		dayTimeUnix := int64(1665360000)

		s := &Storage{
			eventsByID: map[string]*storage.Event{
				id: {ID: id, DateStart: "2022-10-10 00:02:15"},
			},
			eventsByDateStart: map[string]*storage.Event{
				"2022-10-10 00:02:15": {ID: id, DateStart: "2022-10-10 00:02:15"},
			},
			eventsByDay: map[int64][]storage.Event{
				dayTimeUnix: {
					{ID: "1", DateStart: "2022-10-10 00:02:15"},
				},
			},
		}

		id = "2"
		err := s.DeleteEvent(id)
		require.Error(t, err, ErrEventNotExist)
	})
}

func TestStorageReadMethods(t *testing.T) {
	t.Run("list day success", func(t *testing.T) {
		s := &Storage{
			eventsByDay: map[int64][]storage.Event{
				1665360000: {
					{ID: "1", DateStart: "2022-10-10 00:02:15"},
					{ID: "2", DateStart: "2022-10-10 00:04:15"},
				},
				1665522000: {
					{ID: "3", DateStart: "2022-10-12 00:02:15"},
				},
			},
		}

		events, err := s.ListEventDay("2022-10-10")
		require.NoError(t, err)

		require.Truef(t, len(events) == 2, "expected length: %d, actual: %d", 2, len(events))
	})

	t.Run("list day success: no events", func(t *testing.T) {
		s := &Storage{
			eventsByDay: map[int64][]storage.Event{
				1665360000: {
					{ID: "1", DateStart: "2022-10-10 00:02:15"},
					{ID: "2", DateStart: "2022-10-10 00:04:15"},
				},
				1665522000: {
					{ID: "3", DateStart: "2022-10-12 00:02:15"},
				},
			},
		}

		events, err := s.ListEventDay("2022-11-10")
		require.NoError(t, err)

		require.Truef(t, len(events) == 0, "expected length: %d, actual: %d", 0, len(events))
	})

	t.Run("list week success", func(t *testing.T) {
		s := &Storage{
			eventsByDay: map[int64][]storage.Event{
				1665360000: {
					{ID: "1", DateStart: "2022-10-10 00:02:15"},
					{ID: "2", DateStart: "2022-10-10 00:04:15"},
				},
				1665522000: {
					{ID: "3", DateStart: "2022-10-12 00:02:15"},
				},
				1666051200: {
					{ID: "4", DateStart: "2022-10-18 00:02:15"},
				},
			},
		}

		events, err := s.ListEventWeek("2022-10-10")
		require.NoError(t, err)

		require.Truef(t, len(events) == 3, "expected length: %d, actual: %d", 3, len(events))
	})

	t.Run("list week success: no events", func(t *testing.T) {
		s := &Storage{
			eventsByDay: map[int64][]storage.Event{
				1665360000: {
					{ID: "1", DateStart: "2022-10-10 00:02:15"},
					{ID: "2", DateStart: "2022-10-10 00:04:15"},
				},
				1665522000: {
					{ID: "3", DateStart: "2022-10-12 00:02:15"},
				},
				1666051200: {
					{ID: "4", DateStart: "2022-10-18 00:02:15"},
				},
			},
		}

		events, err := s.ListEventWeek("2022-11-10")
		require.NoError(t, err)

		require.Truef(t, len(events) == 0, "expected length: %d, actual: %d", 0, len(events))
	})

	t.Run("list month success", func(t *testing.T) {
		s := &Storage{
			eventsByDay: map[int64][]storage.Event{
				1665360000: {
					{ID: "1", DateStart: "2022-10-10 00:02:15"},
					{ID: "2", DateStart: "2022-10-10 00:04:15"},
				},
				1665522000: {
					{ID: "3", DateStart: "2022-10-12 00:02:15"},
				},
				1667952000: {
					{ID: "4", DateStart: "2022-11-09 00:02:15"},
				},
				1668038400: {
					{ID: "5", DateStart: "2022-11-10 00:00:00"},
				},
			},
		}

		events, err := s.ListEventMonth("2022-10-10")
		require.NoError(t, err)

		require.Truef(t, len(events) == 4, "expected length: %d, actual: %d", 4, len(events))
	})

	t.Run("list month success: no events", func(t *testing.T) {
		s := &Storage{
			eventsByDay: map[int64][]storage.Event{
				1665360000: {
					{ID: "1", DateStart: "2022-10-10 00:02:15"},
					{ID: "2", DateStart: "2022-10-10 00:04:15"},
				},
				1665522000: {
					{ID: "3", DateStart: "2022-10-12 00:02:15"},
				},
				1667952000: {
					{ID: "4", DateStart: "2022-11-09 00:02:15"},
				},
				1668038400: {
					{ID: "5", DateStart: "2022-11-10 00:00:00"},
				},
			},
		}

		events, err := s.ListEventMonth("2022-12-10")
		require.NoError(t, err)

		require.Truef(t, len(events) == 0, "expected length: %d, actual: %d", 0, len(events))
	})
}

func TestStorageMethodsConcurrency(t *testing.T) {
	t.Run("create concurrency success", func(t *testing.T) {
		s := New()

		numGoroutines := 10

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(i int) {
				defer wg.Done()

				dateStartTime := time.Date(2022, time.January, 1, 0, 0, i, 0, time.UTC)

				_ = s.CreateEvent(storage.Event{
					Title:       "Concurrent Event",
					DateStart:   dateStartTime.Format(time.RFC3339),
					DateEnd:     dateStartTime.Add(2 * time.Hour).Format(time.RFC3339),
					Description: "Concurrent Event Description",
					UserID:      1,
					DatePost:    time.Now().Format(time.RFC3339),
				})
			}(i)
		}

		wg.Wait()

		require.Truef(t,
			numGoroutines == len(s.eventsByID),
			"events expected: %d, actual: %d", numGoroutines, len(s.eventsByID))
	})

	t.Run("delete concurrency test", func(t *testing.T) {
		s := &Storage{
			eventsByID: map[string]*storage.Event{
				"1": {ID: "1", DateStart: "2022-10-10 00:02:15"},
				"2": {ID: "2", DateStart: "2022-10-10 00:04:15"},
			},
		}

		numGoroutines := 10

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		var errCnt, noErrCnt int

		for i := 0; i < numGoroutines; i++ {
			go func(i int) {
				defer wg.Done()

				err := s.DeleteEvent(strconv.Itoa(i))

				s.mu.Lock()
				if err == nil {
					noErrCnt++
				} else {
					errCnt++
				}
				s.mu.Unlock()
			}(i)
		}

		wg.Wait()

		require.Truef(t, noErrCnt == 2, "no errors expected: %d, actual: %d", 2, noErrCnt)
		require.Truef(t, errCnt == 8, "errors expected: %d, actual: %d", 8, noErrCnt)
	})

	t.Run("list concurrency test", func(t *testing.T) {
		s := &Storage{
			eventsByDay: map[int64][]storage.Event{
				1665360000: {
					{ID: "1", DateStart: "2022-10-10 00:02:15"},
					{ID: "2", DateStart: "2022-10-10 00:04:15"},
				},
			},
		}

		numGoroutines := 10

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		results := make([][]storage.Event, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(i int) {
				defer wg.Done()
				result, _ := s.ListEventDay("2022-11-10")
				results[i] = result
			}(i)
		}

		wg.Wait()

		require.Truef(t,
			numGoroutines == len(results),
			"events lists expected: %d, actual: %d", numGoroutines, len(results))
	})
}
