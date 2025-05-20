package ghcapi_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *GhcAPISuite) TestUploads() {

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
		_, err := suite.NewHandlerConfig().FileStorer().Store(uploadUser1.Upload.StorageKey, file.Data, "somehash", nil)
		suite.NoError(err)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(),
			[]roles.RoleType{roles.RoleTypeTOO})
		req := suite.NewAuthenticatedOfficeRequest("GET", "/ghc/v1/uploads/"+uploadUser1.Upload.ID.String()+"/status", nil, officeUser)
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
		_, err := suite.NewHandlerConfig().FileStorer().Store(uploadUser1.Upload.StorageKey, file.Data, "somehash", nil)
		suite.NoError(err)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(),
			[]roles.RoleType{roles.RoleTypeTOO})
		req := suite.NewAuthenticatedOfficeRequest("GET", "/ghc/v1/uploads/"+uploadUser1.Upload.ID.String()+"/status", nil, officeUser)
		rr := httptest.NewRecorder()

		fakeS3, ok := suite.NewHandlerConfig().FileStorer().(*storageTest.FakeS3Storage)
		suite.True(ok)
		suite.NotNil(fakeS3, "FileStorer should be fakeS3")

		fakeS3.EmptyTags = true
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
		suite.Equal("text/event-stream", rr.Header().Get("content-type"))

		suite.Contains(rr.Body.String(), "PROCESSING")
		suite.Contains(rr.Body.String(), "CLEAN")
		suite.Contains(rr.Body.String(), "Connection closed")
	})
}
