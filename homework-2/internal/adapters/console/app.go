package console

import (
	"errors"
	"fmt"
	"io"

	"pos/internal/application"
	"pos/internal/domain"
)

const mainMenu = `
1. Создать курс
2. Добавить студента
3. Подготовить рубрикатор
4. Выставить оценку
5. Исправить оценку
6. Отменить оценку
7. Показать студентов
8. Показать критерии
9. Показать оценки студента
10. Рассчитать итог студента
11. Экспортировать итоги в CSV
12. Показать курсы
0. Выход`

type App struct {
	service *application.Service
	input   *Input
	output  io.Writer
}

func NewApp(service *application.Service, reader io.Reader, output io.Writer) (*App, error) {
	if service == nil {
		return nil, errors.New("application service must not be nil")
	}
	if reader == nil || output == nil {
		return nil, errors.New("console input and output must not be nil")
	}
	return &App{service: service, input: NewInput(reader, output), output: output}, nil
}

func (a *App) Run() error {
	for {
		if _, err := fmt.Fprintln(a.output, mainMenu); err != nil {
			return fmt.Errorf("write menu: %w", err)
		}
		option, err := a.input.ReadInt("> ")
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		switch option {
		case 0:
			return nil
		case 1:
			err = a.createCourse()
		case 2:
			err = a.addStudent()
		case 3:
			err = a.prepareRubricator()
		case 4:
			err = a.assignGrade()
		case 5:
			err = a.correctGrade()
		case 6:
			err = a.cancelGrade()
		case 7:
			err = a.showStudents()
		case 8:
			err = a.showCriteria()
		case 9:
			err = a.showGrades()
		case 10:
			err = a.calculateFinalScore()
		case 11:
			err = a.exportFinalScores()
		case 12:
			err = a.showCourses()
		default:
			_, err = fmt.Fprintln(a.output, "Неизвестный пункт меню.")
			if err != nil {
				return err
			}
			continue
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		if errors.Is(err, ErrInputRead) {
			return err
		}
		if err != nil {
			if _, writeErr := fmt.Fprintln(a.output, "Ошибка:", err); writeErr != nil {
				return writeErr
			}
			continue
		}
	}
}

func (a *App) readExistingCourseID() (domain.CourseID, error) {
	value, err := a.input.ReadLine("ID курса: ")
	if err != nil {
		return "", err
	}
	id := domain.CourseID(value)
	if _, err := a.service.Criteria(id); err != nil {
		return "", err
	}
	return id, nil
}

func (a *App) createCourse() error {
	id, err := a.input.ReadLine("ID курса: ")
	if err != nil {
		return err
	}
	if err := a.service.CreateCourse(domain.CourseID(id)); err != nil {
		return err
	}
	_, err = fmt.Fprintln(a.output, "Курс создан.")
	return err
}

func (a *App) addStudent() error {
	courseID, err := a.readExistingCourseID()
	if err != nil {
		return err
	}
	value, err := a.input.ReadLine("ID студента: ")
	if err != nil {
		return err
	}
	studentID := domain.StudentID(value)
	if studentID == "" {
		return fmt.Errorf("student ID must not be empty")
	}
	students, err := a.service.Students(courseID)
	if err != nil {
		return err
	}
	for _, student := range students {
		if student.ID() == studentID {
			return fmt.Errorf(
				"%w: %q", domain.ErrStudentExists, studentID,
			)
		}
	}
	name, err := a.input.ReadLine("Имя студента: ")
	if err != nil {
		return err
	}
	if err := a.service.AddStudent(courseID, studentID, name); err != nil {
		return err
	}
	_, err = fmt.Fprintln(a.output, "Студент добавлен.")
	return err
}
