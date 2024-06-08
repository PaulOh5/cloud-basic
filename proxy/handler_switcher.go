package proxy

import (
	"net/http"
	"sync"
)

type HandlerSwitcher struct {
	lock    sync.RWMutex
	handler map[string]http.Handler
}

func NewHandlerSwitcher() *HandlerSwitcher {
	return &HandlerSwitcher{
		handler: make(map[string]http.Handler),
	}
}

func (hs *HandlerSwitcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hs.lock.RLock()
	handler, ok := hs.handler[r.URL.Path]
	hs.lock.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	handler.ServeHTTP(w, r)
}

func (hs *HandlerSwitcher) AddHandler(path string, handler http.Handler) {
	hs.lock.Lock()
	hs.handler[path] = handler
	hs.lock.Unlock()
}

func (hs *HandlerSwitcher) RemoveHandler(path string) {
	hs.lock.Lock()
	delete(hs.handler, path)
	hs.lock.Unlock()
}
