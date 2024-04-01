package services

import (
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/spf13/afero"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/paperwork"
)

// WeightEstimatorPage1 is an object representing fields from Page 1 of the pdf
type WeightEstimatorPage1 struct {
	LivingRoomCuFt1    string
	LivingRoomPieces1  string
	LivingRoomTotal1   string
	LivingRoomCuFt2    string
	LivingRoomPieces2  string
	LivingRoomTotal2   string
	LivingRoomCuFt3    string
	LivingRoomPieces3  string
	LivingRoomTotal3   string
	LivingRoomCuFt4    string
	LivingRoomPieces4  string
	LivingRoomTotal4   string
	LivingRoomCuFt5    string
	LivingRoomPieces5  string
	LivingRoomTotal5   string
	LivingRoomCuFt6    string
	LivingRoomPieces6  string
	LivingRoomTotal6   string
	LivingRoomCuFt7    string
	LivingRoomPieces7  string
	LivingRoomTotal7   string
	LivingRoomCuFt8    string
	LivingRoomPieces8  string
	LivingRoomTotal8   string
	LivingRoomCuFt9    string
	LivingRoomPieces9  string
	LivingRoomTotal9   string
	LivingRoomCuFt10   string
	LivingRoomPieces10 string
	LivingRoomTotal10  string
	LivingRoomCuFt11   string
	LivingRoomPieces11 string
	LivingRoomTotal11  string
	LivingRoomCuFt12   string
	LivingRoomPieces12 string
	LivingRoomTotal12  string
	LivingRoomCuFt13   string
	LivingRoomPieces13 string
	LivingRoomTotal13  string
	LivingRoomCuFt14   string
	LivingRoomPieces14 string
	LivingRoomTotal14  string
	LivingRoomCuFt15   string
	LivingRoomPieces15 string
	LivingRoomTotal15  string
	LivingRoomCuFt16   string
	LivingRoomPieces16 string
	LivingRoomTotal16  string
	LivingRoomCuFt17   string
	LivingRoomPieces17 string
	LivingRoomTotal17  string
	LivingRoomCuFt18   string
	LivingRoomPieces18 string
	LivingRoomTotal18  string
	LivingRoomCuFt19   string
	LivingRoomPieces19 string
	LivingRoomTotal19  string
	LivingRoomCuFt20   string
	LivingRoomPieces20 string
	LivingRoomTotal20  string
	LivingRoomCuFt21   string
	LivingRoomPieces21 string
	LivingRoomTotal21  string
	LivingRoomCuFt22   string
	LivingRoomPieces22 string
	LivingRoomTotal22  string
	LivingRoomCuFt23   string
	LivingRoomPieces23 string
	LivingRoomTotal23  string
	LivingRoomCuFt24   string
	LivingRoomPieces24 string
	LivingRoomTotal24  string
	LivingRoomCuFt25   string
	LivingRoomPieces25 string
	LivingRoomTotal25  string
	LivingRoomCuFt26   string
	LivingRoomPieces26 string
	LivingRoomTotal26  string
	LivingRoomCuFt27   string
	LivingRoomPieces27 string
	LivingRoomTotal27  string
	LivingRoomCuFt28   string
	LivingRoomPieces28 string
	LivingRoomTotal28  string
	LivingRoomCuFt29   string
	LivingRoomPieces29 string
	LivingRoomTotal29  string
	LivingRoomCuFt30   string
	LivingRoomPieces30 string
	LivingRoomTotal30  string
	LivingRoomCuFt31   string
	LivingRoomPieces31 string
	LivingRoomTotal31  string
	LivingRoomCuFt32   string
	LivingRoomPieces32 string
	LivingRoomTotal33  string
}

// WeightEstimatorPage2 is an object representing fields from Page 2 of the pdf
type WeightEstimatorPage2 struct {
	LivingRoomCuFt34       string
	LivingRoomPieces34     string
	LivingRoomTotal34      string
	LivingRoomCuFt35       string
	LivingRoomPieces35     string
	LivingRoomTotal35      string
	LivingRoomCuFt36       string
	LivingRoomPieces36     string
	LivingRoomTotal36      string
	LivingRoomCuFt37       string
	LivingRoomPieces37     string
	LivingRoomTotal37      string
	LivingRoomCuFt38       string
	LivingRoomPieces38     string
	LivingRoomTotal38      string
	LivingRoomCuFt39       string
	LivingRoomPieces39     string
	LivingRoomTotal39      string
	LivingRoomCuFt40       string
	LivingRoomPieces40     string
	LivingRoomTotal40      string
	LivingRoomCuFt41       string
	LivingRoomPieces41     string
	LivingRoomTotal41      string
	LivingRoomPiecesTotal1 string
	LivingRoomCuFtTotal1   string
	LivingRoomCuFt42       string
	LivingRoomPieces42     string
	LivingRoomTotal42      string
	LivingRoomPiecesTotal2 string
	LivingRoomCuFtTotal2   string
	LivingRoomTotalItems   string
	LivingRoomTotalCube    string
	LivingRoomWeight       string
}

//go:generate mockery --name SSWPPMComputer
type WeightTicketParserComputer interface {
	ParseWeightEstimatorExcelFile(appCtx appcontext.AppContext, path string, weightGenerator paperwork.Generator) (string, error)
}

//go:generate mockery --name SSWPPMGenerator
type WeightTicketParserGenerator interface {
	FillWeightEstimatorPDFForm(Page1Values WeightEstimatorPage1, Page2Values WeightEstimatorPage2) (weightEstimatorFile afero.File, pdfInfo *pdfcpu.PDFInfo, err error)
}
