package domain

import (
	"testing"
	"time"
)

func TestTodoService_AddTodo(t *testing.T) {
	service := NewTodoService()

	tests := []struct {
		name         string
		text         string
		nextID       int
		expectedID   int
		expectedNext int
	}{
		{"add first todo", "First task", 1, 1, 2},
		{"add second todo", "Second task", 5, 5, 6},
		{"empty text allowed", "", 10, 10, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todos := []*Todo{}

			result, newNextID := service.AddTodo(todos, tt.nextID, tt.text)

			if len(result) != 1 {
				t.Fatalf("Expected 1 todo, got %d", len(result))
			}

			newTodo := result[0]
			if newTodo.ID != tt.expectedID {
				t.Errorf("Expected ID %d, got %d", tt.expectedID, newTodo.ID)
			}

			if newTodo.Text != tt.text {
				t.Errorf("Expected text %q, got %q", tt.text, newTodo.Text)
			}

			if newTodo.Completed {
				t.Error("Expected new todo to be incomplete")
			}

			if newNextID != tt.expectedNext {
				t.Errorf("Expected NextID %d, got %d", tt.expectedNext, newNextID)
			}
		})
	}
}

func TestTodoService_ToggleTodo(t *testing.T) {
	service := NewTodoService()
	now := time.Now()

	tests := []struct {
		name              string
		initialTodos      []*Todo
		toggleID          int
		expectedCompleted map[int]bool
	}{
		{
			"toggle incomplete",
			[]*Todo{
				{ID: 1, Text: "First", Completed: false, CreatedAt: now},
				{ID: 2, Text: "Second", Completed: true, CreatedAt: now},
				{ID: 3, Text: "Third", Completed: false, CreatedAt: now},
			},
			1,
			map[int]bool{1: true, 2: true, 3: false},
		},
		{
			"toggle complete",
			[]*Todo{
				{ID: 1, Text: "First", Completed: false, CreatedAt: now},
				{ID: 2, Text: "Second", Completed: true, CreatedAt: now},
				{ID: 3, Text: "Third", Completed: false, CreatedAt: now},
			},
			2,
			map[int]bool{1: false, 2: false, 3: false},
		},
		{
			"nonexistent id no change",
			[]*Todo{
				{ID: 1, Text: "First", Completed: false, CreatedAt: now},
				{ID: 2, Text: "Second", Completed: true, CreatedAt: now},
				{ID: 3, Text: "Third", Completed: false, CreatedAt: now},
			},
			999,
			map[int]bool{1: false, 2: true, 3: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ToggleTodo(tt.initialTodos, tt.toggleID)

			for id, expectedComplete := range tt.expectedCompleted {
				var found bool
				for _, todo := range result {
					if todo.ID == id {
						found = true
						if todo.Completed != expectedComplete {
							t.Errorf("Todo ID %d: expected completed=%v, got %v", id, expectedComplete, todo.Completed)
						}
						break
					}
				}
				if !found {
					t.Errorf("Todo ID %d not found", id)
				}
			}
		})
	}
}

func TestTodoService_ClearCompleted(t *testing.T) {
	service := NewTodoService()
	now := time.Now()

	todos := []*Todo{
		{ID: 1, Text: "First", Completed: false, CreatedAt: now},
		{ID: 2, Text: "Second", Completed: true, CreatedAt: now},
		{ID: 3, Text: "Third", Completed: true, CreatedAt: now},
		{ID: 4, Text: "Fourth", Completed: false, CreatedAt: now},
	}

	result := service.ClearCompleted(todos)

	if len(result) != 2 {
		t.Fatalf("Expected 2 remaining todos, got %d", len(result))
	}

	expectedIDs := []int{1, 4}
	for i, todo := range result {
		if todo.ID != expectedIDs[i] {
			t.Errorf("Position %d: expected ID %d, got %d", i, expectedIDs[i], todo.ID)
		}
		if todo.Completed {
			t.Errorf("Todo ID %d should be incomplete", todo.ID)
		}
	}
}

func TestTodoService_FilterTodos(t *testing.T) {
	service := NewTodoService()
	now := time.Now()

	todos := []*Todo{
		{ID: 1, Text: "First", Completed: false, CreatedAt: now},
		{ID: 2, Text: "Second", Completed: true, CreatedAt: now},
		{ID: 3, Text: "Third", Completed: false, CreatedAt: now},
		{ID: 4, Text: "Fourth", Completed: true, CreatedAt: now},
	}

	tests := []struct {
		name        string
		filter      string
		expectedIDs []int
	}{
		{"all shows all", "all", []int{1, 2, 3, 4}},
		{"active shows incomplete", "active", []int{1, 3}},
		{"completed shows complete", "completed", []int{2, 4}},
		{"invalid filter shows empty", "invalid", []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.FilterTodos(todos, tt.filter)

			if len(result) != len(tt.expectedIDs) {
				t.Fatalf("Expected %d todos, got %d", len(tt.expectedIDs), len(result))
			}

			for i, id := range tt.expectedIDs {
				if result[i].ID != id {
					t.Errorf("Position %d: expected ID %d, got %d", i, id, result[i].ID)
				}
			}
		})
	}
}
