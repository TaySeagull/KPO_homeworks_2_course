package domain

import (
	"fmt"
	"sort"
)

// Rubric is the aggregate root that manages assessment criteria
type Rubric struct {
	criteria map[string]Criterion
}

func NewRubric() *Rubric {
	return &Rubric{
		criteria: make(map[string]Criterion),
	}
}

func (rubric *Rubric) AddCriterion(criterion Criterion) error {
	if err := validateCriterion(criterion); err != nil {
		return err
	}

	if _, exists := rubric.criteria[criterion.ID()]; exists {
		return fmt.Errorf("criterion with ID %s already exists", criterion.ID())
	}

	rubric.criteria[criterion.ID()] = criterion
	return nil
}

func (rubric *Rubric) RemoveCriterion(id string) error {
	if _, exists := rubric.criteria[id]; !exists {
		return fmt.Errorf("criterion with ID %s does not exist", id)
	}

	delete(rubric.criteria, id)
	return nil
}

func (rubric *Rubric) UpdateCriterion(criterion Criterion) error {
	if err := validateCriterion(criterion); err != nil {
		return err
	}

	if _, exists := rubric.criteria[criterion.ID()]; !exists {
		return fmt.Errorf("criterion with ID %s does not exist", criterion.ID())
	}

	rubric.criteria[criterion.ID()] = criterion
	return nil
}

func (rubric *Rubric) GetCriterion(id string) (Criterion, bool) {
	criterion, exists := rubric.criteria[id]
	return criterion, exists
}

func (rubric *Rubric) GetCriteria() []Criterion {
	result := make([]Criterion, 0, len(rubric.criteria))

	for _, criterion := range rubric.criteria {
		result = append(result, criterion)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID() < result[j].ID()
	})

	return result
}

func (rubric *Rubric) ValidateWeights() error {
	if len(rubric.criteria) == 0 {
		return fmt.Errorf("rubric must contain at least one criterion")
	}

	total := 0

	for _, criterion := range rubric.criteria {
		total += criterion.Weight().Value()
	}

	if total != 100 {
		return fmt.Errorf(
			"total criterion weights must equal 100 percent, got %d",
			total,
		)
	}

	return nil
}

func validateCriterion(criterion Criterion) error {
	if criterion.ID() == "" {
		return fmt.Errorf("criterion ID cannot be empty")
	}

	if criterion.Name() == "" {
		return fmt.Errorf("criterion name cannot be empty")
	}

	if criterion.MaxScore().Value() <= 0 {
		return fmt.Errorf("criterion maximum score must be positive")
	}

	if criterion.Weight().Value() < 1 || criterion.Weight().Value() > 100 {
		return fmt.Errorf("criterion weight must be between 1 and 100")
	}

	return nil
}
