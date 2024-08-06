package main

import (
	"time"

	"github.com/anton2920/gofa/net/html"
	"github.com/anton2920/gofa/net/http"
)

const (
	WidthMedium = 6
)

func DisplayHTMLStart(w *http.Response) {
	w.WriteString(html.Header)
}

func DisplayHeadStart(w *http.Response) {
	w.WriteString(`<head>`)
}

func DisplayHeadEnd(w *http.Response) {
	w.WriteString(`</head>`)
}

func DisplayBodyStart(w *http.Response) {
	w.WriteString(`<body>`)
}

func DisplayFormattedTime(w *http.Response, t int64) {
	w.Write(time.Unix(t, 0).AppendFormat(make([]byte, 0, 20), "2006/01/02 15:04:05"))
}

func DisplayLabel(w *http.Response, l Language, label string) {
	w.WriteString(`<label>`)
	w.WriteString(label)
	w.WriteString(`:<br></label>`)
}

func DisplayConstraintInput(w *http.Response, t string, minLength int, maxLength int, name string, value string, required bool) {
	w.WriteString(`<input type="`)
	w.WriteString(t)
	w.WriteString(`" minlength="`)
	w.WriteInt(minLength)
	w.WriteString(`" maxlength="`)
	w.WriteInt(maxLength)
	w.WriteString(`" name="`)
	w.WriteString(name)
	w.WriteString(`" value="`)
	w.WriteHTMLString(value)
	w.WriteString(`"`)
	if required {
		w.WriteString(` required`)
	}
	w.WriteString(`>`)
}

func DisplaySubmit(w *http.Response, l Language, name string, value string) {
	w.WriteString(`<input type="submit" name="`)
	w.WriteString(name)
	w.WriteString(`" value="`)
	w.WriteString(Ls(l, value))
	w.WriteString(`">`)
}

func DisplayBodyEnd(w *http.Response) {
	w.WriteString(`</body>`)
}

func DisplayHTMLEnd(w *http.Response) {
	w.WriteString(`</html>`)
}
