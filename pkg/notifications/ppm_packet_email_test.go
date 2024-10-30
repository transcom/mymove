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
}

var destinationAddressModel = models.Address{
	ID:             uuid.Must(uuid.NewV4()),
	StreetAddress1: "2 Second St",
	StreetAddress2: models.StringPointer("Bldg 2"),
	City:           "Key West",
	State:          "FL",
	PostalCode:     "33040",
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
		ID:                           uuid.Must(uuid.NewV4()),
		ShipmentID:                   shipment.ID,
		Status:                       models.PPMShipmentStatusWaitingOnCustomer,
		PickupAddressID:              &pickupAddress.ID,
		DestinationAddressID:         &destinationAddress.ID,
		IsActualExpenseReimbursement: models.BoolPointer(false),
	}

	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, []factory.Customization{
		{Model: customPPM},
	})
	notification := NewPpmPacketEmail(ppmShipment.ID)

	ppmEmailData, _, err := notification.GetEmailData(suite.AppContextForTest())
	suite.NoError(err)
	suite.NotNil(ppmEmailData)

	suite.EqualExportedValues(ppmEmailData, PpmPacketEmailData{
		OriginCity:                        &pickupAddress.City,
		OriginState:                       &pickupAddress.State,
		OriginZIP:                         &pickupAddress.PostalCode,
		DestinationCity:                   &destinationAddress.City,
		DestinationState:                  &destinationAddress.State,
		DestinationZIP:                    &destinationAddress.PostalCode,
		SubmitLocation:                    allOtherSubmitLocation,
		ServiceBranch:                     affiliationDisplayValue[*serviceMember.Affiliation],
		Locator:                           move.Locator,
		IsActualExpenseReimbursement:      notification.ConvertBoolToString(ppmShipment.IsActualExpenseReimbursement),
		OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
		WashingtonHQServicesLink:          WashingtonHQServicesLink,
		MyMoveLink:                        MyMoveLink,
		SmartVoucherLink:                  SmartVoucherLink,
	})

	expectedHTMLContent := "<p>*** DO NOT REPLY directly to this email ***</p>\n<p>This is a confirmation that your Personally Procured Move (PPM) with the <strong>assigned move code " + move.Locator + "</strong> from <strong>" + pickupAddress.City + ", " + pickupAddress.State + "</strong> to <strong>" + destinationAddress.City + ", " + destinationAddress.State + "</strong> has been processed in MilMove.</p>\n<h4>Next steps:</h4>\n\n<p>For Air Force and Space Force personnel (FURTHER ACTION REQUIRED):</p>\n<p>Log in to SmartVoucher at <a href=\"https://smartvoucher.dfas.mil/\">https://smartvoucher.dfas.mil/</a> using your CAC or myPay username and password. This will allow you to edit your voucher, and complete and sign DD Form 1351-2.</p>\n\n<p>You can now log into MilMove <a href=\"https://my.move.mil/\">https://my.move.mil/</a> and download your payment packet to submit to your local finance office. <strong>You must complete this step to receive final settlement of your PPM.</strong></p>\n<p>Note: The Transportation Office does not determine claimable expenses. Claimable expenses will be determined by finance.</p>\n<p>Please be advised, your local finance office may require a DD Form 1351-2 to process payment. You can obtain a copy of this form by utilizing the search feature at <a href=\"https://www.esd.whs.mil\">https://www.esd.whs.mil</a>.</p>\n\n<p>If you have any questions, contact a government transportation office. You can see a listing of transportation offices on Military One Source here: <a href=\"https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL\">https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL</a></p>\n<p>Thank you,</p>\n<p>USTRANSCOM MilMove Team</p>\n<p>\n  The information contained in this email may contain Privacy Act information and is therefore protected under the\n  Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.\n</p>\n\n"

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
		ID:                           uuid.Must(uuid.NewV4()),
		ShipmentID:                   shipment.ID,
		Status:                       models.PPMShipmentStatusWaitingOnCustomer,
		PickupAddressID:              &pickupAddress.ID,
		DestinationAddressID:         &destinationAddress.ID,
		IsActualExpenseReimbursement: models.BoolPointer(false),
	}

	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, []factory.Customization{
		{Model: customPPM},
	})
	notification := NewPpmPacketEmail(ppmShipment.ID)

	ppmEmailData, _, err := notification.GetEmailData(suite.AppContextForTest())
	suite.NoError(err)
	suite.NotNil(ppmEmailData)

	suite.EqualExportedValues(ppmEmailData, PpmPacketEmailData{
		OriginCity:                        &pickupAddress.City,
		OriginState:                       &pickupAddress.State,
		OriginZIP:                         &pickupAddress.PostalCode,
		DestinationCity:                   &destinationAddress.City,
		DestinationState:                  &destinationAddress.State,
		DestinationZIP:                    &destinationAddress.PostalCode,
		SubmitLocation:                    armySubmitLocation,
		ServiceBranch:                     affiliationDisplayValue[*serviceMember.Affiliation],
		Locator:                           move.Locator,
		IsActualExpenseReimbursement:      notification.ConvertBoolToString(ppmShipment.IsActualExpenseReimbursement),
		OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
		WashingtonHQServicesLink:          WashingtonHQServicesLink,
		MyMoveLink:                        MyMoveLink,
		SmartVoucherLink:                  SmartVoucherLink,
	})

	expectedHTMLContent := "<p>*** DO NOT REPLY directly to this email ***</p>\n<p>This is a confirmation that your Personally Procured Move (PPM) with the <strong>assigned move code " + move.Locator + "</strong> from <strong>" + pickupAddress.City + ", " + pickupAddress.State + "</strong> to <strong>" + destinationAddress.City + ", " + destinationAddress.State + "</strong> has been processed in MilMove.</p>\n<h4>Next steps:</h4>\n\n<p>For Army personnel (FURTHER ACTION REQUIRED):</p>\n<p>Log in to SmartVoucher at <a href=\"https://smartvoucher.dfas.mil/\">https://smartvoucher.dfas.mil/</a> using your CAC or myPay username and password. This will allow you to edit your voucher, and complete and sign DD Form 1351-2.</p>\n\n<p>Please be advised, your local finance office may require a DD Form 1351-2 to process payment. You can obtain a copy of this form by utilizing the search feature at <a href=\"https://www.esd.whs.mil\">https://www.esd.whs.mil</a>.</p>\n\n<p>If you have any questions, contact a government transportation office. You can see a listing of transportation offices on Military One Source here: <a href=\"https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL\">https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL</a></p>\n<p>Thank you,</p>\n<p>USTRANSCOM MilMove Team</p>\n<p>\n  The information contained in this email may contain Privacy Act information and is therefore protected under the\n  Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.\n</p>\n\n"

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
		ID:                   uuid.Must(uuid.NewV4()),
		ShipmentID:           shipment.ID,
		Status:               models.PPMShipmentStatusWaitingOnCustomer,
		PickupAddressID:      &pickupAddress.ID,
		DestinationAddressID: &destinationAddress.ID,
	}

	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, []factory.Customization{
		{Model: customPPM},
	})
	notification := NewPpmPacketEmail(ppmShipment.ID)

	ppmEmailData, _, err := notification.GetEmailData(suite.AppContextForTest())
	suite.NoError(err)
	suite.NotNil(ppmEmailData)

	suite.EqualExportedValues(ppmEmailData, PpmPacketEmailData{
		OriginCity:                        &pickupAddress.City,
		OriginState:                       &pickupAddress.State,
		OriginZIP:                         &pickupAddress.PostalCode,
		DestinationCity:                   &destinationAddress.City,
		DestinationState:                  &destinationAddress.State,
		DestinationZIP:                    &destinationAddress.PostalCode,
		SubmitLocation:                    allOtherSubmitLocation,
		ServiceBranch:                     affiliationDisplayValue[*serviceMember.Affiliation],
		Locator:                           move.Locator,
		IsActualExpenseReimbursement:      notification.ConvertBoolToString(ppmShipment.IsActualExpenseReimbursement),
		OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
		WashingtonHQServicesLink:          WashingtonHQServicesLink,
		MyMoveLink:                        MyMoveLink,
		SmartVoucherLink:                  SmartVoucherLink,
	})

	expectedHTMLContent := "<p>*** DO NOT REPLY directly to this email ***</p>\n<p>This is a confirmation that your Personally Procured Move (PPM) with the <strong>assigned move code " + move.Locator + "</strong> from <strong>" + pickupAddress.City + ", " + pickupAddress.State + "</strong> to <strong>" + destinationAddress.City + ", " + destinationAddress.State + "</strong> has been processed in MilMove.</p>\n<h4>Next steps:</h4>\n\n<p>For Marine Corps, Navy, and Coast Guard personnel:</p>\n<p>You can now log into MilMove <a href=\"https://my.move.mil/\">https://my.move.mil/</a> and view your payment packet; however, you do not need to forward your payment packet to finance as your closeout location is associated with your finance office and they will handle this step for you.</p>\n<p>Note: Not all claimed expenses may have been accepted during PPM closeout if they did not meet the definition of a valid expense.</p>\n\n<p>If you have any questions, contact a government transportation office. You can see a listing of transportation offices on Military One Source here: <a href=\"https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL\">https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL</a></p>\n<p>Thank you,</p>\n<p>USTRANSCOM MilMove Team</p>\n<p>\n  The information contained in this email may contain Privacy Act information and is therefore protected under the\n  Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.\n</p>\n\n"

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
		ID:                   uuid.Must(uuid.NewV4()),
		ShipmentID:           shipment.ID,
		Status:               models.PPMShipmentStatusWaitingOnCustomer,
		PickupAddressID:      &pickupAddress.ID,
		DestinationAddressID: &destinationAddress.ID,
	}

	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, []factory.Customization{
		{Model: customPPM},
	})

	notification := NewPpmPacketEmail(ppmShipment.ID)

	ppmEmailData, _, err := notification.GetEmailData(suite.AppContextForTest())
	suite.NoError(err)

	expectedTextContent := "*** DO NOT REPLY directly to this email ***\n\nThis is a confirmation that your Personally Procured Move (PPM) with the assigned move code " + move.Locator + " from " + pickupAddress.City + ", " + pickupAddress.State + " to " + destinationAddress.City + ", " + destinationAddress.State + " has been processed in MilMove.\n\nNext steps:\n\nFor Army personnel (FURTHER ACTION REQUIRED):\n\nLog in to SmartVoucher at https://smartvoucher.dfas.mil/ using your CAC or myPay username and password. This will allow you to edit your voucher, and complete and sign DD Form 1351-2.\n\n\n\nPlease be advised, your local finance office may require a DD Form 1351-2 to process payment. You can obtain a copy of this form by utilizing the search feature at https://www.esd.whs.mil.\n\nIf you have any questions, contact a government transportation office. You can see a listing of transportation offices on Military One Source here: https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL\n\nThank you,\n\nUSTRANSCOM MilMove Team\n\nThe information contained in this email may contain Privacy Act information and is therefore protected under the\nPrivacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.\n\n"

	textContent, err := notification.RenderText(suite.AppContextForTest(), ppmEmailData)

	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}

func (suite *NotificationSuite) TestPpmPacketEmailTextTemplateRenderForArmyWithActualExpenseReimbursement() {
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
		ID:                           uuid.Must(uuid.NewV4()),
		ShipmentID:                   shipment.ID,
		Status:                       models.PPMShipmentStatusWaitingOnCustomer,
		PickupAddressID:              &pickupAddress.ID,
		DestinationAddressID:         &destinationAddress.ID,
		IsActualExpenseReimbursement: models.BoolPointer(true),
	}

	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, []factory.Customization{
		{Model: customPPM},
	})

	notification := NewPpmPacketEmail(ppmShipment.ID)

	ppmEmailData, _, err := notification.GetEmailData(suite.AppContextForTest())
	suite.NoError(err)

	expectedTextContent := "*** DO NOT REPLY directly to this email ***\n\nThis is a confirmation that your Personally Procured Move (PPM) with the assigned move code " + move.Locator + " from " + pickupAddress.City + ", " + pickupAddress.State + " to " + destinationAddress.City + ", " + destinationAddress.State + " has been processed in MilMove.\n\nNext steps:\n\nFor Army personnel (FURTHER ACTION REQUIRED):\n\nLog in to SmartVoucher at https://smartvoucher.dfas.mil/ using your CAC or myPay username and password. This will allow you to edit your voucher, and complete and sign DD Form 1351-2.\n\n\n\nPlease be advised, your local finance office may require a DD Form 1351-2 to process payment. You can obtain a copy of this form by utilizing the search feature at https://www.esd.whs.mil.\n\nPlease Note: Your PPM has been designated as Actual Expense Reimbursement. This is the standard entitlement for Civilian employees. For uniformed Service Members, your PPM may have been designated as Actual Expense Reimbursement due to failure to receive authorization prior to movement or failure to obtain certified weight tickets. Actual Expense Reimbursement means reimbursement for expenses not to exceed the Government Constructed Cost (GCC).\n\nIf you have any questions, contact a government transportation office. You can see a listing of transportation offices on Military One Source here: https://installations.militaryonesource.mil/search?program-service=2/view-by=ALL\n\nThank you,\n\nUSTRANSCOM MilMove Team\n\nThe information contained in this email may contain Privacy Act information and is therefore protected under the\nPrivacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine.\n\n"

	textContent, err := notification.RenderText(suite.AppContextForTest(), ppmEmailData)

	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}
