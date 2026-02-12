package model

import (
	"github.com/ariashabry/boilerplate-go/internal/domains/product/model/dto"
	"github.com/ariashabry/boilerplate-go/internal/migration"
)

// RegisterProductModels registers product domain models for migration
func RegisterProductModels(migrationService *migration.MigrationServiceImpl) {
	migrationService.RegisterModels(&dto.Product{})
}
