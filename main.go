package main

import (
	"runtime"

	"github.com/anton2920/gofa/debug"
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/event"
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/net/html"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/session"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/trace"
)

func Router(w *http.Response, r *http.Request, s session.Session) error {
	defer trace.End(trace.Begin(""))

	h := html.New(w, nil, s, Styles)
	db := Database{Language: s.Language}

	if err := func(h *html.HTML, r *http.Request, db *Database) error {
		defer trace.End(trace.Begin(""))

		if r.Error != nil {
			return r.Error
		}

		path := r.URL.Path
		switch path {
		default:
			switch {
			default:
				switch path {
				case "/":
					return IndexPage(h)
				}
			case strings.StartsWith(path, "/user"):
				switch path[len("/user"):] {
				case "/signin":
					return UserSignin(h, r, db)
				case "/signout":
				case "/signup":
					return UserSignup(h, r, db)
				}
			}
		case "/error":
			return http.ServerError(errors.New("test error"))
		case "/panic":
			panic("test panic")
		}

		return http.NotFound("requested item does not exist")
	}(&h, r, &db); err != nil {
		return ErrorPage(&h, err)
	}

	return nil
}

func main1() {
	const addr = "0.0.0.0:7075"

	if debug.Debug {
		log.SetLevel(log.LevelDebug)
	}

	trace.BeginProfile()
	defer trace.EndAndPrintProfile()

	if err := OpenDB("db.sqlite"); err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	l, err := http.Listen(addr)
	if err != nil {
		log.Fatalf("Failed to listen for HTTP connections: %v", err)
	}
	log.Infof("Listening on %s... (%s, GOMAXPROCS=%d)", addr, runtime.Version(), runtime.GOMAXPROCS(0))

	q, err := event.NewQueue()
	if err != nil {
		log.Fatalf("Failed to create new event queue: %v", err)
	}

	_ = syscall.IgnoreSignals(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	_ = q.AddSocket(int32(l.Socket), event.RequestRead, event.TriggerEdge, nil)
	_ = q.AddSignals(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	events := make([]event.Event, 64)

	var quit bool
	for !quit {
		n, err := q.GetEvents(events)
		if err != nil {
			log.Errorf("Failed to get events from queue: %v", err)
			continue
		}

		for i := 0; i < n; i++ {
			e := &events[i]

			switch e.Type {
			case event.TypeRead:
				c, err := l.Accept()
				if err != nil {
					log.Errorf("Failed to accept new HTTP connection: %v", err)
					continue
				}
				go http.ConnectionHandler(c, Router)
			case event.TypeSignal:
				sig := syscall.Signal(e.Identifier)
				log.Infof("Received %d (%s), exitting...", sig, sig)
				quit = true
			}
		}
	}

	if err := l.Close(); err != nil {
		log.Warnf("Failed to close HTTP listener: %v", err)
	}
}

func main() {
	main1()
}
