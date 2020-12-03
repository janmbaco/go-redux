package core

type ReduceFunc func(state interface{}, action Action) interface{}

type Reducer interface {
	Reduce(state interface{}, action Action) interface{}
}

type reducer struct {
	reducer ReduceFunc
}

func (rx *reducer) Reduce(state interface{}, action Action) interface{} {
	return rx.reducer(state, action)
}
