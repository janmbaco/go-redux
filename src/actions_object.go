package redux

import (
	"fmt"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"reflect"
	"strings"
)

type ActionsObject interface {
	GetActions() []Action
	GetActionsNames() []string
	Contains(action Action) bool
	ContainsByName(actionName string) bool
	GetActionByName(name string) Action
}

type actionsObject struct {
	actions      []Action
	actionsNames []string
	actionByName map[string]Action
	nameByAction map[Action]string
}

func NewActionsObject(catcher errors.ErrorCatcher, actions interface{}) ActionsObject {
	errorschecker.CheckNilParameter(map[string]interface{}{"actions": actions})
	result := &actionsObject{
		actions:      getActionsIn(catcher, actions),
		actionsNames: make([]string, 0),
		actionByName: make(map[string]Action),
		nameByAction: make(map[Action]string),
	}
	for _, b := range result.actions {
		result.actionsNames = append(result.actionsNames, b.GetType())
		result.actionByName[b.GetType()] = b
		result.nameByAction[b] = b.GetType()
	}
	return result
}

func (ao *actionsObject) GetActions() []Action {
	return ao.actions
}

func (ao *actionsObject) Contains(action Action) bool {
	_, ok := ao.nameByAction[action]
	return ok
}

func (ao *actionsObject) GetActionsNames() []string {
	return ao.actionsNames
}

func (ao *actionsObject) ContainsByName(actionName string) bool {
	_, ok := ao.actionByName[actionName]
	return ok
}

func (ao *actionsObject) GetActionByName(name string) Action {
	return ao.actionByName[name]
}

func (ao *actionsObject) GetNameByAction(action Action) string {
	return ao.nameByAction[action]
}

func getActionsIn(catcher errors.ErrorCatcher, object interface{}) []Action {
	result := make([]Action, 0)
	rv := reflect.Indirect(reflect.ValueOf(object))
	rt := rv.Type()
	actionType := reflect.TypeOf((*Action)(nil)).Elem()
	panicMessage := make([]string, 1)
	for i := 0; i < rt.NumField(); i++ {
		if rt.Field(i).Type.Implements(actionType) {
			if rv.Field(i).IsNil() {
				catcher.TryCatchError(func() {
					rv.Field(i).Set(reflect.ValueOf(&action{name: rt.Field(i).Name}))
					result = append(result, rv.Field(i).Elem().Interface().(Action))
				}, func(err error) {
					panicMessage = append(panicMessage, fmt.Sprintf("The custom action '%v' of type '%v' in '%v' can't be nil!", rt.Field(i).Name, rt.Field(i).Type.String(), rt.String()))
				})
			} else {
				result = append(result, rv.Field(i).Interface().(Action))
			}
		}
	}
	if len(panicMessage) > 1 {
		panic(strings.Join(panicMessage, "\n"))
	}
	if len(result) == 0 {
		panic("There isn`t any action on the actionsObject object!")
	}
	return result
}
