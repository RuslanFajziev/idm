package employee

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert" // импортируем библиотеку с ассерт-функциями
	"github.com/stretchr/testify/mock"   // импортируем пакет для создания моков
)

// объявляем структуру мок-репозитория
type MockRepo struct {
	mock.Mock
}

// реализуем интерфейс репозитория у мока
func (m *MockRepo) Save(entity *Entity) (id int64, err error) {

	// Общая конфигурация поведения мок-объекта
	args := m.Called(entity)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepo) FindById(id int64) (employee Entity, err error) {
	args := m.Called(id)
	return args.Get(0).(Entity), args.Error(1)
}

func (m *MockRepo) GetAll() (entities []Entity, err error) {
	args := m.Called()
	return args.Get(0).([]Entity), args.Error(1)
}

func (m *MockRepo) FindByIds(ids []int64) (entities []Entity, err error) {
	args := m.Called(ids)
	return args.Get(0).([]Entity), args.Error(1)
}

func (m *MockRepo) DeleteById(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) DeleteByIds(ids []int64) error {
	args := m.Called(ids)
	return args.Error(0)
}

type StubRepo struct {
	entities map[int]*Entity
}

func NewStubRepo() *StubRepo {
	return &StubRepo{
		entities: map[int]*Entity{
			1: {Name: "Pupkin Vasia", Id: 1},
			2: {Name: "John Doe", Id: 2},
		},
	}
}

func (s *StubRepo) Save(entity *Entity) (id int64, err error) {
	if strings.EqualFold("Error Name", entity.Name) {
		return 0, fmt.Errorf("cannot save an bad object")
	}

	entity.Id = 3
	s.entities[3] = entity
	return entity.Id, nil
}

func (s *StubRepo) FindById(id int64) (employee Entity, err error) {
	if id == 0 {
		return Entity{}, fmt.Errorf("not found entity by %d", id)
	}
	return *s.entities[1], nil
}

func (s *StubRepo) GetAll() (entities []Entity, err error) {
	return []Entity{}, nil
}

func (s *StubRepo) FindByIds(ids []int64) (entities []Entity, err error) {
	return []Entity{}, nil
}

func (s *StubRepo) DeleteById(id int64) error {
	return nil
}

func (s *StubRepo) DeleteByIds(ids []int64) error {
	return nil
}

func TestSubSave(t *testing.T) {
	a := assert.New(t)

	t.Run("should return the id of the saved entity", func(t *testing.T) {
		repo := NewStubRepo()
		srv := NewService(repo)
		var request = Request{
			Name:   "Van Dam",
			Create: time.Now(),
			Update: time.Now(),
		}
		newId, err := srv.Save(request)

		a.Nil(err)
		a.Equal(int64(3), newId)
	})
	t.Run("should return error of the saved entity", func(t *testing.T) {
		repo := NewStubRepo()
		srv := NewService(repo)
		var request = Request{
			Name:   "Error Name",
			Create: time.Now(),
			Update: time.Now(),
		}
		newId, err := srv.Save(request)

		a.NotNil(err)
		a.Equal(int64(0), newId)
		a.True(strings.Contains(err.Error(), "cannot save an bad object"))
	})
}

func TestSubFindById(t *testing.T) {
	a := assert.New(t)

	t.Run("should return found employee", func(t *testing.T) {
		repo := NewStubRepo()
		srv := NewService(repo)
		response, err := srv.FindById(99)

		a.Nil(err)
		a.Equal("Pupkin Vasia", response.Name)
	})
}

func TestSave(t *testing.T) {
	a := assert.New(t)
	var request = Request{
		Name:   "John Doe",
		Create: time.Now(),
		Update: time.Now(),
	}

	t.Run("should return the id of the saved entity", func(t *testing.T) {
		repo := new(MockRepo)
		srv := NewService(repo)
		var id int64 = 5
		entity := request.toEntity()
		repo.On("Save", entity).Return(id, nil)
		got, err := srv.Save(request)

		a.Nil(err)
		a.Equal(id, got)
	})
	t.Run("should return error of the saved entity", func(t *testing.T) {
		repo := new(MockRepo)
		srv := NewService(repo)
		var id int64 = 0
		entity := request.toEntity()

		var err = errors.New("database error")
		var want = fmt.Errorf("error save employee: %w", err)

		repo.On("Save", entity).Return(id, err)
		newId, got := srv.Save(request)

		a.NotNil(err)
		a.Equal(id, newId)
		a.Equal(want, got)
	})
}

func TestFindById(t *testing.T) {

	// создаём экземпляр объекта с ассерт-функциями
	var a = assert.New(t)

	t.Run("should return found employee", func(t *testing.T) {
		// создаём экземпляр мок-объекта
		var repo = new(MockRepo)

		// создаём экземпляр сервиса, который собираемся тестировать. Передаём в его конструктор мок вместо реального репозитория
		var svc = NewService(repo)

		// создаём Entity, которую должен вернуть репозиторий
		var entity = Entity{
			Id:     1,
			Name:   "John Doe",
			Create: time.Now(),
			Update: time.Now(),
		}

		// создаём Response, который ожидаем получить от сервиса
		var want = entity.toResponse()

		// конфигурируем поведение мок-репозитория (при вызове метода FindById с аргументом 1 вернуть Entity, созданную нами выше)
		repo.On("FindById", int64(1)).Return(entity, nil)

		// вызываем сервис с аргументом id = 1
		var got, err = svc.FindById(1)

		// проверяем, что сервис не вернул ошибку
		a.Nil(err)

		// проверяем, что сервис вернул нам тот employee.Response, который мы ожилали получить
		a.Equal(want, got)
		// проверяем, что сервис вызвал репозиторий ровно 1 раз
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})

	t.Run("should return wrapped error", func(t *testing.T) {

		// Создаём для теста новый экземпляр мока репозитория.
		// Мы собираемся проверить счётчик вызовов, поэтому хотим, чтобы счётчик содержал количество вызовов к репозиторию,
		// выполненных в рамках одного нашего теста.
		// Ели сделать мок общим для нескольких тестов, то он посчитает вызовы, которые сделали все тесты
		var repo = new(MockRepo)

		// создаём новый экземпляр сервиса (чтобы передать ему новый мок репозитория)
		var svc = NewService(repo)

		// создаём пустую структуру employee.Entity, которую сервис вернёт вместе с ошибкой
		var entity = Entity{}

		// ошибка, которую вернёт репозиторий
		var err = errors.New("database error")

		// ошибка, которую должен будет вернуть сервис
		var want = fmt.Errorf("error finding employee with id 1: %w", err)

		repo.On("FindById", int64(1)).Return(entity, err)

		var response, got = svc.FindById(1)

		// проверяем результаты теста
		a.Empty(response)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})
}

func TestGetAll(t *testing.T) {
	a := assert.New(t)
	t.Run("return all entities", func(t *testing.T) {
		repo := new(MockRepo)
		srv := NewService(repo)
		listEntity := []Entity{Entity{Name: "name1"}, Entity{Name: "name2"}}
		repo.On("GetAll").Return(listEntity, nil)
		result, err := srv.GetAll()

		a.Nil(err)
		a.NotNil(result)
		a.Equal(len(listEntity), len(result))
		a.Equal(listEntity[0].Name, result[0].Name)
	})
	t.Run("return error when called return all entities", func(t *testing.T) {
		repo := new(MockRepo)
		srv := NewService(repo)

		err := errors.New("database error")
		want := fmt.Errorf("error GetAll employees: %w", err)

		repo.On("GetAll").Return([]Entity{}, err)
		result, err := srv.GetAll()

		a.Equal(result, []Response{})
		a.NotNil(err)
		a.Equal(want, err)
	})
}

func TestFindByIds(t *testing.T) {
	var a = assert.New(t)

	t.Run("should return found employees", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		var entity1 = Entity{
			Id:     1,
			Name:   "Pupkin Vasia",
			Create: time.Now(),
			Update: time.Now(),
		}
		var entity2 = Entity{
			Id:     2,
			Name:   "John Doe",
			Create: time.Now(),
			Update: time.Now(),
		}
		entities := []Entity{entity1, entity2}

		var want = toResponses(entities)
		var ids = []int64{1, 2}

		repo.On("FindByIds", ids).Return(entities, nil)
		result, err := svc.FindByIds(ids)

		a.Nil(err)
		a.Equal(want, result)
		a.Equal(want[1].Name, result[1].Name)
		a.True(repo.AssertNumberOfCalls(t, "FindByIds", 1))
	})

	t.Run("should return wrapped error", func(t *testing.T) {

		var repo = new(MockRepo)
		var svc = NewService(repo)
		entities := []Entity{}
		var ids = []int64{1, 2}

		var err = errors.New("database error")
		var want = fmt.Errorf("error finding employee with ids %d: %w", ids, err)

		repo.On("FindByIds", ids).Return(entities, err)

		response, err := svc.FindByIds(ids)

		a.Equal(response, []Response{})
		a.NotNil(err)
		a.Equal(want, err)
		a.True(repo.AssertNumberOfCalls(t, "FindByIds", 1))
	})
}

func TestDeleteById(t *testing.T) {
	var a = assert.New(t)
	t.Run("return nil when called DeleteById", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		var id int64 = 7
		repo.On("DeleteById", id).Return(nil).Once()
		err := svc.DeleteById(id)

		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "DeleteById", 1))
	})

	t.Run("return error when called DeleteById", func(t *testing.T) {

		var repo = new(MockRepo)
		var svc = NewService(repo)
		var id int64 = 7

		var err = errors.New("database error")
		var want = fmt.Errorf("error delete employee by id %d:  %w", id, err)

		repo.On("DeleteById", id).Return(want)
		err = svc.DeleteById(id)

		a.NotNil(err)
		a.True(strings.Contains(err.Error(), want.Error()))
		a.True(repo.AssertNumberOfCalls(t, "DeleteById", 1))
	})
}

func TestDeleteByIds(t *testing.T) {
	var a = assert.New(t)
	t.Run("return nil when called DeleteByIds", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		var ids = []int64{1, 2}
		repo.On("DeleteByIds", ids).Return(nil).Once()
		err := svc.DeleteByIds(ids)

		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "DeleteByIds", 1))
	})

	t.Run("return error when called DeleteByIds", func(t *testing.T) {

		var repo = new(MockRepo)
		var svc = NewService(repo)
		var ids = []int64{1, 2}

		var err = errors.New("database error")
		var want = fmt.Errorf("error delete employee by ids %d: %w", ids, err)

		repo.On("DeleteByIds", ids).Return(want)
		err = svc.DeleteByIds(ids)

		a.NotNil(err)
		a.True(strings.Contains(err.Error(), want.Error()))
		a.True(repo.AssertNumberOfCalls(t, "DeleteByIds", 1))
	})
}
