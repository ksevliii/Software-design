package console

import (
	"fmt"

	"pos/internal/domain"
)

func (a *App) showStudents() error {
	id, err := a.input.ReadLine("ID курса: ")
	if err != nil {
		return err
	}
	students, err := a.service.Students(domain.CourseID(id))
	if err != nil {
		return err
	}
	if len(students) == 0 {
		_, err := fmt.Fprintln(a.output, "Студентов пока нет.")
		return err
	}
	for _, student := range students {
		if _, err := fmt.Fprintf(a.output, "%s - %s\n", student.ID(), student.Name()); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) showCriteria() error {
	id, err := a.input.ReadLine("ID курса: ")
	if err != nil {
		return err
	}
	criteria, err := a.service.Criteria(domain.CourseID(id))
	if err != nil {
		return err
	}
	if len(criteria) == 0 {
		_, err := fmt.Fprintln(a.output, "Рубрикатор ещё не подготовлен.")
		return err
	}
	for _, criterion := range criteria {
		if _, err := fmt.Fprintf(a.output,
			"%s - %s; максимум: %d; вес: %d\n",
			criterion.ID(),
			criterion.Name(),
			criterion.MaxScore(),
			criterion.Weight()); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) showGrades() error {
	courseID, err := a.readExistingCourseID()
	if err != nil {
		return err
	}
	studentID, err := a.input.ReadLine("ID студента: ")
	if err != nil {
		return err
	}
	grades, err := a.service.Grades(courseID, domain.StudentID(studentID))
	if err != nil {
		return err
	}
	if len(grades) == 0 {
		_, err := fmt.Fprintln(a.output, "Оценок пока нет.")
		return err
	}
	for _, grade := range grades {
		if _, err := fmt.Fprintf(a.output, "Критерий %s: %d\n", grade.CriterionID(), grade.Score()); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) calculateFinalScore() error {
	courseID, err := a.readExistingCourseID()
	if err != nil {
		return err
	}
	studentID, err := a.input.ReadLine("ID студента: ")
	if err != nil {
		return err
	}
	score, err := a.service.CalculateFinalScore(courseID, domain.StudentID(studentID))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(a.output, "Итог: %.2f из 10\n", score)
	return err
}

func (a *App) exportFinalScores() error {
	courseID, err := a.readExistingCourseID()
	if err != nil {
		return err
	}
	path, err := a.input.ReadLine("Имя CSV-файла (с путем): ")
	if err != nil {
		return err
	}
	if err := a.service.ExportFinalScores(courseID, path); err != nil {
		return err
	}
	_, err = fmt.Fprintln(a.output, "Отчёт сохранён:", path)
	return err
}

func (a *App) showCourses() error {
	ids, err := a.service.CourseIDs()
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		_, err := fmt.Fprintln(a.output, "Курсов пока нет.")
		return err
	}
	for _, id := range ids {
		if _, err := fmt.Fprintln(a.output, id); err != nil {
			return err
		}
	}
	return nil
}
