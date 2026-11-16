package container_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"hw_2/internal/container"
)

func TestNewWithCustomPaths(t *testing.T) {
	tempDir := t.TempDir()

	studentsFile := filepath.Join(tempDir, "students", "students.json")
	rubricFile := filepath.Join(tempDir, "rubric", "rubric.json")

	appContainer, err := container.New(container.Config{
		StudentsFile: studentsFile,
		RubricFile:   rubricFile,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if appContainer == nil {
		t.Fatal("expected container to be created")
	}

	if appContainer.Facade == nil {
		t.Fatal("expected application facade to be initialized")
	}

	if _, err := os.Stat(filepath.Dir(studentsFile)); err != nil {
		t.Errorf("expected students directory to exist: %v", err)
	}

	if _, err := os.Stat(filepath.Dir(rubricFile)); err != nil {
		t.Errorf("expected rubric directory to exist: %v", err)
	}

	appContainer.Shutdown()
}

func TestNewWithDefaultPaths(t *testing.T) {
	tempDir := t.TempDir()

	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("change working directory: %v", err)
	}

	defer func() {
		if err := os.Chdir(previousDir); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	}()

	appContainer, err := container.New(container.Config{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if appContainer == nil {
		t.Fatal("expected container to be created")
	}

	if appContainer.Facade == nil {
		t.Fatal("expected application facade to be initialized")
	}

	expectedDirectories := []string{
		filepath.Join(tempDir, "data"),
	}

	for _, directory := range expectedDirectories {
		if _, err := os.Stat(directory); err != nil {
			t.Errorf("expected directory %q to exist: %v", directory, err)
		}
	}

	appContainer.Shutdown()
}

func TestNewWithEmptyStudentsDirectoryPath(t *testing.T) {
	tempDir := t.TempDir()

	blockedPath := filepath.Join(tempDir, "blocked")
	if err := os.WriteFile(blockedPath, []byte("file"), 0644); err != nil {
		t.Fatalf("create blocking file: %v", err)
	}

	appContainer, err := container.New(container.Config{
		StudentsFile: filepath.Join(blockedPath, "students.json"),
		RubricFile:   filepath.Join(tempDir, "rubric.json"),
	})

	if err == nil {
		if appContainer != nil {
			appContainer.Shutdown()
		}
		t.Fatal("expected an error when creating students directory")
	}

	if !strings.Contains(err.Error(), "create students data directory") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestNewWithInvalidRubricDirectoryPath(t *testing.T) {
	tempDir := t.TempDir()

	blockedPath := filepath.Join(tempDir, "blocked")
	if err := os.WriteFile(blockedPath, []byte("file"), 0644); err != nil {
		t.Fatalf("create blocking file: %v", err)
	}

	appContainer, err := container.New(container.Config{
		StudentsFile: filepath.Join(tempDir, "students.json"),
		RubricFile:   filepath.Join(blockedPath, "rubric.json"),
	})

	if err == nil {
		if appContainer != nil {
			appContainer.Shutdown()
		}
		t.Fatal("expected an error when creating rubric directory")
	}

	if !strings.Contains(err.Error(), "create rubric data directory") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestShutdown(t *testing.T) {
	appContainer, err := container.New(container.Config{
		StudentsFile: filepath.Join(t.TempDir(), "students.json"),
		RubricFile:   filepath.Join(t.TempDir(), "rubric.json"),
	})
	if err != nil {
		t.Fatalf("create container: %v", err)
	}

	appContainer.Shutdown()
}

func TestShutdownNilContainer(t *testing.T) {
	var appContainer *container.Container

	appContainer.Shutdown()
}
