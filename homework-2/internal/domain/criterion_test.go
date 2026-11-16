package domain

import "testing"

func TestNewCriterion(t *testing.T) {
	tests := []struct {
		name              string
		id                CriterionID
		criterionName     string
		maxScore          int
		weight            int
		wantID            CriterionID
		wantCriterionName string
		wantErr           bool
	}{
		{
			name:              "valid criterion",
			id:                "hw",
			criterionName:     "Домашняя работа",
			maxScore:          20,
			weight:            40,
			wantID:            "hw",
			wantCriterionName: "Домашняя работа",
		},
		{
			name:              "trim spaces",
			id:                " hw ",
			criterionName:     " Домашняя работа ",
			maxScore:          20,
			weight:            40,
			wantID:            "hw",
			wantCriterionName: "Домашняя работа",
		},
		{
			name:              "minimum positive values",
			id:                "hw",
			criterionName:     "Домашняя работа",
			maxScore:          1,
			weight:            1,
			wantID:            "hw",
			wantCriterionName: "Домашняя работа",
		},
		{
			name:              "weight above 100",
			id:                "hw",
			criterionName:     "Домашняя работа",
			maxScore:          20,
			weight:            200,
			wantID:            "hw",
			wantCriterionName: "Домашняя работа",
		},
		{
			name:          "empty ID",
			criterionName: "Домашняя работа",
			maxScore:      20,
			weight:        40,
			wantErr:       true,
		},
		{
			name:          "blank ID",
			id:            " \t ",
			criterionName: "Домашняя работа",
			maxScore:      20,
			weight:        40,
			wantErr:       true,
		},
		{
			name:     "empty name",
			id:       "hw",
			maxScore: 20,
			weight:   40,
			wantErr:  true,
		},
		{
			name:          "blank name",
			id:            "hw",
			criterionName: " \t ",
			maxScore:      20,
			weight:        40,
			wantErr:       true,
		},
		{
			name:          "zero maximum",
			id:            "hw",
			criterionName: "Домашняя работа",
			maxScore:      0,
			weight:        40,
			wantErr:       true,
		},
		{
			name:          "negative maximum",
			id:            "hw",
			criterionName: "Домашняя работа",
			maxScore:      -1,
			weight:        40,
			wantErr:       true,
		},
		{
			name:          "zero weight",
			id:            "hw",
			criterionName: "Домашняя работа",
			maxScore:      20,
			weight:        0,
			wantErr:       true,
		},
		{
			name:          "negative weight",
			id:            "hw",
			criterionName: "Домашняя работа",
			maxScore:      20,
			weight:        -1,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCriterion(tt.id, tt.criterionName, tt.maxScore, tt.weight)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewCriterion() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got.ID() != tt.wantID {
				t.Errorf("ID() = %q, want %q", got.ID(), tt.wantID)
			}
			if got.Name() != tt.wantCriterionName {
				t.Errorf("Name() = %q, want %q", got.Name(), tt.wantCriterionName)
			}
			if got.MaxScore() != tt.maxScore {
				t.Errorf("MaxScore() = %d, want %d", got.MaxScore(), tt.maxScore)
			}
			if got.Weight() != tt.weight {
				t.Errorf("Weight() = %d, want %d", got.Weight(), tt.weight)
			}
		})
	}
}
