package tests

import (
	"idm/inner/employee"
)

type Fixture struct {
	employee *employee.EmployeeRepository
}

func NewFixture(employee *employee.EmployeeRepository) *Fixture {
	return &Fixture{employee}
}

func (f *Fixture) Employee(name string) int64 {
	var entity = employee.EmployeeEntity{
		Name: name,
	}
	newId, err := f.employee.Save(&entity)
	if err != nil {
		panic(err)
	}
	return newId
}
