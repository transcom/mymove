package handlers

import (
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/messages"
	form1299op "github.com/transcom/mymove/pkg/gen/restapi/operations/form1299s"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForForm1299Model(form1299 models.Form1299) messages.Form1299Payload {

	form1299Payload := messages.Form1299Payload{
		CreatedAt:                        fmtDateTime(form1299.CreatedAt),
		ID:                               fmtUUID(form1299.ID),
		UpdatedAt:                        fmtDateTime(form1299.UpdatedAt),
		DatePrepared:                     (*strfmt.Date)(form1299.DatePrepared),
		ShipmentNumber:                   form1299.ShipmentNumber,
		NameOfPreparingOffice:            form1299.NameOfPreparingOffice,
		DestOfficeName:                   form1299.DestOfficeName,
		OriginOfficeAddressName:          form1299.OriginOfficeAddressName,
		OriginOfficeAddress:              form1299.OriginOfficeAddress,
		ServiceMemberFirstName:           form1299.ServiceMemberFirstName,
		ServiceMemberMiddleInitial:       form1299.ServiceMemberMiddleInitial,
		ServiceMemberLastName:            form1299.ServiceMemberLastName,
		ServiceMemberSsn:                 form1299.ServiceMemberSsn,
		ServiceMemberAgency:              form1299.ServiceMemberAgency,
		HhgTotalPounds:                   form1299.HhgTotalPounds,
		HhgProgearPounds:                 form1299.HhgProgearPounds,
		HhgValuableItemsCartons:          form1299.HhgValuableItemsCartons,
		MobileHomeSerialNumber:           form1299.MobileHomeSerialNumber,
		MobileHomeLengthFt:               form1299.MobileHomeLengthFt,
		MobileHomeLengthInches:           form1299.MobileHomeLengthInches,
		MobileHomeWidthFt:                form1299.MobileHomeWidthFt,
		MobileHomeWidthInches:            form1299.MobileHomeWidthInches,
		MobileHomeHeightFt:               form1299.MobileHomeHeightFt,
		MobileHomeHeightInches:           form1299.MobileHomeHeightInches,
		MobileHomeTypeExpando:            form1299.MobileHomeTypeExpando,
		MobileHomeServicesRequested:      form1299.MobileHomeServicesRequested,
		StationOrdersType:                form1299.StationOrdersType,
		StationOrdersIssuedBy:            form1299.StationOrdersIssuedBy,
		StationOrdersNewAssignment:       form1299.StationOrdersNewAssignment,
		StationOrdersDate:                (*strfmt.Date)(form1299.StationOrdersDate),
		StationOrdersNumber:              form1299.StationOrdersNumber,
		StationOrdersParagraphNumber:     form1299.StationOrdersParagraphNumber,
		StationOrdersInTransitTelephone:  form1299.StationOrdersInTransitTelephone,
		InTransitAddress:                 form1299.InTransitAddress,
		PickupAddress:                    form1299.PickupAddress,
		PickupAddressMobileCourtName:     form1299.PickupAddressMobileCourtName,
		PickupTelephone:                  form1299.PickupTelephone,
		DestAddress:                      form1299.DestAddress,
		DestAddressMobileCourtName:       form1299.DestAddressMobileCourtName,
		AgentToReceiveHhg:                form1299.AgentToReceiveHhg,
		ExtraAddress:                     form1299.ExtraAddress,
		PackScheduledDate:                (*strfmt.Date)(form1299.PackScheduledDate),
		PickupScheduledDate:              (*strfmt.Date)(form1299.PickupScheduledDate),
		DeliveryScheduledDate:            (*strfmt.Date)(form1299.DeliveryScheduledDate),
		Remarks:                          form1299.Remarks,
		OtherMoveFrom:                    form1299.OtherMoveFrom,
		OtherMoveTo:                      form1299.OtherMoveTo,
		OtherMoveNetPounds:               form1299.OtherMoveNetPounds,
		OtherMoveProgearPounds:           form1299.OtherMoveProgearPounds,
		ServiceMemberSignature:           form1299.ServiceMemberSignature,
		DateSigned:                       (*strfmt.Date)(form1299.DateSigned),
		ContractorAddress:                form1299.ContractorAddress,
		ContractorName:                   form1299.ContractorName,
		NonavailabilityOfSignatureReason: form1299.NonavailabilityOfSignatureReason,
		CertifiedBySignature:             form1299.CertifiedBySignature,
		TitleOfCertifiedBySignature:      form1299.TitleOfCertifiedBySignature,
	}
	return form1299Payload
}

// ShowForm1299Handler fetches a single form1299 by id
func ShowForm1299Handler(params form1299op.ShowForm1299Params) middleware.Responder {
	fmt.Println("WEOINWEFOWFNWOEFNN")
	fmt.Println(params.Form1299ID)
	formID := params.Form1299ID
	form := models.Form1299{}
	var response middleware.Responder
	if err := dbConnection.Find(&form, formID); err != nil {
		fmt.Println(err)
		response = form1299op.NewShowForm1299NotFound()
		// TODO: read the err
	} else {
		formPayload := payloadForForm1299Model(form)
		response = form1299op.NewShowForm1299OK().WithPayload(&formPayload)
	}

	return response
}

// CreateForm1299Handler creates a new form1299 via POST /form1299
func CreateForm1299Handler(params form1299op.CreateForm1299Params) middleware.Responder {
	newForm1299 := models.Form1299{
		DatePrepared:                     (*time.Time)(params.CreateForm1299Payload.DatePrepared),
		ShipmentNumber:                   params.CreateForm1299Payload.ShipmentNumber,
		NameOfPreparingOffice:            params.CreateForm1299Payload.NameOfPreparingOffice,
		DestOfficeName:                   params.CreateForm1299Payload.DestOfficeName,
		OriginOfficeAddressName:          params.CreateForm1299Payload.OriginOfficeAddressName,
		OriginOfficeAddress:              params.CreateForm1299Payload.OriginOfficeAddress,
		ServiceMemberFirstName:           params.CreateForm1299Payload.ServiceMemberFirstName,
		ServiceMemberMiddleInitial:       params.CreateForm1299Payload.ServiceMemberMiddleInitial,
		ServiceMemberLastName:            params.CreateForm1299Payload.ServiceMemberLastName,
		ServiceMemberSsn:                 params.CreateForm1299Payload.ServiceMemberSsn,
		ServiceMemberAgency:              params.CreateForm1299Payload.ServiceMemberAgency,
		HhgTotalPounds:                   params.CreateForm1299Payload.HhgTotalPounds,
		HhgProgearPounds:                 params.CreateForm1299Payload.HhgProgearPounds,
		HhgValuableItemsCartons:          params.CreateForm1299Payload.HhgValuableItemsCartons,
		MobileHomeSerialNumber:           params.CreateForm1299Payload.MobileHomeSerialNumber,
		MobileHomeLengthFt:               params.CreateForm1299Payload.MobileHomeLengthFt,
		MobileHomeLengthInches:           params.CreateForm1299Payload.MobileHomeWidthInches,
		MobileHomeWidthFt:                params.CreateForm1299Payload.MobileHomeWidthFt,
		MobileHomeWidthInches:            params.CreateForm1299Payload.MobileHomeWidthInches,
		MobileHomeHeightFt:               params.CreateForm1299Payload.MobileHomeHeightFt,
		MobileHomeHeightInches:           params.CreateForm1299Payload.MobileHomeHeightInches,
		MobileHomeTypeExpando:            params.CreateForm1299Payload.MobileHomeTypeExpando,
		MobileHomeServicesRequested:      params.CreateForm1299Payload.MobileHomeServicesRequested,
		StationOrdersType:                params.CreateForm1299Payload.StationOrdersType,
		StationOrdersIssuedBy:            params.CreateForm1299Payload.StationOrdersIssuedBy,
		StationOrdersNewAssignment:       params.CreateForm1299Payload.StationOrdersNewAssignment,
		StationOrdersDate:                (*time.Time)(params.CreateForm1299Payload.StationOrdersDate),
		StationOrdersNumber:              params.CreateForm1299Payload.StationOrdersNumber,
		StationOrdersParagraphNumber:     params.CreateForm1299Payload.StationOrdersParagraphNumber,
		StationOrdersInTransitTelephone:  params.CreateForm1299Payload.StationOrdersInTransitTelephone,
		InTransitAddress:                 params.CreateForm1299Payload.InTransitAddress,
		PickupAddress:                    params.CreateForm1299Payload.PickupAddress,
		PickupAddressMobileCourtName:     params.CreateForm1299Payload.PickupAddressMobileCourtName,
		PickupTelephone:                  params.CreateForm1299Payload.PickupTelephone,
		DestAddress:                      params.CreateForm1299Payload.DestAddress,
		DestAddressMobileCourtName:       params.CreateForm1299Payload.DestAddressMobileCourtName,
		AgentToReceiveHhg:                params.CreateForm1299Payload.AgentToReceiveHhg,
		ExtraAddress:                     params.CreateForm1299Payload.ExtraAddress,
		PackScheduledDate:                (*time.Time)(params.CreateForm1299Payload.PackScheduledDate),
		PickupScheduledDate:              (*time.Time)(params.CreateForm1299Payload.PickupScheduledDate),
		DeliveryScheduledDate:            (*time.Time)(params.CreateForm1299Payload.DeliveryScheduledDate),
		Remarks:                          params.CreateForm1299Payload.Remarks,
		OtherMoveFrom:                    params.CreateForm1299Payload.OtherMoveFrom,
		OtherMoveTo:                      params.CreateForm1299Payload.OtherMoveTo,
		OtherMoveNetPounds:               params.CreateForm1299Payload.OtherMoveNetPounds,
		OtherMoveProgearPounds:           params.CreateForm1299Payload.OtherMoveProgearPounds,
		ServiceMemberSignature:           params.CreateForm1299Payload.ServiceMemberSignature,
		DateSigned:                       (*time.Time)(params.CreateForm1299Payload.DateSigned),
		ContractorAddress:                params.CreateForm1299Payload.ContractorAddress,
		ContractorName:                   params.CreateForm1299Payload.ContractorName,
		NonavailabilityOfSignatureReason: params.CreateForm1299Payload.NonavailabilityOfSignatureReason,
		CertifiedBySignature:             params.CreateForm1299Payload.CertifiedBySignature,
		TitleOfCertifiedBySignature:      params.CreateForm1299Payload.TitleOfCertifiedBySignature,
	}
	var response middleware.Responder
	if err := dbConnection.Create(&newForm1299); err != nil {
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
	if err := dbConnection.All(&form1299s); err != nil {
		zap.L().Error("DB Query", zap.Error(err))
		response = form1299op.NewIndexForm1299sBadRequest()
	} else {
		form1299Payloads := make(messages.IndexForm1299sPayload, len(form1299s))
		for i, form1299 := range form1299s {
			form1299Payload := payloadForForm1299Model(form1299)
			form1299Payloads[i] = &form1299Payload
		}
		response = form1299op.NewIndexForm1299sOK().WithPayload(form1299Payloads)
	}
	return response
}
