package repository

import (
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

	GetByField(target interface{}, field string, value interface{}, preloads ...string) error
	GetByFields(target interface{}, filters map[string]interface{}, preloads ...string) error
	GetByFieldBatch(target interface{}, field string, value interface{}, limit, offset int, preloads ...string) error
	GetByFieldsBatch(target interface{}, filters map[string]interface{}, limit, offset int, preloads ...string) error

	GetOneByField(target interface{}, field string, value interface{}, preloads ...string) error
	GetOneByFields(target interface{}, filters map[string]interface{}, preloads ...string) error

	GetOneByID(target interface{}, id string, preloads ...string) error

	Create(target interface{}) error
	Save(target interface{}) error
	Delete(target interface{}) error

	DBWithPreloads(preloads []string) *gorm.DB
}

type GormRepository struct {
	logger       logger.ILogger
	db           *gorm.DB
	defaultJoins []string
}

// NewGormRepository returns a new base repository that implements TransactionRepository
func NewGormRepository(db *gorm.DB, logger logger.ILogger, defaultJoins ...string) IRepository {
	return &GormRepository{
		defaultJoins: defaultJoins,
		logger:       logger,
		db:           db,
	}
}

func (r *GormRepository) GetAll(target interface{}, preloads ...string) error {
	r.logger.Debugf("Executing GetAll on %T", target)

	res := r.DBWithPreloads(preloads).
		Unscoped().
		Find(target)

	return r.HandleError(res)
}

func (r *GormRepository) GetBatch(target interface{}, limit, offset int, preloads ...string) error {
	r.logger.Debugf("Executing GetBatch on %T", target)

	res := r.DBWithPreloads(preloads).
		Unscoped().
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *GormRepository) GetWhere(target interface{}, condition string, preloads ...string) error {
	r.logger.Debugf("Executing GetWhere on %T with %v ", target, condition)

	res := r.DBWithPreloads(preloads).
		Where(condition).
		Find(target)

	return r.HandleError(res)
}

func (r *GormRepository) GetWhereBatch(target interface{}, condition string, limit, offset int, preloads ...string) error {
	r.logger.Debugf("Executing GetWhere on %T with %v ", target, condition)

	res := r.DBWithPreloads(preloads).
		Where(condition).
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *GormRepository) GetByField(target interface{}, field string, value interface{}, preloads ...string) error {
	r.logger.Debugf("Executing GetByField on %T with %v = %v", target, field, value)

	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Find(target)

	return r.HandleError(res)
}

func (r *GormRepository) GetByFields(target interface{}, filters map[string]interface{}, preloads ...string) error {
	r.logger.Debugf("Executing GetByField on %T with filters = %+v", target, filters)

	db := r.DBWithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.Find(target)

	return r.HandleError(res)
}

func (r *GormRepository) GetByFieldBatch(target interface{}, field string, value interface{}, limit, offset int, preloads ...string) error {
	r.logger.Debugf("Executing GetByField on %T with %v = %v", target, field, value)

	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *GormRepository) GetByFieldsBatch(target interface{}, filters map[string]interface{}, limit, offset int, preloads ...string) error {
	r.logger.Debugf("Executing GetByField on %T with filters = %+v", target, filters)

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

func (r *GormRepository) GetOneByField(target interface{}, field string, value interface{}, preloads ...string) error {
	r.logger.Debugf("Executing GetOneByField on %T with %v = %v", target, field, value)

	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		First(target)

	return r.HandleOneError(res)
}

func (r *GormRepository) GetOneByFields(target interface{}, filters map[string]interface{}, preloads ...string) error {
	r.logger.Debugf("Executing FindOneByField on %T with filters = %+v", target, filters)

	db := r.DBWithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.First(target)
	return r.HandleOneError(res)
}

func (r *GormRepository) GetOneByID(target interface{}, id string, preloads ...string) error {
	r.logger.Debugf("Executing GetOneByID on %T with ID %v", target, id)

	res := r.DBWithPreloads(preloads).
		Where("id = ?", id).
		First(target)

	return r.HandleOneError(res)
}

func (r *GormRepository) Create(target interface{}) error {
	r.logger.Debugf("Executing Create on %T", target)

	res := r.db.Create(target)
	return r.HandleError(res)
}

func (r *GormRepository) CreateTx(target interface{}, tx *gorm.DB) error {
	r.logger.Debugf("Executing Create on %T", target)

	res := tx.Create(target)
	return r.HandleError(res)
}

func (r *GormRepository) Save(target interface{}) error {
	r.logger.Debugf("Executing Save on %T", target)

	res := r.db.Save(target)
	return r.HandleError(res)
}

func (r *GormRepository) SaveTx(target interface{}, tx *gorm.DB) error {
	r.logger.Debugf("Executing Save on %T", target)

	res := tx.Save(target)
	return r.HandleError(res)
}

func (r *GormRepository) Delete(target interface{}) error {
	r.logger.Debugf("Executing Delete on %T", target)

	res := r.db.Delete(target)
	return r.HandleError(res)
}

func (r *GormRepository) DeleteTx(target interface{}, tx *gorm.DB) error {
	r.logger.Debugf("Executing Delete on %T", target)

	res := tx.Delete(target)
	return r.HandleError(res)
}

func (r *GormRepository) HandleError(res *gorm.DB) error {
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		err := fmt.Errorf("Error: %w", res.Error)
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *GormRepository) HandleOneError(res *gorm.DB) error {
	if err := r.HandleError(res); err != nil {
		return err
	}

	if res.RowsAffected != 1 {
		return ErrNotFound
	}

	return nil
}

func (r *GormRepository) DBWithPreloads(preloads []string) *gorm.DB {
	dbConn := r.db

	for _, join := range r.defaultJoins {
		dbConn = dbConn.Joins(join)
	}

	for _, preload := range preloads {
		dbConn = dbConn.Preload(preload)
	}

	return dbConn
}
