package domain

import (
	"errors"
	"fmt"
)

const MaxFinalScore = 10

type ScoreComponent struct {
	MaxScore int
	Weight   int
	Score    int
	HasGrade bool
}

type FinalScoreStrategy interface {
	Calculate(components []ScoreComponent) (float64, error)
}

type StrictStrategy struct{}

func (StrictStrategy) Calculate(components []ScoreComponent) (float64, error) {
	return calculateWeighted(components, false)
}

type MissingAsZeroStrategy struct{}

func (MissingAsZeroStrategy) Calculate(components []ScoreComponent) (float64, error) {
	return calculateWeighted(components, true)
}

func calculateWeighted(components []ScoreComponent, missingAsZero bool) (float64, error) {
	if len(components) == 0 {
		return 0, errors.New("cannot calculate final score without criteria")
	}
	for i, component := range components {
		if component.MaxScore <= 0 || component.Weight <= 0 {
			return 0, fmt.Errorf("component %d: maximum score and weight must be positive", i)
		}
		if component.HasGrade && (component.Score < 0 || component.Score > component.MaxScore) {
			return 0, fmt.Errorf("component %d: %w", i, ErrScoreOutOfRange)
		}
	}
	var weightedSum, totalWeight float64
	for _, component := range components {
		if !component.HasGrade && !missingAsZero {
			return 0, ErrIncompleteGrades
		}
		weight := float64(component.Weight)
		totalWeight += weight
		if component.HasGrade {
			weightedSum += weight * float64(component.Score) / float64(component.MaxScore)
		}
	}
	return MaxFinalScore * weightedSum / totalWeight, nil
}

type FinalScoreService struct {
	strategy FinalScoreStrategy
}

func NewFinalScoreService(strategy FinalScoreStrategy) (*FinalScoreService, error) {
	if strategy == nil {
		return nil, errors.New("final score strategy must not be nil")
	}
	return &FinalScoreService{strategy: strategy}, nil
}

func (s *FinalScoreService) Calculate(course *Course, studentID StudentID) (float64, error) {
	if course == nil {
		return 0, errors.New("course must not be nil")
	}
	grades, err := course.Grades(studentID)
	if err != nil {
		return 0, err
	}
	gradesByCriterion := make(map[CriterionID]Grade, len(grades))
	for _, grade := range grades {
		gradesByCriterion[grade.CriterionID()] = grade
	}
	criteria := course.Rubricator().Criteria()
	components := make([]ScoreComponent, 0, len(criteria))
	for _, criterion := range criteria {
		grade, exists := gradesByCriterion[criterion.ID()]
		components = append(components, ScoreComponent{
			MaxScore: criterion.MaxScore(),
			Weight:   criterion.Weight(),
			Score:    grade.Score(),
			HasGrade: exists,
		})
	}
	return s.strategy.Calculate(components)
}
