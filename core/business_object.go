package redux

type BusinessObject struct {
	ActionsObject ActionsObject
	Reducer       Reducer
	StateEnity    *StateEntity
	StateManager  StateManager
}
