package events

import "github.com/gxjakkap/reception/store"

type EventContext struct {
	gs *store.GuildsStore
	ps *store.PendingStore
}

func NewEventContext(gs *store.GuildsStore, ps *store.PendingStore) *EventContext {
	return &EventContext{
		gs: gs,
		ps: ps,
	}
}
