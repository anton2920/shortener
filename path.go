package main

import (
	"fmt"
	"strconv"

	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace"
)

type Path string

func GetValidID(si string) (database.ID, error) {
	id, err := strconv.Atoi(si)
	if err != nil {
		return 0, err
	}
	if (id < database.MinValidID+1) || (id > database.MaxValidID) {
		return 0, errors.New("ID out of range")
	}
	return database.ID(id), nil
}

func GetIDFromPath(path Path, skip int) database.ID {
	var id database.ID
	var err error
	var i int

	for (len(path) > 0) && (i < skip+1) {
		id = 0

		slash := strings.FindChar(string(path), '/')
		if slash == -1 {
			break
		}
		path = path[slash+1:]

		nextSlash := strings.FindChar(string(path), '/')
		if nextSlash == -1 {
			nextSlash = len(path)
		}

		id, err = GetValidID(string(path[:nextSlash]))
		if err == nil {
			i++
		}
	}

	return id
}

/* Match returns `true` if `p` matches format described in `format`. Additionally it slices `p` with the length of the matched string. */
func (p *Path) Match(format string, args ...interface{}) bool {
	defer trace.End(trace.Begin(""))

	var narg int
	var ok bool

	path := string(*p)
	for {
		percent := strings.FindChar(format, '%')
		if percent == -1 {
			const ellipsis = "..."

			if !strings.EndsWith(format, ellipsis) {
				ok = path == format
			} else {
				format = format[:len(format)-len(ellipsis)]
				ok = strings.StartsWith(path, format)
			}
			if ok {
				*p = Path(path[len(format):])
			}

			return ok
		}

		match := format[:percent]
		if !strings.StartsWith(path, match) {
			return false
		}
		path = path[len(match):]
		format = format[len(match):]

		slashP := strings.FindChar(path, '/')
		if slashP == -1 {
			slashP = len(path)
		}

		slashF := strings.FindChar(format, '/')
		if slashF == -1 {
			slashF = len(format)
		}

		n, err := fmt.Sscanf(path[:slashP], format[:slashF], args[narg:]...)
		if (n == 0) && (err != nil) {
			return false
		}
		narg += n

		path = path[slashP:]
		format = format[2:]
	}
}
