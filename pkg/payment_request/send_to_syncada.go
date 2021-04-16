package paymentrequest

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services"
)

// SendToSyncada send EDI file to Syncada for processing
func SendToSyncada(edi string, gexSender services.GexSender, sftpSender services.SyncadaSFTPSender, sendEDIFile bool, logger Logger) error {
	filenameFormat := "2006_01_02T15_04_05Z07_00"
	if (gexSender == nil) && (sftpSender == nil) {
		logger.Error("cannot send to Syncada, SendToSyncada() senders are nil")
		return fmt.Errorf("cannot send to Syncada, SendToSyncada() senders are nil")
	}
	if gexSender != nil {
		// Send to Syncada via GEX
		txName := fmt.Sprintf("%s_edi858", time.Now().Format(filenameFormat))
		if sendEDIFile == true {
			logger.Info("SendToSyncada() is in send mode using GEX, sending syncadaFileName: " + txName + "")
			resp, err := gexSender.SendToGex(edi, txName)
			if err != nil {
				return err
			}
			if resp.StatusCode != http.StatusOK {
				logger.Error("func SendToSyncada() failed send to GEX with", zap.Int("StatusCode", resp.StatusCode), zap.String("Status", resp.Status))
				return fmt.Errorf("func SendToSyncada() failed send to GEX: StatusCode %d, Status %s", resp.StatusCode, resp.Status)
			}
			logger.Info("SUCCESS: 858 Processor sent new file to syncada for Payment Request, using GEX", zap.String("syncadaFileName", txName))
		} else {
			logger.Info("SendToSyncada() is in do not send mode using GEX, txName: " + txName + "")
		}
	} else if sftpSender != nil {
		// Send to Syncada via SFTP
		edi858String := strings.NewReader(edi)
		syncadaFileName := fmt.Sprintf("%s_edi858.txt", time.Now().Format(filenameFormat))

		if sendEDIFile == true {
			logger.Info("SendToSyncada() is in send mode using Syncada SFTP, sending syncadaFileName: " + syncadaFileName + "")
			_, err := sftpSender.SendToSyncadaViaSFTP(edi858String, syncadaFileName)
			if err != nil {
				return err
			}
			logger.Info("SUCCESS: 858 Processor sent new file to syncada for Payment Request, using Syncada SFTP", zap.String("syncadaFileName", syncadaFileName))
		} else {
			logger.Info("SendToSyncada() is in do not send mode using Syncada SFTP, syncadaFileName: " + syncadaFileName + "")
		}
	}
	return nil
}
