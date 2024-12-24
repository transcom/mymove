package order

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *OrderServiceSuite) TestUpdateBillableWeightAsTOO() {
	suite.Run("Returns an error when order is not found", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())
		newAuthorizedWeight := int(10000)
		eTag := ""

		_, _, err := excessWeightRiskManager.UpdateBillableWeightAsTOO(suite.AppContextForTest(), nonexistentUUID, &newAuthorizedWeight, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when the etag does not match", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		order := factory.BuildMove(suite.DB(), nil, nil).Orders
		newAuthorizedWeight := int(10000)
		eTag := ""

		_, _, err := excessWeightRiskManager.UpdateBillableWeightAsTOO(suite.AppContextForTest(), order.ID, &newAuthorizedWeight, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), order.ID.String())
	})

	suite.Run("Updates the BillableWeight and approves the move when all fields are valid", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		now := time.Now()
		move := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ExcessWeightQualifiedAt: &now,
				},
			},
		}, nil)
		order := move.Orders
		newAuthorizedWeight := int(12345)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		updatedOrder, _, err := excessWeightRiskManager.UpdateBillableWeightAsTOO(suite.AppContextForTest(), order.ID, &newAuthorizedWeight, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)
		suite.NoError(err)
		err = suite.DB().Load(&orderInDB, "Moves", "Entitlement")
		suite.NoError(err)
		moveInDB := orderInDB.Moves[0]
		acknowledgedAt := moveInDB.ExcessWeightAcknowledgedAt

		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(newAuthorizedWeight, *updatedOrder.Entitlement.DBAuthorizedWeight)
		suite.Equal(newAuthorizedWeight, *orderInDB.Entitlement.DBAuthorizedWeight)
		suite.Equal(models.MoveStatusAPPROVED, moveInDB.Status)
		suite.WithinDuration(time.Now(), *acknowledgedAt, 2*time.Second)
	})

	suite.Run("Updates the BillableWeight but does not approve the move if unacknowledged amended orders exist", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		storer := storageTest.NewFakeS3Storage(true)
		userUploader, err := uploader.NewUserUploader(storer, 100*uploader.MB)
		suite.NoError(err)
		amendedDocument := factory.BuildDocument(suite.DB(), nil, nil)
		amendedUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    amendedDocument,
				LinkOnly: true,
			},
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   suite.AppContextForTest(),
				},
			},
		}, nil)

		amendedDocument.UserUploads = append(amendedDocument.UserUploads, amendedUpload)
		now := time.Now()
		approvalsRequestedMove := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ExcessWeightQualifiedAt: &now,
				},
			},
			{
				Model:    amendedDocument,
				LinkOnly: true,
				Type:     &factory.Documents.UploadedAmendedOrders,
			},
			{
				Model:    amendedDocument.ServiceMember,
				LinkOnly: true,
			},
		}, nil)
		order := approvalsRequestedMove.Orders
		newAuthorizedWeight := int(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		updatedOrder, _, err := excessWeightRiskManager.UpdateBillableWeightAsTOO(suite.AppContextForTest(), order.ID, &newAuthorizedWeight, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)
		suite.NoError(err)
		err = suite.DB().Load(&orderInDB, "Moves", "Entitlement")
		suite.NoError(err)
		moveInDB := orderInDB.Moves[0]
		acknowledgedAt := moveInDB.ExcessWeightAcknowledgedAt

		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(newAuthorizedWeight, *updatedOrder.Entitlement.DBAuthorizedWeight)
		suite.Equal(newAuthorizedWeight, *orderInDB.Entitlement.DBAuthorizedWeight)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moveInDB.Status)
		suite.WithinDuration(time.Now(), *acknowledgedAt, 2*time.Second)
	})

	suite.Run("Updates the BillableWeight but does not approve the move if unreviewed service items exist", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)

		_, _, move := suite.createServiceItem()
		order := move.Orders
		newAuthorizedWeight := int(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		updatedOrder, _, err := excessWeightRiskManager.UpdateBillableWeightAsTOO(suite.AppContextForTest(), order.ID, &newAuthorizedWeight, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)
		suite.NoError(err)
		err = suite.DB().Load(&orderInDB, "Moves", "Entitlement")
		suite.NoError(err)
		moveInDB := orderInDB.Moves[0]
		acknowledgedAt := moveInDB.ExcessWeightAcknowledgedAt

		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(newAuthorizedWeight, *updatedOrder.Entitlement.DBAuthorizedWeight)
		suite.Equal(newAuthorizedWeight, *orderInDB.Entitlement.DBAuthorizedWeight)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moveInDB.Status)
		suite.WithinDuration(time.Now(), *acknowledgedAt, 2*time.Second)
	})

	suite.Run("Updates the BillableWeight but does not acknowledge the risk if there is no excess weight risk", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		order := move.Orders
		newAuthorizedWeight := int(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		updatedOrder, _, err := excessWeightRiskManager.UpdateBillableWeightAsTOO(suite.AppContextForTest(), order.ID, &newAuthorizedWeight, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)
		suite.NoError(err)
		err = suite.DB().Load(&orderInDB, "Moves", "Entitlement")
		suite.NoError(err)
		moveInDB := orderInDB.Moves[0]

		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(newAuthorizedWeight, *updatedOrder.Entitlement.DBAuthorizedWeight)
		suite.Equal(newAuthorizedWeight, *orderInDB.Entitlement.DBAuthorizedWeight)
		suite.Equal(models.MoveStatusAPPROVED, moveInDB.Status)
		suite.Nil(moveInDB.ExcessWeightAcknowledgedAt)
	})

	suite.Run("Updates the BillableWeight but does not acknowledge the risk if the risk was already acknowledged", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		now := time.Now()
		move := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ExcessWeightAcknowledgedAt: &now,
					ExcessWeightQualifiedAt:    &now,
				},
			},
		}, nil)
		order := move.Orders
		newAuthorizedWeight := int(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		updatedOrder, _, err := excessWeightRiskManager.UpdateBillableWeightAsTOO(suite.AppContextForTest(), order.ID, &newAuthorizedWeight, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)
		suite.NoError(err)
		err = suite.DB().Load(&orderInDB, "Moves", "Entitlement")
		suite.NoError(err)
		moveInDB := orderInDB.Moves[0]

		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(newAuthorizedWeight, *updatedOrder.Entitlement.DBAuthorizedWeight)
		suite.Equal(newAuthorizedWeight, *orderInDB.Entitlement.DBAuthorizedWeight)
		suite.Equal(models.MoveStatusAPPROVED, moveInDB.Status)
	})
}

func (suite *OrderServiceSuite) TestUpdateBillableWeightAsTIO() {
	suite.Run("Returns an error when order is not found", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())
		newAuthorizedWeight := int(10000)
		newTIOremarks := "TIO remarks"
		eTag := ""

		_, _, err := excessWeightRiskManager.UpdateMaxBillableWeightAsTIO(suite.AppContextForTest(), nonexistentUUID, &newAuthorizedWeight, &newTIOremarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when the etag does not match", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		order := factory.BuildMove(suite.DB(), nil, nil).Orders
		newAuthorizedWeight := int(10000)
		newTIOremarks := "TIO remarks"
		eTag := ""

		_, _, err := excessWeightRiskManager.UpdateMaxBillableWeightAsTIO(suite.AppContextForTest(), order.ID, &newAuthorizedWeight, &newTIOremarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), order.ID.String())
	})

	suite.Run("Updates the MaxBillableWeight and TIO remarks and approves the move when all fields are valid", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		now := time.Now()
		move := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ExcessWeightQualifiedAt: &now,
				},
			},
		}, nil)
		order := move.Orders
		newAuthorizedWeight := int(12345)
		newTIOremarks := "TIO remarks"
		eTag := etag.GenerateEtag(order.UpdatedAt)

		updatedOrder, _, err := excessWeightRiskManager.UpdateMaxBillableWeightAsTIO(suite.AppContextForTest(), order.ID, &newAuthorizedWeight, &newTIOremarks, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)
		suite.NoError(err)
		err = suite.DB().Load(&orderInDB, "Moves", "Entitlement")
		suite.NoError(err)
		moveInDB := orderInDB.Moves[0]
		acknowledgedAt := moveInDB.ExcessWeightAcknowledgedAt

		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(newAuthorizedWeight, *updatedOrder.Entitlement.DBAuthorizedWeight)
		suite.Equal(newAuthorizedWeight, *orderInDB.Entitlement.DBAuthorizedWeight)
		suite.Equal(newTIOremarks, *moveInDB.TIORemarks)
		suite.Equal(models.MoveStatusAPPROVED, moveInDB.Status)
		suite.WithinDuration(time.Now(), *acknowledgedAt, 2*time.Second)
	})

	suite.Run("Updates the MaxBillableWeight and TIO remarks but does not approve the move if unacknowledged amended orders exist", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		storer := storageTest.NewFakeS3Storage(true)
		userUploader, err := uploader.NewUserUploader(storer, 100*uploader.MB)
		suite.NoError(err)
		amendedDocument := factory.BuildDocument(suite.DB(), nil, nil)
		amendedUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    amendedDocument,
				LinkOnly: true,
			},
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   suite.AppContextForTest(),
				},
			},
		}, nil)

		amendedDocument.UserUploads = append(amendedDocument.UserUploads, amendedUpload)
		now := time.Now()
		approvalsRequestedMove := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ExcessWeightQualifiedAt: &now,
				},
			},
			{
				Model:    amendedDocument,
				LinkOnly: true,
				Type:     &factory.Documents.UploadedAmendedOrders,
			},
			{
				Model:    amendedDocument.ServiceMember,
				LinkOnly: true,
			},
		}, nil)
		order := approvalsRequestedMove.Orders
		newAuthorizedWeight := int(10000)
		newTIOremarks := "TIO remarks"
		eTag := etag.GenerateEtag(order.UpdatedAt)

		updatedOrder, _, err := excessWeightRiskManager.UpdateMaxBillableWeightAsTIO(suite.AppContextForTest(), order.ID, &newAuthorizedWeight, &newTIOremarks, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)
		suite.NoError(err)
		err = suite.DB().Load(&orderInDB, "Moves", "Entitlement")
		suite.NoError(err)
		moveInDB := orderInDB.Moves[0]
		acknowledgedAt := moveInDB.ExcessWeightAcknowledgedAt

		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(newAuthorizedWeight, *updatedOrder.Entitlement.DBAuthorizedWeight)
		suite.Equal(newAuthorizedWeight, *orderInDB.Entitlement.DBAuthorizedWeight)
		suite.Equal(newTIOremarks, *moveInDB.TIORemarks)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moveInDB.Status)
		suite.WithinDuration(time.Now(), *acknowledgedAt, 2*time.Second)
	})

	suite.Run("Updates the MaxBillableWeight and TIO remarks but does not approve the move if unreviewed service items exist", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)

		_, _, move := suite.createServiceItem()
		order := move.Orders
		newAuthorizedWeight := int(10000)
		newTIOremarks := "TIO remarks"
		eTag := etag.GenerateEtag(order.UpdatedAt)

		updatedOrder, _, err := excessWeightRiskManager.UpdateMaxBillableWeightAsTIO(suite.AppContextForTest(), order.ID, &newAuthorizedWeight, &newTIOremarks, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)
		suite.NoError(err)
		err = suite.DB().Load(&orderInDB, "Moves", "Entitlement")
		suite.NoError(err)
		moveInDB := orderInDB.Moves[0]
		acknowledgedAt := moveInDB.ExcessWeightAcknowledgedAt

		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(newAuthorizedWeight, *updatedOrder.Entitlement.DBAuthorizedWeight)
		suite.Equal(newAuthorizedWeight, *orderInDB.Entitlement.DBAuthorizedWeight)
		suite.Equal(newTIOremarks, *moveInDB.TIORemarks)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moveInDB.Status)
		suite.WithinDuration(time.Now(), *acknowledgedAt, 2*time.Second)
	})

	suite.Run("Updates the MaxBillableWeight and TIO remarks but does not acknowledge the risk if there is no excess weight risk", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		order := move.Orders
		newAuthorizedWeight := int(10000)
		newTIOremarks := "TIO remarks"
		eTag := etag.GenerateEtag(order.UpdatedAt)

		updatedOrder, _, err := excessWeightRiskManager.UpdateMaxBillableWeightAsTIO(suite.AppContextForTest(), order.ID, &newAuthorizedWeight, &newTIOremarks, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)
		suite.NoError(err)
		err = suite.DB().Load(&orderInDB, "Moves", "Entitlement")
		suite.NoError(err)
		moveInDB := orderInDB.Moves[0]

		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(newAuthorizedWeight, *updatedOrder.Entitlement.DBAuthorizedWeight)
		suite.Equal(newAuthorizedWeight, *orderInDB.Entitlement.DBAuthorizedWeight)
		suite.Equal(newTIOremarks, *moveInDB.TIORemarks)
		suite.Equal(models.MoveStatusAPPROVED, moveInDB.Status)
		suite.Nil(moveInDB.ExcessWeightAcknowledgedAt)
	})

	suite.Run("Updates the MaxBillableWeight and TIO remarks but does not acknowledge the risk if the risk was already acknowledged", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		now := time.Now()
		move := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ExcessWeightAcknowledgedAt: &now,
					ExcessWeightQualifiedAt:    &now,
				},
			},
		}, nil)
		order := move.Orders
		newAuthorizedWeight := int(10000)
		newTIOremarks := "TIO remarks"
		eTag := etag.GenerateEtag(order.UpdatedAt)

		updatedOrder, _, err := excessWeightRiskManager.UpdateMaxBillableWeightAsTIO(suite.AppContextForTest(), order.ID, &newAuthorizedWeight, &newTIOremarks, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)
		suite.NoError(err)
		err = suite.DB().Load(&orderInDB, "Moves", "Entitlement")
		suite.NoError(err)
		moveInDB := orderInDB.Moves[0]

		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(newAuthorizedWeight, *updatedOrder.Entitlement.DBAuthorizedWeight)
		suite.Equal(newAuthorizedWeight, *orderInDB.Entitlement.DBAuthorizedWeight)
		suite.Equal(newTIOremarks, *moveInDB.TIORemarks)
		suite.Equal(models.MoveStatusAPPROVED, moveInDB.Status)
	})
}

func (suite *OrderServiceSuite) TestAcknowledgeExcessWeightRisk() {
	suite.Run("Returns an error when move is not found", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())
		eTag := ""

		_, err := excessWeightRiskManager.AcknowledgeExcessWeightRisk(suite.AppContextForTest(), nonexistentUUID, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when the etag does not match", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		move := factory.BuildMove(suite.DB(), nil, nil)
		order := move.Orders
		eTag := ""

		_, err := excessWeightRiskManager.AcknowledgeExcessWeightRisk(suite.AppContextForTest(), order.ID, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), move.ID.String())
	})

	suite.Run("Updates the ExcessWeightAcknowledgedAt field and approves the move when all fields are valid", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		now := time.Now()
		move := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ExcessWeightQualifiedAt: &now,
				},
			},
		}, nil)
		order := move.Orders
		eTag := etag.GenerateEtag(move.UpdatedAt)

		updatedMove, err := excessWeightRiskManager.AcknowledgeExcessWeightRisk(suite.AppContextForTest(), order.ID, eTag)
		suite.NoError(err)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)

		suite.Equal(move.ID.String(), updatedMove.ID.String())
		suite.Equal(models.MoveStatusAPPROVED, moveInDB.Status)
		suite.Equal(models.MoveStatusAPPROVED, updatedMove.Status)
		suite.NotNil(moveInDB.ExcessWeightAcknowledgedAt)
		suite.NotNil(updatedMove.ExcessWeightAcknowledgedAt)
	})

	suite.Run("Updates the ExcessWeightAcknowledgedAt field but does not approve the move if unacknowledged amended orders exist", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		storer := storageTest.NewFakeS3Storage(true)
		userUploader, err := uploader.NewUserUploader(storer, 100*uploader.MB)
		suite.NoError(err)
		amendedDocument := factory.BuildDocument(suite.DB(), nil, nil)
		amendedUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    amendedDocument,
				LinkOnly: true,
			},
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   suite.AppContextForTest(),
				},
			},
		}, nil)

		amendedDocument.UserUploads = append(amendedDocument.UserUploads, amendedUpload)
		now := time.Now()
		approvalsRequestedMove := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ExcessWeightQualifiedAt: &now,
				},
			},
			{
				Model:    amendedDocument,
				LinkOnly: true,
				Type:     &factory.Documents.UploadedAmendedOrders,
			},
			{
				Model:    amendedDocument.ServiceMember,
				LinkOnly: true,
			},
		}, nil)
		order := approvalsRequestedMove.Orders
		eTag := etag.GenerateEtag(approvalsRequestedMove.UpdatedAt)

		updatedMove, err := excessWeightRiskManager.AcknowledgeExcessWeightRisk(suite.AppContextForTest(), order.ID, eTag)
		suite.NoError(err)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, approvalsRequestedMove.ID)
		suite.NoError(err)

		suite.Equal(approvalsRequestedMove.ID.String(), updatedMove.ID.String())
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moveInDB.Status)
		suite.NotNil(moveInDB.ExcessWeightAcknowledgedAt)
	})

	suite.Run("Updates the ExcessWeightAcknowledgedAt field but does not approve the move if unreviewed service items exist", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)

		_, _, move := suite.createServiceItem()
		eTag := etag.GenerateEtag(move.UpdatedAt)
		order := move.Orders

		updatedMove, err := excessWeightRiskManager.AcknowledgeExcessWeightRisk(suite.AppContextForTest(), order.ID, eTag)
		suite.NoError(err)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)

		suite.Equal(move.ID.String(), updatedMove.ID.String())
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moveInDB.Status)
		suite.NotNil(moveInDB.ExcessWeightAcknowledgedAt)
	})

	suite.Run("Does not update the ExcessWeightAcknowledgedAt field if there is no risk of excess weight", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)

		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(move.UpdatedAt)
		order := move.Orders

		updatedMove, err := excessWeightRiskManager.AcknowledgeExcessWeightRisk(suite.AppContextForTest(), order.ID, eTag)
		suite.NoError(err)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)

		suite.Equal(move.ID.String(), updatedMove.ID.String())
		suite.Equal(models.MoveStatusAPPROVED, moveInDB.Status)
		suite.Nil(moveInDB.ExcessWeightAcknowledgedAt)
		suite.Nil(updatedMove.ExcessWeightAcknowledgedAt)
	})

	suite.Run("Does not update the ExcessWeightAcknowledgedAt field if the risk was already acknowledged", func() {
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)

		date := time.Now().Add(30 * time.Minute)
		move := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ExcessWeightAcknowledgedAt: &date,
					ExcessWeightQualifiedAt:    &date,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(move.UpdatedAt)
		order := move.Orders

		updatedMove, err := excessWeightRiskManager.AcknowledgeExcessWeightRisk(suite.AppContextForTest(), order.ID, eTag)
		suite.NoError(err)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)

		suite.Equal(move.ID.String(), updatedMove.ID.String())
		suite.Equal(models.MoveStatusAPPROVED, moveInDB.Status)
		suite.WithinDuration(*move.ExcessWeightAcknowledgedAt, *moveInDB.ExcessWeightAcknowledgedAt, 1*time.Second)
	})
}

func (suite *OrderServiceSuite) createServiceItem() (string, models.MTOServiceItem, models.Move) {
	now := time.Now()
	move := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				ExcessWeightQualifiedAt: &now,
			},
		},
	}, nil)
	serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

	return eTag, serviceItem, move
}
