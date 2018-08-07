package paperwork

import (
	"path/filepath"

	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

// GenerateAdvancePaperwork generates the advance paperwork for a move.
// Outputs to a tempfile
func GenerateAdvancePaperwork(g *Generator, moveID uuid.UUID, build string) (string, error) {
	move, err := models.FetchMoveForAdvancePaperwork(g.db, moveID)
	if err != nil {
		return "", err
	}

	summary := NewShipmentSummary(&move)
	outfile, err := g.newTempFile()
	if err != nil {
		return "", err
	}
	if err := summary.DrawForm(outfile); err != nil {
		return "", err
	}
	outfile.Close()

	generatedPath := outfile.Name()

	order, err := models.FetchOrderForPDFConversion(g.db, move.OrdersID)
	if err != nil {
		return "", err
	}

	ordersPaths, err := g.ConvertUploadsToPDF(order.UploadedOrders.Uploads)
	if err != nil {
		return "", err
	}

	var inputFiles []string
	g.logger.Debug("adding orders and shipment summary to packet", zap.Any("inputFiles", inputFiles))
	inputFiles = append(ordersPaths, generatedPath)

	for _, ppm := range move.PersonallyProcuredMoves {
		if ppm.Advance != nil && ppm.Advance.MethodOfReceipt == models.MethodOfReceiptOTHERDD {
			g.logger.Debug("adding direct deposit form to packet", zap.Any("inputFiles", inputFiles))
			ddFormPath := filepath.Join(build, "/downloads/direct_deposit_form.pdf")
			inputFiles = append(inputFiles, ddFormPath)
			break
		}
	}

	mergedFile, err := g.MergePDFFiles(inputFiles)
	if err != nil {
		return "", err
	}

	return mergedFile.Name(), nil

}
