package tests

import (
	"idm/inner/database"
	"idm/inner/role"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryСase1(t *testing.T) {
	var env = ".env"
	Init(env)
	a := assert.New(t)
	db, err := database.ConnectDb(env)
	var clearDataBase = func() {
		db.MustExec("delete from role")
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

	var roleRepository = role.NewRoleRepository(db)
	var fixture = NewFixture(roleRepository)

	var testName = "test name"

	t.Run("Check FindById DeleteById GetAll", func(t *testing.T) {
		var newRoleId = fixture.Role(testName)

		entity, err := roleRepository.FindById(newRoleId)

		a.Nil(err)
		a.NotEmpty(entity)
		a.NotEmpty(entity.Id)
		a.NotEmpty(entity.Create)
		a.NotEmpty(entity.Update)
		a.Equal(testName, entity.Name)

		err = roleRepository.DeleteById(entity.Id)
		a.Nil(err)
		entities, err := roleRepository.GetAll()
		a.Nil(err)
		a.Equal(0, len(entities))
	})
}

func TestRepositoryСase2(t *testing.T) {
	var env = ".env"
	Init(env)
	a := assert.New(t)
	db, err := database.ConnectDb(env)
	var clearDataBase = func() {
		db.MustExec("delete from role")
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

	var roleRepository = role.NewRoleRepository(db)
	var fixture = NewFixture(roleRepository)
	var testName = "test name"
	var testName2 = "test name 2"

	t.Run("Check FindByIds GetAll DeleteByIds", func(t *testing.T) {
		var newRoleId = fixture.Role(testName)
		var newRoleId2 = fixture.Role(testName2)

		entities, err := roleRepository.FindByIds([]int64{newRoleId, newRoleId2})

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

		entities, err = roleRepository.GetAll()
		a.Nil(err)
		a.NotEmpty(entities)
		a.Equal(2, len(entities))
		a.Equal(testName, entities[0].Name)
		a.Equal(testName2, entities[1].Name)

		err = roleRepository.DeleteByIds([]int64{entities[0].Id, entities[1].Id})
		a.Nil(err)
		entities, err = roleRepository.GetAll()
		a.Nil(err)
		a.Equal(0, len(entities))
	})
}
