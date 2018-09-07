package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
)

type fakeModel struct {
	FieldName string
}

func noErr(err error) {
	if err != nil {
		log.Panic("oops ", err)
	}
}

func stringPtr(s string) *string {
	return &s
}

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	shipmentID := flag.String("shipment", "", "The shipment ID to generate 1203 form for")
	flag.Parse()

	// DB connection
	err := pop.AddLookupPaths(*config)
	if err != nil {
		log.Fatal(err)
	}
	db, err := pop.Connect(*env)
	if err != nil {
		log.Fatal(err)
	}

	if *shipmentID == "" {
		log.Fatal("Usage: paperwork -shipment <29cb984e-c70d-46f0-926d-cd89e07a6ec3>")
	}

	// This is the path to an image you want to use as a form template
	templateImagePath := "./cmd/generate_1203_form/FORM_1203.png"

	f, err := os.Open(templateImagePath)
	noErr(err)
	defer f.Close()

	// Define your field positions here, it should be a mapping from a struct field name
	// to a FieldPos, which encodes the x and y location, and width of a form field
	var fields = map[string]paperwork.FieldPos{
		"ServiceAgentName":         paperwork.NewFieldPos(28, 12, 79),
		"StandardCarrierAlphaCode": paperwork.NewFieldPos(109, 16, 19),
		"CodeOfService":            paperwork.NewFieldPos(131, 16, 19),
		"ShipmentNumber":           paperwork.NewFieldPos(152, 16, 19),
		"DateIssued":               paperwork.NewFieldPos(173, 16, 40),
		"RequestedPackDate":        paperwork.NewFieldPos(3, 29.5, 19),
		"RequestedPickupDate":      paperwork.NewFieldPos(24, 29.5, 19),
		"RequiredDeliveryDate":     paperwork.NewFieldPos(45, 29.5, 19),
		"ServiceMemberFullName":    paperwork.NewFieldPos(109, 26.5, 30),
		"ServiceMemberEdipi":       paperwork.NewFieldPos(140, 26.5, 25),
		"ServiceMemberRank":        paperwork.NewFieldPos(165, 26.5, 50),

		// "TSPName": paperwork.NewFieldPos(),
		// "ServiceMemberStatus": paperwork.NewFieldPos(),
		// "ServiceMemberDependentStatus": paperwork.NewFieldPos(),

		"AuthorityForShipment":        paperwork.NewFieldPos(110, 37.5, 60),
		"OrdersIssueDate":             paperwork.NewFieldPos(174, 37.5, 25),
		"SecondaryPickupAddress":      paperwork.NewFieldPos(3, 39, 60),
		"ServiceMemberAffiliation":    paperwork.NewFieldPos(110, 47, 60),
		"TransportationControlNumber": paperwork.NewFieldPos(174, 47, 25),
		"FullNameOfShipper":           paperwork.NewFieldPos(110, 58, 100),
		"ConsigneeName":               paperwork.NewFieldPos(3, 75, 100),
		"ConsigneeAddress":            paperwork.NewFieldPos(3, 78, 100),
		"PickupAddress":               paperwork.NewFieldPos(110, 75, 100),

		// "NTSDetails": paperwork.NewFieldPos(),

		"ResponsibleDestinationOffice": paperwork.NewFieldPos(3, 92, 80),
		"DestinationGbloc":             paperwork.NewFieldPos(95, 89, 17),
		"BillChargesToName":            paperwork.NewFieldPos(110, 92, 80),
		"BillChargesToAddress":         paperwork.NewFieldPos(110, 96, 80),

		// "FreightBillNumber": paperwork.NewFieldPos(),

		"AppropriationsChargeable": paperwork.NewFieldPos(110, 109, 80),

		"Remarks": paperwork.NewFieldPos(3, 125, 160),
		// "PackagesNumber": paperwork.NewFieldPos(),
		// "PackagesKind": paperwork.NewFieldPos(),
		"DescriptionOfShipment": paperwork.NewFieldPos(45, 151, 60),
		// "WeightGrossPounds": paperwork.NewFieldPos(),
		// "WeightTarePounds": paperwork.NewFieldPos(),
		// "WeightNetPounds": paperwork.NewFieldPos(),
		// "LineHaulTransportationRate": paperwork.NewFieldPos(),
		// "LineHaulTransportationCharges": paperwork.NewFieldPos(),
		// "PackingUnpackingCharges": paperwork.NewFieldPos(),
		// "OtherAccessorialServices": paperwork.NewFieldPos(),
		// "TariffOrSpecialRateAuthorities": paperwork.NewFieldPos(),
		// "IssuingOfficerFullName": paperwork.NewFieldPos(),
		// "IssuingOfficerTitle": paperwork.NewFieldPos(),
		// "IssuingOfficeName": paperwork.NewFieldPos(),
		// "IssuingOfficeAddress": paperwork.NewFieldPos(),
		// "IssuingOfficeGBLOC": paperwork.NewFieldPos(),
		// "DateOfReceiptOfShipment": paperwork.NewFieldPos(),
		// "SignatureOfAgentOrDriver": paperwork.NewFieldPos(),
		// "PerInitials": paperwork.NewFieldPos(),
		// "ForUsePayingOfficerUnauthorizedItems": paperwork.NewFieldPos(),
		// "ForUsePayingOfficerExcessDistance": paperwork.NewFieldPos(),
		// "ForUsePayingOfficerExcessValuation": paperwork.NewFieldPos(),
		// "ForUsePayingOfficerExcessWeight": paperwork.NewFieldPos(),
		// "ForUsePayingOfficerOther": paperwork.NewFieldPos(),
		// "CertOfTSPBillingDate": paperwork.NewFieldPos(),
		// "CertOfTSPBillingDeliveryPoint": paperwork.NewFieldPos(),
		// "CertOfTSPBillingNameOfDeliveringCarrier": paperwork.NewFieldPos(),
		// "CertOfTSPBillingPlaceDelivered": paperwork.NewFieldPos(),
		// "CertOfTSPBillingShortage": paperwork.NewFieldPos(),
		// "CertOfTSPBillingDamage": paperwork.NewFieldPos(),
		// "CertOfTSPBillingCarrierOSD": paperwork.NewFieldPos(),
		// "CertOfTSPBillingDestinationCarrierName": paperwork.NewFieldPos(),
		// "CertOfTSPBillingAuthorizedAgentSignature": paperwork.NewFieldPos(),
	}

	// Define the data here that you want to populate the form with. Data will only be populated
	// in the form if the field name exist BOTH in the fields map and your data below
	parsedID := uuid.Must(uuid.FromString(*shipmentID))
	gbl, err := models.FetchGovBillOfLadingExtractor(db, parsedID)
	noErr(err)

	// Build our form with a template image and field placement
	form, err := paperwork.NewTemplateForm(f, fields)
	noErr(err)

	// Uncomment the below line if you want to draw borders around field boxes, very useful
	// for getting field positioning right initially
	// form.UseBorders()

	// Populate form fields with provided data
	err = form.DrawData(gbl)
	noErr(err)

	output, _ := os.Create("./cmd/generate_1203_form/test-output.pdf")
	err = form.Output(output)
	noErr(err)

	fmt.Println("done!")
}
