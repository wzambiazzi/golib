package dao

import (
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
)

//CheckProduct func returns if a product exists or not based on tenant_id, ref1 and ref2
func CheckProduct(conn *sqlx.DB, tid int, ref1, ref2 string) (exists bool, err error) {
	const query = "SELECT EXISTS(SELECT 1 FROM public.resource_metadata WHERE tenant_id = $1 AND trim(ref1) = $2 AND trim(ref2) = $3) AS \"exists\";"

	if err = conn.Get(&exists, query, tid, ref1, ref2); err != nil {
		if err.Error() != "sql no results" {
			return false, db.WrapError(err, "conn.Get()")
		}
	}

	return
}
