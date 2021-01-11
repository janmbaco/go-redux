package redux

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/janmbaco/go-infrastructure/errorhandler"
)

type ActionsObject interface {
	GetActions() []Action
	GetActionsNames() []string
	Contains(action Action) bool
	ContainsByName(actionName string) bool
	GetActionByName(name string) Action
}

type actionObject struct {
	actions      []Action
	actionsNames []string
	actionByName map[string]Action
	nameByAction map[Action]string
}

func NewActionsObject(object interface{}) ActionsObject {
	errorhandler.CheckNilParameter(map[string]interface{}{"object": object})
	if object == nil {
		panic("The object parameter can`t be nil!")
	}
	result := &actionObject{
		actions:      getActionsIn(object),
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

func (actionObject *actionObject) GetActions() []Action {
	return actionObject.actions
}

func (actionObject *actionObject) Contains(action Action) bool {
	_, ok := actionObject.nameByAction[action]
	return ok
}

func (actionObject *actionObject) GetActionsNames() []string {
	return actionObject.actionsNames
}

func (actionObject *actionObject) ContainsByName(actionName string) bool {
	_, ok := actionObject.actionByName[actionName]
	return ok
}

func (actionObject *actionObject) GetActionByName(name string) Action {
	return actionObject.actionByName[name]
}

func (actionObject *actionObject) GetNameByAction(action Action) string {
	return actionObject.nameByAction[action]
}

func getActionsIn(object interface{}) []Action {
	result := make([]Action, 0)
	rv := reflect.Indirect(reflect.ValueOf(object))
	rt := rv.Type()
	actionType := reflect.TypeOf((*Action)(nil)).Elem()
	panicMessage := make([]string, 1)
	for i := 0; i < rt.NumField(); i++ {
		if rt.Field(i).Type.Implements(actionType) {
			if rv.Field(i).IsNil() {
				errorhandler.TryCatchError(func() {
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
