package basichttprouter

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type StaticMuxer struct {
	routes  map[string]Muxer
	handler http.Handler
}

func (m *StaticMuxer) Handle(path string, handler http.Handler) error {
	// keeps only non-empty strings between slashes
	spath := strings.FieldsFunc(path, func(c rune) bool { return c == '/' })
	fmt.Printf("path: %s, len: %d\n", spath, len(spath))

	m.handle(spath, handler)

	return nil
}

func (m *StaticMuxer) handle(spath []string, handler http.Handler) error {
	if len(spath) == 0 {
		m.handler = handler
		return nil
	}

	subm, found := m.routes[spath[0]]
	if !found {
		subm = NewStaticMuxer()
		m.routes[spath[0]] = subm
	}

	subm.handle(spath[1:], handler)

	return nil
}

func (m *StaticMuxer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var spath []string
	spathctx := r.Context().Value("spath")

	if spathctx == nil {
		spath = strings.FieldsFunc(r.URL.Path, func(c rune) bool { return c == '/' })
	} else {
		spath = spathctx.([]string)
	}

	if len(spath) == 0 {
		m.handler.ServeHTTP(w, r)
		return
	}

	subm, found := m.routes[spath[0]]
	if !found {
		http.NotFound(w, r)
		return
	}

	subm.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "spath", spath[1:])))
}

func NewStaticMuxer() *StaticMuxer {
	return &StaticMuxer{
		routes:  map[string]Muxer{},
		handler: http.NotFoundHandler(),
	}
}
