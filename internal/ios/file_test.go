package ios

import (
	"os"
	"testing"
)

func TestCompute(t *testing.T) {
	tests := []struct {
		Size         int64
		ExpectedSize int64
		Order        string
	}{
		{Size: 4096, Order: "KB", ExpectedSize: 4},
		{Size: 8192, Order: "KB", ExpectedSize: 8},
	}

	for _, test := range tests {
		size, order := compute(test.Size)

		if order != test.Order {
			t.Errorf("want %s order, got %s", test.Order, order)
		}
		if int64(size) != test.ExpectedSize {
			t.Errorf("wrong expected size, want %d, got %f", test.ExpectedSize, size)
		}
	}
}

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
		t.Log(file.Name)
	}

	t.Cleanup(func() {
		for _, file := range files {
			os.Remove(file.Name())
		}

		os.Remove(directory)
	})
}
