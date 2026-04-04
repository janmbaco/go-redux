package application

import (
	"fmt"

	"github.com/janmbaco/go-redux/v2/examples/todo-ddd/domain"
	"github.com/janmbaco/go-redux/v2/handlers"
)

// TodoHandler implements the application logic for todo operations
type TodoHandler struct {
	todoService *domain.TodoService
}

// NewTodoHandler creates a new TodoHandler with injected dependencies
func NewTodoHandler(todoService *domain.TodoService) handlers.ActionHandler[AppState] {
	h := &TodoHandler{
		todoService: todoService,
	}

	builder := handlers.NewActionHandlerBuilder[AppState]()

	builder.On(AddTodoAction, func(state AppState, text string) AppState {
		todoState := state["todos"].(TodoState)
		entities := ToEntities(todoState.Todos)

		// Apply domain logic
		updatedEntities, newNextID := h.todoService.AddTodo(entities, todoState.NextID, text)

		// Convert back to state
		todoState.Todos = ToDTOs(updatedEntities)
		todoState.NextID = newNextID
		state["todos"] = todoState

		fmt.Printf("✅ Added todo: %s (ID: %d)\n", text, newNextID-1)
		return state
	})

	builder.On(ToggleTodoAction, func(state AppState, id int) AppState {
		todoState := state["todos"].(TodoState)
		entities := ToEntities(todoState.Todos)

		// Apply domain logic
		updatedEntities := h.todoService.ToggleTodo(entities, id)

		// Convert back to state
		todoState.Todos = ToDTOs(updatedEntities)
		state["todos"] = todoState

		// Find toggled todo for logging
		for _, todo := range updatedEntities {
			if todo.ID == id {
				status := "active"
				if !todo.IsActive() {
					status = "completed"
				}
				fmt.Printf("🔄 Toggled todo ID %d: %s -> %s\n", id, todo.Text, status)
				break
			}
		}

		return state
	})

	builder.On(RemoveTodoAction, func(state AppState, id int) AppState {
		todoState := state["todos"].(TodoState)
		entities := ToEntities(todoState.Todos)

		// Find todo for logging before removal
		var removedText string
		for _, todo := range entities {
			if todo.ID == id {
				removedText = todo.Text
				break
			}
		}

		// Apply domain logic
		updatedEntities := h.todoService.RemoveTodo(entities, id)

		// Convert back to state
		todoState.Todos = ToDTOs(updatedEntities)
		state["todos"] = todoState

		fmt.Printf("🗑️  Removed todo ID %d: %s\n", id, removedText)
		return state
	})

	builder.On(ClearCompletedAction, func(state AppState) AppState {
		todoState := state["todos"].(TodoState)
		entities := ToEntities(todoState.Todos)

		completedCount := len(entities) - len(h.todoService.ClearCompleted(entities))

		// Apply domain logic
		updatedEntities := h.todoService.ClearCompleted(entities)

		// Convert back to state
		todoState.Todos = ToDTOs(updatedEntities)
		state["todos"] = todoState

		fmt.Printf("🧹 Cleared %d completed todos\n", completedCount)
		return state
	})

	builder.SetSelector("todos")
	builder.SetInitialState(AppState{
		"todos": TodoState{Todos: []TodoDTO{}, NextID: 1},
	})

	return builder.Build()
}
