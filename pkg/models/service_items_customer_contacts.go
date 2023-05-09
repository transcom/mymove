package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// ServiceItemsCustomerContacts represents an MTOServiceItem and a MTOServiceItemCustomerContact
type ServiceItemsCustomerContacts struct {
	ID                              uuid.UUID                     `db:"id"`
	MTOServiceItemID                uuid.UUID                     `db:"mto_service_item_id"`
	MTOServiceItem                  MTOServiceItem                `belongs_to:"mto_service_items" db:"-"`
	MTOServiceItemCustomerContactID uuid.UUID                     `db:"mto_service_item_customer_contact_id"`
	MTOServiceItemCustomerContact   MTOServiceItemCustomerContact `belongs_to:"mto_service_items_customer_contacts" db:"-"`
	CreatedAt                       time.Time                     `db:"created_at"`
	UpdateAt                        time.Time                     `db:"updated_at"`
	DeletedAt                       *time.Time                    `db:"deleted_at"`
}

// TableName overrides the table name used by Pop.
func (s ServiceItemsCustomerContacts) TableName() string {
	return "service_items_customer_contacts"
}
