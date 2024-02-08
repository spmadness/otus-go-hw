package grpc

import (
	"context"

	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/server"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateEvent(_ context.Context, event *pb.Event) (*pb.Result, error) {
	e := storage.Event{
		Title:       event.GetTitle(),
		DateStart:   event.GetDateStart(),
		DateEnd:     event.GetDateEnd(),
		Description: event.GetDescription(),
		UserID:      event.GetUserId(),
		DatePost:    event.GetDatePost(),
	}

	err := server.ValidateCreateEvent(e)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	err = s.storage.CreateEvent(e)
	if err != nil {
		return &pb.Result{}, err
	}
	return &pb.Result{}, nil
}

func (s *Server) UpdateEvent(_ context.Context, ur *pb.UpdateRequest) (*pb.Result, error) {
	id := ur.GetId().GetId()
	update := ur.GetEvent()

	e, err := s.storage.GetEvent(id)
	if err != nil {
		return nil, err
	}

	updateEventFields(&e, update)

	err = server.ValidateUpdateEvent(e)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	err = s.storage.UpdateEvent(id, e)
	if err != nil {
		return &pb.Result{}, err
	}
	return &pb.Result{}, nil
}

func (s *Server) DeleteEvent(_ context.Context, eventID *pb.EventId) (*pb.Result, error) {
	e := storage.Event{
		ID: eventID.GetId(),
	}

	err := server.ValidateDeleteEvent(e)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	err = s.storage.DeleteEvent(e.ID)
	if err != nil {
		return &pb.Result{}, err
	}
	return &pb.Result{}, nil
}

func (s *Server) ListEventDay(_ context.Context, in *pb.ListDate) (*pb.Result, error) {
	return listEvent(in.GetDateStart(), s.storage.ListEventDay)
}

func (s *Server) ListEventWeek(_ context.Context, in *pb.ListDate) (*pb.Result, error) {
	return listEvent(in.GetDateStart(), s.storage.ListEventWeek)
}

func (s *Server) ListEventMonth(_ context.Context, in *pb.ListDate) (*pb.Result, error) {
	return listEvent(in.GetDateStart(), s.storage.ListEventMonth)
}

func listEvent(date string, f func(date string) ([]storage.Event, error)) (*pb.Result, error) {
	lm := storage.ListEventValidation{
		DateStart: date,
	}

	result := &pb.Result{}

	err := server.ValidateListEvent(lm)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	events, err := f(lm.DateStart)
	if err != nil {
		return result, err
	}

	if len(events) == 0 {
		return result, nil
	}

	storageEvents := make([]*pb.Event, len(events))

	for i, event := range events {
		e := &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			DateStart:   event.DateStart,
			DateEnd:     event.DateEnd,
			Description: event.Description,
			UserId:      event.UserID,
			DatePost:    event.DatePost,
		}
		storageEvents[i] = e
	}

	result.Events = storageEvents

	return result, nil
}

func updateEventFields(e *storage.Event, changed *pb.Event) {
	if changed.Title != "" && changed.Title != e.Title {
		e.Title = changed.Title
	}
	if changed.DateStart != "" && changed.DateStart != e.DateStart {
		e.DateStart = changed.DateStart
	}
	if changed.DateEnd != "" && changed.DateEnd != e.DateEnd {
		e.DateEnd = changed.DateEnd
	}
	if changed.Description != "" && changed.Description != e.Description {
		e.Description = changed.Description
	}
	if changed.UserId != "" && changed.UserId != e.UserID {
		e.UserID = changed.UserId
	}
	if changed.DatePost != "" && changed.DatePost != e.DatePost {
		e.DatePost = changed.DatePost
	}
}
