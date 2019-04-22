import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { get } from 'lodash';

import { calculateEntitlementsForShipment } from 'shared/Entities/modules/shipments';
import { formatWeight } from 'shared/formatters';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faComments from '@fortawesome/fontawesome-free-solid/faComments';
import faEmail from '@fortawesome/fontawesome-free-solid/faEnvelope';

function renderEntitlements(entitlements) {
  const weightEntitlement = formatWeight(get(entitlements, 'weight', '0'));
  const proGearEntitlement = formatWeight(get(entitlements, 'pro_gear', '0'));
  const spouseProGearEntitlement = formatWeight(get(entitlements, 'pro_gear_spouse', '0'));

  return (
    <React.Fragment>
      <b>Entitlements</b>
      <br />
      {weightEntitlement} <br />
      Pro-gear: {proGearEntitlement} / Spouse: {spouseProGearEntitlement} <br />
    </React.Fragment>
  );
}

export const CustomerInfo = ({ serviceMember, backupContact, entitlements }) => {
  return (
    <div>
      <div className="usa-grid">
        <div className="extras content">
          <p>
            <b>
              {serviceMember.last_name}, {serviceMember.first_name}
            </b>
            <br />
            DoD ID#: {serviceMember.edipi} - {serviceMember.affiliation} - {serviceMember.rank}
            <br />
            {serviceMember.telephone}
            {serviceMember.secondary_telephone && <span>- {serviceMember.secondary_telephone}</span>}
            <br />
            <a href={`mailto:${serviceMember.personal_email}`}>{serviceMember.personal_email}</a>
            <br />
            Preferred contact method:{' '}
            {serviceMember.phone_is_preferred && (
              <span>
                <FontAwesomeIcon className="icon icon-grey" icon={faPhone} flip="horizontal" /> <span>Phone</span>
              </span>
            )}
            {serviceMember.text_message_is_preferred && (
              <span>
                <FontAwesomeIcon className="icon icon-grey" icon={faComments} /> <span>Text</span>
              </span>
            )}
            {serviceMember.email_is_preferred && (
              <span>
                <FontAwesomeIcon className="icon icon-grey" icon={faEmail} /> <span>Email</span>
              </span>
            )}
          </p>
          <p>
            {backupContact.name && (
              <span>
                <b>Backup Contacts</b>
                <br />
                {backupContact.name} ({backupContact.permission})<br />
                {backupContact.telephone && (
                  <span>
                    {backupContact.telephone}
                    <br />
                  </span>
                )}
                {backupContact.email && (
                  <span>
                    <a href={`mailto:${backupContact.email}`}>{backupContact.email}</a>
                    <br />
                  </span>
                )}
              </span>
            )}
          </p>
          <p>{renderEntitlements(entitlements)}</p>
        </div>
      </div>
    </div>
  );
};

const { bool, object, shape, string, arrayOf } = PropTypes;

CustomerInfo.propTypes = {
  serviceMember: shape({
    backupContacts: arrayOf(
      shape({
        service_member_id: string,
        service_member: object,
        created_at: string,
        permission: string.isRequired,
        id: string,
        updated_at: string,
        name: string.isRequired,
        email: string.isRequired,
        phone: string,
      }),
    ),
    id: string,
    created_at: string,
    updated_at: string,
    user: object,
    user_id: string,
    edipi: string.isRequired,
    rank: string.isRequired,
    affiliation: string.isRequired,
    secondary_telephone: string,
    last_name: string.isRequired,
    telephone: string.isRequired,
    first_name: string.isRequired,
    personal_email: string.isRequired,
    phone_is_preferred: bool,
    text_message_is_preferred: bool,
    email_is_preferred: bool,
  }).isRequired,
};

const mapStateToProps = (state, ownProps) => {
  const defaultServiceMember = {
    backupContacts: [],
    id: '',
    created_at: '',
    updated_at: '',
    user: {},
    user_id: '',
    edipi: '',
    rank: '',
    affiliation: '',
    secondary_telephone: '',
    last_name: '',
    telephone: '',
    first_name: '',
    personal_email: '',
  };
  const defaultBackupContact = {
    service_member_id: '',
    service_member: {},
    created_at: '',
    permission: '',
    id: '',
    updated_at: '',
    name: '',
    email: '',
    phone: '',
  };
  const serviceMember = ownProps.shipment.service_member || defaultServiceMember;
  const backupContact = ownProps.shipment.service_member.backup_contacts[0] || defaultBackupContact;

  return {
    serviceMember,
    backupContact,
    entitlements: calculateEntitlementsForShipment(state, ownProps.shipment.id),
  };
};

export default connect(mapStateToProps)(CustomerInfo);
