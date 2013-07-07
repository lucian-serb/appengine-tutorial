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

type User struct {
	IsLoggedIn bool
	Greetings  []Greeting
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/sign", sign)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
}

const guestbookTemplateHTML = `
<!doctype html>
<html>
	<body>
		{{range .Greetings}}
			{{with .Author}}
				<p><b>{{.}}</b> wrote:</p>
			{{else}}
				<p><b>An anonymous person</b> wrote:</p>
			{{end}}
			<pre>{{.Content}}</pre>
		{{end}}
		<form action="/sign" method="post">
			<div><textarea name="content" rows="3" cols="60"></textarea></div>
			<div>
				<input type="submit" value="Sign Guestbook">
				{{if .IsLoggedIn}}
					<input type="submit" value="Logout" form="logout">
				{{else}}
					<input type="submit" value="Login" form="login">
				{{end}}
			</div>
		</form>
		{{if .IsLoggedIn}}
			<form action="/logout" method="post" id="logout"></form>
		{{else}}
			<form action="/login" method="post" id="login"></form>
		{{end}}
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

	isLoggedIn := user.Current(c) != nil
	usr := User{
		isLoggedIn,
		greetings,
	}

	if err := guestbookTemplate.Execute(w, usr); err != nil {
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

func login(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if url, err := user.LoginURL(c, "/"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if url, err := user.LogoutURL(c, "/"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}
