package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync/atomic"
	"unsafe"

	"github.com/anton2920/gofa/alloc"
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/event"
	"github.com/anton2920/gofa/intel"
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/net/http/http1"
	"github.com/anton2920/gofa/net/tcp"
	"github.com/anton2920/gofa/prof"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/time"
)

const (
	APIPrefix = "/api"
	FSPrefix  = "/fs"
)

var (
	BuildMode string
	Debug     bool
)

var WorkingDirectory string

var DateBufferPtr unsafe.Pointer

func HandlePageRequest(w *http.Response, r *http.Request, path string) error {
	switch {
	case path == "/":
		return IndexPage(w, r)
	case strings.StartsWith(path, "/user"):
		switch path[len("/user"):] {
		default:
			return UserPage(w, r)
		case "/signin":
			return UserSigninPage(w, r, nil)
		case "/signup":
			return UserSignupPage(w, r, nil)
		}
	}

	return http.NotFound(Ls(GL, "requested page does not exist"))
}

func HandleAPIRequest(w *http.Response, r *http.Request, path string) error {
	switch {
	case strings.StartsWith(path, "/user"):
		switch path[len("/user"):] {
		case "/signin":
			return UserSigninHandler(w, r)
		case "/signout":
			return UserSignoutHandler(w, r)
		case "/signup":
			return UserSignupHandler(w, r)
		}
	}

	return http.NotFound(Ls(GL, "requested API endpoint does not exist"))
}

/* TODO(anton2920): maybe switch to sendfile(2)? */
func HandleFSRequest(w *http.Response, r *http.Request, path string) error {
	return http.NotFound(Ls(GL, "requested file does not exist"))
}

func RouterFunc(w *http.Response, r *http.Request) (err error) {
	defer prof.End(prof.Begin(""))

	defer func() {
		if p := recover(); p != nil {
			err = errors.NewPanic(p)
		}
	}()

	path := r.URL.Path
	switch {
	default:
		return HandlePageRequest(w, r, path)
	case strings.StartsWith(path, APIPrefix):
		return HandleAPIRequest(w, r, path[len(APIPrefix):])
	case strings.StartsWith(path, FSPrefix):
		return HandleFSRequest(w, r, path[len(FSPrefix):])

	case path == "/error":
		return http.ServerError(errors.New(Ls(GL, "test error")))
	case path == "/panic":
		panic(Ls(GL, "test panic"))
	}
}

func Router(ctx *http.Context, ws []http.Response, rs []http.Request) {
	defer prof.End(prof.Begin(""))

	for i := 0; i < len(rs); i++ {
		w := &ws[i]
		r := &rs[i]

		start := intel.RDTSC()
		w.Headers.Set("Content-Type", `text/html; charset="UTF-8"`)
		level := log.LevelDebug

		err := RouterFunc(w, r)
		if err != nil {
			ErrorPageHandler(w, r, GL, err)
			if (w.StatusCode >= http.StatusBadRequest) && (w.StatusCode < http.StatusInternalServerError) {
				level = log.LevelWarn
			} else {
				level = log.LevelError
			}
			http.CloseAfterWrite(ctx)
		}

		if r.Headers.Get("Connection") == "close" {
			w.Headers.Set("Connection", "close")
			http.CloseAfterWrite(ctx)
		}

		addr := ctx.ClientAddress
		if r.Headers.Has("X-Forwarded-For") {
			addr = r.Headers.Get("X-Forwarded-For")
		}
		end := intel.RDTSC()
		elapsed := end - start

		log.Logf(level, "[%21s] %7s %s -> %v (%v), %4dµs", addr, r.Method, r.URL.Path, w.StatusCode, err, elapsed.ToUsec())
	}
}

func GetDateHeader() []byte {
	defer prof.End(prof.Begin(""))

	return unsafe.Slice((*byte)(atomic.LoadPointer(&DateBufferPtr)), time.RFC822Len)
}

func UpdateDateHeader(now int) {
	buffer := make([]byte, time.RFC822Len)
	time.PutTmRFC822(buffer, time.ToTm(now))
	atomic.StorePointer(&DateBufferPtr, unsafe.Pointer(&buffer[0]))
}

func ServerWorker(q *event.Queue) {
	events := make([]event.Event, 64)

	const batchSize = 32
	ws := make([]http.Response, batchSize)
	rs := make([]http.Request, batchSize)

	getEvents := func(q *event.Queue, events []event.Event) (int, error) {
		defer prof.End(prof.Begin("github.com/anton2920/gofa/event.(*Queue).GetEvents"))
		return q.GetEvents(events)
	}

	for {
		n, err := getEvents(q, events)
		if err != nil {
			log.Errorf("Failed to get events from client queue: %v", err)
			continue
		}
		dateBuffer := GetDateHeader()

		for i := 0; i < n; i++ {
			e := &events[i]
			if errno := e.Error(); errno != 0 {
				log.Errorf("Event for %v returned code %d (%s)", e.Identifier, errno, errno)
				continue
			}

			ctx, ok := http.GetContextFromPointer(e.UserData)
			if !ok {
				continue
			}
			if e.EndOfFile() {
				http.Close(ctx)
				continue
			}

			switch e.Type {
			case event.Read:
				var read int
				for read < e.Data {
					n, err := http.Read(ctx)
					if err != nil {
						if err == http.NoSpaceLeft {
							http1.FillError(ctx, err, dateBuffer)
							http.CloseAfterWrite(ctx)
							break
						}
						log.Errorf("Failed to read data from client: %v", err)
						http.Close(ctx)
						break
					}
					read += n

					for n > 0 {
						n, err = http1.ParseRequestsUnsafe(ctx, rs)
						if err != nil {
							http1.FillError(ctx, err, dateBuffer)
							http.CloseAfterWrite(ctx)
							break
						}
						Router(ctx, ws[:n], rs[:n])
						http1.FillResponses(ctx, ws[:n], dateBuffer)
					}
				}
				fallthrough
			case event.Write:
				_, err = http.Write(ctx)
				if err != nil {
					log.Errorf("Failed to write data to client: %v", err)
					http.Close(ctx)
					continue
				}
			}
		}
	}
}

func main() {
	var err error

	nworkers := min(runtime.GOMAXPROCS(0)/2, runtime.NumCPU())
	switch BuildMode {
	default:
		BuildMode = "Release"
	case "Debug":
		Debug = true
		log.SetLevel(log.LevelDebug)
	case "Profiling":
		f, err := os.Create(fmt.Sprintf("masters-cpu.pprof"))
		if err != nil {
			log.Fatalf("Failed to create a profiling file: %v", err)
		}
		defer f.Close()

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	case "gofa/prof":
		nworkers = 1

		prof.BeginProfile()
		defer prof.EndAndPrintProfile()
	}
	log.Infof("Starting Shortener in %q mode...", BuildMode)

	if err := RestoreSessionsFromFile(SessionsFile); err != nil {
		log.Warnf("Failed to restore sessions from file: %v", err)
	}

	const address = "0.0.0.0:7075"
	l, err := tcp.Listen(address, 128)
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	defer syscall.Close(l)

	log.Infof("Listening on %s...", address)

	q, err := event.NewQueue()
	if err != nil {
		log.Fatalf("Failed to create listener event queue: %v", err)
	}
	defer q.Close()

	_ = q.AddSocket(l, event.RequestRead, event.TriggerEdge, nil)
	_ = q.AddTimer(1, 1, event.Seconds, nil)

	_ = syscall.IgnoreSignals(syscall.SIGINT, syscall.SIGTERM)
	_ = q.AddSignals(syscall.SIGINT, syscall.SIGTERM)

	ctxPool := alloc.NewSyncPool[http.Context](nworkers * 512)
	qs := make([]*event.Queue, nworkers)
	for i := 0; i < nworkers; i++ {
		qs[i], err = event.NewQueue()
		if err != nil {
			log.Fatalf("Failed to create new client queue: %v", err)
		}
		go ServerWorker(qs[i])
	}
	now := time.Unix()
	UpdateDateHeader(now)

	events := make([]event.Event, 64)
	var counter int

	var quit bool
	for !quit {
		n, err := q.GetEvents(events)
		if err != nil {
			log.Errorf("Failed to get events: %v", err)
			continue
		}

		for i := 0; i < n; i++ {
			e := &events[i]

			switch e.Type {
			default:
				log.Panicf("Unhandled event: %#v", e)
			case event.Read:
				ctx, err := http.Accept(l, &ctxPool, 1024)
				if err != nil {
					if err == http.TooManyClients {
						http1.FillError(ctx, err, GetDateHeader())
						http.Write(ctx)
						http.Close(ctx)
					}
					log.Errorf("Failed to accept new HTTP connection: %v", err)
					continue
				}
				_ = qs[counter%len(qs)].AddHTTP(ctx, event.RequestRead, event.TriggerEdge)
				counter++
			case event.Timer:
				now += e.Data
				UpdateDateHeader(now)
			case event.Signal:
				log.Infof("Received signal %d, exitting...", e.Identifier)
				quit = true
				break
			}
		}
	}

	if err := StoreSessionsToFile(SessionsFile); err != nil {
		log.Warnf("Failed to store sessions to file: %v", err)
	}
}
