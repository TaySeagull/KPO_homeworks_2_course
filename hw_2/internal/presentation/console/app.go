package console

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"hw_2/internal/application"
)

// App represents the console interface of the POS application
type App struct {
	facade *application.ApplicationFacade
	reader *bufio.Reader
	output io.Writer
}

// New creates a console application
func New(
	facade *application.ApplicationFacade,
	input io.Reader,
	output io.Writer,
) *App {
	return &App{
		facade: facade,
		reader: bufio.NewReader(input),
		output: output,
	}
}

// Run starts the main application menu
func (app *App) Run() {
	for {
		app.printMenu()

		choice, err := app.readLine("Выберите действие: ")
		if err != nil {
			fmt.Fprintln(app.output, "Ошибка чтения ввода:", err)
			return
		}

		switch choice {
		case "1":
			app.addStudent()
		case "2":
			app.listStudents()
		case "3":
			app.addCriterion()
		case "4":
			app.listCriteria()
		case "5":
			app.updateCriterion()
		case "6":
			app.removeCriterion()
		case "7":
			app.setEvaluation()
		case "8":
			app.listEvaluations()
		case "9":
			app.calculateScore()
		case "10":
			app.exportScore()
		case "0":
			fmt.Fprintln(app.output, "Работа программы завершена")
			return
		default:
			fmt.Fprintln(app.output, "Неизвестная команда")
		}

		fmt.Fprintln(app.output)
	}
}

// printMenu displays the available application operations
func (app *App) printMenu() {
	fmt.Fprintln(app.output, "======================================")
	fmt.Fprintln(app.output, "       ПОС — Преподавательская        ")
	fmt.Fprintln(app.output, "           Оценка Студентов           ")
	fmt.Fprintln(app.output, "======================================")
	fmt.Fprintln(app.output, "1. Добавить студента")
	fmt.Fprintln(app.output, "2. Показать всех студентов")
	fmt.Fprintln(app.output, "3. Добавить критерий")
	fmt.Fprintln(app.output, "4. Показать рубрикатор")
	fmt.Fprintln(app.output, "5. Изменить критерий")
	fmt.Fprintln(app.output, "6. Удалить критерий")
	fmt.Fprintln(app.output, "7. Выставить оценку")
	fmt.Fprintln(app.output, "8. Показать оценки студента")
	fmt.Fprintln(app.output, "9. Рассчитать итоговый балл")
	fmt.Fprintln(app.output, "10. Экспортировать итог в CSV")
	fmt.Fprintln(app.output, "0. Выход")
	fmt.Fprintln(app.output, "======================================")
}

// addStudent handles student creation
func (app *App) addStudent() {
	fmt.Fprintln(app.output, "\nДобавление студента")

	id, err := app.readLine("ID: ")
	if err != nil {
		app.printError(err)
		return
	}

	name, err := app.readLine("Имя: ")
	if err != nil {
		app.printError(err)
		return
	}

	surname, err := app.readLine("Фамилия: ")
	if err != nil {
		app.printError(err)
		return
	}

	patronymic, err := app.readLine("Отчество: ")
	if err != nil {
		app.printError(err)
		return
	}

	group, err := app.readLine("Группа: ")
	if err != nil {
		app.printError(err)
		return
	}

	faculty, err := app.readLine("Факультет: ")
	if err != nil {
		app.printError(err)
		return
	}

	if err := app.facade.AddStudent(
		id,
		name,
		surname,
		patronymic,
		group,
		faculty,
	); err != nil {
		app.printError(err)
		return
	}

	fmt.Fprintln(app.output, "Студент успешно добавлен")
}

// listStudents displays all registered students
func (app *App) listStudents() {
	students, err := app.facade.GetAllStudents()
	if err != nil {
		app.printError(err)
		return
	}

	fmt.Fprintln(app.output, "\nСписок студентов")

	if len(students) == 0 {
		fmt.Fprintln(app.output, "Студентов пока нет")
		return
	}

	for _, student := range students {
		fmt.Fprintf(
			app.output,
			"ID: %s | %s %s %s | Группа: %s | Факультет: %s\n",
			student.ID(),
			student.Surname(),
			student.Name(),
			student.Patronymic(),
			student.Group(),
			student.Faculty(),
		)
	}
}

// addCriterion handles criterion creation
func (app *App) addCriterion() {
	fmt.Fprintln(app.output, "\nДобавление критерия")

	id, err := app.readLine("ID критерия: ")
	if err != nil {
		app.printError(err)
		return
	}

	name, err := app.readLine("Название: ")
	if err != nil {
		app.printError(err)
		return
	}

	maxScore, err := app.readInt("Максимальный балл: ")
	if err != nil {
		app.printError(err)
		return
	}

	weight, err := app.readInt("Вес критерия (в процентах): ")
	if err != nil {
		app.printError(err)
		return
	}

	if err := app.facade.AddCriterion(id, name, maxScore, weight); err != nil {
		app.printError(err)
		return
	}

	fmt.Fprintln(app.output, "Критерий успешно добавлен")
}

// listCriteria displays all rubric criteria
func (app *App) listCriteria() {
	rubric, err := app.facade.GetRubric()
	if err != nil {
		app.printError(err)
		return
	}

	criteria := rubric.GetCriteria()

	fmt.Fprintln(app.output, "\nРубрикатор")

	if len(criteria) == 0 {
		fmt.Fprintln(app.output, "Критериев пока нет")
		return
	}

	for _, criterion := range criteria {
		fmt.Fprintf(
			app.output,
			"ID: %s | %s | Максимум: %d | Вес: %d%%\n",
			criterion.ID(),
			criterion.Name(),
			criterion.MaxScore().Value(),
			criterion.Weight().Value(),
		)
	}
}

// updateCriterion handles criterion modification
func (app *App) updateCriterion() {
	fmt.Fprintln(app.output, "\nИзменение критерия")

	id, err := app.readLine("ID критерия: ")
	if err != nil {
		app.printError(err)
		return
	}

	name, err := app.readLine("Новое название: ")
	if err != nil {
		app.printError(err)
		return
	}

	maxScore, err := app.readInt("Новый максимальный балл: ")
	if err != nil {
		app.printError(err)
		return
	}

	weight, err := app.readInt("Новый вес (в процентах): ")
	if err != nil {
		app.printError(err)
		return
	}

	if err := app.facade.UpdateCriterion(id, name, maxScore, weight); err != nil {
		app.printError(err)
		return
	}

	fmt.Fprintln(app.output, "Критерий успешно изменён")
}

// removeCriterion handles criterion deletion
func (app *App) removeCriterion() {
	fmt.Fprintln(app.output, "\nУдаление критерия")

	id, err := app.readLine("ID критерия: ")
	if err != nil {
		app.printError(err)
		return
	}

	if err := app.facade.RemoveCriterion(id); err != nil {
		app.printError(err)
		return
	}

	fmt.Fprintln(app.output, "Критерий успешно удалён")
}

// setEvaluation handles creating or updating a student's evaluation
func (app *App) setEvaluation() {
	fmt.Fprintln(app.output, "\nВыставление оценки")

	studentID, err := app.readLine("ID студента: ")
	if err != nil {
		app.printError(err)
		return
	}

	criterionID, err := app.readLine("ID критерия: ")
	if err != nil {
		app.printError(err)
		return
	}

	score, err := app.readInt("Полученный балл: ")
	if err != nil {
		app.printError(err)
		return
	}

	comment, err := app.readLine("Комментарий: ")
	if err != nil {
		app.printError(err)
		return
	}

	if err := app.facade.SetEvaluation(
		studentID,
		criterionID,
		score,
		comment,
	); err != nil {
		app.printError(err)
		return
	}

	fmt.Fprintln(app.output, "Оценка успешно сохранена")
}

// listEvaluations displays all evaluations of a student
func (app *App) listEvaluations() {
	fmt.Fprintln(app.output, "\nОценки студента")

	studentID, err := app.readLine("ID студента: ")
	if err != nil {
		app.printError(err)
		return
	}

	student, err := app.facade.GetStudent(studentID)
	if err != nil {
		app.printError(err)
		return
	}

	evaluations := student.GetEvaluations()

	if len(evaluations) == 0 {
		fmt.Fprintln(app.output, "У студента пока нет оценок")
		return
	}

	for _, evaluation := range evaluations {
		fmt.Fprintf(
			app.output,
			"Критерий: %s | Балл: %d | Комментарий: %s | Обновлено: %s\n",
			evaluation.CriterionID(),
			evaluation.Score().Value(),
			evaluation.Comment(),
			evaluation.UpdatedAt().Format("02.01.2006 15:04:05"),
		)
	}
}

// calculateScore calculates and displays a student's final score
func (app *App) calculateScore() {
	fmt.Fprintln(app.output, "\nРасчёт итогового балла")

	studentID, err := app.readLine("ID студента: ")
	if err != nil {
		app.printError(err)
		return
	}

	score, err := app.facade.CalculateStudentScore(studentID)
	if err != nil {
		app.printError(err)
		return
	}

	fmt.Fprintf(
		app.output,
		"Итоговый балл студента: %d из 10\n",
		score.Value(),
	)
}

// exportScore exports a student's final score to a CSV file
func (app *App) exportScore() {
	fmt.Fprintln(app.output, "\nЭкспорт итогового балла")

	studentID, err := app.readLine("ID студента: ")
	if err != nil {
		app.printError(err)
		return
	}

	destination, err := app.readLine("Путь к CSV-файлу: ")
	if err != nil {
		app.printError(err)
		return
	}

	if err := app.facade.ExportStudentScore(studentID, destination); err != nil {
		app.printError(err)
		return
	}

	fmt.Fprintf(app.output, "Результат экспортирован в %s\n", destination)
}

// readLine reads a line from the console
func (app *App) readLine(prompt string) (string, error) {
	fmt.Fprint(app.output, prompt)

	value, err := app.reader.ReadString('\n')
	if err != nil && len(value) == 0 {
		return "", err
	}

	return strings.TrimSpace(value), nil
}

// readInt reads an integer from the console
func (app *App) readInt(prompt string) (int, error) {
	value, err := app.readLine(prompt)
	if err != nil {
		return 0, err
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("ожидалось целое число")
	}

	return number, nil
}

// printError displays an application error
func (app *App) printError(err error) {
	fmt.Fprintln(app.output, "Ошибка:", err)
}
