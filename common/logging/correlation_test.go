package logging

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

const (
	testCorrIDKey = "1111-2222-3333"
)

func TestFromContext(t *testing.T) {
	md := metadata.Pairs(correlationIDKey, testCorrIDKey)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	logger := FromContext(ctx)

	cid := logger.Data[correlationIDKey].(string)
	if cid != testCorrIDKey {
		t.Errorf("correlation id in logger from context is not expected '1111-2222-3333': %s", cid)
	}
}

func TestWithContext(t *testing.T) {
	md := metadata.Pairs(correlationIDKey, testCorrIDKey)
	ctx := metadata.NewIncomingContext(context.Background(), md)
	logger := WithContext(ctx, NewLogger())
	cid := logger.Data[correlationIDKey].(string)
	if cid != testCorrIDKey {
		t.Errorf("correlation id in logger from context is not expected '1111-2222-3333': %s", cid)
	}
}

func TestCorrelationIDFromContext(t *testing.T) {
	t.Run("WithCorrelationIDInContextMetadata", func(t *testing.T) {
		md := metadata.Pairs(correlationIDKey, testCorrIDKey)
		ctx := metadata.NewIncomingContext(context.Background(), md)
		cid := CorrelationIDFromContext(ctx)
		if cid != testCorrIDKey {
			t.Errorf("correlation id in context is not expected '1111-2222-3333': %s", cid)
		}
	})

	t.Run("WithoutCorrelationIDInContextMetadata", func(t *testing.T) {
		md := metadata.Pairs("not-correlation-key", "1234")
		ctx := metadata.NewIncomingContext(context.Background(), md)

		cid := CorrelationIDFromContext(ctx)
		if cid != "" {
			t.Errorf("correlation id in context is not empty: %s", cid)
		}
	})

	t.Run("WithoutMetadataInContext", func(t *testing.T) {
		ctx := ContextWithCorrelationID(context.Background())
		cid := CorrelationIDFromContext(ctx)
		if cid == "" {
			t.Errorf("correlation id in context is empty")
		}
	})
}

func TestContextWithCorrelationID(t *testing.T) {
	t.Run("WithCorrelationIDInContextMetadata", func(t *testing.T) {
		md := metadata.Pairs(correlationIDKey, testCorrIDKey)
		ctx := metadata.NewIncomingContext(context.Background(), md)
		ctx = ContextWithCorrelationID(ctx)
		cid := CorrelationIDFromContext(ctx)
		if cid != testCorrIDKey {
			t.Errorf("correlation id in context is not expected '1111-2222-3333': %s", cid)
		}
	})

	t.Run("WithoutCorrelationIDInContextMetadata", func(t *testing.T) {
		md := metadata.Pairs("not-correlation-key", "1234")
		ctx := metadata.NewIncomingContext(context.Background(), md)
		ctx = ContextWithCorrelationID(ctx)
		cid := CorrelationIDFromContext(ctx)
		if cid == "" {
			t.Errorf("correlation id in context empty")
		}
	})

	t.Run("WithoutMetadataInContext", func(t *testing.T) {
		ctx := ContextWithCorrelationID(context.Background())
		cid := CorrelationIDFromContext(ctx)
		if cid == "" {
			t.Errorf("correlation id in context empty")
		}
	})
}
