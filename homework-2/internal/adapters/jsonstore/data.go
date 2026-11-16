package jsonstore

type courseData struct {
	ID           string          `json:"id"`
	Students     []studentData   `json:"students"`
	Criteria     []criterionData `json:"criteria"`
	Grades       []gradeData     `json:"grades"`
	RubricLocked bool            `json:"rubric_locked"`
}

type studentData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type criterionData struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	MaxScore int    `json:"max_score"`
	Weight   int    `json:"weight"`
}

type gradeData struct {
	StudentID   string `json:"student_id"`
	CriterionID string `json:"criterion_id"`
	Score       int    `json:"score"`
}
