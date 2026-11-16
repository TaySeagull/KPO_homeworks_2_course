package application_test

import (
	"errors"
	"strings"
	"testing"

	"hw_2/internal/application"
)

func TestStudentServiceAddStudent(t *testing.T) {
	repository := newInMemoryStudentRepository()
	service := application.NewStudentService(repository)

	err := service.AddStudent(
		"student-1",
		"Иван",
		"Иванов",
		"Иванович",
		"БПИ-241",
		"ФКН",
	)
	if err != nil {
		t.Fatalf("AddStudent returned unexpected error: %v", err)
	}

	student, err := repository.GetByID("student-1")
	if err != nil {
		t.Fatalf("student was not saved: %v", err)
	}

	if student.Name() != "Иван" {
		t.Errorf("expected name Иван, got %s", student.Name())
	}

	if student.Surname() != "Иванов" {
		t.Errorf("expected surname Иванов, got %s", student.Surname())
	}

	if student.Group() != "БПИ-241" {
		t.Errorf("expected group БПИ-241, got %s", student.Group())
	}
}

func TestStudentServiceAddDuplicateStudent(t *testing.T) {
	repository := newInMemoryStudentRepository()
	service := application.NewStudentService(repository)

	err := service.AddStudent(
		"student-1",
		"Иван",
		"Иванов",
		"",
		"БПИ-241",
		"ФКН",
	)
	if err != nil {
		t.Fatalf("first AddStudent returned unexpected error: %v", err)
	}

	err = service.AddStudent(
		"student-1",
		"Пётр",
		"Петров",
		"",
		"БПИ-242",
		"ФКН",
	)
	if err == nil {
		t.Fatal("expected error when adding duplicate student")
	}

	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected duplicate ID error, got: %v", err)
	}
}

func TestStudentServiceAddStudentWithInvalidData(t *testing.T) {
	repository := newInMemoryStudentRepository()
	service := application.NewStudentService(repository)

	err := service.AddStudent(
		"student-1",
		"",
		"Иванов",
		"",
		"БПИ-241",
		"ФКН",
	)
	if err == nil {
		t.Fatal("expected error for empty student name")
	}

	if len(repository.students) != 0 {
		t.Error("invalid student must not be saved")
	}
}

func TestStudentServiceGetStudent(t *testing.T) {
	repository := newInMemoryStudentRepository()
	service := application.NewStudentService(repository)

	student, err := createTestStudent("student-1")
	if err != nil {
		t.Fatalf("createTestStudent returned error: %v", err)
	}

	if err := repository.Save(student); err != nil {
		t.Fatalf("repository Save returned error: %v", err)
	}

	result, err := service.GetStudent("student-1")
	if err != nil {
		t.Fatalf("GetStudent returned unexpected error: %v", err)
	}

	if result.ID() != "student-1" {
		t.Errorf("expected student ID student-1, got %s", result.ID())
	}
}

func TestStudentServiceGetStudentNotFound(t *testing.T) {
	repository := newInMemoryStudentRepository()
	service := application.NewStudentService(repository)

	_, err := service.GetStudent("unknown")
	if err == nil {
		t.Fatal("expected error for nonexistent student")
	}
}

func TestStudentServiceGetAllStudents(t *testing.T) {
	repository := newInMemoryStudentRepository()
	service := application.NewStudentService(repository)

	first, err := createTestStudent("student-1")
	if err != nil {
		t.Fatalf("create first student: %v", err)
	}

	second, err := createTestStudent("student-2")
	if err != nil {
		t.Fatalf("create second student: %v", err)
	}

	if err := repository.Save(first); err != nil {
		t.Fatalf("save first student: %v", err)
	}

	if err := repository.Save(second); err != nil {
		t.Fatalf("save second student: %v", err)
	}

	students, err := service.GetAllStudents()
	if err != nil {
		t.Fatalf("GetAllStudents returned unexpected error: %v", err)
	}

	if len(students) != 2 {
		t.Errorf("expected 2 students, got %d", len(students))
	}
}

func TestStudentServiceRepositoryErrors(t *testing.T) {
	t.Run("exists error", func(t *testing.T) {
		repository := newInMemoryStudentRepository()
		repository.existsErr = errors.New("storage failure")
		service := application.NewStudentService(repository)

		err := service.AddStudent(
			"student-1",
			"Иван",
			"Иванов",
			"",
			"БПИ-241",
			"ФКН",
		)
		if err == nil {
			t.Fatal("expected repository error")
		}
	})

	t.Run("save error", func(t *testing.T) {
		repository := newInMemoryStudentRepository()
		repository.saveErr = errors.New("storage failure")
		service := application.NewStudentService(repository)

		err := service.AddStudent(
			"student-1",
			"Иван",
			"Иванов",
			"",
			"БПИ-241",
			"ФКН",
		)
		if err == nil {
			t.Fatal("expected repository error")
		}
	})

	t.Run("get all error", func(t *testing.T) {
		repository := newInMemoryStudentRepository()
		repository.getAllErr = errors.New("storage failure")
		service := application.NewStudentService(repository)

		_, err := service.GetAllStudents()
		if err == nil {
			t.Fatal("expected repository error")
		}
	})
}
