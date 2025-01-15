package internalapi_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *InternalAPISuite) TestUploads() {

	suite.Run("Received status for upload, read tag without event queue", func() {
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

	suite.Run("Received statuses for upload, receiving multiple statuses with event queue", func() {
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

		fakeS3, ok := suite.HandlerConfig().FileStorer().(*storageTest.FakeS3Storage)
		suite.True(ok)
		suite.NotNil(fakeS3, "FileStorer should be fakeS3")

		fakeS3.EmptyTags = true
		go func() {
			time.Sleep(5 * time.Second)
			fakeS3.EmptyTags = false
		}()

		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
		suite.Equal("text/event-stream", rr.Header().Get("content-type"))

		message1 := "id: 0\nevent: message\ndata: PROCESSING\n\n"
		message2 := "id: 1\nevent: message\ndata: CLEAN\n\n"
		messageClose := "id: 2\nevent: close\ndata: Connection closed\n\n"

		suite.Equal(message1+message2+messageClose, rr.Body.String())
	})
}
