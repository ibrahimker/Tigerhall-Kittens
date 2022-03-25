package logging

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	testValue = "value"
)

func TestLogInterceptor(t *testing.T) {
	var gotCtx context.Context

	// client test invoker
	inv := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		gotCtx = ctx
		return nil
	}
	// server test handler
	handler := func(ctx context.Context, _ interface{}) (interface{}, error) {
		gotCtx = ctx
		return nil, nil
	}

	log := NewLogger().WithFields(logrus.Fields{
		"test-key":   testValue,
		"ignore-key": "ignore",
	})

	ctx := ContextWithCorrelationID(context.Background())
	corr := CorrelationIDFromContext(ctx)
	callCtx := WithLogger(ctx, log)

	// Call client interceptor
	_ = UnaryClientInterceptor("ignore-key")(callCtx, "", nil, nil, nil, inv)

	if md := metautils.ExtractOutgoing(gotCtx); md.Get("x-log-test-key") != testValue {
		t.Error("missing log key", md)
	}
	if c := CorrelationIDFromContext(gotCtx); c != corr {
		t.Errorf("incorrect correlationID want %q got %q", corr, c)
	}
	outCtx := gotCtx

	inCtx := metautils.ExtractOutgoing(outCtx).Clone().ToIncoming(ctx)
	// Call server interceptor with clear
	_, err := UnaryServerInterceptor(true)(inCtx, nil, nil, handler)
	if err != nil {
		return
	}
	if lg := FromContext(gotCtx); lg.Data["test-key"] == testValue {
		t.Error("log key not cleared", lg.Data)
	}
	if lg := FromContext(gotCtx); lg.Data["ignore-key"] != nil {
		t.Error("ignore key forwarded", lg.Data)
	}
	if c := CorrelationIDFromContext(gotCtx); c != corr {
		t.Errorf("incorrect correlationID want %q got %q", corr, c)
	}

	inCtx = metautils.ExtractOutgoing(outCtx).Clone().ToIncoming(ctx)
	// Call server interceptor with clear
	_, err = UnaryServerInterceptor(false)(inCtx, nil, nil, handler)
	if err != nil {
		return
	}
	if lg := FromContext(gotCtx); lg.Data["test-key"] != testValue {
		t.Error("log key not loaded", lg.Data)
	}
	if lg := FromContext(gotCtx); lg.Data["ignore-key"] != nil {
		t.Error("ignore key forwarded", lg.Data)
	}
	if c := CorrelationIDFromContext(gotCtx); c != corr {
		t.Errorf("incorrect correlationID want %q got %q", corr, c)
	}
}
