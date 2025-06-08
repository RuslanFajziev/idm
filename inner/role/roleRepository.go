package role

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type RoleRepository struct {
	db *sqlx.DB
}

type RoleEntity struct {
	Id     int64     `db:"id"`
	Name   string    `db:"name"`
	Create time.Time `db:"create_at"`
	Update time.Time `db:"update_at"`
}

func NewRoleRepository(database *sqlx.DB) *RoleRepository {
	return &RoleRepository{db: database}
}

func (rep *RoleRepository) SaveRole(entity RoleEntity) (roleId int64, err error) {
	query := "INSERT INTO role (name) VALUES ($1) RETURN id"
	err = rep.db.Get(&roleId, query, entity.Name)
	return roleId, err
}

func (rep *RoleRepository) FindById(id int64) (entity RoleEntity, err error) {
	query := "SELECT * FROM role WHERE id = $1"
	err = rep.db.Get(&entity, query, id)
	return entity, err
}

func (rep *RoleRepository) GetAllRoles() (entities []RoleEntity, err error) {
	query := "SELECT * FROM role"
	err = rep.db.Select(&entities, query)
	return entities, err
}

func (rep *RoleRepository) FindRolesByIds(ids []int64) (entities []RoleEntity, err error) {
	for _, value := range ids {
		ent, err := rep.FindById(value)
		if err != nil {
			return entities, err
		}
		entities = append(entities, ent)
	}
	return entities, nil
}

func (rep *RoleRepository) DeleteRoleById(id int64) error {
	query := "DELETE FROM role WHERE id = $1"
	_, err := rep.db.Exec(query, id)
	return err
}

func (rep *RoleRepository) DeleteRoleByIds(ids []int64) error {
	for _, value := range ids {
		err := rep.DeleteRoleById(value)
		if err != nil {
			return err
		}
	}
	return nil
}
