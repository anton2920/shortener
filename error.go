package main

import "github.com/anton2920/gofa/net/html"

func ErrorPage(h *html.HTML, ierr error) error {
	h.Begin()

	h.HeadBegin()
	{
		DisplayStyles(h)
		h.Title("Error")
	}
	h.HeadEnd()

	h.BodyBegin()
	{
		h.H1("Error!")
		h.Error(ierr)
	}
	h.BodyEnd()

	h.End()
	return ierr
}
