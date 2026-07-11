package main

import (
	"github.com/anton2920/gofa/net/html"
	"github.com/anton2920/gofa/net/http"
)

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
		h.H2("World fastest shortener!")

		if h.ID == 0 {
			h.A("/user/signup", "Sign up")
			h.A("/user/signin", "Sign in")
		} else {
			h.A(h.PathID("/user/", h.ID), "Profile")
			h.A("/user/signout", "Sign out")
		}
	}
	h.BodyEnd()

	h.End()
	return nil
}

func IndexHandler(w *http.Response, r *http.Request) error {
	h := html.New(w, r, &Styles)

	switch r.Method {
	default:
		return IndexPage(&h)
	}
}
