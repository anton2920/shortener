package main

import (
	"net/url"
	"sync"
	"unsafe"

	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/time"
)

type URL struct {
	ID    database.ID
	Flags int32

	RawURL    string
	ExpiresAt int64

	RedirectCounts map[int64]int64
	RedirectFrom   map[string]int64
}

const (
	FlagActive  int32 = 0
	FlagDeleted       = 1
	FlagPrivate       = 2
)

var (
	URLs     = make(map[string]URL)
	URLsLock sync.RWMutex
)

func GetURLByID(id database.ID, url *URL) error {
	URLsLock.RLock()
	defer URLsLock.RUnlock()

	for _, v := range URLs {
		if v.ID == id {
			*url = v
			return nil
		}
	}

	return database.NotFound
}

func GetURLByPath(path string, url *URL) error {
	URLsLock.RLock()
	u, ok := URLs[path]
	URLsLock.RUnlock()
	if !ok {
		return database.NotFound
	}

	*url = u
	return nil
}

func CreateURL(path string, url *URL) error {
	URLsLock.Lock()

	url.ID = database.ID(len(URLs) + 1)
	URLs[path] = *url

	URLsLock.Unlock()
	return nil
}

func SaveURL(path string, url *URL) error {
	URLsLock.Lock()

	URLs[path] = *url

	URLsLock.Unlock()
	return nil
}

func URLCreateHandler(w *http.Response, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return http.ClientError(err)
	}

	rawURL := r.Form.Get("URL")
	if len(rawURL) == 0 {
		return IndexPage(w, r, "", http.BadRequest("provided URL is empty"))
	}
	_, err := url.Parse(rawURL)
	if err != nil {
		return IndexPage(w, r, "", http.BadRequest("provided URL is incorrect: %v", err))
	}

	/* TODO(anton2920): verify URL is not shortened by our system before. */

	buffer := make([]byte, len(rawURL))
	copy(buffer, rawURL)
	rawURL = string(buffer)

	for {
		if false {
			const shortenedLen = 7
			buffer = make([]byte, shortenedLen)
			SlicePutRandomBase52(buffer)
			buffer[shortenedLen/2] = '-'
		} else {
			const shortenedLen = 12
			buffer = make([]byte, shortenedLen)
			SlicePutRandomBase26(buffer)
			buffer[3] = '-'
			buffer[8] = '-'
		}

		URLsLock.RLock()
		_, ok := URLs[unsafe.String(unsafe.SliceData(buffer), len(buffer))]
		URLsLock.RUnlock()
		if !ok {
			break
		}
	}
	shortened := string(buffer)

	var url URL

	url.RawURL = rawURL
	url.RedirectCounts = make(map[int64]int64)
	url.RedirectFrom = make(map[string]int64)

	if err := CreateURL(shortened, &url); err != nil {
		return http.ServerError(err)
	}

	return IndexPage(w, r, shortened, nil)
}

func URLRedirectHandler(w *http.Response, r *http.Request, path string) error {
	var url URL

	if err := GetURLByPath(path, &url); err != nil {
		if err == database.NotFound {
			return http.NotFound("shortened URL does not exist")
		}
		return http.ServerError(err)
	}
	defer SaveURL(path, &url)

	now := int64(time.Unix() / 60 * 60 * 24)
	url.RedirectCounts[now]++
	url.RedirectFrom[r.Headers.Get("Referer")]++

	w.Redirect(url.RawURL, http.StatusSeeOther)
	return nil
}
