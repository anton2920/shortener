package main

import (
	"runtime"

	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/net/html"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/session"
)

func Router(w *http.Response, r *http.Request, s session.Session) error {
	h := html.New(w, r, s, Styles)

	if err := func(h *html.HTML) error {
		if r.Error != nil {
			return r.Error
		}

		switch r.URL.Path {
		case "/":
			IndexPage(h)
		case "/plaintext":
			w.WriteString("Hello, world!")
		case "/error":
			return http.ServerError(errors.New("test error"))
		case "/panic":
			panic("test panic")
		}

		return nil
	}(&h); err != nil {
		return ErrorPage(&h, err)
	}

	return nil
}

func main1() {
	const addr = "0.0.0.0:7075"

	l, err := http.Listen(addr)
	if err != nil {
		log.Fatalf("Failed to listen for HTTP connections: %v", err)
	}

	log.Infof("Listening on %s... (%s)", addr, runtime.Version())

	for {
		c, err := l.Accept()
		if err != nil {
			log.Errorf("Failed to accept new HTTP connection: %v", err)
			continue
		}
		go http.ConnectionHandler(c, Router)
	}

	if err := l.Close(); err != nil {
		log.Warnf("Failed to close HTTP listener: %v", err)
	}
}

func main() {
	main1()
}
