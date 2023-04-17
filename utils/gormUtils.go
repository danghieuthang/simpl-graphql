package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// Build preload association with given select
func PreLoadWithSelect(sourceTx *gorm.DB, fields []string, root string) *gorm.DB {
	groupFields := GroupField(fields, root)
	for key, value := range groupFields {
		if key == root {
			continue
		}
		sourceTx = sourceTx.Preload(key, func(db *gorm.DB) *gorm.DB {
			return db.Select(value)
		})
	}
	return sourceTx.Select(groupFields[root])
}

func GroupField(fields []string, root string) map[string][]string {
	group := make(map[string][]string)
	for _, field := range fields {
		if strings.Contains(field, ".") {
			s := strings.Split(field, ".")
			caser := cases.Title(language.AmericanEnglish)
			entityName := caser.String(s[0])
			if fieldInGroup, ok := group[entityName]; ok {
				group[entityName] = append(fieldInGroup, s[1])
				continue
			}
			group[entityName] = []string{s[1]}
			continue
		}
		if fieldInGroup, ok := group[root]; ok {
			group[root] = append(fieldInGroup, field)
			continue
		}
		group[root] = []string{field}
	}
	return group
}
