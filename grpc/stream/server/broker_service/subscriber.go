package broker_service

import (
	"grpc/stream/prototypes"
	"sync"
)

var (
	globalSubscriberMappingMtx sync.RWMutex
	globalSubscriberMapping    = make(map[string]chan<- *prototypes.MessageEntry)
)

func pushMessageToSubscriber(msg *prototypes.MessageEntry) {
	globalSubscriberMappingMtx.RLock()
	defer globalSubscriberMappingMtx.RUnlock()
	for _, ch := range globalSubscriberMapping {
		ch <- msg
	}
}

func addSubscriber(name string, ch chan<- *prototypes.MessageEntry) {
	globalSubscriberMappingMtx.Lock()
	defer globalSubscriberMappingMtx.Unlock()
	_, ok := globalSubscriberMapping[name]
	if !ok {
		globalSubscriberMapping[name] = ch
	}
}

func delSubscriber(name string) {
	globalSubscriberMappingMtx.Lock()
	defer globalSubscriberMappingMtx.Unlock()
	_, ok := globalSubscriberMapping[name]
	if ok {
		delete(globalSubscriberMapping, name)
	}
}
