package reducers

// Reducer represents a pure function that transforms state without payload.
// This is the fundamental concept from Redux philosophy v1:
// A reducer is a pure function (state) => newState
//
// Reducers must follow these rules:
//   - Pure: same input always produces same output
//   - No side effects: no I/O, no mutations, no random values
//   - Immutable: return new state, don't modify input
//
// Example:
//
//	func ResetReducer(state CounterState) CounterState {
//	    return CounterState{Count: 0}
//	}
type Reducer[S any] func(state S) S

// PayloadReducer represents a pure function that transforms state with a typed payload.
// This is the fundamental concept from Redux philosophy v1:
// A reducer is a pure function (state, payload) => newState
//
// The generic parameter P ensures type-safety at compile time.
//
// Reducers must follow these rules:
//   - Pure: same input always produces same output
//   - No side effects: no I/O, no mutations, no random values
//   - Immutable: return new state, don't modify input
//
// Example:
//
//	func IncrementReducer(state CounterState, amount int) CounterState {
//	    state.Count += amount
//	    return state
//	}
type PayloadReducer[S any, P any] func(state S, payload P) S
