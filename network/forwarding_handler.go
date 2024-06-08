package network

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ForwardingList map[string]string

func NewForwardingHandler(pair ForwardingList) (http.Handler, error) {
	mux := http.NewServeMux()

	for path, targetURL := range pair {
		parsedURL, err := url.Parse(targetURL)
		if err != nil {
			return nil, err
		}
		proxy := httputil.NewSingleHostReverseProxy(parsedURL)
		mux.Handle("/"+path, proxy)
	}

	return mux, nil
}
