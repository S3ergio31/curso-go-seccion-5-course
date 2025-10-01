package course

import (
	"errors"
	"fmt"
)

var ErrorNameRequired = errors.New("name is required")

var ErrorStartDateRequired = errors.New("start date is required")

var ErrorEndDateRequired = errors.New("end date is required")

var ErrorEndLesserStart = errors.New("start date must be lower than end date")

type ErrorCourseNotFound struct {
	CourseID string
}

func (e ErrorCourseNotFound) Error() string {
	return fmt.Sprintf("course '%s' does not found", e.CourseID)
}
