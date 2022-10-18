package mvacation

import (
	"database/sql"
	"fmt"
	"log"
	"vsystem/config"
	"vsystem/pkg/model"
	"vsystem/pkg/model/mtype"
	"vsystem/pkg/model/muser"
)

type Vacation struct {
	Nn          int
	Id          int
	TypeId      int
	TypeTitle   string
	UserId      int
	UserTitle   string
	UserEmail   string
	StartDate   string
	EndDate     string
	Duration    float64
	Status      int
	StatusTitle string
	Spent       float64
	AsOfNow     string
	AsOfDec31   string
	Partially   string
	Part        int //1 first part of BD, 2 - second part of BD
}

func NewVacation() *Vacation {
	return &Vacation{}
}

type VacationList struct {
	List []Vacation
}

func NewVacationList() *VacationList {
	return &VacationList{}
}

type CalendarView struct {
	Title    string `json:"title"`
	Start    string `json:"start"`
	End      string `json:"end"`
	duration float64
	part     int
}

type CalendarViewList struct {
	Events []CalendarView
}

// return list of vacation filtered by user userId
// for current and previos year
// sorted by Start vacation DESC
func (l *VacationList) GetAllById(uId int) {
	db := model.GetDB()
	defer db.Close()
	sqlS := `SELECT
					v.id,
					u.email,
					v.typeId,
					v.startDate,
					v.duration,
					v.status,
					v.partOfBd
			FROM vacations AS v
			LEFT JOIN users AS u ON v.userId = u.id
			WHERE u.id = ?
			AND  strftime('%Y',v.startDate) > (strftime('%Y','now') - 1)
			ORDER BY v.startDate DESC;`
	rows, err := db.Query(sqlS, uId)
	model.CheckErr(err)
	defer rows.Close()

	nn := 1
	for rows.Next() {
		v := NewVacation()
		err := rows.Scan(&v.Id, &v.UserEmail, &v.TypeId, &v.StartDate, &v.Duration, &v.Status, &v.Part)
		if v.Part > 0 {
			v.Partially = `yes`
		}
		model.CheckErr(err)
		v.Nn = nn
		nn++
		vType := mtype.NewVacationTypeList()
		v.TypeTitle = vType.GetById(v.TypeId).TypeTitle
		v.EndDate = model.GetEndDate(v.StartDate, v.Duration)
		v.StartDate = model.ReformatDate(v.StartDate)
		v.StatusTitle = config.Statuses[v.Status]
		l.List = append(l.List, *v)
	}

}

func (v *Vacation) GetById(vId int) {
	if vId > 0 {
		db := model.GetDB()
		defer db.Close()
		sqlS := `SELECT v.id,v.typeId, v.startDate, v.duration, status FROM vacations AS v WHERE v.id = ?`
		err := db.QueryRow(sqlS, vId).Scan(&v.Id, &v.TypeId, &v.StartDate, &v.Duration, &v.Status)
		if err != nil {
			log.Println(err)
		}
	}
}

func (v *Vacation) Save2DB(uId int, status int) {
	db := model.GetDB()
	defer db.Close()
	sqlI := `INSERT INTO vacations(userId, typeId, startDate, duration, status, partOfBd) VALUES(?,?,?,?,?,?)`
	_, err := db.Exec(sqlI, uId, v.TypeId, v.StartDate, v.Duration, status, v.Part)
	model.CheckErr(err)
}

func (l *VacationList) GetAllByFLM(uId int) {
	sqlS := `SELECT
					v.id,
					v.userId,
					u.title,
					v.typeId,
					v.startDate,
					v.duration,
					v.status,
					v.partOfBd
			FROM vacations AS v
			LEFT JOIN users AS u ON v.userId = u.id
			WHERE u.flm = ?
			AND v.status = 1
			ORDER BY v.startDate DESC;`
	l.getAllByMng(uId, sqlS, `FLM`, 0, 0, 0, ``, ``)
}

func (l *VacationList) GetAllByHR(uId int, EmployeeId int, VacationTypeId int, Status string, StartDate string, EndDate string) {
	var state int
	switch Status {
	case `Any`:
		state = 0
	case `New`:
		state = config.AcceptedByFLM
	case `Approved`:
		state = config.AcceptedByHR
	case `Rejected`:
		state = config.RejectedByHR
	}

	sqlS := `SELECT
					v.id,
					v.userId,
					u.title,
					v.typeId,
					v.startDate,
					v.duration,
					v.status,
					v.partOfBd
			FROM vacations AS v
			LEFT JOIN users AS u ON v.userId = u.id
			WHERE v.status > 1 
			AND IIF(?, u.id = ?, 1 )
			AND IIF(?, v.typeId = ?, 1 )
			AND IIF(?, v.status = ?, 1 )
			AND v.StartDate > ?
			AND v.StartDate < ?
			ORDER BY v.startDate DESC;`
	l.getAllByMng(uId, sqlS, `HR`, EmployeeId, VacationTypeId, state, StartDate, EndDate)
}

func (l *VacationList) getAllByMng(uId int, sqlS string, mngType string, EmployeeId int, VacationTypeId int, state int, StartDate string, EndDate string) {
	db := model.GetDB()
	defer db.Close()

	var rows *sql.Rows
	var err error
	if mngType == `FLM` {
		rows, err = db.Query(sqlS, uId)
	} else {
		rows, err = db.Query(sqlS, EmployeeId, EmployeeId, VacationTypeId, VacationTypeId, state, state, StartDate, EndDate)

	}
	model.CheckErr(err)
	defer rows.Close()
	usermap := make(map[int][]muser.VacancyBalance)
	nn := 1
	for rows.Next() {
		v := NewVacation()
		err := rows.Scan(&v.Id, &v.UserId, &v.UserTitle, &v.TypeId, &v.StartDate, &v.Duration, &v.Status, &v.Part)
		model.CheckErr(err)
		if v.Part > 0 {
			v.Partially = `yes`
		}
		v.Nn = nn
		nn++
		vType := mtype.NewVacationTypeList()
		v.TypeTitle = vType.GetById(v.TypeId).TypeTitle
		v.EndDate = model.GetEndDate(v.StartDate, v.Duration)
		v.StartDate = model.ReformatDate(v.StartDate)
		v.StatusTitle = config.Statuses[v.Status]
		// Get info by user if map empty
		u := muser.NewUser()
		if _, ok := usermap[v.UserId]; !ok {
			u.GetById(v.UserId)
			u.GetVacationBalabce(v.UserId, u.StartDate, u.ExtraDays, u.SpillOver)
			usermap[v.UserId] = u.VacancyBalance
		}
		for _, w := range usermap[v.UserId] {
			if w.Id == v.TypeId {
				v.Spent = w.Spent
				if w.IsUnLim {
					v.AsOfNow = `-`
					v.AsOfDec31 = `-`
				} else {
					v.AsOfNow = fmt.Sprintf("%.2f", w.AsOfNow)
					v.AsOfDec31 = fmt.Sprintf("%.2f", w.AsOfDec31)
				}
			}
		}
		l.List = append(l.List, *v)
	}

}

func UpdateStatus(vId int, mng string, responce string) {
	db := model.GetDB()
	defer db.Close()
	sqlU := `UPDATE vacations SET status = ? WHERE id = ?`
	var tmpId int
	switch {
	case responce == `yes` && mng == config.UTypeFLM:
		tmpId = config.AcceptedByFLM
	case responce == `no` && mng == config.UTypeFLM:
		tmpId = config.RejectedByFLM
	case responce == `yes` && mng == config.UTypeHR:
		tmpId = config.AcceptedByHR
	case responce == `no` && mng == config.UTypeHR:
		tmpId = config.RejectedByHR
	}
	_, err := db.Exec(sqlU, tmpId, vId)
	model.CheckErr(err)
}

func (l *CalendarViewList) GetFor3Month() {
	db := model.GetDB()
	defer db.Close()
	sqlS := `SELECT
					u.title,
					v.startDate,
					v.duration,
					v.partOfBd
			FROM vacations AS v
			LEFT JOIN users AS u ON v.userId = u.id
			WHERE v.status = ?
			AND v.startDate BETWEEN DATE('now', '-1 months') AND DATE('now', '+1 months')
			ORDER BY v.startDate;`
	rows, err := db.Query(sqlS, config.AcceptedByHR)
	model.CheckErr(err)
	defer rows.Close()

	for rows.Next() {
		cv := CalendarView{}
		err := rows.Scan(&cv.Title, &cv.Start, &cv.duration, &cv.part)
		model.CheckErr(err)
		switch cv.part {
		case 1:
			cv.End = cv.Start
			cv.Start += `T10:00:00`
			cv.End += `T14:00:00`
		case 2:
			cv.End = cv.Start
			cv.Start += `T15:00:00`
			cv.End += `T19:00:00`
		case 0:
			cv.End = model.GetEndDate4Cal(cv.Start, cv.duration)
		}
		l.Events = append(l.Events, cv)
	}
}
