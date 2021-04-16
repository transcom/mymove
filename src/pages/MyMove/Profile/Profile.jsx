import React from 'react';
import { connect } from 'react-redux';
import { arrayOf } from 'prop-types';

import ContactInfoDisplay from 'components/Customer/Profile/ContactInfoDisplay/ContactInfoDisplay';
import { BackupContactShape, OrdersShape, ServiceMemberShape } from 'types/customerShapes';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentOrders,
  selectCurrentMove,
  selectBackupContacts,
} from 'store/entities/selectors';
import SectionWrapper from 'components/Customer/SectionWrapper';
import ServiceInfoDisplay from 'components/Customer/Review/ServiceInfoDisplay/ServiceInfoDisplay';
import { customerRoutes } from 'constants/routes';
import formStyles from 'styles/form.module.scss';

const Profile = ({ serviceMember, currentOrders, currentBackupContacts }) => {
  const rank = currentOrders ? currentOrders.grade : serviceMember.rank;
  const currentStation = currentOrders ? currentOrders.origin_duty_station : serviceMember.current_station;
  const backupContact = {
    name: currentBackupContacts[0]?.name || '',
    telephone: currentBackupContacts[0]?.telephone || '',
    email: currentBackupContacts[0]?.email || '',
  };

  return (
    <div className="grid-container usa-prose">
      <div className="grid-row">
        <div className="grid-col-12">
          <h1>Profile</h1>
          <SectionWrapper className={formStyles.formSection}>
            <ContactInfoDisplay
              telephone={serviceMember?.telephone || ''}
              personalEmail={serviceMember?.personal_email || ''}
              emailIsPreferred={serviceMember?.email_is_preferred}
              phoneIsPreferred={serviceMember?.phone_is_preferred}
              residentialAddress={serviceMember?.residential_address || ''}
              backupMailingAddress={serviceMember?.backup_mailing_address || ''}
              backupContact={backupContact}
              editURL={customerRoutes.EDIT_PROFILE_PATH}
            />
          </SectionWrapper>
          <SectionWrapper className={formStyles.formSection}>
            <ServiceInfoDisplay
              firstName={serviceMember?.first_name || ''}
              lastName={serviceMember?.last_name || ''}
              currentDutyStationName={currentStation?.name || ''}
              affiliation={serviceMember?.affiliation || ''}
              rank={rank || ''}
              edipi={serviceMember?.edipi || ''}
              editURL={customerRoutes.EDIT_PROFILE_PATH}
            />
          </SectionWrapper>
        </div>
      </div>
    </div>
  );
};

Profile.propTypes = {
  serviceMember: ServiceMemberShape.isRequired,
  currentOrders: OrdersShape.isRequired,
  currentBackupContacts: arrayOf(BackupContactShape).isRequired,
};

function mapStateToProps(state) {
  return {
    serviceMember: selectServiceMemberFromLoggedInUser(state),
    move: selectCurrentMove(state) || {},
    currentOrders: selectCurrentOrders(state),
    currentBackupContacts: selectBackupContacts(state),
  };
}

export default connect(mapStateToProps)(Profile);
