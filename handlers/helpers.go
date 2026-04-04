package handlers

import (
	"fmt"
	"reflect"
)

// Helper functions to invoke module methods using reflection

// InvokeCanHandle calls the CanHandle method on a module
func InvokeCanHandle(module any, action any) (bool, error) {
	method := reflect.ValueOf(module).MethodByName("CanHandle")
	if !method.IsValid() {
		return false, fmt.Errorf("module does not have CanHandle method")
	}

	results := method.Call([]reflect.Value{reflect.ValueOf(action)})
	if len(results) != 1 {
		return false, fmt.Errorf("CanHandle must return 1 value")
	}

	return results[0].Bool(), nil
}

// InvokeHandle calls the Handle method on a module
func InvokeHandle(module any, state any, action any) (any, error) {
	method := reflect.ValueOf(module).MethodByName("Handle")
	if !method.IsValid() {
		return nil, fmt.Errorf("module does not have Handle method")
	}

	results := method.Call([]reflect.Value{
		reflect.ValueOf(state),
		reflect.ValueOf(action),
	})

	if len(results) != 2 {
		return nil, fmt.Errorf("Handle must return 2 values (state, error)")
	}

	var err error
	if !results[1].IsNil() {
		err = results[1].Interface().(error)
	}

	return results[0].Interface(), err
}

// InvokeGetSelector calls the GetSelector method on a module
func InvokeGetSelector(module any) string {
	method := reflect.ValueOf(module).MethodByName("GetSelector")
	if !method.IsValid() {
		return ""
	}

	results := method.Call(nil)
	if len(results) != 1 {
		return ""
	}

	return results[0].String()
}

// InvokeGetInitialState calls the GetInitialState method on a module
func InvokeGetInitialState(module any) any {
	method := reflect.ValueOf(module).MethodByName("GetInitialState")
	if !method.IsValid() {
		return nil
	}

	results := method.Call(nil)
	if len(results) != 1 {
		return nil
	}

	return results[0].Interface()
}
