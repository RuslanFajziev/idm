package employee

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewEmployeeRepository(database *sqlx.DB) *Repository {
	return &Repository{db: database}
}

func (rep *Repository) BeginTransaction() (tx *sqlx.Tx, err error) {
	return rep.db.Beginx()
}

func (rep *Repository) FindByNameTx(tx *sqlx.Tx, name string) (isExists bool, err error) {
	err = tx.Get(&isExists, "SELECT EXISTS(SELECT 1 FROM employee WHERE name = $1)", name)
	return isExists, err
}

func (rep *Repository) SaveTx(tx *sqlx.Tx, entity *Entity) (id int64, err error) {
	query := "INSERT INTO employee (name, create_at, update_at) VALUES ($1, $2, $3) RETURNING id"
	err = tx.Get(&id, query, entity.Name, entity.Create, entity.Update)
	return id, err
}

func (rep *Repository) Save(entity *Entity) (id int64, err error) {
	query := "INSERT INTO employee (name) VALUES ($1) RETURNING id"
	err = rep.db.Get(&id, query, entity.Name)
	return id, err
}

func (rep *Repository) FindById(id int64) (entity Entity, err error) {
	query := "SELECT * FROM employee WHERE id = $1"
	err = rep.db.Get(&entity, query, id)
	return entity, err
}

func (rep *Repository) GetAll() (entities []Entity, err error) {
	query := "SELECT * FROM employee"
	err = rep.db.Select(&entities, query)
	return entities, err
}

func (rep *Repository) FindByIds(ids []int64) (entities []Entity, err error) {
	query := "SELECT * FROM employee WHERE id IN (?)"
	query, args, err := sqlx.In(query, ids)

	if err != nil {
		return nil, err
	}

	query = sqlx.Rebind(2, query)
	err = rep.db.Select(&entities, query, args...)
	return entities, err
}

func (rep *Repository) DeleteById(id int64) error {
	query := "DELETE FROM employee WHERE id = $1"
	_, err := rep.db.Exec(query, id)
	return err
}

func (rep *Repository) DeleteByIds(ids []int64) error {
	query := "DELETE FROM employee WHERE id IN (?)"
	query, args, err := sqlx.In(query, ids)

	if err != nil {
		return err
	}

	query = sqlx.Rebind(2, query)
	_, err = rep.db.Exec(query, args...)
	return err
}
