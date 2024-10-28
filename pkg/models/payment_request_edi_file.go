package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
)

type PaymentRequestEdiFile struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
	EdiString            string    `json:"edi_string" db:"edi_string"`
	Filename             string    `json:"file_name" db:"file_name"`
	PaymentRequestNumber string    `json:"payment_request_number" db:"payment_request_number"`
}

func (p PaymentRequestEdiFile) TableName() string {
	return "payment_request_edi_files"
}

type PaymentRequestEdiFiles []PaymentRequestEdiFile

func CreatePaymentRequestEdiFile(db *pop.Connection, fileName string, ediString string, paymentRequestNumber string) error {
	paymentRequestEdiFile := &PaymentRequestEdiFile{
		Filename:             fileName,
		EdiString:            ediString,
		PaymentRequestNumber: paymentRequestNumber,
	}

	if paymentRequestEdiFile.EdiString == "" {
		return nil
	}

	if paymentRequestEdiFile.Filename == "" {
		return nil
	}

	if paymentRequestEdiFile.PaymentRequestNumber == "" {
		return nil
	}

	verrs, err := db.ValidateAndCreate(paymentRequestEdiFile)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		return verrs
	}
	return nil
}

func FetchAllPaymentRequestEdiFiles(db *pop.Connection) (PaymentRequestEdiFiles, error) {
	var paymentRequestEdiFiles PaymentRequestEdiFiles
	err := db.All(&paymentRequestEdiFiles)
	if err != nil {
		return nil, err
	}
	return paymentRequestEdiFiles, nil
}

func FetchPaymentRequestEdiByPaymentRequestNumber(db *pop.Connection, paymentRequestNumber string) (PaymentRequestEdiFile, error) {
	var paymentRequestEdiFile PaymentRequestEdiFile
	err := db.Where("payment_request_number = ?", paymentRequestNumber).First(&paymentRequestEdiFile)
	if err != nil {
		return PaymentRequestEdiFile{}, err
	}
	return paymentRequestEdiFile, nil
}
