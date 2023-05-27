package repository

import (
	"context"
	"example/web-service-gin/pkg/audit"
	"example/web-service-gin/pkg/entity"
	"example/web-service-gin/pkg/logger"
	"example/web-service-gin/pkg/utils"
	"fmt"

	"gorm.io/gorm"
)

var (
	ErrNotFound = gorm.ErrRecordNotFound
)

type IBaseRepository[T entity.IEntity] interface {
	GetOne(ctx context.Context, condition string, preloads ...string) (T, error)
	GetOneAsNoTracking(ctx context.Context, condition string, preloads ...string) (T, error)
	GetAll(ctx context.Context, condition string, limit, offset int, preloads ...string) ([]T, error)
	GetAllAsNoTracking(ctx context.Context, condition string, limit, offset int, preloads ...string) ([]T, error)
	Count(ctx context.Context, condition string) (int64, error)
	Update(ctx context.Context, t T) error
	Create(ctx context.Context, t T) (T, error)
	Delete(ctx context.Context, t T) error
}

type BaseRepository[T entity.IEntity] struct {
	db     *gorm.DB
	logger logger.ILogger
}

func NewRepository[T entity.IEntity](db *gorm.DB, logger logger.ILogger) IBaseRepository[T] {
	return &BaseRepository[T]{
		db:     db,
		logger: logger,
	}
}
func (r BaseRepository[T]) GetOne(ctx context.Context, condition string, preloads ...string) (T, error) {
	var t T
	res := r.DBWithPreloads(ctx, preloads).
		Where(condition).
		Find(&t)
	err := r.HandleError(res)
	if err == nil {
		r.trackingData(ctx, t)
		return t, nil
	}
	return t, err
}
func (r BaseRepository[T]) GetOneAsNoTracking(ctx context.Context, condition string, preloads ...string) (T, error) {
	var t T
	res := r.DBWithPreloads(ctx, preloads).
		Where(condition).
		Find(&t)
	err := r.HandleError(res)
	if err == nil {
		return t, nil
	}
	return t, err
}
func (r BaseRepository[T]) GetAll(ctx context.Context, condition string, limit, offset int, preloads ...string) ([]T, error) {
	var t []T
	res := r.DBWithPreloads(ctx, preloads).
		Where(condition).
		Limit(limit).
		Offset(offset).
		Find(&t)
	err := r.HandleError(res)
	if err == nil {
		return t, nil
	}
	return t, err
}

func (r BaseRepository[T]) GetAllAsNoTracking(ctx context.Context, condition string, limit, offset int, preloads ...string) ([]T, error) {
	var t []T
	res := r.DBWithPreloads(ctx, preloads).
		Where(condition).
		Limit(limit).
		Offset(offset).
		Find(&t)
	err := r.HandleError(res)
	if err == nil {
		return t, nil
	}
	return t, err
}

func (r BaseRepository[T]) Count(ctx context.Context, condition string) (int64, error) {
	var t T
	var total int64
	res := r.db.Model(&t).Where(condition).
		Count(&total)
	err := r.HandleError(res)
	if err == nil {
		return total, nil
	}
	return total, err
}

func (r BaseRepository[T]) Create(ctx context.Context, t T) (T, error) {
	db := r.getTxFromContext(ctx)
	res := db.WithContext(ctx).Create(&t)
	return t, r.HandleError(res)
}

func (r BaseRepository[T]) Update(ctx context.Context, t T) error {
	db := r.getTxFromContext(ctx)
	res := db.WithContext(ctx).Save(&t)
	return r.HandleError(res)
}

func (r BaseRepository[T]) Delete(ctx context.Context, t T) error {
	db := r.getTxFromContext(ctx)
	res := db.WithContext(ctx).Delete(&t)
	return r.HandleError(res)
}

func (r BaseRepository[T]) HandleError(res *gorm.DB) error {
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		err := fmt.Errorf("DB: %w", res.Error)
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r BaseRepository[T]) DBWithPreloads(ctx context.Context, preloads []string) *gorm.DB {
	dbConn := r.getTxFromContext(ctx)

	logger.Logger.Infof("Preload:%s", preloads)
	for _, preload := range preloads {
		dbConn = dbConn.Preload(preload)
	}

	return dbConn
}

// Get the executing transaction from the context.
// If not exist in context then return default.
func (r BaseRepository[T]) getTxFromContext(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		return tx
	}
	return r.db
}

func (r BaseRepository[T]) trackingData(ctx context.Context, data T) {
	changeTracker, ok := ctx.Value("changeTracker").(audit.IChangeTracker)
	if ok {
		key := fmt.Sprintf("%s.%d", utils.GetType(data), data.GetId())
		changeTracker.Set(key, data)
	}
}
