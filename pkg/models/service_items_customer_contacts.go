package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// ServiceItemsCustomerContacts represents an MTOServiceItem and a MTOServiceItemCustomerContact
type ServiceItemsCustomerContacts struct {
	ID                uuid.UUID  `db:"id"`
	MTOServiceItemID  uuid.UUID  `db:"user_id"`
	CustomerContactID uuid.UUID  `db:"role_id"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdateAt          time.Time  `db:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at"`
}

// TableName overrides the table name used by Pop.
func (s ServiceItemsCustomerContacts) TableName() string {
	return "service_items_customer_contacts"
}
