package internalhttp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var createUpdateCases = []struct {
	event storage.Event
	err   bool
}{
	{
		event: storage.Event{
			Title:       "Test Event",
			DateStart:   "2022-10-11 12:00:00",
			DateEnd:     "2022-10-12 13:00:00",
			Description: "Test Description",
			UserID:      "d5095366-ea13-4c9d-ae72-9c83d2d93040",
			DatePost:    "2022-10-12 13:00:00",
		},
		err: false,
	},
	{
		event: storage.Event{
			DateStart:   "2022-10-11 12:00:00",
			DateEnd:     "2022-10-12 13:00:00",
			Description: "Test Description",
			UserID:      "d5095366-ea13-4c9d-ae72-9c83d2d93040",
			DatePost:    "2022-10-12 13:00:00",
		},
		err: true,
	},
	{
		event: storage.Event{
			Title:       "Test Event",
			DateStart:   "2022-10-11 12:00:00",
			DateEnd:     "2022-10-12 13:00:00",
			Description: "Test Description",
			UserID:      "invalid-uuid",
			DatePost:    "2022-10-12 13:00:00",
		},
		err: true,
	},
	{
		event: storage.Event{
			Title:       "Test Event",
			DateEnd:     "2022-10-12 13:00:00",
			Description: "Test Description",
			UserID:      "d5095366-ea13-4c9d-ae72-9c83d2d93040",
			DatePost:    "2022-10-12 13:00:00",
		},
		err: true,
	},
	{
		event: storage.Event{
			Title:       "Test Event",
			DateStart:   "invalid-datetime",
			DateEnd:     "2022-10-12 13:00:00",
			Description: "Test Description",
			UserID:      "d5095366-ea13-4c9d-ae72-9c83d2d93040",
			DatePost:    "2022-10-12 13:00:00",
		},
		err: true,
	},
	{
		event: storage.Event{
			Title:       "Test Event",
			DateStart:   "2022-10-11 12:00:00",
			Description: "Test Description",
			UserID:      "d5095366-ea13-4c9d-ae72-9c83d2d93040",
			DatePost:    "2022-10-12 13:00:00",
		},
		err: true,
	},
	{
		event: storage.Event{
			Title:       "Test Event",
			DateStart:   "2022-10-11 12:00:00",
			DateEnd:     "invalid-datetime",
			Description: "Test Description",
			UserID:      "d5095366-ea13-4c9d-ae72-9c83d2d93040",
			DatePost:    "2022-10-12 13:00:00",
		},
		err: true,
	},
	{
		event: storage.Event{
			Title:       "Test Event",
			DateStart:   "2022-10-11 12:00:00",
			DateEnd:     "2022-10-12 13:00:00",
			Description: "Test Description",
			UserID:      "d5095366-ea13-4c9d-ae72-9c83d2d93040",
			DatePost:    "invalid-datetime",
		},
		err: true,
	},
}

func TestMiddleware(t *testing.T) {
	t.Run("location http method test", func(t *testing.T) {
		cases := []struct {
			location string
			method   string
			err      bool
		}{
			{LocationCreate, http.MethodPost, false},
			{LocationCreate, http.MethodPut, true},
			{LocationCreate, http.MethodGet, true},
			{LocationCreate, http.MethodDelete, true},

			{LocationUpdate, http.MethodPost, true},
			{LocationUpdate, http.MethodPut, false},
			{LocationUpdate, http.MethodGet, true},
			{LocationUpdate, http.MethodDelete, true},

			{LocationDelete, http.MethodPost, true},
			{LocationDelete, http.MethodPut, true},
			{LocationDelete, http.MethodGet, true},
			{LocationDelete, http.MethodDelete, false},

			{LocationListDay, http.MethodPost, true},
			{LocationListDay, http.MethodPut, true},
			{LocationListDay, http.MethodGet, false},
			{LocationListDay, http.MethodDelete, true},

			{LocationListWeek, http.MethodPost, true},
			{LocationListWeek, http.MethodPut, true},
			{LocationListWeek, http.MethodGet, false},
			{LocationListWeek, http.MethodDelete, true},

			{LocationListMonth, http.MethodPost, true},
			{LocationListMonth, http.MethodPut, true},
			{LocationListMonth, http.MethodGet, false},
			{LocationListMonth, http.MethodDelete, true},
		}

		for _, tc := range cases {
			r := httptest.NewRequest(tc.method, "/"+tc.location, strings.NewReader(""))
			r.Header.Set("Content-Type", "application/json")

			mw := Middleware{}
			s := mocks.NewStorager(t)

			mux := NewMux(s)
			handler := MiddlewareChain(mw.requestValidatorMiddleware)(mux)

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			if tc.err {
				assert.Equalf(t, http.StatusMethodNotAllowed, w.Code, "code expected: %d, actual: %d",
					http.StatusMethodNotAllowed, w.Code)
				continue
			}
			assert.NotEqual(t, http.StatusMethodNotAllowed, w.Code)
		}
	})

	t.Run("content-type header test", func(t *testing.T) {
		cases := []struct {
			val string
			err bool
		}{
			{"application/json", false},
			{"", true},
			{"application/octet-stream", true},
			{"application/x-www-form-urlencoded", true},
		}

		requestBody := `{"ID": "9723a4b7-4c61-4ae5-97c6-6bf536badf48"}`

		for _, tc := range cases {
			r := httptest.NewRequest(http.MethodDelete, "/delete", strings.NewReader(requestBody))
			r.Header.Set("Content-Type", tc.val)

			mw := Middleware{}
			s := mocks.NewStorager(t)
			if !tc.err {
				s.On("DeleteEvent", mock.AnythingOfType("string")).Return(nil)
			}

			mux := NewMux(s)
			handler := MiddlewareChain(mw.requestValidatorMiddleware)(mux)

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			if tc.err {
				assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
				s.AssertNotCalled(t, "DeleteEvent")
				continue
			}
			assert.Equal(t, http.StatusOK, w.Code)
			s.AssertCalled(t, "DeleteEvent", mock.AnythingOfType("string"))
		}
	})
}

func TestCreateEventHandler(t *testing.T) {
	t.Run("handler validation", func(t *testing.T) {
		s := mocks.NewStorager(t)

		for _, tc := range createUpdateCases {
			jData, err := json.Marshal(tc.event)
			if err != nil {
				t.Errorf("json encode error")
			}

			r := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(string(jData)))
			w := httptest.NewRecorder()

			s.On("CreateEvent", mock.AnythingOfType("storage.Event")).Return(nil)

			_, err = createEvent(w, r, s)
			if tc.err {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Error(t, err)
				s.AssertNotCalled(t, "CreateEvent")
				continue
			}
			assert.Equal(t, http.StatusOK, w.Code)
			assert.NoError(t, err)
			s.AssertCalled(t, "CreateEvent", mock.AnythingOfType("storage.Event"))
		}
	})
}

func TestUpdateEventHandler(t *testing.T) {
	t.Run("handler validation", func(t *testing.T) {
		s := mocks.NewStorager(t)

		for _, tc := range createUpdateCases {
			jData, err := json.Marshal(&UpdateRequest{
				ID:    "eb0af540-6f23-4305-a719-fb65271fca1f",
				Event: tc.event,
			})
			if err != nil {
				t.Errorf("json encode error")
			}

			r := httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(string(jData)))
			w := httptest.NewRecorder()

			s.On("GetEvent", mock.AnythingOfType("string")).Return(storage.Event{
				ID: "eb0af540-6f23-4305-a719-fb65271fca1f",
			}, nil)

			s.On("UpdateEvent", mock.AnythingOfType("string"), mock.AnythingOfType("storage.Event")).Return(nil)

			_, err = updateEvent(w, r, s)
			if tc.err {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Error(t, err)
				s.AssertNotCalled(t, "UpdateEvent")
				continue
			}
			assert.Equal(t, http.StatusOK, w.Code)
			assert.NoError(t, err)
			s.AssertCalled(t, "UpdateEvent", mock.AnythingOfType("string"), mock.AnythingOfType("storage.Event"))
		}
	})
}

func TestDeleteEventHandler(t *testing.T) {
	t.Run("handler validation", func(t *testing.T) {
		cases := []struct {
			val string
			err bool
		}{
			{"9723a4b7-4c61-4ae5-97c6-6bf536badf48", false},
			{"9723a4b7-4c61-4ae5-97c6", true},
			{"", true},
			{"string", true},
		}
		s := mocks.NewStorager(t)

		for _, tc := range cases {
			requestBody := "{\"ID\": \"" + tc.val + "\"}"

			r := httptest.NewRequest(http.MethodDelete, "/delete", strings.NewReader(requestBody))
			w := httptest.NewRecorder()

			s.On("DeleteEvent", mock.AnythingOfType("string")).Return(nil)

			_, err := deleteEvent(w, r, s)

			if tc.err {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Error(t, err)
				s.AssertNotCalled(t, "DeleteEvent")
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
				assert.NoError(t, err)
				s.AssertCalled(t, "DeleteEvent", mock.AnythingOfType("string"))
			}
		}
	})
}

func TestListDayHandler(t *testing.T) {
	t.Run("handler validation", func(t *testing.T) {
		listsHandlerTest(t, "ListEventDay", listEventDay)
	})
}

func TestListWeekHandler(t *testing.T) {
	t.Run("handler validation", func(t *testing.T) {
		listsHandlerTest(t, "ListEventWeek", listEventWeek)
	})
}

func TestListMonthHandler(t *testing.T) {
	t.Run("handler validation", func(t *testing.T) {
		listsHandlerTest(t, "ListEventMonth", listEventMonth)
	})
}

func listsHandlerTest(t *testing.T, method string, f HandlerFunc) {
	t.Helper()

	cases := []struct {
		val string
		err bool
	}{
		{"2022-10-11", false},
		{"", true},
		{"2022-22-22", true},
		{"2022.10.11", true},
	}

	lm := storage.ListEventValidation{}

	s := mocks.NewStorager(t)

	for _, tc := range cases {
		lm.DateStart = tc.val

		jData, err := json.Marshal(lm)
		if err != nil {
			t.Errorf("json encode error")
		}
		r := httptest.NewRequest(http.MethodGet, "/"+LocationListDay, strings.NewReader(string(jData)))
		w := httptest.NewRecorder()

		s.On(method, mock.AnythingOfType("string")).Return([]storage.Event{}, nil)

		_, err = f(w, r, s)
		if tc.err {
			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Error(t, err)
			s.AssertNotCalled(t, method)
			continue
		}
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, err)
		s.AssertCalled(t, method, mock.AnythingOfType("string"))
	}
}
