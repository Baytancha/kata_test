package tracer

import (
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	Configured trace.Tracer
)

func init() {
	Configured = otel.Tracer(os.Getenv("OTEL_SERVICE_NAME"))
}

func Tracer() trace.Tracer {
	return Configured
}
