package internalapi_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *InternalAPISuite) TestUploads() {
	suite.Run("Received message for upload", func() {
		orders := factory.BuildOrder(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		uploadUser1 := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    orders.UploadedOrders,
				LinkOnly: true,
			},
			{
				Model: models.Upload{
					Filename:    "FileName",
					Bytes:       int64(15),
					ContentType: uploader.FileTypePDF,
				},
			},
		}, nil)
		file := suite.Fixture("test.pdf")
		_, err := suite.HandlerConfig().FileStorer().Store(uploadUser1.Upload.StorageKey, file.Data, "somehash", nil)
		suite.NoError(err)

		req := suite.NewAuthenticatedMilRequest("GET", "/internal/uploads/"+uploadUser1.Upload.ID.String()+"/status", nil, orders.ServiceMember)
		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
		suite.Equal("text/event-stream", rr.Header().Get("content-type"))
		suite.Equal("id: 0\nevent: message\ndata: CLEAN\n\nid: 1\nevent: close\ndata: Connection closed\n\n", rr.Body.String())
	})
}
