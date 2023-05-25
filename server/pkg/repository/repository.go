package repository

import (
	"example/web-service-gin/pkg/entity"
	"example/web-service-gin/pkg/logger"

	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB, logger logger.ILogger) IBaseRepository[*entity.User] {
	return NewRepository[*entity.User](db, logger)
}

func NewRoleRepository(db *gorm.DB, logger logger.ILogger) IBaseRepository[*entity.Role] {
	return NewRepository[*entity.Role](db, logger)
}
