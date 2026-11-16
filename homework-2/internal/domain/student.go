package domain

import (
	"errors"
	"strings"
)

type StudentID string

type Student struct {
	id   StudentID
	name string
}

func NewStudent(id StudentID, name string) (Student, error) {
	student := Student{id: StudentID(strings.TrimSpace(string(id))), name: strings.TrimSpace(name)}
	if err := student.validate(); err != nil {
		return Student{}, err
	}
	return student, nil
}

func (s Student) validate() error {
	if s.id == "" {
		return errors.New("student ID must not be empty")
	}
	if s.name == "" {
		return errors.New("student name must not be empty")
	}
	return nil
}

func (s Student) ID() StudentID {
	return s.id
}

func (s Student) Name() string {
	return s.name
}
