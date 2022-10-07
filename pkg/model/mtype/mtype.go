package mtype

import "vsystem/config"

type VacationType struct {
	TypeId    int
	TypeTitle string
	Amount    float64
	Rule      int
	IsUnlim   bool
}

type VacationTypeList struct {
	List []VacationType
}

func newVacationType() *VacationType {
	return &VacationType{}
}
func NewVacationTypeList() *VacationTypeList {
	return &VacationTypeList{}
}
func (l *VacationTypeList) GetVacationTypeList() {
	// Create UnPaid vacation
	v := newVacationType()
	v.TypeId = config.UnPaidVacationId
	v.TypeTitle = config.UnPaidVacationTitle
	v.Amount = 0
	v.Rule = config.AvailableImmediately
	v.IsUnlim = true
	l.List = append(l.List, *v)
	//Create Paid vacation
	v = newVacationType()
	v.TypeId = config.PaidVacationId
	v.TypeTitle = config.PaidVacationTitle
	v.Amount = 14
	v.Rule = config.AvailableGradually
	v.IsUnlim = false
	l.List = append(l.List, *v)

	//Create Sick Leave
	v = newVacationType()
	v.TypeId = config.SickLeaveId
	v.TypeTitle = config.SickLeaveTitle
	v.Amount = 5
	v.Rule = config.AvailableImmediately
	v.IsUnlim = false
	l.List = append(l.List, *v)
}
