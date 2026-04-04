package main

import (
	"fmt"

	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	logsioc "github.com/janmbaco/go-infrastructure/v2/logs/ioc"
	"github.com/janmbaco/go-redux/v2/examples/todo-ddd/application"
	appioc "github.com/janmbaco/go-redux/v2/examples/todo-ddd/application/ioc"
	appresolver "github.com/janmbaco/go-redux/v2/examples/todo-ddd/application/ioc/resolver"
	domainioc "github.com/janmbaco/go-redux/v2/examples/todo-ddd/domain/ioc"
	"github.com/janmbaco/go-redux/v2/ioc"
	"github.com/janmbaco/go-redux/v2/ioc/resolver"
)

func printTodos(state application.AppState) {
	todoState := state["todos"].(application.TodoState)
	filterState := state["filter"].(application.FilterState)

	fmt.Printf("\n📋 Todo List (filter: %s):\n", filterState.Filter)
	fmt.Println("─────────────────────────────────────")

	activeCount := 0
	completedCount := 0

	for _, todo := range todoState.Todos {
		// Apply filter
		if filterState.Filter == "active" && todo.Completed {
			continue
		}
		if filterState.Filter == "completed" && !todo.Completed {
			continue
		}

		status := "⬜"
		if todo.Completed {
			status = "✅"
			completedCount++
		} else {
			activeCount++
		}

		fmt.Printf("%s [%d] %s\n", status, todo.ID, todo.Text)
	}

	fmt.Printf("─────────────────────────────────────\n")
	fmt.Printf("Total: %d | Active: %d | Completed: %d\n\n",
		len(todoState.Todos), activeCount, completedCount)
}

func main() {
	fmt.Println("=== Redux Todo DDD Example ===")
	fmt.Println()
	fmt.Println("Demonstrating Clean Architecture with multiple handlers:")
	fmt.Println()
	fmt.Println("- Domain: Todo entity + TodoService with business logic")
	fmt.Println("- Application: TodoHandler and FilterHandler using domain services")
	fmt.Println("- Infrastructure: IoC setup with multiple modules")
	fmt.Println()

	// Initial state
	initialState := application.AppState{
		"todos":  application.TodoState{Todos: []application.TodoDTO{}, NextID: 1},
		"filter": application.FilterState{Filter: "all"},
	}

	// Build DI container with all modules
	container := dependencyinjection.NewBuilder().
		AddModule(logsioc.NewLogsModule()).
		AddModule(ioc.NewReduxModule[application.AppState]()).
		AddModule(domainioc.NewTodoDomainModule()).
		AddModule(appioc.NewTodoApplicationModule()).
		MustBuild()

	// Resolve store from container
	store := resolver.GetStore[application.AppState](container.Resolver(), initialState)

	// Resolve handlers using application resolver
	todoHandler := appresolver.GetTodoHandler(container.Resolver())
	store.AddModule(todoHandler)

	filterHandler := appresolver.GetFilterHandler(container.Resolver())
	store.AddModule(filterHandler)

	// Subscribe to state changes
	printCallback := printTodos
	store.Subscribe(&printCallback)

	// Dispatch actions
	fmt.Println("Adding todos...")
	store.Dispatch(application.AddTodoAction.With("Learn Go Redux v2"))
	store.Dispatch(application.AddTodoAction.With("Implement DDD pattern"))
	store.Dispatch(application.AddTodoAction.With("Create examples"))

	fmt.Println("\nToggling some todos...")
	store.Dispatch(application.ToggleTodoAction.With(1))
	store.Dispatch(application.ToggleTodoAction.With(2))

	fmt.Println("\nAdding more todos...")
	store.Dispatch(application.AddTodoAction.With("Write documentation"))

	fmt.Println("\nFiltering by active...")
	store.Dispatch(application.SetFilterAction.With("active"))

	fmt.Println("\nFiltering by completed...")
	store.Dispatch(application.SetFilterAction.With("completed"))

	fmt.Println("\nShowing all todos...")
	store.Dispatch(application.SetFilterAction.With("all"))

	fmt.Println("\nRemoving a todo...")
	store.Dispatch(application.RemoveTodoAction.With(3))

	fmt.Println("\nClearing completed todos...")
	store.Dispatch(application.ClearCompletedAction.With(struct{}{}))

	fmt.Printf("\n=== Final State ===\n")
	finalState := store.GetState()
	fmt.Printf("Todos: %+v\n", finalState["todos"])
	fmt.Printf("Filter: %+v\n", finalState["filter"])

	store.Close()
	fmt.Println("\n✅ Todo DDD example completed successfully!")
}
