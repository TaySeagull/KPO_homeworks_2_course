package jsonstore

import (
	"os"
	"path/filepath"
	"testing"

	"hw_2/internal/domain"
	"hw_2/internal/infrastructure/jsonstore"
)

func createTestRubric(t *testing.T) domain.Rubric {
	t.Helper()

	rubric := domain.NewRubric()

	maxScore, err := domain.NewMaxScore(10)
	if err != nil {
		t.Fatalf("create max score: %v", err)
	}

	weight, err := domain.NewWeight(100)
	if err != nil {
		t.Fatalf("create weight: %v", err)
	}

	criterion, err := domain.NewCriterion(
		"criterion-1",
		"Практические задания",
		maxScore,
		weight,
	)
	if err != nil {
		t.Fatalf("create criterion: %v", err)
	}

	if err := rubric.AddCriterion(criterion); err != nil {
		t.Fatalf("add criterion: %v", err)
	}

	return *rubric
}

func TestJSONRubricRepositorySaveAndGet(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "rubric.json")
	repository := jsonstore.NewJSONRubricRepository(filePath)

	expected := createTestRubric(t)

	if err := repository.Save(expected); err != nil {
		t.Fatalf("save rubric: %v", err)
	}

	actual, err := repository.Get()
	if err != nil {
		t.Fatalf("get rubric: %v", err)
	}

	expectedCriteria := expected.GetCriteria()
	actualCriteria := actual.GetCriteria()

	if len(actualCriteria) != len(expectedCriteria) {
		t.Fatalf(
			"expected %d criteria, got %d",
			len(expectedCriteria),
			len(actualCriteria),
		)
	}

	expectedCriterion := expectedCriteria[0]

	actualCriterion, exists := actual.GetCriterion(expectedCriterion.ID())
	if !exists {
		t.Fatalf("criterion %q was not found", expectedCriterion.ID())
	}

	if actualCriterion.Name() != expectedCriterion.Name() {
		t.Errorf(
			"expected criterion name %q, got %q",
			expectedCriterion.Name(),
			actualCriterion.Name(),
		)
	}

	if actualCriterion.MaxScore().Value() != expectedCriterion.MaxScore().Value() {
		t.Errorf(
			"expected max score %d, got %d",
			expectedCriterion.MaxScore().Value(),
			actualCriterion.MaxScore().Value(),
		)
	}

	if actualCriterion.Weight().Value() != expectedCriterion.Weight().Value() {
		t.Errorf(
			"expected weight %d, got %d",
			expectedCriterion.Weight().Value(),
			actualCriterion.Weight().Value(),
		)
	}
}

func TestJSONRubricRepositoryGetMissingFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "missing.json")
	repository := jsonstore.NewJSONRubricRepository(filePath)

	rubric, err := repository.Get()
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}

	if len(rubric.GetCriteria()) != 0 {
		t.Errorf(
			"expected empty rubric, got %d criteria",
			len(rubric.GetCriteria()),
		)
	}
}

func TestJSONRubricRepositoryGetInvalidJSON(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "rubric.json")

	if err := os.WriteFile(filePath, []byte("{invalid json"), 0644); err != nil {
		t.Fatalf("write invalid JSON: %v", err)
	}

	repository := jsonstore.NewJSONRubricRepository(filePath)

	_, err := repository.Get()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
