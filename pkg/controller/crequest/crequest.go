package crequest

import (
	"net/http"
	"strconv"
	"vsystem/config"
	"vsystem/pkg/controller"
	"vsystem/pkg/controller/auth"
	"vsystem/pkg/controller/cvacation"
	"vsystem/pkg/model/mvacation"

	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth/gothic"
)

func GetAllByRole(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := cvacation.NewData(r)
	session, _ := gothic.Store.Get(r, auth.CookiesName)
	if v, ok := session.Values["ID"]; ok || v != nil {

		data.User.GetById(v.(int))
		switch data.User.UserType {
		case config.UTypeFLM:
			data.VacationList.GetAllByFLM(v.(int))
		case config.UTypeHR:
			data.VacationList.GetAllByHR(v.(int))
		case config.UTypeEmp:
			http.Redirect(rw, r, "/", http.StatusFound)
			return
		default:
			data.VacationList.GetAllByHR(v.(int))
		}
		data.UTitle = session.Values["UserTitle"].(string)
		data.UType = session.Values["UserType"].(string)
		tmpls := []string{
			"./ui/tmpl/requests/list.html",
			"./ui/tmpl/layout/links.html",
			"./ui/tmpl/layout/layout.html",
		}

		controller.ExeTemlates(rw, data, tmpls)

	}
}
func UpadetByManager(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	session, _ := gothic.Store.Get(r, auth.CookiesName)
	if v, ok := session.Values["ID"]; ok || v != nil {
		vId, _ := strconv.Atoi(p.ByName("vId"))
		mng := session.Values["UserType"].(string)
		response := r.PostFormValue("response")

		mvacation.UpdateStatus(vId, mng, response)
		http.Redirect(rw, r, "/requests", http.StatusFound)

	}
	http.Redirect(rw, r, "/", http.StatusFound)

}
