package muser

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"vsystem/config"
	"vsystem/pkg/model"
	"vsystem/pkg/model/mtype"
)

type User struct {
	Id             int
	Title          string
	Email          string
	UserType       string
	FLMId          int
	FLMTitle       string
	StartDate      string
	ExtraDays      float64
	SpillOver      float64
	IsActive       bool
	VacancyBalance []VacancyBalance
	FLMs           []FLM
}
type FLM struct {
	Id    int
	Title string
}

type VacancyBalance struct {
	Id          int
	Title       string
	Rule        int
	IsAvailable bool
	AsOfDec31   float64
	AsOfNow     float64
	Spent       float64
	IsUnLim     bool
}

func NewUser() *User {
	return &User{}
}

type UserList struct {
	List []User
}

func NewUserList() *UserList {
	return &UserList{}
}

func (l *UserList) GetAll() {

	db := model.GetDB()
	defer db.Close()
	sqlS := `SELECT u.id, u.title, u.email, u.startDate, u.extraDays, u.spillover, u.IsActive, u.userType, flm.id, flm.title
			FROM users AS u
    			LEFT OUTER JOIN users AS flm ON u.flm = flm.id
			ORDER BY u.IsActive DESC, u.title ASC;`
	rows, err := db.Query(sqlS)
	model.CheckErr(err)
	defer rows.Close()

	for rows.Next() {
		u := NewUser()
		var tmpInt sql.NullInt64
		var tmpStr sql.NullString
		err := rows.Scan(&u.Id, &u.Title, &u.Email, &u.StartDate, &u.ExtraDays, &u.SpillOver, &u.IsActive, &u.UserType, &tmpInt, &tmpStr)
		model.CheckErr(err)
		u.FLMId = int(tmpInt.Int64)
		u.FLMTitle = tmpStr.String
		// u.GetVacationBalabce(u.Id, u.StartDate, u.ExtraDays, u.SpillOver)

		u.StartDate = model.ReformatDate(u.StartDate)
		if err != nil {
			log.Println(err)
		}
		l.List = append(l.List, *u)
	}
}

func (u *User) GetByEmail(email string) {
	db := model.GetDB()
	defer db.Close()
	sqlS := `SELECT u.id, u.title, u.email, u.userType FROM users AS u WHERE u.email = ? AND u.IsActive = true`
	err := db.QueryRow(sqlS, email).Scan(&u.Id, &u.Title, &u.Email, &u.UserType)
	if err != nil {
		u.Id = 0
		log.Println(err)
	}
}

func (u *User) GetById(uId int) {

	if uId != 0 {
		db := model.GetDB()
		defer db.Close()

		sqlS := `SELECT u.id, u.title, u.email, u.startDate, u.extraDays, u.spillover, u.IsActive, u.userType,  flm.id, flm.title 
    				FROM users AS u 
    					LEFT OUTER JOIN users AS flm ON u.flm = flm.id
    				WHERE u.id = ?`
		db.QueryRow(sqlS, uId).Scan(&u.Id, &u.Title, &u.Email, &u.StartDate, &u.ExtraDays, &u.SpillOver, &u.IsActive, &u.UserType, &u.FLMId, &u.FLMTitle)
	} else {
		u.StartDate = fmt.Sprintf("%02d-%02d-%4d", time.Now().Day(), time.Now().Month(), time.Now().Year())
	}
	u.GetVacationBalabce(uId, u.StartDate, u.ExtraDays, u.SpillOver)
	u.getFLMList()

}

func (u *User) getFLMList() {
	db := model.GetDB()
	defer db.Close()
	sqlS := `SELECT id, title FROM users WHERE userType = ? AND IsActive = true ORDER BY title ASC`
	rows, err := db.Query(sqlS, config.UTypeFLM)
	model.CheckErr(err)
	defer rows.Close()

	for rows.Next() {
		flm := FLM{}
		err := rows.Scan(&flm.Id, &flm.Title)
		model.CheckErr(err)
		u.FLMs = append(u.FLMs, flm)
	}

}

func (u *User) GetVacationBalabce(uId int, dataStart string, ExtraD, SD float64) {
	db := model.GetDB()
	defer db.Close()
	vacationList := mtype.NewVacationTypeList()
	vacationList.GetVacationTypeList()
	for _, v := range vacationList.List {

		vb := VacancyBalance{}
		vb.Id = v.TypeId
		vb.Title = v.TypeTitle
		vb.AsOfDec31 = v.Amount
		vb.Rule = v.Rule
		vb.IsUnLim = v.IsUnlim

		sqlS2 := `SELECT COUNT(*) FROM user_type_vacation WHERE userId = ? AND TypeVacationId = ?`
		tmpInt := 0
		err := db.QueryRow(sqlS2, uId, vb.Id).Scan(&tmpInt)
		model.CheckErr(err)
		vb.IsAvailable = false
		if tmpInt > 0 {
			vb.IsAvailable = true
			var tmpF sql.NullFloat64
			sqlS2 = `SELECT SUM(duration) FROM vacations 
						WHERE 
						userId = ? 
							AND 
						typeId = ? 
							AND 
						strftime('%Y',startDate) = strftime('%Y','now')
							AND
						status = ?`
			err = db.QueryRow(sqlS2, uId, vb.Id, config.AcceptedByHR).Scan(&tmpF)
			model.CheckErr(err)
			vb.Spent = tmpF.Float64
			if vb.Rule == config.AvailableImmediately {
				vb.AsOfNow = vb.AsOfDec31
			} else {
				currentTime := time.Now()
				//Check start date
				dStart, err := time.Parse("2006-01-02", dataStart)
				model.CheckErr(err)
				if currentTime.Year() > dStart.Year() {
					vb.AsOfNow = float64(currentTime.Month()-1) * vb.AsOfDec31 / 12
				} else {
					vb.AsOfDec31 = (12 - float64(dStart.Month()) + 1) * vb.AsOfDec31 / 12
					vb.AsOfNow = float64(currentTime.Month()-1) * vb.AsOfDec31 / 12
				}
				if vb.Id == config.PaidVacationId {
					vb.AsOfDec31 += ExtraD + SD
					vb.AsOfNow += ExtraD + SD
				}
			}
			vb.AsOfDec31 -= vb.Spent
			vb.AsOfNow -= vb.Spent
		}
		vb.AsOfDec31 = model.RoundDays(vb.AsOfDec31)
		vb.AsOfNow = model.RoundDays(vb.AsOfNow)
		u.VacancyBalance = append(u.VacancyBalance, vb)
	}

}

func (u *User) Save2DB(uId int) {
	db := model.GetDB()
	defer db.Close()
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	model.CheckErr(err)
	defer tx.Rollback()

	var sql string
	if uId == 0 {
		sql = `INSERT INTO users(title, email, startDate, extraDays, spillover, IsActive, userType, flm) VALUES(?,?,?,?,?,?,?,?)`
		_, err := tx.ExecContext(ctx, sql, u.Title, u.Email, u.StartDate, u.ExtraDays, u.SpillOver, u.IsActive, u.UserType, u.FLMId)
		model.CheckErr(err)
		err = tx.QueryRowContext(ctx, `SELECT last_insert_rowid()`).Scan(&uId)

		model.CheckErr(err)

	} else {
		sql = `UPDATE users SET
		 title = ?, email = ?, startDate = ?, extraDays = ?, spillover = ?, IsActive = ?, userType = ?, flm = ?,
		 updated_at = datetime('now', 'localtime') WHERE id = ?`
		_, err := tx.ExecContext(ctx, sql, u.Title, u.Email, u.StartDate, u.ExtraDays, u.SpillOver, u.IsActive, u.UserType, u.FLMId, uId)
		model.CheckErr(err)
		sql = `DELETE FROM user_type_vacation WHERE userId = ?`
		_, err = tx.ExecContext(ctx, sql, uId)
		model.CheckErr(err)
	}

	for _, v := range u.VacancyBalance {

		sql = `INSERT INTO user_type_vacation(userId, TypeVacationId)VALUES(?,?)`
		_, err := tx.ExecContext(ctx, sql, uId, v.Id)
		model.CheckErr(err)
	}
	u.Id = uId

	err = tx.Commit()
	model.CheckErr(err)
}
