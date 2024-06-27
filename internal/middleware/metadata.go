package middleware

import (
	"context"
	"inventory-service/internal/pkg/app"

	"google.golang.org/grpc"
)

type Metadata struct {
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context method to override the context
func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// Unary interceptor
func (u *Metadata) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var err error
		ctx, err = u.metadata(ctx)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// Stream interceptor
func (u *Metadata) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx, err := u.metadata(stream.Context())
		if err != nil {
			return err
		}

		wrappedStream := &wrappedStream{
			ServerStream: stream,
			ctx:          ctx,
		}

		return handler(srv, wrappedStream)
	}
}

func (u *Metadata) metadata(ctx context.Context) (context.Context, error) {
	ctx, err := app.GetMetadata(ctx)
	if err != nil {
		return ctx, err
	}

	ctx = app.SetMetadata(ctx)
	return ctx, nil
}
