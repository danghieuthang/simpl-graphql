package audit

type IAuditService interface {
	Save(data AuditData) error
	GetByKey(key string) (AuditData, error)
}

type AuditService struct {
}
