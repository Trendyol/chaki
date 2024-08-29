package otelclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var (
	traceID = trace.TraceID([16]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10})
	spanID  = trace.SpanID([8]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77})
)

func buildTestSpanContext() context.Context {
	parentctx, _ := noop.NewTracerProvider().Tracer("test tracer").Start(context.Background(), "test span")
	parentctx = trace.ContextWithRemoteSpanContext(parentctx, trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
	}))

	return parentctx
}

func Test_newClientRoundTripperWrapper(t *testing.T) {
	t.Run("it should set otel information to request header via round tripper", func(t *testing.T) {
		// Given
		wrapper := newClientRoundTripperWrapper(
			noop.NewTracerProvider(),
			propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
		)

		rt := wrapper(roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			// Mock client request with dummy response
			// and assert That trace header is set.
			assert.NotEmpty(t, req.Header.Get("Traceparent"))
			return &http.Response{}, nil
		}))

		ctx := buildTestSpanContext()

		cl := http.DefaultClient
		cl.Transport = rt

		req, _ := http.NewRequest("GET", "http://a.com", nil)

		// When
		_, err := cl.Do(req.WithContext(ctx)) //nolint:bodyclose

		// Then
		assert.NoError(t, err)
	})
}

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (rt roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}
