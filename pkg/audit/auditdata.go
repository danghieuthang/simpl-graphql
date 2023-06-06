package audit

const (
	AuditActionCreate = "Created"
	AuditActionUpdate = "Updated"
	AuditActionDelete = "Deleted"
)

type ChangeData struct {
	From string
	To   string
}
type AuditData struct {
	Key       string
	Action    string
	CreatedBy string
	Data      map[string]ChangeData
}
