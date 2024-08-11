package main

import (
	"net/mail"
	"unicode"
	"unicode/utf8"

	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/trace"
)

type User struct {
	ID    database.ID
	Flags int32

	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedOn int64

	URLs []database.ID
}

const (
	MinUserNameLen = 1
	MaxUserNameLen = 64

	MinEmailLen = 1
	MaxEmailLen = 128

	MinPasswordLen = 5
	MaxPasswordLen = 64
)

func UserNameValid(l Language, name string) error {
	defer trace.End(trace.Begin(""))

	if !strings.LengthInRange(name, MinUserNameLen, MaxUserNameLen) {
		return http.BadRequest(Ls(l, "length of the name must be between %d and %d characters"), MinUserNameLen, MaxUserNameLen)
	}

	/* Fist character must be a letter. */
	r, nbytes := utf8.DecodeRuneInString(name)
	if !unicode.IsLetter(r) {
		return http.BadRequest(Ls(l, "first character of the name must be a letter"))
	}

	/* Latter characters may include: letters, spaces, dots, hyphens and apostrophes. */
	for _, r := range name[nbytes:] {
		if (!unicode.IsLetter(r)) && (r != ' ') && (r != '.') && (r != '-') && (r != '\'') {
			return http.BadRequest(Ls(l, "second and latter characters of the name must be letters, spaces, dots, hyphens or apostrophes"))
		}
	}

	return nil
}

var TestUser = User{
	ID:        1,
	FirstName: "First",
	LastName:  "Last",
	Email:     "test@test.com",
	Password:  "testtest",
	CreatedOn: int64(time.Unix()),
}

func GetUserByEmail(email string, user *User) error {
	if email != TestUser.Email {
		return database.NotFound
	}

	*user = TestUser
	return nil
}

func GetUserByID(id database.ID, user *User) error {
	if id != TestUser.ID {
		return database.NotFound
	}

	*user = TestUser
	return nil
}

func CreateUser(user *User) error {
	return nil
}

func SaveUser(user *User) error {
	return nil
}

func DisplayUserTitle(w *http.Response, user *User) {
	w.WriteHTMLString(user.LastName)
	w.WriteString(` `)
	w.WriteHTMLString(user.FirstName)
	w.WriteString(` (ID: `)
	w.WriteID(user.ID)
	w.WriteString(`)`)
}

func UserPage(w *http.Response, r *http.Request) error {
	defer trace.End(trace.Begin(""))

	var user User

	id, err := GetIDFromURL(GL, r.URL, "/user/")
	if err != nil {
		return err
	}

	if err := GetUserByID(id, &user); err != nil {
		if err == database.NotFound {
			return http.NotFound(Ls(GL, "user with this ID does not exist"))
		}
		return http.ServerError(err)
	}

	DisplayHTMLStart(w)

	DisplayHeadStart(w)
	{
		w.WriteString(`<title>`)
		DisplayUserTitle(w, &user)
		w.WriteString(`</title>`)
	}
	DisplayHeadEnd(w)

	DisplayBodyStart(w)
	{
		w.WriteString(`<h2>`)
		DisplayUserTitle(w, &user)
		w.WriteString(`</h2>`)

		w.WriteString(`<p>`)
		w.WriteString(Ls(GL, "Email"))
		w.WriteString(`: `)
		w.WriteString(user.Email)
		w.WriteString(`</p>`)

		w.WriteString(`<p>`)
		w.WriteString(Ls(GL, "Created on"))
		w.WriteString(`: `)
		DisplayFormattedTime(w, user.CreatedOn)
		w.WriteString(`</p>`)

		w.WriteString(`<h3>`)
		w.WriteString(Ls(GL, "Shortened links"))
		w.WriteString(`</h3>`)
	}
	DisplayBodyEnd(w)

	DisplayHTMLEnd(w)
	return nil
}

func UserSigninPage(w *http.Response, r *http.Request, ierr error) error {
	defer trace.End(trace.Begin(""))

	DisplayHTMLStart(w)

	const title = "Sign in"

	if err := r.ParseForm(); err != nil {
		return http.ClientError(err)
	}

	DisplayHeadStart(w)
	{
		w.WriteString(`<title>`)
		w.WriteString(Ls(GL, title))
		w.WriteString(`</title>`)
	}
	DisplayHeadEnd(w)

	DisplayBodyStart(w)
	{
		w.WriteString(`<h2>`)
		w.WriteString(Ls(GL, title))
		w.WriteString(`</h2>`)

		DisplayError(w, GL, ierr)

		w.WriteString(`<form method="POST" action="` + APIPrefix + `/user/signin">`)
		{
			DisplayLabel(w, GL, "Email")
			DisplayConstraintInput(w, "email", MinEmailLen, MaxEmailLen, "Email", r.Form.Get("Email"), true)
			w.WriteString(`<br><br>`)

			DisplayLabel(w, GL, "Password")
			DisplayConstraintInput(w, "password", MinPasswordLen, MaxPasswordLen, "Password", "", true)
			w.WriteString(`<br><br>`)

			DisplaySubmit(w, GL, "", title)
		}
		w.WriteString(`</form>`)
	}
	DisplayBodyEnd(w)

	DisplayHTMLEnd(w)
	return nil
}

func UserSignupPage(w *http.Response, r *http.Request, ierr error) error {
	defer trace.End(trace.Begin(""))

	DisplayHTMLStart(w)

	const title = "Sign up"

	if err := r.ParseForm(); err != nil {
		return http.ClientError(err)
	}

	DisplayHeadStart(w)
	{
		w.WriteString(`<title>`)
		w.WriteString(Ls(GL, title))
		w.WriteString(`</title>`)
	}
	DisplayHeadEnd(w)

	DisplayBodyStart(w)
	{
		w.WriteString(`<h2>`)
		w.WriteString(Ls(GL, title))
		w.WriteString(`</h2>`)

		DisplayError(w, GL, ierr)

		w.WriteString(`<form method="POST" action="` + APIPrefix + `/user/signup">`)
		{
			DisplayLabel(w, GL, "First Name")
			DisplayConstraintInput(w, "text", MinUserNameLen, MaxUserNameLen, "FirstName", r.Form.Get("FirstName"), true)
			w.WriteString(`<br><br>`)

			DisplayLabel(w, GL, "Last Name")
			DisplayConstraintInput(w, "text", MinUserNameLen, MaxUserNameLen, "LastName", r.Form.Get("LastName"), true)
			w.WriteString(`<br><br>`)

			DisplayLabel(w, GL, "Email")
			DisplayConstraintInput(w, "email", MinEmailLen, MaxEmailLen, "Email", r.Form.Get("Email"), true)
			w.WriteString(`<br><br>`)

			DisplayLabel(w, GL, "Password")
			DisplayConstraintInput(w, "password", MinPasswordLen, MaxPasswordLen, "Password", "", true)
			w.WriteString(`<br><br>`)

			DisplayLabel(w, GL, "Repeat Password")
			DisplayConstraintInput(w, "password", MinPasswordLen, MaxPasswordLen, "RepeatPassword", "", true)
			w.WriteString(`<br><br>`)

			DisplaySubmit(w, GL, "", title)
		}
		w.WriteString(`</form>`)
	}
	DisplayBodyEnd(w)

	DisplayHTMLEnd(w)
	return nil
}

func UserSigninHandler(w *http.Response, r *http.Request) error {
	defer trace.End(trace.Begin(""))

	if err := r.ParseForm(); err != nil {
		return http.ClientError(err)
	}

	address, err := mail.ParseAddress(r.Form.Get("Email"))
	if err != nil {
		return UserSigninPage(w, r, http.BadRequest(Ls(GL, "provided email is not valid")))
	}
	email := address.Address

	var user User
	if err := GetUserByEmail(email, &user); err != nil {
		if err == database.NotFound {
			return UserSigninPage(w, r, http.NotFound(Ls(GL, "user with this email does not exist")))
		}
		return http.ServerError(err)
	}

	password := r.Form.Get("Password")
	if user.Password != password {
		return UserSigninPage(w, r, http.Conflict(Ls(GL, "provided password is incorrect")))
	}

	token, err := GenerateSessionToken()
	if err != nil {
		return http.ServerError(err)
	}
	expiry := time.Unix() + OneWeek

	session := &Session{
		ID:     user.ID,
		Expiry: expiry,
	}

	SessionsLock.Lock()
	Sessions[token] = session
	SessionsLock.Unlock()

	if Debug {
		w.SetCookieUnsafe("Token", token, expiry)
	} else {
		w.SetCookie("Token", token, expiry)
	}
	w.Redirect("/", http.StatusSeeOther)
	return nil
}

func UserSignoutHandler(w *http.Response, r *http.Request) error {
	defer trace.End(trace.Begin(""))

	token := r.Cookie("Token")
	if token == "" {
		return http.UnauthorizedError
	}

	if _, err := GetSessionFromToken(token); err != nil {
		return http.UnauthorizedError
	}

	SessionsLock.Lock()
	delete(Sessions, token)
	SessionsLock.Unlock()

	w.DelCookie("Token")
	w.Redirect("/", http.StatusSeeOther)
	return nil
}

func UserSignupHandler(w *http.Response, r *http.Request) error {
	defer trace.End(trace.Begin(""))

	if err := r.ParseForm(); err != nil {
		return http.ClientError(err)
	}

	firstName := r.Form.Get("FirstName")
	if err := UserNameValid(GL, firstName); err != nil {
		return UserSignupPage(w, r, err)
	}

	lastName := r.Form.Get("LastName")
	if err := UserNameValid(GL, lastName); err != nil {
		return UserSignupPage(w, r, err)
	}

	address, err := mail.ParseAddress(r.Form.Get("Email"))
	if err != nil {
		return UserSignupPage(w, r, http.BadRequest(Ls(GL, "provided email is not valid")))
	}
	email := address.Address

	password := r.Form.Get("Password")
	repeatPassword := r.Form.Get("RepeatPassword")
	if !strings.LengthInRange(password, MinPasswordLen, MaxPasswordLen) {
		return UserSignupPage(w, r, http.BadRequest(Ls(GL, "password length must be between %d and %d characters long"), MinPasswordLen, MaxPasswordLen))
	}
	if password != repeatPassword {
		return UserSignupPage(w, r, http.BadRequest(Ls(GL, "passwords do not match each other")))
	}

	var user User
	if err := GetUserByEmail(email, &user); err == nil {
		return UserSignupPage(w, r, http.Conflict(Ls(GL, "user with this email already exists")))
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Email = email
	user.Password = password
	user.CreatedOn = int64(time.Unix())

	if err := CreateUser(&user); err != nil {
		return http.ServerError(err)
	}

	w.Redirect("/", http.StatusSeeOther)
	return nil

}
