package application_test

import (
	"errors"
	"strings"
	"testing"

	"hw_2/internal/application"
	"hw_2/internal/domain"
)

func TestEvaluationServiceSetEvaluation(t *testing.T) {
	studentRepository := newInMemoryStudentRepository()
	rubricRepository := newInMemoryRubricRepository()

	student, err := createTestStudent("student-1")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	if err := studentRepository.Save(student); err != nil {
		t.Fatalf("save student: %v", err)
	}

	criterion, err := createTestCriterion("criterion-1", "Экзамен", 100, 100)
	if err != nil {
		t.Fatalf("create criterion: %v", err)
	}

	if err := rubricRepository.rubric.AddCriterion(criterion); err != nil {
		t.Fatalf("add criterion: %v", err)
	}

	service := application.NewEvaluationService(studentRepository, rubricRepository)

	err = service.SetEvaluation("student-1", "criterion-1", 85, "Хорошая работа")
	if err != nil {
		t.Fatalf("SetEvaluation returned unexpected error: %v", err)
	}

	savedStudent, err := studentRepository.GetByID("student-1")
	if err != nil {
		t.Fatalf("get saved student: %v", err)
	}

	evaluation, exists := savedStudent.GetEvaluation("criterion-1")
	if !exists {
		t.Fatal("evaluation was not saved")
	}

	if evaluation.Score().Value() != 85 {
		t.Errorf("expected score 85, got %d", evaluation.Score().Value())
	}

	if evaluation.Comment() != "Хорошая работа" {
		t.Errorf("unexpected comment: %s", evaluation.Comment())
	}
}

func TestEvaluationServiceUpdateEvaluation(t *testing.T) {
	studentRepository := newInMemoryStudentRepository()
	rubricRepository := newInMemoryRubricRepository()

	student, err := createTestStudent("student-1")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	if err := studentRepository.Save(student); err != nil {
		t.Fatalf("save student: %v", err)
	}

	criterion, err := createTestCriterion("criterion-1", "Экзамен", 100, 100)
	if err != nil {
		t.Fatalf("create criterion: %v", err)
	}

	if err := rubricRepository.rubric.AddCriterion(criterion); err != nil {
		t.Fatalf("add criterion: %v", err)
	}

	service := application.NewEvaluationService(studentRepository, rubricRepository)

	if err := service.SetEvaluation("student-1", "criterion-1", 60, "Первая оценка"); err != nil {
		t.Fatalf("first SetEvaluation: %v", err)
	}

	if err := service.SetEvaluation("student-1", "criterion-1", 90, "Исправленная оценка"); err != nil {
		t.Fatalf("second SetEvaluation: %v", err)
	}

	savedStudent, err := studentRepository.GetByID("student-1")
	if err != nil {
		t.Fatalf("get saved student: %v", err)
	}

	evaluation, exists := savedStudent.GetEvaluation("criterion-1")
	if !exists {
		t.Fatal("evaluation does not exist")
	}

	if evaluation.Score().Value() != 90 {
		t.Errorf("expected updated score 90, got %d", evaluation.Score().Value())
	}

	if evaluation.Comment() != "Исправленная оценка" {
		t.Errorf("expected updated comment, got %s", evaluation.Comment())
	}
}

func TestEvaluationServiceSetEvaluationForUnknownStudent(t *testing.T) {
	service := application.NewEvaluationService(
		newInMemoryStudentRepository(),
		newInMemoryRubricRepository(),
	)

	err := service.SetEvaluation("unknown", "criterion-1", 80, "")
	if err == nil {
		t.Fatal("expected error for nonexistent student")
	}
}

func TestEvaluationServiceSetEvaluationForUnknownCriterion(t *testing.T) {
	studentRepository := newInMemoryStudentRepository()

	student, err := createTestStudent("student-1")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	if err := studentRepository.Save(student); err != nil {
		t.Fatalf("save student: %v", err)
	}

	service := application.NewEvaluationService(
		studentRepository,
		newInMemoryRubricRepository(),
	)

	err = service.SetEvaluation("student-1", "unknown", 80, "")
	if err == nil {
		t.Fatal("expected error for nonexistent criterion")
	}

	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("expected unknown criterion error, got: %v", err)
	}
}

func TestEvaluationServiceRejectsScoreAboveMaximum(t *testing.T) {
	studentRepository := newInMemoryStudentRepository()
	rubricRepository := newInMemoryRubricRepository()

	student, err := createTestStudent("student-1")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	if err := studentRepository.Save(student); err != nil {
		t.Fatalf("save student: %v", err)
	}

	criterion, err := createTestCriterion("criterion-1", "Экзамен", 100, 100)
	if err != nil {
		t.Fatalf("create criterion: %v", err)
	}

	if err := rubricRepository.rubric.AddCriterion(criterion); err != nil {
		t.Fatalf("add criterion: %v", err)
	}

	service := application.NewEvaluationService(studentRepository, rubricRepository)

	err = service.SetEvaluation("student-1", "criterion-1", 101, "")
	if err == nil {
		t.Fatal("expected error for score above maximum")
	}

	savedStudent, err := studentRepository.GetByID("student-1")
	if err != nil {
		t.Fatalf("get student: %v", err)
	}

	if _, exists := savedStudent.GetEvaluation("criterion-1"); exists {
		t.Error("invalid evaluation must not be saved")
	}
}

func TestEvaluationServiceRejectsNegativeScore(t *testing.T) {
	studentRepository := newInMemoryStudentRepository()
	rubricRepository := newInMemoryRubricRepository()

	student, err := createTestStudent("student-1")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	if err := studentRepository.Save(student); err != nil {
		t.Fatalf("save student: %v", err)
	}

	criterion, err := createTestCriterion("criterion-1", "Экзамен", 100, 100)
	if err != nil {
		t.Fatalf("create criterion: %v", err)
	}

	if err := rubricRepository.rubric.AddCriterion(criterion); err != nil {
		t.Fatalf("add criterion: %v", err)
	}

	service := application.NewEvaluationService(studentRepository, rubricRepository)

	err = service.SetEvaluation("student-1", "criterion-1", -1, "")
	if err == nil {
		t.Fatal("expected error for negative score")
	}
}

func TestEvaluationServiceRepositoryErrors(t *testing.T) {
	t.Run("get student error", func(t *testing.T) {
		studentRepository := newInMemoryStudentRepository()
		studentRepository.getByIDErr = errors.New("storage failure")

		service := application.NewEvaluationService(
			studentRepository,
			newInMemoryRubricRepository(),
		)

		if err := service.SetEvaluation("student-1", "criterion-1", 80, ""); err == nil {
			t.Fatal("expected repository error")
		}
	})

	t.Run("get rubric error", func(t *testing.T) {
		studentRepository := newInMemoryStudentRepository()
		rubricRepository := newInMemoryRubricRepository()
		rubricRepository.getErr = errors.New("storage failure")

		student, err := createTestStudent("student-1")
		if err != nil {
			t.Fatalf("create student: %v", err)
		}

		if err := studentRepository.Save(student); err != nil {
			t.Fatalf("save student: %v", err)
		}

		service := application.NewEvaluationService(studentRepository, rubricRepository)

		if err := service.SetEvaluation("student-1", "criterion-1", 80, ""); err == nil {
			t.Fatal("expected repository error")
		}
	})

	t.Run("save student error", func(t *testing.T) {
		studentRepository := newInMemoryStudentRepository()
		rubricRepository := newInMemoryRubricRepository()

		student, err := createTestStudent("student-1")
		if err != nil {
			t.Fatalf("create student: %v", err)
		}

		if err := studentRepository.Save(student); err != nil {
			t.Fatalf("save student: %v", err)
		}

		criterion, err := createTestCriterion("criterion-1", "Экзамен", 100, 100)
		if err != nil {
			t.Fatalf("create criterion: %v", err)
		}

		if err := rubricRepository.rubric.AddCriterion(criterion); err != nil {
			t.Fatalf("add criterion: %v", err)
		}

		studentRepository.saveErr = errors.New("storage failure")
		service := application.NewEvaluationService(studentRepository, rubricRepository)

		if err := service.SetEvaluation("student-1", "criterion-1", 80, ""); err == nil {
			t.Fatal("expected repository error")
		}
	})
}

func TestScoreServiceCalculateStudentScore(t *testing.T) {
	studentRepository := newInMemoryStudentRepository()
	rubricRepository := newInMemoryRubricRepository()

	student, err := createTestStudent("student-1")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	firstCriterion, err := createTestCriterion("criterion-1", "Контрольная", 10, 50)
	if err != nil {
		t.Fatalf("create first criterion: %v", err)
	}

	secondCriterion, err := createTestCriterion("criterion-2", "Экзамен", 20, 50)
	if err != nil {
		t.Fatalf("create second criterion: %v", err)
	}

	if err := rubricRepository.rubric.AddCriterion(firstCriterion); err != nil {
		t.Fatalf("add first criterion: %v", err)
	}

	if err := rubricRepository.rubric.AddCriterion(secondCriterion); err != nil {
		t.Fatalf("add second criterion: %v", err)
	}

	firstScore, err := domain.NewScoreValue(8)
	if err != nil {
		t.Fatalf("create first score: %v", err)
	}

	secondScore, err := domain.NewScoreValue(10)
	if err != nil {
		t.Fatalf("create second score: %v", err)
	}

	firstMaxScore, err := domain.NewMaxScore(10)
	if err != nil {
		t.Fatalf("create first maximum score: %v", err)
	}

	secondMaxScore, err := domain.NewMaxScore(20)
	if err != nil {
		t.Fatalf("create second maximum score: %v", err)
	}

	if err := student.SetEvaluation("criterion-1", firstScore, firstMaxScore, ""); err != nil {
		t.Fatalf("set first evaluation: %v", err)
	}

	if err := student.SetEvaluation("criterion-2", secondScore, secondMaxScore, ""); err != nil {
		t.Fatalf("set second evaluation: %v", err)
	}

	if err := studentRepository.Save(student); err != nil {
		t.Fatalf("save student: %v", err)
	}

	service := application.NewScoreService(studentRepository, rubricRepository)

	result, err := service.CalculateStudentScore("student-1")
	if err != nil {
		t.Fatalf("CalculateStudentScore returned unexpected error: %v", err)
	}

	if result.Value() != 7 {
		t.Errorf("expected final score 7, got %d", result.Value())
	}
}

func TestScoreServiceCalculateStudentScoreWithMissingEvaluation(t *testing.T) {
	studentRepository := newInMemoryStudentRepository()
	rubricRepository := newInMemoryRubricRepository()

	student, err := createTestStudent("student-1")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	criterion, err := createTestCriterion("criterion-1", "Экзамен", 100, 100)
	if err != nil {
		t.Fatalf("create criterion: %v", err)
	}

	if err := rubricRepository.rubric.AddCriterion(criterion); err != nil {
		t.Fatalf("add criterion: %v", err)
	}

	if err := studentRepository.Save(student); err != nil {
		t.Fatalf("save student: %v", err)
	}

	service := application.NewScoreService(studentRepository, rubricRepository)

	result, err := service.CalculateStudentScore("student-1")
	if err != nil {
		t.Fatalf("CalculateStudentScore returned unexpected error: %v", err)
	}

	if result.Value() != 0 {
		t.Errorf("expected final score 0 for missing evaluation, got %d", result.Value())
	}
}

func TestScoreServiceStudentNotFound(t *testing.T) {
	service := application.NewScoreService(
		newInMemoryStudentRepository(),
		newInMemoryRubricRepository(),
	)

	_, err := service.CalculateStudentScore("unknown")
	if err == nil {
		t.Fatal("expected error for nonexistent student")
	}
}

func TestScoreServiceInvalidRubricWeights(t *testing.T) {
	studentRepository := newInMemoryStudentRepository()
	rubricRepository := newInMemoryRubricRepository()

	student, err := createTestStudent("student-1")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	criterion, err := createTestCriterion("criterion-1", "Экзамен", 100, 50)
	if err != nil {
		t.Fatalf("create criterion: %v", err)
	}

	if err := rubricRepository.rubric.AddCriterion(criterion); err != nil {
		t.Fatalf("add criterion: %v", err)
	}

	if err := studentRepository.Save(student); err != nil {
		t.Fatalf("save student: %v", err)
	}

	service := application.NewScoreService(studentRepository, rubricRepository)

	_, err = service.CalculateStudentScore("student-1")
	if err == nil {
		t.Fatal("expected error for rubric with invalid total weight")
	}
}

func TestScoreServiceRepositoryErrors(t *testing.T) {
	t.Run("get student error", func(t *testing.T) {
		studentRepository := newInMemoryStudentRepository()
		studentRepository.getByIDErr = errors.New("storage failure")

		service := application.NewScoreService(
			studentRepository,
			newInMemoryRubricRepository(),
		)

		if _, err := service.CalculateStudentScore("student-1"); err == nil {
			t.Fatal("expected repository error")
		}
	})

	t.Run("get rubric error", func(t *testing.T) {
		studentRepository := newInMemoryStudentRepository()
		rubricRepository := newInMemoryRubricRepository()
		rubricRepository.getErr = errors.New("storage failure")

		student, err := createTestStudent("student-1")
		if err != nil {
			t.Fatalf("create student: %v", err)
		}

		if err := studentRepository.Save(student); err != nil {
			t.Fatalf("save student: %v", err)
		}

		service := application.NewScoreService(studentRepository, rubricRepository)

		if _, err := service.CalculateStudentScore("student-1"); err == nil {
			t.Fatal("expected repository error")
		}
	})
}
