package payloads

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/handlers"

	"github.com/transcom/mymove/pkg/models"
)

func TestOrder(t *testing.T) {
	order := &models.Order{}
	Order(order)
}

func TestEntitlement(t *testing.T) {

	t.Run("Success - Returns the entitlement payload with only required fields", func(t *testing.T) {
		entitlement := models.Entitlement{
			ID:                             uuid.Must(uuid.NewV4()),
			DependentsAuthorized:           nil,
			TotalDependents:                nil,
			NonTemporaryStorage:            nil,
			PrivatelyOwnedVehicle:          nil,
			DBAuthorizedWeight:             nil,
			StorageInTransit:               nil,
			RequiredMedicalEquipmentWeight: 0,
			OrganizationalClothingAndIndividualEquipment: false,
			ProGearWeight:       0,
			ProGearWeightSpouse: 0,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		payload := Entitlement(&entitlement)

		assert.Equal(t, strfmt.UUID(entitlement.ID.String()), payload.ID)
		assert.Equal(t, int64(0), payload.RequiredMedicalEquipmentWeight)
		assert.Equal(t, false, payload.OrganizationalClothingAndIndividualEquipment)
		assert.Equal(t, int64(0), payload.ProGearWeight)
		assert.Equal(t, int64(0), payload.ProGearWeightSpouse)
		assert.NotEmpty(t, payload.ETag)
		assert.Equal(t, etag.GenerateEtag(entitlement.UpdatedAt), payload.ETag)

		assert.Nil(t, payload.AuthorizedWeight)
		assert.Nil(t, payload.DependentsAuthorized)
		assert.Nil(t, payload.NonTemporaryStorage)
		assert.Nil(t, payload.PrivatelyOwnedVehicle)

		/* These fields are defaulting to zero if they are nil in the model */
		assert.Equal(t, int64(0), payload.StorageInTransit)
		assert.Equal(t, int64(0), payload.TotalDependents)
		assert.Equal(t, int64(0), payload.TotalWeight)
	})

	t.Run("Success - Returns the entitlement payload with all optional fields populated", func(t *testing.T) {
		entitlement := models.Entitlement{
			ID:                             uuid.Must(uuid.NewV4()),
			DependentsAuthorized:           handlers.FmtBool(true),
			TotalDependents:                handlers.FmtInt(2),
			NonTemporaryStorage:            handlers.FmtBool(true),
			PrivatelyOwnedVehicle:          handlers.FmtBool(true),
			DBAuthorizedWeight:             handlers.FmtInt(10000),
			StorageInTransit:               handlers.FmtInt(45),
			RequiredMedicalEquipmentWeight: 500,
			OrganizationalClothingAndIndividualEquipment: true,
			ProGearWeight:       1000,
			ProGearWeightSpouse: 750,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		// TotalWeight needs to read from the internal weightAllotment, in this case 7000 lbs w/o dependents and
		// 9000 lbs with dependents
		entitlement.SetWeightAllotment(string(models.ServiceMemberRankE5))

		payload := Entitlement(&entitlement)

		assert.Equal(t, strfmt.UUID(entitlement.ID.String()), payload.ID)
		assert.True(t, *payload.DependentsAuthorized)
		assert.Equal(t, int64(2), payload.TotalDependents)
		assert.True(t, *payload.NonTemporaryStorage)
		assert.True(t, *payload.PrivatelyOwnedVehicle)
		assert.Equal(t, int64(10000), *payload.AuthorizedWeight)
		assert.Equal(t, int64(9000), payload.TotalWeight)
		assert.Equal(t, int64(45), payload.StorageInTransit)
		assert.Equal(t, int64(500), payload.RequiredMedicalEquipmentWeight)
		assert.Equal(t, true, payload.OrganizationalClothingAndIndividualEquipment)
		assert.Equal(t, int64(1000), payload.ProGearWeight)
		assert.Equal(t, int64(750), payload.ProGearWeightSpouse)
		assert.NotEmpty(t, payload.ETag)
		assert.Equal(t, etag.GenerateEtag(entitlement.UpdatedAt), payload.ETag)
	})

	t.Run("Success - Returns the entitlement payload with total weight self when dependents are not authorized", func(t *testing.T) {
		entitlement := models.Entitlement{
			ID:                             uuid.Must(uuid.NewV4()),
			DependentsAuthorized:           handlers.FmtBool(false),
			TotalDependents:                handlers.FmtInt(2),
			NonTemporaryStorage:            handlers.FmtBool(true),
			PrivatelyOwnedVehicle:          handlers.FmtBool(true),
			DBAuthorizedWeight:             handlers.FmtInt(10000),
			StorageInTransit:               handlers.FmtInt(45),
			RequiredMedicalEquipmentWeight: 500,
			OrganizationalClothingAndIndividualEquipment: true,
			ProGearWeight:       1000,
			ProGearWeightSpouse: 750,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		// TotalWeight needs to read from the internal weightAllotment, in this case 7000 lbs w/o dependents and
		// 9000 lbs with dependents
		entitlement.SetWeightAllotment(string(models.ServiceMemberRankE5))

		payload := Entitlement(&entitlement)

		assert.Equal(t, strfmt.UUID(entitlement.ID.String()), payload.ID)
		assert.False(t, *payload.DependentsAuthorized)
		assert.Equal(t, int64(2), payload.TotalDependents)
		assert.True(t, *payload.NonTemporaryStorage)
		assert.True(t, *payload.PrivatelyOwnedVehicle)
		assert.Equal(t, int64(10000), *payload.AuthorizedWeight)
		assert.Equal(t, int64(7000), payload.TotalWeight)
		assert.Equal(t, int64(45), payload.StorageInTransit)
		assert.Equal(t, int64(500), payload.RequiredMedicalEquipmentWeight)
		assert.Equal(t, true, payload.OrganizationalClothingAndIndividualEquipment)
		assert.Equal(t, int64(1000), payload.ProGearWeight)
		assert.Equal(t, int64(750), payload.ProGearWeightSpouse)
		assert.NotEmpty(t, payload.ETag)
		assert.Equal(t, etag.GenerateEtag(entitlement.UpdatedAt), payload.ETag)
	})
}
