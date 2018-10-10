package basichttprouter

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type Router struct {
	root *Node
}

type Node struct {
	route   *Path
	handler http.Handler
}

type Path struct {
	routes map[string]*Node
	route  *Node
	param  string
}

func (r *Router) Handle(path string, handler http.Handler) error {
	spath := strings.FieldsFunc(path, func(c rune) bool { return c == '/' })

	if r.root == nil {
		r.root = &Node{route: nil, handler: nil}
	}

	return r.root.handle(spath, handler)
}

func (n *Node) handle(spath []string, handler http.Handler) error {
	if len(spath) == 0 {
		n.handler = handler
		return nil
	}

	if n.route == nil {
		n.route = &Path{routes: map[string]*Node{}, route: nil, param: ""}
	}

	return n.route.handle(spath, handler)
}

func (p *Path) handle(spath []string, handler http.Handler) error {
	part := spath[0]
	rest := spath[1:]

	if p.route != nil {
		if part != p.param {
			return errors.New("Always use the same key for parameters")
		}

		return p.route.handle(rest, handler)
	}

	if len(p.routes) > 0 {
		if part[0] == ':' {
			return errors.New("Non parametric route is being used here")
		}

		_, found := p.routes[part]
		if !found {
			p.routes[part] = &Node{route: nil, handler: nil}
		}

		return p.routes[part].handle(rest, handler)
	}

	if part[0] == ':' {
		p.param = part
		p.route = &Node{route: nil, handler: nil}

		return p.route.handle(rest, handler)
	} else {
		p.routes[part] = &Node{route: nil, handler: nil}

		return p.routes[part].handle(rest, handler)
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	spath := strings.FieldsFunc(r.URL.Path, func(c rune) bool { return c == '/' })

	router.root.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "spath", spath)))
}

func (n *Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	spath := r.Context().Value("spath").([]string)

	if len(spath) == 0 {
		n.handler.ServeHTTP(w, r)
		return
	}

	if n.route == nil {
		http.NotFound(w, r)
	}

	n.route.ServeHTTP(w, r)
}

func (p *Path) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	spath := ctx.Value("spath").([]string)
	part := spath[0]
	rest := spath[1:]

	var next *Node

	if p.param != "" {
		ctx = context.WithValue(ctx, p.param, part)
		next = p.route
	} else {
		_, found := p.routes[part]
		if !found {
			http.NotFound(w, r)
			return
		}
		next = p.routes[part]
	}

	next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, "spath", rest)))
}

func NewRouter() *Router {
	return &Router{nil}
}
