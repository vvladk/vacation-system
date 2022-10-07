package main

import (
	"fmt"
	"log"
	"net/http"
	"vsystem/config"
	"vsystem/pkg/controller/auth"
	"vsystem/pkg/controller/cextraday"
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
	// router.GET(`/`, middleware.IsAuthorized(cuser.GetAll)) // TMP should be replaced!
	router.GET(`/users/:uId`, middleware.IsAuthorized(cuser.CreateUpdate))
	router.POST(`/user/:uId`, middleware.IsAuthorized(cuser.Save)) // Create and Update DB
	// router.POST(`/users/:uId`, middleware.IsAuthorized(cuser.Delete)) // Delete from  DB TBD
	router.GET(`/user/:uId`, middleware.IsAuthorized(cuser.ShowById))

	// routes for Extra days
	router.GET(`/extra_days`, middleware.IsAuthorized(cextraday.GetAll))
	router.GET(`/extra_day/:id`, middleware.IsAuthorized(cextraday.Create))
	router.POST(`/extra_day/:id`, middleware.IsAuthorized(cextraday.Save))

	// defaut path  get current list of my vacation
	router.GET(`/`, cvacation.GetAllByUser)
	// router.GET(`/vacation/:vId`, cvacation.GetOneById)               //Show edit form for update/create
	// router.POST(`/preview_vacation/:vId`, cvacation.PreviewVacation) //Preview details

	log.Println("Starting server on ", serverString)
	log.Fatal(http.ListenAndServe(serverString, router))

}
