package app

import (
	"io"
	"strings"

	"github.com/Royal-Linux/hornero/pkg/config"
	"github.com/Royal-Linux/hornero/pkg/i18n"
	"github.com/Royal-Linux/hornero/pkg/log"
	"github.com/Royal-Linux/logrus"
)

// App struct
type App struct {
	closers []io.Closer

	Config        *config.AppConfig
	Log           *logrus.Entry
	Tr            *i18n.TranslationSet
	ErrorChan     chan error
}

// NewApp bootstrap a new application
func NewApp(config *config.AppConfig) (*App, error) {
	app := &App{
		closers:   []io.Closer{},
		Config:    config,
		ErrorChan: make(chan error),
	}
	app.Log = log.NewLogger(config, "23432119147a4367abf7c0de2aa99a2d")
	app.Tr = i18n.NewTranslationSet(app.Log)

	return app, nil
}

func (app *App) Run() error {
	err := app.Gui.RunWithSubprocesses()
	return err
}

type errorMapping struct {
	originalError string
	newError      string
}

// KnownError takes an error and tells us whether it's an error that we know about where we can print a nicely formatted version of it rather than panicking with a stack trace
func (app *App) KnownError(err error) (string, bool) {
	errorMessage := err.Error()

	mappings := []errorMapping{
		{
			originalError: "Got permission denied while trying to connect to the Docker daemon socket",
			newError:      app.Tr.ErrorOccurred,
		},
	}

	for _, mapping := range mappings {
		if strings.Contains(errorMessage, mapping.originalError) {
			return mapping.newError, true
		}
	}

	return "", false
}
