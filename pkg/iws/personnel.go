package iws

// Personnel contains the organizational and pay information corresponding to a Person
type Personnel struct {
	// The code that represents how the DoD personnel and/or finance center views the sponsor based on accountability and reporting strengths. (This attribute is similar to Member Category Code.)
	PnlCatCd PersonnelCategoryCode `xml:"PNL_CAT_CD,omitempty"`
	// The code that represents the organization that "owns" the DEERS population to which the sponsor belongs.
	OrgCd OrgCode `xml:"ORG_CD,omitempty"`
	// The text of a person's or organization's email address in the format xxx@xxxxxx.
	Email string `xml:"EMA_TX,omitempty"`
	// The code that represents the sponsor's rank.
	RankCd string `xml:"RANK_CD,omitempty"`
	// The code that represents the level of pay. (The combination of pay plan code and pay grade code represents the sponsor's pay category.)
	PgCd PayGradeCode `xml:"PG_CD,omitempty"`
	// The code that represents the type of pay category. (The combination of pay plan code and pay grade code represents the sponsor's pay category.)
	PayPlanCd PayPlanCode `xml:"PAY_PLN_CD,omitempty"`
	// The code that represents the branch classification of Service with which the sponsor is affiliated.
	SvcCd ServiceCode `xml:"SVC_CD,omitempty"`
}

// PersonnelCategoryCode represents how the DoD personnel and/or finance center views the sponsor based on accountability and reporting strengths. (This attribute is similar to Member Category Code.)
type PersonnelCategoryCode string

const (
	// PersonnelCategoryCodeActiveDuty indicates an active duty member
	PersonnelCategoryCodeActiveDuty PersonnelCategoryCode = "A"
	// PersonnelCategoryCodeAppointee indicates Presidential Appointees of all Federal Government agencies
	PersonnelCategoryCodeAppointee PersonnelCategoryCode = "B"
	// PersonnelCategoryCodeDODCivilService indicates a DoD and Uniformed Service civil service employee, except Presidential appointee
	PersonnelCategoryCodeDODCivilService PersonnelCategoryCode = "C"
	// PersonnelCategoryCodeDisabled indicates a disabled American veteran
	PersonnelCategoryCodeDisabled PersonnelCategoryCode = "D"
	// PersonnelCategoryCodeDODContractEmployee indicates a DoD and Uniformed Service contract employee
	PersonnelCategoryCodeDODContractEmployee PersonnelCategoryCode = "E"
	// PersonnelCategoryCodeFormer indicates a former member (Reserve service, discharged from RR or SR following notification of retirement eligibility)
	PersonnelCategoryCodeFormer PersonnelCategoryCode = "F"
	// PersonnelCategoryCodeMedalOfHonor indicates a Medal of Honor recipient
	PersonnelCategoryCodeMedalOfHonor PersonnelCategoryCode = "H"
	// PersonnelCategoryCodeNonDODCivilService indicates a Non-DoD civil service employee, except Presidential appointee
	PersonnelCategoryCodeNonDODCivilService PersonnelCategoryCode = "I"
	// PersonnelCategoryCodeAcademyStudent indicates a Service Academy student
	PersonnelCategoryCodeAcademyStudent PersonnelCategoryCode = "J"
	// PersonnelCategoryCodeNAF indicates a non-appropriated fund DoD and Uniformed Service employee (NAF)
	PersonnelCategoryCodeNAF PersonnelCategoryCode = "K"
	// PersonnelCategoryCodeLighthouse indicates Lighthouse service - Obsolete
	PersonnelCategoryCodeLighthouse PersonnelCategoryCode = "L"
	// PersonnelCategoryCodeCivilianAssociate indicates non-federal Agency civilian associates
	PersonnelCategoryCodeCivilianAssociate PersonnelCategoryCode = "M"
	// PersonnelCategoryCodeNationalGuard indicates a National Guard member
	PersonnelCategoryCodeNationalGuard PersonnelCategoryCode = "N"
	// PersonnelCategoryCodeNonDODContractEmployee indicates a Non-DoD contract employee
	PersonnelCategoryCodeNonDODContractEmployee PersonnelCategoryCode = "O"
	// PersonnelCategoryCodeGrayAreaRetiree indicates a Reserve retiree not yet eligible for retired pay ("Gray Area Retiree")
	PersonnelCategoryCodeGrayAreaRetiree PersonnelCategoryCode = "Q"
	// PersonnelCategoryCodeRetiree indicates a retired military member eligible for retired pay
	PersonnelCategoryCodeRetiree PersonnelCategoryCode = "R"
	// PersonnelCategoryCodeForeignAffiliate indicates a Foreign Affiliate
	PersonnelCategoryCodeForeignAffiliate PersonnelCategoryCode = "T"
	// PersonnelCategoryCodeOCONUSHire indicates a DoD OCONUS Hire
	PersonnelCategoryCodeOCONUSHire PersonnelCategoryCode = "U"
	// PersonnelCategoryCodeReservist indicates a Reserve member
	PersonnelCategoryCodeReservist PersonnelCategoryCode = "V"
	// PersonnelCategoryCodeBeneficiary indicates a DoD Beneficiary, a person who receives benefits from the DoD based on prior association, condition or authorization, an example is a former spouse
	PersonnelCategoryCodeBeneficiary PersonnelCategoryCode = "W"
	// PersonnelCategoryCodeCivilianRetiree indicates a Civilian Retiree
	PersonnelCategoryCodeCivilianRetiree PersonnelCategoryCode = "Y"
)

// OrgCode represents the organization that "owns" the DEERS population to which the sponsor belongs.
type OrgCode string

const (
	// OrgCodeArmyAffiliate indicates an Army affiliate (used only for reporting, not on PNL)
	OrgCodeArmyAffiliate OrgCode = "01"
	// OrgCodeAirForceAffiliate indicates an Air Force affiliate (used only for reporting, not on PNL)
	OrgCodeAirForceAffiliate OrgCode = "02"
	// OrgCodeNavyAffiliate indicates a Navy affiliate (used only for reporting, not on PNL)
	OrgCodeNavyAffiliate OrgCode = "03"
	// OrgCodeMarineCorpsAffiliate indicates a Marine Corps affiliate (used only for reporting, not on PNL)
	OrgCodeMarineCorpsAffiliate OrgCode = "04"
	// OrgCodeCoastGuardAffiliate indicates a Coast Guard affiliate (used only for reporting, not on PNL)
	OrgCodeCoastGuardAffiliate OrgCode = "05"
	// OrgCodePublicHealthAffiliate indicates a Public Health affiliate (used only for reporting, not on PNL)
	OrgCodePublicHealthAffiliate OrgCode = "06"
	// OrgCodeNOAAAffiliate indicates a NOAA affiliate (used only for reporting, not on PNL)
	OrgCodeNOAAAffiliate OrgCode = "07"
	// OrgCodeArmyActive indicates Army MILPERCEN DEERS Population: Eligible Army Active Duty
	OrgCodeArmyActive OrgCode = "11"
	// OrgCodeAirForceActive indicates Air Force MILPERCEN DEERS Population: Eligible Air Force Active Duty
	OrgCodeAirForceActive OrgCode = "12"
	// OrgCodeNavyActive indicates Navy MILPERCEN DEERS Population: Eligible Navy Active Duty
	OrgCodeNavyActive OrgCode = "13"
	// OrgCodeMarineCorpsActive indicates Marine Corps MILPERCEN DEERS Population: Eligible Marine Corps Active Duty
	OrgCodeMarineCorpsActive OrgCode = "14"
	// OrgCodeCoastGuardActive indicates Coast Guard MILPERCEN DEERS Population: Eligible Coast Guard Active Duty
	OrgCodeCoastGuardActive OrgCode = "15"
	// OrgCodePublicHealthActive indicates Public Health PERCEN DEERS Population: Eligible Public Health Active
	OrgCodePublicHealthActive OrgCode = "16"
	// OrgCodeNOAAActive indicates NOAA PERCEN DEERS Population: Eligible NOAA Active
	OrgCodeNOAAActive OrgCode = "17"
	// OrgCodeArmyRetired indicates Army Retired Finance Center DEERS Population: Eligible Army Retired and Former Members
	OrgCodeArmyRetired OrgCode = "21"
	// OrgCodeAirForceRetired indicates Air Force Retired Finance Center DEERS Population: Eligible Air Force Retired and Former Members
	OrgCodeAirForceRetired OrgCode = "22"
	// OrgCodeNavyRetired indicates Navy Retired Finance Center DEERS Population: Eligible Navy Retired and Former Members
	OrgCodeNavyRetired OrgCode = "23"
	// OrgCodeMarineCorpsRetired indicates Marine Corps Finance Center DEERS Population: Eligible Marine Corps Retired and Former Members
	OrgCodeMarineCorpsRetired OrgCode = "24"
	// OrgCodeCoastGuardRetired indicates Coast Guard Retired Finance Center DEERS Population: Eligible Coast Guard Retired and Former Members
	OrgCodeCoastGuardRetired OrgCode = "25"
	// OrgCodePublicHealthRetired indicates Public Health Finance Center DEERS Population: Eligible Public Health Retired
	OrgCodePublicHealthRetired OrgCode = "26"
	// OrgCodeNOAARetired indicates NOAA Finance Center DEERS Population: Eligible NOAA Retired
	OrgCodeNOAARetired OrgCode = "27"
	// OrgCodeArmyCadet indicates Army Academy DEERS Population: Eligible Army Cadet
	OrgCodeArmyCadet OrgCode = "31"
	// OrgCodeAirForceCadet indicates Air Force Academy DEERS Population: Eligible Air Force Cadet
	OrgCodeAirForceCadet OrgCode = "32"
	// OrgCodeNavyCadet indicates Navy Academy DEERS Population: Eligible Navy Cadet and OCS
	OrgCodeNavyCadet OrgCode = "33"
	// OrgCodeCoastGuardCadet indicates Coast Guard Academy DEERS Population: Eligible Coast Guard Cadet
	OrgCodeCoastGuardCadet OrgCode = "35"
	// OrgCodeArmyReserve indicates Army Reserve DEERS Population: Eligible Army Reserve
	OrgCodeArmyReserve OrgCode = "41"
	// OrgCodeAirForceReserve indicates Air Force Reserve DEERS Population: Eligible Air Force Reserve
	OrgCodeAirForceReserve OrgCode = "42"
	// OrgCodeNavyReserve indicates Navy Reserve DEERS Population: Eligible Navy Reserve
	OrgCodeNavyReserve OrgCode = "43"
	// OrgCodeMarineCorpsReserve indicates Marine Corps Reserve DEERS Population: Eligible Marine Corps Reserve
	OrgCodeMarineCorpsReserve OrgCode = "44"
	// OrgCodeCoastGuardReserve indicates Coast Guard Reserve DEERS Population: Eligible Coast Guard Reserve
	OrgCodeCoastGuardReserve OrgCode = "45"
	// OrgCodePublicHealthReserve indicates Public Health Reserve DEERS Population: Eligible Public Health Reserve - obsolete
	OrgCodePublicHealthReserve OrgCode = "46"
	// OrgCodeArmyGuard indicates Army Guard DEERS Population: Eligible Army Guard
	OrgCodeArmyGuard OrgCode = "51"
	// OrgCodeAirForceGuard indicates Air Force Guard DEERS Population: Eligible Air Force Guard
	OrgCodeAirForceGuard OrgCode = "52"
	// OrgCodeChampva indicates CHAMPVA DEERS Population: Eligible Disabled American Veteran
	OrgCodeChampva OrgCode = "61"
	// OrgCodeCivilService indicates Civil service DEERS Population: Eligible civil service in DoD
	OrgCodeCivilService OrgCode = "62"
	// OrgCodeCivilianVerificationSystem indicates Civilian Verification System (future use)
	OrgCodeCivilianVerificationSystem OrgCode = "63"
	// OrgCodeCoastGuardCivilian indicates Coast Guard Civilian file
	OrgCodeCoastGuardCivilian OrgCode = "64"
	// OrgCodeNOAACivilian indicates NOAA Civilian Personnel File
	OrgCodeNOAACivilian OrgCode = "65"
	// OrgCodePublicHealthCivilian indicates Public Health Service Civilian Personnel File
	OrgCodePublicHealthCivilian OrgCode = "66"
	// OrgCodeSDVA indicates SDVA - State Offices of Veterans Affairs
	OrgCodeSDVA OrgCode = "67"
	// OrgCodeMerchantMarines indicates Merchant Marines
	OrgCodeMerchantMarines OrgCode = "78"
	// OrgCodeDEERSOnly indicates DEERS Population: Eligible and post-eligible personnel with DEERS online as the sole source (e.g., foreign national, foreign military)
	OrgCodeDEERSOnly OrgCode = "80"
	// OrgCodeVeteransFromMilitaryServiceHistory indicates Veterans from Military Service History Load
	OrgCodeVeteransFromMilitaryServiceHistory OrgCode = "86"
)

// PayGradeCode represents the level of pay. (The combination of pay plan code and pay grade code represents the sponsor's pay category.)
type PayGradeCode string

const (
	// PayGradeCodeUnknown00 means unknown paygrade
	PayGradeCodeUnknown00 PayGradeCode = "00"
	// PayGradeCode01 identifies level 01 in a Civil Service, Cadet, Warrant Officer, Enlisted or Officer pay plan
	PayGradeCode01 PayGradeCode = "01"
	// PayGradeCode02 identifies level 02 in a Civil Service, Warrant Officer, Enlisted or Officer pay plan
	PayGradeCode02 PayGradeCode = "02"
	// PayGradeCode03 identifies level 03 in a Civil Service, Warrant Officer, Enlisted or Officer pay plan
	PayGradeCode03 PayGradeCode = "03"
	// PayGradeCode04 identifies level 04 in a Civil Service, Warrant Officer, Enlisted or Officer pay plan
	PayGradeCode04 PayGradeCode = "04"
	// PayGradeCode05 identifies level 05 in a Civil Service, Warrant Officer, Enlisted or Officer pay plan
	PayGradeCode05 PayGradeCode = "05"
	// PayGradeCode06 identifies level 06 in a Civil Service, Enlisted, or Officer pay plan
	PayGradeCode06 PayGradeCode = "06"
	// PayGradeCode07 identifies level 07 in a Civil Service, Enlisted, or Officer pay plan
	PayGradeCode07 PayGradeCode = "07"
	// PayGradeCode08 identifies level 08 in a Civil Service, Enlisted, or Officer pay plan
	PayGradeCode08 PayGradeCode = "08"
	// PayGradeCode09 identifies level 09 in a Civil Service, Enlisted, or Officer pay plan
	PayGradeCode09 PayGradeCode = "09"
	// PayGradeCode10 identifies level 10 in a Civil Service or Officer pay plan
	PayGradeCode10 PayGradeCode = "10"
	// PayGradeCode11 identifies level 11 in a Civil Service or Officer pay plan
	PayGradeCode11 PayGradeCode = "11"
	// PayGradeCode12 identifies level 12 in a Civil Service Pay Plan
	PayGradeCode12 PayGradeCode = "12"
	// PayGradeCode13 identifies level 13 in a Civil Service Pay Plan
	PayGradeCode13 PayGradeCode = "13"
	// PayGradeCode14 identifies level 14 in a Civil Service Pay Plan
	PayGradeCode14 PayGradeCode = "14"
	// PayGradeCode15 identifies level 15 in a Civil Service Pay Plan
	PayGradeCode15 PayGradeCode = "15"
	// PayGradeCode21 identifies level 21 in a Civil Service pay plan
	PayGradeCode21 PayGradeCode = "21"
	// PayGradeCode22 identifies level 22 in a Civil Service pay plan
	PayGradeCode22 PayGradeCode = "22"
	// PayGradeCode23 identifies level 23 in a Civil Service pay plan
	PayGradeCode23 PayGradeCode = "23"
	// PayGradeCode24 identifies level 24 in a Civil Service pay plan
	PayGradeCode24 PayGradeCode = "24"
	// PayGradeCode25 identifies level 25 in a Civil Service pay plan
	PayGradeCode25 PayGradeCode = "25"
	// PayGradeCode26 identifies level 26 in a Civil Service pay plan
	PayGradeCode26 PayGradeCode = "26"
	// PayGradeCode27 identifies level 27 in a Civil Service pay plan
	PayGradeCode27 PayGradeCode = "27"
	// PayGradeCode28 identifies level 28 in a Civil Service pay plan
	PayGradeCode28 PayGradeCode = "28"
	// PayGradeCode29 identifies level 29 in a Civil Service pay plan
	PayGradeCode29 PayGradeCode = "29"
	// PayGradeCode30 identifies level 30 in a Civil Service pay plan
	PayGradeCode30 PayGradeCode = "30"
	// PayGradeCode31 identifies level 31 in a Civil Service pay plan
	PayGradeCode31 PayGradeCode = "31"
	// PayGradeCode32 identifies level 32 in a Civil Service pay plan
	PayGradeCode32 PayGradeCode = "32"
	// PayGradeCode34 identifies level 34 in a Civil Service pay plan
	PayGradeCode34 PayGradeCode = "34"
	// PayGradeCode36 identifies level 36 in a Civil Service pay plan
	PayGradeCode36 PayGradeCode = "36"
	// PayGradeCode39 identifies level 39 in a Civil Service pay plan
	PayGradeCode39 PayGradeCode = "39"
	// PayGradeCode40 identifies level 40 in a Civil Service pay plan
	PayGradeCode40 PayGradeCode = "40"
	// PayGradeCode41 identifies level 41 in a Civil Service pay plan
	PayGradeCode41 PayGradeCode = "41"
	// PayGradeCode44 identifies level 44 in a Civil Service pay plan
	PayGradeCode44 PayGradeCode = "44"
	// PayGradeCode45 identifies level 45 in a Civil Service pay plan
	PayGradeCode45 PayGradeCode = "45"
	// PayGradeCode47 identifies level 47 in a Civil Service pay plan
	PayGradeCode47 PayGradeCode = "47"
	// PayGradeCode48 identifies level 48 in a Civil Service pay plan
	PayGradeCode48 PayGradeCode = "48"
	// PayGradeCode50 identifies level 50 in a Civil Service pay plan
	PayGradeCode50 PayGradeCode = "50"
	// PayGradeCode64 identifies level 64 in a Civil Service pay plan
	PayGradeCode64 PayGradeCode = "64"
	// PayGradeCode66 identifies level 66 in a Civil Service pay plan
	PayGradeCode66 PayGradeCode = "66"
	// PayGradeCode78 identifies level 78 in a Civil Service pay plan
	PayGradeCode78 PayGradeCode = "78"
	// PayGradeCode79 identifies level 79 in a Civil Service pay plan
	PayGradeCode79 PayGradeCode = "79"
	// PayGradeCode80 identifies level 80 in a Civil Service pay plan
	PayGradeCode80 PayGradeCode = "80"
	// PayGradeCode81 identifies level 81 in a Civil Service pay plan
	PayGradeCode81 PayGradeCode = "81"
	// PayGradeCode82 identifies level 82 in a Civil Service pay plan
	PayGradeCode82 PayGradeCode = "82"
	// PayGradeCode99 identifies level 99 in a Civil Service pay plan
	PayGradeCode99 PayGradeCode = "99"
	// PayGradeCodeAA identifies level AA in a Civil Service pay plan
	PayGradeCodeAA PayGradeCode = "AA"
	// PayGradeCodeBA identifies level BA in a Civil Service pay plan
	PayGradeCodeBA PayGradeCode = "BA"
	// PayGradeCodeCB identifies level CB in a Civil Service pay plan
	PayGradeCodeCB PayGradeCode = "CB"
	// PayGradeCodeCC identifies level CC in a Civil Service pay plan
	PayGradeCodeCC PayGradeCode = "CC"
	// PayGradeCodeCD identifies level CD in a Civil Service pay plan
	PayGradeCodeCD PayGradeCode = "CD"
	// PayGradeCodeCE identifies level CE in a Civil Service pay plan
	PayGradeCodeCE PayGradeCode = "CE"
	// PayGradeCodeCG identifies level CG in a Civil Service pay plan
	PayGradeCodeCG PayGradeCode = "CG"
	// PayGradeCodeCL identifies level CL in a Civil Service pay plan
	PayGradeCodeCL PayGradeCode = "CL"
	// PayGradeCodeCM identifies level CM in a Civil Service pay plan
	PayGradeCodeCM PayGradeCode = "CM"
	// PayGradeCodeDD identifies level DD in a Civil Service pay plan
	PayGradeCodeDD PayGradeCode = "DD"
	// PayGradeCodeDE identifies level DE in a Civil Service pay plan
	PayGradeCodeDE PayGradeCode = "DE"
	// PayGradeCodeDG identifies level DG in a Civil Service pay plan
	PayGradeCodeDG PayGradeCode = "DG"
	// PayGradeCodeED identifies level ED in a Civil Service pay plan
	PayGradeCodeED PayGradeCode = "ED"
	// PayGradeCodeEE identifies level EE in a Civil Service pay plan
	PayGradeCodeEE PayGradeCode = "EE"
	// PayGradeCodeEG identifies level EG in a Civil Service pay plan
	PayGradeCodeEG PayGradeCode = "EG"
	// PayGradeCodeEM identifies level EM in a Civil Service pay plan
	PayGradeCodeEM PayGradeCode = "EM"
	// PayGradeCodeFD identifies level FD in a Civil Service pay plan
	PayGradeCodeFD PayGradeCode = "FD"
	// PayGradeCodeFE identifies level FE in a Civil Service pay plan
	PayGradeCodeFE PayGradeCode = "FE"
	// PayGradeCodeFG identifies level FG in a Civil Service pay plan
	PayGradeCodeFG PayGradeCode = "FG"
	// PayGradeCodeFM identifies level FM in a Civil Service pay plan
	PayGradeCodeFM PayGradeCode = "FM"
	// PayGradeCodeGB identifies level GB in a Civil Service pay plan
	PayGradeCodeGB PayGradeCode = "GB"
	// PayGradeCodeGC identifies level GC in a Civil Service pay plan
	PayGradeCodeGC PayGradeCode = "GC"
	// PayGradeCodeGD identifies level GD in a Civil Service pay plan
	PayGradeCodeGD PayGradeCode = "GD"
	// PayGradeCodeGE identifies level GE in a Civil Service pay plan
	PayGradeCodeGE PayGradeCode = "GE"
	// PayGradeCodeGF identifies level GF in a Civil Service pay plan
	PayGradeCodeGF PayGradeCode = "GF"
	// PayGradeCodeGG identifies level GG in a Civil Service pay plan
	PayGradeCodeGG PayGradeCode = "GG"
	// PayGradeCodeKD identifies level KD in a Civil Service pay plan
	PayGradeCodeKD PayGradeCode = "KD"
	// PayGradeCodeKE identifies level KE in a Civil Service pay plan
	PayGradeCodeKE PayGradeCode = "KE"
	// PayGradeCodeKG identifies level KG in a Civil Service pay plan
	PayGradeCodeKG PayGradeCode = "KG"
	// PayGradeCodeLD identifies level LD in a Civil Service pay plan
	PayGradeCodeLD PayGradeCode = "LD"
	// PayGradeCodeLE identifies level LE in a Civil Service pay plan
	PayGradeCodeLE PayGradeCode = "LE"
	// PayGradeCodeLG identifies level LG in a Civil Service pay plan
	PayGradeCodeLG PayGradeCode = "LG"
	// PayGradeCodeMC identifies level MC in a Civil Service pay plan
	PayGradeCodeMC PayGradeCode = "MC"
	// PayGradeCodeMD identifies level MD in a Civil Service pay plan
	PayGradeCodeMD PayGradeCode = "MD"
	// PayGradeCodeME identifies level ME in a Civil Service pay plan
	PayGradeCodeME PayGradeCode = "ME"
	// PayGradeCodeMG identifies level MG in a Civil Service pay plan
	PayGradeCodeMG PayGradeCode = "MG"
	// PayGradeCodeMM identifies level MM in a Civil Service pay plan
	PayGradeCodeMM PayGradeCode = "MM"
	// PayGradeCodeNC identifies level NC in a Civil Service pay plan
	PayGradeCodeNC PayGradeCode = "NC"
	// PayGradeCodeND identifies level ND in a Civil Service pay plan
	PayGradeCodeND PayGradeCode = "ND"
	// PayGradeCodeNE identifies level NE in a Civil Service pay plan
	PayGradeCodeNE PayGradeCode = "NE"
	// PayGradeCodeNG identifies level NG in a Civil Service pay plan
	PayGradeCodeNG PayGradeCode = "NG"
	// PayGradeCodeNM identifies level NM in a Civil Service pay plan
	PayGradeCodeNM PayGradeCode = "NM"
	// PayGradeCodeNonSupervisory is the non-supervisory paygrade code, used for DoD/non-DoD contractors when the Pay Plan Code is "ZZ"
	PayGradeCodeNonSupervisory PayGradeCode = "NS"
	// PayGradeCodeOC identifies level OC in a Civil Service pay plan
	PayGradeCodeOC PayGradeCode = "OC"
	// PayGradeCodeOD identifies level OD in a Civil Service pay plan
	PayGradeCodeOD PayGradeCode = "OD"
	// PayGradeCodeOE identifies level OE in a Civil Service pay plan
	PayGradeCodeOE PayGradeCode = "OE"
	// PayGradeCodeOG identifies level OG in a Civil Service pay plan
	PayGradeCodeOG PayGradeCode = "OG"
	// PayGradeCodeQD identifies level QD in a Civil Service pay plan
	PayGradeCodeQD PayGradeCode = "QD"
	// PayGradeCodeRH identifies level RH in a Civil Service pay plan
	PayGradeCodeRH PayGradeCode = "RH"
	// PayGradeCodeRI identifies level RI in a Civil Service pay plan
	PayGradeCodeRI PayGradeCode = "RI"
	// PayGradeCodeRJ identifies level RJ in a Civil Service pay plan
	PayGradeCodeRJ PayGradeCode = "RJ"
	// PayGradeCodeRK identifies level RK in a Civil Service pay plan
	PayGradeCodeRK PayGradeCode = "RK"
	// PayGradeCodeSH identifies level SH in a Civil Service pay plan
	PayGradeCodeSH PayGradeCode = "SH"
	// PayGradeCodeSI identifies level SI in a Civil Service pay plan
	PayGradeCodeSI PayGradeCode = "SI"
	// PayGradeCodeSJ identifies level SJ in a Civil Service pay plan
	PayGradeCodeSJ PayGradeCode = "SJ"
	// PayGradeCodeSK identifies level SK in a Civil Service pay plan
	PayGradeCodeSK PayGradeCode = "SK"
	// PayGradeCodeSupervisory is the supervisory paygrade code, used for DoD/non-DoD contractors when the Pay Plan Code is "ZZ"
	PayGradeCodeSupervisory PayGradeCode = "SP"
	// PayGradeCodeTH identifies level TH in a Civil Service pay plan
	PayGradeCodeTH PayGradeCode = "TH"
	// PayGradeCodeTI identifies level TI in a Civil Service pay plan
	PayGradeCodeTI PayGradeCode = "TI"
	// PayGradeCodeTJ identifies level TJ in a Civil Service pay plan
	PayGradeCodeTJ PayGradeCode = "TJ"
	// PayGradeCodeUH identifies level UH in a Civil Service pay plan
	PayGradeCodeUH PayGradeCode = "UH"
	// PayGradeCodeUI identifies level UI in a Civil Service pay plan
	PayGradeCodeUI PayGradeCode = "UI"
	// PayGradeCodeUJ identifies level UJ in a Civil Service pay plan
	PayGradeCodeUJ PayGradeCode = "UJ"
	// PayGradeCodeVH identifies level VH in a Civil Service pay plan
	PayGradeCodeVH PayGradeCode = "VH"
	// PayGradeCodeVI identifies level VI in a Civil Service pay plan
	PayGradeCodeVI PayGradeCode = "VI"
	// PayGradeCodeVJ identifies level VJ in a Civil Service pay plan
	PayGradeCodeVJ PayGradeCode = "VJ"
	// PayGradeCodeWE identifies level WE in a Civil Service pay plan
	PayGradeCodeWE PayGradeCode = "WE"
	// PayGradeCodeWF identifies level WF in a Civil Service pay plan
	PayGradeCodeWF PayGradeCode = "WF"
	// PayGradeCodeWG identifies level WG in a Civil Service pay plan
	PayGradeCodeWG PayGradeCode = "WG"
	// PayGradeCodeWW means unknown paygrade
	PayGradeCodeWW PayGradeCode = "WW"
	// PayGradeCodeXF identifies level XF in a Civil Service pay plan
	PayGradeCodeXF PayGradeCode = "XF"
	// PayGradeCodeXG identifies level XG in a Civil Service pay plan
	PayGradeCodeXG PayGradeCode = "XG"
	// PayGradeCodeYF identifies level YF in a Civil Service pay plan
	PayGradeCodeYF PayGradeCode = "YF"
	// PayGradeCodeYG identifies level YG in a Civil Service pay plan
	PayGradeCodeYG PayGradeCode = "YG"
	// PayGradeCodeZG identifies level ZG in a Civil Service pay plan
	PayGradeCodeZG PayGradeCode = "ZG"
)

// PayPlanCode represents the type of pay category. (The combination of pay plan code and pay grade code represents the sponsor's pay category.)
type PayPlanCode string

const (
	// PayPlanCode999 means Other Civilian Pay Plan
	PayPlanCode999 PayPlanCode = "999"
	// PayPlanCode99999 means Other Civilian Pay Plan
	PayPlanCode99999 PayPlanCode = "99999"
	// PayPlanCodeAA == Administrative Appeals Judges
	PayPlanCodeAA PayPlanCode = "AA"
	// PayPlanCodeAD == Administratively determined not elsewhere specified.
	PayPlanCodeAD PayPlanCode = "AD"
	// PayPlanCodeAF == American Family Members
	PayPlanCodeAF PayPlanCode = "AF"
	// PayPlanCodeAJ == Administrative judges, Nuclear Regulatory Commission
	PayPlanCodeAJ PayPlanCode = "AJ"
	// PayPlanCodeAL == Administrative Law judges
	PayPlanCodeAL PayPlanCode = "AL"
	// PayPlanCodeAS == Non-appropriated fund, administrative support (to be replaced by NF)
	PayPlanCodeAS PayPlanCode = "AS"
	// PayPlanCodeBB == Non supervisory negotiated pay employees
	PayPlanCodeBB PayPlanCode = "BB"
	// PayPlanCodeBL == Leader negotiated pay employees
	PayPlanCodeBL PayPlanCode = "BL"
	// PayPlanCodeBP == Printing and Lithographic negotiated pay employees
	PayPlanCodeBP PayPlanCode = "BP"
	// PayPlanCodeBS == Supervisory negotiated pay employees
	PayPlanCodeBS PayPlanCode = "BS"
	// PayPlanCodeCA == Board of contract appeals
	PayPlanCodeCA PayPlanCode = "CA"
	// PayPlanCodeCC == Commissioned Corps of Public Health Service
	PayPlanCodeCC PayPlanCode = "CC"
	// PayPlanCodeCE == Contract education
	PayPlanCodeCE PayPlanCode = "CE"
	// PayPlanCodeCG == Corporate graded Federal Deposit Insurance Corp.
	PayPlanCodeCG PayPlanCode = "CG"
	// PayPlanCodeCH == Non-appropriated fund, childcare
	PayPlanCodeCH PayPlanCode = "CH"
	// PayPlanCodeCP == U.S. Capitol Police
	PayPlanCodeCP PayPlanCode = "CP"
	// PayPlanCodeCS == Communications Analyst
	PayPlanCodeCS PayPlanCode = "CS"
	// PayPlanCodeCU == Credit Union employees
	PayPlanCodeCU PayPlanCode = "CU"
	// PayPlanCodeCY == Contract education Bureau of Indian Affairs
	PayPlanCodeCY PayPlanCode = "CY"
	// PayPlanCodeCZ == Canal Area General Schedule type positions
	PayPlanCodeCZ PayPlanCode = "CZ"
	// PayPlanCodeDA == Demonstration Administrative (Navy)
	PayPlanCodeDA PayPlanCode = "DA"
	// PayPlanCodeDB == Demonstration Engineers and Scientists (entire DoD)
	PayPlanCodeDB PayPlanCode = "DB"
	// PayPlanCodeDC == Navy Test Program - Clerical
	PayPlanCodeDC PayPlanCode = "DC"
	// PayPlanCodeDE == Demonstration Engineers and Scientists Technicians (entire DoD)
	PayPlanCodeDE PayPlanCode = "DE"
	// PayPlanCodeDG == Demonstration General (Navy)
	PayPlanCodeDG PayPlanCode = "DG"
	// PayPlanCodeDH == Demonstration hourly Air Force logistics command
	PayPlanCodeDH PayPlanCode = "DH"
	// PayPlanCodeDJ == Demonstration Administrative (entire DoD)
	PayPlanCodeDJ PayPlanCode = "DJ"
	// PayPlanCodeDK == Demonstration General Support (entire DoD)
	PayPlanCodeDK PayPlanCode = "DK"
	// PayPlanCodeDN == Defense Nuclear Facilities Safety Board Excepted Service Employees
	PayPlanCodeDN PayPlanCode = "DN"
	// PayPlanCodeDO == Business Management and Professional Career Path, Air Force Research Laboratory. Code is for use by the Department of the Air Force only.
	PayPlanCodeDO PayPlanCode = "DO"
	// PayPlanCodeDP == Demonstration Professional (Navy)
	PayPlanCodeDP PayPlanCode = "DP"
	// PayPlanCodeDQ == Demonstration Artisan Leader (DoD)
	PayPlanCodeDQ PayPlanCode = "DQ"
	// PayPlanCodeDR == Demonstration Air Force Scientist and Engineer
	PayPlanCodeDR PayPlanCode = "DR"
	// PayPlanCodeDS == Demonstration Specialist (Navy)
	PayPlanCodeDS PayPlanCode = "DS"
	// PayPlanCodeDT == Demonstration Technician (Navy)
	PayPlanCodeDT PayPlanCode = "DT"
	// PayPlanCodeDU == Mission Support Career Path, Air Force Research Laboratory. Code is for use by the Department of the Air Force only.
	PayPlanCodeDU PayPlanCode = "DU"
	// PayPlanCodeDV == Demonstration Artisan (DoD)
	PayPlanCodeDV PayPlanCode = "DV"
	// PayPlanCodeDW == Demonstration salaried Air Force and DLA
	PayPlanCodeDW PayPlanCode = "DW"
	// PayPlanCodeDX == Demonstration Supervisory Air Force and DLA
	PayPlanCodeDX PayPlanCode = "DX"
	// PayPlanCodeDZ == Demonstration Artisan (DoD)
	PayPlanCodeDZ PayPlanCode = "DZ"
	// PayPlanCodeEA == Administrative schedule (excluded) Tennessee Valley Authority
	PayPlanCodeEA PayPlanCode = "EA"
	// PayPlanCodeEB == Clerical schedule (excluded) Tennessee Valley Authority
	PayPlanCodeEB PayPlanCode = "EB"
	// PayPlanCodeEC == Engineering and Computing schedule (excluded) Tennessee Valley Authority
	PayPlanCodeEC PayPlanCode = "EC"
	// PayPlanCodeED == Expert
	PayPlanCodeED PayPlanCode = "ED"
	// PayPlanCodeEE == Expert (other)
	PayPlanCodeEE PayPlanCode = "EE"
	// PayPlanCodeEF == Consultant
	PayPlanCodeEF PayPlanCode = "EF"
	// PayPlanCodeEG == Consultant (other)
	PayPlanCodeEG PayPlanCode = "EG"
	// PayPlanCodeEH == Advisory committee member
	PayPlanCodeEH PayPlanCode = "EH"
	// PayPlanCodeEI == Advisory committee member (other)
	PayPlanCodeEI PayPlanCode = "EI"
	// PayPlanCodeEM == Executive schedule Office of the Comptroller of the currency
	PayPlanCodeEM PayPlanCode = "EM"
	// PayPlanCodeEO == FDIC executive pay
	PayPlanCodeEO PayPlanCode = "EO"
	// PayPlanCodeEP == Defense Intelligence Senior Executive Service
	PayPlanCodeEP PayPlanCode = "EP"
	// PayPlanCodeES == Senior Executive Service (SES)
	PayPlanCodeES PayPlanCode = "ES"
	// PayPlanCodeET == General Accounting Office Senior Executive Service
	PayPlanCodeET PayPlanCode = "ET"
	// PayPlanCodeEX == Executive pay
	PayPlanCodeEX PayPlanCode = "EX"
	// PayPlanCodeFA == Foreign Service Chiefs of mission
	PayPlanCodeFA PayPlanCode = "FA"
	// PayPlanCodeFC == Foreign compensation Agency for International Development
	PayPlanCodeFC PayPlanCode = "FC"
	// PayPlanCodeFD == Foreign defense
	PayPlanCodeFD PayPlanCode = "FD"
	// PayPlanCodeFE == Senior Foreign Service
	PayPlanCodeFE PayPlanCode = "FE"
	// PayPlanCodeFH == Members of the Foreign Service employed by the Department of State
	PayPlanCodeFH PayPlanCode = "FH"
	// PayPlanCodeFO == Foreign Service Officers
	PayPlanCodeFO PayPlanCode = "FO"
	// PayPlanCodeFP == Foreign Service personnel
	PayPlanCodeFP PayPlanCode = "FP"
	// PayPlanCodeFZ == Consular Agent Department of State
	PayPlanCodeFZ PayPlanCode = "FZ"
	// PayPlanCodeGD == Skill based pay demonstration project managers (DLA)
	PayPlanCodeGD PayPlanCode = "GD"
	// PayPlanCodeGG == Grades similar to General Schedule
	PayPlanCodeGG PayPlanCode = "GG"
	// PayPlanCodeGH == GG employees converted to performance and management recognition system
	PayPlanCodeGH PayPlanCode = "GH"
	// PayPlanCodeGL == GS Law Enforcement Officers
	PayPlanCodeGL PayPlanCode = "GL"
	// PayPlanCodeGM == Performance Management and Recognition system
	PayPlanCodeGM PayPlanCode = "GM"
	// PayPlanCodeGN == Nurse at Warren G. Magnuson Clinical Center
	PayPlanCodeGN PayPlanCode = "GN"
	// PayPlanCodeGP == GS Physicians and Dentists
	PayPlanCodeGP PayPlanCode = "GP"
	// PayPlanCodeGR == GM Physicians and Dentists
	PayPlanCodeGR PayPlanCode = "GR"
	// PayPlanCodeGS == General Schedule
	PayPlanCodeGS PayPlanCode = "GS"
	// PayPlanCodeGW == Employment under schedule A paid at GS rate Stay-In-School program
	PayPlanCodeGW PayPlanCode = "GW"
	// PayPlanCodeIA == Defense Civilian Intelligence Personnel System pay-banded compensation structure. Code is for use by the Department of Defense only.
	PayPlanCodeIA PayPlanCode = "IA"
	// PayPlanCodeIE == Senior Intelligence Executive Service (SIES) Program
	PayPlanCodeIE PayPlanCode = "IE"
	// PayPlanCodeIJ == Immigration Judge Schedule
	PayPlanCodeIJ PayPlanCode = "IJ"
	// PayPlanCodeIP == Senior Intelligence Professional Program
	PayPlanCodeIP PayPlanCode = "IP"
	// PayPlanCodeJG == Graded tradesmen and craftsmen United States Courts
	PayPlanCodeJG PayPlanCode = "JG"
	// PayPlanCodeJL == Leaders of tradesmen and craftsmen United States Courts
	PayPlanCodeJL PayPlanCode = "JL"
	// PayPlanCodeJP == Non supervisory lithographers and printers United States Courts
	PayPlanCodeJP PayPlanCode = "JP"
	// PayPlanCodeJQ == Lead lithographers and printers United States Courts
	PayPlanCodeJQ PayPlanCode = "JQ"
	// PayPlanCodeJR == Supervisory lithographers and printers United States Courts
	PayPlanCodeJR PayPlanCode = "JR"
	// PayPlanCodeJT == Supervisors for tradesmen and craftsmen United States Courts
	PayPlanCodeJT PayPlanCode = "JT"
	// PayPlanCodeKA == Kleas Act Government Printing Office
	PayPlanCodeKA PayPlanCode = "KA"
	// PayPlanCodeKG == Non-Craft non supervisory Bureau of Engraving and Printing
	PayPlanCodeKG PayPlanCode = "KG"
	// PayPlanCodeKL == Non-Craft leader Bureau of Engraving and Printing
	PayPlanCodeKL PayPlanCode = "KL"
	// PayPlanCodeKS == Non-Craft supervisory Bureau of Engraving and Printing
	PayPlanCodeKS PayPlanCode = "KS"
	// PayPlanCodeLE == United States Secret Service uniformed division Treasury
	PayPlanCodeLE PayPlanCode = "LE"
	// PayPlanCodeLG == Liquidation graded FDIC
	PayPlanCodeLG PayPlanCode = "LG"
	// PayPlanCodeLX == Senior-level Excepted Service Position (GAO)
	PayPlanCodeLX PayPlanCode = "LX"
	// PayPlanCodeMA == Milk Marketing Department of Agriculture
	PayPlanCodeMA PayPlanCode = "MA"
	// PayPlanCodeMC == Cadet (uniformed service only)
	PayPlanCodeMC PayPlanCode = "MC"
	// PayPlanCodeME == Enlisted (uniformed service only)
	PayPlanCodeME PayPlanCode = "ME"
	// PayPlanCodeMO == Officer (uniformed service only)
	PayPlanCodeMO PayPlanCode = "MO"
	// PayPlanCodeMW == Warrant Officer (uniformed service only)
	PayPlanCodeMW PayPlanCode = "MW"
	// PayPlanCodeNA == Non appropriated funds, non supervisory, non leader Federal Wage System
	PayPlanCodeNA PayPlanCode = "NA"
	// PayPlanCodeNC == Naval Research Lab Administrative Support
	PayPlanCodeNC PayPlanCode = "NC"
	// PayPlanCodeND == Demonstration Scientific and Engineering (Navy Only)
	PayPlanCodeND PayPlanCode = "ND"
	// PayPlanCodeNF == Non-appropriated fund, pay band
	PayPlanCodeNF PayPlanCode = "NF"
	// PayPlanCodeNG == Demonstration General Support (Navy Only)
	PayPlanCodeNG PayPlanCode = "NG"
	// PayPlanCodeNH == Business Management and Technical Management Professional. DOD Acquisition Workforce Demonstration Project (entire DoD)
	PayPlanCodeNH PayPlanCode = "NH"
	// PayPlanCodeNJ == Technical Management Support, DOD Acquisition Workforce
	PayPlanCodeNJ PayPlanCode = "NJ"
	// PayPlanCodeNK == Administration Support, DOD Acquisition Workforce Demonstration Project (entire DoD)
	PayPlanCodeNK PayPlanCode = "NK"
	// PayPlanCodeNL == Non-appropriated fund, crafts and trades worker
	PayPlanCodeNL PayPlanCode = "NL"
	// PayPlanCodeNM == Supervisors and Managers. Code is for use by the Department of the Navy only for the Naval Research Laboratory and similar pay demonstration projects.
	PayPlanCodeNM PayPlanCode = "NM"
	// PayPlanCodeNO == Administrative Specialist/Professional
	PayPlanCodeNO PayPlanCode = "NO"
	// PayPlanCodeNP == Science and Engineering Professional
	PayPlanCodeNP PayPlanCode = "NP"
	// PayPlanCodeNR == Science and Engineering Technical
	PayPlanCodeNR PayPlanCode = "NR"
	// PayPlanCodeNS == Non appropriated funds, supervisory, Federal Wage System
	PayPlanCodeNS PayPlanCode = "NS"
	// PayPlanCodeNT == Demonstration Administrative and Technical (Navy Only)
	PayPlanCodeNT PayPlanCode = "NT"
	// PayPlanCodeOC == Office of the Comptroller of the Currency
	PayPlanCodeOC PayPlanCode = "OC"
	// PayPlanCodePA == Attorneys and law clerks General Accounting Office
	PayPlanCodePA PayPlanCode = "PA"
	// PayPlanCodePE == Evaluator and evaluator related General Accounting Office
	PayPlanCodePE PayPlanCode = "PE"
	// PayPlanCodePG == Printing Office grades
	PayPlanCodePG PayPlanCode = "PG"
	// PayPlanCodePS == Non-appropriated fund, patron service (to be replaced by NF)
	PayPlanCodePS PayPlanCode = "PS"
	// PayPlanCodeRS == Senior Biomedical Service
	PayPlanCodeRS PayPlanCode = "RS"
	// PayPlanCodeSA == Administrative schedule Tennessee Valley Authority
	PayPlanCodeSA PayPlanCode = "SA"
	// PayPlanCodeSB == Clerical schedule (excluded) Tennessee Valley Authority
	PayPlanCodeSB PayPlanCode = "SB"
	// PayPlanCodeSC == Engineering and Computing schedule Tennessee Valley Authority
	PayPlanCodeSC PayPlanCode = "SC"
	// PayPlanCodeSD == Scientific and Programming schedule Tennessee Valley Authority
	PayPlanCodeSD PayPlanCode = "SD"
	// PayPlanCodeSE == Aide and Technician schedule Tennessee Valley Authority
	PayPlanCodeSE PayPlanCode = "SE"
	// PayPlanCodeSF == Custodial schedule Tennessee Valley Authority
	PayPlanCodeSF PayPlanCode = "SF"
	// PayPlanCodeSG == Public Safety schedule Tennessee Valley Authority
	PayPlanCodeSG PayPlanCode = "SG"
	// PayPlanCodeSH == Physicians schedule Tennessee Valley Authority
	PayPlanCodeSH PayPlanCode = "SH"
	// PayPlanCodeSJ == Scientific and Programming schedule (excluded) Tennessee Valley Authority
	PayPlanCodeSJ PayPlanCode = "SJ"
	// PayPlanCodeSL == Senior Level Positions
	PayPlanCodeSL PayPlanCode = "SL"
	// PayPlanCodeSM == Management Schedule Tennessee Valley Authority
	PayPlanCodeSM PayPlanCode = "SM"
	// PayPlanCodeSN == Senior Level System Nuclear Regulatory Commission
	PayPlanCodeSN PayPlanCode = "SN"
	// PayPlanCodeSP == Park Police Department of the Interior
	PayPlanCodeSP PayPlanCode = "SP"
	// PayPlanCodeSQ == Physicians and dentists paid under Scientific and Professional (ST) pay
	PayPlanCodeSQ PayPlanCode = "SQ"
	// PayPlanCodeSR == Statutory rates not elsewhere specified
	PayPlanCodeSR PayPlanCode = "SR"
	// PayPlanCodeSS == Senior Staff positions
	PayPlanCodeSS PayPlanCode = "SS"
	// PayPlanCodeST == Scientific and professional
	PayPlanCodeST PayPlanCode = "ST"
	// PayPlanCodeSZ == Canal Area Special category type positions
	PayPlanCodeSZ PayPlanCode = "SZ"
	// PayPlanCodeTA == Construction schedule
	PayPlanCodeTA PayPlanCode = "TA"
	// PayPlanCodeTB == Operating and Maintenance (power facilities) Tennessee Valley Authority
	PayPlanCodeTB PayPlanCode = "TB"
	// PayPlanCodeTC == Chemical Operators Tennessee Valley Authority
	PayPlanCodeTC PayPlanCode = "TC"
	// PayPlanCodeTD == Plant Operators schedule Tennessee Valley Authority
	PayPlanCodeTD PayPlanCode = "TD"
	// PayPlanCodeTE == Operating and Maintenance (nonpower facilities) Tennessee Valley Authority
	PayPlanCodeTE PayPlanCode = "TE"
	// PayPlanCodeTM == Federal Housing Finance board Executive level
	PayPlanCodeTM PayPlanCode = "TM"
	// PayPlanCodeTP == Teaching positions DoD schools only
	PayPlanCodeTP PayPlanCode = "TP"
	// PayPlanCodeTR == Police Forces US Mint and Bureau of Engraving and Printing
	PayPlanCodeTR PayPlanCode = "TR"
	// PayPlanCodeTS == Step System Federal Housing Finance board
	PayPlanCodeTS PayPlanCode = "TS"
	// PayPlanCodeVC == Canteen Service Department of Veterans Affairs
	PayPlanCodeVC PayPlanCode = "VC"
	// PayPlanCodeVE == Canteen Service Executives Department of Veterans Affairs
	PayPlanCodeVE PayPlanCode = "VE"
	// PayPlanCodeVG == Clerical and Administrative support Farm Credit
	PayPlanCodeVG PayPlanCode = "VG"
	// PayPlanCodeVH == Professional, Administrative, and Managerial Farm Credit
	PayPlanCodeVH PayPlanCode = "VH"
	// PayPlanCodeVM == Medical and Dental Department of Veterans Affairs
	PayPlanCodeVM PayPlanCode = "VM"
	// PayPlanCodeVN == Nurses Department of Veterans Affairs
	PayPlanCodeVN PayPlanCode = "VN"
	// PayPlanCodeVP == Clinical Podiatrists and Optometrists Department of Veterans Affairs
	PayPlanCodeVP PayPlanCode = "VP"
	// PayPlanCodeWA == Navigation Lock and Dam Operation and maintenance supervisory USACE
	PayPlanCodeWA PayPlanCode = "WA"
	// PayPlanCodeWB == Wage positions under Federal Wage System otherwise not designated (obsolete)
	PayPlanCodeWB PayPlanCode = "WB"
	// PayPlanCodeWD == Production facilitating non supervisory Federal Wage System (obsolete)
	PayPlanCodeWD PayPlanCode = "WD"
	// PayPlanCodeWE == Currency manufacturing Department of the Treasury
	PayPlanCodeWE PayPlanCode = "WE"
	// PayPlanCodeWF == Motion Picture Production
	PayPlanCodeWF PayPlanCode = "WF"
	// PayPlanCodeWG == Non supervisory pay schedule Federal Wage System (obsolete)
	PayPlanCodeWG PayPlanCode = "WG"
	// PayPlanCodeWI == Printing and Lithographic (D.C.)
	PayPlanCodeWI PayPlanCode = "WI"
	// PayPlanCodeWJ == Hopper Dredge Schedule Supervisory Federal Wage System Dept of Army
	PayPlanCodeWJ PayPlanCode = "WJ"
	// PayPlanCodeWK == Hopper Dredge Schedule non supervisory Federal Wage System Dept of Army
	PayPlanCodeWK PayPlanCode = "WK"
	// PayPlanCodeWL == Leader pay schedules Federal Wage System
	PayPlanCodeWL PayPlanCode = "WL"
	// PayPlanCodeWM == Maritime pay schedules
	PayPlanCodeWM PayPlanCode = "WM"
	// PayPlanCodeWN == Production facilitating supervisory Federal Wage System
	PayPlanCodeWN PayPlanCode = "WN"
	// PayPlanCodeWO == Navigation Lock and Dam Operation and maintenance leader USACE
	PayPlanCodeWO PayPlanCode = "WO"
	// PayPlanCodeWP == Printing and Lithographic (other than D.C.)
	PayPlanCodeWP PayPlanCode = "WP"
	// PayPlanCodeWQ == Aircraft Electronic Equipment and Optical Inst. repair supervisory
	PayPlanCodeWQ PayPlanCode = "WQ"
	// PayPlanCodeWR == Aircraft Electronic Equipment and Optical Inst. repair leader
	PayPlanCodeWR PayPlanCode = "WR"
	// PayPlanCodeWS == Supervisor Federal Wage System (obsolete)
	PayPlanCodeWS PayPlanCode = "WS"
	// PayPlanCodeWT == Apprentices and Shop trainees Federal Wage System (obsolete)
	PayPlanCodeWT PayPlanCode = "WT"
	// PayPlanCodeWU == Aircraft Electronic Equipment and Optical Inst. repair non supervisory
	PayPlanCodeWU PayPlanCode = "WU"
	// PayPlanCodeWW == Wage type excepted Stay-In-School Federal Wage System (obsolete)
	PayPlanCodeWW PayPlanCode = "WW"
	// PayPlanCodeWY == Navigation Lock and Dam Operation and maintenance non supervisory USACE
	PayPlanCodeWY PayPlanCode = "WY"
	// PayPlanCodeWZ == Canal Area Wage System type positions (obsolete)
	PayPlanCodeWZ PayPlanCode = "WZ"
	// PayPlanCodeXA == Special Overlap Area Rate Schedule non supervisory Dept of the Interior
	PayPlanCodeXA PayPlanCode = "XA"
	// PayPlanCodeXB == Special Overlap Area Rate Schedule leader Dept of the Interior
	PayPlanCodeXB PayPlanCode = "XB"
	// PayPlanCodeXC == Special Overlap Area Rate Schedule supervisory Dept of the Interior
	PayPlanCodeXC PayPlanCode = "XC"
	// PayPlanCodeXD == Non supervisory production facilitating special schedule printing employees
	PayPlanCodeXD PayPlanCode = "XD"
	// PayPlanCodeXF == Floating Plant Schedule non supervisory Dept of Army
	PayPlanCodeXF PayPlanCode = "XF"
	// PayPlanCodeXG == Floating Plant Schedule leader Dept of Army
	PayPlanCodeXG PayPlanCode = "XG"
	// PayPlanCodeXH == Floating Plant Schedule supervisory Dept of Army
	PayPlanCodeXH PayPlanCode = "XH"
	// PayPlanCodeXL == Leader special schedule printing employees
	PayPlanCodeXL PayPlanCode = "XL"
	// PayPlanCodeXN == Supervisory production facilitating special schedule printing employees
	PayPlanCodeXN PayPlanCode = "XN"
	// PayPlanCodeXP == Non supervisory special schedule printing employees
	PayPlanCodeXP PayPlanCode = "XP"
	// PayPlanCodeXS == Supervisory special schedule printing employees
	PayPlanCodeXS PayPlanCode = "XS"
	// PayPlanCodeXW == Automotive Mechanic, Non-Supervisory (NAF)
	PayPlanCodeXW PayPlanCode = "XW"
	// PayPlanCodeXY == Automotive Mechanic, Leader (NAF)
	PayPlanCodeXY PayPlanCode = "XY"
	// PayPlanCodeXZ == Automotive Mechanic, Supervisory NAF)
	PayPlanCodeXZ PayPlanCode = "XZ"
	// PayPlanCodeYA == NSPS Standard Career Group - Professional/Analytical Pay Schedule
	PayPlanCodeYA PayPlanCode = "YA"
	// PayPlanCodeYB == NSPS Standard Career Group - Technician/Support Pay Schedule
	PayPlanCodeYB PayPlanCode = "YB"
	// PayPlanCodeYC == NSPS Standard Career Group - Supervisor/Manager Pay Schedule
	PayPlanCodeYC PayPlanCode = "YC"
	// PayPlanCodeYD == NSPS Scientific and Engineering Career Group - Professional Pay Schedule
	PayPlanCodeYD PayPlanCode = "YD"
	// PayPlanCodeYE == NSPS Scientific and Engineering Career Group - Technician/Support Pay Schedule
	PayPlanCodeYE PayPlanCode = "YE"
	// PayPlanCodeYF == NSPS Scientific and Engineering Career Group - Supervisor/Manager Pay Schedule
	PayPlanCodeYF PayPlanCode = "YF"
	// PayPlanCodeYG == NSPS Medical Career Group - Physician/Dentist Pay Schedule
	PayPlanCodeYG PayPlanCode = "YG"
	// PayPlanCodeYH == NSPS Medical Career Group - Professional Pay Schedule
	PayPlanCodeYH PayPlanCode = "YH"
	// PayPlanCodeYI == NSPS Medical Career Group - Technician/Support Pay Schedule
	PayPlanCodeYI PayPlanCode = "YI"
	// PayPlanCodeYJ == NSPS Medical Career Group - Supervisor/Manager Pay Schedule
	PayPlanCodeYJ PayPlanCode = "YJ"
	// PayPlanCodeYK == NSPS Investigative and Protective Career Group - Investigative Pay Schedule
	PayPlanCodeYK PayPlanCode = "YK"
	// PayPlanCodeYL == NSPS Investigative and Protective Career Group - Fire Protection Pay Schedule
	PayPlanCodeYL PayPlanCode = "YL"
	// PayPlanCodeYM == NSPS Investigative and Protective Career Group - Police/Security Guard Pay Schedule
	PayPlanCodeYM PayPlanCode = "YM"
	// PayPlanCodeYN == NSPS Investigative and Protective Career Group - Supervisor/Manager Pay Schedule
	PayPlanCodeYN PayPlanCode = "YN"
	// PayPlanCodeYO == Undefined DoD, DON, DOA, DAF
	PayPlanCodeYO PayPlanCode = "YO"
	// PayPlanCodeYP == NSPS Standard Career Group - Student Educational Employment Program Pay Schedule
	PayPlanCodeYP PayPlanCode = "YP"
	// PayPlanCodeYQ == Undefined DoD, DON, DOA, DAF
	PayPlanCodeYQ PayPlanCode = "YQ"
	// PayPlanCodeYR == Undefined DoD, DON, DOA, DAF
	PayPlanCodeYR PayPlanCode = "YR"
	// PayPlanCodeYS == Undefined DoD, DON, DOA, DAF
	PayPlanCodeYS PayPlanCode = "YS"
	// PayPlanCodeYT == Undefined DoD, DON, DOA, DAF
	PayPlanCodeYT PayPlanCode = "YT"
	// PayPlanCodeYU == Undefined DoD, DON, DOA, DAF
	PayPlanCodeYU PayPlanCode = "YU"
	// PayPlanCodeYV == Temporary summer aid employment
	PayPlanCodeYV PayPlanCode = "YV"
	// PayPlanCodeYW == Student aid employment Stay-In-School
	PayPlanCodeYW PayPlanCode = "YW"
	// PayPlanCodeYX == Undefined DoD, DON, DOA, DAF
	PayPlanCodeYX PayPlanCode = "YX"
	// PayPlanCodeYY == Undefined DoD, DON, DOA, DAF
	PayPlanCodeYY PayPlanCode = "YY"
	// PayPlanCodeYZ == Undefined DoD, DON, DOA, DAF
	PayPlanCodeYZ PayPlanCode = "YZ"
	// PayPlanCodeZA == Administrative (Department of Commerce)
	PayPlanCodeZA PayPlanCode = "ZA"
	// PayPlanCodeZP == Scientific and Engineering Professional National Institute of Standards and Technology
	PayPlanCodeZP PayPlanCode = "ZP"
	// PayPlanCodeZS == Administrative Support (Department of Commerce)
	PayPlanCodeZS PayPlanCode = "ZS"
	// PayPlanCodeZT == Scientific and Engineering Technician (Department of Commerce)
	PayPlanCodeZT PayPlanCode = "ZT"
	// PayPlanCodeZZ == Not applicable. Use only with pay basis WC (without compensation) when other Pay Plan Codes are not applicable.
	PayPlanCodeZZ PayPlanCode = "ZZ"
)

// ServiceCode represents the branch classification of Service with which the sponsor is affiliated.
type ServiceCode string

const (
	// ServiceCodeForeignArmy means Foreign Army
	ServiceCodeForeignArmy ServiceCode = "1"
	// ServiceCodeForeignNavy means Foreign Navy
	ServiceCodeForeignNavy ServiceCode = "2"
	// ServiceCodeForeignMarineCorps means Foreign Marine Corps
	ServiceCodeForeignMarineCorps ServiceCode = "3"
	// ServiceCodeForeignAirForce means Foreign Air Force
	ServiceCodeForeignAirForce ServiceCode = "4"
	// ServiceCodeForeignCoastGuard means Foreign Coast Guard
	ServiceCodeForeignCoastGuard ServiceCode = "6"
	// ServiceCodeArmy means the United States Army
	ServiceCodeArmy ServiceCode = "A"
	// ServiceCodeCoastGuard means the United States Coast Guard
	ServiceCodeCoastGuard ServiceCode = "C"
	// ServiceCodeOSD means the Office of the Secretary of Defense
	ServiceCodeOSD ServiceCode = "D"
	// ServiceCodeAirForce means the United States Air Force
	ServiceCodeAirForce ServiceCode = "F"
	// ServiceCodePublicHealth means the Commissioned Corps of the Public Health Service
	ServiceCodePublicHealth ServiceCode = "H"
	// ServiceCodeMarineCorps means the United States Marine Corps
	ServiceCodeMarineCorps ServiceCode = "M"
	// ServiceCodeNavy means the United States Navy
	ServiceCodeNavy ServiceCode = "N"
	// ServiceCodeNOAA means the Commissioned Corps of the National Oceanic and Atmospheric Administration
	ServiceCodeNOAA ServiceCode = "O"
	// ServiceCodeNotApplicable means Not applicable
	ServiceCodeNotApplicable ServiceCode = "X"
	// ServiceCodeUnknown means Unknown
	ServiceCodeUnknown ServiceCode = "Z"
)
