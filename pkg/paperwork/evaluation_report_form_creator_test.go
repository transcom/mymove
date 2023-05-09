package paperwork

import (
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func setupTestData(suite *PaperworkSuite) (models.EvaluationReport, models.ReportViolations, models.MTOShipments, models.ServiceMember) {
	report := factory.BuildEvaluationReport(suite.DB(), nil, nil)
	shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{{
		Model:    report.Move,
		LinkOnly: true,
	}}, nil)
	violations := testdatagen.MakeReportViolation(suite.DB(), testdatagen.Assertions{Report: report})
	return report, models.ReportViolations{violations}, models.MTOShipments{shipment}, report.Move.Orders.ServiceMember
}
func (suite *PaperworkSuite) TestEvaluationReportFormSmokeTests() {
	suite.Run("Shipment report", func() {
		report, violations, shipments, customer := setupTestData(suite)
		formFiller, err := NewEvaluationReportFormFiller()
		suite.NoError(err)

		err = formFiller.CreateShipmentReport(report, violations, shipments[0], customer)
		suite.NoError(err)

		testFs := afero.NewMemMapFs()

		output, err := testFs.Create("test-output.pdf")
		suite.FatalNil(err)

		err = formFiller.Output(output)
		suite.FatalNil(err)
	})
	suite.Run("Counseling report", func() {
		report, violations, shipments, customer := setupTestData(suite)
		formFiller, err := NewEvaluationReportFormFiller()
		suite.NoError(err)

		err = formFiller.CreateCounselingReport(report, violations, shipments, customer)
		suite.NoError(err)

		testFs := afero.NewMemMapFs()

		output, err := testFs.Create("test-output.pdf")
		suite.FatalNil(err)

		err = formFiller.Output(output)
		suite.FatalNil(err)
	})
}
