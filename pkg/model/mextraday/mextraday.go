package mextraday

import "vsystem/pkg/model"

type ExtraDay struct {
	Date        string
	Description string
	Year        string
}

func NewExtraDay() *ExtraDay {
	return &ExtraDay{}
}

type ExtraDayList struct {
	List []ExtraDay
}

func NewExtraDayList() *ExtraDayList {
	return &ExtraDayList{}
}

// Return list of extradays order by date DESC
func (l *ExtraDayList) GetAll() {
	db := model.GetDB()
	defer db.Close()
	sqlS := `SELECT extra_day, description, strftime('%Y',extra_day)
				FROM extra_days
				WHERE strftime('%Y',extra_day) = strftime('%Y','now')
				OR
					CASE
		 				WHEN strftime('%Y','now') > 6
			 			THEN strftime('%Y',extra_day) = strftime('%Y',date('now','start of year', '+1 year'))
						ELSE strftime('%Y',extra_day) = strftime('%Y',date('now','start of year', '-1 year'))
					END
				ORDER BY extra_day DESC`
	rows, err := db.Query(sqlS)
	model.CheckErr(err)
	defer rows.Close()

	for rows.Next() {
		d := NewExtraDay()
		err := rows.Scan(&d.Date, &d.Description, &d.Year)
		model.CheckErr(err)
		l.List = append(l.List, *d)
	}
}

func Save2DB(id, d, c string) {
	db := model.GetDB()
	defer db.Close()
	if id == `0` {
		sqlI := `INSERT INTO extra_days(extra_day, description) VALUES(?,?)`
		_, err := db.Exec(sqlI, d, c)
		model.CheckErr(err)
	} else {
		sqlD := `DELETE FROM extra_days WHERE extra_day = ?`
		_, err := db.Exec(sqlD, id)
		model.CheckErr(err)
	}
}
