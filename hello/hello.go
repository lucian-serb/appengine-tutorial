package hello

import (
	"fmt"
	"html/template"
	"net/http"
)

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/sign", sign)
}

const guestbookForm = `
<!doctype html>
<html>
	<body>
		<form action="/sign" method="post">
			<div><textarea name="content" rows="3" cols="60"></textarea></div>
			<div><input type="submit" value="Sign Guestbook"></div>
		</form>
	</body>
</html>
`

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, guestbookForm)
}

const signTemplateHTML = `
<!doctype html>
<html>
	<body>
		<p>You wrote:</p>
		<pre>{{.}}</pre>
	</body>
</html>
`

var signTemplate = template.Must(template.New("sign").Parse(signTemplateHTML))

func sign(w http.ResponseWriter, r *http.Request) {
	if err := signTemplate.Execute(w, r.FormValue("content")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
