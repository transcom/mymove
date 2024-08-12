package mobilehomeshipment

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *MobileHomeShipmentSuite) TestMobileHomeShipmentFetcher() {

	fetcher := NewMobileHomeShipmentFetcher()

	suite.Run("GetMobileHomeShipment", func() {
		suite.Run("Can fetch a Mobile Home Shipment if there is no session (e.g. a prime request)", func() {
			appCtx := suite.AppContextWithSessionForTest(nil)

			mobileHomeShipment := factory.BuildMobileHomeShipment(suite.DB(), nil, nil)

			mobileHomeShipmentReturned, err := fetcher.GetMobileHomeShipment(
				appCtx,
				mobileHomeShipment.ID,
				nil,
				nil,
			)

			if suite.NoError(err) && suite.NotNil(mobileHomeShipmentReturned) {
				suite.Equal(mobileHomeShipment.ID, mobileHomeShipmentReturned.ID)
			}
		})

		suite.Run("Can fetch a Mobile Home Shipment if it is an office user making a request from the office app", func() {
			officeUser := factory.BuildOfficeUser(suite.DB(), factory.GetTraitActiveOfficeUser(), nil)

			appCtx := suite.AppContextWithSessionForTest(&auth.Session{
				ApplicationName: auth.OfficeApp,
				UserID:          officeUser.User.ID,
				OfficeUserID:    officeUser.ID,
			})

			mobileHomeShipment := factory.BuildMobileHomeShipment(suite.DB(), nil, nil)

			mobileHomeShipmentReturned, err := fetcher.GetMobileHomeShipment(
				appCtx,
				mobileHomeShipment.ID,
				nil,
				nil,
			)

			if suite.NoError(err) && suite.NotNil(mobileHomeShipmentReturned) {
				suite.Equal(mobileHomeShipment.ID, mobileHomeShipmentReturned.ID)
			}
		})

		suite.Run("Can fetch a Mobile Home Shipment if it is a customer app request by the customer it belongs to", func() {
			mobileHomeShipment := factory.BuildMobileHomeShipment(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
			serviceMember := mobileHomeShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

			appCtx := suite.AppContextWithSessionForTest(&auth.Session{
				ApplicationName: auth.MilApp,
				UserID:          serviceMember.User.ID,
				ServiceMemberID: serviceMember.ID,
			})

			mobileHomeShipmentReturned, err := fetcher.GetMobileHomeShipment(
				appCtx,
				mobileHomeShipment.ID,
				nil,
				nil,
			)

			if suite.NoError(err) && suite.NotNil(mobileHomeShipmentReturned) {
				suite.Equal(mobileHomeShipment.ID, mobileHomeShipmentReturned.ID)
			}
		})

		suite.Run("Returns a not found error if it is a customer app request by a customer that it doesn't belong to", func() {
			maliciousUser := factory.BuildExtendedServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

			appCtx := suite.AppContextWithSessionForTest(&auth.Session{
				ApplicationName: auth.MilApp,
				UserID:          maliciousUser.User.ID,
				ServiceMemberID: maliciousUser.ID,
			})

			mobileHomeShipment := factory.BuildMobileHomeShipment(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

			mobileHomeShipmentReturned, err := fetcher.GetMobileHomeShipment(
				appCtx,
				mobileHomeShipment.ID,
				nil,
				nil,
			)

			if suite.Error(err) && suite.Nil(mobileHomeShipmentReturned) {
				suite.IsType(apperror.NotFoundError{}, err)

				suite.Equal(fmt.Sprintf("ID: %s not found while looking for MobileHome", mobileHomeShipment.ID), err.Error())
			}
		})

		suite.Run("Returns a not found error if the Mobile Home Shipment does not exist", func() {
			nonexistentID := uuid.Must(uuid.NewV4())

			mobileHomeShipmentReturned, err := fetcher.GetMobileHomeShipment(
				suite.AppContextForTest(),
				nonexistentID,
				nil,
				nil,
			)

			if suite.Error(err) && suite.Nil(mobileHomeShipmentReturned) {
				suite.IsType(apperror.NotFoundError{}, err)

				suite.Equal(fmt.Sprintf("ID: %s not found while looking for MobileHome", nonexistentID), err.Error())
			}
		})

		suite.Run("Returns an error if an invalid association is requested", func() {
			mobileHomeShipment := factory.BuildMobileHomeShipment(suite.DB(), nil, nil)

			invalidAssociation := "invalid"
			mobileHomeShipmentReturned, err := fetcher.GetMobileHomeShipment(
				suite.AppContextForTest(),
				mobileHomeShipment.ID,
				[]string{invalidAssociation},
				nil,
			)

			if suite.Error(err) && suite.Nil(mobileHomeShipmentReturned) {
				suite.IsType(apperror.NotImplementedError{}, err)

				suite.Contains(
					err.Error(),
					fmt.Sprintf("Requested eager preload association %s is not implemented", invalidAssociation),
				)
			}
		})

		suite.Run("Returns an error if the shipment has been deleted", func() {
			mobileHomeShipment := factory.BuildMobileHomeShipment(suite.DB(), nil, nil)

			err := utilities.SoftDestroy(suite.DB(), &mobileHomeShipment)
			suite.FatalNoError(err)

			mobileHomeShipmentReturned, err := fetcher.GetMobileHomeShipment(
				suite.AppContextForTest(),
				mobileHomeShipment.ID,
				nil,
				nil,
			)

			if suite.Error(err) && suite.Nil(mobileHomeShipmentReturned) {
				suite.IsType(apperror.NotFoundError{}, err)

				suite.Equal(fmt.Sprintf("ID: %s not found while looking for MobileHome", mobileHomeShipment.ID), err.Error())
			}
		})

		suite.Run("Returns an error if an invalid postload association is passed in", func() {
			mobileHomeShipment := factory.BuildMobileHomeShipment(suite.DB(), nil, nil)

			invalidAssociation := "invalid"
			mobileHomeShipmentReturned, err := fetcher.GetMobileHomeShipment(
				suite.AppContextForTest(),
				mobileHomeShipment.ID,
				nil,
				[]string{invalidAssociation},
			)

			if suite.Error(err) && suite.Nil(mobileHomeShipmentReturned) {
				suite.IsType(apperror.NotImplementedError{}, err)

				suite.Contains(
					err.Error(),
					fmt.Sprintf("Requested post load association %s is not implemented", invalidAssociation),
				)
			}
		})
	})
}

func (suite *MobileHomeShipmentSuite) TestFetchMobileHomeShipment() {

	suite.Run("FindMobileHomeShipment - returns not found for unknown id", func() {
		badID := uuid.Must(uuid.NewV4())
		_, err := FindMobileHomeShipment(suite.AppContextForTest(), badID)
		suite.Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for MobileHome", badID), err.Error())
	})

	suite.Run("FindMobileHomeShipment - returns not found for deleted shipment", func() {
		mobileHomeShipment := factory.BuildMobileHomeShipment(suite.DB(), nil, nil)

		err := utilities.SoftDestroy(suite.DB(), &mobileHomeShipment)
		suite.NoError(err)

		_, err = FindMobileHomeShipment(suite.AppContextForTest(), mobileHomeShipment.ID)
		suite.Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for MobileHome", mobileHomeShipment.ID), err.Error())
	})

	suite.Run("FetchMobileHomeShipmentFromMTOShipmentID - finds records", func() {
		mobileHome := factory.BuildMobileHomeShipment(suite.DB(), nil, nil)

		retrievedMobileHome, _ := FetchMobileHomeShipmentFromMTOShipmentID(suite.AppContextForTest(), mobileHome.ShipmentID)

		suite.Equal(retrievedMobileHome.ID, mobileHome.ID)
		suite.Equal(retrievedMobileHome.ShipmentID, mobileHome.ShipmentID)

	})

	suite.Run("FetchMobileHomeShipmentFromMTOShipmentID  - returns not found for unknown id", func() {
		badID := uuid.Must(uuid.NewV4())
		_, err := FetchMobileHomeShipmentFromMTOShipmentID(suite.AppContextForTest(), badID)
		suite.Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for MobileHome", badID), err.Error())
	})

	suite.Run("FetchMobileHomeShipmentFromMTOShipmentID  - returns not found for deleted shipment", func() {
		mobileHomeShipment := factory.BuildMobileHomeShipment(suite.DB(), nil, nil)

		err := utilities.SoftDestroy(suite.DB(), &mobileHomeShipment)
		suite.NoError(err)

		_, err = FetchMobileHomeShipmentFromMTOShipmentID(suite.AppContextForTest(), mobileHomeShipment.ShipmentID)
		suite.Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for MobileHome", mobileHomeShipment.ShipmentID), err.Error())
	})
}
