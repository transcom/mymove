import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { get } from 'lodash';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faComments from '@fortawesome/fontawesome-free-solid/faComments';
import faEmail from '@fortawesome/fontawesome-free-solid/faEnvelope';

class CustomerInfo extends Component {
  render() {
    const { serviceMember, backupContact } = this.props;

    return (
      <div>
        <div className="usa-grid">
          <div className="extras content">
            <p>
              <b>
                {serviceMember.last_name}, {serviceMember.first_name}
              </b>
              <br />
              DoD ID#: {serviceMember.edipi} - {serviceMember.affiliation} -{' '}
              {serviceMember.rank}
              <br />
              {serviceMember.telephone}
              {serviceMember.secondary_telephone && (
                <span>- {serviceMember.secondary_telephone}</span>
              )}
              <br />
              {serviceMember.personal_email}
              <br />
              Preferred contact method:{' '}
              {serviceMember.phone_is_preferred && (
                <FontAwesomeIcon
                  className="icon"
                  icon={faPhone}
                  flip="horizontal"
                />
              )}
              {serviceMember.text_message_is_preferred && (
                <FontAwesomeIcon className="icon" icon={faComments} />
              )}
              {serviceMember.email_is_preferred && (
                <FontAwesomeIcon className="icon" icon={faEmail} />
              )}
            </p>
            <p>
              {backupContact.name && (
                <span>
                  <b>Backup Contacts</b>
                  <br />
                  {backupContact.name} ({backupContact.permission})
                  <br />
                  {backupContact.telephone && (
                    <span>
                      {backupContact.telephone}
                      <br />
                    </span>
                  )}
                  {backupContact.email && (
                    <span>
                      {backupContact.email}
                      <br />
                    </span>
                  )}
                </span>
              )}
            </p>
          </div>
        </div>
      </div>
    );
  }
}

CustomerInfo.propTypes = {
  serviceMember: PropTypes.object,
};

const { object, shape, string, arrayOf } = PropTypes;

CustomerInfo.propTypes = {
  serviceMember: shape({
    backupContacts: arrayOf(
      shape({
        service_member_id: string.isRequired,
        service_member: object.isRequired,
        created_at: string.isRequired,
        permission: string.isRequired,
        id: string.isRequired,
        updated_at: string,
        name: string.isRequired,
        email: string.isRequired,
        phone: string,
      }),
    ).isRequired,
    id: string.isRequired,
    created_at: string.isRequired,
    updated_at: string.isRequired,
    user: object.isRequired,
    user_id: string.isRequired,
    edipi: string.isRequired,
    rank: string.isRequired,
    affiliation: string.isRequired,
    secondary_telephone: string,
    last_name: string.isRequired,
    telephone: string.isRequired,
    first_name: string.isRequired,
    personal_email: string.isRequired,
  }).isRequired,
};

const mapStateToProps = state => {
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
  const serviceMember = get(
    state,
    'tsp.shipment.service_member',
    defaultServiceMember,
  );
  const backupContact = get(
    state,
    'tsp.shipment.service_member.backup_contacts[0]',
    defaultBackupContact,
  );

  return {
    serviceMember,
    backupContact,
  };
};

export default connect(mapStateToProps)(CustomerInfo);
