# basichttprouter

Project to start playing with golang.

# Done

- Basic static routing

# To do

- Testing
- Parameters in the URL
- Setting handlers by method
- Sanitize URL

# Example

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kingnido/basichttprouter"
	"net/http"
)

func main() {
	m := basichttprouter.NewStaticMuxer()

	m.Handle("", dumpContext)
	m.Handle("/", dumpContext)
	m.Handle("/api", dumpContext)
	m.Handle("/api/go/", dumpContext)
	m.Handle("/static/", dumpContext)

	http.ListenAndServe(":3000", timeLogger(m))
}

var dumpContext = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.Context())
})

func timeLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Println("Elapsed:", time.Since(start))
	})
}
```
