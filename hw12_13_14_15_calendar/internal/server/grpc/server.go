package grpc

import (
	"fmt"
	"net"

	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedEventServiceServer

	logger  Logger
	storage app.Storager
	server  *grpc.Server
	port    int
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func NewServer(logger Logger, storage app.Storager, port int) *Server {
	return &Server{
		logger:  logger,
		storage: storage,
		port:    port,
	}
}

func (s *Server) Start() {
	lsn, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.logger.Error(err.Error())
	}

	s.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryServerRequestLoggingInterceptor(s.logger),
		),
	)
	pb.RegisterEventServiceServer(s.server, s)

	s.logger.Info(fmt.Sprintf("starting grpc server on %s", lsn.Addr().String()))

	if err = s.server.Serve(lsn); err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *Server) Stop() {
	s.server.Stop()
}
