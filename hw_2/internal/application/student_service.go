package application

import (
	"fmt"
	"hw_2/internal/domain"
)

// StudentService handles application operations related to students
type StudentService struct {
	repository StudentRepository
}

func NewStudentService(repository StudentRepository) *StudentService {
	return &StudentService{
		repository: repository,
	}
}

// AddStudent creates and saves a new student
func (service *StudentService) AddStudent(
	id string,
	name string,
	surname string,
	patronymic string,
	group string,
	faculty string,
) error {
	exists, err := service.repository.ExistsByID(id)
	if err != nil {
		return fmt.Errorf("check student existence: %w", err)
	}

	if exists {
		return fmt.Errorf("student with ID %s already exists", id)
	}

	student, err := domain.NewStudent(
		id,
		name,
		surname,
		patronymic,
		group,
		faculty,
	)

	if err != nil {
		return fmt.Errorf("create student: %w", err)
	}

	if err := service.repository.Save(student); err != nil {
		return fmt.Errorf("save student: %w", err)
	}

	return nil
}

// GetStudent returns a student by ID
func (service *StudentService) GetStudent(id string) (domain.Student, error) {
	student, err := service.repository.GetByID(id)

	if err != nil {
		return domain.Student{}, fmt.Errorf("get student: %w", err)
	}

	return student, nil
}

// GetAllStudents returns all students
func (service *StudentService) GetAllStudents() ([]domain.Student, error) {
	students, err := service.repository.GetAll()

	if err != nil {
		return nil, fmt.Errorf("get students: %w", err)
	}

	return students, nil
}
