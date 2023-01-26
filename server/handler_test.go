package server

import (
	"context"
	"crypto/rand"
	"net"
	"time"

	"github.com/AltMax/art-test/config"
	"github.com/AltMax/art-test/models"
	"github.com/AltMax/art-test/services"
	"github.com/AltMax/art-test/units/mocks"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func newMockGrpcConnAndListener(opts ...grpc.DialOption) (listener *bufconn.Listener, conn *grpc.ClientConn) {
	listener = bufconn.Listen(1024 * 1024)

	options := []grpc.DialOption{
		grpc.WithContextDialer(func(i context.Context, s string) (conn net.Conn, e error) {
			return listener.Dial()
		}),
		grpc.WithInsecure(),
	}
	options = append(options, opts...)
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		options...,
	)
	if err != nil {
		panic(err)
	}
	return
}

type handler struct {
	unitsMock         *mocks.Units
	unitServiceClient services.UnitServiceClient
}

func newTestHandler() *handler {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}

	unitsMock := &mocks.Units{}

	handler := &handler{
		unitsMock: unitsMock,
	}

	service := NewUnitService(unitsMock, 1*time.Second)
	listener, conn := newMockGrpcConnAndListener()
	srv := New(&conf)
	services.RegisterUnitServiceServer(srv, service)
	handler.unitServiceClient = services.NewUnitServiceClient(conn)
	go func() {
		if err := srv.Serve(listener); err != nil {
			panic(err)
		}
	}()
	return handler
}

func randomUnit() *models.Unit {
	buf := make([]byte, 50)
	rand.Read(buf)
	return &models.Unit{
		ID:        uuid.New().String(),
		Data:      buf,
		CreatedAt: time.Now().UTC(),
	}
}
