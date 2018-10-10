# basichttprouter

Project to start playing with golang.

# Done

- Static routing
- Parametric routing

# To do

- Docs
- Tests
- Middleware support
- Methods support

# Example

```go
package main

import (
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/kingnido/basichttprouter"
)

func main() {
	m := basichttprouter.NewRouter()

	m.Handle("/api", dumpContext)
	m.Handle("/api/posts", dumpContext)
	m.Handle("/api/posts/:id", dumpVars(":id"))
	m.Handle("/api/posts/:id/comments", dumpVars(":id"))

	http.ListenAndServe(":3000", timeLogger(m))
}

var dumpContext = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, r.Context())
})

func timeLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Println("Elapsed:", time.Since(start))
	})
}

func dumpVars(vars ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := map[string]string{}

		for _, key := range vars {
			obj := r.Context().Value(key)
			key = string(key[1:])

			if obj != nil {
				result[key] = obj.(string)
			} else {
				result[key] = "nil"
			}
		}

		fmt.Fprint(w, result)
	})
}
```
