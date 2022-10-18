package crequest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"vsystem/config"
	"vsystem/pkg/controller"
	"vsystem/pkg/controller/auth"
	"vsystem/pkg/controller/cvacation"
	"vsystem/pkg/model/mtype"
	"vsystem/pkg/model/muser"
	"vsystem/pkg/model/mvacation"

	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth/gothic"
)

type HRView struct {
	EmployeeId     int
	VacationTypeId int
	Status         string
	StartDate      string
	EndDate        string
	StatusList     []string
}

func newHRView() *HRView {
	return &HRView{}
}

type ViewData struct {
	cvacation.ViewData
	UserList *muser.UserList
	HRView   HRView
}

func NewData(r *http.Request) *ViewData {
	return &ViewData{
		*cvacation.NewData(r),
		nil,
		*newHRView(),
	}
}

func GetAllByRole(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := NewData(r)
	session, _ := gothic.Store.Get(r, auth.CookiesName)
	if v, ok := session.Values["ID"]; ok || v != nil {

		data.User.GetById(v.(int))
		switch data.User.UserType {
		case config.UTypeFLM:
			data.VacationList.GetAllByFLM(v.(int))
		case config.UTypeEmp:
			http.Redirect(rw, r, "/", http.StatusFound)
			return
		default: //config.UTypeHR:
			data.HRView.EmployeeId, _ = strconv.Atoi(r.URL.Query().Get("Employee"))
			data.HRView.VacationTypeId, _ = strconv.Atoi(r.URL.Query().Get("VacationType"))
			data.HRView.Status = r.URL.Query().Get("Status")
			if data.HRView.Status == `` {
				data.HRView.Status = `New`
			}
			data.HRView.StartDate = r.URL.Query().Get("StartDate")
			data.HRView.EndDate = r.URL.Query().Get("EndDate")

			// User list
			data.UserList = muser.NewUserList()
			data.UserList.GetAll()
			data.UserList.List = append(data.UserList.List, muser.User{Id: 0, Title: "All"})
			//Vacation list
			data.VacationTypeList.GetVacationTypeList()
			data.VacationTypeList.List = append(data.VacationTypeList.List, mtype.VacationType{TypeId: 0, TypeTitle: "Any"})
			// Status List
			data.HRView.StatusList = []string{`Any`, `New`, `Approved`, `Rejected`}
			//Start Date
			if data.HRView.StartDate == `` {
				currentTime := time.Now()

				startDate := currentTime.AddDate(0, -1, 0)
				data.HRView.StartDate = fmt.Sprintf("%04d-%02d-%02d", startDate.Year(), startDate.Month(), startDate.Day())
			}
			if data.HRView.EndDate == `` {
				currentTime := time.Now()

				EndDate := currentTime.AddDate(0, 2, 0)
				data.HRView.EndDate = fmt.Sprintf("%04d-%02d-%02d", EndDate.Year(), EndDate.Month(), EndDate.Day())
			}
			//Create map with filter params

			data.VacationList.GetAllByHR(v.(int),
				data.HRView.EmployeeId,
				data.HRView.VacationTypeId,
				data.HRView.Status, data.HRView.StartDate,
				data.HRView.EndDate)
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
		if mng != `HR` {
			http.Redirect(rw, r, "/requests", http.StatusFound)
		} else {
			http.Redirect(rw, r, "/requests?Employee="+r.PostFormValue("Employee")+
				"&VacationType="+r.PostFormValue("VacationType")+
				"&Status="+r.PostFormValue("Status")+
				"&StartDate="+r.PostFormValue("StartDate")+
				"&EndDate="+r.PostFormValue("EndDate"), http.StatusFound)
		}
	}
	http.Redirect(rw, r, "/", http.StatusFound)

}
