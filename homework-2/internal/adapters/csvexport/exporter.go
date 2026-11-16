package csvexport

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"pos/internal/application"
)

type Exporter struct{}

func (Exporter) Export(report application.FinalScoreReport, path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("export path must not be empty")
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create CSV file: %w", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	header := []string{"course_id", "student_id", "name", "mode", "score", "status"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("write CSV header: %w", err)
	}
	for _, row := range report.Rows {
		score := ""
		status := "incomplete"
		if row.HasScore {
			score = strconv.FormatFloat(row.Score, 'f', 2, 64)
			status = "calculated"
		}
		record := []string{string(report.CourseID), string(row.StudentID), row.Name, string(report.Mode), score, status}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("write CSV row: %w", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush CSV data: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close CSV file: %w", err)
	}
	return nil
}

var _ application.FinalScoreExporter = Exporter{}
