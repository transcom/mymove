package move

func (suite *MoveServiceSuite) TestMoveUpdate() {
	// move := factory.BuildMove(suite.DB(), nil, nil)

	// updateMove()
}

func (suite *MoveServiceSuite) TestAdditionalDocumentUploader() {
	// fakeFileStorer := test.NewFakeS3Storage(true)
	// uploadCreator := upload.NewUploadCreator(fakeFileStorer)

	// testFileName := "upload-test.pdf"
	// additionalDocumentUploader := NewMoveAdditionalDocumentsUploader(uploadCreator)

	// suite.Run("Success - Additional Document upload is created", func() {
	// 	testFile, fileErr := os.Open("../../testdatagen/testdata/test.pdf")
	// 	suite.Require().NoError(fileErr)

	// 	defer func() {
	// 		closeErr := testFile.Close()
	// 		suite.NoError(closeErr, "Error occurred while closing the test file.")
	// 	}()

	// 	move := factory.BuildMove(suite.DB(), nil, nil)

	// 	updatedMove, _, _, err := additionalDocumentUploader.CreateAdditionalDocumentsUpload(
	// 		suite.AppContextForTest(),
	// 		move.Orders.ServiceMemberID,
	// 		move.ID,
	// 		testFile,
	// 		testFileName,
	// 		fakeFileStorer)
	// 	suite.NoError(err)
	// 	suite.Require().NotNil(updatedMove)

	// 	suite.NotNil(updatedMove.ExcessWeightUploadID)
	// 	suite.Require().NotNil(updatedMove.ExcessWeightUpload)
	// 	suite.Equal(updatedMove.ExcessWeightUpload.ID, *updatedMove.ExcessWeightUploadID)

	// 	suite.Equal(models.UploadTypePRIME, updatedMove.ExcessWeightUpload.UploadType)
	// 	suite.Contains(updatedMove.ExcessWeightUpload.Filename, testFileName)
	// 	suite.Contains(updatedMove.ExcessWeightUpload.StorageKey, testFileName)
	// })

	// suite.Run("Fail - Cannot create upload for move", func() {
	// 	move := factory.BuildMove(suite.DB(), nil, nil)

	// 	testFile, fileErr := os.Open("../../testdatagen/testdata/test.pdf")
	// 	suite.Require().NoError(fileErr)

	// 	defer func() {
	// 		closeErr := testFile.Close()
	// 		suite.NoError(closeErr, "Error occurred while closing the test file.")
	// 	}()

	// 	updatedMove, err := add.CreateExcessWeightUpload(
	// 		suite.AppContextForTest(), move.ID, testFile, testFileName, models.UploadTypePRIME)
	// 	suite.Nil(updatedMove)
	// 	suite.Require().Error(err)

	// 	suite.IsType(apperror.NotFoundError{}, err)
	// 	suite.Contains(err.Error(), move.ID.String())
	// })
}
