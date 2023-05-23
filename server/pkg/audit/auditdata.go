package audit

type ChangeData struct {
	From string
	To   string
}
type AuditData struct {
	CreatedBy string
	Data      map[string]ChangeData
}
