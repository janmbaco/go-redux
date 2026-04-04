package domain

// TodoService encapsulates business logic for todo operations
type TodoService struct{}

// NewTodoService creates a new TodoService
func NewTodoService() *TodoService {
	return &TodoService{}
}

// AddTodo adds a new todo to the list
func (s *TodoService) AddTodo(todos []*Todo, nextID int, text string) ([]*Todo, int) {
	newTodo := NewTodo(nextID, text)
	return append(todos, newTodo), nextID + 1
}

// ToggleTodo toggles the completion status of a todo by ID
func (s *TodoService) ToggleTodo(todos []*Todo, id int) []*Todo {
	for _, todo := range todos {
		if todo.ID == id {
			todo.Toggle()
			break
		}
	}
	return todos
}

// RemoveTodo removes a todo by ID
func (s *TodoService) RemoveTodo(todos []*Todo, id int) []*Todo {
	result := make([]*Todo, 0, len(todos))
	for _, todo := range todos {
		if todo.ID != id {
			result = append(result, todo)
		}
	}
	return result
}

// ClearCompleted removes all completed todos
func (s *TodoService) ClearCompleted(todos []*Todo) []*Todo {
	result := make([]*Todo, 0)
	for _, todo := range todos {
		if todo.IsActive() {
			result = append(result, todo)
		}
	}
	return result
}

// FilterTodos filters todos based on filter type ("all", "active", "completed")
func (s *TodoService) FilterTodos(todos []*Todo, filter string) []*Todo {
	if filter == "all" {
		return todos
	}

	result := make([]*Todo, 0)
	for _, todo := range todos {
		if filter == "active" && todo.IsActive() {
			result = append(result, todo)
		} else if filter == "completed" && !todo.IsActive() {
			result = append(result, todo)
		}
	}
	return result
}
