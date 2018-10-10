package basichttprouter

import (
	"net/http"
)

type Muxer interface {
	http.Handler
	Handle(string, http.Handler) error
	handle([]string, http.Handler) error
}
