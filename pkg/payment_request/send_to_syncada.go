package paymentrequest

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/transcom/mymove/pkg/services/invoice"
)

// SendToSyncada send EDI file to Syncada for processing
func SendToSyncada(edi string, sendEDIFile bool) error {
	//TODO: Send to Syncada via GEX needs to be implemented

	// Send to Syncada via SFTP
	SFTPSession := invoice.NewSyncadaSFTPSession(os.Getenv("SYNCADA_SFTP_PORT"), os.Getenv("SYNCADA_SFTP_USER_ID"), os.Getenv("SYNCADA_SFTP_IP_ADDRESS"), os.Getenv("SYNCADA_SFTP_PASSWORD"), os.Getenv("SYNCADA_SFTP_INBOUND_DIRECTORY"))

	edi858String := strings.NewReader(edi)
	syncadaFileName := fmt.Sprintf("%s_edi858.txt", time.Now().Format("2006_01_02T15_04_05Z07_00"))

	var err error
	if sendEDIFile == true {
		_, err = SFTPSession.SendToSyncadaViaSFTP(edi858String, syncadaFileName)
		if err != nil {
			return err
		}
	} else {
		// TODO add logger
		// TODO log fileneme and say that we are not in send mode
	}
	return err
}
