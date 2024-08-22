package boatshipment

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *BoatShipmentSuite) TestBoatShipmentFetcher() {

	fetcher := NewBoatShipmentFetcher()

	suite.Run("GetBoatShipment", func() {
		suite.Run("Can fetch a Boat Shipment if there is no session (e.g. a prime request)", func() {
			appCtx := suite.AppContextWithSessionForTest(nil)

			boatShipment := factory.BuildBoatShipment(suite.DB(), nil, nil)

			boatShipmentReturned, err := fetcher.GetBoatShipment(
				appCtx,
				boatShipment.ID,
				nil,
				nil,
			)

			if suite.NoError(err) && suite.NotNil(boatShipmentReturned) {
				suite.Equal(boatShipment.ID, boatShipmentReturned.ID)
			}
		})

		suite.Run("Can fetch a Boat Shipment if it is an office user making a request from the office app", func() {
			officeUser := factory.BuildOfficeUser(suite.DB(), factory.GetTraitActiveOfficeUser(), nil)

			appCtx := suite.AppContextWithSessionForTest(&auth.Session{
				ApplicationName: auth.OfficeApp,
				UserID:          officeUser.User.ID,
				OfficeUserID:    officeUser.ID,
			})

			boatShipment := factory.BuildBoatShipment(suite.DB(), nil, nil)

			boatShipmentReturned, err := fetcher.GetBoatShipment(
				appCtx,
				boatShipment.ID,
				nil,
				nil,
			)

			if suite.NoError(err) && suite.NotNil(boatShipmentReturned) {
				suite.Equal(boatShipment.ID, boatShipmentReturned.ID)
			}
		})

		suite.Run("Can fetch a Boat Shipment if it is a customer app request by the customer it belongs to", func() {
			boatShipment := factory.BuildBoatShipment(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
			serviceMember := boatShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

			appCtx := suite.AppContextWithSessionForTest(&auth.Session{
				ApplicationName: auth.MilApp,
				UserID:          serviceMember.User.ID,
				ServiceMemberID: serviceMember.ID,
			})

			boatShipmentReturned, err := fetcher.GetBoatShipment(
				appCtx,
				boatShipment.ID,
				nil,
				nil,
			)

			if suite.NoError(err) && suite.NotNil(boatShipmentReturned) {
				suite.Equal(boatShipment.ID, boatShipmentReturned.ID)
			}
		})

		suite.Run("Returns a not found error if it is a customer app request by a customer that it doesn't belong to", func() {
			maliciousUser := factory.BuildExtendedServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

			appCtx := suite.AppContextWithSessionForTest(&auth.Session{
				ApplicationName: auth.MilApp,
				UserID:          maliciousUser.User.ID,
				ServiceMemberID: maliciousUser.ID,
			})

			boatShipment := factory.BuildBoatShipment(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

			boatShipmentReturned, err := fetcher.GetBoatShipment(
				appCtx,
				boatShipment.ID,
				nil,
				nil,
			)

			if suite.Error(err) && suite.Nil(boatShipmentReturned) {
				suite.IsType(apperror.NotFoundError{}, err)

				suite.Equal(fmt.Sprintf("ID: %s not found while looking for BoatShipment", boatShipment.ID), err.Error())
			}
		})

		suite.Run("Returns a not found error if the Boat Shipment does not exist", func() {
			nonexistentID := uuid.Must(uuid.NewV4())

			boatShipmentReturned, err := fetcher.GetBoatShipment(
				suite.AppContextForTest(),
				nonexistentID,
				nil,
				nil,
			)

			if suite.Error(err) && suite.Nil(boatShipmentReturned) {
				suite.IsType(apperror.NotFoundError{}, err)

				suite.Equal(fmt.Sprintf("ID: %s not found while looking for BoatShipment", nonexistentID), err.Error())
			}
		})

		suite.Run("Returns an error if an invalid association is requested", func() {
			boatShipment := factory.BuildBoatShipment(suite.DB(), nil, nil)

			invalidAssociation := "invalid"
			boatShipmentReturned, err := fetcher.GetBoatShipment(
				suite.AppContextForTest(),
				boatShipment.ID,
				[]string{invalidAssociation},
				nil,
			)

			if suite.Error(err) && suite.Nil(boatShipmentReturned) {
				suite.IsType(apperror.NotImplementedError{}, err)

				suite.Contains(
					err.Error(),
					fmt.Sprintf("Requested eager preload association %s is not implemented", invalidAssociation),
				)
			}
		})

		suite.Run("Returns an error if the shipment has been deleted", func() {
			boatShipment := factory.BuildBoatShipment(suite.DB(), nil, nil)

			err := utilities.SoftDestroy(suite.DB(), &boatShipment)
			suite.FatalNoError(err)

			boatShipmentReturned, err := fetcher.GetBoatShipment(
				suite.AppContextForTest(),
				boatShipment.ID,
				nil,
				nil,
			)

			if suite.Error(err) && suite.Nil(boatShipmentReturned) {
				suite.IsType(apperror.NotFoundError{}, err)

				suite.Equal(fmt.Sprintf("ID: %s not found while looking for BoatShipment", boatShipment.ID), err.Error())
			}
		})

		suite.Run("Returns an error if an invalid postload association is passed in", func() {
			boatShipment := factory.BuildBoatShipment(suite.DB(), nil, nil)

			invalidAssociation := "invalid"
			boatShipmentReturned, err := fetcher.GetBoatShipment(
				suite.AppContextForTest(),
				boatShipment.ID,
				nil,
				[]string{invalidAssociation},
			)

			if suite.Error(err) && suite.Nil(boatShipmentReturned) {
				suite.IsType(apperror.NotImplementedError{}, err)

				suite.Contains(
					err.Error(),
					fmt.Sprintf("Requested post load association %s is not implemented", invalidAssociation),
				)
			}
		})
	})
}

func (suite *BoatShipmentSuite) TestFetchBoatShipment() {

	suite.Run("FindBoatShipment - returns not found for unknown id", func() {
		badID := uuid.Must(uuid.NewV4())
		_, err := FindBoatShipment(suite.AppContextForTest(), badID)
		suite.Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for BoatShipment", badID), err.Error())
	})

	suite.Run("FindBoatShipment - returns not found for deleted shipment", func() {
		boatShipment := factory.BuildBoatShipment(suite.DB(), nil, nil)

		err := utilities.SoftDestroy(suite.DB(), &boatShipment)
		suite.NoError(err)

		_, err = FindBoatShipment(suite.AppContextForTest(), boatShipment.ID)
		suite.Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for BoatShipment", boatShipment.ID), err.Error())
	})

	suite.Run("FetchBoatShipmentFromMTOShipmentID - finds records", func() {
		boat := factory.BuildBoatShipment(suite.DB(), nil, nil)

		retrievedBoat, _ := FetchBoatShipmentFromMTOShipmentID(suite.AppContextForTest(), boat.ShipmentID)

		suite.Equal(retrievedBoat.ID, boat.ID)
		suite.Equal(retrievedBoat.ShipmentID, boat.ShipmentID)

	})

	suite.Run("FetchBoatShipmentFromMTOShipmentID  - returns not found for unknown id", func() {
		badID := uuid.Must(uuid.NewV4())
		_, err := FetchBoatShipmentFromMTOShipmentID(suite.AppContextForTest(), badID)
		suite.Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for BoatShipment", badID), err.Error())
	})

	suite.Run("FetchBoatShipmentFromMTOShipmentID  - returns not found for deleted shipment", func() {
		boatShipment := factory.BuildBoatShipment(suite.DB(), nil, nil)

		err := utilities.SoftDestroy(suite.DB(), &boatShipment)
		suite.NoError(err)

		_, err = FetchBoatShipmentFromMTOShipmentID(suite.AppContextForTest(), boatShipment.ShipmentID)
		suite.Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for BoatShipment", boatShipment.ShipmentID), err.Error())
	})
}
