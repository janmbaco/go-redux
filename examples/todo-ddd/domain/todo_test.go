package domain

import (
	"testing"
	"time"
)

func TestTodo_Toggle(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name             string
		initialComplete  bool
		expectedComplete bool
	}{
		{"toggle incomplete to complete", false, true},
		{"toggle complete to incomplete", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo := Todo{
				ID:        1,
				Text:      "Test",
				Completed: tt.initialComplete,
				CreatedAt: now,
			}

			todo.Toggle()

			if todo.Completed != tt.expectedComplete {
				t.Errorf("Expected completed=%v, got %v", tt.expectedComplete, todo.Completed)
			}
		})
	}
}
