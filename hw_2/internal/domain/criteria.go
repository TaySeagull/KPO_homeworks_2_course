package domain

import "fmt"

// Criterion represents an assessment criterion
type Criterion struct {
	id       string
	name     string
	maxScore MaxScore
	weight   Weight
}

func NewCriterion(
	id string,
	name string,
	maxScore MaxScore,
	weight Weight,
) (Criterion, error) {
	if id == "" {
		return Criterion{}, fmt.Errorf("criterion ID cannot be empty")
	}

	if name == "" {
		return Criterion{}, fmt.Errorf("criterion name cannot be empty")
	}

	if maxScore.Value() <= 0 {
		return Criterion{}, fmt.Errorf("criterion maximum score must be positive")
	}

	if weight.Value() < 1 || weight.Value() > 100 {
		return Criterion{}, fmt.Errorf("criterion weight must be between 1 and 100")
	}

	return Criterion{
		id:       id,
		name:     name,
		maxScore: maxScore,
		weight:   weight,
	}, nil
}

func (criterion Criterion) ID() string {
	return criterion.id
}

func (criterion Criterion) Name() string {
	return criterion.name
}

func (criterion Criterion) MaxScore() MaxScore {
	return criterion.maxScore
}

func (criterion Criterion) Weight() Weight {
	return criterion.weight
}
