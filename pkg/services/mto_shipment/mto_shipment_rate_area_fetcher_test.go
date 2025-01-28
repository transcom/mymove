package mtoshipment

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

const testContractCode = "TEST"
const testContractName = "Test Contract"
const fairbanksAlaskaPostalCode = "99716"
const anchorageAlaskaPostalCode = "99521"
const wasillaAlaskaPostalCode = "99652"
const beverlyHillsCAPostalCode = "90210"
const sanDiegoCAPostalCode = "92075"
const brooklynNYPostalCode = "11220"

func (suite *MTOShipmentServiceSuite) TestGetMoveShipmentRateArea() {
	shipmentRateAreaFetcher := NewMTOShipmentRateAreaFetcher()

	suite.Run("test mapping of one rateArea to many postCodes and one rateArea to one", func() {
		availableToPrimeAtTime := time.Now().Add(-500 * time.Hour)
		testMove := models.Move{
			AvailableToPrimeAt: &availableToPrimeAtTime,
			MTOShipments: models.MTOShipments{
				models.MTOShipment{
					PickupAddress: &models.Address{
						StreetAddress1: "123 Main St",
						City:           "Fairbanks",
						State:          "AK",
						PostalCode:     fairbanksAlaskaPostalCode,
						IsOconus:       models.BoolPointer(true),
					},
					DestinationAddress: &models.Address{
						StreetAddress1: "123 Main St",
						City:           "Anchorage",
						State:          "AK",
						PostalCode:     anchorageAlaskaPostalCode,
						IsOconus:       models.BoolPointer(true),
					},
				},
				models.MTOShipment{
					PickupAddress: &models.Address{
						StreetAddress1: "123 Main St",
						City:           "Beverly Hills",
						State:          "CA",
						PostalCode:     beverlyHillsCAPostalCode,
						IsOconus:       models.BoolPointer(false),
					},
					DestinationAddress: &models.Address{
						StreetAddress1: "123 Main St",
						City:           "San Diego",
						State:          "CA",
						PostalCode:     sanDiegoCAPostalCode,
						IsOconus:       models.BoolPointer(false),
					},
				},
				models.MTOShipment{
					PPMShipment: &models.PPMShipment{
						PickupAddress: &models.Address{
							StreetAddress1: "123 Main St",
							City:           "Wasilla",
							State:          "AK",
							PostalCode:     wasillaAlaskaPostalCode,
							IsOconus:       models.BoolPointer(true),
						},
						DestinationAddress: &models.Address{
							StreetAddress1: "123 Main St",
							City:           "Wasilla",
							State:          "AK",
							PostalCode:     wasillaAlaskaPostalCode,
							IsOconus:       models.BoolPointer(true),
						},
					},
				},
			},
		}

		setupRateArea := func(contract models.ReContract) models.ReRateArea {
			rateAreaCode := uuid.Must(uuid.NewV4()).String()[0:5]
			rateArea := models.ReRateArea{
				ID:         uuid.Must(uuid.NewV4()),
				ContractID: contract.ID,
				IsOconus:   true,
				Code:       rateAreaCode,
				Name:       fmt.Sprintf("Alaska-%s", rateAreaCode),
				Contract:   contract,
			}
			verrs, err := suite.DB().ValidateAndCreate(&rateArea)
			if verrs.HasAny() {
				suite.Fail(verrs.Error())
			}
			if err != nil {
				suite.Fail(err.Error())
			}
			return rateArea
		}

		setupRateAreaToPostalCodeData := func(rateArea models.ReRateArea, postalCode string) models.ReRateArea {
			// fetch US by country id
			us_countryId := uuid.FromStringOrNil("c390ced2-89e1-418d-bbff-f8a79b89c4b6")
			us_country, err := models.FetchCountryByID(suite.DB(), us_countryId)
			suite.NotNil(us_country)
			suite.FatalNoError(err)

			usprc, err := findUsPostRegionCityByZipCode(suite.AppContextForTest(), postalCode)
			suite.NotNil(usprc)
			suite.FatalNoError(err)

			oconusRateArea := testOnlyOconusRateArea{
				ID:                 uuid.Must(uuid.NewV4()),
				RateAreaId:         rateArea.ID,
				CountryId:          us_country.ID,
				UsPostRegionCityId: usprc.ID,
				Active:             true,
			}
			verrs, err := suite.DB().ValidateAndCreate(&oconusRateArea)
			if verrs.HasAny() {
				suite.Fail(verrs.Error())
			}
			if err != nil {
				suite.Fail(err.Error())
			}

			return rateArea
		}

		setupRateAreaToManyPostalCodesData := func(contract models.ReContract, testPostalCode []string) models.ReRateArea {
			rateArea := setupRateArea(contract)
			for _, postalCode := range testPostalCode {
				setupRateAreaToPostalCodeData(rateArea, postalCode)
			}
			return rateArea
		}

		isRateAreaEquals := func(expectedRateArea models.ReRateArea, postalCode string, shipmentPostalCodeRateArea *[]services.ShipmentPostalCodeRateArea) bool {
			var shipmentPostalCodeRateAreaLookupMap = make(map[string]services.ShipmentPostalCodeRateArea)
			for _, i := range *shipmentPostalCodeRateArea {
				shipmentPostalCodeRateAreaLookupMap[i.PostalCode] = i
			}
			if _, ok := shipmentPostalCodeRateAreaLookupMap[postalCode]; !ok {
				return false
			}
			return (shipmentPostalCodeRateAreaLookupMap[postalCode].RateArea.ID == expectedRateArea.ID && shipmentPostalCodeRateAreaLookupMap[postalCode].RateArea.Name == expectedRateArea.Name && shipmentPostalCodeRateAreaLookupMap[postalCode].RateArea.Code == expectedRateArea.Code)
		}

		// create test contract
		contract, err := suite.createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		// setup contract year within availableToPrimeAtTime time
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate:  availableToPrimeAtTime,
				EndDate:    time.Now(),
				ContractID: contract.ID,
			},
		})

		// setup Fairbanks and Anchorage to have same RateArea
		rateArea1 := setupRateAreaToManyPostalCodesData(*contract, []string{fairbanksAlaskaPostalCode, anchorageAlaskaPostalCode})
		// setup Wasilla to have it's own RateArea
		rateArea2 := setupRateAreaToPostalCodeData(setupRateArea(*contract), wasillaAlaskaPostalCode)

		shipmentPostalCodeRateAreas, err := shipmentRateAreaFetcher.GetPrimeMoveShipmentRateAreas(suite.AppContextForTest(), testMove)
		suite.NotNil(shipmentPostalCodeRateAreas)
		suite.FatalNoError(err)
		suite.Equal(3, len(*shipmentPostalCodeRateAreas))

		suite.Equal(true, isRateAreaEquals(rateArea1, fairbanksAlaskaPostalCode, shipmentPostalCodeRateAreas))
		suite.Equal(true, isRateAreaEquals(rateArea1, anchorageAlaskaPostalCode, shipmentPostalCodeRateAreas))
		suite.Equal(true, isRateAreaEquals(rateArea2, wasillaAlaskaPostalCode, shipmentPostalCodeRateAreas))

		suite.Equal(false, isRateAreaEquals(rateArea2, fairbanksAlaskaPostalCode, shipmentPostalCodeRateAreas))
		suite.Equal(false, isRateAreaEquals(rateArea2, anchorageAlaskaPostalCode, shipmentPostalCodeRateAreas))
		suite.Equal(false, isRateAreaEquals(rateArea1, wasillaAlaskaPostalCode, shipmentPostalCodeRateAreas))

		// var shipmentPostalCodeRateAreasLookupMap = make(map[string]services.ShipmentPostalCodeRateArea)
		// for _, pcra := range *shipmentPostalCodeRateAreas {
		// 	shipmentPostalCodeRateAreasLookupMap[pcra.PostalCode] = pcra
		// }

		// suite.Equal("California-South", shipmentPostalCodeRateAreasLookupMap[beverlyHillsCAPostalCode].RateArea.Name)
		// suite.Equal("California-South", shipmentPostalCodeRateAreasLookupMap[sanDiegoCAPostalCode].RateArea.Name)
	})

	suite.Run("Returns matching CONUS rate areas", func() {
		availableToPrimeAtTime := time.Now().Add(-500 * time.Hour)
		testMove := models.Move{
			AvailableToPrimeAt: &availableToPrimeAtTime,
			MTOShipments: models.MTOShipments{
				models.MTOShipment{
					PickupAddress: &models.Address{
						StreetAddress1: "123 Main St",
						City:           "Beverly Hills",
						State:          "CA",
						PostalCode:     beverlyHillsCAPostalCode,
						IsOconus:       models.BoolPointer(false),
					},
					DestinationAddress: &models.Address{
						StreetAddress1: "123 Main St",
						City:           "San Diego",
						State:          "CA",
						PostalCode:     sanDiegoCAPostalCode,
						IsOconus:       models.BoolPointer(false),
					},
				},
				models.MTOShipment{
					PPMShipment: &models.PPMShipment{
						PickupAddress: &models.Address{
							StreetAddress1: "123 Main St",
							City:           "Brooklyn",
							State:          "NY",
							PostalCode:     brooklynNYPostalCode,
							IsOconus:       models.BoolPointer(false),
						},
						DestinationAddress: &models.Address{
							StreetAddress1: "123 Main St",
							City:           "Beverly Hills",
							State:          "CA",
							PostalCode:     beverlyHillsCAPostalCode,
							IsOconus:       models.BoolPointer(false),
						},
					},
				},
			},
		}

		// create test contract
		contract, err := suite.createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		// setup contract year within availableToPrimeAtTime time
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate:  availableToPrimeAtTime,
				EndDate:    time.Now(),
				ContractID: contract.ID,
			},
		})

		rateAreaCA := models.ReRateArea{
			ID:         uuid.Must(uuid.NewV4()),
			ContractID: contract.ID,
			IsOconus:   false,
			Code:       "US88",
			Name:       "California-South",
			Contract:   *contract,
		}
		verrs, err := suite.DB().ValidateAndCreate(&rateAreaCA)
		if verrs.HasAny() {
			suite.Fail(verrs.Error())
		}
		if err != nil {
			suite.Fail(err.Error())
		}

		rateAreaNY := models.ReRateArea{
			ID:         uuid.Must(uuid.NewV4()),
			ContractID: contract.ID,
			IsOconus:   false,
			Code:       "US17",
			Name:       "New York",
			Contract:   *contract,
		}
		verrs, err = suite.DB().ValidateAndCreate(&rateAreaNY)
		if verrs.HasAny() {
			suite.Fail(verrs.Error())
		}
		if err != nil {
			suite.Fail(err.Error())
		}

		domServiceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ContractID: contract.ID,
			},
		})

		brooklynZip3 := models.ReZip3{
			ID:                    uuid.Must(uuid.NewV4()),
			ContractID:            contract.ID,
			Contract:              *contract,
			Zip3:                  brooklynNYPostalCode[0:3],
			RateAreaID:            models.UUIDPointer(rateAreaNY.ID),
			HasMultipleRateAreas:  false,
			BasePointCity:         "Brooklyn",
			State:                 "NY",
			DomesticServiceAreaID: domServiceArea.ID,
		}
		verrs, err = suite.DB().ValidateAndCreate(&brooklynZip3)
		if verrs.HasAny() {
			suite.Fail(verrs.Error())
		}
		if err != nil {
			suite.Fail(err.Error())
		}

		sanDiegoZip3 := models.ReZip3{
			ID:                    uuid.Must(uuid.NewV4()),
			ContractID:            contract.ID,
			Contract:              *contract,
			Zip3:                  sanDiegoCAPostalCode[0:3],
			RateAreaID:            models.UUIDPointer(rateAreaCA.ID),
			HasMultipleRateAreas:  false,
			BasePointCity:         "San Diego",
			State:                 "CA",
			DomesticServiceAreaID: domServiceArea.ID,
		}
		verrs, err = suite.DB().ValidateAndCreate(&sanDiegoZip3)
		if verrs.HasAny() {
			suite.Fail(verrs.Error())
		}
		if err != nil {
			suite.Fail(err.Error())
		}

		beverlyHillsZip3 := models.ReZip3{
			ID:                    uuid.Must(uuid.NewV4()),
			ContractID:            contract.ID,
			Contract:              *contract,
			Zip3:                  beverlyHillsCAPostalCode[0:3],
			RateAreaID:            models.UUIDPointer(rateAreaCA.ID),
			HasMultipleRateAreas:  false,
			BasePointCity:         "Los Angeles",
			State:                 "CA",
			DomesticServiceAreaID: domServiceArea.ID,
		}

		verrs, err = suite.DB().ValidateAndCreate(&beverlyHillsZip3)
		if verrs.HasAny() {
			suite.Fail(verrs.Error())
		}
		if err != nil {
			suite.Fail(err.Error())
		}

		shipmentPostalCodeRateAreas, err := shipmentRateAreaFetcher.GetPrimeMoveShipmentRateAreas(suite.AppContextForTest(), testMove)
		suite.NotNil(shipmentPostalCodeRateAreas)
		suite.Equal(3, len(*shipmentPostalCodeRateAreas))
		suite.Nil(err)

		var shipmentPostalCodeRateAreaLookupMap = make(map[string]services.ShipmentPostalCodeRateArea)
		for _, pcra := range *shipmentPostalCodeRateAreas {
			shipmentPostalCodeRateAreaLookupMap[pcra.PostalCode] = pcra
		}

		suite.Equal("California-South", shipmentPostalCodeRateAreaLookupMap[beverlyHillsCAPostalCode].RateArea.Name)
		suite.Equal("California-South", shipmentPostalCodeRateAreaLookupMap[sanDiegoCAPostalCode].RateArea.Name)
		suite.Equal("New York", shipmentPostalCodeRateAreaLookupMap[brooklynNYPostalCode].RateArea.Name)
	})

	// suite.Run("Returns matching CONUS rate areas V2", func() {
	// 	availableToPrimeAtTime := time.Now().Add(-500 * time.Hour)
	// 	testMove := models.Move{
	// 		AvailableToPrimeAt: &availableToPrimeAtTime,
	// 		MTOShipments: models.MTOShipments{
	// 			models.MTOShipment{
	// 				PickupAddress: &models.Address{
	// 					StreetAddress1: "123 Main St",
	// 					City:           "Beverly Hills",
	// 					State:          "CA",
	// 					PostalCode:     beverlyHillsCAPostalCode,
	// 					IsOconus:       models.BoolPointer(false),
	// 				},
	// 				DestinationAddress: &models.Address{
	// 					StreetAddress1: "123 Main St",
	// 					City:           "San Diego",
	// 					State:          "CA",
	// 					PostalCode:     sanDiegoCAPostalCode,
	// 					IsOconus:       models.BoolPointer(false),
	// 				},
	// 			},
	// 			models.MTOShipment{
	// 				PPMShipment: &models.PPMShipment{
	// 					PickupAddress: &models.Address{
	// 						StreetAddress1: "123 Main St",
	// 						City:           "Brooklyn",
	// 						State:          "NY",
	// 						PostalCode:     brooklynNYPostalCode,
	// 						IsOconus:       models.BoolPointer(false),
	// 					},
	// 					DestinationAddress: &models.Address{
	// 						StreetAddress1: "123 Main St",
	// 						City:           "Beverly Hills",
	// 						State:          "CA",
	// 						PostalCode:     beverlyHillsCAPostalCode,
	// 						IsOconus:       models.BoolPointer(false),
	// 					},
	// 				},
	// 			},
	// 		},
	// 	}

	// 	setupDataForConusRateAreas := func(postalCode string, gbloc string) (models.ReRateArea, models.OconusRateArea, models.UsPostRegionCity, models.DutyLocation) {
	// 		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
	// 		rateAreaCode := uuid.Must(uuid.NewV4()).String()[0:5]
	// 		rateArea := testdatagen.FetchOrMakeReRateArea(suite.DB(), testdatagen.Assertions{
	// 			ReRateArea: models.ReRateArea{
	// 				ContractID: contract.ID,
	// 				IsOconus:   false,
	// 				Name:       fmt.Sprintf("Lower48-%s", rateAreaCode),
	// 				Contract:   contract,
	// 			},
	// 		})
	// 		suite.NotNil(rateArea)
	// 	}
	// })

	suite.Run("not available to prime error", func() {
		testMove := models.Move{
			MTOShipments: models.MTOShipments{
				models.MTOShipment{
					PickupAddress: &models.Address{
						StreetAddress1: "123 Main St",
						City:           "Fairbanks",
						State:          "AK",
						PostalCode:     fairbanksAlaskaPostalCode,
					},
					DestinationAddress: &models.Address{
						StreetAddress1: "123 Main St",
						City:           "Anchorage",
						State:          "AK",
						PostalCode:     anchorageAlaskaPostalCode,
					},
				},
			},
		}

		shipmentPostalCodeRateArea, err := shipmentRateAreaFetcher.GetPrimeMoveShipmentRateAreas(suite.AppContextForTest(), testMove)
		suite.Nil(shipmentPostalCodeRateArea)
		suite.NotNil(err)
		suite.IsType(apperror.UnprocessableEntityError{}, err)
	})

	suite.Run("contract for move not found", func() {
		availableToPrimeAtTime := time.Now().Add(-500 * time.Hour)
		testMove := models.Move{
			AvailableToPrimeAt: &availableToPrimeAtTime,
			MTOShipments: models.MTOShipments{
				models.MTOShipment{
					PickupAddress: &models.Address{
						StreetAddress1: "123 Main St",
						City:           "Fairbanks",
						State:          "AK",
						PostalCode:     fairbanksAlaskaPostalCode,
					},
					DestinationAddress: &models.Address{
						StreetAddress1: "123 Main St",
						City:           "Anchorage",
						State:          "AK",
						PostalCode:     anchorageAlaskaPostalCode,
					},
				},
			},
		}

		// create test contract
		contract, err := suite.createContract(suite.AppContextForTest(), testContractCode, testContractName)
		suite.NotNil(contract)
		suite.FatalNoError(err)

		// setup contract year within availableToPrimeAtTime time
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate:  time.Now(),
				EndDate:    time.Now().Add(5 * time.Hour),
				ContractID: contract.ID,
			},
		})

		shipmentPostalCodeRateArea, err := shipmentRateAreaFetcher.GetPrimeMoveShipmentRateAreas(suite.AppContextForTest(), testMove)
		suite.Nil(shipmentPostalCodeRateArea)
		suite.NotNil(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}

func (suite *MTOShipmentServiceSuite) TestFetchRateAreaByPostalCode() {
	// fetch US by country id
	us_countryId := uuid.FromStringOrNil("c390ced2-89e1-418d-bbff-f8a79b89c4b6")
	us_country, err := models.FetchCountryByID(suite.DB(), us_countryId)
	suite.NotNil(us_country)
	suite.FatalNoError(err)

	// create test contract
	contract, err := suite.createContract(suite.AppContextForTest(), testContractCode, testContractName)
	suite.NotNil(contract)
	suite.FatalNoError(err)

	// create rateArea associated to contract
	rateArea := models.ReRateArea{
		ID:         uuid.Must(uuid.NewV4()),
		ContractID: contract.ID,
		IsOconus:   true,
		Code:       "SomeAlaskaCode",
		Name:       "Alaska",
		Contract:   *contract,
	}
	verrs, err := suite.DB().ValidateAndCreate(&rateArea)
	if verrs.HasAny() {
		suite.Fail(verrs.Error())
	}
	if err != nil {
		suite.Fail(err.Error())
	}

	const alaskaPostalCode = "99506"

	usprc, err := findUsPostRegionCityByZipCode(suite.AppContextForTest(), alaskaPostalCode)
	suite.NotNil(usprc)
	suite.FatalNoError(err)

	oconusRateArea := testOnlyOconusRateArea{
		ID:                 uuid.Must(uuid.NewV4()),
		RateAreaId:         rateArea.ID,
		CountryId:          us_country.ID,
		UsPostRegionCityId: usprc.ID,
		Active:             true,
	}
	verrs, err = suite.DB().ValidateAndCreate(&oconusRateArea)
	if verrs.HasAny() {
		suite.Fail(verrs.Error())
	}
	if err != nil {
		suite.Fail(err.Error())
	}

	match, err := fetchOconusRateAreaByPostalCode(suite.AppContextForTest(), contract.ID, alaskaPostalCode)
	suite.NotNil(match)
	suite.FatalNoError(err)
}

func (suite *MTOShipmentServiceSuite) TestFetchRateAreaByPostalCodeNotFound() {
	_, err := fetchOconusRateAreaByPostalCode(suite.AppContextForTest(), uuid.FromStringOrNil("51393fa4-b31c-40fe-bedf-b692703c46eb"), "90210")
	suite.NotNil(err)
}

func (suite *MTOShipmentServiceSuite) TestFetchContract() {
	// create test contract
	expectedContract, err := suite.createContract(suite.AppContextForTest(), testContractCode, testContractName)
	suite.NotNil(expectedContract)
	suite.FatalNoError(err)

	time := time.Now().Add(-50 * time.Hour)
	testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			StartDate:  time,
			EndDate:    time,
			ContractID: expectedContract.ID,
		},
	})
	contract, err := fetchContract(suite.AppContextForTest(), time)
	suite.NotNil(contract)
	suite.Nil(err)
	suite.Equal(expectedContract.ID, contract.ID)
}

func (suite *MTOShipmentServiceSuite) TestFetchContractNotFound() {
	_, err := fetchContract(suite.AppContextForTest(), time.Now())
	suite.NotNil(err)
}

func (suite *MTOShipmentServiceSuite) createContract(appCtx appcontext.AppContext, contractCode string, contractName string) (*models.ReContract, error) {

	// See if contract code already exists.
	exists, err := appCtx.DB().Where("code = ?", testContractCode).Exists(&models.ReContract{})
	if err != nil {
		return nil, fmt.Errorf("could not determine if contract code [%s] existed: %w", testContractCode, err)
	}
	if exists {
		return nil, fmt.Errorf("the provided contract code [%s] already exists", testContractCode)
	}

	// Contract code is new; insert it.
	contract := models.ReContract{
		Code: contractCode,
		Name: contractName,
	}
	verrs, err := appCtx.DB().ValidateAndSave(&contract)
	if verrs.HasAny() {
		return nil, fmt.Errorf("validation errors when saving contract [%+v]: %w", contract, verrs)
	}
	if err != nil {
		return nil, fmt.Errorf("could not save contract [%+v]: %w", contract, err)
	}

	return &contract, nil
}

func findUsPostRegionCityByZipCode(appCtx appcontext.AppContext, zipCode string) (*models.UsPostRegionCity, error) {
	var usprc models.UsPostRegionCity
	err := appCtx.DB().Where("uspr_zip_id = ?", zipCode).First(&usprc)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, fmt.Errorf("No UsPostRegionCity found for provided zip code %s", zipCode)
		default:
			return nil, err
		}
	}
	return &usprc, nil
}

// **** This model is specifically for testing only to allow both R/W (READ,INSERTS). models.OconusRateArea is (R)READONLY. ***
type testOnlyOconusRateArea struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	RateAreaId         uuid.UUID `json:"rate_area_id" db:"rate_area_id"`
	CountryId          uuid.UUID `json:"country_id" db:"country_id"`
	UsPostRegionCityId uuid.UUID `json:"us_post_region_cities_id" db:"us_post_region_cities_id"`
	Active             bool      `json:"active" db:"active"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

func (o testOnlyOconusRateArea) TableName() string {
	return "re_oconus_rate_areas"
}
