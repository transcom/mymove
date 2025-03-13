package weightticketparser

import (
	"os"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	paperworkgenerator "github.com/transcom/mymove/pkg/paperwork"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *WeightTicketParserServiceSuite) TestFillWeightEstimatorPDFForm() {
	const WeightEstimatorPages = 11
	const WeightEstimatorFileName = "Weight Estimator Full.xlsx"
	const TestPath = "../../testdatagen/testdata/"
	fakeS3 := storageTest.NewFakeS3Storage(true)
	userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, 25*uploader.MB)
	suite.FatalNoError(uploaderErr)
	generator, err := paperworkgenerator.NewGenerator(userUploader.Uploader())
	suite.FatalNil(err)

	weightParserComputer := NewWeightTicketComputer()
	weightParserGenerator, err := NewWeightTicketParserGenerator(generator)
	suite.FatalNoError(err)
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	yuma := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	grade := models.ServiceMemberGradeE9
	ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OrdersType: ordersType,
				Grade:      &grade,
			},
		},
		{
			Model:    fortGordon,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model:    yuma,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.SignedCertification{
				UpdatedAt: time.Now(),
			},
		},
	}, nil)

	_, err = models.SaveMoveDependencies(suite.DB(), &ppmShipment.Shipment.MoveTaskOrder)
	suite.NoError(err)

	excelFile, err := os.Open(TestPath + WeightEstimatorFileName)
	suite.Require().NoError(err)

	defer func() {
		closeErr := excelFile.Close()
		suite.NoError(closeErr, "Error occurred while closing the test excel file.")
	}()

	// Parse our excel file to get data for the pdf
	estimatorPages, err := weightParserComputer.ParseWeightEstimatorExcelFile(suite.AppContextForTest(), excelFile)
	suite.NoError(err)

	// Fill our pdf template with data from parser
	testFile, pdfInfo, err := weightParserGenerator.FillWeightEstimatorPDFForm(*estimatorPages, WeightEstimatorFileName)

	suite.NoError(err)
	println(testFile.Name())                             // ensures was generated with temp filesystem
	suite.Equal(pdfInfo.PageCount, WeightEstimatorPages) // ensures PDF is not corrupted
	err = weightParserGenerator.CleanupFile(testFile)
	suite.NoError(err)
}
