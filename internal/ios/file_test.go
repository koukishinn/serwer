package ios

import (
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	directory, err := os.MkdirTemp(os.TempDir(), "server")
	if err != nil {
		t.Fatalf("failed to create a temporary directory, reason %v", err)
	}

	files := []*os.File{}

	for i := 0; i < 10; i++ {
		file, err := os.CreateTemp(directory, "file*")
		if err != nil {
			t.Fatalf("failed to create a temporary file, reason %v", err)
		}

		files = append(files, file)
	}

	infos, err := Read(directory)

	if len(infos) != len(files) {
		t.Errorf("failed to list files correctly, want %d, got %d", len(files), len(infos))
	}

	for _, file := range infos {
		t.Log(file.Name())
	}

	t.Cleanup(func() {
		for _, file := range files {
			os.Remove(file.Name())
		}

		os.Remove(directory)
	})
}
