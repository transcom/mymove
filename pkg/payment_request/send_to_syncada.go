package paymentrequest

import (
	"fmt"
	"regexp"
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
		edi858Log(edi, logger)
	}
	return err
}

func edi858Log(edi string, logger Logger) {

	// regex: BX\*\S{2}\*.\*\S{2}\*([^\s\*]{1,30})\*
	// group 1: MTO reference number
	reBX, err1 := regexp.Compile(`BX\*\S{2}\*.\*\S{2}\*(?P<MTORefNum>[^\s\*]{1,30})\*`)
	if err1 != nil {
		logger.Error("edi858Log BX compile regexp error", zap.Error(err1))
	}

	logString := ""
	match := reBX.FindStringSubmatch(edi)
	for i, name := range reBX.SubexpNames() {
		if i > 0 && i <= len(match) {
			logString = fmt.Sprintf(" %s: %s,", name, match[i])
		}
	}

	// regex: ISA\*00\*[^\s\*]{1,30}\*00\*[^\s\*]{1,30}\*[^\s\*]{1,30}\*[^*]{1,30}\*[^\s\*]{1,30}\*[^\*]{1,30}\*([^\s\*]{1,30})\*([^\s\*]{1,30})\*[^\s\*]{1,30}\*[^\s\*]{1,30}\*([^\s\*]{1,30})\*
	// group 1: Date
	// group 2: Time
	// group 3: ICN
	reISA, err2 := regexp.Compile(`ISA\*00\*[^\s\*]{1,30}\*00\*[^\s\*]{1,30}\*[^\s\*]{1,30}\*[^*]{1,30}\*[^\s\*]{1,30}\*[^\*]{1,30}\*(?P<Date>[^\s\*]{1,30})\*(?P<Time>[^\s\*]{1,30})\*[^\s\*]{1,30}\*[^\s\*]{1,30}\*(?P<ICN>[^\s\*]{1,30})\*`)
	if err2 != nil {
		logger.Error("edi858Log ISA compile regexp error", zap.Error(err2))
	}

	match2 := reISA.FindStringSubmatch(edi)
	for i, name := range reISA.SubexpNames() {
		if i > 0 && i <= len(match2) {
			logString = logString + fmt.Sprintf(" %s: %s,", name, match2[i])
		}
	}
	logger.Info("EDI 858 sent: ", zap.String("858", logString))
}
