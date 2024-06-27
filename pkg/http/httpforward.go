package httptool

import (
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
	forward.ModifyResponse = func(r *http.Response) error {
		// do nothing now
		return nil
	}
	forward.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Got error while modifying response: %v \n", err)
	}
	return forward, nil
}

// Modify the request to handle the target URL.
func modifyRequest(outReq *http.Request) {
	outReq.Header.Set("X-Proxy", "gopkgs-httptool")
}
