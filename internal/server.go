package internal

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/kokishin/serwer/internal/security"
)

type Server struct {
	e         *echo.Echo
	s         *http.Server
	enclave   *security.Enclave
	directory string
	sessions  map[string]bool

	logger *slog.Logger

	done chan bool
}

type ServerOpts struct {
	Logger    *slog.Logger
	Directory string
}

func NewServer(opts *ServerOpts) *Server {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogMethod:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				opts.Logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.String("method", v.Method),
					slog.Int("status", v.Status),
				)
			} else {
				opts.Logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Any("params", v.QueryParams),
					slog.String("method", v.Method),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))

	s := &http.Server{
		Addr:    "0.0.0.0:8095",
		Handler: e,
	}

	return &Server{
		e:         e,
		s:         s,
		enclave:   security.NewEnclave(),
		directory: opts.Directory,
		sessions:  make(map[string]bool),
		logger:    opts.Logger,
		done:      make(chan bool),
	}
}

func (s *Server) Start() {
	s.e.Static("/", "www")

	s.e.POST("/login", s.handleLogin)

	s.e.GET("/app", s.handleApplication, s.authentication)
	s.e.GET("/settings", s.handleSettings, s.authentication)

	// endpoints to be used by the application
	s.e.GET("/version", s.handleVersion)

	s.e.GET("/files", s.handleFiles, s.authentication)
	s.e.GET("/files/*", s.handleFiles, s.authentication)
	s.e.GET("/preview/*", s.handlePreview, s.authentication)
	s.e.GET("/raw/*", s.handleRaw, s.authentication)

	s.logger.Info("starting server", slog.String("address", s.s.Addr))

	if err := s.s.ListenAndServe(); err != http.ErrServerClosed {
		s.done <- true
		s.logger.Error("failed to start server", slog.Any("error", err))
	}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) Done() <-chan bool {
	return s.done
}

func (s *Server) authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session")
		return next(c)
		if err != nil || !s.sessions[cookie.Value] {
			return c.Redirect(http.StatusMovedPermanently, "/")
		}

		return next(c)
	}
}
