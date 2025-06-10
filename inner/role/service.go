package role

import "fmt"

type Service struct {
	repo Repo
}

type Repo interface {
	Save(entity *Entity) (id int64, err error)
	FindById(id int64) (entity Entity, err error)
	GetAll() (entities []Entity, err error)
	FindByIds(ids []int64) (entities []Entity, err error)
	DeleteById(id int64) error
	DeleteByIds(ids []int64) error
}

func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (serv *Service) Save(req Request) (id int64, err error) {
	id, err = serv.repo.Save(req.toEntity())
	if err != nil {
		return 0, fmt.Errorf("error save role: %w", err)
	}

	return id, nil
}

func (serv *Service) FindById(id int64) (Response, error) {
	resp, err := serv.repo.FindById(id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding role with id %d: %w", id, err)
	}

	return resp.toResponse(), nil
}

func (serv *Service) GetAll() ([]Response, error) {
	resps, err := serv.repo.GetAll()
	if err != nil {
		return []Response{}, fmt.Errorf("error GetAll roles: %w", err)
	}

	return toResponses(resps), nil
}

func (serv *Service) FindByIds(ids []int64) ([]Response, error) {
	resps, err := serv.repo.FindByIds(ids)
	if err != nil {
		return []Response{}, fmt.Errorf("error finding role with ids %d: %w", ids, err)
	}

	return toResponses(resps), nil
}

func (serv *Service) DeleteById(id int64) error {
	err := serv.repo.DeleteById(id)
	if err != nil {
		return fmt.Errorf("error delete role by id %d: %w", id, err)
	}

	return nil
}

func (serv *Service) DeleteByIds(ids []int64) error {
	err := serv.repo.DeleteByIds(ids)
	if err != nil {
		return fmt.Errorf("error delete role by ids %d: %w", ids, err)
	}

	return nil
}
