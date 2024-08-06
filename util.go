package main

import (
	"strconv"

	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/prof"
	"github.com/anton2920/gofa/strings"
)

func GetIDFromURL(l Language, u url.URL, prefix string) (database.ID, error) {
	defer prof.End(prof.Begin(""))

	if !strings.StartsWith(u.Path, prefix) {
		return 0, http.NotFound(Ls(l, "requested page does not exist"))
	}

	id, err := strconv.Atoi(u.Path[len(prefix):])
	if (err != nil) || (id < 0) || (id >= (1 << 31)) {
		return 0, http.BadRequest(Ls(l, "invalid ID for %q"), prefix)
	}

	return database.ID(id), nil
}
