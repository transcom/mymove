package ghcimport

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

type stringSet map[string]struct{}

type serviceAreasSet map[string]domesticServiceAreaData

type domesticServiceAreaData struct {
	serviceArea *models.ReDomesticServiceArea
	zip3s       stringSet
}

func (gre *GHCRateEngineImporter) importREDomesticServiceArea(dbTx *pop.Connection) error {
	// Read all the staged domestic service areas
	var stageDomesticServiceAreas []models.StageDomesticServiceArea
	stageErr := dbTx.All(&stageDomesticServiceAreas)
	if stageErr != nil {
		return errors.Wrap(stageErr, "Could not read staged domestic service areas")
	}

	// Since we appear to have duplicate service areas (and possibly zips too), consolidate data to
	// insert in a separate data structure that we'll then walk for inserting rows into the DB.
	serviceAreasData := make(serviceAreasSet)
	for _, stageArea := range stageDomesticServiceAreas {
		serviceAreaNumber, err := cleanServiceAreaNumber(stageArea.ServiceAreaNumber)
		if err != nil {
			return errors.Wrapf(err, "Could not process service area number [%s]", stageArea.ServiceAreaNumber)
		}

		splitZip3s := strings.Split(stageArea.Zip3s, ",")
		for i, zip3 := range splitZip3s {
			splitZip3s[i], err = cleanZip3(zip3)
			if err != nil {
				return errors.Wrapf(err, "Could not process zip3 [%s]", zip3)
			}
		}

		foundServiceAreaData, serviceAreaFound := serviceAreasData[serviceAreaNumber]
		if serviceAreaFound {
			gre.Logger.Info("Service area already exists",
				zap.String("service area number", serviceAreaNumber),
				zap.String("existing base point city", foundServiceAreaData.serviceArea.BasePointCity),
				zap.String("existing state", foundServiceAreaData.serviceArea.State),
				zap.String("new base point city", stageArea.BasePointCity),
				zap.String("new state", stageArea.State))

			// TODO: We're not storing the additional city/state right now.  We may
			//   want a separate story to consider putting that in the re_zip3s table
			//   if we think it's important (we don't currently use it for pricing).

			// Add zips to existing set
			zip3sSet := foundServiceAreaData.zip3s
			for _, zip3 := range splitZip3s {
				if _, zipFound := zip3sSet[zip3]; zipFound {
					gre.Logger.Info("Zip3 already exists for service area",
						zap.String("service area number", serviceAreaNumber),
						zap.String("zip3", zip3))
				}
				zip3sSet[zip3] = struct{}{}
			}
		} else {
			zip3sSet := make(stringSet)
			for _, zip3 := range splitZip3s {
				zip3sSet[zip3] = struct{}{}
			}

			serviceAreasData[serviceAreaNumber] = domesticServiceAreaData{
				serviceArea: &models.ReDomesticServiceArea{
					BasePointCity: stageArea.BasePointCity,
					State:         stageArea.State,
					ServiceArea:   serviceAreaNumber,
					// Fill in services schedule and SIT PD schedule later from other tab.
				},
				zip3s: zip3sSet,
			}
		}
	}

	// Now walk tab 2b and record the two schedule values.
	var stageDomesticServiceAreaPrices []models.StageDomesticServiceAreaPrice
	stageErr = dbTx.All(&stageDomesticServiceAreaPrices)
	if stageErr != nil {
		return errors.Wrap(stageErr, "Could not read staged domestic service area prices")
	}

	for _, stagePrice := range stageDomesticServiceAreaPrices {
		serviceAreaNumber, err := cleanServiceAreaNumber(stagePrice.ServiceAreaNumber)
		if err != nil {
			return errors.Wrapf(err, "Could not process service area number [%s]", stagePrice.ServiceAreaNumber)
		}

		foundServiceAreaData, serviceAreaFound := serviceAreasData[serviceAreaNumber]
		if !serviceAreaFound {
			return errors.New(fmt.Sprintf("Missing service area [%s] in list of service areas", serviceAreaNumber))
		}

		servicesSchedule, err := stringToInteger(stagePrice.ServicesSchedule)
		if err != nil {
			return errors.Wrapf(err, "Could not process services schedule [%s]", stagePrice.ServicesSchedule)
		}

		sitSchedule, err := stringToInteger(stagePrice.SITPickupDeliverySchedule)
		if err != nil {
			return errors.Wrapf(err, "Could not process SIT P/D schedule [%s]", stagePrice.SITPickupDeliverySchedule)
		}

		serviceArea := foundServiceAreaData.serviceArea
		serviceArea.ServicesSchedule = servicesSchedule
		serviceArea.SITPDSchedule = sitSchedule
	}

	// Now walk our data structure and insert records, recording a map of UUIDs along the way.
	serviceAreaToIDMap := make(map[string]uuid.UUID)
	for _, serviceAreaData := range serviceAreasData {
		// Save the service area record.
		reServiceArea := serviceAreaData.serviceArea
		verrs, err := dbTx.ValidateAndSave(reServiceArea)
		if err != nil {
			return errors.Wrapf(err, "Could not save service area: %+v", *reServiceArea)
		}
		if verrs.HasAny() {
			return errors.Wrapf(verrs, "Validation errors when saving contract: %+v", *reServiceArea)
		}
		serviceAreaToIDMap[reServiceArea.ServiceArea] = reServiceArea.ID

		// Save the associated zips.
		for zip3 := range serviceAreaData.zip3s {
			reZip3 := models.ReZip3{
				Zip3:                  zip3,
				DomesticServiceAreaID: reServiceArea.ID,
			}
			verrs, err := dbTx.ValidateAndSave(&reZip3)
			if err != nil {
				return errors.Wrapf(err, "Could not save zip3: %+v", reZip3)
			}
			if verrs.HasAny() {
				return errors.Wrapf(verrs, "Validation errors when saving zip3: %+v", reZip3)
			}
		}
	}

	gre.serviceAreaToIDMap = serviceAreaToIDMap

	return nil
}
