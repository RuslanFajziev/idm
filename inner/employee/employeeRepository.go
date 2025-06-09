package employee

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type EmployeeRepository struct {
	db *sqlx.DB
}

type EmployeeEntity struct {
	Id     int64     `db:"id"`
	Name   string    `db:"name"`
	Create time.Time `db:"create_at"`
	Update time.Time `db:"update_at"`
}

func NewEmployeeRepository(database *sqlx.DB) *EmployeeRepository {
	return &EmployeeRepository{db: database}
}

func (rep *EmployeeRepository) SaveEmployee(entity EmployeeEntity) (employeeId int64, err error) {
	query := "INSERT INTO employee (name) VALUES ($1) RETURN id"
	err = rep.db.Get(&employeeId, query, entity.Name)
	return employeeId, err
}

func (rep *EmployeeRepository) FindById(id int64) (entity EmployeeEntity, err error) {
	query := "SELECT * FROM employee WHERE id = $1"
	err = rep.db.Get(&entity, query, id)
	return entity, err
}

func (rep *EmployeeRepository) GetAllEmployees() (entities []EmployeeEntity, err error) {
	query := "SELECT * FROM employee"
	err = rep.db.Select(&entities, query)
	return entities, err
}

func (rep *EmployeeRepository) FindEmployeesByIds(ids []int64) (entities []EmployeeEntity, err error) {
	query := "SELECT * FROM employee WHERE id IN (?)"
	query, args, err := sqlx.In(query, ids)

	if err != nil {
		return nil, err
	}

	query = sqlx.Rebind(0, query)
	err = rep.db.Select(&entities, query, args...)
	return entities, err
}

func (rep *EmployeeRepository) DeleteEmployeeById(id int64) error {
	query := "DELETE FROM employee WHERE id = $1"
	_, err := rep.db.Exec(query, id)
	return err
}

func (rep *EmployeeRepository) DeleteEmployeeByIds(ids []int64) error {
	query := "DELETE FROM employee WHERE id IN (?)"
	query, args, err := sqlx.In(query, ids)

	if err != nil {
		return err
	}

	query = sqlx.Rebind(0, query)
	_, err = rep.db.Exec(query, args...)
	return err
}
