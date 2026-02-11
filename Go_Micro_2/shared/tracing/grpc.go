package tracing

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

func WithTracingInterceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.StatsHandler(newServerHandler()),
	}
}

func DialOptionWithTracing() []grpc.DialOption{
	return []grpc.DialOption{
		grpc.WithStatsHandler(newClientHandler()),
	}
}

func newServerHandler() stats.Handler {
	return otelgrpc.NewServerHandler(otelgrpc.WithTraceProvider(otel.GetTracerProvider))
}

func newClientHandler() stats.Handler {
	return otelgrpc.NewClientHandler(otelgrpc.WithTraceProvider(otel.GetTracerProvider))
}