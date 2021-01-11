package redux

type BusinessParam interface {
	GetActionsObject() ActionsObject
	GetReducer() Reducer
	GetInitialState() interface{}
	GetSelector() string
}

type businessParam struct {
	actionObject ActionsObject
	reducer      Reducer
	initialState interface{}
	selector     string
}

func (b businessParam) GetActionsObject() ActionsObject {
	return b.actionObject
}

func (b businessParam) GetReducer() Reducer {
	return b.reducer
}

func (b businessParam) GetInitialState() interface{} {
	return b.initialState
}

func (b businessParam) GetSelector() string {
	return b.selector
}

func NewBusinessParam(initialState interface{}, reducer Reducer, actionObject ActionsObject, selector string) BusinessParam {
	return &businessParam{actionObject: actionObject, reducer: reducer, initialState: initialState, selector: selector}
}
