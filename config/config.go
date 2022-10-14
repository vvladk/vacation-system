package config

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

const (
	// Users parameters
	UTypeHR  = `HR`
	UTypeFLM = `FLM`
	UTypeEmp = `Employee`

	// Vacations parameters
	AvailableImmediately = 1 // All number of days is available from first day
	AvailableGradually   = 2 // Every month added 1/12 part of annual amount
	UnPaidVacationTitle  = `Unpaid Vacation`
	PaidVacationTitle    = `Paid vacation`
	SickLeaveTitle       = `Sick Leave`
	UnPaidVacationId     = 2
	PaidVacationId       = 1
	SickLeaveId          = 3

	//vacations stsatuses
	CreatedByUser      = 1
	CreatedByUserTitle = `Created`
	AcceptedByFLM      = 2
	AcceptedByFLMTitle = `Accepted by FLM`
	RejectedByFLM      = 3
	RejectedByFLMTitle = `Rejected by FLM`
	AcceptedByHR       = 4
	AcceptedByHRTitle  = `Accepted by HR`
	RejectedByHR       = 5
	RejectedByHRTitle  = `Rejected by HR`
)

var Statuses = map[int]string{
	CreatedByUser: CreatedByUserTitle,
	AcceptedByFLM: AcceptedByFLMTitle,
	RejectedByFLM: RejectedByFLMTitle,
	AcceptedByHR:  AcceptedByHRTitle,
	RejectedByHR:  RejectedByHRTitle,
}

var Cfg *ini.File

type sqlConf struct {
	Path   string
	DbName string
}

var SqlConf sqlConf

var Port int
var Host string

var AppTitle string

var UserTypes = []string{UTypeHR, UTypeFLM, UTypeEmp}

func init() {

	cfg, err := ini.Load("./ini/vs.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	Cfg = cfg
	AppTitle = cfg.Section("").Key("AppTitle").String()

	SqlConf.Path = cfg.Section("DB").Key("path").String()
	SqlConf.DbName = cfg.Section("DB").Key("file").String()

	Port = cfg.Section("server").Key("port").MustInt(4000)
	Host = cfg.Section("server").Key("host").MustString("localhost")

	//Type vacation
	for _, v := range cfg.ChildSections("VacationType") {
		fmt.Printf("\n|%+v|\n", v.Key("VacationTitle"))
	}
}
