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
	// Build a data structure that we'll walk for inserting rows into the DB (we have
	// duplicate service areas that we need to handle).
	serviceAreasData, err := gre.buildServiceAreasData(dbTx)
	if err != nil {
		return err
	}

	// TODO: We aren't currently cleaning up any records that may not be part of the latest
	//   import.  That may get tricky with referential integrity.

	// Insert the service area and zip3 records; store a map of service area UUIDs to use
	// in later imports.
	gre.serviceAreaToIDMap, err = gre.saveServiceAreasAndZip3s(dbTx, serviceAreasData)
	if err != nil {
		return err
	}

	return nil
}

func (gre *GHCRateEngineImporter) buildServiceAreasData(dbTx *pop.Connection) (serviceAreasSet, error) {
	// Read all the staged domestic service areas
	var stageDomesticServiceAreas []models.StageDomesticServiceArea
	stageErr := dbTx.All(&stageDomesticServiceAreas)
	if stageErr != nil {
		return nil, errors.Wrap(stageErr, "could not read staged domestic service areas")
	}

	// Since we appear to have duplicate service areas (and possibly zips too), consolidate data to
	// insert in a separate data structure that we'll then walk for inserting rows into the DB.
	serviceAreasData := make(serviceAreasSet)
	for _, stageArea := range stageDomesticServiceAreas {
		serviceAreaNumber, err := cleanServiceAreaNumber(stageArea.ServiceAreaNumber)
		if err != nil {
			return nil, errors.Wrapf(err, "could not process service area number [%s]", stageArea.ServiceAreaNumber)
		}

		splitZip3s := strings.Split(stageArea.Zip3s, ",")
		for i, zip3 := range splitZip3s {
			splitZip3s[i], err = cleanZip3(zip3)
			if err != nil {
				return nil, errors.Wrapf(err, "could not process zip3 [%s]", zip3)
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

	// Now add schedule data (which comes from a different tab) to our service area structure.
	err := gre.addScheduleData(dbTx, serviceAreasData)
	if err != nil {
		return nil, err
	}

	return serviceAreasData, nil
}

func (gre *GHCRateEngineImporter) addScheduleData(dbTx *pop.Connection, serviceAreasData serviceAreasSet) error {
	// Now walk tab 2b and record the two schedule values.
	var stageDomesticServiceAreaPrices []models.StageDomesticServiceAreaPrice
	stageErr := dbTx.All(&stageDomesticServiceAreaPrices)
	if stageErr != nil {
		return errors.Wrap(stageErr, "could not read staged domestic service area prices")
	}

	for _, stagePrice := range stageDomesticServiceAreaPrices {
		serviceAreaNumber, err := cleanServiceAreaNumber(stagePrice.ServiceAreaNumber)
		if err != nil {
			return errors.Wrapf(err, "could not process service area number [%s]", stagePrice.ServiceAreaNumber)
		}

		foundServiceAreaData, serviceAreaFound := serviceAreasData[serviceAreaNumber]
		if !serviceAreaFound {
			return fmt.Errorf("missing service area [%s] in list of service areas", serviceAreaNumber)
		}

		servicesSchedule, err := stringToInteger(stagePrice.ServicesSchedule)
		if err != nil {
			return errors.Wrapf(err, "could not process services schedule [%s]", stagePrice.ServicesSchedule)
		}

		sitSchedule, err := stringToInteger(stagePrice.SITPickupDeliverySchedule)
		if err != nil {
			return errors.Wrapf(err, "could not process SIT P/D schedule [%s]", stagePrice.SITPickupDeliverySchedule)
		}

		serviceArea := foundServiceAreaData.serviceArea
		serviceArea.ServicesSchedule = servicesSchedule
		serviceArea.SITPDSchedule = sitSchedule
	}

	return nil
}

func (gre *GHCRateEngineImporter) saveServiceAreasAndZip3s(dbTx *pop.Connection, serviceAreasData serviceAreasSet) (map[string]uuid.UUID, error) {
	serviceAreaToIDMap := make(map[string]uuid.UUID)
	for _, serviceAreaData := range serviceAreasData {
		reServiceArea := serviceAreaData.serviceArea

		// See if there is an existing service area record.  If so, we may need to update it.
		var existingServiceAreas models.ReDomesticServiceAreas
		err := dbTx.Where("service_area = ?", reServiceArea.ServiceArea).All(&existingServiceAreas)
		if err != nil {
			return nil, errors.Wrapf(err, "could not lookup existing service area [%s]", reServiceArea.ServiceArea)
		}
		doSaveServiceArea := true
		if len(existingServiceAreas) > 0 {
			// Update existing service area with new data.
			existingServiceArea := existingServiceAreas[0]
			doSaveServiceArea = updateExistingServiceArea(&existingServiceArea, *reServiceArea)
			reServiceArea = &existingServiceArea
		}

		if doSaveServiceArea {
			verrs, saveErr := dbTx.ValidateAndSave(reServiceArea)
			if saveErr != nil {
				return nil, errors.Wrapf(saveErr, "could not save service area: %+v", *reServiceArea)
			}
			if verrs.HasAny() {
				return nil, errors.Wrapf(verrs, "validation errors when saving contract: %+v", *reServiceArea)
			}
		}
		serviceAreaToIDMap[reServiceArea.ServiceArea] = reServiceArea.ID

		err = gre.saveZip3sForServiceArea(dbTx, serviceAreaData.zip3s, reServiceArea.ID)
		if err != nil {
			return nil, err
		}
	}

	return serviceAreaToIDMap, nil
}

func (gre *GHCRateEngineImporter) saveZip3sForServiceArea(dbTx *pop.Connection, zip3s stringSet, serviceAreaID uuid.UUID) error {
	// Save the associated zips.
	for zip3 := range zip3s {
		reZip3 := models.ReZip3{
			Zip3:                  zip3,
			DomesticServiceAreaID: serviceAreaID,
		}

		// See if there is an existing zip3 record.  If so, we need to update it.
		var existingZip3s models.ReZip3s
		err := dbTx.Where("zip3 = ?", reZip3.Zip3).All(&existingZip3s)
		if err != nil {
			return errors.Wrapf(err, "could not lookup existing zip3 [%s]", reZip3.Zip3)
		}
		doSaveZip3 := true
		if len(existingZip3s) > 0 {
			// Update existing zip3 with new data.
			existingZip3 := existingZip3s[0]
			doSaveZip3 = updateExistingZip3(&existingZip3, reZip3)
			reZip3 = existingZip3
		}

		if doSaveZip3 {
			verrs, saveErr := dbTx.ValidateAndSave(&reZip3)
			if saveErr != nil {
				return errors.Wrapf(saveErr, "could not save zip3: %+v", reZip3)
			}
			if verrs.HasAny() {
				return errors.Wrapf(verrs, "validation errors when saving zip3: %+v", reZip3)
			}
		}
	}

	return nil
}

func updateExistingServiceArea(existing *models.ReDomesticServiceArea, new models.ReDomesticServiceArea) bool {
	if existing.BasePointCity == new.BasePointCity &&
		existing.State == new.State &&
		existing.ServicesSchedule == new.ServicesSchedule &&
		existing.SITPDSchedule == new.SITPDSchedule {
		return false
	}

	existing.BasePointCity = new.BasePointCity
	existing.State = new.State
	existing.ServicesSchedule = new.ServicesSchedule
	existing.SITPDSchedule = new.SITPDSchedule
	return true
}

func updateExistingZip3(existing *models.ReZip3, new models.ReZip3) bool {
	if existing.DomesticServiceAreaID == new.DomesticServiceAreaID {
		return false
	}

	existing.DomesticServiceAreaID = new.DomesticServiceAreaID
	return true
}
