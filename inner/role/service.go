package role

import (
	"fmt"
	"idm/inner/common"
)

type Service struct {
	repo  Repo
	valid Validator
}

type Repo interface {
	Save(entity *Entity) (id int64, err error)
	FindById(id int64) (entity Entity, err error)
	GetAll() (entities []Entity, err error)
	FindByIds(ids []int64) (entities []Entity, err error)
	DeleteById(id int64) error
	DeleteByIds(ids []int64) error
	FindByName(name string) (isExists bool, err error)
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

func (serv *Service) Save(req Request) (id int64, err error) {
	// валидируем запрос
	err = serv.valid.Validate(req)
	if err != nil {
		// возвращаем кастомную ошибку в случае, если запрос не прошёл валидацию (про кастомные ошибки - дальше)
		return 0, common.RequestValidationError{Message: err.Error()}
	}

	isExists, err := serv.repo.FindByName(req.Name)
	if err != nil {
		return 0, common.DbOperationError{Message: fmt.Errorf("error finding employee by name: %s, %w", req.Name, err).Error()}
	}
	if isExists {
		return 0, common.AlreadyExistsError{Message: fmt.Errorf("employee with name %s already exists", req.Name).Error()}
	}

	id, err = serv.repo.Save(req.toEntity())
	if err != nil {
		return 0, common.DbOperationError{Message: fmt.Errorf("error save role: %w", err).Error()}
	}

	return id, nil
}

func (serv *Service) FindById(id int64) (Response, error) {
	resp, err := serv.repo.FindById(id)
	if err != nil {
		return Response{}, common.DbOperationError{Message: fmt.Errorf("error finding role with id %d: %w", id, err).Error()}
	}

	return resp.toResponse(), nil
}

func (serv *Service) GetAll() ([]Response, error) {
	resps, err := serv.repo.GetAll()
	if err != nil {
		return []Response{}, common.DbOperationError{Message: fmt.Errorf("error GetAll roles: %w", err).Error()}
	}

	return toResponses(resps), nil
}

func (serv *Service) FindByIds(ids []int64) ([]Response, error) {
	resps, err := serv.repo.FindByIds(ids)
	if err != nil {
		return []Response{}, common.DbOperationError{Message: fmt.Errorf("error finding role with ids %d: %w", ids, err).Error()}
	}

	return toResponses(resps), nil
}

func (serv *Service) DeleteById(id int64) error {
	err := serv.repo.DeleteById(id)
	if err != nil {
		return common.DbOperationError{Message: fmt.Errorf("error delete role by id %d: %w", id, err).Error()}
	}

	return nil
}

func (serv *Service) DeleteByIds(ids []int64) error {
	err := serv.repo.DeleteByIds(ids)
	if err != nil {
		return common.DbOperationError{Message: fmt.Errorf("error delete role by ids %d: %w", ids, err).Error()}
	}

	return nil
}
