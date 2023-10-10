package routes

import (
	"sync"
)

var sseClients = make(map[chan string]bool)
var sseClientsMutex = sync.Mutex{}

func addSSEClient(client chan string) {
	sseClientsMutex.Lock()
	defer sseClientsMutex.Unlock()
	sseClients[client] = true
}

func removeSSEClient(client chan string) {
	sseClientsMutex.Lock()
	defer sseClientsMutex.Unlock()
	delete(sseClients, client)
}

func broadcastEvent(event string) {
	sseClientsMutex.Lock()
	defer sseClientsMutex.Unlock()
	for client := range sseClients {
		client <- event
	}
}
