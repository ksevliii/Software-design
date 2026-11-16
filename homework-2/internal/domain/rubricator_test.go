package domain

import "testing"

func TestRubricatorBuilder_Add(t *testing.T) {
	first := mustCriterion(t, "hw", "Домашняя работа")
	second := mustCriterion(t, "hw", "Другое название")
	tests := []struct {
		name      string
		initial   []Criterion
		criterion Criterion
		wantErr   bool
	}{
		{
			name:      "valid criterion",
			criterion: first,
		},
		{
			name:      "zero criterion",
			criterion: Criterion{},
			wantErr:   true,
		},
		{
			name:      "duplicate ID",
			initial:   []Criterion{first},
			criterion: second,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewRubricatorBuilder()

			for _, criterion := range tt.initial {
				if err := builder.Add(criterion); err != nil {
					t.Fatalf("prepare builder: %v", err)
				}
			}
			err := builder.Add(tt.criterion)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Add() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if len(tt.initial) == 0 {
					if _, err := builder.Build(); err == nil {
						t.Fatal("Build() succeeded after rejecting only one criterion")
					}
					return
				}
				rubricator := mustBuild(t, builder)
				got := rubricator.Criteria()

				if len(got) != len(tt.initial) {
					t.Fatalf("criterion count = %d, want %d", len(got), len(tt.initial))
				}
				for _, want := range tt.initial {
					actual, ok := rubricator.Find(want.ID())
					if !ok {
						t.Fatalf("criterion %q disappeared after failed Add()", want.ID())
					}
					if actual != want {
						t.Errorf("criterion = %+v, want %+v", actual, want)
					}
				}
				return
			}
			rubricator := mustBuild(t, builder)
			got, ok := rubricator.Find(tt.criterion.ID())
			if !ok {
				t.Fatal("added criterion not found")
			}
			if got != tt.criterion {
				t.Errorf("criterion = %+v, want %+v", got, tt.criterion)
			}
		})
	}
}

func TestRubricatorBuilder_ZeroValueAdd(t *testing.T) {
	var builder RubricatorBuilder
	criterion := mustCriterion(t, "hw", "Домашняя работа")
	if err := builder.Add(criterion); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	rubricator := mustBuild(t, &builder)
	got, ok := rubricator.Find(criterion.ID())
	if !ok || got != criterion {
		t.Errorf("Find() = (%+v, %v), want (%+v, true)", got, ok, criterion)
	}
}

func TestRubricatorBuilder_BuildEmpty(t *testing.T) {
	tests := []struct {
		name    string
		builder *RubricatorBuilder
	}{
		{
			name:    "created by constructor",
			builder: NewRubricatorBuilder(),
		},
		{
			name:    "zero value",
			builder: &RubricatorBuilder{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.builder.Build(); err == nil {
				t.Fatal("Build() expected error for empty builder")
			}
		})
	}
}

func TestRubricator_Find(t *testing.T) {
	criterion := mustCriterion(t, "hw", "Домашняя работа")
	builder := NewRubricatorBuilder()
	mustAdd(t, builder, criterion)
	rubricator := mustBuild(t, builder)
	tests := []struct {
		name   string
		id     CriterionID
		wantOK bool
	}{
		{name: "existing ID", id: "hw", wantOK: true},
		{name: "unknown ID", id: "exam"},
		{name: "empty ID", id: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := rubricator.Find(tt.id)
			if ok != tt.wantOK {
				t.Fatalf("Find() found = %v, want %v", ok, tt.wantOK)
			}
			if tt.wantOK && got != criterion {
				t.Errorf("Find() = %+v, want %+v", got, criterion)
			}
		})
	}
}

func TestRubricator_CriteriaSorted(t *testing.T) {
	builder := NewRubricatorBuilder()
	for _, id := range []CriterionID{"b", "a", "c"} {
		mustAdd(t, builder, mustCriterion(t, id, string(id)))
	}
	rubricator := mustBuild(t, builder)
	got := rubricator.Criteria()
	want := []CriterionID{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("criterion count = %d, want %d", len(got), len(want))
	}
	for i, id := range want {
		if got[i].ID() != id {
			t.Errorf("criteria[%d].ID() = %q, want %q", i, got[i].ID(), id)
		}
	}
}

func TestRubricator_CriteriaReturnsCopy(t *testing.T) {
	builder := NewRubricatorBuilder()
	original := mustCriterion(t, "hw", "Домашняя работа")
	mustAdd(t, builder, original)
	rubricator := mustBuild(t, builder)
	criteria := rubricator.Criteria()
	if len(criteria) != 1 {
		t.Fatalf("criterion count = %d, want 1", len(criteria))
	}
	criteria[0] = mustCriterion(t, "exam", "Экзамен")
	got := rubricator.Criteria()
	if len(got) != 1 {
		t.Fatalf("criterion count after slice modification = %d, want 1", len(got))
	}
	if got[0] != original {
		t.Errorf("criterion changed: got %+v, want %+v", got[0], original)
	}
}

func TestRubricatorBuilder_BuildReturnsIndependentRubricator(t *testing.T) {
	builder := NewRubricatorBuilder()
	mustAdd(t, builder, mustCriterion(t, "hw", "Домашняя работа"))
	first := mustBuild(t, builder)
	mustAdd(t, builder, mustCriterion(t, "exam", "Экзамен"))
	second := mustBuild(t, builder)
	if _, ok := first.Find("exam"); ok {
		t.Error("first rubricator changed after adding to builder")
	}
	if len(first.Criteria()) != 1 {
		t.Error("first rubricator must contain exactly one criterion")
	}
	if _, ok := second.Find("exam"); !ok {
		t.Error("second rubricator must contain the new criterion")
	}
	if len(second.Criteria()) != 2 {
		t.Error("second rubricator must contain exactly two criteria")
	}
}

func TestRubricator_ZeroValue(t *testing.T) {
	var rubricator Rubricator
	if _, ok := rubricator.Find("hw"); ok {
		t.Error("Find() must return false for zero rubricator")
	}
	if got := rubricator.Criteria(); len(got) != 0 {
		t.Errorf("Criteria() length = %d, want 0", len(got))
	}
}

func mustCriterion(t *testing.T, id CriterionID, name string) Criterion {
	t.Helper()
	criterion, err := NewCriterion(id, name, 10, 1)
	if err != nil {
		t.Fatalf("create criterion: %v", err)
	}
	return criterion
}

func mustAdd(t *testing.T, builder *RubricatorBuilder, criterion Criterion) {
	t.Helper()
	if err := builder.Add(criterion); err != nil {
		t.Fatalf("add criterion: %v", err)
	}
}

func mustBuild(t *testing.T, builder *RubricatorBuilder) Rubricator {
	t.Helper()
	rubricator, err := builder.Build()
	if err != nil {
		t.Fatalf("build rubricator: %v", err)
	}
	return rubricator
}
