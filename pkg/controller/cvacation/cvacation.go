package cvacation

import (
	"net/http"
	"strconv"
	"vsystem/config"
	"vsystem/pkg/controller"
	"vsystem/pkg/controller/auth"
	"vsystem/pkg/model"
	"vsystem/pkg/model/mtype"
	"vsystem/pkg/model/muser"
	"vsystem/pkg/model/mvacation"

	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth/gothic"
)

type ViewData struct {
	controller.ViewData
	mvacation.VacationList
	muser.User
	mvacation.Vacation
	mtype.VacationTypeList
}

func NewData(r *http.Request) *ViewData {
	return &ViewData{
		*controller.NewData(r),
		*mvacation.NewVacationList(),
		*muser.NewUser(),
		*mvacation.NewVacation(),
		*mtype.NewVacationTypeList(),
	}
}

func GetAllByUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := NewData(r)
	data.Links = append(data.Links, controller.Links{Link: `/`, LinkActive: `true`, LinkTitle: `Vacations`})
	data.MenuItem = "Vacations"

	session, _ := gothic.Store.Get(r, auth.CookiesName)
	if v, ok := session.Values["ID"]; ok || v != nil {
		data.VacationList.GetAllById(v.(int))
		data.User.GetById(v.(int))

		tmpls := []string{
			"./ui/tmpl/vacations/list.html",
			"./ui/tmpl/layout/links.html",
			"./ui/tmpl/layout/layout.html",
		}

		controller.ExeTemlates(rw, data, tmpls)
	} else {
		http.Redirect(rw, r, "/login", http.StatusFound)
	}

}
func GetOneById(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := NewData(r)
	data.Links = append(data.Links, controller.Links{Link: `/`, LinkActive: `true`, LinkTitle: `Vacations`})

	data.MenuItem = "Vacations"
	id, err := strconv.Atoi(p.ByName("vId"))
	data.Links = append(data.Links, controller.Links{Link: `/`, LinkActive: `false`, LinkTitle: `Vacation details`})
	model.CheckErr(err)
	data.Vacation.GetById(id)
	data.VacationTypeList.GetVacationTypeList()
	form := "./ui/tmpl/vacations/"
	if r.URL.Query().Get("partially") == `yes` {
		data.Partially = `yes`
		form += `form05.html`

	} else {
		form += `form.html`
	}

	tmpls := []string{
		form,
		"./ui/tmpl/vacations/tabs.html",
		"./ui/tmpl/layout/links.html",
		"./ui/tmpl/layout/layout.html",
	}
	controller.ExeTemlates(rw, data, tmpls)
}
func PreviewVacation(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := NewData(r)
	data.Links = append(data.Links, controller.Links{Link: `/`, LinkActive: `true`, LinkTitle: `Vacations`})
	id, err := strconv.Atoi(p.ByName("vId"))
	model.CheckErr(err)
	data.Links = append(data.Links, controller.Links{Link: `/`, LinkActive: ``, LinkTitle: `Vacation details`})
	data.MenuItem = "Vacations"

	session, _ := gothic.Store.Get(r, auth.CookiesName)
	if v, ok := session.Values["ID"]; ok || v != nil {
		data.VacationList.GetAllById(v.(int))
		data.User.GetById(v.(int))

		data.Vacation.Id = id
		data.Vacation.TypeId, _ = strconv.Atoi(r.URL.Query().Get("TypeId"))
		data.TypeTitle = data.VacationTypeList.GetById(data.Vacation.TypeId).TypeTitle
		data.Vacation.StartDate = r.URL.Query().Get("StartDate")
		data.Vacation.EndDate = r.URL.Query().Get("EndDate")

		if r.URL.Query().Get("partially") != `yes` {
			data.Duration = float64(model.GetDuration(data.Vacation.StartDate, data.Vacation.EndDate))
			data.Vacation.Partially = `no`
		} else {
			data.Duration = 0.5
			data.Vacation.EndDate = data.Vacation.StartDate
			data.Vacation.Part, _ = strconv.Atoi(r.URL.Query().Get("PartOfBd"))
			data.Vacation.Partially = `yes`
		}
		tmpls := []string{
			"./ui/tmpl/vacations/details.html",
			"./ui/tmpl/layout/links.html",
			"./ui/tmpl/layout/layout.html",
		}
		controller.ExeTemlates(rw, data, tmpls)
	} else {
		http.Redirect(rw, r, "/login", http.StatusFound)
	}

}

func Create(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := NewData(r)
	session, _ := gothic.Store.Get(r, auth.CookiesName)
	if v, ok := session.Values["ID"]; ok || v != nil {
		data.Vacation.TypeId, _ = strconv.Atoi(r.PostFormValue("TypeId"))
		data.Vacation.StartDate = r.PostFormValue("StartDate")
		data.Vacation.EndDate = r.PostFormValue("EndDate")
		data.Vacation.Part, _ = strconv.Atoi(r.PostFormValue("PartOfBd"))
		data.Vacation.Partially = r.PostFormValue("partially")
		if data.Vacation.Partially != `yes` {
			data.Duration = float64(model.GetDuration(data.Vacation.StartDate, data.Vacation.EndDate))
		} else {
			data.Duration = 0.5
		}
		var status int
		switch data.UType {
		case config.UTypeEmp:
			status = config.CreatedByUser
		case config.UTypeFLM:
			status = config.AcceptedByFLM
		case config.UTypeHR:
			status = config.AcceptedByHR
		}
		data.Vacation.Save2DB(v.(int), status)
	}
	http.Redirect(rw, r, "/", http.StatusFound)
}
