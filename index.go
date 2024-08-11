package main

import (
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/trace"
)

const (
	MinURLLen = 1
	MaxURLLen = 128
)

func IndexPage(w *http.Response, r *http.Request, shortened string, ierr error) error {
	defer trace.End(trace.Begin(""))

	const title = "URL shortener"

	session, _ := GetSessionFromRequest(r)

	if err := r.ParseForm(); err != nil {
		return http.ClientError(err)
	}

	DisplayHTMLStart(w)

	DisplayHeadStart(w)
	{
		w.WriteString(`<title>`)
		w.WriteString(Ls(GL, title))
		w.WriteString(`</title>`)
	}
	DisplayHeadEnd(w)

	DisplayBodyStart(w)
	{
		w.WriteString(`<h2>`)
		w.WriteString(Ls(GL, title))
		w.WriteString(`</h2>`)

		if session == nil {
			w.WriteString(`<a href="/user/signin">`)
			w.WriteString(Ls(GL, "Sign in"))
			w.WriteString(`</a>`)

			w.WriteString(` <a href="/user/signup">`)
			w.WriteString(Ls(GL, "Sign up"))
			w.WriteString(`</a>`)
		} else {
			var user User
			if err := GetUserByID(session.ID, &user); err != nil {
				return http.ServerError(err)
			}

			w.WriteString(`<a href="/user/`)
			w.WriteID(session.ID)
			w.WriteString(`">`)
			DisplayUserTitle(w, &user)
			w.WriteString(`</a>`)

			w.WriteString(` <a href="` + APIPrefix + `/user/signout">`)
			w.WriteString(Ls(GL, "Sign out"))
			w.WriteString(`</a>`)
		}
		w.WriteString(`<br><br>`)

		DisplayError(w, GL, ierr)

		if len(shortened) > 0 {
			w.WriteString(`<label>`)
			w.WriteString(Ls(GL, "Shortened URL"))
			w.WriteString(`: `)
			w.WriteString(`<a href="/`)
			w.WriteString(shortened)
			w.WriteString(`">`)
			w.WriteString(shortened)
			w.WriteString(`</a>`)
			w.WriteString(`</label>`)
			w.WriteString(`<br><br>`)
		}

		w.WriteString(`<form method="POST" action="` + APIPrefix + `/url/create">`)
		{
			w.WriteString(`<label>`)
			w.WriteString(Ls(GL, "URL"))
			w.WriteString(`: `)
			DisplayConstraintInput(w, "text", MinURLLen, MaxURLLen, "URL", r.Form.Get("URL"), true)
			w.WriteString(`</label>`)
			w.WriteString(`<br><br>`)

			DisplaySubmit(w, GL, "", "Shorten!")
		}
		w.WriteString(`</form>`)
	}
	DisplayBodyEnd(w)

	DisplayHTMLEnd(w)
	return nil
}
