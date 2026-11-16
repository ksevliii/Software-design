package application

import (
	"errors"

	"pos/internal/domain"
)

var ErrCourseNotFound = errors.New("course not found")

type CourseRepository interface {
	Load(id domain.CourseID) (*domain.Course, error)
	Save(course *domain.Course) error
	ListIDs() ([]domain.CourseID, error)
}
