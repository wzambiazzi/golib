package dao

import (
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

//GetResourceMetadataToProcess func
func GetResourceMetadataToProcess(conn *sqlx.DB, tenantID int) (d []model.ResourceMetadata, err error) {
	query := `SELECT tenant_id, sf_content_document_id, sf_content_version_id 
	            FROM public.fn_get_resource_to_process($1);`

	err = conn.Select(&d, query, tenantID)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return
}

//GetProductsWithoutResources func
func GetProductsWithoutResources(conn *sqlx.DB, tenantID int) (d []string, err error) {
	query := `
		WITH t AS (
			SELECT DISTINCT p.ref1, p.ref2
			FROM public.vw_produto_modelo_cor p 
			LEFT JOIN public.resource_metadata r ON p.tenant_id = r.tenant_id 
					AND r.is_deleted = FALSE 
					AND LPAD(p.ref1::text, 5, '0'::text) = LPAD(r.ref1::text, 5, '0'::text) 
					AND LPAD(p.ref2::text, 5, '0'::text) = LPAD(r.ref2::text, 5, '0'::text) 
					AND COALESCE(r.sf_content_document_id, '') = '' 
					AND COALESCE(r.sf_content_version_id, '') = ''
			WHERE p.tenant_id = $1
			AND p.ref1 IS NOT NULL
            AND p.ref2 IS NOT NULL
			AND r.id IS NULL
			ORDER BY p.ref1, p.ref2
		) 
		SELECT DISTINCT LPAD(t.ref1::text, 5, '0'::text) || LPAD(t.ref2::text, 5, '0'::text) AS product_color
		FROM t;`

	err = conn.Select(&d, query, tenantID)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return
}

//SaveResourceMetadata func
func SaveResourceMetadata(conn *sqlx.DB, data *model.ResourceMetadata) (err error) {
	query := `
		INSERT INTO public.resource_metadata AS rm (tenant_id, original_file_name, original_file_extension, content_type, "size", preview_content_b64, full_content_b64, sf_content_document_id, sf_content_version_id, is_downloaded) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, true)
			ON CONFLICT (tenant_id, sf_content_document_id)
			DO UPDATE
		   SET sf_content_version_id   = EXCLUDED.sf_content_version_id,
			   original_file_name      = EXCLUDED.original_file_name, 
			   original_file_extension = EXCLUDED.original_file_extension, 
			   content_type            = EXCLUDED.content_type, 
			   "size"                  = EXCLUDED."size",
			   preview_content_b64     = EXCLUDED.preview_content_b64,
			   full_content_b64        = EXCLUDED.full_content_b64,
			   updated_at              = now()
			WHERE rm.tenant_id = $1;`

	_, err = conn.Exec(query, data.TenantID, data.OriginalFileName, data.OriginalFileExtension, data.ContentType, data.Size, data.PreviewBontentB64, data.FullContentB64, data.SfContentDocumentID, data.SfContentVersionID)
	if err != nil {
		err = db.WrapError(err, "conn.Exec()")
		return
	}

	return
}

//SaveResourceMetadataWithRefs func
func SaveResourceMetadataWithRefs(conn *sqlx.DB, data *model.ResourceMetadata) (rows int64, err error) {
	query := `
		INSERT INTO public.resource_metadata AS rm (tenant_id, original_file_name, original_file_extension, content_type, "size", ref1, ref2, sequence, size_type, full_content_b64) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (tenant_id, ref1, ref2, sequence, size_type)
			DO UPDATE
		   SET original_file_name      = EXCLUDED.original_file_name, 
			   original_file_extension = EXCLUDED.original_file_extension, 
			   content_type            = EXCLUDED.content_type, 
			   "size"                  = EXCLUDED."size",
			   full_content_b64        = EXCLUDED.full_content_b64,
			   is_deleted              = FALSE,
			   deleted_at              = NULL,
			   updated_at              = NOW()
			WHERE rm.tenant_id = $1;`

	result, err := conn.Exec(query, data.TenantID, data.OriginalFileName, data.OriginalFileExtension, data.ContentType, data.Size, data.Ref1, data.Ref2, data.Sequence, data.SizeType, data.FullContentB64)
	if err != nil {
		err = db.WrapError(err, "conn.Exec()")
		return
	}

	rows, err = result.RowsAffected()
	if err != nil {
		err = db.WrapError(err, "result.RowsAffected()")
		return
	}

	return
}

//SoftDeleteImages func
func SoftDeleteImages(conn *sqlx.DB, tenantID int) error {
	query := `
		UPDATE public.resource_metadata 
		   SET is_deleted = TRUE, 
			   deleted_at = NOW()
		  FROM public.resource_metadata r 
		  LEFT JOIN public.vw_produto_modelo_cor p ON r.tenant_id = p.tenant_id 
				AND LPAD(r.ref1::text, 5, '0'::text) = LPAD(p.ref1::text, 5, '0'::text) 
				AND LPAD(r.ref2::text, 5, '0'::text) = LPAD(p.ref2::text, 5, '0'::text)
		 WHERE r.tenant_id = $1
		   AND COALESCE(r.ref1, '') != ''
		   AND COALESCE(r.ref2, '') != ''
		   AND COALESCE(r.sf_content_document_id, '') = ''
		   AND COALESCE(r.sf_content_version_id, '') = ''
		   AND r.is_deleted = FALSE
		   AND p.tenant_id IS NULL;`

	if _, err := conn.Exec(query, tenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
