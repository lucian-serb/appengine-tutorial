package hello

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"html/template"
	"net/http"
	"time"
)

type Greeting struct {
	Author  string
	Content string
	Date    time.Time
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/sign", sign)
}

const guestbookTemplateHTML = `
<!doctype html>
<html>
	<body>
		{{range .}}
			{{with .Author}}
				<p><b>{{.}}</b> wrote:</p>
			{{else}}
				<p><b>An anonymous person</b> wrote:</p>
			{{end}}
			<pre>{{.Content}}</pre>
		{{end}}
		<form action="/sign" method="post">
			<div><textarea name="content" rows="3" cols="60"></textarea></div>
			<div><input type="submit" value="Sign Guestbook"></div>
		</form>
	</body>
</html>
`

var guestbookTemplate = template.Must(template.New("book").Parse(guestbookTemplateHTML))

func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Greeting").Order("-Date").Limit(10)
	greetings := make([]Greeting, 0, 10)

	if _, err := q.GetAll(c, &greetings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := guestbookTemplate.Execute(w, greetings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func sign(w http.ResponseWriter, r *http.Request) {
	if content := r.FormValue("content"); len(content) != 0 {
		c := appengine.NewContext(r)
		g := Greeting{
			Content: content,
			Date:    time.Now(),
		}

		if u := user.Current(c); u != nil {
			g.Author = u.String()
		}

		if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Greeting", nil), &g); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
