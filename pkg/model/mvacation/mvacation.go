package mvacation

import "vsystem/pkg/model"

type Vacation struct {
	Nn        int
	Id        int
	TypeId    int
	TypeTitle string
	UserId    int
	UserEmail string
	StartDate string
	EndDate   string
	Duration  int
	Status    int
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
					v.duration
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
		err := rows.Scan(&v.Id, &v.UserEmail, &v.TypeId, &v.StartDate, &v.Duration)
		model.CheckErr(err)
		v.Nn = nn
		nn++
		v.EndDate = model.GetEndDate(v.StartDate, v.Duration)
		v.StartDate = model.ReformatDate(v.StartDate)
		l.List = append(l.List, *v)
	}

}
