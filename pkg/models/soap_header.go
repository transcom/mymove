package models

type Security struct {
	Text                string              `xml:"chardata"`
	Wsse                string              `xml:"wsse"`
	BinarySecurityToken BinarySecurityToken `xml:"BinarySecurityToken"`
	Signature           Signature           `xml:"Signature"`
}

type BinarySecurityToken struct {
	Text         string `xml:"chardata"`
	EncodingType string `xml:"EncodingType"`
	ValueType    string `xml:"ValueType"`
	ID           string `xml:"Id"`
}

type Signature struct {
	Text           string `xml:"chardata"`
	Ds             string `xml:"ds"`
	SignedInfo     string `xml:"SignedInfo"`
	SignatureValue string `xml:"SignatureValue"`
}
