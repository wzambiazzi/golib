package model

import (
	"time"

	"github.com/lib/pq"
)

//SFDataSync type
type SFDataSync struct {
	ID            int64       `db:"id"`
	TenantID      int         `db:"tenant_id"`
	ExecutionID   int64       `db:"execution_id"`
	StatusID      int16       `db:"status_id"`
	ObjectID      int64       `db:"sf_object_id"`
	ObjectName    string      `db:"sf_object_name"`
	DocOriginalID string      `db:"doc_original_id"`
	IsActive      bool        `db:"is_active"`
	CreatedAt     time.Time   `db:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at"`
	IsDeleted     bool        `db:"is_deleted"`
	DeletedAt     pq.NullTime `db:"deleted_at"`
}
