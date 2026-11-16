package csvexport

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"hw_2/internal/application"
)

// CSVExporter exports student scores to CSV files
type CSVExporter struct{}

// NewCSVExporter creates a CSV exporter
func NewCSVExporter() *CSVExporter {
	return &CSVExporter{}
}

// Export writes a student's final score to a CSV file
func (exporter *CSVExporter) Export(
	data application.ScoreExportData,
	destination string,
) error {
	file, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"student_id",
		"name",
		"surname",
		"patronymic",
		"group",
		"faculty",
		"final_score",
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("write CSV header: %w", err)
	}

	record := []string{
		data.StudentID,
		data.Name,
		data.Surname,
		data.Patronymic,
		data.Group,
		data.Faculty,
		strconv.Itoa(data.FinalScore),
	}

	if err := writer.Write(record); err != nil {
		return fmt.Errorf("write CSV record: %w", err)
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush CSV file: %w", err)
	}

	return nil
}
