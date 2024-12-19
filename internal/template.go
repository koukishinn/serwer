package internal

import (
	"html/template"
	"path/filepath"
	"strings"
)

var (
	TemplateFunctions = template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"until": func(n int) []int {
			var result []int
			for i := 0; i < n; i++ {
				result = append(result, i)
			}
			return result
		},
		"split": func(s string, p string) []string {
			return strings.Split(s, p)
		},
		"join": func(s string, p string) string {
			return filepath.Join(s, p)
		},
	}
)
