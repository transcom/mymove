package invoice

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/services"
)

// SyncadaSenderSFTPSession contains information to create a new Syncada SFTP session
type SyncadaSenderSFTPSession struct {
	userID string
}

// InitNewSyncadaSFTPSession initialize a NewSyncadaSFTPSession and return services.SyncadaSFTPSender
func InitNewSyncadaSFTPSession() (services.SyncadaSFTPSender, error) {
	v := viper.New()
	v.AutomaticEnv()
	userID := v.GetString("GEX_SFTP_USER_ID")
	if userID == "" {
		return nil, fmt.Errorf("Invalid credentials sftp missing GEX_SFTP_USER_ID")
	}

	return &SyncadaSenderSFTPSession{
		userID,
	}, nil
}

// SendToSyncadaViaSFTP copies specified local content to Syncada's SFTP server
func (s *SyncadaSenderSFTPSession) SendToSyncadaViaSFTP(appCtx appcontext.AppContext, localDataReader io.Reader, syncadaFileName string) (int64, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	sshClient, err := cli.InitGEXSSH(appCtx, v)
	if err != nil {
		appCtx.Logger().Fatal("couldn't initialize SSH client", zap.Error(err))
	}
	defer func() {
		if closeErr := sshClient.Close(); closeErr != nil {
			appCtx.Logger().Error("Failed to close connection", zap.Error(closeErr))
		}
	}()

	sftpClient, err := cli.InitGEXSFTP(appCtx, sshClient)
	if err != nil {
		appCtx.Logger().Fatal("couldn't initialize SFTP client", zap.Error(err))
	}

	defer func() {
		if closeErr := sftpClient.Close(); closeErr != nil {
			appCtx.Logger().Error("Failed to close SFTP client", zap.Error(closeErr))
		}
	}()

	// create destination file
	syncadaFilePath := fmt.Sprintf("/%s/%s", s.userID, syncadaFileName)
	syncadaFile, err := sftpClient.Create(syncadaFilePath)
	if err != nil {
		return 0, err
	}

	defer func() {
		if closeErr := syncadaFile.Close(); closeErr != nil {
			appCtx.Logger().Error("Failed to close Syncada destination file", zap.Error(closeErr))
		}
	}()

	// copy source file to destination file
	bytes, err := io.Copy(syncadaFile, localDataReader)
	if err != nil {
		return 0, err
	}

	return bytes, err
}
