package events

import "reflect"

type SelectorSubscribeEvent struct {
	State interface{}
}

func (e *SelectorSubscribeEvent) GetEventArgs() interface{} {
	return e.State
}

func (*SelectorSubscribeEvent) HasEventArgs() bool {
	return true
}

func (*SelectorSubscribeEvent) StopPropagation() bool {
	return false
}

func (*SelectorSubscribeEvent) IsParallelPropagation() bool {
	return true
}

func (e *SelectorSubscribeEvent) GetTypeOfFunc() reflect.Type {
	return reflect.TypeOf(func(state interface{}) {})
}
