package main

import (
	"fmt"

	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	logsioc "github.com/janmbaco/go-infrastructure/v2/logs/ioc"
	"github.com/janmbaco/go-redux/v2/examples/game-server/application"
	appioc "github.com/janmbaco/go-redux/v2/examples/game-server/application/ioc"
	appresolver "github.com/janmbaco/go-redux/v2/examples/game-server/application/ioc/resolver"
	domainioc "github.com/janmbaco/go-redux/v2/examples/game-server/domain/ioc"
	"github.com/janmbaco/go-redux/v2/examples/game-server/infrastructure/subscribers"
	"github.com/janmbaco/go-redux/v2/ioc"
	"github.com/janmbaco/go-redux/v2/ioc/resolver"
)

func main() {
	fmt.Println("=== Redux Game Server Example ===")
	fmt.Println()
	fmt.Println("Demonstrating Clean Architecture with multiple services and subscribers:")
	fmt.Println()
	fmt.Println("- Domain: Player, MovementService, CombatService")
	fmt.Println("- Application: PlayerHandler, CombatHandler")
	fmt.Println("- Infrastructure: BroadcastSubscriber, PersistenceSubscriber, AnalyticsSubscriber")
	fmt.Println()

	// Initial state
	initialState := application.GameState{
		Players: make(map[string]application.PlayerDTO),
	}

	// Build DI container with all modules
	container := dependencyinjection.NewBuilder().
		AddModule(logsioc.NewLogsModule()).
		AddModule(ioc.NewReduxModule[application.GameState]()).
		AddModule(domainioc.NewGameDomainModule()).
		AddModule(appioc.NewGameApplicationModule()).
		MustBuild()

	// Resolve store from container
	store := resolver.GetStore[application.GameState](container.Resolver(), initialState)

	// Resolve and add handlers
	playerHandler := appresolver.GetPlayerHandler(container.Resolver())
	store.AddModule(playerHandler)

	combatHandler := appresolver.GetCombatHandler(container.Resolver())
	store.AddModule(combatHandler)

	// Setup subscribers (each will be called on every state change)
	broadcastSub := subscribers.NewBroadcastSubscriber()
	broadcastCallback := broadcastSub.OnStateChange
	store.Subscribe(&broadcastCallback)

	persistenceSub := subscribers.NewPersistenceSubscriber()
	persistenceCallback := persistenceSub.OnStateChange
	store.Subscribe(&persistenceCallback)

	analyticsSub := subscribers.NewAnalyticsSubscriber()
	analyticsCallback := analyticsSub.OnStateChange
	store.Subscribe(&analyticsCallback)

	fmt.Println()

	// Simulate game events
	fmt.Println("=== Players Joining ===")
	store.Dispatch(application.JoinGameAction.With(application.JoinGamePayload{
		PlayerID: "player1",
		Name:     "Alice",
	}))

	store.Dispatch(application.JoinGameAction.With(application.JoinGamePayload{
		PlayerID: "player2",
		Name:     "Bob",
	}))

	store.Dispatch(application.JoinGameAction.With(application.JoinGamePayload{
		PlayerID: "player3",
		Name:     "Charlie",
	}))

	fmt.Println("\n=== Player Movement ===")
	store.Dispatch(application.MovePlayerAction.With(application.MovePlayerPayload{
		PlayerID: "player1",
		X:        5.0,
		Y:        3.0,
	}))

	store.Dispatch(application.MovePlayerAction.With(application.MovePlayerPayload{
		PlayerID: "player2",
		X:        -2.0,
		Y:        4.0,
	}))

	fmt.Println("\n=== Combat ===")
	store.Dispatch(application.AttackPlayerAction.With(application.AttackPayload{
		AttackerID: "player1",
		VictimID:   "player2",
	}))

	store.Dispatch(application.AttackPlayerAction.With(application.AttackPayload{
		AttackerID: "player1",
		VictimID:   "player2",
	}))

	fmt.Println("\n=== Healing ===")
	store.Dispatch(application.HealPlayerAction.With(application.HealPayload{
		PlayerID: "player2",
		Amount:   30,
	}))

	fmt.Println("\n=== More Combat ===")
	store.Dispatch(application.AttackPlayerAction.With(application.AttackPayload{
		AttackerID: "player3",
		VictimID:   "player1",
	}))

	fmt.Println("\n=== Player Leaving ===")
	store.Dispatch(application.LeaveGameAction.With("player3"))

	fmt.Printf("\n=== Final State ===\n")
	finalState := store.GetState()
	fmt.Printf("Players in game: %d\n", len(finalState.Players))
	for id, player := range finalState.Players {
		fmt.Printf("  %s: %s [HP:%d Score:%d]\n", id, player.Name, player.Health, player.Score)
	}

	// Clean up
	store.Close()
	fmt.Println("\n✅ Game Server example completed successfully!")
}
