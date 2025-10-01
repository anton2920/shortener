package main

import "github.com/anton2920/gofa/net/html"

func IndexPage(h *html.HTML) error {
	h.Begin()

	h.HeadBegin()
	{
		DisplayStyles(h)
		h.Title("Home | Shortener")
	}
	h.HeadEnd()

	h.BodyBegin()
	{
		h.H1("Hello, world!")

	}
	h.BodyEnd()

	h.End()
	return nil
}
