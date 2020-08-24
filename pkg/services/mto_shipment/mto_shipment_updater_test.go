package mtoshipment

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/route/mocks"

	"github.com/transcom/mymove/pkg/services"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/fetch"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MTOShipmentServiceSuite) TestMTOShipmentUpdater() {
	oldMTOShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
	builder := query.NewQueryBuilder(suite.DB())
	fetcher := fetch.NewFetcher(builder)
	planner := &mocks.Planner{}
	planner.On("TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(500, nil)
	mtoShipmentUpdater := NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)

	requestedPickupDate := *oldMTOShipment.RequestedPickupDate
	scheduledPickupDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	firstAvailableDeliveryDate := time.Date(2019, time.March, 10, 0, 0, 0, 0, time.UTC)
	actualPickupDate := time.Date(2020, time.June, 8, 0, 0, 0, 0, time.UTC)
	secondaryPickupAddress := testdatagen.MakeAddress3(suite.DB(), testdatagen.Assertions{})
	secondaryDeliveryAddress := testdatagen.MakeAddress4(suite.DB(), testdatagen.Assertions{})
	primeActualWeight := unit.Pound(1234)
	primeEstimatedWeight := unit.Pound(1234)

	mtoShipment := models.MTOShipment{
		ID:                         oldMTOShipment.ID,
		MoveTaskOrderID:            oldMTOShipment.MoveTaskOrderID,
		DestinationAddress:         oldMTOShipment.DestinationAddress,
		DestinationAddressID:       oldMTOShipment.DestinationAddressID,
		PickupAddress:              oldMTOShipment.PickupAddress,
		PickupAddressID:            oldMTOShipment.PickupAddressID,
		RequestedPickupDate:        &requestedPickupDate,
		ScheduledPickupDate:        &scheduledPickupDate,
		ShipmentType:               "INTERNATIONAL_UB",
		PrimeActualWeight:          &primeActualWeight,
		PrimeEstimatedWeight:       &primeEstimatedWeight,
		FirstAvailableDeliveryDate: &firstAvailableDeliveryDate,
		Status:                     oldMTOShipment.Status,
		ActualPickupDate:           &actualPickupDate,
		ApprovedDate:               &firstAvailableDeliveryDate,
	}

	ghcDomesticTransitTime := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     0,
		WeightLbsUpper:     10000,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

	//now := time.Now()
	primeEstimatedWeight = unit.Pound(4500)

	suite.T().Run("Etag is stale", func(t *testing.T) {
		eTag := etag.GenerateEtag(time.Now())
		_, err := mtoShipmentUpdater.UpdateMTOShipment(&mtoShipment, eTag)
		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	suite.T().Run("If-Unmodified-Since is equal to the updated_at date", func(t *testing.T) {
		eTag := etag.GenerateEtag(oldMTOShipment.UpdatedAt)
		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&mtoShipment, eTag)
		suite.NoError(err)

		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(updatedMTOShipment.MoveTaskOrder.ID, oldMTOShipment.MoveTaskOrder.ID)
		suite.Equal(updatedMTOShipment.ShipmentType, models.MTOShipmentTypeInternationalUB)

		suite.NotZero(updatedMTOShipment.PickupAddressID, oldMTOShipment.PickupAddressID)

		suite.NotZero(updatedMTOShipment.SecondaryPickupAddressID, secondaryPickupAddress.ID)
		suite.NotZero(updatedMTOShipment.SecondaryDeliveryAddressID, secondaryDeliveryAddress.ID)
		suite.Equal(updatedMTOShipment.PrimeActualWeight, &primeActualWeight)
		suite.True(actualPickupDate.Equal(*updatedMTOShipment.ActualPickupDate))
		suite.True(firstAvailableDeliveryDate.Equal(*updatedMTOShipment.FirstAvailableDeliveryDate))
	})

	oldMTOShipment2 := testdatagen.MakeDefaultMTOShipment(suite.DB())
	mtoShipment2 := models.MTOShipment{
		ID:           oldMTOShipment2.ID,
		ShipmentType: "INTERNATIONAL_UB",
	}

	suite.T().Run("Updater can handle optional queries set as nil", func(t *testing.T) {
		eTag := etag.GenerateEtag(oldMTOShipment2.UpdatedAt)

		updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&mtoShipment2, eTag)
		suite.NoError(err)

		suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
		suite.Equal(updatedMTOShipment.MoveTaskOrder.ID, oldMTOShipment2.MoveTaskOrder.ID)
		suite.Equal(updatedMTOShipment.ShipmentType, models.MTOShipmentTypeInternationalUB)
		suite.Nil(updatedMTOShipment.PrimeEstimatedWeight)
	})

	suite.T().Run("Successful update to a minimal MTO shipment", func(t *testing.T) {

		oldShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusDraft,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		//newDestinationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		//	Address: models.Address{
		//		StreetAddress1: "987 Other Avenue",
		//		StreetAddress2: swag.String("P.O. Box 1234"),
		//		StreetAddress3: swag.String("c/o Another Person"),
		//		City:           "Des Moines",
		//		State:          "IA",
		//		PostalCode:     "50309",
		//		Country:        swag.String("US"),
		//	},
		//})
		//newPickupAddress := testdatagen.MakeAddress4(suite.DB(), testdatagen.Assertions{
		//	Address : models.Address{
		//		StreetAddress1: "987 Over There Avenue",
		//		StreetAddress2: swag.String("P.O. Box 1234"),
		//		StreetAddress3: swag.String("c/o Another Person"),
		//		City:           "Houston",
		//		State:          "TX",
		//		PostalCode:     "77083",
		//		Country:        swag.String("US"),
		//	},
		//})
		requestedPickupDate := time.Date(2019, time.March, 15, 0, 0, 0, 0, time.UTC)
		requestedDeliveryDate := time.Date(2019, time.March, 30, 0, 0, 0, 0, time.UTC)
		customerRemarks := "I have a grandfather clock"
		updatedShipment := models.MTOShipment{
			ID: oldShipment.ID,
			//DestinationAddress:         &newDestinationAddress,
			//DestinationAddressID:       &newDestinationAddress.ID,
			//PickupAddress:              &newPickupAddress,
			//PickupAddressID:            &newPickupAddress.ID,
			//SecondaryPickupAddress:     &secondaryPickupAddress,
			//SecondaryDeliveryAddress:   &secondaryDeliveryAddress,
			RequestedPickupDate:        &requestedPickupDate,
			ScheduledPickupDate:        &scheduledPickupDate,
			RequestedDeliveryDate:      &requestedDeliveryDate,
			ActualPickupDate:           &actualPickupDate,
			ApprovedDate:               &firstAvailableDeliveryDate,
			PrimeActualWeight:          &primeActualWeight,
			PrimeEstimatedWeight:       &primeEstimatedWeight,
			FirstAvailableDeliveryDate: &firstAvailableDeliveryDate,
			Status:                     models.MTOShipmentStatusSubmitted,
			CustomerRemarks:            &customerRemarks,
			//MTOAgents:
			// PrimeEstimatedWeightRecordedDate
			//UpdatedAt: oldShipment.UpdatedAt,
		}
		newShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
		suite.NoError(err)
		suite.True(requestedPickupDate.Equal(*newShipment.RequestedPickupDate))
		suite.True(scheduledPickupDate.Equal(*newShipment.ScheduledPickupDate))
		suite.True(requestedDeliveryDate.Equal(*newShipment.RequestedDeliveryDate))
		suite.True(actualPickupDate.Equal(*newShipment.ActualPickupDate))
		suite.True(firstAvailableDeliveryDate.Equal(*newShipment.ApprovedDate))
		suite.True(firstAvailableDeliveryDate.Equal(*newShipment.FirstAvailableDeliveryDate))
		suite.Equal(primeEstimatedWeight, *newShipment.PrimeEstimatedWeight)
		suite.Equal(primeActualWeight, *newShipment.PrimeActualWeight)
		suite.Equal(customerRemarks, *newShipment.CustomerRemarks)
		suite.Equal(models.MTOShipmentStatusSubmitted, newShipment.Status)
	})

	//suite.T().Run("Failed case if not both approved date and estimated weight recorded date is more than ten days prior to scheduled move date", func(t *testing.T) {
	//	eightDaysFromNow := now.AddDate(0, 0, 8)
	//	threeDaysBefore := now.AddDate(0, 0, -3)
	//	oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
	//		MTOShipment: models.MTOShipment{
	//			Status:              "APPROVED",
	//			ScheduledPickupDate: &eightDaysFromNow,
	//			ApprovedDate:        &threeDaysBefore,
	//		},
	//	})
	//	eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
	//	updatedShipment := models.MTOShipment{
	//		ID:                   oldShipment.ID,
	//		PrimeEstimatedWeight: &primeEstimatedWeight,
	//	}
	//
	//	_, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
	//	suite.Error(err)
	//})
	//
	//suite.T().Run("Successful case if both approved date and estimated weight recorded date is more than ten days prior to scheduled move date", func(t *testing.T) {
	//	tenDaysFromNow := now.AddDate(0, 0, 11)
	//	oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
	//		MTOShipment: models.MTOShipment{
	//			Status:              "APPROVED",
	//			ScheduledPickupDate: &tenDaysFromNow,
	//			ApprovedDate:        &now,
	//		},
	//	})
	//	eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
	//	updatedShipment := models.MTOShipment{
	//		ID:                   oldShipment.ID,
	//		PrimeEstimatedWeight: &primeEstimatedWeight,
	//	}
	//	updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
	//	suite.NoError(err)
	//
	//	suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
	//	suite.NotNil(updatedMTOShipment.PrimeEstimatedWeightRecordedDate)
	//})
	//
	//suite.T().Run("Successful case if scheduled pickup is changed. RequiredDeliveryDate should be generated.", func(t *testing.T) {
	//	tenDaysFromNow := now.AddDate(0, 0, 11)
	//	oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
	//		MTOShipment: models.MTOShipment{
	//			Status:       "APPROVED",
	//			ApprovedDate: &now,
	//		},
	//	})
	//	eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
	//	updatedShipment := models.MTOShipment{
	//		ID:                   oldShipment.ID,
	//		PrimeEstimatedWeight: &primeEstimatedWeight,
	//		ScheduledPickupDate:  &tenDaysFromNow,
	//	}
	//	updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
	//	suite.NoError(err)
	//	suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
	//	suite.NotNil(updatedMTOShipment.RequiredDeliveryDate)
	//
	//	// Let's double check our maths.
	//	expectedRDD := updatedShipment.ScheduledPickupDate.AddDate(0, 0, 12)
	//	actualRDD := *updatedMTOShipment.RequiredDeliveryDate
	//	suite.Equal(expectedRDD.Year(), actualRDD.Year())
	//	suite.Equal(expectedRDD.Month(), actualRDD.Month())
	//	suite.Equal(expectedRDD.Day(), actualRDD.Day())
	//
	//})
	//
	//suite.T().Run("Successful case if in Alaska, should add an extra 10 days to required delivery date", func(t *testing.T) {
	//	tenDaysFromNow := now.AddDate(0, 0, 11)
	//	akAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
	//		Address: models.Address{
	//			PostalCode: "12345",
	//			State:      "AK",
	//		},
	//	})
	//	oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
	//		MTOShipment: models.MTOShipment{
	//			Status:               "APPROVED",
	//			ApprovedDate:         &now,
	//			DestinationAddress:   &akAddress,
	//			DestinationAddressID: &akAddress.ID,
	//		},
	//	})
	//	eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
	//	updatedShipment := models.MTOShipment{
	//		ID:                   oldShipment.ID,
	//		PrimeEstimatedWeight: &primeEstimatedWeight,
	//		ScheduledPickupDate:  &tenDaysFromNow,
	//		DestinationAddress:   &akAddress,
	//		DestinationAddressID: &akAddress.ID,
	//	}
	//	updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
	//	suite.NoError(err)
	//	suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
	//	suite.NotNil(updatedMTOShipment.RequiredDeliveryDate)
	//
	//	// Let's double check our maths.
	//	expectedRDD := updatedShipment.ScheduledPickupDate.AddDate(0, 0, 22)
	//	actualRDD := *updatedMTOShipment.RequiredDeliveryDate
	//	suite.Equal(expectedRDD.Year(), actualRDD.Year())
	//	suite.Equal(expectedRDD.Month(), actualRDD.Month())
	//	suite.Equal(expectedRDD.Day(), actualRDD.Day())
	//
	//})
	//
	//suite.T().Run("Successful case in Adak, Alaska, should add 20 days to required delivery date", func(t *testing.T) {
	//	tenDaysFromNow := now.AddDate(0, 0, 11)
	//	adakAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
	//		Address: models.Address{
	//			PostalCode: "99546",
	//			State:      "AK",
	//			City:       "Adak",
	//		},
	//	})
	//	oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
	//		MTOShipment: models.MTOShipment{
	//			Status:               "APPROVED",
	//			ApprovedDate:         &now,
	//			DestinationAddress:   &adakAddress,
	//			DestinationAddressID: &adakAddress.ID,
	//		},
	//	})
	//	eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
	//	updatedShipment := models.MTOShipment{
	//		ID:                   oldShipment.ID,
	//		PrimeEstimatedWeight: &primeEstimatedWeight,
	//		ScheduledPickupDate:  &tenDaysFromNow,
	//		DestinationAddress:   &adakAddress,
	//		DestinationAddressID: &adakAddress.ID,
	//	}
	//	updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
	//	suite.NoError(err)
	//	suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
	//	suite.NotNil(updatedMTOShipment.RequiredDeliveryDate)
	//
	//	// Let's double check our maths.
	//	expectedRDD := updatedShipment.ScheduledPickupDate.AddDate(0, 0, 32)
	//	actualRDD := *updatedMTOShipment.RequiredDeliveryDate
	//	suite.Equal(expectedRDD.Year(), actualRDD.Year())
	//	suite.Equal(expectedRDD.Month(), actualRDD.Month())
	//	suite.Equal(expectedRDD.Day(), actualRDD.Day())
	//
	//})
	//
	//suite.T().Run("Failed case if approved date is 3-9 days from scheduled move date but estimated weight recorded date isn't at least 3 days prior to scheduled move date", func(t *testing.T) {
	//	twoDaysFromNow := now.AddDate(0, 0, 2)
	//	twoDaysBefore := now.AddDate(0, 0, -2)
	//	oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
	//		MTOShipment: models.MTOShipment{
	//			Status:              "APPROVED",
	//			ScheduledPickupDate: &twoDaysFromNow,
	//			ApprovedDate:        &twoDaysBefore,
	//		},
	//	})
	//	eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
	//	updatedShipment := models.MTOShipment{
	//		ID:                   oldShipment.ID,
	//		PrimeEstimatedWeight: &primeEstimatedWeight,
	//	}
	//
	//	_, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
	//	suite.Error(err)
	//})
	//
	//suite.T().Run("Successful case if approved date is 3-9 days from scheduled move date and estimated weight recorded date is at least 3 days prior to scheduled move date", func(t *testing.T) {
	//	sixDaysFromNow := now.AddDate(0, 0, 6)
	//	twoDaysBefore := now.AddDate(0, 0, -2)
	//	oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
	//		MTOShipment: models.MTOShipment{
	//			Status:              "APPROVED",
	//			ScheduledPickupDate: &sixDaysFromNow,
	//			ApprovedDate:        &twoDaysBefore,
	//		},
	//	})
	//	eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
	//	updatedShipment := models.MTOShipment{
	//		ID:                   oldShipment.ID,
	//		PrimeEstimatedWeight: &primeEstimatedWeight,
	//	}
	//	updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
	//	suite.NoError(err)
	//
	//	suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
	//	suite.NotNil(updatedMTOShipment.PrimeEstimatedWeightRecordedDate)
	//})
	//
	//suite.T().Run("Failed case if approved date is less than 3 days from scheduled move date but estimated weight recorded date isn't at least 1 day prior to scheduled move date", func(t *testing.T) {
	//	oneDayFromNow := now.AddDate(0, 0, 1)
	//	oneDayBefore := now.AddDate(0, 0, -1)
	//	oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
	//		MTOShipment: models.MTOShipment{
	//			Status:              "APPROVED",
	//			ScheduledPickupDate: &oneDayFromNow,
	//			ApprovedDate:        &oneDayBefore,
	//		},
	//	})
	//	eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
	//	updatedShipment := models.MTOShipment{
	//		ID:                   oldShipment.ID,
	//		PrimeEstimatedWeight: &primeEstimatedWeight,
	//	}
	//
	//	_, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
	//	suite.Error(err)
	//})
	//
	//suite.T().Run("Successful case if approved date is less than 3 days from scheduled move date and estimated weight recorded date is at least 1 day prior to scheduled move date", func(t *testing.T) {
	//	twoDaysFromNow := now.AddDate(0, 0, 2)
	//	oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
	//		MTOShipment: models.MTOShipment{
	//			Status:              "APPROVED",
	//			ScheduledPickupDate: &twoDaysFromNow,
	//			ApprovedDate:        &now,
	//		},
	//	})
	//	eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
	//	updatedShipment := models.MTOShipment{
	//		ID:                   oldShipment.ID,
	//		PrimeEstimatedWeight: &primeEstimatedWeight,
	//	}
	//	updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
	//	suite.NoError(err)
	//
	//	suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
	//	suite.NotNil(updatedMTOShipment.PrimeEstimatedWeightRecordedDate)
	//})
	//
	//suite.T().Run("Successfully update MTO Agents", func(t *testing.T) {
	//	shipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
	//	mtoAgent1 := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
	//		MTOAgent: models.MTOAgent{
	//			MTOShipment:   shipment,
	//			MTOShipmentID: shipment.ID,
	//			FirstName:     swag.String("Test"),
	//			LastName:      swag.String("Agent"),
	//			Email:         swag.String("test@test.email.com"),
	//			MTOAgentType:  models.MTOAgentReleasing,
	//		},
	//	})
	//	mtoAgent2 := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
	//		MTOAgent: models.MTOAgent{
	//			MTOShipment:   shipment,
	//			MTOShipmentID: shipment.ID,
	//			FirstName:     swag.String("Test"),
	//			LastName:      swag.String("Agent2"),
	//			Email:         swag.String("test2@test.email.com"),
	//			MTOAgentType:  models.MTOAgentReceiving,
	//		},
	//	})
	//	eTag := etag.GenerateEtag(shipment.UpdatedAt)
	//
	//	updatedAgents := make(models.MTOAgents, 2)
	//	updatedAgents[0] = mtoAgent1
	//	updatedAgents[1] = mtoAgent2
	//	newFirstName := "hey this is new"
	//	newLastName := "new thing"
	//	updatedAgents[0].FirstName = &newFirstName
	//	updatedAgents[1].LastName = &newLastName
	//
	//	updatedShipment := models.MTOShipment{
	//		ID:        shipment.ID,
	//		MTOAgents: updatedAgents,
	//	}
	//
	//	updatedMTOShipment, err := mtoShipmentUpdater.UpdateMTOShipment(&updatedShipment, eTag)
	//
	//	suite.NoError(err)
	//	suite.NotZero(updatedMTOShipment.ID, oldMTOShipment.ID)
	//	suite.Equal(*updatedMTOShipment.MTOAgents[0].FirstName, newFirstName)
	//	suite.Equal(*updatedMTOShipment.MTOAgents[1].LastName, newLastName)
	//})
}

func (suite *MTOShipmentServiceSuite) TestUpdateMTOShipmentStatus() {
	mto := testdatagen.MakeDefaultMove(suite.DB())
	estimatedWeight := unit.Pound(2000)
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
			ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
			PrimeEstimatedWeight: &estimatedWeight,
			Status:               models.MTOShipmentStatusSubmitted,
		},
	})
	draftShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusDraft,
		},
	})
	shipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})
	shipment3 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})
	shipment4 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})
	approvedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusApproved,
		},
	})
	rejectedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusRejected,
		},
	})
	shipment.Status = models.MTOShipmentStatusSubmitted
	eTag := etag.GenerateEtag(shipment.UpdatedAt)
	status := models.MTOShipmentStatusApproved
	//Need some values for reServices
	reServiceCodes := []models.ReServiceCode{
		models.ReServiceCodeDSH,
		models.ReServiceCodeDLH,
		models.ReServiceCodeFSC,
		models.ReServiceCodeDOP,
		models.ReServiceCodeDDP,
		models.ReServiceCodeDPK,
		models.ReServiceCodeDUPK,
		models.ReServiceCodeDNPKF,
		models.ReServiceCodeDMHF,
		models.ReServiceCodeDBHF,
		models.ReServiceCodeDBTF,
	}

	for _, serviceCode := range reServiceCodes {
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code:      serviceCode,
				Name:      "test",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		})
	}

	ghcDomesticTransitTime := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     0,
		WeightLbsUpper:     10000,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

	// Let's also create a transit time object with a zero upper bound for weight (this can happen in the table).
	ghcDomesticTransitTime0LbsUpper := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     10001,
		WeightLbsUpper:     0,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime0LbsUpper)

	builder := query.NewQueryBuilder(suite.DB())
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(builder)
	planner := &mocks.Planner{}
	planner.On("TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(500, nil)
	updater := NewMTOShipmentStatusUpdater(suite.DB(), builder, siCreator, planner)

	suite.T().Run("If we get a mto shipment pointer with a status it should update and return no error", func(t *testing.T) {
		_, err := updater.UpdateMTOShipmentStatus(shipment.ID, status, nil, eTag)
		suite.NoError(err)
		serviceItems := models.MTOServiceItems{}
		_ = suite.DB().All(&serviceItems)
		fetchedShipment := models.MTOShipment{}
		err = suite.DB().Find(&fetchedShipment, shipment.ID)
		suite.NoError(err)
		// We also should have a required delivery date
		suite.NotNil(fetchedShipment.RequiredDeliveryDate)
	})

	suite.T().Run("If we are approving a shipment but are missing key information (estimated weight and pickup date) it should fail", func(t *testing.T) {
		_, err := updater.UpdateMTOShipmentStatus(shipment2.ID, status, nil, eTag)
		suite.NotNil(err)
	})

	suite.T().Run("If we act on a shipment with a weight that has a 0 upper weight it should still work", func(t *testing.T) {
		estimatedWeight := unit.Pound(11000)
		shipmentHeavy := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
				ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
				PrimeEstimatedWeight: &estimatedWeight,
				Status:               models.MTOShipmentStatusSubmitted,
			},
		})
		shipmentHeavyEtag := etag.GenerateEtag(shipmentHeavy.UpdatedAt)
		_, err := updater.UpdateMTOShipmentStatus(shipmentHeavy.ID, status, nil, shipmentHeavyEtag)
		suite.NoError(err)
		serviceItems := models.MTOServiceItems{}
		_ = suite.DB().All(&serviceItems)
		fetchedShipment := models.MTOShipment{}
		err = suite.DB().Find(&fetchedShipment, shipment.ID)
		suite.NoError(err)
		// We also should have a required delivery date
		suite.NotNil(fetchedShipment.RequiredDeliveryDate)
	})

	suite.T().Run("Update MTO Shipment status DRAFT to SUBMITTED should return no error", func(t *testing.T) {
		eTag = etag.GenerateEtag(draftShipment.UpdatedAt)
		_, err := updater.UpdateMTOShipmentStatus(draftShipment.ID, "SUBMITTED", nil, eTag)
		suite.NoError(err)
	})

	suite.T().Run("Update MTO Shipment SUBMITTED status to REJECTED with a rejection reason should return no error", func(t *testing.T) {
		eTag = etag.GenerateEtag(shipment2.UpdatedAt)
		rejectionReason := "Rejection reason"
		returnedShipment, err := updater.UpdateMTOShipmentStatus(shipment2.ID, "REJECTED", &rejectionReason, eTag)
		suite.NoError(err)
		suite.NotNil(returnedShipment)
		suite.Equal(models.MTOShipmentStatusRejected, returnedShipment.Status)
		suite.Equal(&rejectionReason, returnedShipment.RejectionReason)
		// Since this is a rejection we should not generate a required delivery date
		suite.Nil(returnedShipment.RequiredDeliveryDate)
	})

	suite.T().Run("Update MTO Shipment status to REJECTED with no rejection reason should return error", func(t *testing.T) {
		eTag = etag.GenerateEtag(shipment3.UpdatedAt)
		_, err := updater.UpdateMTOShipmentStatus(shipment3.ID, "REJECTED", nil, eTag)
		suite.Error(err)
		fmt.Printf("%#v", err)
		suite.IsType(services.InvalidInputError{}, err)
	})

	suite.T().Run("Update MTO Shipment in APPROVED status should return error", func(t *testing.T) {
		rejectionReason := "Rejection reason"
		_, err := updater.UpdateMTOShipmentStatus(approvedShipment.ID, "REJECTED", &rejectionReason, eTag)
		suite.Error(err)
	})

	suite.T().Run("Update MTO Shipment in REJECTED status should return error", func(t *testing.T) {
		_, err := updater.UpdateMTOShipmentStatus(rejectedShipment.ID, "APPROVED", nil, eTag)
		suite.Error(err)
	})

	suite.T().Run("Passing in a stale identifier", func(t *testing.T) {
		staleETag := etag.GenerateEtag(time.Now())

		_, err := updater.UpdateMTOShipmentStatus(shipment4.ID, "APPROVED", nil, staleETag)
		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	suite.T().Run("Passing in an invalid status", func(t *testing.T) {
		eTag = etag.GenerateEtag(shipment4.UpdatedAt)

		_, err := updater.UpdateMTOShipmentStatus(shipment4.ID, "invalid", nil, eTag)
		suite.Error(err)
		fmt.Printf("%#v", err)
		suite.IsType(services.InvalidInputError{}, err)
	})

	suite.T().Run("Passing in a bad shipment id", func(t *testing.T) {
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")

		_, err := updater.UpdateMTOShipmentStatus(badShipmentID, "APPROVED", nil, eTag)
		suite.Error(err)
		fmt.Printf("%#v", err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Changing to APPROVED status records approved_date", func(t *testing.T) {
		shipment5 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		})
		eTag = etag.GenerateEtag(shipment5.UpdatedAt)

		suite.Nil(shipment5.ApprovedDate)
		_, err := updater.UpdateMTOShipmentStatus(shipment5.ID, models.MTOShipmentStatusApproved, nil, eTag)
		suite.NoError(err)
		suite.DB().Find(&shipment5, shipment5.ID)
		suite.Equal(models.MTOShipmentStatusApproved, shipment5.Status)
		suite.NotNil(shipment5.ApprovedDate)
	})

	suite.T().Run("Changing to a non-APPROVED status does not record approved_date", func(t *testing.T) {
		shipment6 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		})
		eTag = etag.GenerateEtag(shipment6.UpdatedAt)
		rejectionReason := "reason"

		suite.Nil(shipment6.ApprovedDate)
		_, err := updater.UpdateMTOShipmentStatus(shipment6.ID, models.MTOShipmentStatusRejected, &rejectionReason, eTag)
		suite.NoError(err)
		suite.DB().Find(&shipment6, shipment6.ID)
		suite.Equal(models.MTOShipmentStatusRejected, shipment6.Status)
		suite.Nil(shipment3.ApprovedDate)
	})
}
