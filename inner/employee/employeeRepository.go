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

func (rep *EmployeeRepository) AddEmployee(entity EmployeeEntity) (employeeId int64, err error) {
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
	for _, value := range ids {
		ent, err := rep.FindById(value)
		if err != nil {
			return entities, err
		}
		entities = append(entities, ent)
	}
	return entities, nil
}

func (rep *EmployeeRepository) DeleteEmployeeById(id int64) error {
	query := "DELETE FROM employee WHERE id = $1"
	_, err := rep.db.Exec(query, id)
	return err
}

func (rep *EmployeeRepository) DeleteEmployeeByIds(ids []int64) error {
	for _, value := range ids {
		err := rep.DeleteEmployeeById(value)
		if err != nil {
			return err
		}
	}
	return nil
}
