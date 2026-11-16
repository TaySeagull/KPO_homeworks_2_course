package application

import "hw_2/internal/domain"

// ApplicationFacade provides a single entry point to application use cases
type ApplicationFacade struct {
	studentService    *StudentService
	rubricService     *RubricService
	evaluationService *EvaluationService
	scoreService      *ScoreService
	exportService     *ExportService
}

// NewApplicationFacade creates a facade for application services
func NewApplicationFacade(
	studentService *StudentService,
	rubricService *RubricService,
	evaluationService *EvaluationService,
	scoreService *ScoreService,
	exportService *ExportService,
) *ApplicationFacade {
	return &ApplicationFacade{
		studentService:    studentService,
		rubricService:     rubricService,
		evaluationService: evaluationService,
		scoreService:      scoreService,
		exportService:     exportService,
	}
}

// AddStudent creates a new student
func (facade *ApplicationFacade) AddStudent(
	id string,
	name string,
	surname string,
	patronymic string,
	group string,
	faculty string,
) error {
	return facade.studentService.AddStudent(
		id,
		name,
		surname,
		patronymic,
		group,
		faculty,
	)
}

// GetStudent returns a student by ID
func (facade *ApplicationFacade) GetStudent(id string) (domain.Student, error) {
	return facade.studentService.GetStudent(id)
}

// GetAllStudents returns all students
func (facade *ApplicationFacade) GetAllStudents() ([]domain.Student, error) {
	return facade.studentService.GetAllStudents()
}

// AddCriterion adds a new criterion to the rubric
func (facade *ApplicationFacade) AddCriterion(
	id string,
	name string,
	maxScore int,
	weight int,
) error {
	return facade.rubricService.AddCriterion(id, name, maxScore, weight)
}

// UpdateCriterion updates an existing criterion
func (facade *ApplicationFacade) UpdateCriterion(
	id string,
	name string,
	maxScore int,
	weight int,
) error {
	return facade.rubricService.UpdateCriterion(id, name, maxScore, weight)
}

// RemoveCriterion removes a criterion from the rubric
func (facade *ApplicationFacade) RemoveCriterion(id string) error {
	return facade.rubricService.RemoveCriterion(id)
}

// GetRubric returns the current rubric
func (facade *ApplicationFacade) GetRubric() (domain.Rubric, error) {
	return facade.rubricService.GetRubric()
}

// SetEvaluation creates or updates a student's evaluation
func (facade *ApplicationFacade) SetEvaluation(
	studentID string,
	criterionID string,
	score int,
	comment string,
) error {
	return facade.evaluationService.SetEvaluation(
		studentID,
		criterionID,
		score,
		comment,
	)
}

// CalculateStudentScore calculates a student's final score
func (facade *ApplicationFacade) CalculateStudentScore(
	studentID string,
) (domain.FinalScore, error) {
	return facade.scoreService.CalculateStudentScore(studentID)
}

// ExportStudentScore calculates and exports a student's final score
func (facade *ApplicationFacade) ExportStudentScore(
	studentID string,
	destination string,
) error {
	return facade.exportService.ExportStudentScore(studentID, destination)
}
