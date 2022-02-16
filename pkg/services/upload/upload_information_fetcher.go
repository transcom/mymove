package upload

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
)

type uploadInformationFetcher struct {
}

// NewUploadInformationFetcher return an implementation of the UploadInformationFetcher interface
func NewUploadInformationFetcher() services.UploadInformationFetcher {
	return &uploadInformationFetcher{}
}

// FetchUploadInformation fetches upload information
func (uif *uploadInformationFetcher) FetchUploadInformation(appCtx appcontext.AppContext, uploadID uuid.UUID) (services.UploadInformation, error) {
	q := `
SELECT uploads.id as upload_id,
       uploads.content_type,
       uploads.created_at,
       uploads.filename,
       uploads.bytes,
       moves.locator,
       sm.id AS service_member_id,
       sm.first_name AS service_member_first_name,
       sm.last_name AS service_member_last_name,
       sm.personal_email AS service_member_email,
       sm.telephone AS service_member_telephone,
       ou.id AS office_user_id,
       ou.first_name AS office_user_first_name,
       ou.last_name AS office_user_last_name,
       ou.email AS office_user_email,
       ou.telephone AS office_user_telephone
FROM uploads
         LEFT JOIN user_uploads ON uploads.id = user_uploads.upload_id
         LEFT JOIN users u ON user_uploads.uploader_id = u.id
         LEFT JOIN documents d ON user_uploads.document_id = d.id
         LEFT JOIN service_members documents_service_members ON d.service_member_id = documents_service_members.id
         LEFT JOIN orders ON documents_service_members.id = orders.service_member_id
         LEFT JOIN moves ON orders.id = moves.orders_id
         LEFT JOIN service_members sm ON u.id = sm.user_id
         LEFT JOIN office_users ou ON u.id = ou.user_id
where uploads.id = $1`
	ui := services.UploadInformation{}
	err := appCtx.DB().RawQuery(q, uploadID).First(&ui)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return services.UploadInformation{}, apperror.NewNotFoundError(uploadID, "")
		default:
			return services.UploadInformation{}, apperror.NewQueryError("UploadInformation", err, "")
		}
	}
	return ui, nil
}
