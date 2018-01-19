### Sample
```
package main

import (
	"net/http"

	"session"
	_ "session/mem"
)

func main() {
	http.ListenAndServe(":8888", http.HandlerFunc(indexHandler))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	m, err := session.NewManager("mem", 1)
	if err != nil {
		panic(err)
	}
	sess, err := m.SessionRead(r)
	if err != nil {
		sess, _ = m.SessionCreate(w)
	}

	sess.Set("name", "giraffate")
	if name, err := sess.Get("name"); err == nil {
		w.Write([]byte("Hello, " + name.(string) + "!\n"))
	}
}
```
