package main

import "github.com/anton2920/gofa/net/html"

var (
	Styles = html.Theme{
		A:      html.Class("text-decoration-none"),
		Body:   html.Class("p-md-4"),
		Button: html.Class("btn"),
		Form:   html.Class("col-lg-4"),
		Input:  html.Class("form-control"),
	}

	StyleButtonSubmit = html.Class("btn-primary w-100")
)

func DisplayStyles(h *html.HTML) {
	h.Link(html.Attributes{Href: "https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css", Rel: "stylesheet"})
}
