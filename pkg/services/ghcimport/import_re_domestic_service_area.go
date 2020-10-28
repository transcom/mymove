package ghcimport

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// A set of service area strings and their associated data
type serviceAreasSet map[string]domesticServiceAreaData

// The ReDomesticServiceArea model and the related zip3 set
type domesticServiceAreaData struct {
	serviceArea *models.ReDomesticServiceArea
	zip3s       zip3Set
}

// A set of zip3s and their associated city/state
type zip3Set map[string]cityState

// City/state data we need for each zip3
type cityState struct {
	city  string
	state string
}

func (gre *GHCRateEngineImporter) importREDomesticServiceArea(dbTx *pop.Connection) error {
	// Build a data structure that we'll walk for inserting rows into the DB (we have
	// duplicate service areas that we need to handle).
	serviceAreasData, err := gre.buildServiceAreasData(dbTx)
	if err != nil {
		return err
	}

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
		return nil, fmt.Errorf("could not read staged domestic service areas: %w", stageErr)
	}

	// Since we appear to have duplicate service areas (and possibly zips too), consolidate data to
	// insert in a separate data structure that we'll then walk for inserting rows into the DB.
	serviceAreasData := make(serviceAreasSet)
	for _, stageArea := range stageDomesticServiceAreas {
		serviceAreaNumber, err := cleanServiceAreaNumber(stageArea.ServiceAreaNumber)
		if err != nil {
			return nil, fmt.Errorf("could not process service area number [%s]: %w", stageArea.ServiceAreaNumber, err)
		}

		splitZip3s := strings.Split(stageArea.Zip3s, ",")
		for i, zip3 := range splitZip3s {
			splitZip3s[i], err = cleanZip3(zip3)
			if err != nil {
				return nil, fmt.Errorf("could not process zip3 [%s]: %w", zip3, err)
			}
		}

		foundServiceAreaData, serviceAreaFound := serviceAreasData[serviceAreaNumber]
		if serviceAreaFound {
			// Service area already encountered; merge new zips into data structure
			zip3sSet := foundServiceAreaData.zip3s
			for _, zip3 := range splitZip3s {
				if _, zipFound := zip3sSet[zip3]; !zipFound {
					zip3sSet[zip3] = cityState{
						city:  stageArea.BasePointCity,
						state: stageArea.State,
					}
				}
			}
		} else {
			// Service area has not been seen yet; create entire data structure
			zip3sSet := make(zip3Set)
			for _, zip3 := range splitZip3s {
				zip3sSet[zip3] = cityState{
					city:  stageArea.BasePointCity,
					state: stageArea.State,
				}
			}

			serviceAreasData[serviceAreaNumber] = domesticServiceAreaData{
				serviceArea: &models.ReDomesticServiceArea{
					ContractID:  gre.ContractID,
					ServiceArea: serviceAreaNumber,
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
		return fmt.Errorf("could not read staged domestic service area prices: %w", stageErr)
	}

	for _, stagePrice := range stageDomesticServiceAreaPrices {
		serviceAreaNumber, err := cleanServiceAreaNumber(stagePrice.ServiceAreaNumber)
		if err != nil {
			return fmt.Errorf("could not process service area number [%s]: %w", stagePrice.ServiceAreaNumber, err)
		}

		foundServiceAreaData, serviceAreaFound := serviceAreasData[serviceAreaNumber]
		if !serviceAreaFound {
			return fmt.Errorf("missing service area [%s] in list of service areas", serviceAreaNumber)
		}

		servicesSchedule, err := stringToInteger(stagePrice.ServicesSchedule)
		if err != nil {
			return fmt.Errorf("could not process services schedule [%s]: %w", stagePrice.ServicesSchedule, err)
		}

		sitSchedule, err := stringToInteger(stagePrice.SITPickupDeliverySchedule)
		if err != nil {
			return fmt.Errorf("could not process SIT P/D schedule [%s]: %w", stagePrice.SITPickupDeliverySchedule, err)
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
		err := dbTx.
			Where("contract_id = ?", gre.ContractID).
			Where("service_area = ?", reServiceArea.ServiceArea).
			All(&existingServiceAreas)
		if err != nil {
			return nil, fmt.Errorf("could not lookup existing service area [%s]: %w", reServiceArea.ServiceArea, err)
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
			if verrs.HasAny() {
				return nil, fmt.Errorf("validation errors when saving service area [%+v]: %w", *reServiceArea, verrs)
			}
			if saveErr != nil {
				return nil, fmt.Errorf("could not save service area [%+v]: %w", *reServiceArea, saveErr)
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

func (gre *GHCRateEngineImporter) saveZip3sForServiceArea(dbTx *pop.Connection, zip3s zip3Set, serviceAreaID uuid.UUID) error {
	// Save the associated zips.
	for zip3, cityState := range zip3s {
		reZip3 := models.ReZip3{
			ContractID:            gre.ContractID,
			Zip3:                  zip3,
			BasePointCity:         cityState.city,
			State:                 cityState.state,
			DomesticServiceAreaID: serviceAreaID,
		}

		// See if there is an existing zip3 record.  If so, we need to update it.
		var existingZip3s models.ReZip3s
		err := dbTx.
			Where("contract_id = ?", gre.ContractID).
			Where("zip3 = ?", reZip3.Zip3).
			All(&existingZip3s)
		if err != nil {
			return fmt.Errorf("could not lookup existing zip3 [%s]: %w", reZip3.Zip3, err)
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
			if verrs.HasAny() {
				return fmt.Errorf("validation errors when saving zip3 [%+v]: %w", reZip3, verrs)
			}
			if saveErr != nil {
				return fmt.Errorf("could not save zip3 [%+v]: %w", reZip3, saveErr)
			}
		}
	}

	return nil
}

func updateExistingServiceArea(existing *models.ReDomesticServiceArea, new models.ReDomesticServiceArea) bool {
	if existing.ServicesSchedule == new.ServicesSchedule &&
		existing.SITPDSchedule == new.SITPDSchedule {
		return false
	}

	existing.ServicesSchedule = new.ServicesSchedule
	existing.SITPDSchedule = new.SITPDSchedule
	return true
}

func updateExistingZip3(existing *models.ReZip3, new models.ReZip3) bool {
	if existing.BasePointCity == new.BasePointCity &&
		existing.State == new.State &&
		existing.DomesticServiceAreaID == new.DomesticServiceAreaID {
		return false
	}

	existing.BasePointCity = new.BasePointCity
	existing.State = new.State
	existing.DomesticServiceAreaID = new.DomesticServiceAreaID
	return true
}
