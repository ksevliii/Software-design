package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type CourseID string

type Course struct {
	id           CourseID
	students     map[StudentID]Student
	rubricator   Rubricator
	grades       map[StudentID]map[CriterionID]Grade
	rubricLocked bool
}

func NewCourse(id CourseID) (*Course, error) {
	id = CourseID(strings.TrimSpace(string(id)))
	if id == "" {
		return nil, errors.New("course ID must not be empty")
	}
	for _, char := range id {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return nil, errors.New("course ID must contain only Latin letters, digits, underscores or hyphens")
		}
	}
	return &Course{
		id:       id,
		students: make(map[StudentID]Student),
		grades:   make(map[StudentID]map[CriterionID]Grade),
	}, nil
}

func (c *Course) ID() CourseID {
	return c.id
}

func (c *Course) Rubricator() Rubricator {
	return c.rubricator
}

func (c *Course) RubricLocked() bool {
	return c.rubricLocked
}

func (c *Course) LockRubricator() error {
	if len(c.rubricator.criteria) == 0 {
		return errors.New("cannot lock an empty rubricator")
	}
	c.rubricLocked = true
	return nil
}

func (c *Course) AddStudent(student Student) error {
	if err := student.validate(); err != nil {
		return fmt.Errorf("add student: %w", err)
	}
	if _, exists := c.students[student.ID()]; exists {
		return fmt.Errorf("%w: %q", ErrStudentExists, student.ID())
	}
	c.students[student.ID()] = student
	return nil
}

func (c *Course) SetRubricator(rubricator Rubricator) error {
	if c.rubricLocked {
		return ErrRubricLocked
	}
	if len(rubricator.criteria) == 0 {
		return errors.New("rubricator must not be empty")
	}
	c.rubricator = rubricator
	return nil
}

func (c *Course) validateGradeParameters(studentID StudentID, criterionID CriterionID) (Criterion, error) {
	if _, exists := c.students[studentID]; !exists {
		return Criterion{}, fmt.Errorf("%w: %q", ErrStudentNotFound, studentID)
	}
	criterion, exists := c.rubricator.Find(criterionID)
	if !exists {
		return Criterion{}, fmt.Errorf("%w: %q", ErrCriterionNotFound, criterionID)
	}
	return criterion, nil
}

func (c *Course) AssignGrade(studentID StudentID, criterionID CriterionID, score int) error {
	criterion, err := c.validateGradeParameters(studentID, criterionID)
	if err != nil {
		return err
	}
	if _, exists := c.grades[studentID][criterionID]; exists {
		return fmt.Errorf("%w: student %q, criterion %q", ErrGradeExists, studentID, criterionID)
	}
	grade, err := newGrade(studentID, criterion, score)
	if err != nil {
		return err
	}
	if c.grades[studentID] == nil {
		c.grades[studentID] = make(map[CriterionID]Grade)
	}
	c.grades[studentID][criterionID] = grade
	c.rubricLocked = true
	return nil
}

func (c *Course) CorrectGrade(studentID StudentID, criterionID CriterionID, score int) error {
	criterion, err := c.validateGradeParameters(studentID, criterionID)
	if err != nil {
		return err
	}
	if _, exists := c.grades[studentID][criterionID]; !exists {
		return fmt.Errorf("%w: student %q, criterion %q", ErrGradeNotFound, studentID, criterionID)
	}
	grade, err := newGrade(studentID, criterion, score)
	if err != nil {
		return err
	}
	c.grades[studentID][criterionID] = grade
	return nil
}

func (c *Course) CancelGrade(studentID StudentID, criterionID CriterionID) error {
	if _, err := c.validateGradeParameters(studentID, criterionID); err != nil {
		return err
	}
	if _, exists := c.grades[studentID][criterionID]; !exists {
		return fmt.Errorf("%w: student %q, criterion %q", ErrGradeNotFound, studentID, criterionID)
	}
	delete(c.grades[studentID], criterionID)
	return nil
}

func (c *Course) FindGrade(studentID StudentID, criterionID CriterionID) (Grade, bool) {
	grade, exists := c.grades[studentID][criterionID]
	return grade, exists
}

func (c *Course) Grades(studentID StudentID) ([]Grade, error) {
	if _, exists := c.students[studentID]; !exists {
		return nil, fmt.Errorf("%w: %q", ErrStudentNotFound, studentID)
	}
	studentGrades := c.grades[studentID]
	grades := make([]Grade, 0, len(studentGrades))
	for _, grade := range studentGrades {
		grades = append(grades, grade)
	}
	sort.Slice(grades, func(idx1, idx2 int) bool {
		return grades[idx1].CriterionID() < grades[idx2].CriterionID()
	})
	return grades, nil
}

func (c *Course) Students() []Student {
	students := make([]Student, 0, len(c.students))
	for _, student := range c.students {
		students = append(students, student)
	}
	sort.Slice(students, func(idx1, idx2 int) bool {
		return students[idx1].ID() < students[idx2].ID()
	})
	return students
}
