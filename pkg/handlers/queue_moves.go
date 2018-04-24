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
		ID:                      fmtUUID(queueMove.ID),
		CreatedAt:               fmtDateTime(queueMove.CreatedAt),
		UpdatedAt:               fmtDateTime(queueMove.UpdatedAt),
		Edipi:                   queueMove.Edipi,
		Branch:                  queueMove.Branch,
		Rank:                    queueMove.Rank,
		FirstName:               queueMove.FirstName,
		MiddleName:              queueMove.MiddleName,
		LastName:                queueMove.LastName,
		Suffix:                  queueMove.Suffix,
		Telephone:               queueMove.Telephone,
		SecondaryTelephone:      queueMove.SecondaryTelephone,
		PhoneIsPreferred:        queueMove.PhoneIsPreferred,
		TextMessageIsPreferred:  queueMove.TextMessageIsPreferred,
		EmailIsPreferred:        queueMove.EmailIsPreferred,
		ResidentialAddress:      payloadForAddressModel(queueMove.ResidentialAddress),
		BackupMailingAddress:    payloadForAddressModel(queueMove.BackupMailingAddress),
		HasSocialSecurityNumber: fmtBool(queueMove.SocialSecurityNumberID != nil),
		IsProfileComplete:       fmtBool(queueMove.IsProfileComplete()),
	}
	return &queueMovePayload
}

// CreateQueueMoveHandler creates a new service member via POST /queueMove
type CreateQueueMoveHandler HandlerContext

// Handle ... creates a new QueueMove from a request payload
func (h CreateQueueMoveHandler) Handle(params queuemoveop.CreateQueueMoveParams) middleware.Responder {
	var response middleware.Responder
	residentialAddress := addressModelFromPayload(params.CreateQueueMove.ResidentialAddress)
	backupMailingAddress := addressModelFromPayload(params.CreateQueueMove.BackupMailingAddress)

	ssnString := params.CreateQueueMove.SocialSecurityNumber
	var ssn *models.SocialSecurityNumber
	verrs := validate.NewErrors()
	if ssnString != nil {
		var err error
		ssn, verrs, err = models.BuildSocialSecurityNumber(ssnString.String())
		if err != nil {
			h.logger.Error("Unexpected error building SSN model", zap.Error(err))
			return queuemoveop.NewCreateQueueMoveInternalServerError()
		}
		// if there are any validation errors, they will get rolled up with the rest of them.
	}

	// User should always be populated by middleware
	user, _ := context.GetUser(params.HTTPRequest.Context())

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

// ShowQueueMoveHandler returns a queueMove for a user and service member ID
type ShowQueueMoveHandler HandlerContext

// Handle retrieves a service member in the system belonging to the logged in user given service member ID
func (h ShowQueueMoveHandler) Handle(params queuemoveop.ShowQueueMoveParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := context.GetUser(params.HTTPRequest.Context())

	queueMoveID, _ := uuid.FromString(params.QueueMoveID.String())
	queueMove, err := models.FetchQueueMove(h.db, user, queueMoveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	queueMovePayload := payloadForQueueMoveModel(user, queueMove)
	return queuemoveop.NewShowQueueMoveOK().WithPayload(queueMovePayload)
}
