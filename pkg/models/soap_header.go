package models

import "encoding/xml"

type Header struct {
	XMLName  xml.Name `xml:"Header"`
	Text     string   `xml:"chardata"`
	Security Security `xml:"Security"`
}

type Security struct {
	Text                string              `xml:"chardata"`
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

type SignedInfo struct {
	Text                   string                 `xml:",chardata"`
	CanonicalizationMethod CanonicalizationMethod `xml:"CanonicalizationMethod"`
	SignatureMethod        SignatureMethod        `xml:"SignatureMethod"`
}
type CanonicalizationMethod struct {
	Text      string `xml:",chardata"`
	Algorithm string `xml:"Algorithm,attr"`
}

type SignatureMethod struct {
	Text      string `xml:",chardata"`
	Algorithm string `xml:"Algorithm,attr"`
}
type Timestamp struct {
	Text    string `xml:",chardata"`
	ID      string `xml:"Id,attr"`
	Created string `xml:"Created"`
	Expires string `xml:"Expires"`
}
