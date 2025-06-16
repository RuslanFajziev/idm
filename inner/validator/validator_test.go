package validator

import (
	"idm/inner/employee"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidatorEmployeeReques(t *testing.T) {
	a := assert.New(t)
	validator := NewRequestValidator()

	t.Run("Check validation employee.Reques", func(t *testing.T) {
		var req = employee.Request{
			Name:   "John Doe",
			Create: time.Now(),
			Update: time.Now(),
		}

		err := validator.Validate(req)
		a.NoError(err)

		req = employee.Request{
			Name:   "J",
			Create: time.Now(),
			Update: time.Now(),
		}

		err = validator.Validate(req)
		a.Error(err)

		req = employee.Request{
			Name:   "John Doe",
			Update: time.Now(),
		}

		err = validator.Validate(req)
		a.Error(err)

		req = employee.Request{
			Name:   "John Doe",
			Create: time.Now(),
		}

		err = validator.Validate(req)
		a.Error(err)
	})
}

func TestValidatorEmployeeRequestById(t *testing.T) {
	a := assert.New(t)
	validator := NewRequestValidator()

	t.Run("Check validation employee.RequestById", func(t *testing.T) {
		var req = employee.RequestById{
			Id: 7777,
		}

		err := validator.Validate(req)
		a.NoError(err)

		req = employee.RequestById{
			Id: 0,
		}

		err = validator.Validate(req)
		a.Error(err)

		req = employee.RequestById{}

		err = validator.Validate(req)
		a.Error(err)
	})
}

func TestValidatorEmployeeRequestByIds(t *testing.T) {
	a := assert.New(t)
	validator := NewRequestValidator()

	t.Run("Check validation employee.RequestByIds", func(t *testing.T) {
		var req = employee.RequestByIds{
			Ids: []int64{777, 25},
		}

		err := validator.Validate(req)
		a.NoError(err)

		req = employee.RequestByIds{}

		err = validator.Validate(req)
		a.Error(err)
	})
}
