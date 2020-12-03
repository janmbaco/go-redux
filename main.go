package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/janmbaco/go-infrastructure/logs"
	"github.com/janmbaco/go-redux/core"
)

type Actions struct {
	actionsObject redux.ActionsObject
	Sumar         redux.Action
	Restar        redux.Action
}

func Sumar(state int, payload *int) int {
	var result int
	if payload == nil {
		result = state + 1
	} else {
		result = state + *payload
	}
	return result
}

type RestarObject struct{}

func (ro *RestarObject) Restar(state int, payload *int) int {
	return state - *payload
}

func main() {
	var myActions = &Actions{}

	myActions.actionsObject = redux.NewActionsObject(myActions)

	businessObjectBuilder := redux.NewBusinessObjectBuilder(0, myActions.actionsObject)

	businessObjectBuilder.On(myActions.Sumar, Sumar)

	businessObjectBuilder.SetActionsLogicByObject(&RestarObject{})

	myStore := redux.NewStore(businessObjectBuilder.GetBusinessObject())

	wg := sync.WaitGroup{}
	pass := 1
	myStore.Subscribe(myActions.actionsObject, func(newState interface{}) {
		logs.Log.Info(strconv.Itoa(newState.(int)))
		var expected int
		switch pass {
		case 1:
			expected = 1
		case 2:
			expected = 6
		case 3:
			expected = 7
		case 4:
			expected = 0
		case 5:
			expected = 7
		}
		if newState.(int) != expected {
			logs.Log.Error(fmt.Sprintf("expected: `%v`, found: `%v`", expected, newState.(int)))
		}
		pass++
		wg.Done()
	})
	wg.Add(1)
	myStore.Dispatch(myActions.Sumar)
	wg.Add(1)
	a := 5
	myStore.Dispatch(myActions.Sumar.With(&a))
	wg.Add(1)
	myStore.Dispatch(myActions.Sumar)
	wg.Add(1)

	a = -7
	myStore.Dispatch(myActions.Sumar.With(&a))
	wg.Add(1)
	myStore.Dispatch(myActions.Restar.With(&a))
	wg.Wait()
}
