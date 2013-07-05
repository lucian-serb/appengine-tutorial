package hello

import (
	"appengine"
	"appengine/user"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)

	if u == nil {
		if url, err := user.LoginURL(c, r.URL.String()); err == nil {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	fmt.Fprintf(w, "Hello %v !!!", u)
}
