package factory

import (
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func BuildMovingExpense(db *pop.Connection, customs []Customization, traits []Trait) models.MovingExpense {
	customs = setupCustomizations(customs, traits)

	// Find upload assertion and convert to models upload
	var cMovingExpense models.MovingExpense
	if result := findValidCustomization(customs, MovingExpense); result != nil {
		cMovingExpense = result.Model.(models.MovingExpense)

		if result.LinkOnly {
			return cMovingExpense
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
				log.Panic("Cannot provide LinkOnly UserUpload to BuildMovingExpense")
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

	packingMaterialType := models.MovingExpenseReceiptTypePackingMaterials
	movingExpense := models.MovingExpense{
		PPMShipmentID:     ppmShipment.ID,
		PPMShipment:       ppmShipment,
		DocumentID:        *upload.DocumentID,
		Document:          upload.Document,
		MovingExpenseType: &packingMaterialType,
		Description:       models.StringPointer("Packing Peanuts"),
		PaidWithGTCC:      models.BoolPointer(true),
		Amount:            models.CentPointer(2345),
		MissingReceipt:    models.BoolPointer(false),
	}

	// MergeModels is not working for PaidWithGTCC so overriding here
	if cMovingExpense.PaidWithGTCC != nil {
		movingExpense.PaidWithGTCC = cMovingExpense.PaidWithGTCC
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&movingExpense, cMovingExpense)

	if db != nil {
		mustCreate(db, &movingExpense)
	}

	return movingExpense

}
