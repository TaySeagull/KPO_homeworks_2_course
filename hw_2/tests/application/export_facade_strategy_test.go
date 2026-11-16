package application_test

import (
	"errors"
	"testing"

	"hw_2/internal/application"
	"hw_2/internal/domain"
)

type fakeFinalScoreExporter struct {
	data        application.ScoreExportData
	destination string
	err         error
	called      bool
}

func (exporter *fakeFinalScoreExporter) Export(
	data application.ScoreExportData,
	destination string,
) error {
	exporter.called = true
	exporter.data = data
	exporter.destination = destination

	return exporter.err
}

func createExportTestServices() (
	*application.StudentService,
	*application.RubricService,
	*application.EvaluationService,
	*application.ScoreService,
	*inMemoryStudentRepository,
	*inMemoryRubricRepository,
) {
	studentRepository := newInMemoryStudentRepository()
	rubricRepository := newInMemoryRubricRepository()

	studentService := application.NewStudentService(studentRepository)
	rubricService := application.NewRubricService(
		rubricRepository,
		studentRepository,
	)
	evaluationService := application.NewEvaluationService(
		studentRepository,
		rubricRepository,
	)
	scoreService := application.NewScoreService(
		studentRepository,
		rubricRepository,
	)

	return studentService, rubricService, evaluationService,
		scoreService, studentRepository, rubricRepository
}

func TestExportServiceExportStudentScore(t *testing.T) {
	studentService, rubricService, evaluationService, scoreService,
		studentRepository, _ := createExportTestServices()

	if err := studentService.AddStudent(
		"student-1",
		"Траляля",
		"Жужужу",
		"Ляляля",
		"БПИ-241",
		"ФКН",
	); err != nil {
		t.Fatalf("AddStudent() error = %v", err)
	}

	if err := rubricService.AddCriterion(
		"criterion-1",
		"Контрольная работа",
		10,
		100,
	); err != nil {
		t.Fatalf("AddCriterion() error = %v", err)
	}

	if err := evaluationService.SetEvaluation(
		"student-1",
		"criterion-1",
		8,
		"Хорошая работа",
	); err != nil {
		t.Fatalf("SetEvaluation() error = %v", err)
	}

	exporter := &fakeFinalScoreExporter{}
	exportService := application.NewExportService(
		studentRepository,
		scoreService,
		exporter,
	)

	err := exportService.ExportStudentScore("student-1", "result.csv")
	if err != nil {
		t.Fatalf("ExportStudentScore() error = %v", err)
	}

	if !exporter.called {
		t.Fatal("Export() was not called")
	}

	if exporter.destination != "result.csv" {
		t.Errorf(
			"destination = %q, want %q",
			exporter.destination,
			"result.csv",
		)
	}

	if exporter.data.StudentID != "student-1" {
		t.Errorf(
			"StudentID = %q, want %q",
			exporter.data.StudentID,
			"student-1",
		)
	}

	if exporter.data.Name != "Траляля" {
		t.Errorf("Name = %q, want %q", exporter.data.Name, "Траляля")
	}

	if exporter.data.Surname != "Жужужу" {
		t.Errorf("Surname = %q, want %q", exporter.data.Surname, "Жужужу")
	}

	if exporter.data.Patronymic != "Ляляля" {
		t.Errorf("Patronymic = %q, want %q", exporter.data.Patronymic, "Ляляля")
	}

	if exporter.data.Group != "БПИ-241" {
		t.Errorf("Group = %q, want %q", exporter.data.Group, "БПИ-241")
	}

	if exporter.data.Faculty != "ФКН" {
		t.Errorf("Faculty = %q, want %q", exporter.data.Faculty, "ФКН")
	}

	if exporter.data.FinalScore != 8 {
		t.Errorf("FinalScore = %d, want 8", exporter.data.FinalScore)
	}
}

func TestExportServiceStudentNotFound(t *testing.T) {
	studentRepository := newInMemoryStudentRepository()
	rubricRepository := newInMemoryRubricRepository()

	expectedErr := errors.New("repository error")
	studentRepository.getByIDErr = expectedErr

	scoreService := application.NewScoreService(
		studentRepository,
		rubricRepository,
	)
	exporter := &fakeFinalScoreExporter{}
	exportService := application.NewExportService(
		studentRepository,
		scoreService,
		exporter,
	)

	err := exportService.ExportStudentScore("unknown", "result.csv")

	if !errors.Is(err, expectedErr) {
		t.Errorf("error = %v, want wrapped %v", err, expectedErr)
	}

	if exporter.called {
		t.Error("Export() must not be called when student retrieval fails")
	}
}

func TestExportServiceScoreCalculationError(t *testing.T) {
	studentRepository := newInMemoryStudentRepository()
	rubricRepository := newInMemoryRubricRepository()

	student, err := domain.NewStudent(
		"student-1",
		"Траляля",
		"Жужужу",
		"Ляляля",
		"БПИ-241",
		"ФКН",
	)
	if err != nil {
		t.Fatalf("NewStudent() error = %v", err)
	}

	if err := studentRepository.Save(student); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	expectedErr := errors.New("rubric repository error")
	rubricRepository.getErr = expectedErr

	scoreService := application.NewScoreService(
		studentRepository,
		rubricRepository,
	)
	exporter := &fakeFinalScoreExporter{}
	exportService := application.NewExportService(
		studentRepository,
		scoreService,
		exporter,
	)

	err = exportService.ExportStudentScore("student-1", "result.csv")

	if !errors.Is(err, expectedErr) {
		t.Errorf("error = %v, want wrapped %v", err, expectedErr)
	}

	if exporter.called {
		t.Error("Export() must not be called when score calculation fails")
	}
}

func TestExportServiceExporterError(t *testing.T) {
	studentRepository := newInMemoryStudentRepository()
	rubricRepository := newInMemoryRubricRepository()

	student, err := domain.NewStudent(
		"student-1",
		"Траляля",
		"Жужужу",
		"Ляляля",
		"БПИ-241",
		"ФКН",
	)
	if err != nil {
		t.Fatalf("NewStudent() error = %v", err)
	}

	if err := studentRepository.Save(student); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	rubricService := application.NewRubricService(
		rubricRepository,
		studentRepository,
	)
	if err := rubricService.AddCriterion(
		"criterion-1",
		"Контрольная работа",
		10,
		100,
	); err != nil {
		t.Fatalf("AddCriterion() error = %v", err)
	}

	scoreService := application.NewScoreService(
		studentRepository,
		rubricRepository,
	)

	expectedErr := errors.New("export error")
	exporter := &fakeFinalScoreExporter{err: expectedErr}
	exportService := application.NewExportService(
		studentRepository,
		scoreService,
		exporter,
	)

	err = exportService.ExportStudentScore("student-1", "result.csv")

	if !errors.Is(err, expectedErr) {
		t.Errorf("error = %v, want wrapped %v", err, expectedErr)
	}

	if !exporter.called {
		t.Error("Export() was not called")
	}
}

func TestApplicationFacadeWorkflow(t *testing.T) {
	studentService, rubricService, evaluationService, scoreService,
		studentRepository, _ := createExportTestServices()

	exporter := &fakeFinalScoreExporter{}
	exportService := application.NewExportService(
		studentRepository,
		scoreService,
		exporter,
	)

	facade := application.NewApplicationFacade(
		studentService,
		rubricService,
		evaluationService,
		scoreService,
		exportService,
	)

	if err := facade.AddStudent(
		"student-1",
		"Траляля",
		"Жужужу",
		"Ляляля",
		"БПИ-241",
		"ФКН",
	); err != nil {
		t.Fatalf("AddStudent() error = %v", err)
	}

	student, err := facade.GetStudent("student-1")
	if err != nil {
		t.Fatalf("GetStudent() error = %v", err)
	}

	if student.Name() != "Траляля" {
		t.Errorf("student name = %q, want %q", student.Name(), "Траляля")
	}

	students, err := facade.GetAllStudents()
	if err != nil {
		t.Fatalf("GetAllStudents() error = %v", err)
	}

	if len(students) != 1 {
		t.Errorf("students count = %d, want 1", len(students))
	}

	if err := facade.AddCriterion(
		"criterion-1",
		"Контрольная работа",
		10,
		100,
	); err != nil {
		t.Fatalf("AddCriterion() error = %v", err)
	}

	if err := facade.UpdateCriterion(
		"criterion-1",
		"Итоговая контрольная",
		10,
		100,
	); err != nil {
		t.Fatalf("UpdateCriterion() error = %v", err)
	}

	rubric, err := facade.GetRubric()
	if err != nil {
		t.Fatalf("GetRubric() error = %v", err)
	}

	if len(rubric.GetCriteria()) != 1 {
		t.Errorf("criteria count = %d, want 1", len(rubric.GetCriteria()))
	}

	if err := facade.SetEvaluation(
		"student-1",
		"criterion-1",
		8,
		"Хорошая работа",
	); err != nil {
		t.Fatalf("SetEvaluation() error = %v", err)
	}

	finalScore, err := facade.CalculateStudentScore("student-1")
	if err != nil {
		t.Fatalf("CalculateStudentScore() error = %v", err)
	}

	if finalScore.Value() != 8 {
		t.Errorf("final score = %d, want 8", finalScore.Value())
	}

	if err := facade.ExportStudentScore("student-1", "result.csv"); err != nil {
		t.Fatalf("ExportStudentScore() error = %v", err)
	}

	if !exporter.called {
		t.Error("ExportStudentScore() did not call exporter")
	}

	if err := facade.RemoveCriterion("criterion-1"); err == nil {
		t.Error("RemoveCriterion() should fail for a criterion with evaluations")
	}

	if err := facade.AddCriterion(
		"criterion-2",
		"Дополнительное задание",
		10,
		100,
	); err != nil {
		t.Fatalf("AddCriterion() for second criterion error = %v", err)
	}

	if err := facade.RemoveCriterion("criterion-2"); err != nil {
		t.Errorf("RemoveCriterion() error = %v", err)
	}
}

func TestWeightedScoreStrategyDelegatesCalculation(t *testing.T) {
	strategy := application.WeightedScoreStrategy{}

	expectedScore, expectedErr := domain.CalculateFinalScore(nil, nil)
	actualScore, actualErr := strategy.Calculate(nil, nil)

	if (actualErr == nil) != (expectedErr == nil) {
		t.Fatalf(
			"Calculate() error = %v, expected error = %v",
			actualErr,
			expectedErr,
		)
	}

	if actualErr != nil && actualErr.Error() != expectedErr.Error() {
		t.Errorf(
			"Calculate() error = %v, want %v",
			actualErr,
			expectedErr,
		)
	}

	if actualErr == nil && actualScore.Value() != expectedScore.Value() {
		t.Errorf(
			"Calculate() score = %d, want %d",
			actualScore.Value(),
			expectedScore.Value(),
		)
	}
}
