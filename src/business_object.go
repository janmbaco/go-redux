package redux

type BusinessObject interface {
	GetActionsObject() ActionsObject
	GetReducer() Reducer
	GetInitialState() interface{}
	GetSelector() string
}

type businessObject struct {
	actionObject ActionsObject
	reducer      Reducer
	initialState interface{}
	selector     string
}

func (b businessObject) GetActionsObject() ActionsObject {
	return b.actionObject
}

func (b businessObject) GetReducer() Reducer {
	return b.reducer
}

func (b businessObject) GetInitialState() interface{} {
	return b.initialState
}

func (b businessObject) GetSelector() string {
	return b.selector
}

func NewBusinessObject(initialState interface{}, reducer Reducer, actionObject ActionsObject, selector string) BusinessObject {
	return &businessObject{actionObject: actionObject, reducer: reducer, initialState: initialState, selector: selector}
}
