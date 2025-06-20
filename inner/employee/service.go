package employee

import (
	"fmt"

	"idm/inner/common"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	repo  Repo
	valid Validator
}

type Repo interface {
	FindByNameTx(tx *sqlx.Tx, name string) (isExists bool, err error)
	BeginTransaction() (tx *sqlx.Tx, err error)
	SaveTx(tx *sqlx.Tx, entity *Entity) (id int64, err error)
	Save(entity *Entity) (id int64, err error)
	FindById(id int64) (entity Entity, err error)
	GetAll() (entities []Entity, err error)
	FindByIds(ids []int64) (entities []Entity, err error)
	DeleteById(id int64) error
	DeleteByIds(ids []int64) error
}

type Validator interface {
	Validate(request any) error
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{
		repo:  repo,
		valid: validator,
	}
}

func (serv *Service) SaveTx(req Request) (id int64, err error) {
	// валидируем запрос (про валидатор расскажу дальше)
	err = serv.valid.Validate(req)
	if err != nil {
		// возвращаем кастомную ошибку в случае, если запрос не прошёл валидацию (про кастомные ошибки - дальше)
		return 0, common.RequestValidationError{Message: err.Error()}
	}

	tx, err := serv.repo.BeginTransaction()
	if err != nil {
		return 0, fmt.Errorf("error creating transaction: %w", err)
	}

	defer func() {
		// проверяем, не было ли паники
		if r := recover(); r != nil {
			err = fmt.Errorf("creating employee panic: %v", r)
			// если была паника, то откатываем транзакцию
			errTx := tx.Rollback()
			if errTx != nil {
				err = fmt.Errorf("creating employee: rolling back transaction errors: %w, %w", err, errTx)
			}
		} else if err != nil {
			// если произошла другая ошибка (не паника), то откатываем транзакцию
			errTx := tx.Rollback()
			if errTx != nil {
				err = fmt.Errorf("creating employee: rolling back transaction errors: %w, %w", err, errTx)
			}
		} else {
			// если ошибок нет, то коммитим транзакцию
			errTx := tx.Commit()
			if errTx != nil {
				err = fmt.Errorf("creating employee: commiting transaction error: %w", errTx)
			}
		}
	}()

	isExists, err := serv.repo.FindByNameTx(tx, req.Name)
	if err != nil {
		return 0, common.DbOperationError{Message: fmt.Errorf("error finding employee by name: %s, %w", req.Name, err).Error()}
	}
	if isExists {
		return 0, common.AlreadyExistsError{Message: fmt.Errorf("employee with name %s already exists", req.Name).Error()}
	}

	newId, err := serv.repo.SaveTx(tx, req.toEntity())
	if err != nil {
		err = fmt.Errorf("error creating employee with name: %s %v", req.Name, err)
	}
	return newId, err
}

func (serv *Service) Save(req Request) (id int64, err error) {
	id, err = serv.repo.Save(req.toEntity())
	if err != nil {
		return 0, fmt.Errorf("error save employee: %w", err)
	}

	return id, nil
}

func (serv *Service) FindById(id int64) (Response, error) {
	resp, err := serv.repo.FindById(id)
	if err != nil {
		return Response{}, common.DbOperationError{Message: fmt.Errorf("error finding employee with id %d: %w", id, err).Error()}

	}

	return resp.toResponse(), nil
}

func (serv *Service) GetAll() ([]Response, error) {
	resps, err := serv.repo.GetAll()
	if err != nil {
		return []Response{}, common.DbOperationError{Message: fmt.Errorf("error get all employees: %w", err).Error()}
	}

	return toResponses(resps), nil
}

func (serv *Service) FindByIds(ids []int64) ([]Response, error) {
	resps, err := serv.repo.FindByIds(ids)
	if err != nil {
		return []Response{}, common.DbOperationError{Message: fmt.Errorf("error finding employee with ids %d: %w", ids, err).Error()}
	}

	return toResponses(resps), nil
}

func (serv *Service) DeleteById(id int64) error {
	err := serv.repo.DeleteById(id)
	if err != nil {
		return common.DbOperationError{Message: fmt.Errorf("error delete employee by id %d: %w", id, err).Error()}
	}

	return nil
}

func (serv *Service) DeleteByIds(ids []int64) error {
	err := serv.repo.DeleteByIds(ids)
	if err != nil {
		return common.DbOperationError{Message: fmt.Errorf("error delete employee by ids %d: %w", ids, err).Error()}
	}

	return nil
}
