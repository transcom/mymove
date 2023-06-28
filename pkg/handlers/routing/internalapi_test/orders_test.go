package internalapi_test

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *InternalAPISuite) TestUploadAmendedOrders() {
	// getCookiesForServiceMember makes a GET request to /internal/users/logged_in and returns the cookies that are set
	//   on the response. These are needed to make unsafe method requests as defined by gorilla/csrf.
	//   https://github.com/gorilla/csrf/blob/master/csrf.go#L32-L33
	getCookiesForServiceMember := func(serviceMember *models.ServiceMember) []*http.Cookie {
		authReq := suite.NewAuthenticatedMilRequest(
			"GET",
			"/internal/users/logged_in",
			nil,
			*serviceMember,
		)

		authRecorder := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(authRecorder, authReq)

		suite.FatalTrue(suite.Equal(http.StatusOK, authRecorder.Code))

		return authRecorder.Result().Cookies()
	}

	// setUpRequestBody sets up the request body for the upload amended orders request, which needs a file.
	setUpRequestBody := func() (*bytes.Buffer, string) {
		buf := new(bytes.Buffer)

		writer := multipart.NewWriter(buf)

		ordersPDF := factory.FixtureOpen("filled-out-orders.pdf")

		defer ordersPDF.Close()

		part, formFileErr := writer.CreateFormFile("file", ordersPDF.Name())

		suite.FatalNoError(formFileErr)

		_, copyErr := io.Copy(part, ordersPDF)

		suite.FatalNoError(copyErr)

		// We need to close the writer so that the trailer is written, otherwise our request will fail.
		suite.FatalNoError(writer.Close())

		return buf, writer.FormDataContentType()
	}

	// setAuthCookies sets the cookies on the request and adds the CSRF token to the header
	setAuthCookies := func(req *http.Request, cookies []*http.Cookie) {
		for _, cookie := range cookies {
			if cookie.Name == auth.MaskedGorillaCSRFToken {
				req.Header.Set("X-CSRF-Token", cookie.Value)
			}

			req.AddCookie(cookie)
		}
	}

	suite.Run("Unauthorized upload to /orders/{ordersId}/upload_amended_orders by another service member", func() {
		move := factory.BuildSubmittedMove(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

		maliciousUser := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

		authCookies := getCookiesForServiceMember(&maliciousUser)

		endpointPath := fmt.Sprintf("/internal/orders/%s/upload_amended_orders", move.Orders.ID.String())

		body, contentType := setUpRequestBody()

		req := suite.NewAuthenticatedMilRequest("PATCH", endpointPath, body, maliciousUser)

		req.Header.Set("Content-Type", contentType)

		setAuthCookies(req, authCookies)

		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusNotFound, rr.Code)
	})

	suite.Run("Unauthorized upload to /orders/{ordersId}/upload_amended_orders by user that isn't logged in", func() {
		orders := factory.BuildOrderWithoutDefaults(suite.DB(), nil, nil)

		endpointPath := fmt.Sprintf("/internal/orders/%s/upload_amended_orders", orders.ID.String())

		body, contentType := setUpRequestBody()

		req := suite.NewMilRequest("PATCH", endpointPath, body)

		req.Header.Set("Content-Type", contentType)

		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)

		// Happens because we don't have a CSRF token, since they aren't logged in.
		suite.Equal(http.StatusForbidden, rr.Code)
	})

	suite.Run("Authorized upload to /orders/{ordersId}/upload_amended_orders", func() {
		move := factory.BuildSubmittedMove(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

		authCookies := getCookiesForServiceMember(&move.Orders.ServiceMember)

		endpointPath := fmt.Sprintf("/internal/orders/%s/upload_amended_orders", move.Orders.ID.String())

		body, contentType := setUpRequestBody()

		req := suite.NewAuthenticatedMilRequest("PATCH", endpointPath, body, move.Orders.ServiceMember)

		req.Header.Set("Content-Type", contentType)

		setAuthCookies(req, authCookies)

		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusCreated, rr.Code)
	})
}
