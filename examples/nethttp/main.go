// net/http example: same component as minimal, served over HTTP. No Fiber or other framework.
package main

import (
	"net/http"

	gc "github.com/abdheshnayak/gohtmlx/examples/nethttp/dist/gohtmlxc"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "GoHTMLX"
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		el := gc.Hello{Name: name, Attrs: nil}.Get()
		if _, err := el.Render(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	http.ListenAndServe(":8080", nil)
}
