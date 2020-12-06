package redux

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/janmbaco/go-infrastructure/errorhandler"
	"github.com/janmbaco/go-infrastructure/events"
	"github.com/janmbaco/go-infrastructure/logs"
)

type businessObjectBuilder struct {
	stateEntity   StateEntity
	stateManager  StateManager
	actionsObject ActionsObject
	blf           map[Action]reflect.Value //business logic funcionality
}

func NewBusinessObjectBuilder(stateEntity StateEntity, actionsObject ActionsObject) *businessObjectBuilder {
	errorhandler.CheckNilParameter(map[string]interface{}{"stateEntity": stateEntity, "actionsObject": actionsObject})

	return &businessObjectBuilder{
		stateEntity:   stateEntity,
		stateManager:  NewStateManager(events.NewPublisher(), stateEntity),
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

	if typeOfState := reflect.TypeOf(builder.stateEntity.GetInitialState()); functionType.NumIn() < 1 || functionType.NumIn() > 2 || functionType.NumOut() != 1 || functionType.In(0) != functionType.Out(0) || functionType.In(0) != typeOfState {
		panic(fmt.Sprintf("The function for action `%v` must to have the contract func(state `%v`, payload *any) `%v`", action.GetName(), typeOfState.Name(), typeOfState.Name()))
	}

	if functionType.NumIn() == 2 {
		action.SetType(functionType.In(1))
	}

	builder.blf[action] = functionValue
	return builder
}

func (builder *businessObjectBuilder) SetActionsLogicByObject(object interface{}) *businessObjectBuilder {
	errorhandler.CheckNilParameter(map[string]interface{}{"object": object})
	typeOfState := reflect.TypeOf(builder.stateEntity.GetInitialState())
	if reflect.TypeOf(object) == typeOfState {
		panic("You cannot create the logic of the actions with the same type as the state!")
	}

	rv := reflect.ValueOf(object)
	rt := reflect.TypeOf(object)
	if rt.Kind() != reflect.Ptr && rt.Kind() != reflect.Struct {
		panic("The object must be a struct")
	}

	for i := 0; i < rt.NumMethod(); i++ {
		m := rt.Method(i)
		mt := m.Type
		if mt.NumIn() >= 2 && mt.NumIn() <= 3 && mt.NumOut() == 1 && mt.In(1) == mt.Out(0) && mt.In(1) == typeOfState {
			if builder.actionsObject.ContainsByName(m.Name) {
				action := builder.actionsObject.GetActionByName(m.Name)
				if mt.NumIn() == 3 {
					action.SetType(mt.In(2))
				}
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
		var result interface{}

		if typeOfState := reflect.TypeOf(state); (function.Type().NumIn() < 3 && function.Type().In(0) != typeOfState) || function.Type().NumIn() == 1 {
			result = function.Call([]reflect.Value{
				reflect.ValueOf(state),
			})[0].Interface()
		} else {
			result = function.Call([]reflect.Value{
				reflect.ValueOf(state),
				action.GetPayload(),
			})[0].Interface()
		}

		return result

	}

	reducer := &reducer{
		reducer: reducerFunc,
	}

	return &BusinessObject{
		ActionsObject: builder.actionsObject,
		Reducer:       reducer,
		StateEnity:    builder.stateEntity,
		StateManager:  builder.stateManager,
	}
}
