package audit

import (
	"fmt"

	"gorm.io/gorm"
)

const (
	// CurrentUserDBScopeKey is the key for current user in db scope
	CurrentUserDBScopeKey = "current_user"

	createCallbackKey        = "audit:assign_created_updated_by"
	updateCallbackKey        = "audit:assign_updated_by"
	gormUpdateAttrs          = "gorm:update_attrs"
	gormBeforeCreate         = "gorm:before_create"
	gormBeforeUpdate         = "gorm:before_update"
	updatedByColumnName      = "updated_by"
	whoAuditFieldCount       = 2
	updatedByObjectFieldName = "UpdatedBy"
	createdByObjectFieldName = "CreatedBy"
)

// isAuditable check if the audit.model exists in the inputObject or not
func isAuditable(scope *gorm.DB) (isAuditable bool) {
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
	user, hasUser := scope.Get(CurrentUserDBScopeKey)
	if hasUser {
		return fmt.Sprintf("%v", user), true
	}
	return "", false
}

// assignUpdatedBy sets the value for updated by column
func assignUpdatedBy(scope *gorm.DB) {
	if isAuditable(scope) {
		if user, ok := GetCurrentUser(scope); ok {
			if attrs, ok := scope.InstanceGet(gormUpdateAttrs); ok {
				updateAttrs := attrs.(map[string]interface{})
				updateAttrs[updatedByColumnName] = user
				scope.InstanceSet(gormUpdateAttrs, updateAttrs)
			} else {
				scope.Statement.SetColumn(updatedByObjectFieldName, user)
			}
		}
	}
}

// assignCreatedAndUpdatedBy sets the value for both updated by and created by columns
func assignCreatedAndUpdatedBy(scope *gorm.DB) {
	if isAuditable(scope) {
		if user, ok := GetCurrentUser(scope); ok {
			scope.Statement.SetColumn(createdByObjectFieldName, user)
		}
		assignUpdatedBy(scope)
	}
}

// RegisterAuditCallbacks register callback into GORM DB
func RegisterAuditCallbacks(db *gorm.DB) {
	callback := db.Callback()
	if callback.Create().Get(createCallbackKey) == nil {
		callback.Create().Before(gormBeforeCreate).Register(createCallbackKey, assignCreatedAndUpdatedBy)
	}
	if callback.Update().Get(updateCallbackKey) == nil {
		callback.Update().Before(gormBeforeUpdate).Register(updateCallbackKey, assignUpdatedBy)
	}
}
