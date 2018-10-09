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

	"net/http"

	"github.com/kingnido/basichttprouter"
)

func main() {
	r := basichttprouter.NewRouter()

	r.Handle("/api/posts/:pid/comments", hello)
	r.Handle("/api/comments/:cid", commentHandler)
	r.Handle("/api/posts/:pid/comments", addCommentHandler)
	r.Handle("/static", hello)
	r.Handle("/", hello)

	fmt.Print(r)

	http.ListenAndServe(":3000", r)
}

var hello = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "hello")
})

var commentHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, r.Context().Value("cid"))
})

var addCommentHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, r.Context().Value("pid"))
})
```
