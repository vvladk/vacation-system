package middleware

import (
	"log"
	"net/http"
	"time"
	"vsystem/pkg/controller/auth"

	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth/gothic"
)

func Logging(router httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		start := time.Now()
		router(rw, r, p)
		log.Printf("[%s]-[%s]-[%s]", r.Method, r.RequestURI, time.Since(start))
	}
}

func IsAuthorized(router httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		session, _ := gothic.Store.Get(r, auth.CookiesName)
		if v, ok := session.Values["ID"]; ok || v != nil {
			router(rw, r, p)
		} else {
			http.Redirect(rw, r, "/login", http.StatusFound)
		}
	}
}
