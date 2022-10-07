package controller

import (
	"html/template"
	"log"
	"net/http"
	"vsystem/config"
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
}

func NewData() *ViewData {
	return &ViewData{
		AppTitle: config.AppTitle,
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
