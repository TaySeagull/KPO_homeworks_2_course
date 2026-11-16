package jsonstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"hw_2/internal/application"
	"hw_2/internal/domain"
)

// JSONStudentRepository stores Student aggregates in a JSON file
type JSONStudentRepository struct {
	filePath         string
	rubricRepository application.RubricRepository
}

// studentDTO represents a student in persistent storage
type studentDTO struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Surname     string          `json:"surname"`
	Patronymic  string          `json:"patronymic"`
	Group       string          `json:"group"`
	Faculty     string          `json:"faculty"`
	Evaluations []evaluationDTO `json:"evaluations"`
}

// evaluationDTO represents an evaluation in persistent storage
type evaluationDTO struct {
	CriterionID string `json:"criterion_id"`
	Score       int    `json:"score"`
	Comment     string `json:"comment"`
	UpdatedAt   string `json:"updated_at"`
}

// NewJSONStudentRepository creates a JSON-backed student repository
func NewJSONStudentRepository(
	filePath string,
	rubricRepository application.RubricRepository,
) *JSONStudentRepository {
	return &JSONStudentRepository{
		filePath:         filePath,
		rubricRepository: rubricRepository,
	}
}

// GetByID returns a student by ID
func (repository *JSONStudentRepository) GetByID(id string) (domain.Student, error) {
	students, err := repository.load()
	if err != nil {
		return domain.Student{}, err
	}

	for _, student := range students {
		if student.ID() == id {
			return student, nil
		}
	}

	return domain.Student{}, fmt.Errorf("student with ID %s not found", id)
}

// GetAll returns all students
func (repository *JSONStudentRepository) GetAll() ([]domain.Student, error) {
	return repository.load()
}

// Save creates or updates a student
func (repository *JSONStudentRepository) Save(student domain.Student) error {
	students, err := repository.load()
	if err != nil {
		return err
	}

	found := false

	for i := range students {
		if students[i].ID() == student.ID() {
			students[i] = student
			found = true
			break
		}
	}

	if !found {
		students = append(students, student)
	}

	return repository.write(students)
}

// ExistsByID checks whether a student with the given ID exists
func (repository *JSONStudentRepository) ExistsByID(id string) (bool, error) {
	students, err := repository.load()
	if err != nil {
		return false, err
	}

	for _, student := range students {
		if student.ID() == id {
			return true, nil
		}
	}

	return false, nil
}

// ExistsEvaluationForCriterion checks whether any student has an evaluation for the criterion
func (repository *JSONStudentRepository) ExistsEvaluationForCriterion(
	criterionID string,
) (bool, error) {
	students, err := repository.load()
	if err != nil {
		return false, err
	}

	for _, student := range students {
		if _, exists := student.GetEvaluation(criterionID); exists {
			return true, nil
		}
	}

	return false, nil
}

// load reads students from the JSON file and restores their domain objects
func (repository *JSONStudentRepository) load() ([]domain.Student, error) {
	data, err := os.ReadFile(repository.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.Student{}, nil
		}

		return nil, fmt.Errorf("read students file: %w", err)
	}

	var storedStudents []studentDTO

	if err := json.Unmarshal(data, &storedStudents); err != nil {
		return nil, fmt.Errorf("decode students file: %w", err)
	}

	rubric, err := repository.rubricRepository.Get()
	if err != nil {
		return nil, fmt.Errorf("load rubric: %w", err)
	}

	students := make([]domain.Student, 0, len(storedStudents))

	for _, storedStudent := range storedStudents {
		student, err := domain.NewStudent(
			storedStudent.ID,
			storedStudent.Name,
			storedStudent.Surname,
			storedStudent.Patronymic,
			storedStudent.Group,
			storedStudent.Faculty,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"restore student %s: %w",
				storedStudent.ID,
				err,
			)
		}

		for _, storedEvaluation := range storedStudent.Evaluations {
			criterion, exists := rubric.GetCriterion(storedEvaluation.CriterionID)
			if !exists {
				return nil, fmt.Errorf(
					"restore evaluation for student %s: criterion %s does not exist",
					storedStudent.ID,
					storedEvaluation.CriterionID,
				)
			}

			score, err := domain.NewScoreValue(storedEvaluation.Score)
			if err != nil {
				return nil, fmt.Errorf(
					"restore evaluation for student %s: %w",
					storedStudent.ID,
					err,
				)
			}

			updatedAt, err := parseTime(storedEvaluation.UpdatedAt)
			if err != nil {
				return nil, fmt.Errorf(
					"restore evaluation for student %s: %w",
					storedStudent.ID,
					err,
				)
			}

			err = student.RestoreEvaluation(
				storedEvaluation.CriterionID,
				score,
				criterion.MaxScore(),
				storedEvaluation.Comment,
				updatedAt,
			)
			if err != nil {
				return nil, fmt.Errorf(
					"restore evaluation for student %s: %w",
					storedStudent.ID,
					err,
				)
			}
		}

		students = append(students, student)
	}

	return students, nil
}

// write saves students to the JSON file
func (repository *JSONStudentRepository) write(students []domain.Student) error {
	storedStudents := make([]studentDTO, 0, len(students))

	for _, student := range students {
		storedStudent := studentDTO{
			ID:          student.ID(),
			Name:        student.Name(),
			Surname:     student.Surname(),
			Patronymic:  student.Patronymic(),
			Group:       student.Group(),
			Faculty:     student.Faculty(),
			Evaluations: make([]evaluationDTO, 0),
		}

		for _, evaluation := range student.GetEvaluations() {
			storedStudent.Evaluations = append(
				storedStudent.Evaluations,
				evaluationDTO{
					CriterionID: evaluation.CriterionID(),
					Score:       evaluation.Score().Value(),
					Comment:     evaluation.Comment(),
					UpdatedAt:   evaluation.UpdatedAt().Format("2006-01-02T15:04:05.999999999Z07:00"),
				},
			)
		}

		storedStudents = append(storedStudents, storedStudent)
	}

	data, err := json.MarshalIndent(storedStudents, "", "  ")
	if err != nil {
		return fmt.Errorf("encode students: %w", err)
	}

	directory := filepath.Dir(repository.filePath)

	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("create students directory: %w", err)
	}

	if err := os.WriteFile(repository.filePath, data, 0644); err != nil {
		return fmt.Errorf("write students file: %w", err)
	}

	return nil
}

// parseTime restores a timestamp from JSON
func parseTime(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid evaluation timestamp %q: %w", value, err)
	}

	return parsed, nil
}
