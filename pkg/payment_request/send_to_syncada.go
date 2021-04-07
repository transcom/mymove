package paymentrequest

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services"
)

// SendToSyncada send EDI file to Syncada for processing
func SendToSyncada(edi string, gexSender services.GexSender, sftpSender services.SyncadaSFTPSender, sendEDIFile bool, logger Logger) error {
	var err error

	if (gexSender == nil) && (sftpSender == nil) {
		return fmt.Errorf("cannot send to Syncada, SendToSyncada() senders are nil")
	}
	if gexSender != nil {
		//TODO: Send to Syncada via GEX needs to be implemented
		logger.Warn("func SendToSyncada() -- GEX Sender NOT IMPLEMENTED")
	} else if sftpSender != nil {
		// Send to Syncada via SFTP
		edi858String := strings.NewReader(edi)
		syncadaFileName := fmt.Sprintf("%s_edi858.txt", time.Now().Format("2006_01_02T15_04_05Z07_00"))

		if sendEDIFile == true {
			logger.Info("SendToSyncada() is in send mode, sending syncadaFileName: " + syncadaFileName + "")
			_, err = sftpSender.SendToSyncadaViaSFTP(edi858String, syncadaFileName)
			if err != nil {
				return err
			}
			logger.Info("SUCCESS: 858 Processor sent new file to syncada for Payment Request", zap.String("syncadaFileName", syncadaFileName))
		} else {
			logger.Info("SendToSyncada() is in do not send mode, syncadaFileName: " + syncadaFileName + "")
		}
	}
	return err
}
