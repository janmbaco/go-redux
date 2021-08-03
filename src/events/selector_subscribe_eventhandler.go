package events

import (
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

type SelectorSubscribeEventHandler struct {
	subscriptions eventsmanager.Subscriptions
}

func NewSelectorSubscribeEventHandler(subscriptions eventsmanager.Subscriptions) *SelectorSubscribeEventHandler {
	return &SelectorSubscribeEventHandler{subscriptions: subscriptions}
}

func (m *SelectorSubscribeEventHandler) Subscribe(subscription *func(state interface{})) {
	m.subscriptions.Add(&SelectorSubscribeEvent{}, subscription)
}

func (m *SelectorSubscribeEventHandler) UnSubscribe(subscription *func(state interface{})) {
	m.subscriptions.Remove(&SelectorSubscribeEvent{}, subscription)
}
