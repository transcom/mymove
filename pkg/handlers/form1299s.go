package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	form1299op "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/form1299s"
	"github.com/transcom/mymove/pkg/gen/internalmodel"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForAddressModel(a *models.Address) *internalmodel.Address {
	if a != nil {
		return &internalmodel.Address{
			StreetAddress1: swag.String(a.StreetAddress1),
			StreetAddress2: a.StreetAddress2,
			City:           swag.String(a.City),
			State:          swag.String(a.State),
			Zip:            swag.String(a.Zip),
		}
	}
	return nil
}

func addressModelFromPayload(rawAddress *internalmodel.Address) *models.Address {
	if rawAddress == nil {
		return nil
	}
	address := models.Address{
		StreetAddress1: *rawAddress.StreetAddress1,
		StreetAddress2: rawAddress.StreetAddress2,
		City:           *rawAddress.City,
		State:          *rawAddress.State,
		Zip:            *rawAddress.Zip,
	}
	return &address
}

func payloadForForm1299Model(form1299 models.Form1299) internalmodel.Form1299Payload {
	form1299Payload := internalmodel.Form1299Payload{
		CreatedAt:                              fmtDateTime(form1299.CreatedAt),
		ID:                                     fmtUUID(form1299.ID),
		UpdatedAt:                              fmtDateTime(form1299.UpdatedAt),
		DatePrepared:                           (*strfmt.Date)(form1299.DatePrepared),
		ShipmentNumber:                         form1299.ShipmentNumber,
		NameOfPreparingOffice:                  form1299.NameOfPreparingOffice,
		DestOfficeName:                         form1299.DestOfficeName,
		OriginOfficeAddressName:                form1299.OriginOfficeAddressName,
		OriginOfficeAddress:                    payloadForAddressModel(form1299.OriginOfficeAddress),
		ServiceMemberFirstName:                 form1299.ServiceMemberFirstName,
		ServiceMemberMiddleInitial:             form1299.ServiceMemberMiddleInitial,
		ServiceMemberLastName:                  form1299.ServiceMemberLastName,
		ServiceMemberSsn:                       form1299.ServiceMemberSsn,
		ServiceMemberAgency:                    form1299.ServiceMemberAgency,
		ServiceMemberRank:                      form1299.ServiceMemberRank,
		HhgTotalPounds:                         form1299.HhgTotalPounds,
		HhgProgearPounds:                       form1299.HhgProgearPounds,
		HhgValuableItemsCartons:                form1299.HhgValuableItemsCartons,
		MobileHomeSerialNumber:                 form1299.MobileHomeSerialNumber,
		MobileHomeLengthFt:                     form1299.MobileHomeLengthFt,
		MobileHomeLengthInches:                 form1299.MobileHomeLengthInches,
		MobileHomeWidthFt:                      form1299.MobileHomeWidthFt,
		MobileHomeWidthInches:                  form1299.MobileHomeWidthInches,
		MobileHomeHeightFt:                     form1299.MobileHomeHeightFt,
		MobileHomeHeightInches:                 form1299.MobileHomeHeightInches,
		MobileHomeTypeExpando:                  form1299.MobileHomeTypeExpando,
		MobileHomeContentsPackedRequested:      &form1299.MobileHomeContentsPackedRequested,
		MobileHomeBlockedRequested:             &form1299.MobileHomeBlockedRequested,
		MobileHomeUnblockedRequested:           &form1299.MobileHomeUnblockedRequested,
		MobileHomeStoredAtOriginRequested:      &form1299.MobileHomeStoredAtOriginRequested,
		MobileHomeStoredAtDestinationRequested: &form1299.MobileHomeStoredAtDestinationRequested,
		StationOrdersType:                      form1299.StationOrdersType,
		StationOrdersIssuedBy:                  form1299.StationOrdersIssuedBy,
		StationOrdersNewAssignment:             form1299.StationOrdersNewAssignment,
		StationOrdersDate:                      (*strfmt.Date)(form1299.StationOrdersDate),
		StationOrdersNumber:                    form1299.StationOrdersNumber,
		StationOrdersParagraphNumber:           form1299.StationOrdersParagraphNumber,
		StationOrdersInTransitTelephone:        form1299.StationOrdersInTransitTelephone,
		InTransitAddress:                       payloadForAddressModel(form1299.InTransitAddress),
		PickupAddress:                          payloadForAddressModel(form1299.PickupAddress),
		PickupTelephone:                        form1299.PickupTelephone,
		DestAddress:                            payloadForAddressModel(form1299.DestAddress),
		AgentToReceiveHhg:                      form1299.AgentToReceiveHhg,
		ExtraAddress:                           payloadForAddressModel(form1299.ExtraAddress),
		PackScheduledDate:                      (*strfmt.Date)(form1299.PackScheduledDate),
		PickupScheduledDate:                    (*strfmt.Date)(form1299.PickupScheduledDate),
		DeliveryScheduledDate:                  (*strfmt.Date)(form1299.DeliveryScheduledDate),
		Remarks:                                form1299.Remarks,
		OtherMove1From:                         form1299.OtherMove1From,
		OtherMove1To:                           form1299.OtherMove1To,
		OtherMove1NetPounds:                    form1299.OtherMove1NetPounds,
		OtherMove1ProgearPounds:                form1299.OtherMove1ProgearPounds,
		OtherMove2From:                         form1299.OtherMove2From,
		OtherMove2To:                           form1299.OtherMove2To,
		OtherMove2NetPounds:                    form1299.OtherMove2NetPounds,
		OtherMove2ProgearPounds:                form1299.OtherMove2ProgearPounds,
		ServiceMemberSignature:                 form1299.ServiceMemberSignature,
		DateSigned:                             (*strfmt.Date)(form1299.DateSigned),
		ContractorAddress:                      payloadForAddressModel(form1299.ContractorAddress),
		ContractorName:                         form1299.ContractorName,
		NonavailabilityOfSignatureReason:       form1299.NonavailabilityOfSignatureReason,
		CertifiedBySignature:                   form1299.CertifiedBySignature,
		TitleOfCertifiedBySignature:            form1299.TitleOfCertifiedBySignature,
	}
	return form1299Payload
}

// ShowForm1299Handler fetches a single form1299 by id
func ShowForm1299Handler(params form1299op.ShowForm1299Params) middleware.Responder {
	formID := params.Form1299ID

	var response middleware.Responder
	// remove this validation when https://github.com/go-swagger/go-swagger/pull/1394 is merged.
	if strfmt.IsUUID(string(formID)) {
		form, err := models.FetchForm1299ByID(dbConnection, formID)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				response = form1299op.NewShowForm1299NotFound()
			} else {
				// This is an unknown error from the db, nothing to do but log and 500
				zap.L().Error("DB Insertion error", zap.Error(err))
				response = form1299op.NewShowForm1299InternalServerError()
			}
		} else {
			formPayload := payloadForForm1299Model(form)
			response = form1299op.NewShowForm1299OK().WithPayload(&formPayload)
		}
	} else {
		return form1299op.NewShowForm1299BadRequest()
	}

	return response
}

// CreateForm1299Handler creates a new form1299 via POST /form1299
func CreateForm1299Handler(params form1299op.CreateForm1299Params) middleware.Responder {
	originOfficeAddress := addressModelFromPayload(params.CreateForm1299Payload.OriginOfficeAddress)
	inTransitAddress := addressModelFromPayload(params.CreateForm1299Payload.InTransitAddress)
	pickupAddress := addressModelFromPayload(params.CreateForm1299Payload.PickupAddress)
	destAddress := addressModelFromPayload(params.CreateForm1299Payload.DestAddress)
	extraAddress := addressModelFromPayload(params.CreateForm1299Payload.ExtraAddress)
	contractorAddress := addressModelFromPayload(params.CreateForm1299Payload.ContractorAddress)

	newForm1299 := models.Form1299{
		DatePrepared:                           (*time.Time)(params.CreateForm1299Payload.DatePrepared),
		ShipmentNumber:                         params.CreateForm1299Payload.ShipmentNumber,
		NameOfPreparingOffice:                  params.CreateForm1299Payload.NameOfPreparingOffice,
		DestOfficeName:                         params.CreateForm1299Payload.DestOfficeName,
		OriginOfficeAddressName:                params.CreateForm1299Payload.OriginOfficeAddressName,
		OriginOfficeAddress:                    originOfficeAddress,
		ServiceMemberFirstName:                 params.CreateForm1299Payload.ServiceMemberFirstName,
		ServiceMemberMiddleInitial:             params.CreateForm1299Payload.ServiceMemberMiddleInitial,
		ServiceMemberLastName:                  params.CreateForm1299Payload.ServiceMemberLastName,
		ServiceMemberSsn:                       params.CreateForm1299Payload.ServiceMemberSsn,
		ServiceMemberAgency:                    params.CreateForm1299Payload.ServiceMemberAgency,
		ServiceMemberRank:                      params.CreateForm1299Payload.ServiceMemberRank,
		HhgTotalPounds:                         params.CreateForm1299Payload.HhgTotalPounds,
		HhgProgearPounds:                       params.CreateForm1299Payload.HhgProgearPounds,
		HhgValuableItemsCartons:                params.CreateForm1299Payload.HhgValuableItemsCartons,
		MobileHomeSerialNumber:                 params.CreateForm1299Payload.MobileHomeSerialNumber,
		MobileHomeLengthFt:                     params.CreateForm1299Payload.MobileHomeLengthFt,
		MobileHomeLengthInches:                 params.CreateForm1299Payload.MobileHomeWidthInches,
		MobileHomeWidthFt:                      params.CreateForm1299Payload.MobileHomeWidthFt,
		MobileHomeWidthInches:                  params.CreateForm1299Payload.MobileHomeWidthInches,
		MobileHomeHeightFt:                     params.CreateForm1299Payload.MobileHomeHeightFt,
		MobileHomeHeightInches:                 params.CreateForm1299Payload.MobileHomeHeightInches,
		MobileHomeTypeExpando:                  params.CreateForm1299Payload.MobileHomeTypeExpando,
		MobileHomeContentsPackedRequested:      *params.CreateForm1299Payload.MobileHomeContentsPackedRequested,
		MobileHomeBlockedRequested:             *params.CreateForm1299Payload.MobileHomeBlockedRequested,
		MobileHomeUnblockedRequested:           *params.CreateForm1299Payload.MobileHomeUnblockedRequested,
		MobileHomeStoredAtOriginRequested:      *params.CreateForm1299Payload.MobileHomeStoredAtOriginRequested,
		MobileHomeStoredAtDestinationRequested: *params.CreateForm1299Payload.MobileHomeStoredAtDestinationRequested,
		StationOrdersType:                      params.CreateForm1299Payload.StationOrdersType,
		StationOrdersIssuedBy:                  params.CreateForm1299Payload.StationOrdersIssuedBy,
		StationOrdersNewAssignment:             params.CreateForm1299Payload.StationOrdersNewAssignment,
		StationOrdersDate:                      (*time.Time)(params.CreateForm1299Payload.StationOrdersDate),
		StationOrdersNumber:                    params.CreateForm1299Payload.StationOrdersNumber,
		StationOrdersParagraphNumber:           params.CreateForm1299Payload.StationOrdersParagraphNumber,
		StationOrdersInTransitTelephone:        params.CreateForm1299Payload.StationOrdersInTransitTelephone,
		InTransitAddress:                       inTransitAddress,
		PickupAddress:                          pickupAddress,
		PickupTelephone:                        params.CreateForm1299Payload.PickupTelephone,
		DestAddress:                            destAddress,
		AgentToReceiveHhg:                      params.CreateForm1299Payload.AgentToReceiveHhg,
		ExtraAddress:                           extraAddress,
		PackScheduledDate:                      (*time.Time)(params.CreateForm1299Payload.PackScheduledDate),
		PickupScheduledDate:                    (*time.Time)(params.CreateForm1299Payload.PickupScheduledDate),
		DeliveryScheduledDate:                  (*time.Time)(params.CreateForm1299Payload.DeliveryScheduledDate),
		Remarks:                                params.CreateForm1299Payload.Remarks,
		OtherMove1From:                         params.CreateForm1299Payload.OtherMove1From,
		OtherMove1To:                           params.CreateForm1299Payload.OtherMove1To,
		OtherMove1NetPounds:                    params.CreateForm1299Payload.OtherMove1NetPounds,
		OtherMove1ProgearPounds:                params.CreateForm1299Payload.OtherMove1ProgearPounds,
		OtherMove2From:                         params.CreateForm1299Payload.OtherMove2From,
		OtherMove2To:                           params.CreateForm1299Payload.OtherMove2To,
		OtherMove2NetPounds:                    params.CreateForm1299Payload.OtherMove2NetPounds,
		OtherMove2ProgearPounds:                params.CreateForm1299Payload.OtherMove2ProgearPounds,
		ServiceMemberSignature:                 params.CreateForm1299Payload.ServiceMemberSignature,
		DateSigned:                             (*time.Time)(params.CreateForm1299Payload.DateSigned),
		ContractorAddress:                      contractorAddress,
		ContractorName:                         params.CreateForm1299Payload.ContractorName,
		NonavailabilityOfSignatureReason:       params.CreateForm1299Payload.NonavailabilityOfSignatureReason,
		CertifiedBySignature:                   params.CreateForm1299Payload.CertifiedBySignature,
		TitleOfCertifiedBySignature:            params.CreateForm1299Payload.TitleOfCertifiedBySignature,
	}
	var response middleware.Responder
	verrs, err := models.CreateForm1299WithAddresses(dbConnection, &newForm1299)
	if verrs.HasAny() {
		zap.L().Error("DB Validation", zap.Error(verrs))
		response = form1299op.NewCreateForm1299BadRequest()
	} else if err != nil {
		zap.L().Error("DB Insertion", zap.Error(err))
		response = form1299op.NewCreateForm1299BadRequest()
	} else {
		form1299Payload := payloadForForm1299Model(newForm1299)
		response = form1299op.NewCreateForm1299Created().WithPayload(&form1299Payload)
	}
	return response
}

// IndexForm1299sHandler returns a list of all form1299s
func IndexForm1299sHandler(params form1299op.IndexForm1299sParams) middleware.Responder {
	var form1299s models.Form1299s
	var response middleware.Responder
	form1299s, err := models.FetchAllForm1299s(dbConnection)
	if err != nil {
		zap.L().Error("DB Query", zap.Error(err))
		response = form1299op.NewIndexForm1299sBadRequest()
	} else {
		form1299Payloads := make(internalmodel.IndexForm1299sPayload, len(form1299s))
		for i, form1299 := range form1299s {
			form1299Payload := payloadForForm1299Model(form1299)
			form1299Payloads[i] = &form1299Payload
		}
		response = form1299op.NewIndexForm1299sOK().WithPayload(form1299Payloads)
	}
	return response
}
