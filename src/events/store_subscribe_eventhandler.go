package events

import (
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

type StoreSubscribeEventHandler struct {
	subscriptions eventsmanager.Subscriptions
}

func NewStoreSubscribeEventHandler(subscriptions eventsmanager.Subscriptions) *StoreSubscribeEventHandler {
	return &StoreSubscribeEventHandler{subscriptions: subscriptions}
}

func (m *StoreSubscribeEventHandler) Unsubscribe(subscription *func()) {
	m.subscriptions.Remove(&StoreSubscribeEvent{}, subscription)
}

func (m *StoreSubscribeEventHandler) Subscribe(subscription *func()) {
	m.subscriptions.Add(&StoreSubscribeEvent{}, subscription)
}
