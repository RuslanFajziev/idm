package employee

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock" // библиотека для мокирования SQL-запросов в тестах
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert" // импортируем библиотеку с ассерт-функциями
	"github.com/stretchr/testify/mock"   // импортируем пакет для создания моков
)

func NewSqlmock() (*sqlx.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock, nil
}

// не удалось создать транзакцию
func TestBeginTransactionError(t *testing.T) {
	a := assert.New(t)
	db, mock, err := NewSqlmock()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()
	err = errors.New("connection error")
	mock.ExpectBegin().WillReturnError(err)

	t.Run("check error begin transation", func(t *testing.T) {
		repo := NewEmployeeRepository(db)
		srv := NewService(repo)

		id, errIn := srv.SaveTx(Request{Name: "Pupkin"})
		a.Equal(int64(0), id)
		a.Error(errIn)
		a.Equal(errIn.Error(), fmt.Errorf("error creating transaction: %w", err).Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// ошибка при проверке наличия работника с таким именем
func TestFindByName(t *testing.T) {
	a := assert.New(t)
	db, mock, err := NewSqlmock()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	err = errors.New("db error while searching by name")
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT 1 FROM employee").WillReturnError(err)
	mock.ExpectRollback()

	t.Run("check error while searching by name", func(t *testing.T) {
		repo := NewEmployeeRepository(db)
		srv := NewService(repo)

		id, errIn := srv.SaveTx(Request{Name: "Pupkin"})
		a.Equal(int64(0), id)
		a.Error(errIn)
		a.True(strings.Contains(errIn.Error(), fmt.Errorf("db error while searching by name").Error()))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// работник с таким именем уже есть в базе данных
func TestFindByName2(t *testing.T) {
	a := assert.New(t)
	db, mock, err := NewSqlmock()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(rows)
	mock.ExpectRollback()

	t.Run("check save employee, when a employee with that name exists", func(t *testing.T) {
		repo := NewEmployeeRepository(db)
		srv := NewService(repo)

		id, errIn := srv.SaveTx(Request{Name: "Pupkin"})
		a.Equal(int64(0), id)
		a.Error(errIn)
		a.True(strings.Contains(errIn.Error(), "employee already exists"))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// работника с таким именем нет в базе данных, но создание нового работника завершилось ошибкой
func TestFindByNameAndSave(t *testing.T) {
	a := assert.New(t)
	db, mock, err := NewSqlmock()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	req := Request{
		Name:   "Pupkin",
		Create: time.Now(),
		Update: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
	err = fmt.Errorf("db error when add new employee")
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(rows)
	mock.ExpectQuery("INSERT INTO employee").WillReturnError(err)
	mock.ExpectRollback()

	t.Run("check save employee, when save error", func(t *testing.T) {
		repo := NewEmployeeRepository(db)
		srv := NewService(repo)

		id, errIn := srv.SaveTx(req)
		a.Equal(int64(0), id)
		a.Error(errIn)
		a.True(strings.Contains(errIn.Error(), err.Error()))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// успешное создание нового работника
func TestFindByNameAndSave2(t *testing.T) {
	a := assert.New(t)
	db, mock, err := NewSqlmock()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	req := Request{
		Name:   "Pupkin",
		Create: time.Now(),
		Update: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
	rowsIds := sqlmock.NewRows([]string{"id"}).AddRow(int64(777))

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(rows)
	mock.ExpectQuery("INSERT INTO employee").WillReturnRows(rowsIds)
	mock.ExpectCommit()

	t.Run("check save employee, when save employee success", func(t *testing.T) {
		repo := NewEmployeeRepository(db)
		srv := NewService(repo)

		id, errIn := srv.SaveTx(req)
		a.Equal(int64(777), id)
		a.NoError(errIn)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// объявляем структуру мок-репозитория
type MockRepo struct {
	mock.Mock
}

func (rep *MockRepo) FindByNameTx(tx *sqlx.Tx, name string) (isExists bool, err error) {
	return true, nil
}

func (rep *MockRepo) BeginTransaction() (tx *sqlx.Tx, err error) {
	return nil, nil
}

func (s *MockRepo) SaveTx(tx *sqlx.Tx, entity *Entity) (id int64, err error) {
	return 99, nil
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

func (rep *StubRepo) FindByNameTx(tx *sqlx.Tx, name string) (isExists bool, err error) {
	return true, nil
}

func (rep *StubRepo) BeginTransaction() (tx *sqlx.Tx, err error) {
	return nil, nil
}

func (s *StubRepo) SaveTx(tx *sqlx.Tx, entity *Entity) (id int64, err error) {
	return 99, nil
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
