package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// ServiceItemsCustomerContacts represents a relationship between a MTOServiceItem and a MTOServiceItemCustomerContact.
type ServiceItemsCustomerContacts struct {
	ID uuid.UUID `db:"id"`
	// For many-to-many relationships, Pop will query for an id field on the join table based off of the names of the originating models.
	// since MTOServiceItem and MTOServiceItemCustomerContact have four capital letters in front,
	// Pop will not put an underscore between the "mto" and "service" in both ids. The columns on this table must exclude that underscore.
	// This is why the db and fk_id flags are different.
	MTOServiceItemID                uuid.UUID                     `db:"mtoservice_item_id" fk_id:"mto_service_item_id"`
	MTOServiceItem                  MTOServiceItem                `belongs_to:"mto_service_items" db:"-"`
	MTOServiceItemCustomerContactID uuid.UUID                     `db:"mtoservice_item_customer_contact_id" fk_id:"mto_service_item_customer_contact_id"`
	MTOServiceItemCustomerContact   MTOServiceItemCustomerContact `belongs_to:"mto_service_items_customer_contacts" db:"-"`
	CreatedAt                       time.Time                     `db:"created_at"`
	UpdateAt                        time.Time                     `db:"updated_at"`
	DeletedAt                       *time.Time                    `db:"deleted_at"`
}

// TableName overrides the table name used by Pop.
func (s ServiceItemsCustomerContacts) TableName() string {
	return "service_items_customer_contacts"
}
