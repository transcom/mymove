import React from 'react';
import { connect } from 'react-redux';
import { arrayOf, bool } from 'prop-types';
import { Alert } from '@trussworks/react-uswds';

import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import ContactInfoDisplay from 'components/Customer/Profile/ContactInfoDisplay/ContactInfoDisplay';
import { BackupContactShape, OrdersShape, ServiceMemberShape } from 'types/customerShapes';
import {
  selectServiceMemberFromLoggedInUser,
  selectMoveIsInDraft,
  selectCurrentOrders,
  selectBackupContacts,
} from 'store/entities/selectors';
import SectionWrapper from 'components/Customer/SectionWrapper';
import ServiceInfoDisplay from 'components/Customer/Review/ServiceInfoDisplay/ServiceInfoDisplay';
import { customerRoutes } from 'constants/routes';
import formStyles from 'styles/form.module.scss';
import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';

const Profile = ({ serviceMember, currentOrders, currentBackupContacts, moveIsInDraft }) => {
  const showMessages = currentOrders.id && !moveIsInDraft;
  const rank = currentOrders.grade ?? serviceMember.rank;
  const originStation = currentOrders.origin_duty_location ?? serviceMember.current_station;
  const transportationOfficePhoneLines = originStation?.transportation_office?.phone_lines;
  const transportationOfficePhone = transportationOfficePhoneLines ? transportationOfficePhoneLines[0] : '';
  const backupContact = {
    name: currentBackupContacts[0]?.name || '',
    telephone: currentBackupContacts[0]?.telephone || '',
    email: currentBackupContacts[0]?.email || '',
  };

  return (
    <div className="grid-container usa-prose">
      <ConnectedFlashMessage />
      <div className="grid-row">
        <div className="grid-col-12">
          <h1>Profile</h1>
          {showMessages && <Alert type="info">Contact your movers if you need to make changes to your move.</Alert>}
          <SectionWrapper className={formStyles.formSection}>
            <ContactInfoDisplay
              telephone={serviceMember?.telephone || ''}
              secondaryTelephone={serviceMember?.secondary_telephone || ''}
              personalEmail={serviceMember?.personal_email || ''}
              emailIsPreferred={serviceMember?.email_is_preferred}
              phoneIsPreferred={serviceMember?.phone_is_preferred}
              residentialAddress={serviceMember?.residential_address || ''}
              backupMailingAddress={serviceMember?.backup_mailing_address || ''}
              backupContact={backupContact}
              editURL={customerRoutes.CONTACT_INFO_EDIT_PATH}
            />
          </SectionWrapper>
          <SectionWrapper className={formStyles.formSection}>
            <ServiceInfoDisplay
              firstName={serviceMember?.first_name || ''}
              lastName={serviceMember?.last_name || ''}
              originDutyStationName={originStation?.name || ''}
              originTransportationOfficeName={originStation?.transportation_office?.name || ''}
              originTransportationOfficePhone={transportationOfficePhone}
              affiliation={ORDERS_BRANCH_OPTIONS[serviceMember?.affiliation] || ''}
              rank={ORDERS_RANK_OPTIONS[rank] || ''}
              edipi={serviceMember?.edipi || ''}
              editURL={customerRoutes.SERVICE_INFO_EDIT_PATH}
              isEditable={moveIsInDraft}
              showMessage={showMessages}
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
  moveIsInDraft: bool.isRequired,
};

function mapStateToProps(state) {
  return {
    serviceMember: selectServiceMemberFromLoggedInUser(state),
    moveIsInDraft: selectMoveIsInDraft(state),
    currentOrders: selectCurrentOrders(state) || {},
    currentBackupContacts: selectBackupContacts(state),
  };
}

export default connect(mapStateToProps)(Profile);
