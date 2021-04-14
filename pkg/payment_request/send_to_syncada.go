package paymentrequest

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services"
)

// SendToSyncada send EDI file to Syncada for processing
func SendToSyncada(edi string, icn int64, gexSender services.GexSender, sftpSender services.SyncadaSFTPSender, sendEDIFile bool, logger Logger) error {
	syncadaFileName := fmt.Sprintf("%s_%s_edi858.txt", time.Now().Format("2006_01_02T15_04_05Z07_00"), string(icn))

	if (gexSender == nil) && (sftpSender == nil) {
		return fmt.Errorf("cannot send to Syncada, SendToSyncada() senders are nil")
	}
	if gexSender != nil {
		resp, err := gexSender.SendToGex(edi, syncadaFileName)
		if err != nil {
			logger.Error("GEX Sender encountered an error", zap.Error(err))
			return errors.Wrap(err, "GEX sender encountered an error")
		}
		if resp == nil {
			return fmt.Errorf("no response when sending EDI to GEX")
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("received error response when sending EDI to GEX %v", resp)
		}
		logger.Info(
			"Posted to GEX",
			zap.String("filename", syncadaFileName),
			zap.Int("statusCode", resp.StatusCode))
	} else if sftpSender != nil {
		// Send to Syncada via SFTP
		edi858String := strings.NewReader(edi)

		if sendEDIFile == true {
			logger.Info("SendToSyncada() is in send mode, sending syncadaFileName: " + syncadaFileName + "")
			_, err := sftpSender.SendToSyncadaViaSFTP(edi858String, syncadaFileName)
			if err != nil {
				return err
			}
			logger.Info("SUCCESS: 858 Processor sent new file to syncada for Payment Request", zap.String("syncadaFileName", syncadaFileName))
		} else {
			logger.Info("SendToSyncada() is in do not send mode, syncadaFileName: " + syncadaFileName + "")
		}
	}
	return nil
}
