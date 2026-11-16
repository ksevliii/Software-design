package domain

import (
	"errors"
	"strings"
)

type CriterionID string

type Criterion struct {
	id       CriterionID
	name     string
	maxScore int
	weight   int
}

func NewCriterion(id CriterionID, name string, maxScore, weight int) (Criterion, error) {
	c := Criterion{
		id:       CriterionID(strings.TrimSpace(string(id))),
		name:     strings.TrimSpace(name),
		maxScore: maxScore,
		weight:   weight,
	}
	if err := c.validate(); err != nil {
		return Criterion{}, err
	}
	return c, nil
}

func (c Criterion) validate() error {
	if c.id == "" {
		return errors.New("criterion ID must not be empty")
	}
	if c.name == "" {
		return errors.New("criterion name must not be empty")
	}
	if c.maxScore <= 0 {
		return errors.New("criterion maximum score must be positive")
	}
	if c.weight <= 0 {
		return errors.New("criterion weight must be positive")
	}
	return nil
}

func (c Criterion) ID() CriterionID {
	return c.id
}

func (c Criterion) Name() string {
	return c.name
}

func (c Criterion) MaxScore() int {
	return c.maxScore
}

func (c Criterion) Weight() int {
	return c.weight
}
