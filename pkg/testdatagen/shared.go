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
)

// Assertions defines assertions about what the data contains
type Assertions struct {
	User                   models.User
	OfficeUser             models.OfficeUser
	ServiceMember          models.ServiceMember
	Order                  models.Order
	Move                   models.Move
	PersonallyProcuredMove models.PersonallyProcuredMove
	Document               models.Document
	BackupContact          models.BackupContact
	Upload                 models.Upload
	Address                models.Address
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
		mergo.Merge(dst, src, mergo.WithTransformers(customTransformer{})),
	)
}

// customTransformer handles testing for zero values in structs that mergo can't normally deal with
type customTransformer struct {
}

// Checks if dst is a zero value, then overwrites with src
func (t customTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	// UUID comparison
	if typ == reflect.TypeOf(uuid.UUID{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				dstUUID := dst.Interface().(uuid.UUID)
				if isZeroUUID(dstUUID) {
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
				isZero := dst.MethodByName("IsZero")
				result := isZero.Call([]reflect.Value{})
				if result[0].Bool() {
					dst.Set(src)
				}
			}
			return nil
		}
	}
	return nil
}
