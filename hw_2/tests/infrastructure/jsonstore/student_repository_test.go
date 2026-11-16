package jsonstore

import (
	"path/filepath"
	"testing"
	"time"

	"hw_2/internal/domain"
	"hw_2/internal/infrastructure/jsonstore"
)

func createTestStudentRepository(
	t *testing.T,
) (*jsonstore.JSONStudentRepository, *jsonstore.JSONRubricRepository) {
	t.Helper()

	directory := t.TempDir()

	rubricRepository := jsonstore.NewJSONRubricRepository(
		filepath.Join(directory, "rubric.json"),
	)

	rubric := createTestRubric(t)

	if err := rubricRepository.Save(rubric); err != nil {
		t.Fatalf("save test rubric: %v", err)
	}

	studentRepository := jsonstore.NewJSONStudentRepository(
		filepath.Join(directory, "students.json"),
		rubricRepository,
	)

	return studentRepository, rubricRepository
}

func createTestStudent(t *testing.T) domain.Student {
	t.Helper()

	student, err := domain.NewStudent(
		"student-1",
		"Ляляля",
		"Жужужу",
		"Ивановна",
		"БПИ-241",
		"ФКН",
	)
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	score, err := domain.NewScoreValue(8)
	if err != nil {
		t.Fatalf("create score: %v", err)
	}

	maxScore, err := domain.NewMaxScore(10)
	if err != nil {
		t.Fatalf("create max score: %v", err)
	}

	if err := student.SetEvaluation(
		"criterion-1",
		score,
		maxScore,
		"Хорошая работа",
	); err != nil {
		t.Fatalf("set evaluation: %v", err)
	}

	return student
}

func TestJSONStudentRepositorySaveAndGet(t *testing.T) {
	repository, _ := createTestStudentRepository(t)

	expected := createTestStudent(t)

	if err := repository.Save(expected); err != nil {
		t.Fatalf("save student: %v", err)
	}

	actual, err := repository.GetByID(expected.ID())
	if err != nil {
		t.Fatalf("get student: %v", err)
	}

	if actual.ID() != expected.ID() {
		t.Errorf("expected ID %q, got %q", expected.ID(), actual.ID())
	}

	if actual.Name() != expected.Name() {
		t.Errorf("expected name %q, got %q", expected.Name(), actual.Name())
	}

	if actual.Surname() != expected.Surname() {
		t.Errorf("expected surname %q, got %q", expected.Surname(), actual.Surname())
	}

	evaluation, exists := actual.GetEvaluation("criterion-1")
	if !exists {
		t.Fatal("expected evaluation to be restored")
	}

	if evaluation.Score().Value() != 8 {
		t.Errorf("expected score 8, got %d", evaluation.Score().Value())
	}

	if evaluation.Comment() != "Хорошая работа" {
		t.Errorf(
			"expected comment %q, got %q",
			"Хорошая работа",
			evaluation.Comment(),
		)
	}

	if evaluation.UpdatedAt().IsZero() {
		t.Error("expected non-zero evaluation update time")
	}
}

func TestJSONStudentRepositoryGetAll(t *testing.T) {
	repository, _ := createTestStudentRepository(t)

	firstStudent := createTestStudent(t)

	secondStudent, err := domain.NewStudent(
		"student-2",
		"Улюлю",
		"Жиши",
		"Пиши",
		"БПИ-241",
		"ФКН",
	)
	if err != nil {
		t.Fatalf("create second student: %v", err)
	}

	if err := repository.Save(firstStudent); err != nil {
		t.Fatalf("save first student: %v", err)
	}

	if err := repository.Save(secondStudent); err != nil {
		t.Fatalf("save second student: %v", err)
	}

	students, err := repository.GetAll()
	if err != nil {
		t.Fatalf("get all students: %v", err)
	}

	if len(students) != 2 {
		t.Fatalf("expected 2 students, got %d", len(students))
	}
}

func TestJSONStudentRepositoryExistsByID(t *testing.T) {
	repository, _ := createTestStudentRepository(t)

	student := createTestStudent(t)

	if err := repository.Save(student); err != nil {
		t.Fatalf("save student: %v", err)
	}

	exists, err := repository.ExistsByID(student.ID())
	if err != nil {
		t.Fatalf("check existing student: %v", err)
	}

	if !exists {
		t.Error("expected student to exist")
	}

	exists, err = repository.ExistsByID("missing-student")
	if err != nil {
		t.Fatalf("check missing student: %v", err)
	}

	if exists {
		t.Error("expected missing student not to exist")
	}
}

func TestJSONStudentRepositoryExistsEvaluationForCriterion(t *testing.T) {
	repository, _ := createTestStudentRepository(t)

	student := createTestStudent(t)

	if err := repository.Save(student); err != nil {
		t.Fatalf("save student: %v", err)
	}

	exists, err := repository.ExistsEvaluationForCriterion("criterion-1")
	if err != nil {
		t.Fatalf("check existing evaluation: %v", err)
	}

	if !exists {
		t.Error("expected evaluation to exist")
	}

	exists, err = repository.ExistsEvaluationForCriterion("missing-criterion")
	if err != nil {
		t.Fatalf("check missing evaluation: %v", err)
	}

	if exists {
		t.Error("expected evaluation not to exist")
	}
}

func TestJSONStudentRepositoryPreservesEvaluationTimestamp(t *testing.T) {
	repository, rubricRepository := createTestStudentRepository(t)

	student, err := domain.NewStudent(
		"student-3",
		"Устал",
		"Усталов",
		"Усталович",
		"БПИ-241",
		"ФКН",
	)
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	score, err := domain.NewScoreValue(9)
	if err != nil {
		t.Fatalf("create score: %v", err)
	}

	maxScore, err := domain.NewMaxScore(10)
	if err != nil {
		t.Fatalf("create max score: %v", err)
	}

	updatedAt := time.Date(
		2026, time.November, 16,
		12, 30, 0, 0,
		time.UTC,
	)

	if err := student.RestoreEvaluation(
		"criterion-1",
		score,
		maxScore,
		"Отлично",
		updatedAt,
	); err != nil {
		t.Fatalf("restore evaluation: %v", err)
	}

	if err := repository.Save(student); err != nil {
		t.Fatalf("save student: %v", err)
	}

	_ = rubricRepository

	actual, err := repository.GetByID(student.ID())
	if err != nil {
		t.Fatalf("get student: %v", err)
	}

	evaluation, exists := actual.GetEvaluation("criterion-1")
	if !exists {
		t.Fatal("expected evaluation to exist")
	}

	if !evaluation.UpdatedAt().Equal(updatedAt) {
		t.Errorf(
			"expected timestamp %v, got %v",
			updatedAt,
			evaluation.UpdatedAt(),
		)
	}
}
