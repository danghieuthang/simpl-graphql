package audit

import (
	"example/web-service-gin/pkg/logger"
	"example/web-service-gin/pkg/middleware/auth"
	"example/web-service-gin/pkg/utils"
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	// CurrentUserDBScopeKey is the key for current user in db scope
	CurrentUserDBScopeKey = "current_user"

	createCallbackKey        = "audit:assign_created_updated_by"
	updateCallbackKey        = "audit:assign_updated_by"
	gormUpdateAttrs          = "gorm:update_attrs"
	gormBeforeCreate         = "gorm:before_create"
	gormBeforeUpdate         = "gorm:before_update"
	gormAfterCreate          = "gorm:after_create"
	gormAfterUpdate          = "gorm:after_update"
	updatedByColumnName      = "updated_by"
	whoAuditFieldCount       = 2
	updatedByObjectFieldName = "UpdatedBy"
	createdByObjectFieldName = "CreatedBy"
	createdAtFieldName       = "CreatedAt"
	lastModifiedFieldName    = "LastModifiedAt"

	// audit trail
	auditCreatedCallbackKey = "audit:created"
	auditUpdatedCallbackKey = "audit:modified"
)

// isAuditableDefault check if the audit.model exists in the inputObject or not
func isAuditableDefault(scope *gorm.DB) (isAuditable bool) {
	if scope.Statement.Schema.ModelType == nil {
		return false
	}
	auditFieldCount := 0
	fields := scope.Statement.Schema.Fields
	for _, field := range fields {
		if field.Name == updatedByObjectFieldName || field.Name == createdByObjectFieldName {
			auditFieldCount++
		}
	}
	return auditFieldCount == whoAuditFieldCount
}

// GetCurrentUser gets the current user from db scope
func GetCurrentUser(scope *gorm.DB) (string, bool) {
	user, hasUser := scope.Statement.Context.Value("currentUser").(*auth.AuthenticatedUser)
	if hasUser {
		return fmt.Sprintf("%v", user.Email), true
	}
	return "", false
}

// assignUpdatedBy sets the value for updated by column
func assignUpdatedBy(scope *gorm.DB) {
	if isAuditableDefault(scope) {
		if user, ok := GetCurrentUser(scope); ok {
			if attrs, ok := scope.InstanceGet(gormUpdateAttrs); ok {
				updateAttrs := attrs.(map[string]interface{})
				updateAttrs[updatedByColumnName] = user
				scope.InstanceSet(gormUpdateAttrs, updateAttrs)
			} else {
				scope.Statement.SetColumn(updatedByObjectFieldName, user)
				scope.Statement.SetColumn(lastModifiedFieldName, time.Now())
			}
		}
	}
}

func assignCreated(scope *gorm.DB) {
	if isAuditableDefault(scope) {
		if user, ok := GetCurrentUser(scope); ok {
			scope.Statement.SetColumn(createdByObjectFieldName, user)
			scope.Statement.SetColumn(createdAtFieldName, time.Now())
		}
	}
}

// isAuditableDefault check if the audit.model exists in the inputObject or not
func isAuditTrail(scope *gorm.DB) (isAuditable bool) {
	if scope.Statement.Schema.ModelType == nil {
		return false
	}
	if !AuditConfig.Enable {
		return false
	}
	// if !scope.Statement.Changed() {
	// 	return false
	// }
	tableName := scope.Statement.Schema.Name
	return IsAuditable(tableName)
}

// Check a schema field is exist in list audit field
func isSchemaFieldAudit(auditFields []string, field *schema.Field) bool {
	for _, v := range auditFields {
		if field.Name == v {
			return true
		}
	}
	return false
}

// assignCreatedAndUpdatedBy sets the value for both updated by and created by columns
func assignCreatedAndUpdatedBy(scope *gorm.DB, action string) {
	if isAuditTrail(scope) {
		primaryKey, _ := scope.Statement.Schema.PrimaryFields[0].ValueOf(scope.Statement.Context, scope.Statement.ReflectValue)
		key := fmt.Sprintf("%s.%v", scope.Statement.Schema.Name, primaryKey)
		var auditData = AuditData{
			Data:   make(map[string]ChangeData),
			Key:    key,
			Action: action,
		}
		if user, ok := GetCurrentUser(scope); ok {
			auditData.CreatedBy = user
		}
		auditFields := GetFieldTrails(scope.Statement.Schema.Name)

		ctx := scope.Statement.Context
		changeTracker, hasTracking := ctx.Value("changeTracker").(IChangeTracker)
		if hasTracking {
			if oldValue, ok := changeTracker.Get(key); ok {
				for _, field := range scope.Statement.Schema.Fields {
					if !isSchemaFieldAudit(auditFields, field) {
						continue
					}

					fromValue := fmt.Sprintf("%v", utils.GetFieldValue(oldValue, field.Name))
					toValue, isZero := field.ValueOf(scope.Statement.Context, scope.Statement.ReflectValue)
					if !isZero && !reflect.DeepEqual(fromValue, toValue) {
						var changeData = &ChangeData{
							From: fmt.Sprintf("%v", fromValue),
							To:   fmt.Sprintf("%v", toValue),
						}
						auditData.Data[field.Name] = *changeData
					}
				}
			} else {
				for _, field := range scope.Statement.Schema.Fields {
					if !isSchemaFieldAudit(auditFields, field) {
						continue
					}
					var changeData = &ChangeData{
						From: "",
					}
					fieldValue, isZero := field.ValueOf(scope.Statement.Context, scope.Statement.ReflectValue)
					if !isZero {
						changeData.To = fmt.Sprintf("%v", fieldValue)
					} else {
						changeData.From = ""
					}
					auditData.Data[field.Name] = *changeData
				}
			}
		}

		logger.Logger.Info(auditData)
	}
}

func auditCreate(scope *gorm.DB) {
	assignCreatedAndUpdatedBy(scope, AuditActionCreate)
}
func auditUpdate(scope *gorm.DB) {
	assignCreatedAndUpdatedBy(scope, AuditActionUpdate)
}

// RegisterAuditCallbacks register callback into GORM DB
func RegisterAuditCallbacks(tx *gorm.DB) {
	callback := tx.Callback()
	if callback.Create().Get(createCallbackKey) == nil {
		callback.Create().Before(gormBeforeCreate).Register(createCallbackKey, assignCreated)
	}
	if callback.Update().Get(updateCallbackKey) == nil {
		callback.Update().Before(gormBeforeUpdate).Register(updateCallbackKey, assignUpdatedBy)
	}

	if callback.Create().Get(auditCreatedCallbackKey) == nil {
		callback.Create().Before(gormAfterCreate).Register(auditCreatedCallbackKey, auditCreate)
	}
	if callback.Update().Get(auditUpdatedCallbackKey) == nil {
		callback.Update().Before(gormAfterUpdate).Register(auditUpdatedCallbackKey, auditUpdate)
	}
}
