package employee

import (
	"encoding/json"
	"fmt"
	"idm/inner/common"
	"idm/inner/web"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Объявляем структуру мока сервиса employee.Service
type MockService struct {
	mock.Mock
}

// Реализуем функции мок-сервиса
func (srv *MockService) FindById(id int64) (Response, error) {
	args := srv.Called(id)
	return args.Get(0).(Response), args.Error(1)
}

func (srv *MockService) SaveTx(req Request) (id int64, err error) {
	args := srv.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func (srv *MockService) FindByIds(ids []int64) ([]Response, error) {
	args := srv.Called(ids)
	return args.Get(0).([]Response), args.Error(1)
}

func (srv *MockService) GetAll() ([]Response, error) {
	args := srv.Called()
	return args.Get(0).([]Response), args.Error(1)
}

func (srv *MockService) DeleteById(id int64) error {
	args := srv.Called(id)
	return args.Error(0)
}

func (srv *MockService) DeleteByIds(ids []int64) error {
	args := srv.Called(ids)
	return args.Error(0)
}

func GetLogger(t *testing.T) *common.Logger {
	loggerTest := zaptest.NewLogger(t)
	return &common.Logger{Logger: loggerTest}
}

func TestCreateEmployee(t *testing.T) {
	var a = assert.New(t)

	// тестируем положительный сценарий: работника создали и получили его id
	t.Run("should return created employee id", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		var svc = new(MockService)
		var controller = NewController(server, svc, GetLogger(t))
		controller.RegisterRoutes()

		var body = strings.NewReader("{\"name\": \"john doe\"}")
		var req = httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		svc.On("SaveTx", mock.AnythingOfType("Request")).Return(int64(123), nil)

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.ResponseBody[int64]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(int64(123), responseBody.Data)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})

	t.Run("should return error when saving employee", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		var svc = new(MockService)
		var controller = NewController(server, svc, GetLogger(t))
		controller.RegisterRoutes()
		// Готовим тестовое окружение
		var body = strings.NewReader("{\"name\": \"john doe\"}")
		var req = httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		var errMess = fmt.Errorf("employee with name %s already exists", "john doe").Error()
		svc.On("SaveTx", mock.AnythingOfType("Request")).Return(int64(0),
			common.AlreadyExistsError{Message: errMess})

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.ResponseBody[string]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
		a.Equal(errMess, responseBody.Message)
	})

	t.Run("should return error when saving employee v2", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		var svc = new(MockService)
		var controller = NewController(server, svc, GetLogger(t))
		controller.RegisterRoutes()

		var body = strings.NewReader("{\"name\": \"john doe\"}")
		var req = httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		var errMess1 = fmt.Errorf("database error")
		var errMess2 = fmt.Errorf("error finding employee by name: %s, %w", "john doe", errMess1).Error()
		svc.On("SaveTx", mock.AnythingOfType("Request")).Return(int64(0),
			common.DbOperationError{Message: errMess2})

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.ResponseBody[string]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
		a.Equal(errMess2, responseBody.Message)
	})
}

func TestContrlFindById(t *testing.T) {
	var a = assert.New(t)

	// тестируем положительный сценарий: работника создали и получили его id
	t.Run("should return employee by id", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		var svc = new(MockService)
		var controller = NewController(server, svc, GetLogger(t))
		controller.RegisterRoutes()
		// Готовим тестовое окружение
		var req = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees/id/123", nil)
		// req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		var entity = Response{
			Id:     123,
			Name:   "Pupkin",
			Create: time.Now(),
			Update: time.Now(),
		}
		svc.On("FindById", int64(123)).Return(entity, nil)

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.ResponseBody[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(entity.Name, responseBody.Data.Name)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})

	t.Run("should exception FindById", func(t *testing.T) {
		server := web.NewServer()
		var svc = new(MockService)
		var controller = NewController(server, svc, GetLogger(t))
		controller.RegisterRoutes()

		var req = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees/id/123", nil)
		var errMess1 = fmt.Errorf("database error")
		var errMess2 = fmt.Errorf("error finding employee by id: %s, %w", "123", errMess1).Error()
		svc.On("FindById", int64(123)).Return(Response{}, common.DbOperationError{Message: errMess2})

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.ResponseBody[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
		a.Equal(errMess2, responseBody.Message)
	})
}

func TestContrlGetAll(t *testing.T) {
	var a = assert.New(t)

	// тестируем положительный сценарий: работника создали и получили его id
	t.Run("should return employee all", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		var svc = new(MockService)
		var controller = NewController(server, svc, GetLogger(t))
		controller.RegisterRoutes()
		// Готовим тестовое окружение
		var req = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees", nil)
		// req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		var entity1 = Response{
			Id:     123,
			Name:   "Pupkin",
			Create: time.Now(),
			Update: time.Now(),
		}
		var entity2 = Response{
			Id:     321,
			Name:   "AnyName",
			Create: time.Now(),
			Update: time.Now(),
		}
		svc.On("GetAll").Return([]Response{entity1, entity2}, nil)

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.ResponseBody[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.True(len(responseBody.Data) == 2)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})

	t.Run("should exception GetAll", func(t *testing.T) {
		server := web.NewServer()
		var svc = new(MockService)
		var controller = NewController(server, svc, GetLogger(t))
		controller.RegisterRoutes()

		var req = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees", nil)
		var errMess1 = fmt.Errorf("database error")
		var errMess2 = fmt.Errorf("error finding employee by id: %s, %w", "123", errMess1).Error()
		svc.On("GetAll").Return([]Response{}, common.DbOperationError{Message: errMess2})

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.ResponseBody[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
		a.Equal(errMess2, responseBody.Message)
	})
}

func TestContrlDeleteById(t *testing.T) {
	var a = assert.New(t)

	// тестируем положительный сценарий: работника создали и получили его id
	t.Run("should DeleteById", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		var svc = new(MockService)
		var controller = NewController(server, svc, GetLogger(t))
		controller.RegisterRoutes()
		// Готовим тестовое окружение
		var req = httptest.NewRequest(fiber.MethodDelete, "/api/v1/employees/id/123", nil)

		svc.On("DeleteById", int64(123)).Return(nil)

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.ResponseBody[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})

	t.Run("should exception DeleteById", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		var svc = new(MockService)
		var controller = NewController(server, svc, GetLogger(t))
		controller.RegisterRoutes()
		// Готовим тестовое окружение
		var req = httptest.NewRequest(fiber.MethodDelete, "/api/v1/employees/id/123", nil)

		var errMess1 = fmt.Errorf("database error")
		var errMess2 = fmt.Errorf("error finding employee by id: %s, %w", "123", errMess1).Error()
		svc.On("DeleteById", int64(123)).Return(common.DbOperationError{Message: errMess2})

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.ResponseBody[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
		a.Equal(errMess2, responseBody.Message)
	})
}

func TestContrlFindByIds(t *testing.T) {
	var a = assert.New(t)

	// тестируем положительный сценарий: работника создали и получили его id
	t.Run("should return employees by ids", func(t *testing.T) {
		server := web.NewServer()
		var svc = new(MockService)
		var controller = NewController(server, svc, GetLogger(t))
		controller.RegisterRoutes()
		var req = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees/ids?ids=1,2,3", nil)
		var entity1 = Response{
			Id:     123,
			Name:   "Pupkin",
			Create: time.Now(),
			Update: time.Now(),
		}
		var entity2 = Response{
			Id:     321,
			Name:   "AnyName",
			Create: time.Now(),
			Update: time.Now(),
		}
		svc.On("FindByIds", []int64{1, 2, 3}).Return([]Response{entity1, entity2}, nil)

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.ResponseBody[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.True(len(responseBody.Data) == 2)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})

	t.Run("should exception by ids", func(t *testing.T) {
		server := web.NewServer()
		var svc = new(MockService)
		var controller = NewController(server, svc, GetLogger(t))
		controller.RegisterRoutes()
		var req = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees/ids?ids=1,2,3", nil)

		var errMess1 = fmt.Errorf("database error")
		var errMess2 = fmt.Errorf("error finding employees by ids: %s, %w", "1,2,3", errMess1).Error()
		svc.On("FindByIds", []int64{1, 2, 3}).Return([]Response{}, common.DbOperationError{Message: errMess2})

		resp, err := server.App.Test(req)

		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.ResponseBody[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
		a.Equal(errMess2, responseBody.Message)
	})
}

func TestContrlDeleteByIds(t *testing.T) {
	var a = assert.New(t)

	// тестируем положительный сценарий: работника создали и получили его id
	t.Run("should DeleteByIds", func(t *testing.T) {
		server := web.NewServer()
		var svc = new(MockService)
		var controller = NewController(server, svc, GetLogger(t))
		controller.RegisterRoutes()
		var req = httptest.NewRequest(fiber.MethodDelete, "/api/v1/employees/ids?ids=1,2,3", nil)
		svc.On("DeleteByIds", []int64{1, 2, 3}).Return(nil)

		resp, err := server.App.Test(req)

		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.ResponseBody[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})

	t.Run("should exception DeleteByIds", func(t *testing.T) {
		server := web.NewServer()
		var svc = new(MockService)
		var controller = NewController(server, svc, GetLogger(t))
		controller.RegisterRoutes()
		var req = httptest.NewRequest(fiber.MethodDelete, "/api/v1/employees/ids?ids=1,2,3", nil)

		var errMess1 = fmt.Errorf("database error")
		var errMess2 = fmt.Errorf("error finding employee by id: %s, %w", "123", errMess1).Error()
		svc.On("DeleteByIds", []int64{1, 2, 3}).Return(common.DbOperationError{Message: errMess2})

		resp, err := server.App.Test(req)

		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.ResponseBody[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.NotEmpty(responseBody.Message)
		a.Equal(errMess2, responseBody.Message)
	})
}
