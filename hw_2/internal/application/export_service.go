package application

import (
	"fmt"
)

// ExportService executes the use case of exporting a student's final score
type ExportService struct {
	studentRepository StudentRepository
	scoreService      *ScoreService
	exporter          FinalScoreExporter
}

// NewExportService creates an export service
func NewExportService(
	studentRepository StudentRepository,
	scoreService *ScoreService,
	exporter FinalScoreExporter,
) *ExportService {
	return &ExportService{
		studentRepository: studentRepository,
		scoreService:      scoreService,
		exporter:          exporter,
	}
}

// ExportStudentScore calculates and exports a student's final score
func (service *ExportService) ExportStudentScore(
	studentID string,
	destination string,
) error {
	student, err := service.studentRepository.GetByID(studentID)
	if err != nil {
		return fmt.Errorf("get student: %w", err)
	}

	finalScore, err := service.scoreService.CalculateStudentScore(studentID)
	if err != nil {
		return fmt.Errorf("calculate student score: %w", err)
	}

	data := ScoreExportData{
		StudentID:  student.ID(),
		Name:       student.Name(),
		Surname:    student.Surname(),
		Patronymic: student.Patronymic(),
		Group:      student.Group(),
		Faculty:    student.Faculty(),
		FinalScore: finalScore.Value(),
	}

	if err := service.exporter.Export(data, destination); err != nil {
		return fmt.Errorf("export student score: %w", err)
	}

	return nil
}
