package console_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"hw_2/internal/container"
	"hw_2/internal/presentation/console"
)

// newTestApp creates a console application with a real DI container
func newTestApp(t *testing.T, input string) (*console.App, *bytes.Buffer) {
	t.Helper()

	directory := t.TempDir()

	appContainer, err := container.New(container.Config{
		StudentsFile: filepath.Join(directory, "students.json"),
		RubricFile:   filepath.Join(directory, "rubric.json"),
	})
	if err != nil {
		t.Fatalf("failed to create container: %v", err)
	}

	t.Cleanup(func() {
		appContainer.Shutdown()
	})

	var output bytes.Buffer

	app := console.New(
		appContainer.Facade,
		strings.NewReader(input),
		&output,
	)

	return app, &output
}

func TestAppCompleteStudentWorkflow(t *testing.T) {
	exportPath := filepath.Join(t.TempDir(), "result.csv")

	input := strings.Join([]string{
		"1",
		"s001",
		"Траляля",
		"Жужужу",
		"Ляляля",
		"KPO-01",
		"Факультет информатики",

		"3",
		"project",
		"Проект",
		"10",
		"50",

		"3",
		"practice",
		"Практика",
		"10",
		"50",

		"7",
		"s001",
		"project",
		"8",
		"Хорошая работа",

		"7",
		"s001",
		"practice",
		"10",
		"Всё выполнено отлично",

		"2",
		"4",

		"8",
		"s001",

		"9",
		"s001",

		"10",
		"s001",
		exportPath,

		"0",
	}, "\n") + "\n"

	app, output := newTestApp(t, input)
	app.Run()

	result := output.String()

	expectedMessages := []string{
		"Студент успешно добавлен",
		"Критерий успешно добавлен",
		"Оценка успешно сохранена",
		"Список студентов",
		"Траляля",
		"Жужужу",
		"Проект",
		"Практика",
		"Хорошая работа",
		"Всё выполнено отлично",
		"Итоговый балл студента: 9 из 10",
		"Результат экспортирован",
		"Работа программы завершена",
	}

	for _, message := range expectedMessages {
		if !strings.Contains(result, message) {
			t.Errorf("expected output to contain %q", message)
		}
	}

	if _, err := os.Stat(exportPath); err != nil {
		t.Fatalf("expected CSV file to be created: %v", err)
	}

	csvContent, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("failed to read CSV file: %v", err)
	}

	if len(csvContent) == 0 {
		t.Error("expected CSV file to contain data")
	}
}

func TestAppCriterionManagement(t *testing.T) {
	input := strings.Join([]string{
		"3",
		"project",
		"Проект",
		"10",
		"100",

		"5",
		"project",
		"Улучшенный проект",
		"20",
		"100",

		"4",

		"6",
		"project",

		"4",

		"0",
	}, "\n") + "\n"

	app, output := newTestApp(t, input)
	app.Run()

	result := output.String()

	expectedMessages := []string{
		"Критерий успешно добавлен",
		"Критерий успешно изменён",
		"Улучшенный проект",
		"Критерий успешно удалён",
		"Критериев пока нет",
	}

	for _, message := range expectedMessages {
		if !strings.Contains(result, message) {
			t.Errorf("expected output to contain %q", message)
		}
	}
}

func TestAppInvalidCommand(t *testing.T) {
	input := "999\n0\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "Неизвестная команда") {
		t.Error("expected an error message for an unknown command")
	}

	if !strings.Contains(output.String(), "Работа программы завершена") {
		t.Error("expected application to exit normally")
	}
}

func TestAppInvalidIntegerInput(t *testing.T) {
	input := strings.Join([]string{
		"3",
		"project",
		"Проект",
		"invalid",
		"0",
	}, "\n") + "\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "ожидалось целое число") {
		t.Error("expected an error message for invalid integer input")
	}
}

func TestAppMissingStudent(t *testing.T) {
	input := "8\nunknown-student\n0\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "Ошибка:") {
		t.Error("expected an error when requesting evaluations of a nonexistent student")
	}
}

func TestAppEmptyStudentList(t *testing.T) {
	input := "2\n0\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "Студентов пока нет") {
		t.Error("expected an empty student list message")
	}
}

func TestAppEmptyCriteriaList(t *testing.T) {
	input := "4\n0\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "Критериев пока нет") {
		t.Error("expected an empty criteria list message")
	}
}

func TestAppExit(t *testing.T) {
	input := "0\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "ПОС") {
		t.Error("expected application menu to be displayed")
	}

	if !strings.Contains(output.String(), "Работа программы завершена") {
		t.Error("expected application to exit")
	}
}

func TestAppReadInputErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "error while reading menu choice",
			input: "",
		},
		{
			name:  "error while reading student ID",
			input: "1\n",
		},
		{
			name:  "error while reading student name",
			input: "1\ns001\n",
		},
		{
			name:  "error while reading student surname",
			input: "1\ns001\nТраляля\n",
		},
		{
			name:  "error while reading student patronymic",
			input: "1\ns001\nТраляля\nЖужужу\n",
		},
		{
			name:  "error while reading student group",
			input: "1\ns001\nТраляля\nЖужужу\nЛяляля\n",
		},
		{
			name:  "error while reading student faculty",
			input: "1\ns001\nТраляля\nЖужужу\nЛяляля\nKPO-01\n",
		},
		{
			name:  "error while reading criterion ID",
			input: "3\n",
		},
		{
			name:  "error while reading criterion name",
			input: "3\nproject\n",
		},
		{
			name:  "error while reading criterion weight",
			input: "3\nproject\nПроект\n10\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, output := newTestApp(t, test.input)
			app.Run()

			if !strings.Contains(output.String(), "Ошибка") {
				t.Error("expected input error to be displayed")
			}
		})
	}
}

func TestAppInvalidStudentCreation(t *testing.T) {
	input := strings.Join([]string{
		"1",
		"",
		"Траляля",
		"Жужужу",
		"Ляляля",
		"KPO-01",
		"Факультет информатики",
		"0",
	}, "\n") + "\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "Ошибка:") {
		t.Error("expected an error when creating a student with empty ID")
	}
}

func TestAppInvalidCriterionCreation(t *testing.T) {
	input := strings.Join([]string{
		"3",
		"project",
		"Проект",
		"-10",
		"100",
		"0",
	}, "\n") + "\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "Ошибка:") {
		t.Error("expected an error when creating a criterion with invalid maximum score")
	}
}

func TestAppInvalidCriterionUpdate(t *testing.T) {
	input := strings.Join([]string{
		"3",
		"project",
		"Проект",
		"10",
		"100",

		"5",
		"unknown-criterion",
		"Обновлённый проект",
		"20",
		"100",

		"0",
	}, "\n") + "\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "Ошибка:") {
		t.Error("expected an error when updating a nonexistent criterion")
	}
}

func TestAppInvalidCriterionRemoval(t *testing.T) {
	input := "6\nunknown-criterion\n0\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "Ошибка:") {
		t.Error("expected an error when removing a nonexistent criterion")
	}
}

func TestAppInvalidEvaluation(t *testing.T) {
	input := strings.Join([]string{
		"7",
		"unknown-student",
		"unknown-criterion",
		"5",
		"Комментарий",
		"0",
	}, "\n") + "\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "Ошибка:") {
		t.Error("expected an error when setting evaluation for a nonexistent student")
	}
}

func TestAppEmptyEvaluations(t *testing.T) {
	input := strings.Join([]string{
		"1",
		"s001",
		"Траляля",
		"Жужужу",
		"Ляляля",
		"KPO-01",
		"Факультет информатики",

		"8",
		"s001",
		"0",
	}, "\n") + "\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "У студента пока нет оценок") {
		t.Error("expected a message about empty evaluations")
	}
}

func TestAppInvalidScoreCalculation(t *testing.T) {
	input := "9\nunknown-student\n0\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "Ошибка:") {
		t.Error("expected an error when calculating score for a nonexistent student")
	}
}

func TestAppInvalidExport(t *testing.T) {
	input := strings.Join([]string{
		"10",
		"unknown-student",
		"result.csv",
		"0",
	}, "\n") + "\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "Ошибка:") {
		t.Error("expected an error when exporting score for a nonexistent student")
	}
}

func TestAppInvalidIntegerInCriterionWeight(t *testing.T) {
	input := strings.Join([]string{
		"3",
		"project",
		"Проект",
		"10",
		"invalid",
		"0",
	}, "\n") + "\n"

	app, output := newTestApp(t, input)
	app.Run()

	if !strings.Contains(output.String(), "ожидалось целое число") {
		t.Error("expected an integer parsing error")
	}
}
