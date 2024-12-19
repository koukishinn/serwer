package ios

import (
	"os"
	"path/filepath"
)

type File struct {
	Name  string
	IsDir bool
	Size  int64
	Path  string
}

// Read a directory, giving back
func Read(dir string) ([]File, error) {
	var files []File

	entry, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range entry {
		info, _ := file.Info()
		files = append(files, File{
			Name:  file.Name(),
			IsDir: file.IsDir(),
			Size:  info.Size(),
			Path:  filepath.Join(dir, file.Name()),
		})
	}

	return files, err
}
