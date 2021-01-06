package main

import (
	"github.com/janmbaco/go-infrastructure/logs"
	"github.com/janmbaco/go-redux/core"
	"strconv"
)

type IntState struct {
	initialState *int
}

func (s *IntState) GetInitialState() interface{} {
	return s.initialState
}

type Actions struct {
	actionsObject redux.ActionsObject
	Sumar         redux.Action
	Restar        redux.Action
}

func Sumar(state *int, payload *int) *int {

	if payload == nil {
		*state = *state + 1
	} else {
		*state = *state + *payload
	}
	return state
}

type RestarObject struct{}

func (ro *RestarObject) Restar(state *int, payload *int) *int {
	if payload == nil {
		*state = *state - 1
	} else {
		*state = *state - *payload
	}
	return state
}

type subsciption struct {
	fn1 redux.SubscribeFunc
}

func (s *subsciption) Initialize(name string) {
	s.fn1 = func(newState interface{}) {
		logs.Log.Info(name + " Current state: " + strconv.Itoa(*newState.(*int)))
	}
}

func main() {

	var myActions = &Actions{}

	myActions.actionsObject = redux.NewActionsObject(myActions)
	initialState := 0
	var myState = &IntState{initialState: &initialState}

	businessObjectBuilder := redux.NewBusinessObjectBuilder(myState, myActions.actionsObject)
	businessObjectBuilder.On(myActions.Sumar, Sumar)
	businessObjectBuilder.SetActionsLogicByObject(&RestarObject{})

	myStore := redux.NewStore(businessObjectBuilder.GetBusinessObject())

	sub1 := &subsciption{}
	sub1.Initialize("sub1")

	logs.Log.Info("Get current state on subscribe")
	myStore.Subscribe(myState, &sub1.fn1)

	sub2 := &subsciption{}
	sub2.Initialize("sub2")

	logs.Log.Info("Get current state on subscribe")
	myStore.Subscribe(myState, &sub2.fn1)

	logs.Log.Info("Dispatch subscribed")

	logs.Log.Info("Sum 1")
	myStore.Dispatch(myActions.Sumar)

	logs.Log.Info("Sum 5")
	a := 5
	myStore.Dispatch(myActions.Sumar.With(&a))

	logs.Log.Info("Sum 1")
	myStore.Dispatch(myActions.Sumar)

	logs.Log.Info("Sum -7")
	a = -7
	myStore.Dispatch(myActions.Sumar.With(&a))

	logs.Log.Info("Substract -7")
	myStore.Dispatch(myActions.Restar.With(&a))

	logs.Log.Info("Uunsubscribe")
	myStore.UnSubscribe(myState, &sub1.fn1)

	logs.Log.Info("Dispatch unsubscribed")

	logs.Log.Info("Sum -7")
	myStore.Dispatch(myActions.Sumar.With(&a))

	logs.Log.Info("Substract 1")
	myStore.Dispatch(myActions.Restar)

	logs.Log.Info("Get current state on resubscribe")
	myStore.Subscribe(myState, &sub1.fn1)

	logs.Log.Info("finish")

}
