package redux

import (
	"github.com/janmbaco/go-infrastructure/errorhandler"
)

type SubscribeFunc func(newState interface{})

type Store interface {
	Dispatch(Action)
	Subscribe(ActionsObject, SubscribeFunc)
}

type store struct {
	bo map[ActionsObject]*BusinessObject
}

func NewStore(businessObjects ...*BusinessObject) Store {
	if len(businessObjects) == 0 {
		panic("There must be at least one businessObject parameter")
	}
	newStore := &store{
		bo: make(map[ActionsObject]*BusinessObject),
	}

	for _, bo := range businessObjects {
		if _, ko := newStore.bo[bo.ActionsObject]; ko {
			panic("Cannot add multiple BusinessObject with the same ActionsObject!")
		}
		newStore.bo[bo.ActionsObject] = bo
	}

	return newStore
}

func (s *store) Dispatch(action Action) {
	errorhandler.CheckNilParameter(map[string]interface{}{"action": action})
	for _, bo := range s.bo {
		if bo.ActionsObject.Contains(action) {
			bo.StateManager.SetState(bo.Reducer.Reduce(bo.StateManager.GetState(), action))
			break
		}
	}
}

func (s *store) Subscribe(actionsObject ActionsObject, subscribeFunc SubscribeFunc) {
	errorhandler.CheckNilParameter(map[string]interface{}{"actionsObject": actionsObject, "subscribeFunc": subscribeFunc})
	if _, ok := s.bo[actionsObject]; !ok {
		panic("There is no BusinessObject for that ActionsObject!")
	}

	s.bo[actionsObject].StateManager.Subscribe(func() {
		subscribeFunc(s.bo[actionsObject].StateManager.GetState())
	})
}
