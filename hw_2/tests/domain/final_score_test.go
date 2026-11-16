package domain_test

import (
	"testing"

	"hw_2/internal/domain"
)

func TestCalculateFinalScore(t *testing.T) {
	tests := []struct {
		name          string
		firstScore    int
		secondScore   int
		expectedScore int
	}{
		{
			name:          "maximum scores",
			firstScore:    10,
			secondScore:   20,
			expectedScore: 10,
		},
		{
			name:          "zero scores",
			firstScore:    0,
			secondScore:   0,
			expectedScore: 0,
		},
		{
			name:          "weighted average",
			firstScore:    8,
			secondScore:   10,
			expectedScore: 7,
		},
		{
			name:          "half of maximum",
			firstScore:    5,
			secondScore:   10,
			expectedScore: 5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			firstCriterion := createCriterion(t, "1", "Exam", 10, 50)
			secondCriterion := createCriterion(t, "2", "Project", 20, 50)

			firstValue, _ := domain.NewScoreValue(test.firstScore)
			secondValue, _ := domain.NewScoreValue(test.secondScore)

			firstEvaluation, err := domain.NewEvaluation(
				"1",
				firstValue,
				"",
				firstCriterion.MaxScore(),
			)
			if err != nil {
				t.Fatalf("failed to create first evaluation: %v", err)
			}

			secondEvaluation, err := domain.NewEvaluation(
				"2",
				secondValue,
				"",
				secondCriterion.MaxScore(),
			)
			if err != nil {
				t.Fatalf("failed to create second evaluation: %v", err)
			}

			result, err := domain.CalculateFinalScore(
				[]domain.Evaluation{firstEvaluation, secondEvaluation},
				[]domain.Criterion{firstCriterion, secondCriterion},
			)
			if err != nil {
				t.Fatalf("failed to calculate final score: %v", err)
			}

			if result.Value() != test.expectedScore {
				t.Errorf("got %d, want %d", result.Value(), test.expectedScore)
			}
		})
	}
}

func TestCalculateFinalScoreWithoutSomeEvaluations(t *testing.T) {
	firstCriterion := createCriterion(t, "1", "Exam", 10, 50)
	secondCriterion := createCriterion(t, "2", "Project", 20, 50)

	score, _ := domain.NewScoreValue(10)

	evaluation, err := domain.NewEvaluation(
		"1",
		score,
		"",
		firstCriterion.MaxScore(),
	)
	if err != nil {
		t.Fatalf("failed to create evaluation: %v", err)
	}

	result, err := domain.CalculateFinalScore(
		[]domain.Evaluation{evaluation},
		[]domain.Criterion{firstCriterion, secondCriterion},
	)
	if err != nil {
		t.Fatalf("failed to calculate final score: %v", err)
	}

	if result.Value() != 5 {
		t.Errorf("got %d, want 5", result.Value())
	}
}

func TestCalculateFinalScoreWithoutCriteria(t *testing.T) {
	_, err := domain.CalculateFinalScore(nil, nil)

	if err == nil {
		t.Fatal("expected error for empty rubric")
	}
}

func TestCalculateFinalScoreWithInvalidWeights(t *testing.T) {
	firstCriterion := createCriterion(t, "1", "Exam", 10, 30)
	secondCriterion := createCriterion(t, "2", "Project", 20, 30)

	_, err := domain.CalculateFinalScore(
		nil,
		[]domain.Criterion{firstCriterion, secondCriterion},
	)

	if err == nil {
		t.Fatal("expected error when weights do not sum to 100")
	}
}

func TestCalculateFinalScoreWithUnknownCriterion(t *testing.T) {
	criterion := createCriterion(t, "1", "Exam", 10, 100)

	score, _ := domain.NewScoreValue(5)
	unknownEvaluation, err := domain.NewEvaluation(
		"unknown",
		score,
		"",
		criterion.MaxScore(),
	)
	if err != nil {
		t.Fatalf("failed to create evaluation: %v", err)
	}

	_, err = domain.CalculateFinalScore(
		[]domain.Evaluation{unknownEvaluation},
		[]domain.Criterion{criterion},
	)

	if err == nil {
		t.Fatal("expected error for evaluation with unknown criterion")
	}
}

func TestCalculateFinalScoreWithDuplicateEvaluations(t *testing.T) {
	criterion := createCriterion(t, "1", "Exam", 10, 100)

	firstScore, _ := domain.NewScoreValue(5)
	secondScore, _ := domain.NewScoreValue(7)

	firstEvaluation, _ := domain.NewEvaluation(
		"1",
		firstScore,
		"",
		criterion.MaxScore(),
	)

	secondEvaluation, _ := domain.NewEvaluation(
		"1",
		secondScore,
		"",
		criterion.MaxScore(),
	)

	_, err := domain.CalculateFinalScore(
		[]domain.Evaluation{firstEvaluation, secondEvaluation},
		[]domain.Criterion{criterion},
	)

	if err == nil {
		t.Fatal("expected error for duplicate evaluations")
	}
}
