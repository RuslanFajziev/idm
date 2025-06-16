package employee

import (
	"time"

	_ "github.com/lib/pq"
)

type Entity struct {
	Id     int64     `db:"id"`
	Name   string    `db:"name"`
	Create time.Time `db:"create_at"`
	Update time.Time `db:"update_at"`
}

type Response struct {
	Id     int64     `json:"id"`
	Name   string    `json:"name"`
	Create time.Time `json:"create_at"`
	Update time.Time `json:"update_at"`
}

type Request struct {
	Name   string    `json:"name" validate:"required,min=2,max=155"`
	Create time.Time `json:"create_at" validate:"required"`
	Update time.Time `json:"update_at" validate:"required"`
}

type RequestById struct {
	Id int64 `json:"id" validate:"required,gt=0"`
}

type RequestByIds struct {
	Ids []int64 `json:"ids" validate:"required"`
}

func (e *Entity) toResponse() Response {
	return Response{
		Id:     e.Id,
		Name:   e.Name,
		Create: e.Create,
		Update: e.Update,
	}
}

func toResponses(entities []Entity) (responses []Response) {
	for _, e := range entities {
		responses = append(responses, e.toResponse())
	}

	return responses
}

func (r *Request) toEntity() *Entity {
	return &Entity{
		Name:   r.Name,
		Create: r.Create,
		Update: r.Update,
	}
}
