package services

import (
	"time"

	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/models"
)

// StorageInTransitCreator is the service object for creating a Storage In Transit
//go:generate mockery -name StorageInTransitCreator
type StorageInTransitCreator interface {
	CreateStorageInTransit(storageInTransitPayload apimessages.StorageInTransit, shipmentID uuid.UUID, session *auth.Session) (*models.StorageInTransit, *validate.Errors, error)
}

// StorageInTransitsIndexer is the service object for retrieving a collection of StorageInTransits
//go:generate mockery -name StorageInTransitsIndexer
type StorageInTransitsIndexer interface {
	IndexStorageInTransits(shipmentID uuid.UUID, session *auth.Session) ([]models.StorageInTransit, error)
}

// StorageInTransitApprover is the service object for approving a Storage In Transit
//go:generate mockery -name StorageInTransitApprover
type StorageInTransitApprover interface {
	ApproveStorageInTransit(payload apimessages.StorageInTransitApprovalPayload, shipmentID uuid.UUID, session *auth.Session, storageInTransitID uuid.UUID) (*models.StorageInTransit, *validate.Errors, error)
}

// StorageInTransitDenier is the service object for denying a Storage In Transit
//go:generate mockery -name StorageInTransitDenier
type StorageInTransitDenier interface {
	DenyStorageInTransit(payload apimessages.StorageInTransitDenialPayload, shipmentID uuid.UUID, session *auth.Session, storageInTransitID uuid.UUID) (*models.StorageInTransit, *validate.Errors, error)
}

// StorageInTransitInSITPlacer is the object for placing a Storage In Transit into SIT status
//go:generate mockery -name StorageInTransitInSITPlacer
type StorageInTransitInSITPlacer interface {
	PlaceIntoSITStorageInTransit(payload apimessages.StorageInTransitInSitPayload, shipmentID uuid.UUID, session *auth.Session, storageInTransitID uuid.UUID) (*models.StorageInTransit, *validate.Errors, error)
}

// StorageInTransitDeliverer is the service object for delivering a Storage In Transit
//go:generate mockery -name StorageInTransitDeliverer
type StorageInTransitDeliverer interface {
	DeliverStorageInTransits(shipmentID uuid.UUID, tspID uuid.UUID) ([]models.StorageInTransit, *validate.Errors, error)
}

// StorageInTransitReleaser is the service object for releasing a Storage In Transit
//go:generate mockery -name StorageInTransitReleaser
type StorageInTransitReleaser interface {
	ReleaseStorageInTransit(payload apimessages.StorageInTransitReleasePayload, shipmentID uuid.UUID, session *auth.Session, storageInTransitID uuid.UUID) (*models.StorageInTransit, *validate.Errors, error)
}

// StorageInTransitDeleter is the service object for deleting a Storage In Transit
//go:generate mockery -name StorageInTransitDeleter
type StorageInTransitDeleter interface {
	DeleteStorageInTransit(shipmentID uuid.UUID, storageInTransitID uuid.UUID, session *auth.Session) (*models.StorageInTransit, error)
}

// StorageInTransitPatcher is the service object for editing a Storage In Transit
//go:generate mockery -name StorageInTransitPatcher
type StorageInTransitPatcher interface {
	PatchStorageInTransit(payload apimessages.StorageInTransit, shipmentID uuid.UUID, storageInTransitID uuid.UUID, session *auth.Session) (*models.StorageInTransit, *validate.Errors, error)
}

// StorageInTransitByIDFetcher is the service object for fetching a Storage In Transit
//go:generate mockery -name StorageInTransitByIDFetcher
type StorageInTransitByIDFetcher interface {
	FetchStorageInTransitByID(storageInTransitID uuid.UUID, shipmentID uuid.UUID, session *auth.Session) (*models.StorageInTransit, error)
}

// StorageInTransitNumberGenerator is an interface for generating a storage in transit number
type StorageInTransitNumberGenerator interface {
	GenerateStorageInTransitNumber(placeInSitTime time.Time) (string, error)
}
