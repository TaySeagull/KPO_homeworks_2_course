package domain

import "fmt"

// ScoreValue represents a student's score for a criterion
type ScoreValue struct {
	value int
}

func NewScoreValue(value int) (ScoreValue, error) {
	if value < 0 {
		return ScoreValue{}, fmt.Errorf("score cannot be negative")
	}

	return ScoreValue{value: value}, nil
}

func (score ScoreValue) Value() int {
	return score.value
}

// MaxScore represents the maximum possible score for a criterion
type MaxScore struct {
	value int
}

func NewMaxScore(value int) (MaxScore, error) {
	if value <= 0 {
		return MaxScore{}, fmt.Errorf("maximum score must be positive")
	}

	return MaxScore{value: value}, nil
}

func (maxScore MaxScore) Value() int {
	return maxScore.value
}

// Weight represents a criterion weight as a percentage
type Weight struct {
	value int
}

func NewWeight(value int) (Weight, error) {
	if value < 1 || value > 100 {
		return Weight{}, fmt.Errorf("weight must be between 1 and 100")
	}

	return Weight{value: value}, nil
}

func (weight Weight) Value() int {
	return weight.value
}

// FinalScore represents the final score on a 0-10 scale
type FinalScore struct {
	value int
}

func NewFinalScore(value int) (FinalScore, error) {
	if value < 0 || value > 10 {
		return FinalScore{}, fmt.Errorf("final score must be between 0 and 10")
	}

	return FinalScore{value: value}, nil
}

func (score FinalScore) Value() int {
	return score.value
}
