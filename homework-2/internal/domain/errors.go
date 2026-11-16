package domain

import "errors"

var (
	ErrScoreOutOfRange   = errors.New("score out of range")
	ErrStudentExists     = errors.New("student already exists")
	ErrStudentNotFound   = errors.New("student not found")
	ErrCriterionNotFound = errors.New("criterion not found")
	ErrRubricLocked      = errors.New("rubricator is locked")
	ErrGradeExists       = errors.New("grade already exists")
	ErrGradeNotFound     = errors.New("grade not found")
	ErrIncompleteGrades  = errors.New("student has missing grades")
)
