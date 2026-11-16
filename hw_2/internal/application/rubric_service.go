package application

import (
	"fmt"
	"hw_2/internal/domain"
)

// RubricService handles application operations related to the rubric
type RubricService struct {
	rubricRepository  RubricRepository
	studentRepository StudentRepository
}

func NewRubricService(
	rubricRepository RubricRepository,
	studentRepository StudentRepository,
) *RubricService {
	return &RubricService{
		rubricRepository:  rubricRepository,
		studentRepository: studentRepository,
	}
}

// AddCriterion creates and adds a new criterion to the rubric
func (service *RubricService) AddCriterion(
	id string,
	name string,
	maxScore int,
	weight int,
) error {
	rubric, err := service.rubricRepository.Get()
	if err != nil {
		return fmt.Errorf("get rubric: %w", err)
	}

	criterion, err := createCriterion(id, name, maxScore, weight)
	if err != nil {
		return err
	}

	if err := rubric.AddCriterion(criterion); err != nil {
		return fmt.Errorf("add criterion: %w", err)
	}

	if err := service.rubricRepository.Save(rubric); err != nil {
		return fmt.Errorf("save rubric: %w", err)
	}

	return nil
}

// UpdateCriterion updates an existing criterion if no evaluations reference it
func (service *RubricService) UpdateCriterion(
	id string,
	name string,
	maxScore int,
	weight int,
) error {
	rubric, err := service.rubricRepository.Get()

	if err != nil {
		return fmt.Errorf("get rubric: %w", err)
	}

	if _, exists := rubric.GetCriterion(id); !exists {
		return fmt.Errorf("criterion with ID %s does not exist", id)
	}

	hasEvaluations, err := service.studentRepository.ExistsEvaluationForCriterion(id)

	if err != nil {
		return fmt.Errorf("check criterion evaluations: %w", err)
	}

	if hasEvaluations {
		return fmt.Errorf("cannot update criterion %s: evaluations already exist", id)
	}

	criterion, err := createCriterion(id, name, maxScore, weight)

	if err != nil {
		return err
	}

	if err := rubric.UpdateCriterion(criterion); err != nil {
		return fmt.Errorf("update criterion: %w", err)
	}

	if err := service.rubricRepository.Save(rubric); err != nil {
		return fmt.Errorf("save rubric: %w", err)
	}

	return nil
}

// RemoveCriterion removes a criterion if no evaluations reference it
func (service *RubricService) RemoveCriterion(id string) error {
	rubric, err := service.rubricRepository.Get()

	if err != nil {
		return fmt.Errorf("get rubric: %w", err)
	}

	if _, exists := rubric.GetCriterion(id); !exists {
		return fmt.Errorf("criterion with ID %s does not exist", id)
	}

	hasEvaluations, err := service.studentRepository.ExistsEvaluationForCriterion(id)

	if err != nil {
		return fmt.Errorf("check criterion evaluations: %w", err)
	}

	if hasEvaluations {
		return fmt.Errorf("cannot remove criterion %s: evaluations already exist", id)
	}

	if err := rubric.RemoveCriterion(id); err != nil {
		return fmt.Errorf("remove criterion: %w", err)
	}

	if err := service.rubricRepository.Save(rubric); err != nil {
		return fmt.Errorf("save rubric: %w", err)
	}

	return nil
}

// GetRubric returns the current rubric
func (service *RubricService) GetRubric() (domain.Rubric, error) {
	rubric, err := service.rubricRepository.Get()

	if err != nil {
		return domain.Rubric{}, fmt.Errorf("get rubric: %w", err)
	}

	return rubric, nil
}

// createCriterion constructs a criterion from primitive input values
func createCriterion(
	id string,
	name string,
	maxScore int,
	weight int,
) (domain.Criterion, error) {
	maxScoreValue, err := domain.NewMaxScore(maxScore)

	if err != nil {
		return domain.Criterion{}, err
	}

	weightValue, err := domain.NewWeight(weight)

	if err != nil {
		return domain.Criterion{}, err
	}

	criterion, err := domain.NewCriterion(
		id,
		name,
		maxScoreValue,
		weightValue,
	)

	if err != nil {
		return domain.Criterion{}, err
	}

	return criterion, nil
}
