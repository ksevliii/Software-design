package application

import (
	"errors"
	"fmt"
	"pos/internal/domain"
)

type FinalScoreRow struct {
	StudentID domain.StudentID
	Name      string
	Score     float64
	HasScore  bool
}

type FinalScoreReport struct {
	CourseID domain.CourseID
	Mode     CalculationMode
	Rows     []FinalScoreRow
}

type FinalScoreExporter interface {
	Export(report FinalScoreReport, path string) error
}

func (s *Service) ExportFinalScores(courseID domain.CourseID, path string) error {
	course, err := s.repository.Load(courseID)
	if err != nil {
		return fmt.Errorf("load course: %w", err)
	}
	if len(course.Rubricator().Criteria()) == 0 {
		return errors.New("cannot export final scores without criteria")
	}
	report := FinalScoreReport{CourseID: course.ID(), Mode: s.mode, Rows: make([]FinalScoreRow, 0)}
	for _, student := range course.Students() {
		score, err := s.scoring.Calculate(course, student.ID())
		row := FinalScoreRow{StudentID: student.ID(), Name: student.Name()}
		switch {
		case err == nil:
			row.Score = score
			row.HasScore = true
		case errors.Is(err, domain.ErrIncompleteGrades):
			row.HasScore = false
		default:
			return fmt.Errorf("calculate score for student %q: %w", student.ID(), err)
		}
		report.Rows = append(report.Rows, row)
	}
	if err := s.exporter.Export(report, path); err != nil {
		return fmt.Errorf("export final scores: %w", err)
	}
	return nil
}
