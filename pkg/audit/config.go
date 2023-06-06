package audit

import (
	"encoding/json"
	"example/web-service-gin/pkg/database"
	"example/web-service-gin/pkg/logger"
	"os"
)

type AuditFieldConfiguration struct {
	Name string `json:"Name"`
	Code string `json:"Code"`
}

type AuditConfiguration struct {
	Enable         bool                                          `json:"Enable"`
	Configurations map[string]map[string]AuditFieldConfiguration `json:"Configurations"`
}

var AuditConfig AuditConfiguration

func RegisterAudit() {
	raw, err := os.ReadFile("config/audit.json")
	if err != nil {
		logger.Logger.Info("Error occured while reading config/audit.json")
		return
	}
	json.Unmarshal(raw, &AuditConfig)
	// if enable audit trail
	if AuditConfig.Enable {
		RegisterAuditCallbacks(database.DB)
	}
}

// Check a table is register audit in config
func IsAuditable(tableName string) bool {
	if _, ok := AuditConfig.Configurations[tableName]; ok {
		return true
	}
	return false
}

// Get all field that need to audit trail of table from config
func GetFieldTrails(tableName string) []string {
	var fields []string
	for k := range AuditConfig.Configurations[tableName] {
		fields = append(fields, k)
	}
	return fields
}
