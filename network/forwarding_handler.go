package network

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type ForwardingHandler struct {
	lock    sync.RWMutex
	handler map[string]http.Handler
}

func (fh *ForwardingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fh.lock.RLock()
	key := fh.getKey(r.URL.Path)
	handler, ok := fh.handler[key]
	fh.lock.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	handler.ServeHTTP(w, r)
}

func (fh *ForwardingHandler) AddForwarding(key string, targetURL string) error {
	if key == "" || targetURL == "" {
		return fmt.Errorf("key and targetURL should not be empty")
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return err
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy := httputil.NewSingleHostReverseProxy(parsedURL)

		proxy.ModifyResponse = func(resp *http.Response) error {
			if location, ok := resp.Header["Location"]; ok && len(location) > 0 {
				modifedLocation := key + location[0]
				resp.Header.Set("Location", modifedLocation)
			}
			return nil
		}

		r.URL.Host = parsedURL.Host
		r.URL.Scheme = parsedURL.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = parsedURL.Host
		r.URL.Path = r.URL.Path[len(key):]

		proxy.ServeHTTP(w, r)
	})

	fh.lock.Lock()
	fh.handler[key] = handler
	fh.lock.Unlock()

	return nil
}

func (fh *ForwardingHandler) RemoveForwarding(key string) {
	fh.lock.Lock()
	delete(fh.handler, key)
	fh.lock.Unlock()
}

func (fh *ForwardingHandler) getKey(urlPath string) string {
	for key := range fh.handler {
		if len(key) <= len(urlPath) && key == urlPath[:len(key)] {
			return key
		}
	}
	return ""
}

func NewForwardingHandler() *ForwardingHandler {
	return &ForwardingHandler{
		handler: make(map[string]http.Handler),
	}
}
