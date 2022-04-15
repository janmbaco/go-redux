package events

import (
	"reflect"
)

type StoreSubscribeEvent struct {
}

func (*StoreSubscribeEvent) IsParallelPropagation() bool {
	return true
}

func (*StoreSubscribeEvent) GetTypeOfFunc() reflect.Type {
	return reflect.TypeOf(func() {})
}

func (*StoreSubscribeEvent) StopPropagation() bool {
	return false
}

func (*StoreSubscribeEvent) HasEventArgs() bool {
	return false
}

func (*StoreSubscribeEvent) GetEventArgs() interface{} {
	return nil
}
