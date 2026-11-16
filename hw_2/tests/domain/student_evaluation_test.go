package domain_test

import (
	"testing"
	"time"

	"hw_2/internal/domain"
)

func createStudent(t *testing.T) domain.Student {
	t.Helper()

	student, err := domain.NewStudent(
		"student-1",
		"Tatyana",
		"Zharkova",
		"",
		"БПИ-231",
		"ФКН",
	)
	if err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	return student
}

func TestNewStudent(t *testing.T) {
	_, err := domain.NewStudent(
		"1",
		"Anna",
		"Ivanova",
		"",
		"БПИ-231",
		"ФКН",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewStudentWithEmptyRequiredFields(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		firstName string
		surname   string
		group     string
		faculty   string
	}{
		{
			name:      "empty ID",
			id:        "",
			firstName: "Anna",
			surname:   "Ivanova",
			group:     "БПИ-231",
			faculty:   "ФКН",
		},
		{
			name:      "empty name",
			id:        "1",
			firstName: "",
			surname:   "Ivanova",
			group:     "БПИ-231",
			faculty:   "ФКН",
		},
		{
			name:      "empty surname",
			id:        "1",
			firstName: "Anna",
			surname:   "",
			group:     "БПИ-231",
			faculty:   "ФКН",
		},
		{
			name:      "empty group",
			id:        "1",
			firstName: "Anna",
			surname:   "Ivanova",
			group:     "",
			faculty:   "ФКН",
		},
		{
			name:      "empty faculty",
			id:        "1",
			firstName: "Anna",
			surname:   "Ivanova",
			group:     "БПИ-231",
			faculty:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := domain.NewStudent(
				test.id,
				test.firstName,
				test.surname,
				"",
				test.group,
				test.faculty,
			)

			if err == nil {
				t.Fatal("expected error for empty required field")
			}
		})
	}
}

func TestStudentSetEvaluation(t *testing.T) {
	student := createStudent(t)

	score, err := domain.NewScoreValue(8)
	if err != nil {
		t.Fatalf("failed to create score: %v", err)
	}

	maxScore, err := domain.NewMaxScore(10)
	if err != nil {
		t.Fatalf("failed to create maximum score: %v", err)
	}

	if err := student.SetEvaluation("criterion-1", score, maxScore, "Good"); err != nil {
		t.Fatalf("failed to set evaluation: %v", err)
	}

	evaluation, exists := student.GetEvaluation("criterion-1")
	if !exists {
		t.Fatal("evaluation was not saved")
	}

	if evaluation.Score().Value() != 8 {
		t.Errorf("got score %d, want 8", evaluation.Score().Value())
	}

	if evaluation.Comment() != "Good" {
		t.Errorf("got comment %q, want %q", evaluation.Comment(), "Good")
	}
}

func TestStudentSetEvaluationAboveMaximum(t *testing.T) {
	student := createStudent(t)

	score, _ := domain.NewScoreValue(11)
	maxScore, _ := domain.NewMaxScore(10)

	err := student.SetEvaluation("criterion-1", score, maxScore, "")

	if err == nil {
		t.Fatal("expected error for score above maximum")
	}
}

func TestStudentSetEvaluationUpdatesExisting(t *testing.T) {
	student := createStudent(t)

	maxScore, _ := domain.NewMaxScore(10)

	firstScore, _ := domain.NewScoreValue(5)
	if err := student.SetEvaluation("criterion-1", firstScore, maxScore, "First"); err != nil {
		t.Fatalf("failed to set first evaluation: %v", err)
	}

	secondScore, _ := domain.NewScoreValue(9)
	if err := student.SetEvaluation("criterion-1", secondScore, maxScore, "Updated"); err != nil {
		t.Fatalf("failed to update evaluation: %v", err)
	}

	evaluation, exists := student.GetEvaluation("criterion-1")
	if !exists {
		t.Fatal("evaluation does not exist")
	}

	if evaluation.Score().Value() != 9 {
		t.Errorf("got score %d, want 9", evaluation.Score().Value())
	}

	if evaluation.Comment() != "Updated" {
		t.Errorf("got comment %q, want %q", evaluation.Comment(), "Updated")
	}
}

func TestStudentRestoreEvaluation(t *testing.T) {
	student := createStudent(t)

	score, _ := domain.NewScoreValue(7)
	maxScore, _ := domain.NewMaxScore(10)
	updatedAt := time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)

	err := student.RestoreEvaluation(
		"criterion-1",
		score,
		maxScore,
		"Restored",
		updatedAt,
	)
	if err != nil {
		t.Fatalf("failed to restore evaluation: %v", err)
	}

	evaluation, exists := student.GetEvaluation("criterion-1")
	if !exists {
		t.Fatal("evaluation was not restored")
	}

	if !evaluation.UpdatedAt().Equal(updatedAt) {
		t.Errorf("got timestamp %v, want %v", evaluation.UpdatedAt(), updatedAt)
	}
}

func TestStudentRestoreDuplicateEvaluation(t *testing.T) {
	student := createStudent(t)

	score, _ := domain.NewScoreValue(7)
	maxScore, _ := domain.NewMaxScore(10)
	updatedAt := time.Now()

	if err := student.RestoreEvaluation(
		"criterion-1",
		score,
		maxScore,
		"",
		updatedAt,
	); err != nil {
		t.Fatalf("failed to restore first evaluation: %v", err)
	}

	if err := student.RestoreEvaluation(
		"criterion-1",
		score,
		maxScore,
		"",
		updatedAt,
	); err == nil {
		t.Fatal("expected error when restoring duplicate evaluation")
	}
}

func TestStudentRemoveEvaluation(t *testing.T) {
	student := createStudent(t)

	score, _ := domain.NewScoreValue(7)
	maxScore, _ := domain.NewMaxScore(10)

	if err := student.SetEvaluation("criterion-1", score, maxScore, ""); err != nil {
		t.Fatalf("failed to set evaluation: %v", err)
	}

	if err := student.RemoveEvaluation("criterion-1"); err != nil {
		t.Fatalf("failed to remove evaluation: %v", err)
	}

	if _, exists := student.GetEvaluation("criterion-1"); exists {
		t.Fatal("evaluation still exists after removal")
	}
}

func TestStudentRemoveNonexistentEvaluation(t *testing.T) {
	student := createStudent(t)

	if err := student.RemoveEvaluation("unknown"); err == nil {
		t.Fatal("expected error when removing nonexistent evaluation")
	}
}
