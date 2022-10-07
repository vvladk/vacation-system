package cvacation

import (
	"net/http"
	"vsystem/pkg/controller"
	"vsystem/pkg/controller/auth"
	"vsystem/pkg/model/mvacation"

	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth/gothic"
)

type ViewData struct {
	controller.ViewData
	mvacation.VacationList
	// vacation.Vacation
	// mtype.VacationTypeList
}

func NewData() *ViewData {
	return &ViewData{
		*controller.NewData(),
		*mvacation.NewVacationList(),
		// *vacation.NewVacation(),
		// *mtype.NewVacationTypeList(),
	}
}

func GetAllByUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := NewData()
	data.Links = append(data.Links, controller.Links{Link: `/`, LinkActive: `true`, LinkTitle: `Vacations`})
	data.MenuItem = "Vacations"

	session, _ := gothic.Store.Get(r, auth.CookiesName)
	if v, ok := session.Values["ID"]; ok || v != nil {
		data.VacationList.GetAllById(v.(int))
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
	// data := NewData()
	// data.Links = append(data.Links, controller.Links{Link: `/`, LinkActive: `true`, LinkTitle: `Vacations`})
	// data.MenuItem = "Vacations"

	// data.VacationTypeList.GetVacationTypeList()

	// tmpls := []string{
	// 	"./ui/tmpl/vacations/form.tmpl.html",
	// 	"./ui/tmpl/base/links.tmpl.html",
	// 	"./ui/tmpl/base/layout.tmpl.html",
	// }
	// controller.ExeTemlates(rw, data, tmpls)
}
func PreviewVacation(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// data := NewData()
	// data.Links = append(data.Links, controller.Links{Link: `/`, LinkActive: `true`, LinkTitle: `Vacations`})
	// // Id, err := strconv.Atoi(p.ByName("vId"))
	// // common.CheckErr(err)
	// data.Links = append(data.Links, controller.Links{Link: `/preview_vacation/` + p.ByName("vId"), LinkActive: `true`, LinkTitle: `Vacation details`})
	// data.MenuItem = "Vacations"

	// data.Id, _ = strconv.Atoi(p.ByName("vId"))
	// data.TypeId, _ = strconv.Atoi(r.PostFormValue("TypeId"))
	// // data.VacationTypeList.GetById(data.TypeId)
	// // data.TypeTitle = data.VacationTypeList.VacationTypeList[0].TypeTitle
	// data.UserId = 1
	// data.UserEmail = `vladislav.kondratyuk@gmail.com`
	// data.StartDate = r.PostFormValue("StartDate")
	// data.EndDate = r.PostFormValue("EndDate")
	// data.Duration = model.GetDuration(data.StartDate, data.EndDate)

	// tmpls := []string{
	// 	"./ui/tmpl/vacations/preview.tmpl.html",
	// 	"./ui/tmpl/base/links.tmpl.html",
	// 	"./ui/tmpl/base/layout.tmpl.html",
	// }
	// controller.ExeTemlates(rw, data, tmpls)

}
