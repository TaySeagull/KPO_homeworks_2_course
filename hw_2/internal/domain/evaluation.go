package domain

import (
	"fmt"
	"time"
)

// Evaluation represents a student's assessment for a criterion
type Evaluation struct {
	criterionID string
	score       ScoreValue
	comment     string
	updatedAt   time.Time
}

func NewEvaluation(
	criterionID string,
	score ScoreValue,
	comment string,
	maxScore MaxScore,
) (Evaluation, error) {
	return NewEvaluationAt(
		criterionID,
		score,
		comment,
		maxScore,
		time.Now(),
	)
}

// NewEvaluationAt creates an evaluation with a specified update time
func NewEvaluationAt(
	criterionID string,
	score ScoreValue,
	comment string,
	maxScore MaxScore,
	updatedAt time.Time,
) (Evaluation, error) {
	if criterionID == "" {
		return Evaluation{}, fmt.Errorf("criterion ID cannot be empty")
	}

	if score.Value() < 0 {
		return Evaluation{}, fmt.Errorf("score cannot be negative")
	}

	if maxScore.Value() <= 0 {
		return Evaluation{}, fmt.Errorf("maximum score must be positive")
	}

	if score.Value() > maxScore.Value() {
		return Evaluation{}, fmt.Errorf(
			"score %d exceeds maximum score %d",
			score.Value(),
			maxScore.Value(),
		)
	}

	if updatedAt.IsZero() {
		return Evaluation{}, fmt.Errorf("evaluation update time cannot be zero")
	}

	return Evaluation{
		criterionID: criterionID,
		score:       score,
		comment:     comment,
		updatedAt:   updatedAt,
	}, nil
}

func (evaluation Evaluation) CriterionID() string {
	return evaluation.criterionID
}

func (evaluation Evaluation) Score() ScoreValue {
	return evaluation.score
}

func (evaluation Evaluation) Comment() string {
	return evaluation.comment
}

func (evaluation Evaluation) UpdatedAt() time.Time {
	return evaluation.updatedAt
}
