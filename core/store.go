package redux

import (
	"github.com/janmbaco/go-infrastructure/errorhandler"
)

type Store interface {
	Dispatch(Action)
	Subscribe(StateEntity, SubscribeFunc)
	UnSubscribe(StateEntity, SubscribeFunc)
}

type store struct {
	bo map[StateEntity]*BusinessObject
}

func NewStore(businessObjects ...*BusinessObject) Store {
	if len(businessObjects) == 0 {
		panic("There must be at least one businessObject parameter")
	}
	newStore := &store{
		bo: make(map[StateEntity]*BusinessObject),
	}
	actions := make(map[ActionsObject]bool)
	for _, bo := range businessObjects {
		if _, ko := newStore.bo[bo.StateEnity]; ko {
			panic("Cannot add multiple BusinessObject with the same StateEntity!")
		}
		if _, ko := actions[bo.ActionsObject]; ko {
			panic("Cannot add multiple BusinessObject with the same ActionsObject!")
		}
		newStore.bo[bo.StateEnity] = bo
		actions[bo.ActionsObject] = true
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

func (s *store) Subscribe(stateEntity StateEntity, subscribeFunc SubscribeFunc) {
	errorhandler.CheckNilParameter(map[string]interface{}{"stateEntity": stateEntity, "subscribeFunc": subscribeFunc})
	if _, ok := s.bo[stateEntity]; !ok {
		panic("There is no BusinessObject for that StateEntity!")
	}

	s.bo[stateEntity].StateManager.Subscribe(subscribeFunc)

	errorhandler.OnErrorContinue(func() {
		subscribeFunc(s.bo[stateEntity].StateManager.GetState())
	})
}

func (s *store) UnSubscribe(stateEntity StateEntity, subscribeFunc SubscribeFunc) {
	errorhandler.CheckNilParameter(map[string]interface{}{"stateEntity": stateEntity, "subscribeFunc": subscribeFunc})
	if _, ok := s.bo[stateEntity]; !ok {
		panic("There is no BusinessObject for that StateEntity!")
	}

	s.bo[stateEntity].StateManager.UnSubscribe(subscribeFunc)
}
