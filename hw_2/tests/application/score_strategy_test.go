package application_test

import (
	"errors"
	"testing"

	"hw_2/internal/application"
	"hw_2/internal/domain"
)

type fakeScoreStrategy struct {
	result           domain.FinalScore
	err              error
	called           bool
	evaluationsCount int
	criteriaCount    int
}

func (strategy *fakeScoreStrategy) Calculate(
	evaluations []domain.Evaluation,
	criteria []domain.Criterion,
) (domain.FinalScore, error) {
	strategy.called = true
	strategy.evaluationsCount = len(evaluations)
	strategy.criteriaCount = len(criteria)

	return strategy.result, strategy.err
}

func TestScoreServiceUsesInjectedStrategy(t *testing.T) {
	studentRepository := newInMemoryStudentRepository()
	rubricRepository := newInMemoryRubricRepository()

	studentService := application.NewStudentService(studentRepository)
	rubricService := application.NewRubricService(
		rubricRepository,
		studentRepository,
	)
	evaluationService := application.NewEvaluationService(
		studentRepository,
		rubricRepository,
	)

	if err := studentService.AddStudent(
		"student-1",
		"Траляля",
		"Жужужу",
		"Ляляля",
		"БПИ-241",
		"ФКН",
	); err != nil {
		t.Fatalf("AddStudent() error = %v", err)
	}

	if err := rubricService.AddCriterion(
		"criterion-1",
		"Контрольная работа",
		10,
		100,
	); err != nil {
		t.Fatalf("AddCriterion() error = %v", err)
	}

	if err := evaluationService.SetEvaluation(
		"student-1",
		"criterion-1",
		8,
		"Хорошая работа",
	); err != nil {
		t.Fatalf("SetEvaluation() error = %v", err)
	}

	strategy := &fakeScoreStrategy{
		result: domain.FinalScore{},
	}

	scoreService := application.NewScoreServiceWithStrategy(
		studentRepository,
		rubricRepository,
		strategy,
	)

	_, err := scoreService.CalculateStudentScore("student-1")
	if err != nil {
		t.Fatalf("CalculateStudentScore() error = %v", err)
	}

	if !strategy.called {
		t.Fatal("injected strategy was not called")
	}

	if strategy.evaluationsCount != 1 {
		t.Errorf(
			"evaluations count = %d, want 1",
			strategy.evaluationsCount,
		)
	}

	if strategy.criteriaCount != 1 {
		t.Errorf(
			"criteria count = %d, want 1",
			strategy.criteriaCount,
		)
	}
}

func TestScoreServicePropagatesStrategyError(t *testing.T) {
	studentRepository := newInMemoryStudentRepository()
	rubricRepository := newInMemoryRubricRepository()

	studentService := application.NewStudentService(studentRepository)
	rubricService := application.NewRubricService(
		rubricRepository,
		studentRepository,
	)

	if err := studentService.AddStudent(
		"student-1",
		"Траляля",
		"Жужужу",
		"Ляляля",
		"БПИ-241",
		"ФКН",
	); err != nil {
		t.Fatalf("AddStudent() error = %v", err)
	}

	if err := rubricService.AddCriterion(
		"criterion-1",
		"Контрольная работа",
		10,
		100,
	); err != nil {
		t.Fatalf("AddCriterion() error = %v", err)
	}

	expectedErr := errors.New("strategy calculation error")
	strategy := &fakeScoreStrategy{
		err: expectedErr,
	}

	scoreService := application.NewScoreServiceWithStrategy(
		studentRepository,
		rubricRepository,
		strategy,
	)

	_, err := scoreService.CalculateStudentScore("student-1")

	if !errors.Is(err, expectedErr) {
		t.Errorf("error = %v, want wrapped %v", err, expectedErr)
	}

	if !strategy.called {
		t.Error("injected strategy was not called")
	}
}

func TestScoreServiceUsesDefaultStrategyWhenNil(t *testing.T) {
	studentRepository := newInMemoryStudentRepository()
	rubricRepository := newInMemoryRubricRepository()

	studentService := application.NewStudentService(studentRepository)
	rubricService := application.NewRubricService(
		rubricRepository,
		studentRepository,
	)
	evaluationService := application.NewEvaluationService(
		studentRepository,
		rubricRepository,
	)

	if err := studentService.AddStudent(
		"student-1",
		"Траляля",
		"Жужужу",
		"Ляляля",
		"БПИ-241",
		"ФКН",
	); err != nil {
		t.Fatalf("AddStudent() error = %v", err)
	}

	if err := rubricService.AddCriterion(
		"criterion-1",
		"Контрольная работа",
		10,
		100,
	); err != nil {
		t.Fatalf("AddCriterion() error = %v", err)
	}

	if err := evaluationService.SetEvaluation(
		"student-1",
		"criterion-1",
		8,
		"Хорошая работа",
	); err != nil {
		t.Fatalf("SetEvaluation() error = %v", err)
	}

	scoreService := application.NewScoreServiceWithStrategy(
		studentRepository,
		rubricRepository,
		nil,
	)

	score, err := scoreService.CalculateStudentScore("student-1")
	if err != nil {
		t.Fatalf("CalculateStudentScore() error = %v", err)
	}

	if score.Value() != 8 {
		t.Errorf("score = %d, want 8", score.Value())
	}
}
