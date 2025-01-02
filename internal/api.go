package internal

import (
	"html/template"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.com/kokishin/serwer/internal/functional"
	"gitlab.com/kokishin/serwer/internal/ios"
)

func (s *Server) handleLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "admin" && password == "kokishin" {
		id, err := s.enclave.Nonce()
		if err != nil {
			return c.HTML(http.StatusInternalServerError, "500")
		}

		s.sessions[id] = true
		c.SetCookie(&http.Cookie{
			Name:     "session",
			Value:    id,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		return c.HTML(http.StatusOK, `<script>window.location.href='/app';</script>`)
	}

	return c.String(http.StatusUnauthorized, "credenciais inválidas")
}

func (s *Server) handleApplication(c echo.Context) error {
	return c.File("www/application.html")
}

func (s *Server) handleSettings(c echo.Context) error {
	return c.File("www/settings.html")
}

func (s *Server) handleVersion(c echo.Context) error {
	return c.String(http.StatusOK, version)
}

func (s *Server) handleFiles(c echo.Context) error {
	path := c.Param("*")

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	if !strings.HasPrefix(path, s.directory) {
		s.logger.Warn("", slog.String("old-path", path), slog.String("prefix", s.directory))
		path = filepath.Join(s.directory, path)
	}

	const items = 30

	parameter := c.QueryParam("page")
	page, err := strconv.Atoi(parameter)
	if err != nil || page < 1 {
		page = 1
	}

	s.logger.Info("reading directory", slog.String("path", path))

	files, err := ios.Read(path)
	if err != nil {
		s.logger.Error("failed to read directory", slog.Any("error", err))
		return c.String(http.StatusInternalServerError, "unable to read directory")
	}

	functional.ForEach(files, func(e ios.File) {
		s.logger.Debug("", slog.String("name", e.Name), slog.String("path", e.Path))
	})

	totalFiles := len(files)
	totalPages := (totalFiles + items - 1) / items

	startIndex := (page - 1) * items
	endIndex := startIndex + items

	if endIndex > totalFiles {
		endIndex = totalFiles
	}
	pagedFiles := files[startIndex:endIndex]

	data := map[string]interface{}{
		"Files":       pagedFiles,
		"CurrentPage": page,
		"CurrentPath": path,
		"TotalPages":  totalPages,
	}

	tmpl := template.Must(
		template.New("files").Funcs(TemplateFunctions).ParseFiles(TemplateFiles.Path()),
	)

	return tmpl.Execute(c.Response().Writer, data)
}

func (s *Server) handlePreview(c echo.Context) error {
	path := c.Param("*")
	if !strings.HasPrefix(path, s.directory) {
		path = filepath.Join(s.directory, path)
	}

	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return c.NoContent(http.StatusNotFound)
	}

	data := map[string]interface{}{
		"Path":     filepath.Dir(path),
		"FileName": info.Name(),
		"FilePath": path,
		"FileType": mime.TypeByExtension(filepath.Ext(path)),
	}

	tmpl := template.Must(
		template.New("preview").Funcs(TemplateFunctions).ParseFiles(TemplatePreview.Path()),
	)

	return tmpl.Execute(c.Response().Writer, data)
}

func (s *Server) handleRaw(c echo.Context) error {
	path := c.Param("*")
	if !strings.HasPrefix(path, s.directory) {
		path = filepath.Join(s.directory, path)
	}

	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return c.NoContent(http.StatusNotFound)
	}

	return c.File(path)
}
