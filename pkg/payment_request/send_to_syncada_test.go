package paymentrequest

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func (suite *PaymentRequestHelperSuite) TestSendToSyncada() {
	suite.T().Run("returns no error if send is false", func(t *testing.T) {
		sftpSender := services.SyncadaSFTPSender(nil)
		gexSender := services.GexSender(nil)
		err := SendToSyncada(suite.AppContextForTest(), "edi string", 12345, gexSender, sftpSender, false)
		suite.NoError(err)
	})

	suite.T().Run("returns error if no sender", func(t *testing.T) {
		sftpSender := services.SyncadaSFTPSender(nil)
		gexSender := services.GexSender(nil)
		err := SendToSyncada(suite.AppContextForTest(), "edi string", 12345, gexSender, sftpSender, true)
		suite.Error(err)
		suite.Equal("cannot send to Syncada, SendToSyncada() senders are nil", err.Error())
	})

	suite.T().Run("successful on gex sender", func(t *testing.T) {
		fakeEdi := "pretend this is an edi"
		response := &http.Response{StatusCode: http.StatusOK}
		sftpSender := services.SyncadaSFTPSender(nil)
		gexSender := &mocks.GexSender{}
		expectedFilename := fmt.Sprintf("%s_%d_edi858.txt", time.Now().Format("2006_01_02T15_04_05Z07_00"), 12345)
		gexSender.
			On("SendToGex", services.GEXChannelInvoice, fakeEdi, expectedFilename).Return(response, nil)
		err := SendToSyncada(suite.AppContextForTest(), fakeEdi, 12345, gexSender, sftpSender, true)
		suite.NoError(err)
	})

	suite.T().Run("unsuccessful on gex sender", func(t *testing.T) {
		fakeEdi := "pretend this is an edi"
		response := &http.Response{StatusCode: http.StatusForbidden}
		sftpSender := services.SyncadaSFTPSender(nil)
		gexSender := &mocks.GexSender{}
		expectedFilename := fmt.Sprintf("%s_%d_edi858.txt", time.Now().Format("2006_01_02T15_04_05Z07_00"), 12345)
		gexSender.
			On("SendToGex", services.GEXChannelInvoice, fakeEdi, expectedFilename).Return(response, nil)
		err := SendToSyncada(suite.AppContextForTest(), fakeEdi, 12345, gexSender, sftpSender, true)
		suite.Error(err)
		suite.Contains("received error response when sending EDI to GEX &{ 403  0 0 map[] <nil> 0 [] false false map[] <nil> <nil>}", err.Error())
	})

	suite.T().Run("no response on gex sender", func(t *testing.T) {
		fakeEdi := "pretend this is an edi"
		sftpSender := services.SyncadaSFTPSender(nil)
		gexSender := &mocks.GexSender{}
		expectedFilename := fmt.Sprintf("%s_%d_edi858.txt", time.Now().Format("2006_01_02T15_04_05Z07_00"), 12345)
		gexSender.
			On("SendToGex", services.GEXChannelInvoice, fakeEdi, expectedFilename).Return(nil, nil)
		err := SendToSyncada(suite.AppContextForTest(), fakeEdi, 12345, gexSender, sftpSender, true)
		suite.Error(err)
		suite.Contains("no response when sending EDI to GEX", err.Error())
	})

	suite.T().Run("error response on gex sender", func(t *testing.T) {
		fakeEdi := "pretend this is an edi"
		sftpSender := services.SyncadaSFTPSender(nil)
		gexSender := &mocks.GexSender{}
		expectedFilename := fmt.Sprintf("%s_%d_edi858.txt", time.Now().Format("2006_01_02T15_04_05Z07_00"), 12345)
		gexSender.
			On("SendToGex", services.GEXChannelInvoice, fakeEdi, expectedFilename).Return(nil, fmt.Errorf("gex send threw error"))
		err := SendToSyncada(suite.AppContextForTest(), fakeEdi, 12345, gexSender, sftpSender, true)
		suite.Error(err)
		suite.Contains("GEX sender encountered an error: gex send threw error", err.Error())
	})

	suite.T().Run("error on sftp sender", func(t *testing.T) {
		bytesSent := int64(0)
		// int64, error
		sftpSender := &mocks.SyncadaSFTPSender{}
		sftpSender.
			On("SendToSyncadaViaSFTP", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(bytesSent, fmt.Errorf("test error"))
		gexSender := services.GexSender(nil)
		err := SendToSyncada(suite.AppContextForTest(), "edi string", 12345, gexSender, sftpSender, true)
		suite.Error(err)
		suite.Equal("test error", err.Error())
	})

	suite.T().Run("successful on sftp sender", func(t *testing.T) {
		fakeEdi := "pretend this is an edi"
		bytesSent := int64(10)
		// int64, error
		sftpSender := &mocks.SyncadaSFTPSender{}
		expectedFilename := fmt.Sprintf("%s_%d_edi858.txt", time.Now().Format("2006_01_02T15_04_05Z07_00"), 12345)
		sftpSender.
			On("SendToSyncadaViaSFTP", mock.AnythingOfType("*appcontext.appContext"), strings.NewReader(fakeEdi), expectedFilename).Return(bytesSent, nil)
		gexSender := services.GexSender(nil)
		err := SendToSyncada(suite.AppContextForTest(), fakeEdi, 12345, gexSender, sftpSender, true)
		suite.NoError(err)
	})
}
