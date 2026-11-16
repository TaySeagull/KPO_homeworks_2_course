package application

import "hw_2/internal/domain"

// ScoreCalculationStrategy defines an algorithm for calculating a student's final score
type ScoreCalculationStrategy interface {
	Calculate(
		evaluations []domain.Evaluation,
		criteria []domain.Criterion,
	) (domain.FinalScore, error)
}

// WeightedScoreStrategy calculates the final score using criterion weights
type WeightedScoreStrategy struct{}

// Calculate applies the weighted scoring algorithm
func (strategy WeightedScoreStrategy) Calculate(
	evaluations []domain.Evaluation,
	criteria []domain.Criterion,
) (domain.FinalScore, error) {
	return domain.CalculateFinalScore(evaluations, criteria)
}
