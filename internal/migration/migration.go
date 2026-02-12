package migration

import (
	"github.com/ariashabry/boilerplate-go/helpers/log"
	"github.com/ariashabry/boilerplate-go/infras"
)

// MigrationService handles database migrations for all domains
type MigrationService interface {
	RunMigrations() error
}

type MigrationServiceImpl struct {
	db     *infras.PostgresConn
	log    *log.AppLog
	models []interface{}
}

// ProvideMigrationService creates a new migration service
func ProvideMigrationService(db *infras.PostgresConn, log *log.AppLog) *MigrationServiceImpl {
	service := &MigrationServiceImpl{
		db:     db,
		log:    log,
		models: []interface{}{},
	}
	return service
}

// RunMigrations runs all registered models migration
func (m *MigrationServiceImpl) RunMigrations() error {
	if len(m.models) == 0 {
		m.log.Warn("[Migration] No models to migrate")
		return nil
	}

	m.log.Infof("[Migration] Running migrations for %d models", len(m.models))

	if err := m.db.Write.AutoMigrate(m.models...); err != nil {
		m.log.WithError(err).Error("[Migration] Failed to run migrations")
		return err
	}

	m.log.Info("[Migration] All migrations completed successfully")
	return nil
}

// RegisterModels registers models for migration
func (m *MigrationServiceImpl) RegisterModels(models ...interface{}) {
	m.models = append(m.models, models...)
}
