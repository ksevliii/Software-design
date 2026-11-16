package domain

import "testing"

func TestNewStudent(t *testing.T) {
	tests := []struct {
		name        string
		id          StudentID
		studentName string
		wantID      StudentID
		wantName    string
		wantErr     bool
	}{
		{
			name:        "valid student",
			id:          "s1",
			studentName: "Ксюша",
			wantID:      "s1",
			wantName:    "Ксюша",
		},
		{
			name:        "trim spaces",
			id:          " s1 ",
			studentName: " Ксюшаша ",
			wantID:      "s1",
			wantName:    "Ксюшаша",
		},
		{
			name:        "empty ID",
			studentName: "Ксюша",
			wantErr:     true,
		},
		{
			name:        "blank ID",
			id:          " \t ",
			studentName: "Ксюшик",
			wantErr:     true,
		},
		{
			name:    "empty name",
			id:      "s1",
			wantErr: true,
		},
		{
			name:        "blank name",
			id:          "s1",
			studentName: " \t ",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStudent(tt.id, tt.studentName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewStudent() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got.ID() != tt.wantID {
				t.Errorf("ID() = %q, want %q", got.ID(), tt.wantID)
			}
			if got.Name() != tt.wantName {
				t.Errorf("Name() = %q, want %q", got.Name(), tt.wantName)
			}
		})
	}
}
