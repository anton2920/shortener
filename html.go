package main

import "github.com/anton2920/gofa/net/html"

var Styles = html.Theme{
	Body: html.Class("p-md-4"),
}

func DisplayStyles(h *html.HTML) {
	h.Link(html.Attributes{Href: "https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css", Rel: "stylesheet"})
}
