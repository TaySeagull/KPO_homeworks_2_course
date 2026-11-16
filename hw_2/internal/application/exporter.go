package application

// ScoreExportData contains all information required to export a student's final score
type ScoreExportData struct {
	StudentID  string
	Name       string
	Surname    string
	Patronymic string
	Group      string
	Faculty    string
	FinalScore int
}

// FinalScoreExporter defines an abstraction for exporting a student's final score
type FinalScoreExporter interface {
	Export(data ScoreExportData, destination string) error
}
