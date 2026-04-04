package domain

import "time"

// Todo represents a todo item entity in the domain
type Todo struct {
	ID        int
	Text      string
	Completed bool
	CreatedAt time.Time
}

// NewTodo creates a new Todo entity
func NewTodo(id int, text string) *Todo {
	return &Todo{
		ID:        id,
		Text:      text,
		Completed: false,
		CreatedAt: time.Now(),
	}
}

// Toggle changes the completion status of the todo
func (t *Todo) Toggle() {
	t.Completed = !t.Completed
}

// IsActive returns true if the todo is not completed
func (t *Todo) IsActive() bool {
	return !t.Completed
}
