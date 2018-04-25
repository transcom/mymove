package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth/context"
	queuemoveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/queue_moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForQueueMoveModel(user models.User, queueMove models.QueueMove) *internalmessages.QueueMove {
	queueMovePayload := internalmessages.QueueMove{
		ID:               fmtUUID(queueMove.ID),
		CreatedAt:        fmtDateTime(queueMove.CreatedAt),
		UpdatedAt:        fmtDateTime(queueMove.UpdatedAt),
		Edipi:            queueMove.Edipi,
		Rank:             queueMove.Rank,
		CustomerName:     queueMove.CustomerName,
		LocatorNumber:    queueMove.LocatorNumber,
		Status:           queueMove.Status,
		MoveType:         queueMove.MoveType,
		MoveDate:         fmtDate(queueMove.MoveDate),
		CustomerDeadline: fmtDate(queueMove.CustomerDeadline),
		LastModified:     queueMove.LastModified,
	}
	return &queueMovePayload
}

// IndexQueueNewMovesHandler returns a list of all queueMoves in the new moves queue
type IndexQueueNewMovesHandler HandlerContext

// Handle retrieves a list of all queueMoves in the system in the new moves queue
func (h IndexQueueNewMovesHandler) Handle(params queuemoveop.IndexQueueNewMovesParams) middleware.Responder {
	var response middleware.Responder

	queueMoves, err := models.GetMovesForUserID(h.db, user.ID)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		response = moveop.NewIndexMovesBadRequest()
	} else {
		movePayloads := make(internalmessages.IndexMovesPayload, len(moves))
		for i, move := range moves {
			movePayload := payloadForMoveModel(user, move)
			movePayloads[i] = &movePayload
		}
		response = moveop.NewIndexQueueMovesOK().WithPayload(movePayloads)
	}
	return response
}

// CreateQueueMoveHandler creates a new service member via POST /queueMove
type CreateQueueMoveHandler HandlerContext

// Handle ... creates a new QueueMove from a request payload
func (h CreateQueueMoveHandler) Handle(params queuemoveop.CreateQueueMoveParams) middleware.Responder {
	var response middleware.Responder
	verrs := validate.NewErrors()

	// Create a new queueMove for an authenticated user
	newQueueMove := models.QueueMove{
		UserID:                 user.ID,
		Edipi:                  params.CreateQueueMove.Edipi,
		Branch:                 params.CreateQueueMove.Branch,
		Rank:                   params.CreateQueueMove.Rank,
		FirstName:              params.CreateQueueMove.FirstName,
		MiddleName:             params.CreateQueueMove.MiddleName,
		LastName:               params.CreateQueueMove.LastName,
		Suffix:                 params.CreateQueueMove.Suffix,
		Telephone:              params.CreateQueueMove.Telephone,
		SecondaryTelephone:     params.CreateQueueMove.SecondaryTelephone,
		PersonalEmail:          stringFromEmail(params.CreateQueueMove.PersonalEmail),
		PhoneIsPreferred:       params.CreateQueueMove.PhoneIsPreferred,
		TextMessageIsPreferred: params.CreateQueueMove.TextMessageIsPreferred,
		EmailIsPreferred:       params.CreateQueueMove.EmailIsPreferred,
		ResidentialAddress:     residentialAddress,
		BackupMailingAddress:   backupMailingAddress,
		SocialSecurityNumber:   ssn,
	}
	smVerrs, err := models.SaveQueueMove(h.db, &newQueueMove)
	verrs.Append(smVerrs)
	if verrs.HasAny() {
		h.logger.Info("DB Validation", zap.Error(verrs))
		response = queuemoveop.NewCreateQueueMoveBadRequest()
	} else if err != nil {
		if err == models.ErrCreateViolatesUniqueConstraint {
			h.logger.Info("Attempted to create a second SM when one already exists")
			response = queuemoveop.NewCreateQueueMoveBadRequest()
		} else {
			h.logger.Error("DB Insertion", zap.Error(err))
			response = queuemoveop.NewCreateQueueMoveInternalServerError()
		}
	} else {
		queuemovePayload := payloadForQueueMoveModel(user, newQueueMove)
		response = queuemoveop.NewCreateQueueMoveCreated().WithPayload(queuemovePayload)
	}
	return response
}
