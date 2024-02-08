package internalhttp

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/app"
)

type Response struct {
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request, s app.Storager) (interface{}, error)

const (
	LocationCreate    = "create"
	LocationUpdate    = "update"
	LocationDelete    = "delete"
	LocationListDay   = "list-day"
	LocationListWeek  = "list-week"
	LocationListMonth = "list-month"
)

func NewMux(s app.Storager) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/"+LocationCreate, handleRequest(createEvent, s))

	mux.Handle("/"+LocationUpdate, handleRequest(updateEvent, s))

	mux.Handle("/"+LocationDelete, handleRequest(deleteEvent, s))

	mux.Handle("/"+LocationListDay, handleRequest(listEventDay, s))

	mux.Handle("/"+LocationListWeek, handleRequest(listEventWeek, s))

	mux.Handle("/"+LocationListMonth, handleRequest(listEventMonth, s))

	return mux
}

func handleRequest(handler HandlerFunc, s app.Storager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := Response{}

		data, err := handler(w, r, s)
		if err != nil {
			resp.Error = err.Error()
		}
		if data != nil {
			resp.Data = data
		}

		writeJSONResponse(w, resp)
	})
}

func writeJSONResponse(w http.ResponseWriter, resp Response) {
	jData, err := json.Marshal(resp)
	if err != nil {
		log.Printf("json response marshal error")
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jData)
	if err != nil {
		log.Printf("json write response failed")
	}
}
