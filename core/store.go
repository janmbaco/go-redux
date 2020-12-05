package redux

import (
	"github.com/janmbaco/go-infrastructure/errorhandler"
)

type SubscribeFunc func(newState interface{})

type Store interface {
	Dispatch(Action)
	Subscribe(*StateEntity, SubscribeFunc)
}

type store struct {
	bo map[*StateEntity]*BusinessObject
}

func NewStore(businessObjects ...*BusinessObject) Store {
	if len(businessObjects) == 0 {
		panic("There must be at least one businessObject parameter")
	}
	newStore := &store{
		bo: make(map[*StateEntity]*BusinessObject),
	}

	for _, bo := range businessObjects {
		if _, ko := newStore.bo[bo.StateEnity]; ko {
			panic("Cannot add multiple BusinessObject with the same StateEntity!")
		}
		newStore.bo[bo.StateEnity] = bo
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

func (s *store) Subscribe(stateEntity *StateEntity, subscribeFunc SubscribeFunc) {
	errorhandler.CheckNilParameter(map[string]interface{}{"stateEntity": stateEntity, "subscribeFunc": subscribeFunc})
	if _, ok := s.bo[stateEntity]; !ok {
		panic("There is no BusinessObject for that StateEntity!")
	}

	s.bo[stateEntity].StateManager.Subscribe(func() {
		subscribeFunc(s.bo[stateEntity].StateManager.GetState())
	})
}
