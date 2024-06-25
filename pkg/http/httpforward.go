package httptool

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewFwd(addr string) (*httputil.ReverseProxy, error) {
	remote, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	forward := httputil.NewSingleHostReverseProxy(remote)
	originalDirector := forward.Director
	forward.Director = func(r *http.Request) {
		originalDirector(r)
		modifyRequest(r)
	}
	forward.ModifyResponse = modifyResponse()
	forward.ErrorHandler = errorHandler()
	return forward, nil
}

// Modify the request to handle the target URL.
func modifyRequest(outReq *http.Request) {
	outReq.Header.Set("X-Proxy", "gopkgs-httptool")
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		log.Printf("Got error while modifying response: %v \n", err)
	}
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		return errors.New("response body is invalid")
	}
}
