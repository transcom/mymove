package ordersapi

import (
	"errors"
	"fmt"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/models"
)

var issuerDict = map[string]ordersmessages.Issuer{
	string(internalmessages.DeptIndicatorAIRFORCE): ordersmessages.IssuerAirForce,
	string(internalmessages.DeptIndicatorMARINES):  ordersmessages.IssuerMarineCorps,
}

func deptIndicatorToAPIIssuer(di *string) (ordersmessages.Issuer, error) {
	if di == nil {
		return "", errors.New("DeptIndicator cannot be nil")
	}
	issuer := issuerDict[*di]
	if issuer == "" {
		return "", fmt.Errorf("Unknown issuer %v", di)
	}
	return issuer, nil
}

var toAPIAffiliationMap = map[models.ServiceMemberAffiliation]ordersmessages.Affiliation{
	models.AffiliationARMY:       ordersmessages.AffiliationArmy,
	models.AffiliationNAVY:       ordersmessages.AffiliationNavy,
	models.AffiliationMARINES:    ordersmessages.AffiliationMarineCorps,
	models.AffiliationAIRFORCE:   ordersmessages.AffiliationAirForce,
	models.AffiliationCOASTGUARD: ordersmessages.AffiliationCoastGuard,
	// currently missing - models.AffiliationCIVILIANAGENCY: ordersmessages.AffiliationCivilianAgency
}

func toAPIAffiliation(sma *models.ServiceMemberAffiliation) (ordersmessages.Affiliation, error) {
	if sma == nil {
		return "", errors.New("ServiceMemberAffiliation cannot be nil")
	}
	aff := toAPIAffiliationMap[*sma]
	if aff == "" {
		return "", fmt.Errorf("Unknown affiliation %v", sma)
	}
	return aff, nil
}

var fromAPIAffiliationMap = map[ordersmessages.Affiliation]models.ServiceMemberAffiliation{
	ordersmessages.AffiliationArmy:        models.AffiliationARMY,
	ordersmessages.AffiliationNavy:        models.AffiliationNAVY,
	ordersmessages.AffiliationMarineCorps: models.AffiliationMARINES,
	ordersmessages.AffiliationAirForce:    models.AffiliationAIRFORCE,
	ordersmessages.AffiliationCoastGuard:  models.AffiliationCOASTGUARD,
	// currently missing - ordersmessages.AffiliationCivilianAgency: models.AffiliationCIVILIANAGENCY,
}

func fromAPIAffiliation(oma ordersmessages.Affiliation) (models.ServiceMemberAffiliation, error) {
	aff := fromAPIAffiliationMap[oma]
	if aff == "" {
		return "", fmt.Errorf("Unknown affiliation %v", oma)
	}
	return aff, nil
}

var toAPIRankMap = map[models.ServiceMemberRank]ordersmessages.Rank{
	models.ServiceMemberRankE1:                ordersmessages.RankE1,
	models.ServiceMemberRankE2:                ordersmessages.RankE2,
	models.ServiceMemberRankE3:                ordersmessages.RankE3,
	models.ServiceMemberRankE4:                ordersmessages.RankE4,
	models.ServiceMemberRankE5:                ordersmessages.RankE5,
	models.ServiceMemberRankE6:                ordersmessages.RankE6,
	models.ServiceMemberRankE7:                ordersmessages.RankE7,
	models.ServiceMemberRankE8:                ordersmessages.RankE8,
	models.ServiceMemberRankE9:                ordersmessages.RankE9,
	models.ServiceMemberRankO1ACADEMYGRADUATE: ordersmessages.RankO1,
	models.ServiceMemberRankO2:                ordersmessages.RankO2,
	models.ServiceMemberRankO3:                ordersmessages.RankO3,
	models.ServiceMemberRankO4:                ordersmessages.RankO4,
	models.ServiceMemberRankO5:                ordersmessages.RankO5,
	models.ServiceMemberRankO6:                ordersmessages.RankO6,
	models.ServiceMemberRankO7:                ordersmessages.RankO7,
	models.ServiceMemberRankO8:                ordersmessages.RankO8,
	models.ServiceMemberRankO9:                ordersmessages.RankO9,
	models.ServiceMemberRankO10:               ordersmessages.RankO10,
	models.ServiceMemberRankW1:                ordersmessages.RankW1,
	models.ServiceMemberRankW2:                ordersmessages.RankW2,
	models.ServiceMemberRankW3:                ordersmessages.RankW3,
	models.ServiceMemberRankW4:                ordersmessages.RankW4,
	models.ServiceMemberRankW5:                ordersmessages.RankW5,
	models.ServiceMemberRankAVIATIONCADET:     ordersmessages.RankAviationCadet,
	models.ServiceMemberRankCIVILIANEMPLOYEE:  ordersmessages.RankCivilian,
	models.ServiceMemberRankACADEMYCADET:      ordersmessages.RankCadet,
	models.ServiceMemberRankMIDSHIPMAN:        ordersmessages.RankMidshipman,
}

func toAPIRank(smr *models.ServiceMemberRank) (ordersmessages.Rank, error) {
	if smr == nil {
		return "", errors.New("ServiceMemberRank cannot be nil")
	}
	rank := toAPIRankMap[*smr]
	if rank == "" {
		return "", fmt.Errorf("Unknown rank %v", smr)
	}
	return rank, nil
}

var fromAPIRankMap = map[ordersmessages.Rank]models.ServiceMemberRank{
	ordersmessages.RankE1:            models.ServiceMemberRankE1,
	ordersmessages.RankE2:            models.ServiceMemberRankE2,
	ordersmessages.RankE3:            models.ServiceMemberRankE3,
	ordersmessages.RankE4:            models.ServiceMemberRankE4,
	ordersmessages.RankE5:            models.ServiceMemberRankE5,
	ordersmessages.RankE6:            models.ServiceMemberRankE6,
	ordersmessages.RankE7:            models.ServiceMemberRankE7,
	ordersmessages.RankE8:            models.ServiceMemberRankE8,
	ordersmessages.RankE9:            models.ServiceMemberRankE9,
	ordersmessages.RankO1:            models.ServiceMemberRankO1ACADEMYGRADUATE,
	ordersmessages.RankO2:            models.ServiceMemberRankO2,
	ordersmessages.RankO3:            models.ServiceMemberRankO3,
	ordersmessages.RankO4:            models.ServiceMemberRankO4,
	ordersmessages.RankO5:            models.ServiceMemberRankO5,
	ordersmessages.RankO6:            models.ServiceMemberRankO6,
	ordersmessages.RankO7:            models.ServiceMemberRankO7,
	ordersmessages.RankO8:            models.ServiceMemberRankO8,
	ordersmessages.RankO9:            models.ServiceMemberRankO9,
	ordersmessages.RankO10:           models.ServiceMemberRankO10,
	ordersmessages.RankW1:            models.ServiceMemberRankW1,
	ordersmessages.RankW2:            models.ServiceMemberRankW2,
	ordersmessages.RankW3:            models.ServiceMemberRankW3,
	ordersmessages.RankW4:            models.ServiceMemberRankW4,
	ordersmessages.RankW5:            models.ServiceMemberRankW5,
	ordersmessages.RankAviationCadet: models.ServiceMemberRankAVIATIONCADET,
	ordersmessages.RankCivilian:      models.ServiceMemberRankCIVILIANEMPLOYEE,
	ordersmessages.RankCadet:         models.ServiceMemberRankACADEMYCADET,
	ordersmessages.RankMidshipman:    models.ServiceMemberRankMIDSHIPMAN,
}

func fromAPIRank(omr ordersmessages.Rank) (models.ServiceMemberRank, error) {
	rank := fromAPIRankMap[omr]
	if rank == "" {
		return "", fmt.Errorf("Unknown rank %v", omr)
	}
	return rank, nil
}

func serviceMemberToAPIMember(sm models.ServiceMember) (*ordersmessages.Member, error) {
	member := ordersmessages.Member{}
	member.GivenName = sm.FirstName
	if sm.MiddleName != nil {
		member.MiddleName = *sm.MiddleName
	}
	member.FamilyName = sm.LastName
	if sm.Suffix != nil {
		member.Suffix = *sm.Suffix
	}
	var err error
	member.Affiliation, err = toAPIAffiliation(sm.Affiliation)
	if err != nil {
		return nil, err
	}
	member.Rank, err = toAPIRank(sm.Rank)
	if err != nil {
		return nil, err
	}
	if sm.Title != nil {
		member.Title = *sm.Title
	}
	return &member, nil
}

var toAPIStatusMap = map[models.OrderStatus]ordersmessages.Status{
	// OrderStatusAPPROVED and OrderStatusCANCELED are not currently used by the manual Orders entry flow.
	models.OrderStatusAPPROVED:  ordersmessages.StatusAuthorized,
	models.OrderStatusCANCELED:  ordersmessages.StatusCanceled,
	models.OrderStatusDRAFT:     ordersmessages.StatusRfo,
	models.OrderStatusSUBMITTED: ordersmessages.StatusAuthorized,
}

func toAPIStatus(os models.OrderStatus) (ordersmessages.Status, error) {
	apiStatus := toAPIStatusMap[os]
	if apiStatus == "" {
		return "", fmt.Errorf("Unknown order status %s", os)
	}
	return apiStatus, nil
}

var tourTypeDict = map[models.TourType]ordersmessages.TourType{
	models.TourTypeAccompanied:                       ordersmessages.TourTypeAccompanied,
	models.TourTypeUnaccompanied:                     ordersmessages.TourTypeUnaccompanied,
	models.TourTypeUnaccompaniedDependentsRestricted: ordersmessages.TourTypeUnaccompaniedDependentsRestricted,
}

func tourTypeToAPITourType(tt models.TourType) (ordersmessages.TourType, error) {
	apiTourType := tourTypeDict[tt]
	if apiTourType == "" {
		return "", fmt.Errorf("Unknown tour type %s", tt)
	}
	return apiTourType, nil
}

var toAPIOrdersTypeMap = map[internalmessages.OrdersType]ordersmessages.OrdersType{
	internalmessages.OrdersTypeACCESSION:           ordersmessages.OrdersTypeAccession,
	internalmessages.OrdersTypeBETWEENDUTYSTATIONS: ordersmessages.OrdersTypeBetweenDutyStations,
	internalmessages.OrdersTypeBRAC:                ordersmessages.OrdersTypeBrac,
	internalmessages.OrdersTypeCOT:                 ordersmessages.OrdersTypeCot,
	internalmessages.OrdersTypeEMERGENCYEVAC:       ordersmessages.OrdersTypeEmergencyEvac,
	internalmessages.OrdersTypeIPCOT:               ordersmessages.OrdersTypeIpcot,
	internalmessages.OrdersTypeLOWCOSTTRAVEL:       ordersmessages.OrdersTypeLowCostTravel,
	internalmessages.OrdersTypeOPERATIONAL:         ordersmessages.OrdersTypeOperational,
	internalmessages.OrdersTypeOTEIP:               ordersmessages.OrdersTypeOteip,
	internalmessages.OrdersTypeROTATIONAL:          ordersmessages.OrdersTypeRotational,
	internalmessages.OrdersTypeSEPARATION:          ordersmessages.OrdersTypeSeparation,
	internalmessages.OrdersTypeSPECIALPURPOSE:      ordersmessages.OrdersTypeSpecialPurpose,
	internalmessages.OrdersTypeTRAINING:            ordersmessages.OrdersTypeTraining,
	internalmessages.OrdersTypeUNITMOVE:            ordersmessages.OrdersTypeUnitMove,
	// FIXME: NO MATCH for internalmessages.OrdersTypePERMANENTCHANGEOFSTATION - PCS happens for many kinds of orders
}

func toAPIOrdersType(ot internalmessages.OrdersType) (ordersmessages.OrdersType, error) {
	ordersType := toAPIOrdersTypeMap[ot]
	if ordersType == "" {
		return "", fmt.Errorf("Unknown orders type %s", ot)
	}
	return ordersType, nil
}

var fromAPIOrdersTypeMap = map[ordersmessages.OrdersType]internalmessages.OrdersType{
	ordersmessages.OrdersTypeAccession:           internalmessages.OrdersTypeACCESSION,
	ordersmessages.OrdersTypeBetweenDutyStations: internalmessages.OrdersTypeBETWEENDUTYSTATIONS,
	ordersmessages.OrdersTypeBrac:                internalmessages.OrdersTypeBRAC,
	ordersmessages.OrdersTypeCot:                 internalmessages.OrdersTypeCOT,
	ordersmessages.OrdersTypeEmergencyEvac:       internalmessages.OrdersTypeEMERGENCYEVAC,
	ordersmessages.OrdersTypeIpcot:               internalmessages.OrdersTypeIPCOT,
	ordersmessages.OrdersTypeLowCostTravel:       internalmessages.OrdersTypeLOWCOSTTRAVEL,
	ordersmessages.OrdersTypeOperational:         internalmessages.OrdersTypeOPERATIONAL,
	ordersmessages.OrdersTypeOteip:               internalmessages.OrdersTypeOTEIP,
	ordersmessages.OrdersTypeRotational:          internalmessages.OrdersTypeROTATIONAL,
	ordersmessages.OrdersTypeSeparation:          internalmessages.OrdersTypeSEPARATION,
	ordersmessages.OrdersTypeSpecialPurpose:      internalmessages.OrdersTypeSPECIALPURPOSE,
	ordersmessages.OrdersTypeTraining:            internalmessages.OrdersTypeTRAINING,
	ordersmessages.OrdersTypeUnitMove:            internalmessages.OrdersTypeUNITMOVE,
}

func fromAPIOrdersType(ot ordersmessages.OrdersType) (internalmessages.OrdersType, error) {
	ordersType := fromAPIOrdersTypeMap[ot]
	if ordersType == "" {
		return "", fmt.Errorf("Unknown orders type %s", ot)
	}
	return ordersType, nil
}
