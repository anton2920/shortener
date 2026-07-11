package main

import (
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/database/kv"
	"github.com/anton2920/gofa/encoding/wire"
	"github.com/anton2920/gofa/l10n"
	"github.com/anton2920/gofa/net/html"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/trace"
)

//gpp:generate encoding(wire), database(kv)
type User struct {
	database.RecordHeader

	Email    string
	Password string

	CreatedAt int64
}

/* TODO(anton2920): generate via gpp. */
func PutUserWire(s *wire.Serializer, user *User) {
	/* TODO(anton2920): implement. */
}

/* TODO(anton2920): generate via gpp. */
func CreateUser(tx *kv.Tx, user *User) error {
	var s wire.Serializer
	s.Buffer = make([]byte, 0, 1024)
	PutUserWire(&s, user)

	id, err := tx.Add(s.Buffer)
	if err != nil {
		return err
	}
	user.ID = id

	return nil
}

/* TODO(anton2920): generate via gpp. */
func GetUserByEmail(tx *kv.Tx, email string, user *User) error {
	/* Query email index for user. */
	/*

	 */
	return nil
}

/* TODO(anton2920): generate via gpp. */
func FillUserFromRequest(vs url.Values, user *User) {
	user.Email = vs.Get("Email")
	user.Password = vs.Get("Password")
}

/* TODO(anton2920): generate via gpp. */
func VerifyUser(l l10n.Language, user *User, repeatPassword string) error {
	/* TODO(anton2920): verify email. */

	if user.Password != repeatPassword {
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

func UserSigninHandler(w *http.Response, r *http.Request) error {
	var user User

	h := html.New(w, r, &Styles)

	switch r.Method {
	default:
		return UserSigninPage(&h, &user, nil)
	}
}

func UserSignoutHandler(w *http.Response, r *http.Request) error {
	return nil
}

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

func UserSignupHandler(w *http.Response, r *http.Request, db *kv.Database) error {
	var user User

	tx, err := db.Begin(r.Language)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	h := html.New(w, r, &Styles)

	switch r.Method {
	default:
		return UserSignupPage(&h, &user, nil)
	case http.MethodPost:
		FillUserFromRequest(r.Form, &user)
		if err := VerifyUser(r.Language, &user, r.Form.Get("RepeatPassword")); err != nil {
			return UserSignupPage(&h, &user, err)
		}

		if err := GetUserByEmail(tx, user.Email, nil); err != nil {
			return UserSignupPage(&h, &user, http.Conflict("%s", r.L("user with this email already exists")))
		}
		if err := CreateUser(tx, &user); err != nil {
			return UserSignupPage(&h, &user, err)
		}

		w.Redirect("/user/signin", http.StatusSeeOther)
		return nil
	}
}
