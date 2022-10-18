package ccalendar

import (
	"fmt"
	"net/http"
	"time"
	"vsystem/pkg/controller"
	"vsystem/pkg/controller/auth"
	"vsystem/pkg/controller/cvacation"
	"vsystem/pkg/model/mvacation"

	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth/gothic"
)

type CalendarData struct {
	CurrenDate string
	Events     mvacation.CalendarViewList
}

func NewCalendarData() *CalendarData {
	return &CalendarData{}
}

type ViewData struct {
	cvacation.ViewData
	Data *CalendarData
}

func NewData(r *http.Request) *ViewData {
	return &ViewData{
		*cvacation.NewData(r),
		NewCalendarData(),
	}
}

func GetCalendar(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	data := NewData(r)

	session, _ := gothic.Store.Get(r, auth.CookiesName)
	if v, ok := session.Values["ID"]; ok || v != nil {
		// data.User.GetById(v.(int))
		data.UTitle = session.Values["UserTitle"].(string)
		data.UType = session.Values["UserType"].(string)
		data.Data.CurrenDate = fmt.Sprintf("%04d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day())
		cv := mvacation.CalendarViewList{}
		cv.GetFor3Month()
		data.Data.Events = cv

		tmpls := []string{
			"./ui/tmpl/calendar/calendar.html",
			"./ui/tmpl/layout/links.html",
			"./ui/tmpl/layout/layout.html",
		}
		controller.ExeTemlates(rw, data, tmpls)
	}
}
