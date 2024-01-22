package notifications

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

var pickupAddressModel = models.Address{
	ID:             uuid.Must(uuid.NewV4()),
	StreetAddress1: "1 First St",
	StreetAddress2: models.StringPointer("Apt 1"),
	City:           "Miami Gardens",
	State:          "FL",
	PostalCode:     "33169",
	Country:        models.StringPointer("US"),
}

var destinationAddressModel = models.Address{
	ID:             uuid.Must(uuid.NewV4()),
	StreetAddress1: "2 Second St",
	StreetAddress2: models.StringPointer("Bldg 2"),
	City:           "Key West",
	State:          "FL",
	PostalCode:     "33040",
	Country:        models.StringPointer("US"),
}

var affiliationDisplayValue = map[models.ServiceMemberAffiliation]string{
	models.AffiliationARMY:       "Army",
	models.AffiliationNAVY:       "Marine Corps, Navy, and Coast Guard",
	models.AffiliationMARINES:    "Marine Corps, Navy, and Coast Guard",
	models.AffiliationAIRFORCE:   "Air Force and Space Force",
	models.AffiliationSPACEFORCE: "Air Force and Space Force",
	models.AffiliationCOASTGUARD: "Marine Corps, Navy, and Coast Guard",
}

var armySubmitLocation = `the Defense Finance and Accounting Service (DFAS)`
var allOtherSubmitLocation = `your local finance office`

func (suite *NotificationSuite) TestPpmPacketEmail() {
	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, nil)
	notification := NewPpmPacketEmail(ppmShipment.ID)

	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{
		ServiceMemberID: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.ID,
		ApplicationName: auth.MilApp,
	}))
	subject := "Your Personally Procured Move (PPM) closeout has been processed and is now available for your review."

	suite.NoError(err)
	suite.Equal(len(emails), 1)

	email := emails[0]
	sm := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
	suite.Equal(email.recipientEmail, *sm.PersonalEmail)
	suite.Equal(email.subject, subject)
	suite.NotEmpty(email.htmlBody)
	suite.NotEmpty(email.textBody)
}

func (suite *NotificationSuite) TestPpmPacketEmailHTMLTemplateRenderForAirAndSpaceForce() {
	var pickupAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: pickupAddressModel,
		},
	}, nil)
	var destinationAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: destinationAddressModel,
		},
	}, nil)

	customAffiliation := models.AffiliationAIRFORCE
	serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
		{Model: models.ServiceMember{
			Affiliation: &customAffiliation,
		}},
	}, nil)
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
	}, nil)
	shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	customPPM := models.PPMShipment{
		ID:                         uuid.Must(uuid.NewV4()),
		ShipmentID:                 shipment.ID,
		Status:                     models.PPMShipmentStatusWaitingOnCustomer,
		PickupPostalAddressID:      &pickupAddress.ID,
		DestinationPostalAddressID: &destinationAddress.ID,
		PickupPostalCode:           "79329",
		DestinationPostalCode:      "90210",
	}

	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, []factory.Customization{
		{Model: customPPM},
	})
	notification := NewPpmPacketEmail(ppmShipment.ID)

	ppmEmailData, _, err := GetEmailData(*notification, suite.AppContextForTest())
	suite.NoError(err)
	suite.NotNil(ppmEmailData)

	suite.EqualExportedValues(ppmEmailData, PpmPacketEmailData{
		OriginCity:       &pickupAddress.City,
		OriginState:      &pickupAddress.State,
		DestinationCity:  &destinationAddress.City,
		DestinationState: &destinationAddress.State,
		SubmitLocation:   allOtherSubmitLocation,
		ServiceBranch:    affiliationDisplayValue[*serviceMember.Affiliation],
		Locator:          move.Locator,
	})

	expectedHTMLContent := `<p>*** DO NOT REPLY directly to this email ***</p>
<p>This is a confirmation that your Personally Procured Move (PPM) with the <strong>assigned move code ` + move.Locator + `</strong> from <strong>` + pickupAddress.City + `, ` + pickupAddress.State + ` to ` + destinationAddress.City + `, ` + destinationAddress.State + ` </strong>has been processed in MilMove. </p>
<h4>Next steps:</h4>

<p>For ` + affiliationDisplayValue[*serviceMember.Affiliation] + ` personnel (FURTHER ACTION REQUIRED):</p>
<p>You can now log into MilMove <a href="https://my.move.mil/">https://my.move.mil/</a> and download your payment packet to submit to ` + allOtherSubmitLocation + `. <strong>You must complete this step to receive final settlement of your PPM.</strong></p>
<p>Note: The Transportation Office does not determine claimable expenses. Claimable expenses will be determined by finance.</p>

<p>If you have any questions, contact a government transportation office. You can see a listing of transportation offices on Military One Source here: <a href="https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL">https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL</a></p>

<p>Thank you,</p>

<p>USTRANSCOM MilMove Team</p>

<p>
  The information contained in this email may contain Privacy Act information and is therefore protected under the
  Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.
</p>
`

	htmlContent, err := notification.RenderHTML(suite.AppContextForTest(), ppmEmailData)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)
}

func (suite *NotificationSuite) TestPpmPacketEmailHTMLTemplateRenderForArmy() {
	var pickupAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: pickupAddressModel,
		},
	}, nil)
	var destinationAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: destinationAddressModel,
		},
	}, nil)

	customAffiliation := models.AffiliationARMY
	serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
		{Model: models.ServiceMember{
			Affiliation: &customAffiliation,
		}},
	}, nil)
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
	}, nil)
	shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	customPPM := models.PPMShipment{
		ID:                         uuid.Must(uuid.NewV4()),
		ShipmentID:                 shipment.ID,
		Status:                     models.PPMShipmentStatusWaitingOnCustomer,
		PickupPostalAddressID:      &pickupAddress.ID,
		DestinationPostalAddressID: &destinationAddress.ID,
		PickupPostalCode:           "79329",
		DestinationPostalCode:      "90210",
	}

	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, []factory.Customization{
		{Model: customPPM},
	})
	notification := NewPpmPacketEmail(ppmShipment.ID)

	ppmEmailData, _, err := GetEmailData(*notification, suite.AppContextForTest())
	suite.NoError(err)
	suite.NotNil(ppmEmailData)

	suite.EqualExportedValues(ppmEmailData, PpmPacketEmailData{
		OriginCity:       &pickupAddress.City,
		OriginState:      &pickupAddress.State,
		DestinationCity:  &destinationAddress.City,
		DestinationState: &destinationAddress.State,
		SubmitLocation:   armySubmitLocation,
		ServiceBranch:    affiliationDisplayValue[*serviceMember.Affiliation],
		Locator:          move.Locator,
	})

	expectedHTMLContent := `<p>*** DO NOT REPLY directly to this email ***</p>
<p>This is a confirmation that your Personally Procured Move (PPM) with the <strong>assigned move code ` + move.Locator + `</strong> from <strong>` + pickupAddress.City + `, ` + pickupAddress.State + ` to ` + destinationAddress.City + `, ` + destinationAddress.State + ` </strong>has been processed in MilMove. </p>
<h4>Next steps:</h4>

<p>For ` + affiliationDisplayValue[*serviceMember.Affiliation] + ` personnel (FURTHER ACTION REQUIRED):</p>
<p>You can now log into MilMove <a href="https://my.move.mil/">https://my.move.mil/</a> and download your payment packet to submit to ` + armySubmitLocation + `. <strong>You must complete this step to receive final settlement of your PPM.</strong></p>
<p>Note: Not all claimed expenses may have been accepted during PPM closeout if they did not meet the definition of a valid expense</p>

<p>If you have any questions, contact a government transportation office. You can see a listing of transportation offices on Military One Source here: <a href="https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL">https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL</a></p>

<p>Thank you,</p>

<p>USTRANSCOM MilMove Team</p>

<p>
  The information contained in this email may contain Privacy Act information and is therefore protected under the
  Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.
</p>
`

	htmlContent, err := notification.RenderHTML(suite.AppContextForTest(), ppmEmailData)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)
}

func (suite *NotificationSuite) TestPpmPacketEmailHTMLTemplateRenderForNavalBranches() {
	var pickupAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: pickupAddressModel,
		},
	}, nil)
	var destinationAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: destinationAddressModel,
		},
	}, nil)

	customAffiliation := models.AffiliationMARINES
	serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
		{Model: models.ServiceMember{
			Affiliation: &customAffiliation,
		}},
	}, nil)
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
	}, nil)
	shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	customPPM := models.PPMShipment{
		ID:                         uuid.Must(uuid.NewV4()),
		ShipmentID:                 shipment.ID,
		Status:                     models.PPMShipmentStatusWaitingOnCustomer,
		PickupPostalAddressID:      &pickupAddress.ID,
		DestinationPostalAddressID: &destinationAddress.ID,
		PickupPostalCode:           "79329",
		DestinationPostalCode:      "90210",
	}

	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, []factory.Customization{
		{Model: customPPM},
	})
	notification := NewPpmPacketEmail(ppmShipment.ID)

	ppmEmailData, _, err := GetEmailData(*notification, suite.AppContextForTest())
	suite.NoError(err)
	suite.NotNil(ppmEmailData)

	suite.EqualExportedValues(ppmEmailData, PpmPacketEmailData{
		OriginCity:       &pickupAddress.City,
		OriginState:      &pickupAddress.State,
		DestinationCity:  &destinationAddress.City,
		DestinationState: &destinationAddress.State,
		SubmitLocation:   allOtherSubmitLocation,
		ServiceBranch:    affiliationDisplayValue[*serviceMember.Affiliation],
		Locator:          move.Locator,
	})

	expectedHTMLContent := `<p>*** DO NOT REPLY directly to this email ***</p>
<p>This is a confirmation that your Personally Procured Move (PPM) with the <strong>assigned move code ` + move.Locator + `</strong> from <strong>` + pickupAddress.City + `, ` + pickupAddress.State + ` to ` + destinationAddress.City + `, ` + destinationAddress.State + ` </strong>has been processed in MilMove. </p>
<h4>Next steps:</h4>

<p>For ` + affiliationDisplayValue[*serviceMember.Affiliation] + ` personnel:</p>
<p>You can now log into MilMove <a href="https://my.move.mil/">https://my.move.mil/</a> and view your payment packet; however, you do not need to forward your packet to finance as your closeout location is associated with your finance office and they will handle this step for you.</p>
<p>Note: Not all claimed expenses may have been accepted during PPM closeout if they did not meet the definition of a valid expense.</p>

<p>If you have any questions, contact a government transportation office. You can see a listing of transportation offices on Military One Source here: <a href="https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL">https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL</a></p>

<p>Thank you,</p>

<p>USTRANSCOM MilMove Team</p>

<p>
  The information contained in this email may contain Privacy Act information and is therefore protected under the
  Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.
</p>
`

	htmlContent, err := notification.RenderHTML(suite.AppContextForTest(), ppmEmailData)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)
}

func (suite *NotificationSuite) TestPpmPacketEmailTextTemplateRender() {

	var pickupAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: pickupAddressModel,
		},
	}, nil)
	var destinationAddress = factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: destinationAddressModel,
		},
	}, nil)

	customAffiliation := models.AffiliationARMY
	serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
		{Model: models.ServiceMember{
			Affiliation: &customAffiliation,
		}},
	}, nil)
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
	}, nil)
	shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	customPPM := models.PPMShipment{
		ID:                         uuid.Must(uuid.NewV4()),
		ShipmentID:                 shipment.ID,
		Status:                     models.PPMShipmentStatusWaitingOnCustomer,
		PickupPostalAddressID:      &pickupAddress.ID,
		DestinationPostalAddressID: &destinationAddress.ID,
		PickupPostalCode:           "79329",
		DestinationPostalCode:      "90210",
	}

	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, []factory.Customization{
		{Model: customPPM},
	})

	notification := NewPpmPacketEmail(ppmShipment.ID)

	ppmEmailData, _, err := GetEmailData(*notification, suite.AppContextForTest())
	suite.NoError(err)

	expectedTextContent := `*** DO NOT REPLY directly to this email ***

This is a confirmation that your Personally Procured Move (PPM) with the assigned move code ` + move.Locator + ` from ` + pickupAddress.City + `, ` + pickupAddress.State + ` to ` + destinationAddress.City + `, ` + destinationAddress.State + ` has been processed in MilMove.

Next steps:

For ` + affiliationDisplayValue[*serviceMember.Affiliation] + ` personnel (FURTHER ACTION REQUIRED):

You can now log into MilMove <https://my.move.mil/> and download your payment packet to submit to ` + armySubmitLocation + `. You must complete this step to receive final settlement of your PPM.

Note: Not all claimed expenses may have been accepted during PPM closeout if they did not meet the definition of a valid expense.

If you have any questions, contact a government transportation office. You can see a listing of transportation offices on Military One Source here: <https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL>

Thank you,

USTRANSCOM MilMove Team


The information contained in this email may contain Privacy Act information and is therefore protected under the
Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.
`

	textContent, err := notification.RenderText(suite.AppContextForTest(), ppmEmailData)

	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}

func (suite *NotificationSuite) TestPpmPacketEmailZipcodeFallback() {
	customAffiliation := models.AffiliationAIRFORCE
	serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
		{Model: models.ServiceMember{
			Affiliation: &customAffiliation,
		}},
	}, nil)
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
	}, nil)
	shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	customPPM := models.PPMShipment{
		ID:                    uuid.Must(uuid.NewV4()),
		ShipmentID:            shipment.ID,
		Status:                models.PPMShipmentStatusWaitingOnCustomer,
		PickupPostalCode:      "79329",
		DestinationPostalCode: "90210",
	}

	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, []factory.Customization{
		{Model: customPPM},
	})
	notification := NewPpmPacketEmail(ppmShipment.ID)

	ppmEmailData, _, err := GetEmailData(*notification, suite.AppContextForTest())
	suite.NoError(err)
	suite.NotNil(ppmEmailData)

	suite.EqualExportedValues(ppmEmailData, PpmPacketEmailData{
		OriginZIP:      &customPPM.PickupPostalCode,
		DestinationZIP: &customPPM.DestinationPostalCode,
		SubmitLocation: allOtherSubmitLocation,
		ServiceBranch:  affiliationDisplayValue[*serviceMember.Affiliation],
		Locator:        move.Locator,
	})

	expectedHTMLContent := `<p>*** DO NOT REPLY directly to this email ***</p>
<p>This is a confirmation that your Personally Procured Move (PPM) with the <strong>assigned move code ` + move.Locator + `</strong> from <strong>` + *ppmEmailData.OriginZIP + ` to ` + *ppmEmailData.DestinationZIP + ` </strong>has been processed in MilMove. </p>
<h4>Next steps:</h4>

<p>For ` + affiliationDisplayValue[*serviceMember.Affiliation] + ` personnel (FURTHER ACTION REQUIRED):</p>
<p>You can now log into MilMove <a href="https://my.move.mil/">https://my.move.mil/</a> and download your payment packet to submit to ` + allOtherSubmitLocation + `. <strong>You must complete this step to receive final settlement of your PPM.</strong></p>
<p>Note: The Transportation Office does not determine claimable expenses. Claimable expenses will be determined by finance.</p>

<p>If you have any questions, contact a government transportation office. You can see a listing of transportation offices on Military One Source here: <a href="https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL">https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL</a></p>

<p>Thank you,</p>

<p>USTRANSCOM MilMove Team</p>

<p>
  The information contained in this email may contain Privacy Act information and is therefore protected under the
  Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.
</p>
`

	htmlContent, err := notification.RenderHTML(suite.AppContextForTest(), ppmEmailData)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)
}
