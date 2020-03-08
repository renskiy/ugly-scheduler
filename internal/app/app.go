package app

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib" // DB driver
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func init() {
	_ = viper.BindEnv("db_host")
	_ = viper.BindEnv("db_port")
	_ = viper.BindEnv("db_name")
	_ = viper.BindEnv("db_user")
	_ = viper.BindEnv("db_password")

	viper.SetDefault("db_driver", "pgx")

	viper.SetDefault("db_host", "localhost")
	viper.SetDefault("db_port", "5432")
	viper.SetDefault("db_name", "ugly-scheduler")
	viper.SetDefault("db_user", "postgres")
}

// App represents application interface
type App interface {
	DB() *sqlx.DB
	Logger() *zap.Logger
	Router() *mux.Router
	Run()
	RunTask(task func(shutdown <-chan struct{}) error)
}

// New returns new instance of App
func New() App {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(errors.Wrap(err, "could not init logger"))
	}

	app := &app{
		logger:   logger,
		router:   mux.NewRouter(),
		db:       sqlx.MustOpen(viper.GetString("db_driver"), DBConnectionString()),
		shutdown: make(chan struct{}),
	}

	app.router.Use(handlers.RecoveryHandler(handlers.RecoveryLogger(&recoveryLogger{logger.Sugar()})))

	return app
}

type app struct {
	logger *zap.Logger
	router *mux.Router
	db     *sqlx.DB

	tasks    sync.WaitGroup
	shutdown chan struct{}
}

// DB returns DB instance
func (app *app) DB() *sqlx.DB {
	return app.db
}

// Logger returns logger instance
func (app *app) Logger() *zap.Logger {
	return app.logger
}

// Router returns router instance
func (app *app) Router() *mux.Router {
	return app.router
}

// Run application
func (app *app) Run() {
	timeout := context.Background()

	defer app.logger.Sync()

	server := http.Server{
		Addr:    viper.GetString("addr"),
		Handler: app.router,
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		timeout, _ = context.WithTimeout(timeout, 10*time.Second)
		close(app.shutdown)
		if err := server.Shutdown(timeout); err != nil {
			app.logger.Error("can't shutdown HTTP server", zap.Error(err))
		}
	}()
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		app.logger.Fatal("can't start server", zap.Error(err))
	}

	tasksDone := make(chan struct{})
	go func() {
		defer close(tasksDone)
		app.tasks.Wait()
	}()

	select {
	case <-timeout.Done():
		app.logger.Error("graceful shutdown failed", zap.Error(timeout.Err()))
	case <-tasksDone:
		app.logger.Info("graceful shutdown completed")
	}
}

// RunTask runs background task
func (app *app) RunTask(task func(shutdown <-chan struct{}) error) {
	app.tasks.Add(1)
	go func() {
		defer app.tasks.Done()
		if err := task(app.shutdown); err != nil {
			app.logger.Error("can't perform background task", zap.Error(err))
		}
	}()
}

type recoveryLogger struct {
	*zap.SugaredLogger
}

// Println is used to comfort RecoveryHandlerLogger interface of gorilla/handlers package
func (logger *recoveryLogger) Println(args ...interface{}) {
	logger.Error(args...)
}
