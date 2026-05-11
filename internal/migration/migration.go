package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"doctor-go/internal/config"
	"doctor-go/migrations"
)

func Up(cfg config.MySQLConfig) error {
	m, err := newMigrator(cfg)
	if err != nil {
		return err
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func Down(cfg config.MySQLConfig, steps int) error {
	m, err := newMigrator(cfg)
	if err != nil {
		return err
	}
	if steps <= 0 {
		return errors.New("steps must be greater than 0")
	}
	err = m.Steps(-steps)
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func Version(cfg config.MySQLConfig) (uint, bool, error) {
	m, err := newMigrator(cfg)
	if err != nil {
		return 0, false, err
	}
	return m.Version()
}

func Force(cfg config.MySQLConfig, version int) error {
	m, err := newMigrator(cfg)
	if err != nil {
		return err
	}
	return m.Force(version)
}

func newMigrator(cfg config.MySQLConfig) (*migrate.Migrate, error) {
	db, err := sql.Open("mysql", withMultiStatements(cfg.DSN))
	if err != nil {
		return nil, err
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	instance, err := migrate.NewWithInstance("iofs", source, "mysql", driver)
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return instance, nil
}

func withMultiStatements(dsn string) string {
	if strings.Contains(dsn, "multiStatements=") {
		return dsn
	}
	if strings.Contains(dsn, "?") {
		return dsn + "&multiStatements=true"
	}
	return dsn + "?multiStatements=true"
}

func PrintVersion(cfg config.MySQLConfig) error {
	version, dirty, err := Version(cfg)
	if errors.Is(err, migrate.ErrNilVersion) {
		fmt.Fprintln(os.Stdout, "version: none")
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "version: %d dirty: %v\n", version, dirty)
	return nil
}
