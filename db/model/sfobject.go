package model

import (
	"database/sql"

	m "bitbucket.org/everymind/evmd-golib/modelbase"
)

//SFObject type
type SFObject struct {
	ID          int     `db:"id"`
	TenantID    int     `db:"tenant_id"`
	ExecutionID int64   `db:"execution_id"`
	Name        string  `db:"sf_object_name"`
	DocDescribe m.JSONB `db:"doc_describe"`
	DocMetaData m.JSONB `db:"doc_meta_data"` // In DB is JSONB
}

//SFObjectToProcess type
type SFObjectToProcess struct {
	ID         int64          `db:"id"`
	ObjectName sql.NullString `db:"sf_object_name"`
	TenantID   int            `db:"tenant_id"`
	TenantName string         `db:"tenant_name"`
	Filter     sql.NullString `db:"filter"`
	SfaPks     m.JSONB        `db:"sfa_pks"`
}

//SFObjectToProcesses type
type SFObjectToProcesses []SFObjectToProcess
