package redux

type Reducer *func(state interface{}, action Action) interface{}
