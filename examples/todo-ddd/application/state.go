package application

import "time"

// TodoDTO represents a todo item in the application state
type TodoDTO struct {
	ID        int
	Text      string
	Completed bool
	CreatedAt time.Time
}

// TodoState represents the Redux state for todos
type TodoState struct {
	Todos  []TodoDTO
	NextID int
}

// FilterState represents the Redux state for filtering
type FilterState struct {
	Filter string // "all", "active", "completed"
}

// AppState is the combined application state
type AppState map[string]any
