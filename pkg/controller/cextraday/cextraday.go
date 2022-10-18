package cextraday

import (
	"net/http"
	"vsystem/config"
	"vsystem/pkg/controller"
	"vsystem/pkg/model/mextraday"

	"github.com/julienschmidt/httprouter"
)

type ViewData struct {
	controller.ViewData
	mextraday.ExtraDayList
}

func NewData(r *http.Request) *ViewData {
	return &ViewData{
		*controller.NewData(r),
		*mextraday.NewExtraDayList(),
	}
}

func GetAll(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := NewData(r)
	data.Links = append(data.Links, controller.Links{Link: `/extra_days`, LinkActive: `true`, LinkTitle: `Holidays`})
	data.MenuItem = "Holidays"

	data.ExtraDayList.GetAll()
	templ := "./ui/tmpl/extra_days/"
	if data.UType != config.UTypeHR {
		templ += "list_ro.html"
	} else {
		templ += "list.html"
	}

	tmpls := []string{
		templ,
		"./ui/tmpl/layout/links.html",
		"./ui/tmpl/layout/layout.html",
	}

	controller.ExeTemlates(rw, data, tmpls)
}

func Create(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := NewData(r)
	data.Links = append(data.Links, controller.Links{Link: `/extra_days`, LinkActive: `true`, LinkTitle: `Holidays`})
	data.MenuItem = "Holidays"

	data.Links = append(data.Links, controller.Links{Link: `/extra_day/0`, LinkActive: `true`, LinkTitle: `New Holidays`})
	tmpls := []string{
		"./ui/tmpl/extra_days/form.html",
		"./ui/tmpl/layout/links.html",
		"./ui/tmpl/layout/layout.html",
	}
	controller.ExeTemlates(rw, data, tmpls)
}

func Save(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	mextraday.Save2DB(p.ByName("id"), r.PostFormValue("HDate"), r.PostFormValue("description"))
	http.Redirect(rw, r, `/extra_days`, http.StatusFound)
}
