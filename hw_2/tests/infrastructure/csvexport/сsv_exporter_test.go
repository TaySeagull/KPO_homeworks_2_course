package csvexport_test

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"hw_2/internal/application"
	"hw_2/internal/infrastructure/csvexport"
)

func TestCSVExporterExport(t *testing.T) {
	exporter := csvexport.NewCSVExporter()
	destination := filepath.Join(t.TempDir(), "result.csv")

	data := application.ScoreExportData{
		StudentID:  "student-1",
		Name:       "Ляляля",
		Surname:    "Жужужу",
		Patronymic: "Ивановна",
		Group:      "БПИ231",
		Faculty:    "ФКН",
		FinalScore: 9,
	}

	err := exporter.Export(data, destination)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	file, err := os.Open(destination)
	if err != nil {
		t.Fatalf("expected CSV file to exist, got %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("expected valid CSV content, got %v", err)
	}

	expected := [][]string{
		{
			"student_id",
			"name",
			"surname",
			"patronymic",
			"group",
			"faculty",
			"final_score",
		},
		{
			"student-1",
			"Ляляля",
			"Жужужу",
			"Ивановна",
			"БПИ231",
			"ФКН",
			"9",
		},
	}

	if !reflect.DeepEqual(records, expected) {
		t.Errorf("unexpected CSV content:\nexpected: %v\ngot: %v", expected, records)
	}
}

func TestCSVExporterExportInvalidDestination(t *testing.T) {
	exporter := csvexport.NewCSVExporter()

	destination := filepath.Join(
		t.TempDir(),
		"missing",
		"result.csv",
	)

	data := application.ScoreExportData{
		StudentID:  "student-1",
		Name:       "Test",
		Surname:    "Student",
		FinalScore: 8,
	}

	err := exporter.Export(data, destination)
	if err == nil {
		t.Fatal("expected error when destination directory does not exist")
	}
}
