package controller

import (
	"html/template"
	"log"
	"net/http"
	"vsystem/config"
	"vsystem/pkg/controller/auth"

	"github.com/markbates/goth/gothic"
)

type Links struct {
	Link       string
	LinkActive string
	LinkTitle  string
}

type ViewData struct {
	AppTitle    string
	MenuItem    string
	SubMenuItem string
	Method      string
	Links       []Links
	UTitle      string
	UType       string
}

func NewData(r *http.Request) *ViewData {
	session, _ := gothic.Store.Get(r, auth.CookiesName)
	return &ViewData{
		AppTitle: config.AppTitle,
		UTitle:   session.Values["UserTitle"].(string),
		UType:    session.Values["UserType"].(string),
	}
}

func ExeTemlates(rw http.ResponseWriter, data any, tmpls []string) {

	ts, err := template.ParseFiles(tmpls...)
	if err != nil {
		log.Println(err.Error())
		http.Error(rw, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(rw, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(rw, "Internal Server Error", 500)
	}

}
