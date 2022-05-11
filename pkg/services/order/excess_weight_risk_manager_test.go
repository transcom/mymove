package order

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *OrderServiceSuite) TestUpdateBillableWeightAsTOO() {
	suite.T().Run("Returns an error when order is not found", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())
		newAuthorizedWeight := int(10000)
		eTag := ""

		_, _, err := excessWeightRiskManager.UpdateBillableWeightAsTOO(suite.AppContextForTest(), nonexistentUUID, &newAuthorizedWeight, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when the etag does not match", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		order := testdatagen.MakeDefaultMove(suite.DB()).Orders
		newAuthorizedWeight := int(10000)
		eTag := ""

		_, _, err := excessWeightRiskManager.UpdateBillableWeightAsTOO(suite.AppContextForTest(), order.ID, &newAuthorizedWeight, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), order.ID.String())
	})

	suite.T().Run("Updates the BillableWeight and approves the move when all fields are valid", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		now := time.Now()
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{ExcessWeightQualifiedAt: &now},
		})
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

	suite.T().Run("Updates the BillableWeight but does not approve the move if unacknowledged amended orders exist", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		storer := storageTest.NewFakeS3Storage(true)
		userUploader, err := uploader.NewUserUploader(storer, 100*uploader.MB)
		suite.NoError(err)
		amendedDocument := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{})
		amendedUpload := testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
			UserUpload: models.UserUpload{
				DocumentID: &amendedDocument.ID,
				Document:   amendedDocument,
				UploaderID: amendedDocument.ServiceMember.UserID,
			},
			UserUploader: userUploader,
		})

		amendedDocument.UserUploads = append(amendedDocument.UserUploads, amendedUpload)
		now := time.Now()
		approvalsRequestedMove := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				UploadedAmendedOrders:   &amendedDocument,
				UploadedAmendedOrdersID: &amendedDocument.ID,
				ServiceMember:           amendedDocument.ServiceMember,
				ServiceMemberID:         amendedDocument.ServiceMemberID,
			},
			Move: models.Move{ExcessWeightQualifiedAt: &now},
		})
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

	suite.T().Run("Updates the BillableWeight but does not approve the move if unreviewed service items exist", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
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

	suite.T().Run("Updates the BillableWeight but does not acknowledge the risk if there is no excess weight risk", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{})
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

	suite.T().Run("Updates the BillableWeight but does not acknowledge the risk if the risk was already acknowledged", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		now := time.Now()
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				ExcessWeightAcknowledgedAt: &now,
				ExcessWeightQualifiedAt:    &now,
			},
		})
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
	suite.T().Run("Returns an error when order is not found", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())
		newAuthorizedWeight := int(10000)
		newTIOremarks := "TIO remarks"
		eTag := ""

		_, _, err := excessWeightRiskManager.UpdateMaxBillableWeightAsTIO(suite.AppContextForTest(), nonexistentUUID, &newAuthorizedWeight, &newTIOremarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when the etag does not match", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		order := testdatagen.MakeDefaultMove(suite.DB()).Orders
		newAuthorizedWeight := int(10000)
		newTIOremarks := "TIO remarks"
		eTag := ""

		_, _, err := excessWeightRiskManager.UpdateMaxBillableWeightAsTIO(suite.AppContextForTest(), order.ID, &newAuthorizedWeight, &newTIOremarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), order.ID.String())
	})

	suite.T().Run("Updates the MaxBillableWeight and TIO remarks and approves the move when all fields are valid", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		now := time.Now()
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{ExcessWeightQualifiedAt: &now},
		})
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

	suite.T().Run("Updates the MaxBillableWeight and TIO remarks but does not approve the move if unacknowledged amended orders exist", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		storer := storageTest.NewFakeS3Storage(true)
		userUploader, err := uploader.NewUserUploader(storer, 100*uploader.MB)
		suite.NoError(err)
		amendedDocument := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{})
		amendedUpload := testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
			UserUpload: models.UserUpload{
				DocumentID: &amendedDocument.ID,
				Document:   amendedDocument,
				UploaderID: amendedDocument.ServiceMember.UserID,
			},
			UserUploader: userUploader,
		})

		amendedDocument.UserUploads = append(amendedDocument.UserUploads, amendedUpload)
		now := time.Now()
		approvalsRequestedMove := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				UploadedAmendedOrders:   &amendedDocument,
				UploadedAmendedOrdersID: &amendedDocument.ID,
				ServiceMember:           amendedDocument.ServiceMember,
				ServiceMemberID:         amendedDocument.ServiceMemberID,
			},
			Move: models.Move{ExcessWeightQualifiedAt: &now},
		})
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

	suite.T().Run("Updates the MaxBillableWeight and TIO remarks but does not approve the move if unreviewed service items exist", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
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

	suite.T().Run("Updates the MaxBillableWeight and TIO remarks but does not acknowledge the risk if there is no excess weight risk", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{})
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

	suite.T().Run("Updates the MaxBillableWeight and TIO remarks but does not acknowledge the risk if the risk was already acknowledged", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		now := time.Now()
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				ExcessWeightAcknowledgedAt: &now,
				ExcessWeightQualifiedAt:    &now,
			},
		})
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
	suite.T().Run("Returns an error when move is not found", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())
		eTag := ""

		_, err := excessWeightRiskManager.AcknowledgeExcessWeightRisk(suite.AppContextForTest(), nonexistentUUID, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when the etag does not match", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		move := testdatagen.MakeDefaultMove(suite.DB())
		order := move.Orders
		eTag := ""

		_, err := excessWeightRiskManager.AcknowledgeExcessWeightRisk(suite.AppContextForTest(), order.ID, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), move.ID.String())
	})

	suite.T().Run("Updates the ExcessWeightAcknowledgedAt field and approves the move when all fields are valid", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		now := time.Now()
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{ExcessWeightQualifiedAt: &now},
		})
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

	suite.T().Run("Updates the ExcessWeightAcknowledgedAt field but does not approve the move if unacknowledged amended orders exist", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)
		storer := storageTest.NewFakeS3Storage(true)
		userUploader, err := uploader.NewUserUploader(storer, 100*uploader.MB)
		suite.NoError(err)
		amendedDocument := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{})
		amendedUpload := testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
			UserUpload: models.UserUpload{
				DocumentID: &amendedDocument.ID,
				Document:   amendedDocument,
				UploaderID: amendedDocument.ServiceMember.UserID,
			},
			UserUploader: userUploader,
		})

		amendedDocument.UserUploads = append(amendedDocument.UserUploads, amendedUpload)
		now := time.Now()
		approvalsRequestedMove := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				UploadedAmendedOrders:   &amendedDocument,
				UploadedAmendedOrdersID: &amendedDocument.ID,
				ServiceMember:           amendedDocument.ServiceMember,
				ServiceMemberID:         amendedDocument.ServiceMemberID,
			},
			Move: models.Move{ExcessWeightQualifiedAt: &now},
		})
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

	suite.T().Run("Updates the ExcessWeightAcknowledgedAt field but does not approve the move if unreviewed service items exist", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
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

	suite.T().Run("Does not update the ExcessWeightAcknowledgedAt field if there is no risk of excess weight", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)

		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{})
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

	suite.T().Run("Does not update the ExcessWeightAcknowledgedAt field if the risk was already acknowledged", func(t *testing.T) {
		moveRouter := moverouter.NewMoveRouter()
		excessWeightRiskManager := NewExcessWeightRiskManager(moveRouter)

		date := time.Now().Add(30 * time.Minute)
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				ExcessWeightAcknowledgedAt: &date,
				ExcessWeightQualifiedAt:    &date,
			},
		})
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
	move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{ExcessWeightQualifiedAt: &now},
	})

	serviceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: move,
	})

	eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

	return eTag, serviceItem, move
}
