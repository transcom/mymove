package factory

import (
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func BuildProgearWeightTicket(db *pop.Connection, customs []Customization, traits []Trait) models.ProgearWeightTicket {
	customs = setupCustomizations(customs, traits)

	// Find upload assertion and convert to models upload
	var cProgearWeightTicket models.ProgearWeightTicket
	if result := findValidCustomization(customs, ProgearWeightTicket); result != nil {
		cProgearWeightTicket = result.Model.(models.ProgearWeightTicket)

		if result.LinkOnly {
			return cProgearWeightTicket
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
				log.Panic("Cannot provide LinkOnly UserUpload to BuildProgearWeightTicket")
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

	description := "professional equipment"
	progearWeightTicket := models.ProgearWeightTicket{
		PPMShipmentID:    ppmShipment.ID,
		PPMShipment:      ppmShipment,
		DocumentID:       *upload.DocumentID,
		Document:         upload.Document,
		BelongsToSelf:    models.BoolPointer(true),
		Description:      &description,
		HasWeightTickets: models.BoolPointer(true),
		Weight:           models.PoundPointer(unit.Pound(500)),
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&progearWeightTicket, cProgearWeightTicket)

	if cProgearWeightTicket.BelongsToSelf != nil {
		progearWeightTicket.BelongsToSelf = cProgearWeightTicket.BelongsToSelf
	}

	if db != nil {
		mustCreate(db, &progearWeightTicket)
	}

	return progearWeightTicket

}
