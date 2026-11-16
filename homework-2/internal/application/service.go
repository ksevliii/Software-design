package application

import (
	"errors"
	"fmt"

	"pos/internal/domain"
)

var ErrCourseExists = errors.New("course already exists")

type Service struct {
	repository CourseRepository
	scoring    *domain.FinalScoreService
	exporter   FinalScoreExporter
	mode       CalculationMode
}

func NewService(
	repository CourseRepository,
	scoring *domain.FinalScoreService,
	exporter FinalScoreExporter,
	mode CalculationMode) (*Service, error) {
	if repository == nil {
		return nil, errors.New("course repository must not be nil")
	}
	if scoring == nil {
		return nil, errors.New("final score service must not be nil")
	}
	if exporter == nil {
		return nil, errors.New("final score exporter must not be nil")
	}
	if !mode.Valid() {
		return nil, fmt.Errorf("unknown calculation mode: %q", mode)
	}
	return &Service{repository: repository, scoring: scoring, exporter: exporter, mode: mode}, nil
}

func (s *Service) CreateCourse(id domain.CourseID) error {
	course, err := domain.NewCourse(id)
	if err != nil {
		return err
	}
	_, err = s.repository.Load(course.ID())
	if err == nil {
		return fmt.Errorf("%w: %q", ErrCourseExists, course.ID())
	} else if !errors.Is(err, ErrCourseNotFound) {
		return fmt.Errorf("check course existence: %w", err)
	}
	if err := s.repository.Save(course); err != nil {
		return fmt.Errorf("save new course: %w", err)
	}
	return nil
}

func (s *Service) CheckRubricatorEditable(courseID domain.CourseID) error {
	course, err := s.repository.Load(courseID)
	if err != nil {
		return fmt.Errorf("load course: %w", err)
	}
	if course.RubricLocked() {
		return domain.ErrRubricLocked
	}
	return nil
}

func (s *Service) AssignGrade(
	courseID domain.CourseID,
	studentID domain.StudentID,
	criterionID domain.CriterionID,
	score int) error {
	course, err := s.repository.Load(courseID)
	if err != nil {
		return fmt.Errorf("load course: %w", err)
	}
	if err := course.AssignGrade(studentID, criterionID, score); err != nil {
		return fmt.Errorf("assign grade: %w", err)
	}
	if err := s.repository.Save(course); err != nil {
		return fmt.Errorf("save course: %w", err)
	}
	return nil
}

func (s *Service) CalculateFinalScore(courseID domain.CourseID, studentID domain.StudentID) (float64, error) {
	course, err := s.repository.Load(courseID)
	if err != nil {
		return 0, fmt.Errorf("load course: %w", err)
	}
	score, err := s.scoring.Calculate(course, studentID)
	if err != nil {
		return 0, fmt.Errorf("calculate final score: %w", err)
	}
	return score, nil
}

func (s *Service) AddStudent(courseID domain.CourseID, studentID domain.StudentID, name string) error {
	student, err := domain.NewStudent(studentID, name)
	if err != nil {
		return fmt.Errorf("create student: %w", err)
	}
	course, err := s.repository.Load(courseID)
	if err != nil {
		return fmt.Errorf("load course: %w", err)
	}
	if err := course.AddStudent(student); err != nil {
		return fmt.Errorf("add student: %w", err)
	}
	if err := s.repository.Save(course); err != nil {
		return fmt.Errorf("save course: %w", err)
	}
	return nil
}

func (s *Service) SetRubricator(courseID domain.CourseID, rubricator domain.Rubricator) error {
	course, err := s.repository.Load(courseID)
	if err != nil {
		return fmt.Errorf("load course: %w", err)
	}
	if err := course.SetRubricator(rubricator); err != nil {
		return fmt.Errorf("set rubricator: %w", err)
	}
	if err := s.repository.Save(course); err != nil {
		return fmt.Errorf("save course: %w", err)
	}
	return nil
}

func (s *Service) CorrectGrade(
	courseID domain.CourseID,
	studentID domain.StudentID,
	criterionID domain.CriterionID,
	score int) error {
	course, err := s.repository.Load(courseID)
	if err != nil {
		return fmt.Errorf("load course: %w", err)
	}
	if err := course.CorrectGrade(studentID, criterionID, score); err != nil {
		return fmt.Errorf("correct grade: %w", err)
	}
	if err := s.repository.Save(course); err != nil {
		return fmt.Errorf("save course: %w", err)
	}
	return nil
}

func (s *Service) CancelGrade(courseID domain.CourseID, studentID domain.StudentID, criterionID domain.CriterionID) error {
	course, err := s.repository.Load(courseID)
	if err != nil {
		return fmt.Errorf("load course: %w", err)
	}
	if err := course.CancelGrade(studentID, criterionID); err != nil {
		return fmt.Errorf("cancel grade: %w", err)
	}
	if err := s.repository.Save(course); err != nil {
		return fmt.Errorf("save course: %w", err)
	}
	return nil
}

func (s *Service) Students(courseID domain.CourseID) ([]domain.Student, error) {
	course, err := s.repository.Load(courseID)
	if err != nil {
		return nil, fmt.Errorf("load course: %w", err)
	}
	return course.Students(), nil
}

func (s *Service) Criteria(courseID domain.CourseID) ([]domain.Criterion, error) {
	course, err := s.repository.Load(courseID)
	if err != nil {
		return nil, fmt.Errorf("load course: %w", err)
	}
	return course.Rubricator().Criteria(), nil
}

func (s *Service) Grades(courseID domain.CourseID, studentID domain.StudentID) ([]domain.Grade, error) {
	course, err := s.repository.Load(courseID)
	if err != nil {
		return nil, fmt.Errorf("load course: %w", err)
	}
	grades, err := course.Grades(studentID)
	if err != nil {
		return nil, fmt.Errorf("get grades: %w", err)
	}
	return grades, nil
}

func (s *Service) CourseIDs() ([]domain.CourseID, error) {
	ids, err := s.repository.ListIDs()
	if err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	return ids, nil
}
