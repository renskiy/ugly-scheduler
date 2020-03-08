package test

import (
	"github.com/DATA-DOG/go-txdb"
	"github.com/renskiy/ugly-scheduler/internal/app"
	"github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"sync"
)

var once sync.Once

func init() {
	_ = viper.BindEnv("db_name_test")
	viper.SetDefault("db_name_test", "ugly-scheduler-test")
}

// Suite is a base test suite
type Suite struct {
	suite.Suite
	App app.App

	dbDriver string
}

// SetupSuite prepares global context for tests
func (s *Suite) SetupSuite() {
	viper.Set("db_name", viper.GetString("db_name_test"))
	s.dbDriver = viper.GetString("db_driver")

	db := app.New().DB()
	db.MustExec("DROP SCHEMA IF EXISTS public CASCADE")
	db.MustExec("CREATE SCHEMA public")
	migrations := &migrate.FileMigrationSource{Dir: "../../migrations"}
	if _, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up); err != nil {
		s.Fail(err.Error())
	}
	s.NoError(db.Close())

	// register txdb driver as postgres to allow sqlx properly rebind queries
	once.Do(func() {
		txdb.Register("postgres", s.dbDriver, app.DBConnectionString())
	})

	viper.Set("db_driver", "postgres")
}

// TearDownSuite clears context after tests
func (s *Suite) TearDownSuite() {
	viper.Set("db_driver", s.dbDriver)
}

// SetupTest prepares context for test
func (s *Suite) SetupTest() {
	s.App = app.New()
}

// TearDownTest clears local context after test
func (s *Suite) TearDownTest() {
	s.NoError(s.App.DB().Close())
}
