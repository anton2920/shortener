package main

import (
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/prof"
)

const (
	MinURLLen = 1
	MaxURLLen = 128
)

func IndexPage(w *http.Response, r *http.Request) error {
	defer prof.End(prof.Begin(""))

	const title = "URL shortener"

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
		w.WriteString(`<h1>`)
		w.WriteString(Ls(GL, title))
		w.WriteString(`</h1>`)

		w.WriteString(`<form method="POST" action="/">`)
		{
			w.WriteString(`<label>`)
			w.WriteString(Ls(GL, "URL"))
			w.WriteString(`: `)
			DisplayConstraintInput(w, "text", MinURLLen, MaxURLLen, "URL", r.Form.Get("URL")+"awdawdawd", true)
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
