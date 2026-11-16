package console

import (
	"fmt"

	"pos/internal/domain"
)

const rubricatorMenu = `
1. Добавить критерий
2. Подтвердить рубрикатор
0. Отменить подготовку`

func (a *App) prepareRubricator() error {
	courseID, err := a.readExistingCourseID()
	if err != nil {
		return err
	}
	if err := a.service.CheckRubricatorEditable(courseID); err != nil {
		return err
	}
	criteria, err := a.service.Criteria(courseID)
	if err != nil {
		return err
	}
	builder := domain.NewRubricatorBuilder()
	existingIDs := make(map[domain.CriterionID]struct{}, len(criteria))
	for _, criterion := range criteria {
		if err := builder.Add(criterion); err != nil {
			return fmt.Errorf("prepare existing criterion: %w", err)
		}
		existingIDs[criterion.ID()] = struct{}{}
	}
	if _, err := fmt.Fprintf(a.output, "В рубрикаторе уже %d критериев. Они будут сохранены.\n", len(criteria)); err != nil {
		return err
	}

	for {
		if _, err := fmt.Fprintln(a.output, rubricatorMenu); err != nil {
			return err
		}
		option, err := a.input.ReadInt("> ")
		if err != nil {
			return err
		}
		switch option {
		case 0:
			_, err := fmt.Fprintln(a.output, "Подготовка отменена.")
			return err
		case 1:
			if err := a.addCriterionToBuilder(builder, existingIDs); err != nil {
				return err
			}
		case 2:
			rubricator, err := builder.Build()
			if err != nil {
				if _, writeErr := fmt.Fprintln(a.output, "Ошибка:", err); writeErr != nil {
					return writeErr
				}
				continue
			}
			if err := a.service.SetRubricator(courseID, rubricator); err != nil {
				return err
			}
			_, err = fmt.Fprintln(a.output, "Рубрикатор сохранён.")
			return err
		default:
			if _, err := fmt.Fprintln(a.output, "Неизвестный пункт меню."); err != nil {
				return err
			}
		}
	}
}

func (a *App) addCriterionToBuilder(builder *domain.RubricatorBuilder, existingIDs map[domain.CriterionID]struct{}) error {
	id, err := a.input.ReadLine("ID критерия: ")
	if err != nil {
		return err
	}
	if id == "" {
		_, err := fmt.Fprintln(a.output, "ID критерия не должен быть пустым.")
		return err
	}
	criterionID := domain.CriterionID(id)
	if _, exists := existingIDs[criterionID]; exists {
		_, err := fmt.Fprintln(a.output, "Критерий с таким ID уже есть.")
		return err
	}
	name, err := a.input.ReadLine("Название критерия: ")
	if err != nil {
		return err
	}
	if name == "" {
		_, err := fmt.Fprintln(a.output, "Название не должно быть пустым.")
		return err
	}
	maxScore, err := a.input.ReadInt("Максимальный балл: ")
	if err != nil {
		return err
	}
	if maxScore <= 0 {
		_, err := fmt.Fprintln(a.output, "Максимальный балл должен быть положительным.")
		return err
	}
	weight, err := a.input.ReadInt("Вес: ")
	if err != nil {
		return err
	}
	if weight <= 0 {
		_, err := fmt.Fprintln(a.output, "Вес должен быть положительным.")
		return err
	}
	criterion, err := domain.NewCriterion(criterionID, name, maxScore, weight)
	if err != nil {
		_, writeErr := fmt.Fprintln(a.output, "Ошибка:", err)
		return writeErr
	}
	if err := builder.Add(criterion); err != nil {
		_, writeErr := fmt.Fprintln(a.output, "Ошибка:", err)
		return writeErr
	}
	existingIDs[criterion.ID()] = struct{}{}
	_, err = fmt.Fprintln(a.output, "Критерий добавлен в заготовку.")
	return err
}
