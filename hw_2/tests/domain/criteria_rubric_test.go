package domain_test

import (
	"testing"

	"hw_2/internal/domain"
)

func createCriterion(
	t *testing.T,
	id string,
	name string,
	maxScore int,
	weight int,
) domain.Criterion {
	t.Helper()

	max, err := domain.NewMaxScore(maxScore)
	if err != nil {
		t.Fatalf("failed to create maximum score: %v", err)
	}

	w, err := domain.NewWeight(weight)
	if err != nil {
		t.Fatalf("failed to create weight: %v", err)
	}

	criterion, err := domain.NewCriterion(id, name, max, w)
	if err != nil {
		t.Fatalf("failed to create criterion: %v", err)
	}

	return criterion
}

func TestNewCriterion(t *testing.T) {
	max, _ := domain.NewMaxScore(10)
	weight, _ := domain.NewWeight(50)

	tests := []struct {
		name          string
		id            string
		criterionName string
		maxScore      domain.MaxScore
		weight        domain.Weight
		wantErr       bool
	}{
		{
			name:          "valid criterion",
			id:            "1",
			criterionName: "Exam",
			maxScore:      max,
			weight:        weight,
		},
		{
			name:          "empty ID",
			id:            "",
			criterionName: "Exam",
			maxScore:      max,
			weight:        weight,
			wantErr:       true,
		},
		{
			name:          "empty name",
			id:            "1",
			criterionName: "",
			maxScore:      max,
			weight:        weight,
			wantErr:       true,
		},
		{
			name:          "zero maximum score",
			id:            "1",
			criterionName: "Exam",
			maxScore:      domain.MaxScore{},
			weight:        weight,
			wantErr:       true,
		},
		{
			name:          "zero weight",
			id:            "1",
			criterionName: "Exam",
			maxScore:      max,
			weight:        domain.Weight{},
			wantErr:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := domain.NewCriterion(
				test.id,
				test.criterionName,
				test.maxScore,
				test.weight,
			)

			if (err != nil) != test.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestRubricAddCriterion(t *testing.T) {
	rubric := domain.NewRubric()

	criterion := createCriterion(t, "1", "Exam", 10, 100)

	if err := rubric.AddCriterion(criterion); err != nil {
		t.Fatalf("failed to add criterion: %v", err)
	}

	if _, exists := rubric.GetCriterion("1"); !exists {
		t.Fatal("criterion was not added")
	}
}

func TestRubricAddDuplicateCriterion(t *testing.T) {
	rubric := domain.NewRubric()

	first := createCriterion(t, "1", "Exam", 10, 100)
	second := createCriterion(t, "1", "Project", 20, 100)

	if err := rubric.AddCriterion(first); err != nil {
		t.Fatalf("failed to add first criterion: %v", err)
	}

	if err := rubric.AddCriterion(second); err == nil {
		t.Fatal("expected error when adding duplicate criterion")
	}
}

func TestRubricAddInvalidCriterion(t *testing.T) {
	rubric := domain.NewRubric()

	if err := rubric.AddCriterion(domain.Criterion{}); err == nil {
		t.Fatal("expected error when adding invalid criterion")
	}
}

func TestRubricUpdateCriterion(t *testing.T) {
	rubric := domain.NewRubric()

	original := createCriterion(t, "1", "Exam", 10, 100)
	updated := createCriterion(t, "1", "Final exam", 20, 100)

	if err := rubric.AddCriterion(original); err != nil {
		t.Fatalf("failed to add criterion: %v", err)
	}

	if err := rubric.UpdateCriterion(updated); err != nil {
		t.Fatalf("failed to update criterion: %v", err)
	}

	result, exists := rubric.GetCriterion("1")
	if !exists {
		t.Fatal("criterion does not exist")
	}

	if result.Name() != "Final exam" {
		t.Errorf("got name %q, want %q", result.Name(), "Final exam")
	}

	if result.MaxScore().Value() != 20 {
		t.Errorf("got maximum score %d, want 20", result.MaxScore().Value())
	}
}

func TestRubricUpdateNonexistentCriterion(t *testing.T) {
	rubric := domain.NewRubric()

	criterion := createCriterion(t, "1", "Exam", 10, 100)

	if err := rubric.UpdateCriterion(criterion); err == nil {
		t.Fatal("expected error when updating nonexistent criterion")
	}
}

func TestRubricRemoveCriterion(t *testing.T) {
	rubric := domain.NewRubric()

	criterion := createCriterion(t, "1", "Exam", 10, 100)

	if err := rubric.AddCriterion(criterion); err != nil {
		t.Fatalf("failed to add criterion: %v", err)
	}

	if err := rubric.RemoveCriterion("1"); err != nil {
		t.Fatalf("failed to remove criterion: %v", err)
	}

	if _, exists := rubric.GetCriterion("1"); exists {
		t.Fatal("criterion still exists after removal")
	}
}

func TestRubricRemoveNonexistentCriterion(t *testing.T) {
	rubric := domain.NewRubric()

	if err := rubric.RemoveCriterion("unknown"); err == nil {
		t.Fatal("expected error when removing nonexistent criterion")
	}
}

func TestRubricValidateWeights(t *testing.T) {
	tests := []struct {
		name    string
		weights []int
		wantErr bool
	}{
		{
			name:    "weights equal 100",
			weights: []int{40, 60},
		},
		{
			name:    "weights less than 100",
			weights: []int{30, 50},
			wantErr: true,
		},
		{
			name:    "weights greater than 100",
			weights: []int{60, 60},
			wantErr: true,
		},
		{
			name:    "empty rubric",
			weights: []int{},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rubric := domain.NewRubric()

			for i, weight := range test.weights {
				criterion := createCriterion(
					t,
					string(rune('1'+i)),
					"Criterion",
					10,
					weight,
				)

				if err := rubric.AddCriterion(criterion); err != nil {
					t.Fatalf("failed to add criterion: %v", err)
				}
			}

			err := rubric.ValidateWeights()

			if (err != nil) != test.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
