package tests

import (
	"idm/inner/database"
	"idm/inner/employee"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestSaveTx(t *testing.T) {
	var env = ".env"
	Init(env)
	a := assert.New(t)
	db, err := database.ConnectDb(env)
	clearDataBase := func() {
		db.MustExec("delete from employee")
	}
	if err == nil {
		clearDataBase()
		defer db.Close()
	}

	defer func() {
		if r := recover(); r != nil {
			clearDataBase()
		}
	}()

	repo := employee.NewEmployeeRepository(db)

	entity := employee.Entity{
		Name:   "Pupkin",
		Create: time.Now(),
		Update: time.Now(),
	}
	t.Run("Check save employee in trancation", func(t *testing.T) {
		tx, err := repo.BeginTransaction()
		a.NoError(err)
		id, err := repo.SaveTx(tx, &entity)
		a.NoError(err)
		a.True(id > 0)
		err = tx.Commit()
		a.NoError(err)

		tx, err = repo.BeginTransaction()
		a.NoError(err)
		isExists, err := repo.FindByNameTx(tx, entity.Name)
		a.NoError(err)
		a.True(isExists)
		err = tx.Rollback()
		a.NoError(err)
	})
}

func TestEmployeeRepositoryСase1(t *testing.T) {
	var env = ".env"
	Init(env)
	a := assert.New(t)
	db, err := database.ConnectDb(env)
	var clearDataBase = func() {
		db.MustExec("delete from employee")
	}
	if err == nil {
		clearDataBase()
		defer db.Close()
	}

	defer func() {
		if r := recover(); r != nil {
			clearDataBase()
		}
	}()

	var repository = employee.NewEmployeeRepository(db)
	var fixture = NewFixtureEmployee(repository)

	var testName = "test name"

	t.Run("Check FindById DeleteById GetAll", func(t *testing.T) {
		var newRoleId = fixture.Employee(testName)

		entity, err := repository.FindById(newRoleId)

		a.Nil(err)
		a.NotEmpty(entity)
		a.NotEmpty(entity.Id)
		a.NotEmpty(entity.Create)
		a.NotEmpty(entity.Update)
		a.Equal(testName, entity.Name)

		err = repository.DeleteById(entity.Id)
		a.Nil(err)
		entities, err := repository.GetAll()
		a.Nil(err)
		a.Equal(0, len(entities))
	})
}

func TestEmployeeRepositoryСase2(t *testing.T) {
	var env = ".env"
	Init(env)
	a := assert.New(t)
	db, err := database.ConnectDb(env)
	var clearDataBase = func() {
		db.MustExec("delete from employee")
	}
	if err == nil {
		clearDataBase()
		defer db.Close()
	}

	defer func() {
		if r := recover(); r != nil {
			clearDataBase()
		}
	}()

	var repository = employee.NewEmployeeRepository(db)
	var fixture = NewFixtureEmployee(repository)
	var testName = "test name"
	var testName2 = "test name 2"

	t.Run("Check FindByIds GetAll DeleteByIds", func(t *testing.T) {
		var newRoleId = fixture.Employee(testName)
		var newRoleId2 = fixture.Employee(testName2)

		entities, err := repository.FindByIds([]int64{newRoleId, newRoleId2})

		a.Nil(err)
		a.NotEmpty(entities)
		a.NotEmpty(entities[0].Id)
		a.NotEmpty(entities[1].Id)
		a.NotEmpty(entities[0].Create)
		a.NotEmpty(entities[1].Create)
		a.NotEmpty(entities[0].Update)
		a.NotEmpty(entities[1].Update)
		a.Equal(testName, entities[0].Name)
		a.Equal(testName2, entities[1].Name)

		entities, err = repository.GetAll()
		a.Nil(err)
		a.NotEmpty(entities)
		a.Equal(2, len(entities))
		a.Equal(testName, entities[0].Name)
		a.Equal(testName2, entities[1].Name)

		err = repository.DeleteByIds([]int64{entities[0].Id, entities[1].Id})
		a.Nil(err)
		entities, err = repository.GetAll()
		a.Nil(err)
		a.Equal(0, len(entities))
	})
}
