package course

import (
	"context"
	"errors"
	"os"

	"github.com/S3ergio31/curso-go-seccion-5-meta/meta"
	"github.com/S3ergio31/curso-go-seccion-5-response/response"
	"github.com/go-kit/kit/endpoint"
)

type Controller func(ctx context.Context, request any) (any, error)

type Endpoints struct {
	Create endpoint.Endpoint
	Get    endpoint.Endpoint
	GetAll endpoint.Endpoint
	Update endpoint.Endpoint
	Delete endpoint.Endpoint
}

type CreateRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type GetRequest struct {
	ID string
}

type GetAllRequest struct {
	Name  string
	Limit int
	Page  int
}

type DeleteRequest struct {
	ID string
}

type UpdateRequest struct {
	ID        string
	Name      *string `json:"name"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		createRequest := request.(CreateRequest)

		if createRequest.Name == "" {
			return nil, response.BadRequest(ErrorNameRequired.Error())
		}

		if createRequest.StartDate == "" {
			return nil, response.BadRequest(ErrorStartDateRequired.Error())
		}

		if createRequest.EndDate == "" {
			return nil, response.BadRequest(ErrorEndDateRequired.Error())
		}

		course, err := s.Create(
			createRequest.Name,
			createRequest.StartDate,
			createRequest.EndDate,
		)

		if errors.Is(err, ErrorEndLesserStart) {
			return nil, response.BadRequest(err.Error())
		}

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", course, nil), nil
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetRequest)
		course, err := s.Get(req.ID)

		if errors.As(err, &ErrorCourseNotFound{}) {
			return nil, response.NotFound(err.Error())
		}

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Ok("success", course, nil), nil
	}
}

func makeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetAllRequest)
		filters := Filters{
			Name: req.Name,
		}

		count, err := s.Count(filters)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, os.Getenv("PAGINATOR_LIMIT_DEFAULT"))

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		courses, err := s.GetAll(filters, meta.Offset(), meta.Limit())

		if err != nil {
			return nil, response.BadRequest(err.Error())
		}

		return response.Ok("success", courses, nil), nil
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		updateRequest := request.(UpdateRequest)

		if updateRequest.Name != nil && *updateRequest.Name == "" {
			return nil, response.BadRequest(ErrorNameRequired.Error())
		}

		if updateRequest.StartDate != nil && *updateRequest.StartDate == "" {
			return nil, response.BadRequest(ErrorStartDateRequired.Error())
		}

		if updateRequest.EndDate != nil && *updateRequest.EndDate == "" {
			return nil, response.BadRequest(ErrorEndDateRequired.Error())
		}

		err := s.Update(
			updateRequest.ID,
			updateRequest.Name,
			updateRequest.StartDate,
			updateRequest.EndDate,
		)

		if errors.As(err, &ErrorCourseNotFound{}) {
			return nil, response.NotFound(err.Error())
		}

		if errors.Is(err, ErrorEndLesserStart) {
			return nil, response.BadRequest(err.Error())
		}

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Ok("success", nil, nil), nil
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(DeleteRequest)
		err := s.Delete(req.ID)

		if errors.As(err, &ErrorCourseNotFound{}) {
			return nil, response.NotFound(err.Error())
		}

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Ok("success", nil, nil), nil
	}
}
