package testdatagen

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/imdario/mergo"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// Assertions defines assertions about what the data contains
type Assertions struct {
	Address                models.Address
	BackupContact          models.BackupContact
	BlackoutDate           models.BlackoutDate
	Document               models.Document
	DutyStation            models.DutyStation
	Move                   models.Move
	MoveDocument           models.MoveDocument
	OfficeUser             models.OfficeUser
	Order                  models.Order
	PersonallyProcuredMove models.PersonallyProcuredMove
	ServiceMember          models.ServiceMember
	Shipment               models.Shipment
	ShipmentOffer          models.ShipmentOffer
	TspUser                models.TspUser
	Upload                 models.Upload
	User                   models.User
}

func stringPointer(s string) *string {
	return &s
}

func poundPointer(p unit.Pound) *unit.Pound {
	return &p
}

func uuidPointer(u uuid.UUID) *uuid.UUID {
	return &u
}

func timePointer(t time.Time) *time.Time {
	return &t
}

func mustCreate(db *pop.Connection, model interface{}) {
	verrs, err := db.ValidateAndCreate(model)
	if err != nil {
		log.Panic(fmt.Errorf("Errors encountered saving %v: %v", model, err))
	}
	if verrs.HasAny() {
		log.Panic(fmt.Errorf("Validation errors encountered saving %v: %v", model, verrs))
	}
}

func mustSave(db *pop.Connection, model interface{}) {
	verrs, err := db.ValidateAndSave(model)
	if err != nil {
		log.Panic(fmt.Errorf("Errors encountered saving %v: %v", model, err))
	}
	if verrs.HasAny() {
		log.Panic(fmt.Errorf("Validation errors encountered saving %v: %v", model, verrs))
	}
}

func noErr(err error) {
	if err != nil {
		log.Panic(fmt.Errorf("Error encountered: %v", err))
	}
}

// isZeroUUID determines whether a UUID is its zero value
func isZeroUUID(testID uuid.UUID) bool {
	return uuid.Equal(testID, uuid.UUID{})
}

// mergeModels merges src into dst, if non-zero values are present
// dst should be a pointer the struct you are merging into
func mergeModels(dst, src interface{}) {
	noErr(
		mergo.Merge(dst, src, mergo.WithOverride, mergo.WithTransformers(customTransformer{})),
	)
}

// customTransformer handles testing for zero values in structs that mergo can't normally deal with
type customTransformer struct {
}

// Checks if src is not a zero value, then overwrites dst
func (t customTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	// UUID comparison
	if typ == reflect.TypeOf(uuid.UUID{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				srcID := src.Interface().(uuid.UUID)
				if !isZeroUUID(srcID) {
					dst.Set(src)
				}
			}
			return nil
		}
	}
	// time.Time comparison
	if typ == reflect.TypeOf(time.Time{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				isZero := src.MethodByName("IsZero")
				result := isZero.Call([]reflect.Value{})
				if !result[0].Bool() {
					dst.Set(src)
				}
			}
			return nil
		}
	}
	return nil
}
