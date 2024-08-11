package main

import (
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/trace"
)

type Language int32

const (
	EN Language = iota
	RU
	FR
	XX
)

var Language2String = [...]string{
	EN: "English",
	RU: "Русский",
	FR: "Français",
}

/* TODO(anton2920): remove '([A-Z]|[a-z])[a-z]+' duplicates. */
var Localizations = map[string]*[XX]string{
	"Shorten!": {
		RU: "Sokraryt'",
	},
	"URL": {
		RU: "Ssylka",
	},
	"URL shortener": {
		RU: "Sokrashshyatel' ssylok",
	},
}

var GL = EN

func (l Language) String() string {
	defer trace.End(trace.Begin(""))

	return Language2String[l]
}

func Ls(l Language, s string) string {
	defer trace.End(trace.Begin(""))

	if l == EN {
		return s
	}

	ls := Localizations[s]
	if (ls == nil) || (ls[l] == "") {

		switch s {
		default:
			log.Errorf("Not localized %q", s)
		case "↑", "↓", "^|", "|v", "-", "Command":
		}

		return s
	}

	return ls[l]
}
