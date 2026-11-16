package domain

import (
	"fmt"
	"time"
)

// Student represents a student and is the aggregate root for evaluations
type Student struct {
	id          string
	name        string
	surname     string
	patronymic  string
	group       string
	faculty     string
	evaluations map[string]Evaluation
}

func NewStudent(
	id string,
	name string,
	surname string,
	patronymic string,
	group string,
	faculty string,
) (Student, error) {
	if id == "" {
		return Student{}, fmt.Errorf("student ID cannot be empty")
	}

	if name == "" {
		return Student{}, fmt.Errorf("student name cannot be empty")
	}

	if surname == "" {
		return Student{}, fmt.Errorf("student surname cannot be empty")
	}

	if group == "" {
		return Student{}, fmt.Errorf("student group cannot be empty")
	}

	if faculty == "" {
		return Student{}, fmt.Errorf("student faculty cannot be empty")
	}

	return Student{
		id:          id,
		name:        name,
		surname:     surname,
		patronymic:  patronymic,
		group:       group,
		faculty:     faculty,
		evaluations: make(map[string]Evaluation),
	}, nil
}

func (student Student) ID() string {
	return student.id
}

func (student Student) Name() string {
	return student.name
}

func (student Student) Surname() string {
	return student.surname
}

func (student Student) Patronymic() string {
	return student.patronymic
}

func (student Student) Group() string {
	return student.group
}

func (student Student) Faculty() string {
	return student.faculty
}

// SetEvaluation creates or updates an evaluation for the given criterion
func (student *Student) SetEvaluation(
	criterionID string,
	score ScoreValue,
	maxScore MaxScore,
	comment string,
) error {
	evaluation, err := NewEvaluation(
		criterionID,
		score,
		comment,
		maxScore,
	)
	if err != nil {
		return err
	}

	if student.evaluations == nil {
		student.evaluations = make(map[string]Evaluation)
	}

	student.evaluations[criterionID] = evaluation
	return nil
}

// RestoreEvaluation restores an evaluation from persistent storage
func (student *Student) RestoreEvaluation(
	criterionID string,
	score ScoreValue,
	maxScore MaxScore,
	comment string,
	updatedAt time.Time,
) error {
	if student.evaluations == nil {
		student.evaluations = make(map[string]Evaluation)
	}

	if _, exists := student.evaluations[criterionID]; exists {
		return fmt.Errorf("evaluation for criterion %s already exists", criterionID)
	}

	evaluation, err := NewEvaluationAt(
		criterionID,
		score,
		comment,
		maxScore,
		updatedAt,
	)
	if err != nil {
		return err
	}

	student.evaluations[criterionID] = evaluation
	return nil
}

// RemoveEvaluation removes an evaluation for the given criterion
func (student *Student) RemoveEvaluation(criterionID string) error {
	if _, exists := student.evaluations[criterionID]; !exists {
		return fmt.Errorf("evaluation for criterion %s does not exist", criterionID)
	}

	delete(student.evaluations, criterionID)
	return nil
}

func (student Student) GetEvaluation(criterionID string) (Evaluation, bool) {
	evaluation, exists := student.evaluations[criterionID]
	return evaluation, exists
}

func (student Student) GetEvaluations() []Evaluation {
	evaluations := make([]Evaluation, 0, len(student.evaluations))

	for _, evaluation := range student.evaluations {
		evaluations = append(evaluations, evaluation)
	}

	return evaluations
}
