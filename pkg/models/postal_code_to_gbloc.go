package models

// PostalCodeToGBLOC is a mapping from Postal Codes to GBLOCs
type PostalCodeToGBLOC struct {
	PostalCode string `db:"postal_code"`
	GBLOC      string `db:"gbloc"`
}

// PostalCodeToGBLOCs is not required by pop and may be deleted
type PostalCodeToGBLOCs []PostalCodeToGBLOC
