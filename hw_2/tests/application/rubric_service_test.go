package application_test

import (
	"errors"
	"strings"
	"testing"

	"hw_2/internal/application"
	"hw_2/internal/domain"
)

func TestRubricServiceAddCriterion(t *testing.T) {
	rubricRepository := newInMemoryRubricRepository()
	studentRepository := newInMemoryStudentRepository()
	service := application.NewRubricService(rubricRepository, studentRepository)

	err := service.AddCriterion("criterion-1", "Экзамен", 100, 100)
	if err != nil {
		t.Fatalf("AddCriterion returned unexpected error: %v", err)
	}

	rubric, err := service.GetRubric()
	if err != nil {
		t.Fatalf("GetRubric returned unexpected error: %v", err)
	}

	criterion, exists := rubric.GetCriterion("criterion-1")
	if !exists {
		t.Fatal("criterion was not added")
	}

	if criterion.Name() != "Экзамен" {
		t.Errorf("expected criterion name Экзамен, got %s", criterion.Name())
	}

	if criterion.MaxScore().Value() != 100 {
		t.Errorf("expected maximum score 100, got %d", criterion.MaxScore().Value())
	}

	if criterion.Weight().Value() != 100 {
		t.Errorf("expected weight 100, got %d", criterion.Weight().Value())
	}
}

func TestRubricServiceAddDuplicateCriterion(t *testing.T) {
	rubricRepository := newInMemoryRubricRepository()
	service := application.NewRubricService(
		rubricRepository,
		newInMemoryStudentRepository(),
	)

	if err := service.AddCriterion("criterion-1", "Экзамен", 100, 100); err != nil {
		t.Fatalf("first AddCriterion returned error: %v", err)
	}

	err := service.AddCriterion("criterion-1", "Контрольная", 50, 100)
	if err == nil {
		t.Fatal("expected error when adding duplicate criterion")
	}
}

func TestRubricServiceAddInvalidCriterion(t *testing.T) {
	rubricRepository := newInMemoryRubricRepository()
	service := application.NewRubricService(
		rubricRepository,
		newInMemoryStudentRepository(),
	)

	err := service.AddCriterion("criterion-1", "Экзамен", 0, 100)
	if err == nil {
		t.Fatal("expected error for zero maximum score")
	}

	rubric, err := rubricRepository.Get()
	if err != nil {
		t.Fatalf("Get rubric: %v", err)
	}

	if len(rubric.GetCriteria()) != 0 {
		t.Error("invalid criterion must not be added")
	}
}

func TestRubricServiceUpdateCriterion(t *testing.T) {
	rubricRepository := newInMemoryRubricRepository()
	service := application.NewRubricService(
		rubricRepository,
		newInMemoryStudentRepository(),
	)

	if err := service.AddCriterion("criterion-1", "Экзамен", 100, 100); err != nil {
		t.Fatalf("AddCriterion: %v", err)
	}

	err := service.UpdateCriterion("criterion-1", "Итоговый экзамен", 50, 100)
	if err != nil {
		t.Fatalf("UpdateCriterion returned unexpected error: %v", err)
	}

	rubric, err := rubricRepository.Get()
	if err != nil {
		t.Fatalf("Get rubric: %v", err)
	}

	criterion, exists := rubric.GetCriterion("criterion-1")
	if !exists {
		t.Fatal("criterion does not exist after update")
	}

	if criterion.Name() != "Итоговый экзамен" {
		t.Errorf("expected updated name, got %s", criterion.Name())
	}

	if criterion.MaxScore().Value() != 50 {
		t.Errorf("expected updated maximum score 50, got %d", criterion.MaxScore().Value())
	}
}

func TestRubricServiceUpdateNonexistentCriterion(t *testing.T) {
	service := application.NewRubricService(
		newInMemoryRubricRepository(),
		newInMemoryStudentRepository(),
	)

	err := service.UpdateCriterion("unknown", "Экзамен", 100, 100)
	if err == nil {
		t.Fatal("expected error when updating nonexistent criterion")
	}
}

func TestRubricServiceCannotUpdateCriterionWithEvaluations(t *testing.T) {
	rubricRepository := newInMemoryRubricRepository()
	studentRepository := newInMemoryStudentRepository()
	service := application.NewRubricService(rubricRepository, studentRepository)

	if err := service.AddCriterion("criterion-1", "Экзамен", 100, 100); err != nil {
		t.Fatalf("AddCriterion: %v", err)
	}

	student, err := createTestStudent("student-1")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	score, err := domain.NewScoreValue(80)
	if err != nil {
		t.Fatalf("create score: %v", err)
	}

	maxScore, err := domain.NewMaxScore(100)
	if err != nil {
		t.Fatalf("create maximum score: %v", err)
	}

	if err := student.SetEvaluation("criterion-1", score, maxScore, ""); err != nil {
		t.Fatalf("SetEvaluation: %v", err)
	}

	if err := studentRepository.Save(student); err != nil {
		t.Fatalf("save student: %v", err)
	}

	err = service.UpdateCriterion("criterion-1", "Новый экзамен", 50, 100)
	if err == nil {
		t.Fatal("expected error when updating criterion with evaluations")
	}

	if !strings.Contains(err.Error(), "evaluations already exist") {
		t.Errorf("expected evaluations restriction error, got: %v", err)
	}
}

func TestRubricServiceRemoveCriterion(t *testing.T) {
	rubricRepository := newInMemoryRubricRepository()
	service := application.NewRubricService(
		rubricRepository,
		newInMemoryStudentRepository(),
	)

	if err := service.AddCriterion("criterion-1", "Экзамен", 100, 100); err != nil {
		t.Fatalf("AddCriterion: %v", err)
	}

	if err := service.RemoveCriterion("criterion-1"); err != nil {
		t.Fatalf("RemoveCriterion returned unexpected error: %v", err)
	}

	rubric, err := rubricRepository.Get()
	if err != nil {
		t.Fatalf("Get rubric: %v", err)
	}

	if _, exists := rubric.GetCriterion("criterion-1"); exists {
		t.Error("criterion still exists after removal")
	}
}

func TestRubricServiceRemoveNonexistentCriterion(t *testing.T) {
	service := application.NewRubricService(
		newInMemoryRubricRepository(),
		newInMemoryStudentRepository(),
	)

	err := service.RemoveCriterion("unknown")
	if err == nil {
		t.Fatal("expected error when removing nonexistent criterion")
	}
}

func TestRubricServiceCannotRemoveCriterionWithEvaluations(t *testing.T) {
	rubricRepository := newInMemoryRubricRepository()
	studentRepository := newInMemoryStudentRepository()
	service := application.NewRubricService(rubricRepository, studentRepository)

	if err := service.AddCriterion("criterion-1", "Экзамен", 100, 100); err != nil {
		t.Fatalf("AddCriterion: %v", err)
	}

	student, err := createTestStudent("student-1")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	score, err := domain.NewScoreValue(80)
	if err != nil {
		t.Fatalf("create score: %v", err)
	}

	maxScore, err := domain.NewMaxScore(100)
	if err != nil {
		t.Fatalf("create maximum score: %v", err)
	}

	if err := student.SetEvaluation("criterion-1", score, maxScore, ""); err != nil {
		t.Fatalf("SetEvaluation: %v", err)
	}

	if err := studentRepository.Save(student); err != nil {
		t.Fatalf("save student: %v", err)
	}

	err = service.RemoveCriterion("criterion-1")
	if err == nil {
		t.Fatal("expected error when removing criterion with evaluations")
	}

	if !strings.Contains(err.Error(), "evaluations already exist") {
		t.Errorf("expected evaluations restriction error, got: %v", err)
	}
}

func TestRubricServiceRepositoryErrors(t *testing.T) {
	t.Run("get rubric error", func(t *testing.T) {
		rubricRepository := newInMemoryRubricRepository()
		rubricRepository.getErr = errors.New("storage failure")

		service := application.NewRubricService(
			rubricRepository,
			newInMemoryStudentRepository(),
		)

		if err := service.AddCriterion("criterion-1", "Экзамен", 100, 100); err == nil {
			t.Fatal("expected repository error")
		}
	})

	t.Run("check evaluations error", func(t *testing.T) {
		rubricRepository := newInMemoryRubricRepository()
		studentRepository := newInMemoryStudentRepository()
		studentRepository.existsEvaluationForErr = errors.New("storage failure")

		criterion, err := createTestCriterion("criterion-1", "Экзамен", 100, 100)
		if err != nil {
			t.Fatalf("create criterion: %v", err)
		}

		if err := rubricRepository.rubric.AddCriterion(criterion); err != nil {
			t.Fatalf("add criterion: %v", err)
		}

		service := application.NewRubricService(rubricRepository, studentRepository)

		if err := service.RemoveCriterion("criterion-1"); err == nil {
			t.Fatal("expected repository error")
		}
	})

	t.Run("save rubric error", func(t *testing.T) {
		rubricRepository := newInMemoryRubricRepository()
		rubricRepository.saveErr = errors.New("storage failure")

		service := application.NewRubricService(
			rubricRepository,
			newInMemoryStudentRepository(),
		)

		if err := service.AddCriterion("criterion-1", "Экзамен", 100, 100); err == nil {
			t.Fatal("expected repository error")
		}
	})
}
