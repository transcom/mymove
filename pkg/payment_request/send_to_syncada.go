package paymentrequest

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
)

// SendToSyncada send EDI file to Syncada for processing
func SendToSyncada(appCtx appcontext.AppContext, edi string, icn int64, gexSender services.GexSender, sftpSender services.SyncadaSFTPSender, sendEDIFile bool) error {
	logger := appCtx.Logger()
	syncadaFileName := fmt.Sprintf("%s_%d_edi858.txt", time.Now().Format("2006_01_02T15_04_05Z07_00"), icn)

	if !sendEDIFile {
		logger.Info("SendToSyncada() is in do not send mode, syncadaFileName: " + syncadaFileName + "")
		return nil
	}

	if (gexSender == nil) && (sftpSender == nil) {
		return fmt.Errorf("cannot send to Syncada, SendToSyncada() senders are nil")
	}
	if gexSender != nil {
		logger.Info("SendToSyncada() is in send mode using GEX, sending syncadaFileName: " + syncadaFileName)
		resp, err := gexSender.SendToGex(services.GEXChannelInvoice, edi, syncadaFileName)
		if err != nil {
			logger.Error("GEX Sender encountered an error", zap.Error(err))
			return fmt.Errorf("GEX sender encountered an error: %w", err)
		}
		if resp == nil {
			logger.Error("GEX Sender receieved no response from GEX")
			return fmt.Errorf("no response when sending EDI to GEX")
		}
		if resp.StatusCode != http.StatusOK {
			logger.Error("func SendToSyncada() failed send to GEX with", zap.Int("StatusCode", resp.StatusCode), zap.String("Status", resp.Status))
			return fmt.Errorf("received error response when sending EDI to GEX %v", resp)
		}
		logger.Info(
			"SUCCESS: EDI858 Processor sent a new file to syncada for Payment Request, using GEX",
			zap.String("filename", syncadaFileName))
	} else if sftpSender != nil {
		// Send to Syncada via SFTP
		edi858String := strings.NewReader(edi)

		logger.Info("SendToSyncada() is in send mode, sending syncadaFileName: " + syncadaFileName + "")
		_, err := sftpSender.SendToSyncadaViaSFTP(appCtx, edi858String, syncadaFileName)
		if err != nil {
			return err
		}
		logger.Info("SUCCESS: EDI858 Processor sent new file to syncada for Payment Request", zap.String("syncadaFileName", syncadaFileName))
	}
	return nil
}
