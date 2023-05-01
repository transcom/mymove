package factory

import (
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func BuildWeightTicket(db *pop.Connection, customs []Customization, traits []Trait) models.WeightTicket {
	customs = setupCustomizations(customs, traits)

	// Find upload assertion and convert to models upload
	var cWeightTicket models.WeightTicket
	if result := findValidCustomization(customs, WeightTicket); result != nil {
		cWeightTicket = result.Model.(models.WeightTicket)

		if result.LinkOnly {
			return cWeightTicket
		}
	}

	ppmShipment := BuildPPMShipment(db, customs, traits)

	emptyUploadCustoms := []Customization{}
	fullUploadCustoms := []Customization{}
	trailerUploadCustoms := []Customization{}

	if db != nil {
		emptyUploadCustoms = append(emptyUploadCustoms, Customization{
			Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
			LinkOnly: true,
		})
		fullUploadCustoms = append(fullUploadCustoms, Customization{
			Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
			LinkOnly: true,
		})
		trailerUploadCustoms = append(trailerUploadCustoms, Customization{
			Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
			LinkOnly: true,
		})

		// Find upload assertion and convert to models upload
		var cUserUploadParams *UserUploadExtendedParams
		result := findValidCustomization(customs, UserUpload)
		if result != nil {
			if result.LinkOnly {
				log.Panic("Cannot provide LinkOnly UserUpload to BuildWeightTicket")
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
			// need to build our own params to override the File
			emptyUploadCustoms = append(emptyUploadCustoms,
				Customization{
					Model: models.UserUpload{},
					ExtendedParams: &UserUploadExtendedParams{
						UserUploader: cUserUploadParams.UserUploader,
						AppContext:   cUserUploadParams.AppContext,
						File:         testdatagen.Fixture("empty-weight-ticket.png"),
					},
				})
			fullUploadCustoms = append(fullUploadCustoms,
				Customization{
					Model: models.UserUpload{},
					ExtendedParams: &UserUploadExtendedParams{
						UserUploader: cUserUploadParams.UserUploader,
						AppContext:   cUserUploadParams.AppContext,
						File:         testdatagen.Fixture("full-weight-ticket.png"),
					},
				})

			trailerUploadCustoms = append(trailerUploadCustoms,
				Customization{
					Model:          models.UserUpload{},
					ExtendedParams: cUserUploadParams,
				})
		}

	}
	emptyUpload := BuildUserUpload(db, emptyUploadCustoms, nil)
	fullUpload := BuildUserUpload(db, fullUploadCustoms, nil)
	trailerUpload := BuildUserUpload(db, trailerUploadCustoms, nil)

	emptyUpload.Document.UserUploads = []models.UserUpload{emptyUpload}
	fullUpload.Document.UserUploads = []models.UserUpload{fullUpload}
	trailerUpload.Document.UserUploads = []models.UserUpload{trailerUpload}

	emptyWeight := unit.Pound(14500)
	fullWeight := emptyWeight + unit.Pound(4000)

	weightTicket := models.WeightTicket{
		VehicleDescription:                models.StringPointer("2022 Honda CR-V Hybrid"),
		EmptyWeight:                       &emptyWeight,
		MissingEmptyWeightTicket:          models.BoolPointer(false),
		FullWeight:                        &fullWeight,
		MissingFullWeightTicket:           models.BoolPointer(false),
		OwnsTrailer:                       models.BoolPointer(false),
		TrailerMeetsCriteria:              models.BoolPointer(false),
		PPMShipmentID:                     ppmShipment.ID,
		PPMShipment:                       ppmShipment,
		EmptyDocumentID:                   *emptyUpload.DocumentID,
		EmptyDocument:                     emptyUpload.Document,
		FullDocumentID:                    *fullUpload.DocumentID,
		FullDocument:                      fullUpload.Document,
		ProofOfTrailerOwnershipDocumentID: *trailerUpload.DocumentID,
		ProofOfTrailerOwnershipDocument:   trailerUpload.Document,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&weightTicket, cWeightTicket)

	if db != nil {
		mustCreate(db, &weightTicket)
	}

	return weightTicket

}
