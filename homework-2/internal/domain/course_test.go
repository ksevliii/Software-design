package domain

import (
	"errors"
	"testing"
)

func TestCourse_AssignGrade(t *testing.T) {
	tests := []struct {
		name        string
		studentID   StudentID
		criterionID CriterionID
		score       int
		wantErr     error
	}{
		{
			name:        "valid grade",
			studentID:   "s1",
			criterionID: "hw",
			score:       7,
		},
		{
			name:        "zero grade",
			studentID:   "s1",
			criterionID: "hw",
			score:       0,
		},
		{
			name:        "maximum grade",
			studentID:   "s1",
			criterionID: "hw",
			score:       10,
		},
		{
			name:        "negative grade",
			studentID:   "s1",
			criterionID: "hw",
			score:       -1,
			wantErr:     ErrScoreOutOfRange,
		},
		{
			name:        "grade above maximum",
			studentID:   "s1",
			criterionID: "hw",
			score:       11,
			wantErr:     ErrScoreOutOfRange,
		},
		{
			name:        "unknown student",
			studentID:   "unknown",
			criterionID: "hw",
			score:       7,
			wantErr:     ErrStudentNotFound,
		},
		{
			name:        "unknown criterion",
			studentID:   "s1",
			criterionID: "unknown",
			score:       7,
			wantErr:     ErrCriterionNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			course := newTestCourse(t)
			err := course.AssignGrade(tt.studentID, tt.criterionID, tt.score)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("AssignGrade() error = %v, want %v", err, tt.wantErr)
			}
			grade, exists := course.FindGrade(tt.studentID, tt.criterionID)
			if tt.wantErr != nil {
				if exists {
					t.Error("grade was saved after failed assignment")
				}
				if err := course.SetRubricator(course.Rubricator()); err != nil {
					t.Errorf("rubricator locked after failed assignment: %v", err)
				}
				return
			}
			if !exists {
				t.Fatal("assigned grade not found")
			}
			if grade.StudentID() != tt.studentID {
				t.Errorf("StudentID() = %q, want %q", grade.StudentID(), tt.studentID)
			}
			if grade.CriterionID() != tt.criterionID {
				t.Errorf("CriterionID() = %q, want %q", grade.CriterionID(), tt.criterionID)
			}
			if grade.Score() != tt.score {
				t.Errorf("Score() = %d, want %d", grade.Score(), tt.score)
			}
		})
	}
}

func TestCourse_AssignGradeRejectsDuplicate(t *testing.T) {
	course := newTestCourse(t)
	if err := course.AssignGrade("s1", "hw", 7); err != nil {
		t.Fatalf("make grade: %v", err)
	}
	err := course.AssignGrade("s1", "hw", 9)
	if !errors.Is(err, ErrGradeExists) {
		t.Fatalf("AssignGrade() error = %v, want %v", err, ErrGradeExists)
	}
	grade, exists := course.FindGrade("s1", "hw")
	if !exists || grade.Score() != 7 {
		t.Fatalf("original grade changed: got (%+v, %v)", grade, exists)
	}
}

func TestCourse_CorrectGrade(t *testing.T) {
	tests := []struct {
		name      string
		score     int
		wantScore int
		wantErr   error
	}{
		{name: "valid correction", score: 9, wantScore: 9},
		{name: "correct to zero", score: 0, wantScore: 0},
		{
			name:      "invalid correction preserves old grade",
			score:     11,
			wantScore: 7,
			wantErr:   ErrScoreOutOfRange,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			course := newTestCourse(t)
			if err := course.AssignGrade("s1", "hw", 7); err != nil {
				t.Fatalf("make grade: %v", err)
			}
			err := course.CorrectGrade("s1", "hw", tt.score)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("CorrectGrade() error = %v, want %v", err, tt.wantErr)
			}
			grade, exists := course.FindGrade("s1", "hw")
			if !exists {
				t.Fatal("grade disappeared")
			}
			if grade.Score() != tt.wantScore {
				t.Errorf("Score() = %d, want %d", grade.Score(), tt.wantScore)
			}
		})
	}
}

func TestCourse_CorrectGradeRejectsMissingGrade(t *testing.T) {
	course := newTestCourse(t)
	err := course.CorrectGrade("s1", "hw", 8)
	if !errors.Is(err, ErrGradeNotFound) {
		t.Fatalf("CorrectGrade() error = %v, want %v", err, ErrGradeNotFound)
	}
	if _, exists := course.FindGrade("s1", "hw"); exists {
		t.Error("correction created a missing grade")
	}
}

func TestCourse_CancelGrade(t *testing.T) {
	course := newTestCourse(t)
	if err := course.AssignGrade("s1", "hw", 7); err != nil {
		t.Fatalf("make grade: %v", err)
	}
	if err := course.CancelGrade("s1", "hw"); err != nil {
		t.Fatalf("CancelGrade() error = %v", err)
	}
	if _, exists := course.FindGrade("s1", "hw"); exists {
		t.Error("cancelled grade still exists")
	}
	err := course.CancelGrade("s1", "hw")
	if !errors.Is(err, ErrGradeNotFound) {
		t.Errorf("second CancelGrade() error = %v, want %v", err, ErrGradeNotFound)
	}
	err = course.SetRubricator(course.Rubricator())
	if !errors.Is(err, ErrRubricLocked) {
		t.Errorf("SetRubricator() error = %v, want %v", err, ErrRubricLocked)
	}
}

func TestCourse_SetRubricatorRejectsEmpty(t *testing.T) {
	course := newTestCourse(t)
	if err := course.SetRubricator(Rubricator{}); err == nil {
		t.Fatal("SetRubricator() expected an error for empty rubricator")
	}
	if _, exists := course.Rubricator().Find("hw"); !exists {
		t.Error("original rubricator changed after failed replacement")
	}
}

func TestCourse_SetRubricatorRejectsChangeAfterGrade(t *testing.T) {
	course := newTestCourse(t)
	if err := course.AssignGrade("s1", "hw", 7); err != nil {
		t.Fatalf("make grade: %v", err)
	}
	err := course.SetRubricator(course.Rubricator())
	if !errors.Is(err, ErrRubricLocked) {
		t.Fatalf("SetRubricator() error = %v, want %v", err, ErrRubricLocked)
	}
}

func newTestCourse(t *testing.T) *Course {
	t.Helper()
	course, err := NewCourse("course1")
	if err != nil {
		t.Fatalf("create course: %v", err)
	}
	student, err := NewStudent("s1", "Анна")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	if err := course.AddStudent(student); err != nil {
		t.Fatalf("add student: %v", err)
	}
	criterion, err := NewCriterion("hw", "Домашняя работа", 10, 1)
	if err != nil {
		t.Fatalf("create criterion: %v", err)
	}
	builder := NewRubricatorBuilder()
	if err := builder.Add(criterion); err != nil {
		t.Fatalf("add criterion: %v", err)
	}
	rubricator, err := builder.Build()
	if err != nil {
		t.Fatalf("build rubricator: %v", err)
	}
	if err := course.SetRubricator(rubricator); err != nil {
		t.Fatalf("set rubricator: %v", err)
	}
	return course
}
