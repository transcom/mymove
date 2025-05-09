package lockmove

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/lib/pq"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveLocker struct {
}

// NewMoveLocker creates a new moveLocker service
func NewMoveLocker() services.MoveLocker {
	return &moveLocker{}
}

// LockMove updates a move with relevant values of who has a move locked and the expiration of the lock pending it isn't unlocked before then
func (m moveLocker) LockMove(appCtx appcontext.AppContext, move *models.Move, officeUserID uuid.UUID) (*models.Move, error) {

	if officeUserID == uuid.Nil {
		return &models.Move{}, apperror.NewQueryError("OfficeUserID", nil, "No office user provided in request to lock move")
	}

	// fetching office user
	officeUser, err := models.FetchOfficeUserByID(appCtx.DB(), officeUserID)
	if err != nil {
		return nil, err
	}

	// fetching transportation office that office user belongs to
	// this data will be used to display to read-only viewers in the UI
	var transportationOffice models.TransportationOffice
	err = appCtx.DB().Q().
		Join("office_users", "transportation_offices.id = office_users.transportation_office_id").
		Join("addresses", "transportation_offices.address_id = addresses.id").
		Where("office_users.id = ?", officeUserID).
		EagerPreload("Address", "Address.Country").
		First(&transportationOffice)

	if err != nil {
		return nil, err
	}

	if move.LockedByOfficeUserID != models.UUIDPointer(officeUserID) {
		move.LockedByOfficeUserID = models.UUIDPointer(officeUserID)
	}

	if officeUser != nil {
		move.LockedByOfficeUser = officeUser
	}

	if transportationOffice.ID != uuid.Nil {
		move.LockedByOfficeUser.TransportationOffice = transportationOffice
	}

	// the lock will have a default expiration time of 30 minutes from initial opening
	// this will reset with valid user activity
	now := time.Now()
	expirationTime := now.Add(30 * time.Minute)
	move.LockExpiresAt = &expirationTime

	// Store move before update
	var moveBeforeUpdate = *move

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if err := appCtx.DB().RawQuery("UPDATE moves SET locked_by=?, lock_expires_at=?, updated_at=? WHERE id=?", officeUserID, expirationTime, moveBeforeUpdate.UpdatedAt, move.ID).Exec(); err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return move, nil
}

// BulkLockMove updates a multiple moves with relevant values of who has a move locked and the expiration of the lock pending it isn't unlocked before then
func (m moveLocker) LockMoves(appCtx appcontext.AppContext, moveIds []uuid.UUID, officeUserID uuid.UUID) error {

	if officeUserID == uuid.Nil {
		return apperror.NewQueryError("OfficeUserID", nil, "No office user provided in request to lock move")
	}

	// fetching office user
	officeUser, err := models.FetchOfficeUserByID(appCtx.DB(), officeUserID)
	if err != nil {
		return err
	}

	now := time.Now()
	expirationTime := now.Add(30 * time.Minute)

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if err := appCtx.DB().RawQuery(
			"UPDATE moves SET locked_by=?, lock_expires_at=? WHERE id=ANY(?)",
			officeUser.ID, expirationTime, pq.Array(moveIds),
		).Exec(); err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return transactionError
	}

	return nil
}
