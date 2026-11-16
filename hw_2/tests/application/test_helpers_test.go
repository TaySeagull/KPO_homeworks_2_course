package application_test

import (
	"errors"

	"hw_2/internal/domain"
)

type inMemoryStudentRepository struct {
	students               map[string]domain.Student
	getByIDErr             error
	getAllErr              error
	saveErr                error
	existsErr              error
	existsEvaluationForErr error
}

func newInMemoryStudentRepository() *inMemoryStudentRepository {
	return &inMemoryStudentRepository{
		students: make(map[string]domain.Student),
	}
}

func (repository *inMemoryStudentRepository) GetByID(id string) (domain.Student, error) {
	if repository.getByIDErr != nil {
		return domain.Student{}, repository.getByIDErr
	}

	student, exists := repository.students[id]
	if !exists {
		return domain.Student{}, errors.New("student not found")
	}

	return student, nil
}

func (repository *inMemoryStudentRepository) GetAll() ([]domain.Student, error) {
	if repository.getAllErr != nil {
		return nil, repository.getAllErr
	}

	students := make([]domain.Student, 0, len(repository.students))

	for _, student := range repository.students {
		students = append(students, student)
	}

	return students, nil
}

func (repository *inMemoryStudentRepository) Save(student domain.Student) error {
	if repository.saveErr != nil {
		return repository.saveErr
	}

	repository.students[student.ID()] = student
	return nil
}

func (repository *inMemoryStudentRepository) ExistsByID(id string) (bool, error) {
	if repository.existsErr != nil {
		return false, repository.existsErr
	}

	_, exists := repository.students[id]
	return exists, nil
}

func (repository *inMemoryStudentRepository) ExistsEvaluationForCriterion(
	criterionID string,
) (bool, error) {
	if repository.existsEvaluationForErr != nil {
		return false, repository.existsEvaluationForErr
	}

	for _, student := range repository.students {
		if _, exists := student.GetEvaluation(criterionID); exists {
			return true, nil
		}
	}

	return false, nil
}

type inMemoryRubricRepository struct {
	rubric  domain.Rubric
	getErr  error
	saveErr error
}

func newInMemoryRubricRepository() *inMemoryRubricRepository {
	return &inMemoryRubricRepository{
		rubric: *domain.NewRubric(),
	}
}

func (repository *inMemoryRubricRepository) Get() (domain.Rubric, error) {
	if repository.getErr != nil {
		return domain.Rubric{}, repository.getErr
	}

	return repository.rubric, nil
}

func (repository *inMemoryRubricRepository) Save(rubric domain.Rubric) error {
	if repository.saveErr != nil {
		return repository.saveErr
	}

	repository.rubric = rubric
	return nil
}

func createTestCriterion(
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

	return domain.NewCriterion(id, name, maxScoreValue, weightValue)
}

func createTestStudent(id string) (domain.Student, error) {
	return domain.NewStudent(
		id,
		"Иван",
		"Иванов",
		"Иванович",
		"БПИ-241",
		"ФКН",
	)
}
