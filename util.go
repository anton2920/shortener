package main

import (
	"math/rand/v2"
	"strconv"

	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace"
)

func GetIDFromURL(l Language, u url.URL, prefix string) (database.ID, error) {
	defer trace.End(trace.Begin(""))

	if !strings.StartsWith(u.Path, prefix) {
		return 0, http.NotFound(Ls(l, "requested page does not exist"))
	}

	id, err := strconv.Atoi(u.Path[len(prefix):])
	if (err != nil) || (id < 0) || (id >= (1 << 31)) {
		return 0, http.BadRequest(Ls(l, "invalid ID for %q"), prefix)
	}

	return database.ID(id), nil
}

func SlicePutRandomBase26(buffer []byte) {
	const letters = "abcdefghijklmnopqrstuvwxyz"

	for i := 0; i < len(buffer); i++ {
		buffer[i] = letters[rand.Int()%len(letters)]
	}
}

func SlicePutRandomBase52(buffer []byte) {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	for i := 0; i < len(buffer); i++ {
		buffer[i] = letters[rand.Int()%len(letters)]
	}
}
