package models_test

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicMovingDocumentInstantiation() {
	expenseDoc := &models.MovingExpenseDocument{}

	expErrors := map[string][]string{
		"move_document_id":       {"MoveDocumentID can not be blank."},
		"payment_method":         {"PaymentMethod can not be blank."},
		"moving_expense_type":    {"MovingExpenseType can not be blank."},
		"requested_amount_cents": {"0 is not greater than 0."},
	}

	suite.verifyValidationErrors(expenseDoc, expErrors)
}

func (suite *ModelSuite) TestStorageExpenseDaysInStorage() {
	documentType := models.MovingExpenseTypeSTORAGE
	startDate := time.Date(2018, 05, 12, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2018, 05, 15, 0, 0, 0, 0, time.UTC)
	storageExpense := models.MovingExpenseDocument{
		MovingExpenseType:    documentType,
		RequestedAmountCents: 1000,
		PaymentMethod:        "GTCC",
		ReceiptMissing:       true,
		StorageStartDate:     &startDate,
		StorageEndDate:       &endDate,
	}
	daysInStorage, err := storageExpense.DaysInStorage()
	suite.Nil(err)
	suite.Equal(daysInStorage, 3)
}

func (suite *ModelSuite) TestStorageExpenseDaysInStorageSameDay() {
	documentType := models.MovingExpenseTypeSTORAGE
	startDate := time.Date(2018, 05, 12, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2018, 05, 12, 1, 0, 0, 0, time.UTC)
	storageExpense := models.MovingExpenseDocument{
		MovingExpenseType:    documentType,
		RequestedAmountCents: 1000,
		PaymentMethod:        "GTCC",
		ReceiptMissing:       true,
		StorageStartDate:     &startDate,
		StorageEndDate:       &endDate,
	}
	daysInStorage, err := storageExpense.DaysInStorage()
	suite.Nil(err)
	suite.Equal(daysInStorage, 0)
}

func (suite *ModelSuite) TestStorageExpenseDaysInStorageStartGreaterThanEnd() {
	documentType := models.MovingExpenseTypeSTORAGE
	startDate := time.Date(2018, 05, 13, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2018, 05, 12, 0, 0, 0, 0, time.UTC)
	storageExpense := models.MovingExpenseDocument{
		MovingExpenseType:    documentType,
		RequestedAmountCents: 1000,
		PaymentMethod:        "GTCC",
		ReceiptMissing:       true,
		StorageStartDate:     &startDate,
		StorageEndDate:       &endDate,
	}
	daysInStorage, err := storageExpense.DaysInStorage()
	suite.NotNil(err)
	suite.Equal(daysInStorage, 0)
}

func (suite *ModelSuite) TestStorageExpenseDaysInStorageMissingDate() {
	documentType := models.MovingExpenseTypeSTORAGE
	storageExpense := models.MovingExpenseDocument{
		MovingExpenseType:    documentType,
		RequestedAmountCents: 1000,
		PaymentMethod:        "GTCC",
		ReceiptMissing:       true,
	}
	daysInStorage, err := storageExpense.DaysInStorage()
	suite.NotNil(err)
	suite.Equal(daysInStorage, 0)
}

func (suite *ModelSuite) TestStorageExpenseDaysNonStorageExpense() {
	documentType := models.MovingExpenseTypeRENTALEQUIPMENT
	storageExpense := models.MovingExpenseDocument{
		MovingExpenseType:    documentType,
		RequestedAmountCents: 1000,
		PaymentMethod:        "GTCC",
		ReceiptMissing:       true,
	}
	daysInStorage, err := storageExpense.DaysInStorage()
	suite.NotNil(err)
	suite.Equal(daysInStorage, 0)
}
