package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"vsystem/config"
	"vsystem/pkg/model"
	"vsystem/pkg/model/muser"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type ctxKey string

var provider string

var CookiesName string

func init() {
	CookiesName = config.Cfg.Section("session").Key("cookiesName").String()
	store := sessions.NewCookieStore([]byte(config.Cfg.Section("session").Key("key").String()))
	sec, _ := config.Cfg.Section("session").Key("maxAgeSec").Int()
	maxAge, _ := config.Cfg.Section("session").Key("maxAgeDays").Int()
	store.MaxAge(maxAge * sec)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure, _ = config.Cfg.Section("session").Key("isProd").Bool()

	gothic.Store = store
	provider = `google`
	goth.UseProviders(
		google.New(config.Cfg.Section("auth").Key("client_id").String(),
			config.Cfg.Section("auth").Key("client_secret").String(),
			config.Cfg.Section("auth").Key("callback").String(),
			"email", "profile"),
	)
}

func Login(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	gothic.BeginAuthHandler(rw, r)
}

func Loguot(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	session, _ := gothic.Store.Get(r, CookiesName)
	delete(session.Values, "ID")
	delete(session.Values, "email")
	delete(session.Values, "UserType")
	gothic.Store.Save(r, rw, session)
	gothic.Logout(rw, r)

	t, _ := template.ParseFiles("./ui/tmpl/layout/index.html")
	err := t.Execute(rw, nil)
	model.CheckErr(err)
}
func AuthCallBack(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	user, err := gothic.CompleteUserAuth(rw, r)
	if err != nil {
		fmt.Fprintln(rw, err)
		return
	}
	session, _ := gothic.Store.Get(r, CookiesName)
	u := muser.NewUser()
	u.GetByEmail(user.Email)
	if u.Id != 0 {
		session.Values["ID"] = u.Id
		session.Values["email"] = u.Email
		session.Values["UserType"] = u.UserType
		session.Values["UserTitle"] = u.Title
		gothic.Store.Save(r, rw, session)
		http.Redirect(rw, r, "/users/", http.StatusFound)
	} else {
		http.Redirect(rw, r, "/AccessDenied", http.StatusFound)
	}
}

func AccessDenied(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, err := template.ParseFiles("./ui/tmpl/layout/notregistered.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(rw, "Internal Server Error", 500)
		return
	}
	err = t.Execute(rw, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(rw, "Internal Server Error", 500)
	}
}
