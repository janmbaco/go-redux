package core

type BusinessObject struct {
	ActionsObject ActionsObject
	Reducer       Reducer
	StateManager  StateManager
}
