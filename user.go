package main

import (
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/l10n"
	"github.com/anton2920/gofa/net/html"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/trace"
)

//gpp:generate encoding(wire)
type User struct {
	database.RecordHeader

	Email          string
	Password       string
	RepeatPassword string

	CreatedAt int64
}

func (db *Database) CreateUser(user *User) error {
	user.CreatedAt = time.Now()

	if _, err := DB.Exec(`INSERT INTO users(email, password, created_at) VALUES(?, ?, ?)`, user.Email, user.Password, user.CreatedAt); err != nil {
		if UniqueViolation(err) {
			return http.Conflict("user with this email already exists")
		}
		return http.ServerError(err)
	}

	return nil
}

func FillUserFromValues(vs url.Values, user *User) {
	user.Email = vs.Get("Email")
	user.Password = vs.Get("Password")
	user.RepeatPassword = vs.Get("RepeatPassword")
}

func VerifyUser(l l10n.Language, user *User) error {
	/* TODO(anton2920): verify email. */

	if user.Password != user.RepeatPassword {
		return http.BadRequest(l.L("passwords do not match"))
	}

	return nil
}

func UserSigninPage(h *html.HTML, user *User, ierr error) error {
	h.Begin()

	h.HeadBegin()
	{
		DisplayStyles(h)
		h.Title("Sign in | Shortener")
	}
	h.HeadEnd()

	h.BodyBegin()
	{
		h.H2("Sign in", html.Class("mb-3"))

		h.Error(ierr, html.Class("mb-3"))

		h.FormBegin("POST")
		{
			h.Label("Email")
			h.Input("email", html.Attributes{Class: "mb-3", Name: "Email", Value: user.Email, Required: true})

			h.Label("Password")
			h.Input("password", html.Attributes{Class: "mb-3", Name: "Password", Required: true})

			h.Button("Sign in", StyleButtonSubmit)
		}
		h.FormEnd()
	}
	h.BodyEnd()

	h.End()
	return nil
}

func UserSignin(h *html.HTML, r *http.Request, db *Database) error {
	var user User

	switch r.Method {
	default:
		return UserSigninPage(h, &user, nil)
	}
}

func UserSignout() {}

func UserSignupPage(h *html.HTML, user *User, ierr error) error {
	defer trace.End(trace.Begin(""))

	h.Begin()

	h.HeadBegin()
	{
		DisplayStyles(h)
		h.Title("Sign up | Shortener")
	}
	h.HeadEnd()

	h.BodyBegin()
	{
		h.H2("Sign up", html.Class("mb-3"))

		h.Error(ierr, html.Class("mb-3"))

		h.FormBegin("POST")
		{
			h.Label("Email")
			h.Input("email", html.Attributes{Class: "mb-3", Name: "Email", Value: user.Email, Required: true})

			h.Label("Password")
			h.Input("password", html.Attributes{Class: "mb-3", Name: "Password", Required: true})

			h.Label("Repeat password")
			h.Input("password", html.Attributes{Class: "mb-3", Name: "RepeatPassword", Required: true})

			h.Button("Sign up", StyleButtonSubmit)
		}
		h.FormEnd()
	}
	h.BodyEnd()

	h.End()
	return nil
}

func UserSignup(h *html.HTML, r *http.Request, db *Database) error {
	defer trace.End(trace.Begin(""))

	var user User

	switch r.Method {
	default:
		return UserSignupPage(h, &user, nil)
	case http.MethodPost:
		FillUserFromValues(r.Form, &user)
		if err := VerifyUser(h.Language, &user); err != nil {
			return UserSignupPage(h, &user, err)
		}

		if err := db.CreateUser(&user); err != nil {
			return UserSignupPage(h, &user, err)
		}

		h.Redirect("/user/signin", http.StatusSeeOther)
		return nil
	}
}
