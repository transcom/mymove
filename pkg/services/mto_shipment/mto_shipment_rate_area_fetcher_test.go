package mtoshipment

import (
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

const fairbanksAlaskaPostalCode = "99716"
const anchorageAlaskaPostalCode = "99521"
const wasillaAlaskaPostalCode = "99652"
const beverlyHillsCAPostalCode = "90210"
const sanDiegoCAPostalCode = "92075"
const brooklynNYPostalCode = "11220"
const rateAreaAK1 = "US8190100"
const rateAreaAK2 = "US8101000"
const rateAreaCA = "US88"

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
					MarketCode: models.MarketCodeInternational,
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
					MarketCode: models.MarketCodeDomestic,
				},
				models.MTOShipment{
					PPMShipment: &models.PPMShipment{
						PickupAddress: &models.Address{
							StreetAddress1: "123 Main St",
							City:           "Beverly Hills",
							State:          "CA",
							PostalCode:     beverlyHillsCAPostalCode,
							IsOconus:       models.BoolPointer(false),
						},
						DestinationAddress: &models.Address{
							StreetAddress1: "123 Main St",
							City:           "Wasilla",
							State:          "AK",
							PostalCode:     wasillaAlaskaPostalCode,
							IsOconus:       models.BoolPointer(true),
						},
					},
					MarketCode: models.MarketCodeInternational,
				},
			},
		}

		isRateAreaEquals := func(expectedRateArea string, postalCode string, shipmentPostalCodeRateArea *[]services.ShipmentPostalCodeRateArea) bool {
			var shipmentPostalCodeRateAreaLookupMap = make(map[string]services.ShipmentPostalCodeRateArea)
			for _, i := range *shipmentPostalCodeRateArea {
				shipmentPostalCodeRateAreaLookupMap[i.PostalCode] = i
			}
			if _, ok := shipmentPostalCodeRateAreaLookupMap[postalCode]; !ok {
				return false
			}
			return (shipmentPostalCodeRateAreaLookupMap[postalCode].RateArea.Code == expectedRateArea)
		}

		shipmentPostalCodeRateAreas, err := shipmentRateAreaFetcher.GetPrimeMoveShipmentRateAreas(suite.AppContextForTest(), testMove)
		suite.NotNil(shipmentPostalCodeRateAreas)
		suite.FatalNoError(err)
		suite.Equal(4, len(*shipmentPostalCodeRateAreas))

		suite.Equal(true, isRateAreaEquals(rateAreaAK1, fairbanksAlaskaPostalCode, shipmentPostalCodeRateAreas))
		suite.Equal(false, isRateAreaEquals(rateAreaAK1, anchorageAlaskaPostalCode, shipmentPostalCodeRateAreas))
		suite.Equal(true, isRateAreaEquals(rateAreaAK2, wasillaAlaskaPostalCode, shipmentPostalCodeRateAreas))
		suite.Equal(true, isRateAreaEquals(rateAreaCA, beverlyHillsCAPostalCode, shipmentPostalCodeRateAreas))

		// Postal code used only in a CONUS shipment should not have been fetched
		i := slices.IndexFunc(*shipmentPostalCodeRateAreas, func(pcra services.ShipmentPostalCodeRateArea) bool {
			return pcra.PostalCode == sanDiegoCAPostalCode
		})
		suite.Equal(-1, i)

		suite.Equal(false, isRateAreaEquals(rateAreaAK1, beverlyHillsCAPostalCode, shipmentPostalCodeRateAreas))
		suite.Equal(false, isRateAreaEquals(rateAreaAK1, wasillaAlaskaPostalCode, shipmentPostalCodeRateAreas))
		suite.Equal(false, isRateAreaEquals(rateAreaAK2, fairbanksAlaskaPostalCode, shipmentPostalCodeRateAreas))
		suite.Equal(true, isRateAreaEquals(rateAreaAK2, anchorageAlaskaPostalCode, shipmentPostalCodeRateAreas))
		suite.Equal(false, isRateAreaEquals(rateAreaAK2, beverlyHillsCAPostalCode, shipmentPostalCodeRateAreas))
		suite.Equal(false, isRateAreaEquals(rateAreaCA, fairbanksAlaskaPostalCode, shipmentPostalCodeRateAreas))
		suite.Equal(false, isRateAreaEquals(rateAreaCA, anchorageAlaskaPostalCode, shipmentPostalCodeRateAreas))
		suite.Equal(false, isRateAreaEquals(rateAreaCA, wasillaAlaskaPostalCode, shipmentPostalCodeRateAreas))
	})

	suite.Run("Does not return rate areas for CONUS only shipments", func() {
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
					MarketCode: models.MarketCodeDomestic,
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
					MarketCode: models.MarketCodeDomestic,
				},
			},
		}

		domServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ServiceArea:      "004",
				ServicesSchedule: 2,
			},
			ReContract: testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{}),
		})
		suite.NotNil(domServiceArea)
		suite.NotNil(domServiceArea.Contract)

		// setup contract year within availableToPrimeAtTime time
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: testdatagen.ContractStartDate,
				EndDate:   testdatagen.ContractEndDate,
			},
		})

		setupDomesticRateAreaAndZip3s := func(rateAreaCode string, rateAreaName string, postalCodes map[string]string, domesticServiceArea models.ReDomesticServiceArea) (models.ReRateArea, error) {
			rateArea := testdatagen.FetchOrMakeReRateArea(suite.DB(), testdatagen.Assertions{
				ReRateArea: models.ReRateArea{
					ContractID: domesticServiceArea.ContractID,
					IsOconus:   false,
					Code:       rateAreaCode,
					Name:       rateAreaName,
					Contract:   domesticServiceArea.Contract,
				},
			})

			for postalCode, basePointCity := range postalCodes {
				testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
					ReZip3: models.ReZip3{
						Contract:              domesticServiceArea.Contract,
						ContractID:            domesticServiceArea.ContractID,
						DomesticServiceAreaID: domesticServiceArea.ID,
						DomesticServiceArea:   domesticServiceArea,
						Zip3:                  postalCode[0:3],
						BasePointCity:         basePointCity,
					},
				})
			}

			return rateArea, nil
		}

		_, err := setupDomesticRateAreaAndZip3s("US88", "California-South", map[string]string{beverlyHillsCAPostalCode: "Beverly Hills", sanDiegoCAPostalCode: "San Diego"}, domServiceArea)
		if err != nil {
			suite.Fail(err.Error())
		}

		_, err = setupDomesticRateAreaAndZip3s("US17", "New York", map[string]string{brooklynNYPostalCode: "Brooklyn"}, domServiceArea)
		if err != nil {
			suite.Fail(err.Error())
		}

		shipmentPostalCodeRateAreas, err := shipmentRateAreaFetcher.GetPrimeMoveShipmentRateAreas(suite.AppContextForTest(), testMove)
		suite.NotNil(shipmentPostalCodeRateAreas)
		suite.Equal(0, len(*shipmentPostalCodeRateAreas))
		suite.Nil(err)
	})

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
