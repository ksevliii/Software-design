package console

import (
	"fmt"

	"pos/internal/domain"
)

func (a *App) readGradeTarget() (domain.CourseID, domain.StudentID, domain.CriterionID, error) {
	courseID, err := a.readExistingCourseID()
	if err != nil {
		return "", "", "", err
	}
	students, err := a.service.Students(courseID)
	if err != nil {
		return "", "", "", err
	}
	studentInput, err := a.input.ReadLine("ID студента: ")
	if err != nil {
		return "", "", "", err
	}
	studentID := domain.StudentID(studentInput)
	studentExists := false
	for _, student := range students {
		if student.ID() == studentID {
			studentExists = true
			break
		}
	}
	if !studentExists {
		return "", "", "", fmt.Errorf("%w: %q", domain.ErrStudentNotFound, studentID)
	}
	criteria, err := a.service.Criteria(courseID)
	if err != nil {
		return "", "", "", err
	}
	if len(criteria) == 0 {
		return "", "", "", fmt.Errorf("рубрикатор курса ещё не подготовлен")
	}
	criterionInput, err := a.input.ReadLine("ID критерия: ")
	if err != nil {
		return "", "", "", err
	}
	criterionID := domain.CriterionID(criterionInput)
	for _, criterion := range criteria {
		if criterion.ID() == criterionID {
			return courseID, studentID, criterionID, nil
		}
	}
	return "", "", "", fmt.Errorf("%w: %q", domain.ErrCriterionNotFound, criterionID)
}

func (a *App) assignGrade() error {
	courseID, studentID, criterionID, err := a.readGradeTarget()
	if err != nil {
		return err
	}
	score, err := a.input.ReadInt("Оценка: ")
	if err != nil {
		return err
	}
	if err := a.service.AssignGrade(courseID, studentID, criterionID, score); err != nil {
		return err
	}
	_, err = fmt.Fprintln(a.output, "Оценка выставлена.")
	return err
}

func (a *App) correctGrade() error {
	courseID, studentID, criterionID, err := a.readGradeTarget()
	if err != nil {
		return err
	}
	score, err := a.input.ReadInt("Исправленная оценка: ")
	if err != nil {
		return err
	}
	if err := a.service.CorrectGrade(courseID, studentID, criterionID, score); err != nil {
		return err
	}
	_, err = fmt.Fprintln(a.output, "Оценка исправлена.")
	return err
}

func (a *App) cancelGrade() error {
	courseID, studentID, criterionID, err := a.readGradeTarget()
	if err != nil {
		return err
	}
	if err := a.service.CancelGrade(courseID, studentID, criterionID); err != nil {
		return err
	}
	_, err = fmt.Fprintln(a.output, "Оценка отменена.")
	return err
}
