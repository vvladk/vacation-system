package cuser

import (
	"fmt"
	"net/http"
	"strconv"
	"vsystem/config"
	"vsystem/pkg/controller"
	"vsystem/pkg/model"
	"vsystem/pkg/model/muser"

	"github.com/julienschmidt/httprouter"
)

type ViewData struct {
	controller.ViewData
	muser.UserList
	muser.User
	EmployeeTypes []string
}

func NewData() *ViewData {
	return &ViewData{
		*controller.NewData(),
		*muser.NewUserList(),
		*muser.NewUser(),
		config.UserTypes,
	}
}

// Return all existing users
func GetAll(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := NewData()
	data.Links = append(data.Links, controller.Links{Link: `/users`, LinkActive: `true`, LinkTitle: `Employees`})
	data.MenuItem = "Users"

	data.UserList.GetAll()

	tmpls := []string{
		"./ui/tmpl/user/list.html",
		"./ui/tmpl/layout/links.html",
		"./ui/tmpl/layout/layout.html",
	}

	controller.ExeTemlates(rw, data, tmpls)

}

// Display form for creation or edit user
func CreateUpdate(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := NewData()
	data.Links = append(data.Links, controller.Links{Link: `/users`, LinkActive: `true`, LinkTitle: `Employees`})
	data.MenuItem = "Users"
	data.Method = "POST"
	id, err := strconv.Atoi(p.ByName("uId"))
	model.CheckErr(err)

	data.User.GetById(id)
	if id == 0 {
		data.Links = append(data.Links, controller.Links{Link: `/user/` + p.ByName("uId"), LinkActive: `true`, LinkTitle: `New Employee`})
	} else {
		data.Links = append(data.Links, controller.Links{Link: `/user/` + p.ByName("uId"), LinkActive: `true`, LinkTitle: `Edit Employee [` + data.User.Title + `]`})
	}
	tmpls := []string{
		"./ui/tmpl/user/form.html",
		"./ui/tmpl/layout/links.html",
		"./ui/tmpl/layout/layout.html",
	}

	controller.ExeTemlates(rw, data, tmpls)
}

// resolve 2 post requests Create - Update
func Save(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.Atoi(p.ByName("uId"))
	model.CheckErr(err)
	u := muser.NewUser()
	u.Title = r.PostFormValue("title")
	u.Email = r.PostFormValue("email")
	u.StartDate = r.PostFormValue("StartDate")
	u.SpillOver, _ = strconv.ParseFloat(r.PostFormValue("SpillOver"), 64)
	u.ExtraDays, _ = strconv.ParseFloat(r.PostFormValue("ExtraDays"), 64)
	u.IsActive, _ = strconv.ParseBool(r.PostFormValue("IsActive"))
	u.UserType = r.PostFormValue("UserType")
	u.FLMId, _ = strconv.Atoi(r.PostFormValue("FLM"))

	for _, v := range r.Form["vacationType"] {
		i, _ := strconv.Atoi(v)
		u.VacancyBalance = append(u.VacancyBalance, muser.VacancyBalance{Id: i, IsAvailable: true})
	}

	u.Save2DB(id)

	http.Redirect(rw, r, fmt.Sprintf("/user/%d", u.Id), http.StatusFound)
}

// Show full information about employee
func ShowById(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := NewData()
	data.Links = append(data.Links, controller.Links{Link: `/users`, LinkActive: `true`, LinkTitle: `Employees`})
	data.MenuItem = "Users"
	id, err := strconv.Atoi(p.ByName("uId"))
	model.CheckErr(err)

	data.User.GetById(id)
	data.Links = append(data.Links, controller.Links{Link: `/user/` + p.ByName("uId"), LinkActive: `true`, LinkTitle: `Details of Employee [` + data.User.Title + `]`})
	tmpls := []string{
		"./ui/tmpl/user/details.html",
		"./ui/tmpl/layout/links.html",
		"./ui/tmpl/layout/layout.html",
	}

	controller.ExeTemlates(rw, data, tmpls)
}
