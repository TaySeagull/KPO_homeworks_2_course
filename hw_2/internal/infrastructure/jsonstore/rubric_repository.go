package jsonstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"hw_2/internal/domain"
)

// JSONRubricRepository stores the Rubric aggregate in a JSON file
type JSONRubricRepository struct {
	filePath string
}

// criterionDTO represents a criterion in persistent storage
type criterionDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	MaxScore int    `json:"max_score"`
	Weight   int    `json:"weight"`
}

// NewJSONRubricRepository creates a JSON-backed rubric repository
func NewJSONRubricRepository(filePath string) *JSONRubricRepository {
	return &JSONRubricRepository{
		filePath: filePath,
	}
}

// Get returns the current rubric
func (repository *JSONRubricRepository) Get() (domain.Rubric, error) {
	data, err := os.ReadFile(repository.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return *domain.NewRubric(), nil
		}

		return domain.Rubric{}, fmt.Errorf("read rubric file: %w", err)
	}

	var storedCriteria []criterionDTO

	if err := json.Unmarshal(data, &storedCriteria); err != nil {
		return domain.Rubric{}, fmt.Errorf("decode rubric file: %w", err)
	}

	rubric := domain.NewRubric()

	for _, storedCriterion := range storedCriteria {
		maxScore, err := domain.NewMaxScore(storedCriterion.MaxScore)
		if err != nil {
			return domain.Rubric{}, fmt.Errorf(
				"restore criterion %s max score: %w",
				storedCriterion.ID,
				err,
			)
		}

		weight, err := domain.NewWeight(storedCriterion.Weight)
		if err != nil {
			return domain.Rubric{}, fmt.Errorf(
				"restore criterion %s weight: %w",
				storedCriterion.ID,
				err,
			)
		}

		criterion, err := domain.NewCriterion(
			storedCriterion.ID,
			storedCriterion.Name,
			maxScore,
			weight,
		)
		if err != nil {
			return domain.Rubric{}, fmt.Errorf(
				"restore criterion %s: %w",
				storedCriterion.ID,
				err,
			)
		}

		if err := rubric.AddCriterion(criterion); err != nil {
			return domain.Rubric{}, fmt.Errorf(
				"add restored criterion %s: %w",
				storedCriterion.ID,
				err,
			)
		}
	}

	return *rubric, nil
}

// Save stores the rubric in a JSON file
func (repository *JSONRubricRepository) Save(rubric domain.Rubric) error {
	criteria := rubric.GetCriteria()
	storedCriteria := make([]criterionDTO, 0, len(criteria))

	for _, criterion := range criteria {
		storedCriteria = append(storedCriteria, criterionDTO{
			ID:       criterion.ID(),
			Name:     criterion.Name(),
			MaxScore: criterion.MaxScore().Value(),
			Weight:   criterion.Weight().Value(),
		})
	}

	data, err := json.MarshalIndent(storedCriteria, "", "  ")
	if err != nil {
		return fmt.Errorf("encode rubric: %w", err)
	}

	directory := filepath.Dir(repository.filePath)

	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("create rubric directory: %w", err)
	}

	if err := os.WriteFile(repository.filePath, data, 0644); err != nil {
		return fmt.Errorf("write rubric file: %w", err)
	}

	return nil
}
