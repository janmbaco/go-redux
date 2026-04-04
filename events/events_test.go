package events_test

import (
	"sync/atomic"
	"testing"

	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-redux/v2/events"
)

func TestStateEventManagerPublish_ShouldNotifySubscriber_WhenSubscribed(t *testing.T) {
	// Arrange
	logger := logs.NewLogger()
	logger.Mute()
	manager := events.NewStateEventManager[int](logger)
	var calls atomic.Int32
	var received atomic.Int32
	callback := func(state int) {
		calls.Add(1)
		received.Store(int32(state))
	}
	if err := manager.Subscribe(&callback); err != nil {
		t.Fatalf("expected no subscribe error, got %v", err)
	}

	// Act
	manager.Publish(42)

	// Assert
	if calls.Load() != 1 {
		t.Fatalf("expected subscriber to be called once, got %d", calls.Load())
	}

	if received.Load() != 42 {
		t.Fatalf("expected state %d, got %d", 42, received.Load())
	}
}

func TestStateEventManagerPublish_ShouldNotNotifySubscriber_WhenUnsubscribed(t *testing.T) {
	// Arrange
	logger := logs.NewLogger()
	logger.Mute()
	manager := events.NewStateEventManager[int](logger)
	var calls atomic.Int32
	callback := func(state int) {
		calls.Add(1)
	}
	if err := manager.Subscribe(&callback); err != nil {
		t.Fatalf("expected no subscribe error, got %v", err)
	}
	if err := manager.Unsubscribe(&callback); err != nil {
		t.Fatalf("expected no unsubscribe error, got %v", err)
	}

	// Act
	manager.Publish(42)

	// Assert
	if calls.Load() != 0 {
		t.Fatalf("expected subscriber not to be called, got %d", calls.Load())
	}
}

func TestStateEventManagerSubscribe_ShouldIgnoreDuplicateCallback_WhenSubscribedTwice(t *testing.T) {
	// Arrange
	logger := logs.NewLogger()
	logger.Mute()
	manager := events.NewStateEventManager[int](logger)
	var calls atomic.Int32
	callback := func(state int) {
		calls.Add(1)
	}
	if err := manager.Subscribe(&callback); err != nil {
		t.Fatalf("expected no subscribe error, got %v", err)
	}
	if err := manager.Subscribe(&callback); err != nil {
		t.Fatalf("expected no duplicate subscribe error, got %v", err)
	}

	// Act
	manager.Publish(42)

	// Assert
	if calls.Load() != 1 {
		t.Fatalf("expected duplicate subscription to be ignored, got %d calls", calls.Load())
	}
}
