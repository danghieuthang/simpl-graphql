package repository

import (
	"context"
	"example/web-service-gin/pkg/logger"
	"fmt"

	"gorm.io/gorm"
)

var (
	ErrNotFound = gorm.ErrRecordNotFound
)

type IRepository interface {
	GetAll(target interface{}, preloads ...string) error
	GetBatch(target interface{}, limit, offset int, preloads ...string) error

	GetWhere(target interface{}, condition string, preloads ...string) error
	GetWhereBatch(target interface{}, condition string, limit, offset int, preloads ...string) error

	CountWhere(target interface{}, total *int64, condition string) error

	GetByField(target interface{}, field string, value interface{}, preloads ...string) error
	GetByFields(target interface{}, filters map[string]interface{}, preloads ...string) error
	GetByFieldBatch(target interface{}, field string, value interface{}, limit, offset int, preloads ...string) error
	GetByFieldsBatch(target interface{}, filters map[string]interface{}, limit, offset int, preloads ...string) error

	GetOneByField(target interface{}, field string, value interface{}, preloads ...string) error
	GetOneByFields(target interface{}, filters map[string]interface{}, preloads ...string) error

	GetOneByID(target interface{}, id string, preloads ...string) error

	Create(target interface{}) error
	CreateWithContext(target interface{}, ctx context.Context) error
	Save(target interface{}) error
	SaveWithContext(target interface{}, ctx context.Context) error
	Delete(target interface{}) error

	DBWithPreloads(preloads []string) *gorm.DB
}

type gormRepository struct {
	logger       logger.ILogger
	db           *gorm.DB
	defaultJoins []string
}

// NewGormRepository returns a new base repository that implements TransactionRepository
func NewGormRepository(db *gorm.DB, logger logger.ILogger, defaultJoins ...string) IRepository {
	return &gormRepository{
		defaultJoins: defaultJoins,
		logger:       logger,
		db:           db,
	}
}

func (r *gormRepository) GetAll(target interface{}, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Unscoped().
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) GetBatch(target interface{}, limit, offset int, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Unscoped().
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) GetWhere(target interface{}, condition string, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(condition).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) GetWhereBatch(target interface{}, condition string, limit, offset int, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(condition).
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) CountWhere(target interface{}, total *int64, condition string) error {
	res := r.db.Model(target).
		Where(condition).
		Count(total)

	return r.HandleError(res)
}

func (r *gormRepository) GetByField(target interface{}, field string, value interface{}, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) GetByFields(target interface{}, filters map[string]interface{}, preloads ...string) error {
	db := r.DBWithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) GetByFieldBatch(target interface{}, field string, value interface{}, limit, offset int, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) GetByFieldsBatch(target interface{}, filters map[string]interface{}, limit, offset int, preloads ...string) error {
	db := r.DBWithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) GetOneByField(target interface{}, field string, value interface{}, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		First(target)

	return r.HandleOneError(res)
}

func (r *gormRepository) GetOneByFields(target interface{}, filters map[string]interface{}, preloads ...string) error {
	db := r.DBWithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.First(target)
	return r.HandleOneError(res)
}

func (r *gormRepository) GetOneByID(target interface{}, id string, preloads ...string) error {
	res := r.DBWithPreloads(preloads).
		Where("id = ?", id).
		First(target)

	return r.HandleOneError(res)
}
func (r *gormRepository) Create(target interface{}) error {
	res := r.db.Create(target)
	return r.HandleError(res)
}

func (r *gormRepository) CreateWithContext(target interface{}, ctx context.Context) error {
	db := r.getTxFromContext(ctx)
	res := db.WithContext(ctx).Create(target)
	return r.HandleError(res)
}

func (r *gormRepository) Save(target interface{}) error {
	res := r.db.Save(target)
	return r.HandleError(res)
}
func (r *gormRepository) SaveWithContext(target interface{}, ctx context.Context) error {
	tx := r.getTxFromContext(ctx)
	res := tx.WithContext(ctx).Save(target)
	return r.HandleError(res)
}

func (r *gormRepository) Delete(target interface{}) error {
	res := r.db.Delete(target)
	return r.HandleError(res)
}

func (r *gormRepository) DeleteWithContext(ctx context.Context, target interface{}) error {
	tx := r.getTxFromContext(ctx)
	res := tx.Delete(target)
	return r.HandleError(res)
}
func (r *gormRepository) HandleError(res *gorm.DB) error {
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		err := fmt.Errorf("DB: %w", res.Error)
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *gormRepository) HandleOneError(res *gorm.DB) error {
	if err := r.HandleError(res); err != nil {
		return err
	}

	if res.RowsAffected != 1 {
		return ErrNotFound
	}

	return nil
}

func (r *gormRepository) DBWithPreloads(preloads []string) *gorm.DB {
	dbConn := r.db

	for _, join := range r.defaultJoins {
		dbConn = dbConn.Joins(join)
	}

	logger.Logger.Infof("Preload:%s", preloads)
	for _, preload := range preloads {
		dbConn = dbConn.Preload(preload)
	}

	return dbConn
}

// Get the executing transaction from the context.
// If not exist in context then return default.
func (r *gormRepository) getTxFromContext(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		return tx
	}
	return r.db
}
