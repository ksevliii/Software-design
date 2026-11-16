package domain

import (
	"errors"
	"fmt"
	"strings"
)

type Grade struct {
	studentID   StudentID
	criterionID CriterionID
	score       int
}

func newGrade(studentID StudentID, criterion Criterion, score int) (Grade, error) {
	if strings.TrimSpace(string(studentID)) == "" {
		return Grade{}, errors.New("student ID must not be empty")
	}
	if err := criterion.validate(); err != nil {
		return Grade{}, err
	}
	if score < 0 || score > criterion.MaxScore() {
		return Grade{}, fmt.Errorf(
			"%w: got %d, allowed range [0, %d]",
			ErrScoreOutOfRange, score, criterion.MaxScore(),
		)
	}
	return Grade{studentID: studentID, criterionID: criterion.ID(), score: score}, nil
}

func (g Grade) StudentID() StudentID {
	return g.studentID
}

func (g Grade) CriterionID() CriterionID {
	return g.criterionID
}

func (g Grade) Score() int {
	return g.score
}
