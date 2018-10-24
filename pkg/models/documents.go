package models

import (
	"github.com/gobuffalo/uuid"
	"time"
)

// A Document represents a physical artifact such as a multipage form that was
// filled out by hand. A Document can have many associated Uploads, which allows
// for handling multiple files that belong to the same document.
type Document struct {
	ID              uuid.UUID     `db:"id"`
	ServiceMemberID uuid.UUID     `db:"service_member_id"`
	ServiceMember   ServiceMember `belongs_to:"service_members"`
	CreatedAt       time.Time     `db:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at"`
	Uploads         Uploads       `has_many:"uploads" order_by:"created_at asc"`
}

// Documents is not required by pop and may be deleted
type Documents []Document

// An Upload represents an uploaded file, such as an image or PDF.
type Upload struct {
	ID          uuid.UUID  `db:"id"`
	DocumentID  *uuid.UUID `db:"document_id"`
	Document    Document   `belongs_to:"documents"`
	UploaderID  uuid.UUID  `db:"uploader_id"`
	Filename    string     `db:"filename"`
	Bytes       int64      `db:"bytes"`
	ContentType string     `db:"content_type"`
	Checksum    string     `db:"checksum"`
	StorageKey  string     `db:"storage_key"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

// Uploads is not required by pop and may be deleted
type Uploads []Upload

// DocumentDB defines the functions needed from the DB for accessing Document objects
type DocumentDB interface {
	Fetch(id uuid.UUID) (*Document, error)
	FetchUpload(id uuid.UUID) (*Upload, error)
}
