package redux

import (
	"fmt"
	"reflect"
)

type Action interface {
	With(interface{}) Action
	GetPayload() reflect.Value
	SetType(reflect.Type)
	GetName() string
}

type action struct {
	payload   reflect.Value
	typ       reflect.Type
	payloaded bool
	name      string
}

func (action *action) With(payload interface{}) Action {

	if action.typ == nil {
		panic("no payload can be assigned to this action!")
	}

	if reflect.TypeOf(payload) != action.typ {
		panic(fmt.Sprintf("The type of payload must be '%v'", action.typ.String()))
	}

	action.payload = reflect.ValueOf(payload)
	action.payloaded = true
	return action
}

func (action *action) GetPayload() reflect.Value {
	var result reflect.Value
	if action.payloaded {
		result = action.payload
		action.payloaded = false
	} else {
		result = reflect.Zero(action.typ)
	}
	return result
}

func (action *action) SetType(typ reflect.Type) {
	action.typ = typ
}

func (action *action) GetName() string {
	return action.name
}
