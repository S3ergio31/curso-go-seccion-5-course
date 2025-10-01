package course

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/S3ergio31/curso-go-seccion-5-domain/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(course *domain.Course) error
	GetAll(filters Filters, offset, limit int) ([]domain.Course, error)
	Get(id string) (*domain.Course, error)
	Delete(id string) error
	Update(id string, name *string, startDate, endDate *time.Time) error
	Count(filters Filters) (int, error)
}

type repository struct {
	logger *log.Logger
	db     *gorm.DB
}

func (r repository) Create(course *domain.Course) error {
	if err := r.db.Create(course).Error; err != nil {
		r.logger.Println(err)
		return err
	}

	r.logger.Println("course created with id: ", course.ID)
	return nil
}

func (r repository) GetAll(filters Filters, offset, limit int) ([]domain.Course, error) {
	var courses []domain.Course

	tx := r.db.Model(&courses)

	tx = applyFilters(tx, filters)

	tx = tx.Limit(limit).Offset(offset)

	if err := tx.Order("created_at desc").Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r repository) Get(id string) (*domain.Course, error) {
	course := domain.Course{ID: id}
	err := r.db.First(&course).Error

	if err == gorm.ErrRecordNotFound {
		return nil, ErrorCourseNotFound{id}
	}

	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (r repository) Delete(id string) error {
	course := domain.Course{ID: id}

	result := r.db.Delete(&course)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrorCourseNotFound{id}
	}

	return nil
}

func (r repository) Update(id string, name *string, startDate, endDate *time.Time) error {
	values := make(map[string]interface{}, 0)

	if name != nil {
		values["name"] = *name
	}

	if startDate != nil {
		values["start_date"] = *startDate
	}

	if endDate != nil {
		values["end_date"] = *endDate
	}

	if err := r.db.Model(&domain.Course{}).Where("id = ?", id).Updates(values).Error; err != nil {
		return err
	}
	return nil
}

func (r repository) Count(filters Filters) (int, error) {
	var count int64

	tx := r.db.Model(domain.Course{})

	tx = applyFilters(tx, filters)

	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func NewRepository(logger *log.Logger, db *gorm.DB) Repository {
	return &repository{logger: logger, db: db}
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)
	}

	return tx
}
