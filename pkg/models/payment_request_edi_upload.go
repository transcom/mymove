package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type PaymentRequestEdiUpload struct {
	ID        uuid.UUID  `db:"id"`
	UploadID  uuid.UUID  `db:"upload_id"`
	Upload    Upload     `belongs_to:"uploads" fk_id:"upload_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
