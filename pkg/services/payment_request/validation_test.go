package paymentrequest

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PaymentRequestServiceSuite) Test_validatePaymentRequest() {
	type args struct {
		appCtx            appcontext.AppContext
		newPaymentRequest models.PaymentRequest
		oldPaymentRequest *models.PaymentRequest
		checks            []paymentRequestValidator
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid payment request",
			args: args{
				appCtx: appcontext.NewAppContext(suite.DB(), suite.Logger(), nil),
				newPaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000000")),
					Status:          models.PaymentRequestStatusPending,
				},
				oldPaymentRequest: nil,
				checks: []paymentRequestValidator{
					paymentRequestValidatorFunc(func(appCtx appcontext.AppContext, newPaymentRequest models.PaymentRequest, oldPaymentRequest *models.PaymentRequest) error {
						return nil
					}),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := validatePaymentRequest(tt.args.appCtx, tt.args.newPaymentRequest, tt.args.oldPaymentRequest, tt.args.checks...)
			if (err != nil) != tt.wantErr {
				suite.T().Errorf("validatePaymentRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
