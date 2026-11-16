package domain

import (
	"errors"
	"math"
	"testing"
)

func TestFinalScoreStrategies(t *testing.T) {
	const epsilon = 1e-8
	tests := []struct {
		name       string
		strategy   FinalScoreStrategy
		components []ScoreComponent
		want       float64
		wantErr    error
		anyErr     bool
	}{
		{
			name:     "weighted score with different maximum",
			strategy: StrictStrategy{},
			components: []ScoreComponent{
				{MaxScore: 20, Weight: 40, Score: 16, HasGrade: true},
				{MaxScore: 10, Weight: 60, Score: 7, HasGrade: true},
			},
			want: 7.4,
		},
		{
			name:     "relative weights",
			strategy: StrictStrategy{},
			components: []ScoreComponent{
				{MaxScore: 20, Weight: 2, Score: 16, HasGrade: true},
				{MaxScore: 10, Weight: 3, Score: 7, HasGrade: true},
			},
			want: 7.4,
		},
		{
			name:     "maximum score",
			strategy: StrictStrategy{},
			components: []ScoreComponent{
				{MaxScore: 20, Weight: 40, Score: 20, HasGrade: true},
				{MaxScore: 10, Weight: 60, Score: 10, HasGrade: true},
			},
			want: 10,
		},
		{
			name:     "assigned zero is not missing",
			strategy: StrictStrategy{},
			components: []ScoreComponent{
				{MaxScore: 10, Weight: 1, Score: 0, HasGrade: true},
			},
			want: 0,
		},
		{
			name:     "strict rejects missing grade",
			strategy: StrictStrategy{},
			components: []ScoreComponent{
				{MaxScore: 20, Weight: 40, Score: 16, HasGrade: true},
				{MaxScore: 10, Weight: 60, HasGrade: false},
			},
			wantErr: ErrIncompleteGrades,
		},
		{
			name:     "missing grade contributes zero",
			strategy: MissingAsZeroStrategy{},
			components: []ScoreComponent{
				{MaxScore: 20, Weight: 40, Score: 16, HasGrade: true},
				{MaxScore: 10, Weight: 60, HasGrade: false},
			},
			want: 3.2,
		},
		{
			name:     "all grades missing",
			strategy: MissingAsZeroStrategy{},
			components: []ScoreComponent{
				{MaxScore: 10, Weight: 1},
				{MaxScore: 20, Weight: 2},
			},
			want: 0,
		},
		{
			name:     "zero policy also calculates complete grades",
			strategy: MissingAsZeroStrategy{},
			components: []ScoreComponent{
				{MaxScore: 10, Weight: 1, Score: 8, HasGrade: true},
			},
			want: 8,
		},
		{
			name:     "empty criteria",
			strategy: StrictStrategy{},
			anyErr:   true,
		},
		{
			name:     "zero maximum",
			strategy: StrictStrategy{},
			components: []ScoreComponent{
				{MaxScore: 0, Weight: 1, HasGrade: true},
			},
			anyErr: true,
		},
		{
			name:     "negative maximum",
			strategy: StrictStrategy{},
			components: []ScoreComponent{
				{MaxScore: -1, Weight: 1, HasGrade: true},
			},
			anyErr: true,
		},
		{
			name:     "zero weight",
			strategy: StrictStrategy{},
			components: []ScoreComponent{
				{MaxScore: 10, Weight: 0, HasGrade: true},
			},
			anyErr: true,
		},
		{
			name:     "negative weight",
			strategy: StrictStrategy{},
			components: []ScoreComponent{
				{MaxScore: 10, Weight: -1, HasGrade: true},
			},
			anyErr: true,
		},
		{
			name:     "negative score",
			strategy: StrictStrategy{},
			components: []ScoreComponent{
				{MaxScore: 10, Weight: 1, Score: -1, HasGrade: true},
			},
			wantErr: ErrScoreOutOfRange,
		},
		{
			name:     "score above maximum",
			strategy: StrictStrategy{},
			components: []ScoreComponent{
				{MaxScore: 10, Weight: 1, Score: 11, HasGrade: true},
			},
			wantErr: ErrScoreOutOfRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.strategy.Calculate(tt.components)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Calculate() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if tt.anyErr {
				if err == nil {
					t.Fatal("Calculate() expected an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("Calculate() unexpected error: %v", err)
			}
			if math.IsNaN(got) || math.Abs(got-tt.want) > epsilon {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
