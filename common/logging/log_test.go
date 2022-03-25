package logging

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()

	if logger.Out != os.Stdout {
		t.Errorf("logger.Out is not os.Stdout")
	}
}

func TestNewTestLogger(t *testing.T) {
	logger := NewTestLogger()

	if len(logger.Data) != 1 {
		t.Errorf("len(logger.Data)(%d) does not equal 1", len(logger.Data))
	}
}

type testStruct struct {
	ID string
}

func TestNewHandlerLogger(t *testing.T) {
	request := &testStruct{ID: "test"}
	logger, ctx := NewHandlerLogger(context.Background(), NewTestLogger(), "handler-name-test", request)
	if len(logger.Data) != 4 {
		t.Errorf("len(logger.Data)(%d) does not equal 4", len(logger.Data))
	}
	if CorrelationIDFromContext(ctx) == "" {
		t.Fatalf("correlation-id is mandatory to have")
	}

	logger2, ctx2 := NewHandlerLogger(context.Background(), NewTestLogger(), "test", nil)
	if len(logger2.Data) != 3 {
		t.Errorf("len(logger2.Data)(%d) does not equal 3", len(logger2.Data))
	}
	if CorrelationIDFromContext(ctx2) == "" {
		t.Fatalf("correlation-id is mandatory to have")
	}
}

func TestNewServiceLogger(t *testing.T) {
	logger := NewServiceLogger(context.Background(), "test", logrus.Fields{
		"field-1": "field1",
	})

	if len(logger.Data) != 3 {
		t.Errorf("len(logger.Data)(%d) does not equal 3", len(logger.Data))
	}
}

func TestNewRepoLogger(t *testing.T) {
	logger := NewRepoLogger(context.Background(), "test", logrus.Fields{
		"field-1": "field1",
	})

	if len(logger.Data) != 3 {
		t.Errorf("len(logger.Data)(%d) does not equal 3", len(logger.Data))
	}
}

func TestWithError(t *testing.T) {
	logger := NewLogger().WithFields(logrus.Fields{
		"field1":      "value1",
		"field.key.2": "value2",
	})

	err := fmt.Errorf("This is an error")
	logger = WithError(err, logger)

	if len(logger.Data) != 4 {
		t.Errorf("len(logger.Data)(%d) does not equal 4", len(logger.Data))
	}

	v := logger.Data["field1"].(string)
	if v != "value1" {
		t.Errorf("logger.Data[field1](%s) does not equal 'value1'", v)
	}
	v = logger.Data["field.key.2"].(string)
	if v != "value2" {
		t.Errorf("logger.Data[field.key.2](%s) does not equal 'value2'", v)
	}
	if _, ok := logger.Data["error"].(error); !ok {
		t.Errorf("'error' key missing in logging data")
	}
	if _, ok := logger.Data["stacktrace"]; !ok {
		t.Errorf("'stacktrace' key missing in logging data")
	}
}

func TestContext(t *testing.T) {
	logr, hook := test.NewNullLogger()
	ctx := context.Background()
	tctx := WithLogger(ctx, logr)
	got := FromContext(tctx)
	if got == nil {
		t.Fatal("logger not loaded into context")
	}
	if got.Logger != logr {
		t.Error("retrieved logger does not match loaded")
	}
	got.Info("test")
	e := hook.LastEntry()
	if e.Data[correlationIDKey] == "" {
		t.Error("missing correlation id", e.Data)
	}
}
