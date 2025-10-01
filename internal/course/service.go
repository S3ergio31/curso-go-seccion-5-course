package course

import (
	"log"
	"time"

	"github.com/S3ergio31/curso-go-seccion-5-domain/domain"
)

type Service interface {
	Create(name, startDate, endDate string) (*domain.Course, error)
	GetAll(filters Filters, offset, limit int) ([]domain.Course, error)
	Get(id string) (*domain.Course, error)
	Delete(id string) error
	Update(id string, name, startDate, endDate *string) error
	Count(filters Filters) (int, error)
}

type service struct {
	logger     *log.Logger
	repository Repository
}

type Filters struct {
	Name string
}

func (s service) Create(name, startDate, endDate string) (*domain.Course, error) {

	startDateParsed, err := time.Parse("2006-01-02", startDate)

	if err != nil {
		s.logger.Println(err)
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)

	if err != nil {
		s.logger.Println(err)
		return nil, err
	}

	course := &domain.Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	if err := s.repository.Create(course); err != nil {
		return nil, err
	}

	return course, nil
}

func (s service) GetAll(filters Filters, offset, limit int) ([]domain.Course, error) {
	courses, err := s.repository.GetAll(filters, offset, limit)

	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (s service) Get(id string) (*domain.Course, error) {
	course, err := s.repository.Get(id)

	if err != nil {
		return nil, err
	}

	return course, nil
}

func (s service) Delete(id string) error {
	return s.repository.Delete(id)
}

func (s service) Update(id string, name, startDate, endDate *string) error {
	var startDateParsed, endDateParsed *time.Time

	if startDate != nil {
		date, err := time.Parse("2006-01-02", *startDate)

		if err != nil {
			s.logger.Println(err)
			return err
		}

		startDateParsed = &date
	}

	if endDate != nil {
		date, err := time.Parse("2006-01-02", *endDate)

		if err != nil {
			s.logger.Println(err)
			return err
		}

		endDateParsed = &date
	}

	return s.repository.Update(id, name, startDateParsed, endDateParsed)
}

func (s service) Count(filters Filters) (int, error) {
	return s.repository.Count(filters)
}

func NewService(repository Repository, logger *log.Logger) Service {
	return &service{logger: logger, repository: repository}
}
