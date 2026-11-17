package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"pos/internal/adapters/csvexport"
	"pos/internal/adapters/jsonstore"
	"pos/internal/application"
	"pos/internal/di"
	"pos/internal/domain"
)

func main() {
	if err := runDemo(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка демонстрации:", err)
		os.Exit(1)
	}
}

func runDemo() error {
	fmt.Println("~~~ ДЕМОНСТРАЦИЯ СЦЕНАРИЕВ ПОС ~~~")
	const demoDirectory = "demo-data"
	if err := os.MkdirAll(demoDirectory, 0755); err != nil {
		return fmt.Errorf("create demo directory: %w", err)
	}
	runDirectory, err := os.MkdirTemp(demoDirectory, "run-")
	if err != nil {
		return fmt.Errorf("create demo run directory: %w", err)
	}
	fmt.Println("Папка результатов:", runDirectory)
	repository, err := jsonstore.NewRepository(runDirectory)
	if err != nil {
		return err
	}
	service, err := di.BuildService(repository, csvexport.Exporter{}, application.ModeStrict)
	if err != nil {
		return err
	}
	const courseID domain.CourseID = "demo"
	fmt.Println("\n--- Сценарий 1: создание курса и студентов ---")
	if err := service.CreateCourse(courseID); err != nil {
		return err
	}
	fmt.Println("Создан курс:", courseID)
	students := []struct {
		id   domain.StudentID
		name string
	}{
		{id: "s1", name: "Ксюша"},
		{id: "s2", name: "Даша"},
	}

	for _, student := range students {
		if err := service.AddStudent(courseID, student.id, student.name); err != nil {
			return err
		}
		fmt.Printf("Добавлен студент: %s - %s\n", student.id, student.name)
	}
	fmt.Println("\n--- Сценарий 2: подготовка рубрикатора ---")
	builder := domain.NewRubricatorBuilder()
	criteria := []struct {
		id       domain.CriterionID
		name     string
		maxScore int
		weight   int
	}{
		{id: "hw", name: "Домашняя работа", maxScore: 20, weight: 40},
		{id: "exam", name: "Экзамен", maxScore: 10, weight: 60},
	}
	for _, data := range criteria {
		criterion, err := domain.NewCriterion(data.id, data.name, data.maxScore, data.weight)
		if err != nil {
			return err
		}
		if err := builder.Add(criterion); err != nil {
			return err
		}
		fmt.Printf("%s [%s]: максимум %d, вес %d\n", data.name, data.id, data.maxScore, data.weight)
	}
	rubricator, err := builder.Build()
	if err != nil {
		return err
	}
	if err := service.SetRubricator(courseID, rubricator); err != nil {
		return err
	}
	fmt.Println("Рубрикатор подтверждён и сохранён.")
	fmt.Println("\n--- Сценарий 3: выставление оценок ---")
	grades := []struct {
		studentID   domain.StudentID
		criterionID domain.CriterionID
		score       int
	}{
		{studentID: "s1", criterionID: "hw", score: 16},
		{studentID: "s1", criterionID: "exam", score: 7},
		{studentID: "s2", criterionID: "hw", score: 20},
		{studentID: "s2", criterionID: "exam", score: 9},
	}
	for _, grade := range grades {
		if err := service.AssignGrade(courseID, grade.studentID, grade.criterionID, grade.score); err != nil {
			return err
		}
		fmt.Printf("Студент %s, критерий %s: %d\n", grade.studentID, grade.criterionID, grade.score)
	}
	fmt.Println("\n--- Сценарий 4: отклонение некорректных операций ---")
	if err := expectError(
		"Повторное выставление оценки",
		service.AssignGrade(courseID, "s1", "hw", 10),
		domain.ErrGradeExists); err != nil {
		return err
	}
	if err := expectError(
		"Исправление на 12 баллов из 10",
		service.CorrectGrade(courseID, "s1", "exam", 12),
		domain.ErrScoreOutOfRange); err != nil {
		return err
	}
	if err := expectError(
		"Оценка по отсутствующему критерию",
		service.AssignGrade(courseID, "s1", "unknown", 5),
		domain.ErrCriterionNotFound); err != nil {
		return err
	}
	fmt.Println("\n--- Сценарий 5: расчёт итогов ---")
	for _, student := range students {
		score, err := service.CalculateFinalScore(courseID, student.id)
		if err != nil {
			return err
		}
		fmt.Printf("%s: %.2f из 10\n", student.name, score)
	}
	fmt.Println("\n--- Сценарий 6: исправление и отмена оценки ---")
	fmt.Println("Исправляем Ксюше оценку за экзамен на 10...")
	if err := service.CorrectGrade(courseID, "s1", "exam", 10); err != nil {
		return err
	}
	score, err := service.CalculateFinalScore(courseID, "s1")
	if err != nil {
		return err
	}
	fmt.Printf("Ксюша после исправления экзамена: %.2f\n", score)
	fmt.Println("Отменяем Дашину оценку за экзамен...")
	if err := service.CancelGrade(courseID, "s2", "exam"); err != nil {
		return err
	}
	fmt.Println("Оценка Даши за экзамен отменена.")
	_, calculationErr := service.CalculateFinalScore(courseID, "s2")
	if err := expectError(
		"Строгий расчёт при пропущенной оценке",
		calculationErr,
		domain.ErrIncompleteGrades); err != nil {
		return err
	}
	fmt.Println("\n--- Сценарий 7: экспорт в строгом режиме ---")
	strictPath := filepath.Join(runDirectory, "results-strict.csv")
	if err := service.ExportFinalScores(courseID, strictPath); err != nil {
		return err
	}
	fmt.Println("Отчёт:", strictPath)
	fmt.Println("У Даши итог пустой, статус incomplete.")
	fmt.Println("\n--- Сценарий 8: пропуски как нули ---")
	zeroService, err := di.BuildService(repository, csvexport.Exporter{}, application.ModeMissingAsZero)
	if err != nil {
		return err
	}
	score, err = zeroService.CalculateFinalScore(courseID, "s2")
	if err != nil {
		return err
	}
	fmt.Printf("Даша: %.2f\n", score)
	zeroPath := filepath.Join(runDirectory, "results-zero.csv")
	if err := zeroService.ExportFinalScores(courseID, zeroPath); err != nil {
		return err
	}
	fmt.Println("Отчёт:", zeroPath)
	fmt.Println("\n~~~ ДЕМОНСТРАЦИЯ ЗАВЕРШЕНА ~~~")
	fmt.Println("JSON и CSV сохранены в:", runDirectory)
	return nil
}

func expectError(message string, got, want error) error {
	if !errors.Is(got, want) {
		return fmt.Errorf("%s: получена ошибка %v, ожидалась %v", message, got, want)
	}
	fmt.Printf("%s - отклонено: %v\n", message, got)
	return nil
}
