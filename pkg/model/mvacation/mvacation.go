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
	Nn                 int
	Id                 int
	TypeId             int
	TypeTitle          string
	UserId             int
	UserTitle          string
	UserEmail          string
	StartDate          string
	EndDate            string
	Duration           int
	Status             int
	StatusTitle        string
	Spent              float64
	CurrentlyAvailable string
	AsOfDec31          string
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

// return list of vacation filtered by user userId
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
					v.status
			FROM vacations AS v
			LEFT JOIN users AS u ON v.userId = u.id
			WHERE u.id = ?
			ORDER BY v.startDate DESC;`
	rows, err := db.Query(sqlS, uId)
	model.CheckErr(err)
	defer rows.Close()

	nn := 1
	for rows.Next() {
		v := NewVacation()
		err := rows.Scan(&v.Id, &v.UserEmail, &v.TypeId, &v.StartDate, &v.Duration, &v.Status)
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

func (v *Vacation) Save2DB(uId int) {
	db := model.GetDB()
	defer db.Close()
	sqlI := `INSERT INTO vacations(userId, typeId, startDate, duration, status) VALUES(?,?,?,?,?)`
	_, err := db.Exec(sqlI, uId, v.TypeId, v.StartDate, v.Duration, config.CreatedByUser)
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
					v.status
			FROM vacations AS v
			LEFT JOIN users AS u ON v.userId = u.id
			WHERE u.flm = ?
			AND v.status = 1
			ORDER BY v.startDate DESC;`
	l.getAllByMng(uId, sqlS, `FLM`)
}

func (l *VacationList) GetAllByHR(uId int) {
	sqlS := `SELECT
					v.id,
					v.userId,
					u.title,
					v.typeId,
					v.startDate,
					v.duration,
					v.status
			FROM vacations AS v
			LEFT JOIN users AS u ON v.userId = u.id
			WHERE v.status > 1 
			ORDER BY v.startDate DESC;`

	l.getAllByMng(uId, sqlS, `HR`)
}

func (l *VacationList) getAllByMng(uId int, sqlS string, mngType string) {
	db := model.GetDB()
	defer db.Close()

	var rows *sql.Rows
	var err error
	if mngType == `FLM` {
		rows, err = db.Query(sqlS, uId)
	} else {
		rows, err = db.Query(sqlS)

	}
	model.CheckErr(err)
	defer rows.Close()
	usermap := make(map[int][]muser.VacancyBalance)
	nn := 1
	for rows.Next() {
		v := NewVacation()
		err := rows.Scan(&v.Id, &v.UserId, &v.UserTitle, &v.TypeId, &v.StartDate, &v.Duration, &v.Status)
		model.CheckErr(err)
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
					v.CurrentlyAvailable = `-`
					v.AsOfDec31 = `-`
				} else {
					v.CurrentlyAvailable = fmt.Sprintf("%.2f", w.TillNow)
					v.AsOfDec31 = fmt.Sprintf("%.2f", w.AvailableTillNY)
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
	// sqlI := `INSERT INTO vacations(userId, typeId, startDate, duration, status) VALUES(?,?,?,?,?)`
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
	// config.Statuses[str]
	_, err := db.Exec(sqlU, tmpId, vId)
	model.CheckErr(err)
}
