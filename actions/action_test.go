package actions_test

import (
	"testing"

	"github.com/janmbaco/go-redux/v2/actions"
)

func TestNewAction_ShouldRegisterActionByType_WhenActionIsCreated(t *testing.T) {
	// Arrange
	actionType := "test/register"

	// Act
	action := actions.NewAction[int](actionType)

	// Assert
	registered, ok := actions.GetActionByType(actionType)
	if !ok {
		t.Fatalf("expected action %q to be registered", actionType)
	}

	if registered != action {
		t.Fatalf("expected registered action to match created action")
	}
}

func TestActionWith_ShouldCreateActionInstanceWithPayload_WhenPayloadIsProvided(t *testing.T) {
	// Arrange
	action := actions.NewAction[int]("test/payload")

	// Act
	instance := action.With(7)

	// Assert
	if instance.Type() != "test/payload" {
		t.Fatalf("expected action type %q, got %q", "test/payload", instance.Type())
	}

	if instance.GetPayload() != 7 {
		t.Fatalf("expected payload %d, got %d", 7, instance.GetPayload())
	}

	if instance.GetAction() != action {
		t.Fatalf("expected action instance to keep a reference to the original action")
	}
}
