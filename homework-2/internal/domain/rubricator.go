package domain

import (
	"errors"
	"maps"
	"sort"
)

type Rubricator struct {
	criteria map[CriterionID]Criterion
}

func (r Rubricator) Find(id CriterionID) (Criterion, bool) {
	criterion, ok := r.criteria[id]
	return criterion, ok
}

func (r Rubricator) Criteria() []Criterion {
	criteria := make([]Criterion, 0, len(r.criteria))
	for _, criterion := range r.criteria {
		criteria = append(criteria, criterion)
	}
	sort.Slice(criteria, func(idx1, idx2 int) bool {
		return criteria[idx1].ID() < criteria[idx2].ID()
	})
	return criteria
}

type RubricatorBuilder struct {
	criteria map[CriterionID]Criterion
}

func NewRubricatorBuilder() *RubricatorBuilder {
	return &RubricatorBuilder{criteria: make(map[CriterionID]Criterion)}
}

func (b *RubricatorBuilder) Add(criterion Criterion) error {
	if err := criterion.validate(); err != nil {
		return err
	}
	if _, exists := b.criteria[criterion.id]; exists {
		return errors.New("criterion with this ID already exists")
	}
	if b.criteria == nil {
		b.criteria = make(map[CriterionID]Criterion)
	}
	b.criteria[criterion.id] = criterion
	return nil
}

func (b *RubricatorBuilder) Build() (Rubricator, error) {
	if len(b.criteria) == 0 {
		return Rubricator{}, errors.New("cannot create rubricator with empty criteria")
	}
	copyCriteria := make(map[CriterionID]Criterion, len(b.criteria))
	maps.Copy(copyCriteria, b.criteria)
	return Rubricator{criteria: copyCriteria}, nil
}
