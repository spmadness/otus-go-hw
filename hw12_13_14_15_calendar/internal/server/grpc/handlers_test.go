package grpc

import (
	"context"
	"testing"

	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ListHandlerFunc func(ctx context.Context, in *pb.ListDate) (*pb.Result, error)

var createUpdateCases = []struct {
	event *pb.Event
	err   bool
}{
	{
		event: &pb.Event{
			Title:       "Test Event",
			DateStart:   "2022-10-11 12:00:00",
			DateEnd:     "2022-10-12 13:00:00",
			Description: "Test Description",
			UserId:      "d5095366-ea13-4c9d-ae72-9c83d2d93040",
			DatePost:    "2022-10-12 13:00:00",
		},
		err: false,
	},
	{
		event: &pb.Event{
			DateStart:   "2022-10-11 12:00:00",
			DateEnd:     "2022-10-12 13:00:00",
			Description: "Test Description",
			UserId:      "d5095366-ea13-4c9d-ae72-9c83d2d93040",
			DatePost:    "2022-10-12 13:00:00",
		},
		err: true,
	},
	{
		event: &pb.Event{
			Title:       "Test Event",
			DateStart:   "2022-10-11 12:00:00",
			DateEnd:     "2022-10-12 13:00:00",
			Description: "Test Description",
			UserId:      "invalid-uuid",
			DatePost:    "2022-10-12 13:00:00",
		},
		err: true,
	},
	{
		event: &pb.Event{
			Title:       "Test Event",
			DateEnd:     "2022-10-12 13:00:00",
			Description: "Test Description",
			UserId:      "d5095366-ea13-4c9d-ae72-9c83d2d93040",
			DatePost:    "2022-10-12 13:00:00",
		},
		err: true,
	},
	{
		event: &pb.Event{
			Title:       "Test Event",
			DateStart:   "invalid-datetime",
			DateEnd:     "2022-10-12 13:00:00",
			Description: "Test Description",
			UserId:      "d5095366-ea13-4c9d-ae72-9c83d2d93040",
			DatePost:    "2022-10-12 13:00:00",
		},
		err: true,
	},
	{
		event: &pb.Event{
			Title:       "Test Event",
			DateStart:   "2022-10-11 12:00:00",
			Description: "Test Description",
			UserId:      "d5095366-ea13-4c9d-ae72-9c83d2d93040",
			DatePost:    "2022-10-12 13:00:00",
		},
		err: true,
	},
	{
		event: &pb.Event{
			Title:       "Test Event",
			DateStart:   "2022-10-11 12:00:00",
			DateEnd:     "invalid-datetime",
			Description: "Test Description",
			UserId:      "d5095366-ea13-4c9d-ae72-9c83d2d93040",
			DatePost:    "2022-10-12 13:00:00",
		},
		err: true,
	},
	{
		event: &pb.Event{
			Title:       "Test Event",
			DateStart:   "2022-10-11 12:00:00",
			DateEnd:     "2022-10-12 13:00:00",
			Description: "Test Description",
			UserId:      "d5095366-ea13-4c9d-ae72-9c83d2d93040",
			DatePost:    "invalid-datetime",
		},
		err: true,
	},
}

func TestCreateEventHandler(t *testing.T) {
	t.Run("handler validation", func(t *testing.T) {
		s := mocks.NewStorager(t)
		l := mocks.NewLogger(t)
		server := NewServer(l, s, 8080)

		for _, tc := range createUpdateCases {
			s.On("CreateEvent", mock.AnythingOfType("storage.Event")).Return(nil)

			_, err := server.CreateEvent(context.Background(), tc.event)
			if tc.err {
				assert.Error(t, err)
				s.AssertNotCalled(t, "CreateEvent")
				continue
			}
			assert.NoError(t, err)
			s.AssertCalled(t, "CreateEvent", mock.AnythingOfType("storage.Event"))
		}
	})
}

func TestUpdateEventHandler(t *testing.T) {
	t.Run("handler event validation", func(t *testing.T) {
		s := mocks.NewStorager(t)
		l := mocks.NewLogger(t)
		server := NewServer(l, s, 8080)

		for _, tc := range createUpdateCases {
			s.On("GetEvent", mock.AnythingOfType("string")).Return(storage.Event{
				ID: "eb0af540-6f23-4305-a719-fb65271fca1f",
			}, nil)

			s.On("UpdateEvent", mock.AnythingOfType("string"), mock.AnythingOfType("storage.Event")).Return(nil)

			_, err := server.UpdateEvent(context.Background(), &pb.UpdateRequest{
				Id:    &pb.EventId{Id: "eb0af540-6f23-4305-a719-fb65271fca1f"},
				Event: tc.event,
			})
			if tc.err {
				assert.Error(t, err)
				s.AssertNotCalled(t, "UpdateEvent")
				continue
			}
			assert.NoError(t, err)
			s.AssertCalled(t, "UpdateEvent", mock.AnythingOfType("string"), mock.AnythingOfType("storage.Event"))
		}
	})
}

func TestDeleteEventHandler(t *testing.T) {
	t.Run("handler validation", func(t *testing.T) {
		cases := []struct {
			val *pb.EventId
			err bool
		}{
			{&pb.EventId{Id: "eb0af540-6f23-4305-a719-fb65271fca1f"}, false},
			{&pb.EventId{Id: ""}, true},
			{&pb.EventId{Id: "test"}, true},
			{&pb.EventId{Id: "eb0af540-6f23-4305-a719"}, true},
		}

		s := mocks.NewStorager(t)
		l := mocks.NewLogger(t)
		server := NewServer(l, s, 8080)

		for _, tc := range cases {
			s.On("DeleteEvent", mock.AnythingOfType("string")).Return(nil)

			_, err := server.DeleteEvent(context.Background(), tc.val)
			if tc.err {
				assert.Error(t, err)
				s.AssertNotCalled(t, "DeleteEvent")
				continue
			}
			assert.NoError(t, err)
			s.AssertCalled(t, "DeleteEvent", mock.AnythingOfType("string"))
		}
	})
}

func TestListDayHandler(t *testing.T) {
	l := mocks.NewLogger(t)
	s := mocks.NewStorager(t)
	method := "ListEventDay"

	server := NewServer(l, s, 8080)

	t.Run("handler validation", func(t *testing.T) {
		listsHandlerTest(t, s, method, server.ListEventDay)
	})
}

func TestListWeekHandler(t *testing.T) {
	l := mocks.NewLogger(t)
	s := mocks.NewStorager(t)
	method := "ListEventWeek"

	server := NewServer(l, s, 8080)

	t.Run("handler validation", func(t *testing.T) {
		listsHandlerTest(t, s, method, server.ListEventWeek)
	})
}

func TestListMonthHandler(t *testing.T) {
	l := mocks.NewLogger(t)
	s := mocks.NewStorager(t)
	method := "ListEventMonth"

	server := NewServer(l, s, 8080)

	t.Run("handler validation", func(t *testing.T) {
		listsHandlerTest(t, s, method, server.ListEventMonth)
	})
}

func listsHandlerTest(t *testing.T, s *mocks.Storager, method string, f ListHandlerFunc) {
	t.Helper()

	cases := []struct {
		val string
		err bool
	}{
		{"2022-10-11", false},
		{"", true},
		{"2022-22-22", true},
		{"2022.12.11", true},
	}

	lm := pb.ListDate{}

	for _, tc := range cases {
		s.On(method, mock.AnythingOfType("string")).Return([]storage.Event{}, nil)

		lm.DateStart = tc.val

		_, err := f(context.Background(), &lm)
		if tc.err {
			assert.Error(t, err)
			s.AssertNotCalled(t, method)
			continue
		}
		assert.NoError(t, err)
		s.AssertCalled(t, method, mock.AnythingOfType("string"))
	}
}
