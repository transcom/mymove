package models

import "encoding/xml"

type Header struct {
	XMLName  xml.Name `xml:"Header"`
	Security Security `xml:"Security"`
}

type Security struct {
	Wsse                string              `xml:"wsse,attr"`
	Wsu                 string              `xml:"wsu,attr"`
	BinarySecurityToken BinarySecurityToken `xml:"BinarySecurityToken"`
	Signature           Signature           `xml:"Signature"`
	Timestamp           Timestamp           `xml:"TimeStamp"`
}
type Signature struct {
	Text       string     `xml:",chardata"`
	ID         string     `xml:"Id,attr"`
	Ds         string     `xml:"ds,attr"`
	SignedInfo SignedInfo `xml:"SignedInfo"`
}

type BinarySecurityToken struct {
	Text         string `xml:",chardata"`
	EncodingType string `xml:"EncodingType,attr"`
	ValueType    string `xml:"ValueType,attr"`
	ID           string `xml:"Id,attr"`
}

type Reference []struct {
	URI        string `xml:"URI,attr"`
	Transforms struct {
		Text      string `xml:",chardata"`
		Transform struct {
			Text      string `xml:",chardata"`
			Algorithm string `xml:"Algorithm,attr"`
		} `xml:"Transform"`
	} `xml:"Transforms"`
	DigestMethod struct {
		Text      string `xml:",chardata"`
		Algorithm string `xml:"Algorithm,attr"`
	} `xml:"DigestMethod"`
	DigestValue string `xml:"DigestValue"`
}

type SignedInfo struct {
	CanonicalizationMethod CanonicalizationMethod `xml:"CanonicalizationMethod"`
	SignatureMethod        SignatureMethod        `xml:"SignatureMethod"`
	Reference              Reference              `xml:"Reference"`
}
type CanonicalizationMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type SignatureMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}
type Timestamp struct {
	ID      string `xml:"Id,attr"`
	Created string `xml:"Created"`
	Expires string `xml:"Expires"`
}
