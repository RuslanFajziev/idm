package tests

import (
	"idm/inner/database"
	"idm/inner/employee"
	"idm/inner/role"
	"os"
)

func ClearEnv() {
	driverName := "DB_DRIVER_NAME"
	dsnName := "DB_DSN"
	os.Unsetenv(driverName)
	os.Unsetenv(dsnName)
}

type FixtureRole struct {
	roles *role.RoleRepository
}

type FixtureEmployee struct {
	employee *employee.EmployeeRepository
}

func NewFixtureRole(roles *role.RoleRepository) *FixtureRole {
	return &FixtureRole{roles}
}

func NewFixtureEmployee(employee *employee.EmployeeRepository) *FixtureEmployee {
	return &FixtureEmployee{employee}
}

func (f *FixtureRole) Role(name string) int64 {
	var entity = role.RoleEntity{
		Name: name,
	}
	newId, err := f.roles.Save(&entity)
	if err != nil {
		panic(err)
	}
	return newId
}

func (f *FixtureEmployee) Employee(name string) int64 {
	var entity = employee.EmployeeEntity{
		Name: name,
	}
	newId, err := f.employee.Save(&entity)
	if err != nil {
		panic(err)
	}
	return newId
}

func Init(envFile string) {
	db, err := database.ConnectDb(envFile)
	if err == nil {
		defer db.Close()
	}
	if err != nil {
		panic(err)
	}

	content, err := os.ReadFile("InitTestTables.sql")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(string(content))
	if err != nil {
		panic(err)
	}
}
