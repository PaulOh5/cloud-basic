package network

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxyHandler(targetURL string) (http.Handler, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	proxyHandler := httputil.NewSingleHostReverseProxy(parsedURL)
	return proxyHandler, nil
}
