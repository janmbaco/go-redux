package redux

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"github.com/janmbaco/go-infrastructure/logs"
)

type BusinesParamBuilder interface {
	SetInitialState(interface{}) BusinesParamBuilder
	SetActions(interface{}) BusinesParamBuilder
	SetSelector(selector string) BusinesParamBuilder
	On(action Action, function interface{}) BusinesParamBuilder
	SetActionsLogicByObject(object interface{}) BusinesParamBuilder
	GetBusinessParam() BusinessParam
}

type businessParamBuilder struct {
	logger        logs.Logger
	initialState  interface{}
	selector      string
	businessParamFactory BusinessParamFactory
	actionsObjectFactory ActionsObjectFactory
	actionsObject ActionsObject
	blf           map[Action]reflect.Value // business logic funcionality
}
type redueActions struct {
	blf map[Action]reflect.Value
}

func NewBusinessParamBuilder(logger logs.Logger, aactionsObjectFactory ActionsObjectFactory, businessParamFactory BusinessParamFactory) BusinesParamBuilder {
	errorschecker.CheckNilParameter(map[string]interface{}{"logger": logger, "aactionsObjectFactory": aactionsObjectFactory, "businessParamFactory":businessParamFactory})
	return &businessParamBuilder{blf: make(map[Action]reflect.Value), logger: logger, actionsObjectFactory: aactionsObjectFactory, businessParamFactory: businessParamFactory}
}

func (builder *businessParamBuilder) SetInitialState(initialState interface{}) BusinesParamBuilder {
	builder.initialState = initialState
	return builder
}

func (builder *businessParamBuilder) SetActions(actions interface{}) BusinesParamBuilder {
	errorschecker.CheckNilParameter(map[string]interface{}{"actions": actions})
	builder.actionsObject = builder.actionsObjectFactory.Create(actions)
	return builder
}

func (builder *businessParamBuilder) SetSelector(selector string) BusinesParamBuilder {
	if selector == "" {
		panic("The selector can not be string empty!")
	}

	builder.selector = selector
	return builder
}

func (builder *businessParamBuilder) On(action Action, function interface{}) BusinesParamBuilder {
	errorschecker.CheckNilParameter(map[string]interface{}{"action": action, "function": function})

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

	if typeOfState := reflect.TypeOf(builder.initialState); functionType.NumIn() < 1 || functionType.NumIn() > 2 || functionType.NumOut() != 1 || functionType.In(0) != functionType.Out(0) || functionType.In(0) != typeOfState {
		panic(fmt.Errorf("the function for action `%v` must to have the contract func(state `%v`, payload *any) `%v`", action.GetType(), typeOfState.Name(), typeOfState.Name()))
	}

	if functionType.NumIn() == 2 {
		action.SetPayloadType(functionType.In(1))
	}

	builder.blf[action] = functionValue
	return builder
}

func (builder *businessParamBuilder) SetActionsLogicByObject(object interface{}) BusinesParamBuilder {
	errorschecker.CheckNilParameter(map[string]interface{}{"object": object})
	typeOfState := reflect.TypeOf(builder.initialState)
	if reflect.TypeOf(object) == typeOfState {
		panic("You cannot create the logic of the actionsObject with the same type as the state!")
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
					action.SetPayloadType(mt.In(2))
				}
				builder.blf[action] = rv.Method(i)
			} else {
				builder.logger.Warning(fmt.Sprintf("The func`%v` in the object `%v` has not a action asociated in the ActionsObject! ActionObject:`%v`", m.Name, rt.String(), builder.actionsObject.GetActionsNames()))
			}
		}
	}
	return builder
}

func (builder *businessParamBuilder) GetBusinessParam() BusinessParam {

	if builder.selector == "" {
		panic("The selector can not be string empty!")
	}

	if builder.actionsObject == nil {
		panic("There isn´t any ActionsObject to load to the BusinessParam!")
	}

	panicMessage := strings.Builder{}
	for _, action := range builder.actionsObject.GetActions() {
		if _, ok := builder.blf[action]; !ok {
			panicMessage.WriteString(fmt.Sprintf("The logic for the actionsObject '%v' is not defined!\n", action.GetType()))
		}
	}
	if panicMessage.Len() > 0 {
		panic(panicMessage.String())
	}
	reducerActions := &redueActions{blf: make(map[Action]reflect.Value)}
	for key, value := range builder.blf {
		reducerActions.blf[key] = value
	}
	reducer := reducerActions.Reducer
	if builder.selector == "" {
		builder.selector = strconv.Itoa(int(reflect.ValueOf(builder.initialState).Pointer()))
	}

	businessParam := builder.businessParamFactory.Create(
		BusinessParamFactoryParamter{
			builder.initialState,
			&reducer,
			builder.actionsObject,
			builder.selector,
	})

	builder.initialState = nil
	builder.actionsObject = nil
	builder.selector = ""
	for k := range builder.blf {
		delete(builder.blf, k)
	}

	return businessParam
}

func (ra *redueActions) Reducer(state interface{}, action Action) interface{} {
	function, exists := ra.blf[action]
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
