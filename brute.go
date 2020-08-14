package brute

import (
	"net/http"
	"strings"
)

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

type Handle func(http.ResponseWriter, *http.Request, Params)

type Endpoint struct {
	Method  string
	Pattern string
	Handler Handle

	parts []string
}

type router struct {
	endpoints []Endpoint
}

func New(endpoints ...Endpoint) http.Handler {
	r := &router{}
	for _, endpoint := range endpoints {
		endpoint.Pattern = strings.Trim(endpoint.Pattern, "/")
		endpoint.parts = strings.Split(endpoint.Pattern, "/")
		r.endpoints = append(r.endpoints, endpoint)
	}
	return r
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, endpoint := range r.endpoints {
		params, ok := match(req, endpoint)
		if ok {
			endpoint.Handler(w, req, params)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func match(req *http.Request, endpoint Endpoint) (Params, bool) {
	if req.Method != endpoint.Method {
		return nil, false
	}

	path := strings.Trim(req.URL.Path, "/")
	reqParts := strings.Split(path, "/")

	if len(reqParts) != len(endpoint.parts) {
		return nil, false
	}

	var params Params
	for i, endpointPart := range endpoint.parts {
		if endpointPart != "" && endpointPart[0] == ':' {
			params = append(params, Param{endpointPart[1:], reqParts[i]})
		} else if reqParts[i] != endpointPart {
			return nil, false
		}
	}

	return params, true
}
