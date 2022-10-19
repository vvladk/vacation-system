package main

import (
	"fmt"
	"log"
	"net/http"
	"vsystem/config"
	"vsystem/pkg/controller/auth"
	"vsystem/pkg/controller/ccalendar"
	"vsystem/pkg/controller/cextraday"
	"vsystem/pkg/controller/crequest"
	"vsystem/pkg/controller/cuser"
	"vsystem/pkg/controller/cvacation"
	"vsystem/pkg/middleware"

	"github.com/julienschmidt/httprouter"
)

func main() {

	router := httprouter.New()

	serverString := fmt.Sprintf("%s:%d", config.Host, config.Port)
	//define a static path
	router.ServeFiles("/assets/*filepath", http.Dir("./ui/assets"))
	//routes for authorisation
	router.GET(`/login`, auth.Login)
	router.GET(`/logout`, auth.Loguot)
	router.GET(`/auth/google/callback`, auth.AuthCallBack)
	router.GET(`/AccessDenied`, auth.AccessDenied)

	// routes for users
	router.GET(`/users`, middleware.IsAuthorized(cuser.GetAll))
	router.GET(`/users/:uId`, middleware.IsAuthorized(cuser.CreateUpdate))
	router.POST(`/user/:uId`, middleware.IsAuthorized(cuser.Save)) // Create and Update DB
	router.GET(`/user/:uId`, middleware.IsAuthorized(cuser.ShowById))

	// routes for Extra days
	router.GET(`/extra_days`, middleware.IsAuthorized(cextraday.GetAll))
	router.GET(`/extra_day/:id`, middleware.IsAuthorized(cextraday.Create))
	router.POST(`/extra_day/:id`, middleware.IsAuthorized(cextraday.Save))

	// defaut path get current list of my vacation
	router.GET(`/`, middleware.Logging(middleware.IsAuthorized(cvacation.GetAllByUser)))
	router.GET(`/vacation/:vId`, middleware.Logging(middleware.IsAuthorized(cvacation.GetOneById)))       //Show edit form for update/create
	router.GET(`/vacations/:vId`, middleware.Logging(middleware.IsAuthorized(cvacation.PreviewVacation))) //Preview details
	router.POST(`/vacations/:vId`, middleware.Logging(middleware.IsAuthorized(cvacation.Create)))         //Cretae a vacation

	// routes for manage requests
	router.GET(`/requests`, middleware.Logging(middleware.IsAuthorized(crequest.GetAllByRole)))
	router.POST(`/requests/:vId`, middleware.Logging(middleware.IsAuthorized(crequest.UpadetByManager)))

	//routes for calendar view
	router.GET(`/calendar`, middleware.Logging(middleware.IsAuthorized(ccalendar.GetCalendar)))

	log.Println("Starting server on ", serverString)
	log.Fatal(http.ListenAndServe(serverString, router))

}
