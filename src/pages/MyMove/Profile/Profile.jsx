import React from 'react';
import { connect } from 'react-redux';
import { arrayOf, bool } from 'prop-types';
import { Alert } from '@trussworks/react-uswds';
import { Link } from 'react-router-dom';

import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import ContactInfoDisplay from 'components/Customer/Profile/ContactInfoDisplay/ContactInfoDisplay';
import { BackupContactShape, OrdersShape, ServiceMemberShape } from 'types/customerShapes';
import {
  selectServiceMemberFromLoggedInUser,
  selectMoveIsInDraft,
  selectCurrentOrders,
  selectBackupContacts,
  selectOktaUser,
} from 'store/entities/selectors';
import SectionWrapper from 'components/Customer/SectionWrapper';
import ServiceInfoDisplay from 'components/Customer/Review/ServiceInfoDisplay/ServiceInfoDisplay';
import OktaInfoDisplay from 'components/Customer/Profile/OktaInfoDisplay/OktaInfoDisplay';
import { customerRoutes, generalRoutes } from 'constants/routes';
import formStyles from 'styles/form.module.scss';
import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';
import { OktaUserInfoShape } from 'types/user';

const Profile = ({ serviceMember, currentOrders, currentBackupContacts, moveIsInDraft, oktaUser }) => {
  const showMessages = currentOrders.id && !moveIsInDraft;
  const rank = currentOrders.grade ?? serviceMember.rank;
  const originDutyLocation = currentOrders.origin_duty_location ?? serviceMember.current_location;
  const transportationOfficePhoneLines = originDutyLocation?.transportation_office?.phone_lines;
  const transportationOfficePhone = transportationOfficePhoneLines ? transportationOfficePhoneLines[0] : '';
  const backupContact = {
    name: currentBackupContacts[0]?.name || '',
    telephone: currentBackupContacts[0]?.telephone || '',
    email: currentBackupContacts[0]?.email || '',
  };

  // displays the profile data for MilMove & Okta
  // Profile w/contact info for servicemember & backup contact
  // Service info that displays name, branch, rank, DoDID/EDIPI, and current duty location
  // okta profile information: username, email, first name, last name, and DoDID/EDIPI
  return (
    <div className="grid-container usa-prose">
      <ConnectedFlashMessage />
      <div className="grid-row">
        <div className="grid-col-12">
          <Link to={generalRoutes.HOME_PATH}>Return to Move</Link>
          <h1>Profile</h1>
          {showMessages && (
            <Alert headingLevel="h4" type="info">
              You can change these details later by talking to a move counselor or customer care representative.
            </Alert>
          )}
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
              originDutyLocationName={originDutyLocation?.name || ''}
              originTransportationOfficeName={originDutyLocation?.transportation_office?.name || ''}
              originTransportationOfficePhone={transportationOfficePhone}
              affiliation={ORDERS_BRANCH_OPTIONS[serviceMember?.affiliation] || ''}
              rank={ORDERS_RANK_OPTIONS[rank] || ''}
              edipi={serviceMember?.edipi || ''}
              editURL={customerRoutes.SERVICE_INFO_EDIT_PATH}
              isEditable={moveIsInDraft}
              showMessage={showMessages}
            />
          </SectionWrapper>
          <SectionWrapper className={formStyles.formSection}>
            <OktaInfoDisplay
              oktaUsername={oktaUser?.login || 'Not Provided'}
              oktaEmail={oktaUser?.email || 'Not Provided'}
              oktaFirstName={oktaUser?.firstName || 'Not Provided'}
              oktaLastName={oktaUser?.lastName || 'Not Provided'}
              oktaEdipi={oktaUser?.cac_edipi || 'Not Provided'}
              editURL={customerRoutes.EDIT_OKTA_PROFILE_PATH}
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
  oktaUser: OktaUserInfoShape.isRequired,
};

function mapStateToProps(state) {
  return {
    serviceMember: selectServiceMemberFromLoggedInUser(state),
    moveIsInDraft: selectMoveIsInDraft(state),
    currentOrders: selectCurrentOrders(state) || {},
    currentBackupContacts: selectBackupContacts(state),
    oktaUser: selectOktaUser(state),
  };
}

export default connect(mapStateToProps)(Profile);
