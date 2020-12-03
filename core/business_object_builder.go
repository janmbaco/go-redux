package core

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/janmbaco/go-infrastructure/errorhandler"
	"github.com/janmbaco/go-infrastructure/events"
	"github.com/janmbaco/go-infrastructure/logs"
)

type businessObjectBuilder struct {
	initialState  interface{}
	stateManager  StateManager
	actionsObject ActionsObject
	blf           map[Action]reflect.Value //business logic funcionality
}

func NewBusinessObjectBuilder(initialState interface{}, actionsObject ActionsObject) *businessObjectBuilder {
	errorhandler.CheckNilParameter(map[string]interface{}{"initialState": initialState, "actionsObject": actionsObject})

	return &businessObjectBuilder{
		initialState:  initialState,
		stateManager:  NewStateManager(events.NewEventPublisher(), initialState),
		actionsObject: actionsObject,
		blf:           make(map[Action]reflect.Value)}
}

func (builder *businessObjectBuilder) SetStateManager(stateManager StateManager) *businessObjectBuilder {
	errorhandler.CheckNilParameter(map[string]interface{}{"stateManager": stateManager})
	builder.stateManager = stateManager
	return builder
}

func (builder *businessObjectBuilder) On(action Action, function interface{}) *businessObjectBuilder {
	errorhandler.CheckNilParameter(map[string]interface{}{"action": action, "function": function})

	if !builder.actionsObject.Contains(action) {
		panic("This action doesn`t belong to this BusinesObject!")
	}

	if _, exists := builder.blf[action]; exists {
		panic("action already reduced!")
	}

	functionValue := reflect.ValueOf(function)
	functionType := reflect.TypeOf(function)
	if functionType.Kind() != reflect.Func {
		panic("The function must be a Func!")
	}

	if typeOfState := reflect.TypeOf(builder.initialState); functionType.NumIn() != 2 || functionType.NumOut() != 1 || functionType.In(0) != functionType.Out(0) || functionType.In(0) != typeOfState {
		panic(fmt.Sprintf("The function for action `%v` must to have the contract func(state `%v`, payload *any) `%v`", action.GetName(), typeOfState.Name(), typeOfState.Name()))
	}

	action.SetType(functionType.In(1))

	builder.blf[action] = functionValue
	return builder
}

func (builder *businessObjectBuilder) SetActionsLogicByObject(object interface{}) *businessObjectBuilder {
	errorhandler.CheckNilParameter(map[string]interface{}{"object": object})

	rv := reflect.ValueOf(object)
	rt := reflect.TypeOf(object)
	if rt.Kind() != reflect.Ptr && rt.Kind() != reflect.Struct {
		panic("The object must be a struct")
	}

	for i := 0; i < rt.NumMethod(); i++ {
		m := rt.Method(i)
		mt := m.Type
		if typeOfState := reflect.TypeOf(builder.initialState); mt.NumIn() == 3 && mt.NumOut() == 1 && mt.In(1) == mt.Out(0) && mt.In(1) == typeOfState {
			if builder.actionsObject.ContainsByName(m.Name) {
				action := builder.actionsObject.GetActionByName(m.Name)
				action.SetType(mt.In(2))
				builder.blf[action] = rv.Method(i)
			} else {
				logs.Log.Warning(fmt.Sprintf("The func`%v` in the object `%v` has not a action asociated in the ActionsObject! ActionObject:`%v`", m.Name, rt.String(), builder.actionsObject.GetActionsNames()))
			}

		}
	}
	return builder
}

func (builder *businessObjectBuilder) GetBusinessObject() *BusinessObject {

	if builder.actionsObject == nil {
		panic("There isn´t any ActionsObject to load to the BusinessObject!")
	}

	panicMessage := strings.Builder{}
	for _, action := range builder.actionsObject.GetActions() {
		if _, ok := builder.blf[action]; !ok {
			panicMessage.WriteString(fmt.Sprintf("The logic for the actionsObject '%v' is not defined!\n", action.GetName()))
		}
	}
	if panicMessage.Len() > 0 {
		panic(panicMessage.String())
	}

	reducerFunc := func(state interface{}, action Action) interface{} {

		function, exists := builder.blf[action]
		if !exists {
			panic("The action is not located in the reducer function!")
		}

		return function.Call([]reflect.Value{
			reflect.ValueOf(state),
			action.GetPayload(),
		})[0].Interface()
	}

	reducer := &reducer{
		reducer: reducerFunc,
	}

	return &BusinessObject{
		ActionsObject: builder.actionsObject,
		Reducer:       reducer,
		StateManager:  builder.stateManager,
	}
}
