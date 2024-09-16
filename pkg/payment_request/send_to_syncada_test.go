package paymentrequest

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func (suite *PaymentRequestHelperSuite) TestSendToSyncada() {
	filenameMatcher := mock.MatchedBy(func(filename string) bool {
		r := regexp.MustCompile(`^\d{4}_\d{2}_\d{2}T\d{2}_\d{2}_\d{2}Z_\d{2}_\d+?_edi858\.txt$`)
		return r.MatchString(filename)
	})

	suite.Run("returns no error if send is false", func() {
		sftpSender := services.SyncadaSFTPSender(nil)
		gexSender := services.GexSender(nil)
		_, err := SendToSyncada(suite.AppContextForTest(), "edi string", 12345, gexSender, sftpSender, false)
		suite.NoError(err)
	})

	suite.Run("returns error if no sender", func() {
		sftpSender := services.SyncadaSFTPSender(nil)
		gexSender := services.GexSender(nil)
		_, err := SendToSyncada(suite.AppContextForTest(), "edi string", 12345, gexSender, sftpSender, true)
		suite.Error(err)
		suite.Equal("cannot send to Syncada, SendToSyncada() senders are nil", err.Error())
	})

	suite.Run("successful on gex sender", func() {
		fakeEdi := "pretend this is an edi"
		response := &http.Response{StatusCode: http.StatusOK}
		sftpSender := services.SyncadaSFTPSender(nil)
		gexSender := &mocks.GexSender{}
		gexSender.
			On("SendToGex", services.GEXChannelInvoice, fakeEdi, filenameMatcher).Return(response, nil)
		_, err := SendToSyncada(suite.AppContextForTest(), fakeEdi, 12345, gexSender, sftpSender, true)
		suite.NoError(err)
	})

	suite.Run("unsuccessful on gex sender", func() {
		fakeEdi := "pretend this is an edi"
		response := &http.Response{StatusCode: http.StatusForbidden}
		sftpSender := services.SyncadaSFTPSender(nil)
		gexSender := &mocks.GexSender{}
		gexSender.
			On("SendToGex", services.GEXChannelInvoice, fakeEdi, filenameMatcher).Return(response, nil)
		_, err := SendToSyncada(suite.AppContextForTest(), fakeEdi, 12345, gexSender, sftpSender, true)
		suite.Error(err)
		suite.Contains("received error response when sending EDI to GEX &{ 403  0 0 map[] <nil> 0 [] false false map[] <nil> <nil>}", err.Error())
	})

	suite.Run("no response on gex sender", func() {
		fakeEdi := "pretend this is an edi"
		sftpSender := services.SyncadaSFTPSender(nil)
		gexSender := &mocks.GexSender{}
		gexSender.
			On("SendToGex", services.GEXChannelInvoice, fakeEdi, filenameMatcher).Return(nil, nil)
		_, err := SendToSyncada(suite.AppContextForTest(), fakeEdi, 12345, gexSender, sftpSender, true)
		suite.Error(err)
		suite.Contains("no response when sending EDI to GEX", err.Error())
	})

	suite.Run("error response on gex sender", func() {
		fakeEdi := "pretend this is an edi"
		sftpSender := services.SyncadaSFTPSender(nil)
		gexSender := &mocks.GexSender{}
		gexSender.
			On("SendToGex", services.GEXChannelInvoice, fakeEdi, filenameMatcher).Return(nil, fmt.Errorf("gex send threw error"))
		_, err := SendToSyncada(suite.AppContextForTest(), fakeEdi, 12345, gexSender, sftpSender, true)
		suite.Error(err)
		suite.Contains("GEX sender encountered an error: gex send threw error", err.Error())
	})

	suite.Run("error on sftp sender", func() {
		bytesSent := int64(0)
		// int64, error
		sftpSender := &mocks.SyncadaSFTPSender{}
		sftpSender.
			On("SendToSyncadaViaSFTP", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, filenameMatcher).Return(bytesSent, fmt.Errorf("test error"))
		gexSender := services.GexSender(nil)
		_, err := SendToSyncada(suite.AppContextForTest(), "edi string", 12345, gexSender, sftpSender, true)
		suite.Error(err)
		suite.Equal("test error", err.Error())
	})

	suite.Run("successful on sftp sender", func() {
		fakeEdi := "pretend this is an edi"
		bytesSent := int64(10)
		// int64, error
		sftpSender := &mocks.SyncadaSFTPSender{}
		sftpSender.
			On("SendToSyncadaViaSFTP", mock.AnythingOfType("*appcontext.appContext"), strings.NewReader(fakeEdi), filenameMatcher).Return(bytesSent, nil)
		gexSender := services.GexSender(nil)
		_, err := SendToSyncada(suite.AppContextForTest(), fakeEdi, 12345, gexSender, sftpSender, true)
		suite.NoError(err)
	})
}
