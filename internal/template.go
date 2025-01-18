package internal

import (
	"html/template"
	"math"
	"path/filepath"
	"strings"
)

var (
	TemplateFunctions = template.FuncMap{
		"add":      add,
		"sub":      sub,
		"until":    until,
		"split":    split,
		"join":     join,
		"truncate": truncate,
	}
)

type TemplatePath string

const (
	TemplateFiles   TemplatePath = "www/templates/files.tmpl"
	TemplatePreview TemplatePath = "www/templates/preview.tmpl"
)

func (t TemplatePath) Path() string {
	return string(t)
}

var (
	add = func(a, b int) int { return a + b }
	sub = func(a, b int) int { return a - b }

	until = func(n int) []int {
		var result []int
		for i := 0; i < n; i++ {
			result = append(result, i)
		}
		return result
	}

	split = func(s string, p string) []string {
		return strings.Split(s, p)
	}

	join = func(s string, p string) string {
		return filepath.Join(s, p)
	}

	truncate = func(n float64) float64 {
		return float64(int(math.Floor(n*1000)) / 1000)
	}
)
