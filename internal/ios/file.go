package ios

import (
	"os"
	"path/filepath"
)

type File struct {
	Name      string
	IsDir     bool
	Size      float64
	SizeOrder string
	Path      string
}

var (
	orders = map[float64]string{
		1 << 0:  "B",
		1 << 10: "KB",
		1 << 20: "MB",
		1 << 30: "GB",
		1 << 40: "TB",
	}
)

// Read a directory, giving back a list of files and subdirectories
func Read(dir string) ([]File, error) {
	var files []File

	entry, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range entry {
		info, _ := file.Info()

		size, order := compute(info.Size())

		files = append(files, File{
			Name:      file.Name(),
			IsDir:     file.IsDir(),
			Size:      size,
			SizeOrder: order,
			Path:      filepath.Join(dir, file.Name()),
		})
	}

	return files, err
}

func compute(i int64) (float64, string) {
	original := float64(i)

	size := original
	order := orders[1]

	for k, v := range orders {
		if original/k >= 1 {
			size = original / k
			order = v
		}
	}

	return size, order
}
