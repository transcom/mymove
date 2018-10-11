package user

import (
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// FetchServiceMember is the interface for a ServiceObject which fetches a ServiceMember object within a session
type FetchServiceMember interface {
	Execute(session *auth.Session, id uuid.UUID) (*models.ServiceMember, error)
}

type fetchServiceMemberService struct {
	smDB models.ServiceMemberDB
}

// NewFetchServiceMemberService is the DI provider to create a FetchServiceMember service object
func NewFetchServiceMemberService(smDB models.ServiceMemberDB) FetchServiceMember {
	return &fetchServiceMemberService{
		smDB,
	}
}

func (s *fetchServiceMemberService) Execute(session *auth.Session, id uuid.UUID) (*models.ServiceMember, error) {
	// TODO: Handle case where more than one user is authorized to modify serviceMember
	if session.IsMyApp() && id != session.ServiceMemberID {
		return nil, services.ErrFetchForbidden
	} else if session.IsTspApp() {
		isManaging, err := s.smDB.IsTspManagingShipment(session.TspUserID, id)
		if err != nil {
			return nil, err
		} else if !isManaging {
			return nil, services.ErrFetchForbidden
		}
	}
	return s.smDB.Fetch(id, true)
}
