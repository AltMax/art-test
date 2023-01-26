package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AltMax/art-test/config"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcValidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func recoveryHandler(ctx context.Context, p interface{}) (err error) {
	log.Error().
		Str("panic", fmt.Sprintf("%+v", p)).
		Str("ctx", fmt.Sprintf("%+v", ctx)).
		Msg("PANIC")
	if err, ok := p.(error); ok {
		return fmt.Errorf("panic: %w", err)
	}
	return fmt.Errorf("panic: %+v", p)
}

// ErrorToInternalErrorMiddleware converts errors to google grpc errors with code Internsl
func ErrorToInternalErrorMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (result interface{}, err error) {

	result, err = handler(ctx, req)

	_, ok := status.FromError(err)
	if ok {
		return
	}

	err = status.Error(codes.Internal, "")

	return
}

func New(conf *config.Config, middlewares ...grpc.UnaryServerInterceptor) *grpc.Server {
	interceptors := []grpc.UnaryServerInterceptor{
		grpcRecovery.UnaryServerInterceptor(grpcRecovery.WithRecoveryHandlerContext(recoveryHandler)),
		ErrorToInternalErrorMiddleware,
		logIncomingRequestsMiddleware,
		grpcValidator.UnaryServerInterceptor(),
	}

	interceptors = append(interceptors, middlewares...)
	return grpc.NewServer(grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(interceptors...)))
}

func logIncomingRequestsMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	requestJSON, _ := json.Marshal(req)
	result, err := handler(ctx, req)
	responseJSON, _ := json.Marshal(result)
	var logEvent *zerolog.Event
	if err != nil {
		logEvent = log.Error().Str("error", fmt.Sprintf("%+v", err))
	} else {
		logEvent = log.Info()
	}
	logEvent.
		Dur("duration", time.Since(start)).
		RawJSON("json_response", responseJSON).
		RawJSON("json_request", requestJSON).
		Str("url", info.FullMethod).
		Str("ctx", fmt.Sprintf("%+v", ctx)).
		Msg("complete")

	return result, err
}
