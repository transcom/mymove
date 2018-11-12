package user

import (
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/services"
)

type fetchServiceMemberService struct {
	smDB models.ServiceMemberDB
}

// NewFetchServiceMemberService is the DI provider to create a FetchServiceMember service objec
func NewFetchServiceMemberService(smDB models.ServiceMemberDB) services.FetchServiceMember {
	return &fetchServiceMemberService{
		smDB,
	}
}

func (s *fetchServiceMemberService) Execute(id uuid.UUID, session *server.Session) (*models.ServiceMember, error) {
	if session != nil {
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
	}
	return s.smDB.Fetch(id, true)
}
