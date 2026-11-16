package jsonstore

import (
	"errors"
	"fmt"

	"pos/internal/domain"
)

func courseToData(course *domain.Course) (courseData, error) {
	if course == nil {
		return courseData{}, errors.New("course must not be nil")
	}
	data := courseData{
		ID:           string(course.ID()),
		Students:     make([]studentData, 0),
		Criteria:     make([]criterionData, 0),
		Grades:       make([]gradeData, 0),
		RubricLocked: course.RubricLocked(),
	}
	for _, student := range course.Students() {
		data.Students = append(data.Students, studentData{ID: string(student.ID()), Name: student.Name()})
		grades, err := course.Grades(student.ID())
		if err != nil {
			return courseData{}, fmt.Errorf("get grades for student %q: %w", student.ID(), err)
		}
		for _, grade := range grades {
			data.Grades = append(data.Grades, gradeData{
				StudentID:   string(grade.StudentID()),
				CriterionID: string(grade.CriterionID()),
				Score:       grade.Score(),
			})
		}
	}
	for _, criterion := range course.Rubricator().Criteria() {
		data.Criteria = append(data.Criteria, criterionData{
			ID:       string(criterion.ID()),
			Name:     criterion.Name(),
			MaxScore: criterion.MaxScore(),
			Weight:   criterion.Weight(),
		})
	}
	return data, nil
}

func courseFromData(data courseData) (*domain.Course, error) {
	course, err := domain.NewCourse(domain.CourseID(data.ID))
	if err != nil {
		return nil, fmt.Errorf("restore course: %w", err)
	}
	for _, studentData := range data.Students {
		student, err := domain.NewStudent(domain.StudentID(studentData.ID), studentData.Name)
		if err != nil {
			return nil, fmt.Errorf("restore student %q: %w", studentData.ID, err)
		}
		if err := course.AddStudent(student); err != nil {
			return nil, fmt.Errorf("restore students: %w", err)
		}
	}
	if len(data.Criteria) > 0 {
		builder := domain.NewRubricatorBuilder()
		for _, criterionData := range data.Criteria {
			criterion, err := domain.NewCriterion(
				domain.CriterionID(criterionData.ID),
				criterionData.Name,
				criterionData.MaxScore,
				criterionData.Weight,
			)
			if err != nil {
				return nil, fmt.Errorf("restore criterion %q: %w", criterionData.ID, err)
			}
			if err := builder.Add(criterion); err != nil {
				return nil, fmt.Errorf("restore criteria: %w", err)
			}
		}
		rubricator, err := builder.Build()
		if err != nil {
			return nil, fmt.Errorf("restore rubricator: %w", err)
		}
		if err := course.SetRubricator(rubricator); err != nil {
			return nil, fmt.Errorf("set restored rubricator: %w", err)
		}
	}
	if len(data.Grades) > 0 && !data.RubricLocked {
		return nil, errors.New("saved course contains grades but its rubricator is not locked")
	}
	for _, gradeData := range data.Grades {
		if err := course.AssignGrade(
			domain.StudentID(gradeData.StudentID),
			domain.CriterionID(gradeData.CriterionID),
			gradeData.Score,
		); err != nil {
			return nil, fmt.Errorf("restore grade for student %q, criterion %q: %w",
				gradeData.StudentID, gradeData.CriterionID, err)
		}
	}
	if data.RubricLocked {
		if err := course.LockRubricator(); err != nil {
			return nil, fmt.Errorf("restore rubricator lock: %w", err)
		}
	}
	return course, nil
}
