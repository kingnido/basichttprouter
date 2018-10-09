package basichttprouter

import (
	"fmt"
	"strings"

	"net/http"
)

type Router struct {
	routes  map[string]*Router
	handler http.Handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := strings.FieldsFunc(req.URL.Path, func(c rune) bool { return c == '/' })
	node := r
	found := false

	for _, route := range path {
		node, found = node.routes[route]
		if !found {
			http.NotFound(w, req)
			return
		}
	}

	node.handler.ServeHTTP(w, req)
}

func NewRouter() *Router {
	r := &Router{
		routes: map[string]*Router{},
	}

	// show possible routes if no handler defined
	r.handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		routes := make([]string, 0, len(r.routes))
		for k, _ := range r.routes {
			routes = append(routes, k)
		}

		fmt.Fprint(w, routes)
	})

	return r
}

func (r *Router) Handle(p string, h http.Handler) error {
	path := strings.FieldsFunc(p, func(c rune) bool { return c == '/' })

	r.handle(path, h)

	return nil
}

func (r *Router) handle(p []string, h http.Handler) error {
	if len(p) == 0 {
		r.handler = h
		return nil
	}

	subroute := p[0]

	subrouter, found := r.routes[subroute]
	if !found {
		subrouter = NewRouter()
		r.routes[subroute] = subrouter
	}

	subrouter.handle(p[1:], h)

	return nil
}

func (r *Router) String() string {
	return "/\n" + r.string(1)
}

func (r *Router) string(indent int) string {
	result := ""

	for key, value := range r.routes {
		result += strings.Repeat(" ", indent*2) + "/" + key + "\n"
		result += value.string(indent + 1)
	}

	return result
}
