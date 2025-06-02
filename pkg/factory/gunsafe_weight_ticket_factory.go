package factory

import (
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func BuildGunSafeWeightTicket(db *pop.Connection, customs []Customization, traits []Trait) models.GunSafeWeightTicket {
	customs = setupCustomizations(customs, traits)

	// Find upload assertion and convert to models upload
	var cGunSafeWeightTicket models.GunSafeWeightTicket
	if result := findValidCustomization(customs, GunSafeWeightTicket); result != nil {
		cGunSafeWeightTicket = result.Model.(models.GunSafeWeightTicket)

		if result.LinkOnly {
			return cGunSafeWeightTicket
		}
	}

	ppmShipment := BuildPPMShipment(db, customs, traits)

	uploadCustoms := []Customization{}

	if db != nil {
		uploadCustoms = append(uploadCustoms, Customization{
			Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
			LinkOnly: true,
		})

		// Find upload assertion and convert to models upload
		var cUserUploadParams *UserUploadExtendedParams
		result := findValidCustomization(customs, UserUpload)
		if result != nil {
			if result.LinkOnly {
				log.Panic("Cannot provide LinkOnly UserUpload to BuildGunSafeWeightTicket")
			}

			// If extendedParams were provided, extract them
			typedResult, ok := result.ExtendedParams.(*UserUploadExtendedParams)
			if result.ExtendedParams != nil && !ok {
				log.Panic("To create UserUpload model, ExtendedParams must be nil or a pointer to UserUploadExtendedParams")
			}
			cUserUploadParams = typedResult
		}

		// As of 2023-04-28, no caller customizes the documents.
		// Supporting that fully gets pretty complicated, so only
		// support defaults for now
		if cUserUploadParams != nil {
			uploadCustoms = append(uploadCustoms,
				Customization{
					Model:          models.UserUpload{},
					ExtendedParams: cUserUploadParams,
				})
		}

	}
	upload := BuildUserUpload(db, uploadCustoms, nil)

	upload.Document.UserUploads = []models.UserUpload{upload}

	description := "gun safe"
	gunSafeWeightTicket := models.GunSafeWeightTicket{
		PPMShipmentID:    ppmShipment.ID,
		PPMShipment:      ppmShipment,
		DocumentID:       *upload.DocumentID,
		Document:         upload.Document,
		Description:      &description,
		HasWeightTickets: models.BoolPointer(true),
		Weight:           models.PoundPointer(unit.Pound(500)),
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&gunSafeWeightTicket, cGunSafeWeightTicket)

	if db != nil {
		mustCreate(db, &gunSafeWeightTicket)
	}

	return gunSafeWeightTicket

}
