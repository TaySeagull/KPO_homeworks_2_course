package domain

import (
	"fmt"
	"math"
)

// CalculateFinalScore calculates the weighted final score for a student
func CalculateFinalScore(
	evaluations []Evaluation,
	criteria []Criterion,
) (FinalScore, error) {
	if len(criteria) == 0 {
		return FinalScore{}, fmt.Errorf("rubric must contain at least one criterion")
	}

	criteriaByID := make(map[string]Criterion, len(criteria))
	totalWeight := 0

	for _, criterion := range criteria {
		if err := validateCriterion(criterion); err != nil {
			return FinalScore{}, fmt.Errorf(
				"invalid criterion %s: %w",
				criterion.ID(),
				err,
			)
		}

		if _, exists := criteriaByID[criterion.ID()]; exists {
			return FinalScore{}, fmt.Errorf(
				"duplicate criterion ID: %s",
				criterion.ID(),
			)
		}

		criteriaByID[criterion.ID()] = criterion
		totalWeight += criterion.Weight().Value()
	}

	if totalWeight != 100 {
		return FinalScore{}, fmt.Errorf(
			"total criterion weights must equal 100 percent, got %d",
			totalWeight,
		)
	}

	evaluationsByCriterion := make(map[string]Evaluation, len(evaluations))

	for _, evaluation := range evaluations {
		criterionID := evaluation.CriterionID()

		criterion, exists := criteriaByID[criterionID]
		if !exists {
			return FinalScore{}, fmt.Errorf(
				"evaluation references unknown criterion: %s",
				criterionID,
			)
		}

		if _, exists := evaluationsByCriterion[criterionID]; exists {
			return FinalScore{}, fmt.Errorf(
				"duplicate evaluation for criterion: %s",
				criterionID,
			)
		}

		if evaluation.Score().Value() < 0 {
			return FinalScore{}, fmt.Errorf(
				"score for criterion %s cannot be negative",
				criterionID,
			)
		}

		if evaluation.Score().Value() > criterion.MaxScore().Value() {
			return FinalScore{}, fmt.Errorf(
				"score for criterion %s exceeds maximum score",
				criterionID,
			)
		}

		evaluationsByCriterion[criterionID] = evaluation
	}

	weightedSum := 0.0

	for _, criterion := range criteria {
		evaluation, exists := evaluationsByCriterion[criterion.ID()]
		if !exists {
			continue
		}

		normalizedScore := float64(evaluation.Score().Value()) /
			float64(criterion.MaxScore().Value())

		weight := float64(criterion.Weight().Value()) / 100

		weightedSum += normalizedScore * weight
	}

	result := int(math.Round(10 * weightedSum))

	return NewFinalScore(result)
}
