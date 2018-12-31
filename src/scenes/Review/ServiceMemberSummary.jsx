import React, { Fragment } from 'react';
import { Link, withRouter } from 'react-router-dom';
import { get } from 'lodash';
import PropTypes from 'prop-types';

import { formatDateSM } from 'shared/formatters';
import Address from './Address';

import './Review.css';

function getFullName(serviceMember) {
  if (!serviceMember) return;
  return `${serviceMember.first_name} ${serviceMember.middle_name || ''} ${
    serviceMember.last_name
  } ${serviceMember.suffix || ''}`;
}

function getFullContactPreferences(serviceMember) {
  if (!serviceMember) return;
  const prefs = {
    phone_is_preferred: 'Phone',
    text_message_is_preferred: 'Text',
    email_is_preferred: 'Email',
  };
  const preferredMethods = [];
  Object.keys(prefs).forEach(propertyName => {
    /* eslint-disable security/detect-object-injection */
    if (serviceMember[propertyName]) {
      preferredMethods.push(prefs[propertyName]);
    }
    /* eslint-enable security/detect-object-injection */
  });
  return preferredMethods.join(', ');
}

// TODO: Uncomment function below after backup contact auth is implemented.
// function getFullBackupPermission(backup_contact) {
//   const perms = {
//     NONE: '',
//     VIEW: 'View all aspects of this move',
//     EDIT:
//       'Authorized to represent me in all aspects of this move (letter of authorization)',
//   };
//   return `${perms[backup_contact.permission]}`;
// }

function ServiceMemberSummary(props) {
  const {
    backupContacts,
    orders,
    serviceMember,
    schemaRank,
    schemaAffiliation,
    schemaOrdersType,
    moveIsApproved,
    editOrdersPath,
  } = props;

  const rootPath = `/moves/review`;
  const editProfilePath = rootPath + '/edit-profile';
  const editBackupContactPath = rootPath + '/edit-backup-contact';
  const editContactInfoPath = rootPath + '/edit-contact-info';

  const yesNoMap = { true: 'Yes', false: 'No' };

  return (
    <div>
      <h3>Profile and Orders</h3>
      <div className="usa-grid-full review-content">
        <div className="usa-width-one-half review-section">
          <p className="heading">
            Profile
            <span className="edit-section-link">
              <Link to={editProfilePath}>Edit</Link>
            </span>
          </p>
          <table>
            <tbody>
              <tr>
                <td> Name: </td>
                <td>{getFullName(serviceMember)}</td>
              </tr>
              <tr>
                <td>Branch:</td>
                <td>{get(schemaAffiliation['x-display-value'], get(serviceMember, 'affiliation'))}</td>
              </tr>
              <tr>
                <td> Rank/Pay Grade: </td>
                <td>{get(schemaRank['x-display-value'], get(serviceMember, 'rank'))}</td>
              </tr>
              <tr>
                <td> DoD ID#: </td>
                <td>{get(serviceMember, 'edipi')}</td>
              </tr>
              <tr>
                <td> Current Duty Station: </td>
                <td>{get(serviceMember, 'current_station.name')}</td>
              </tr>
            </tbody>
          </table>
          {orders && (
            <Fragment>
              <p className="heading">
                Orders
                {moveIsApproved && '*'}
                {!moveIsApproved && (
                  <span className="edit-section-link">
                    <Link to={editOrdersPath}>Edit</Link>
                  </span>
                )}
              </p>
              <table>
                <tbody>
                  <tr>
                    <td> Orders Type: </td>
                    <td>{get(schemaOrdersType['x-display-value'], get(orders, 'orders_type'))}</td>
                  </tr>
                  <tr>
                    <td> Orders Date: </td>
                    <td> {formatDateSM(get(orders, 'issue_date'))}</td>
                  </tr>
                  <tr>
                    <td> Report-by Date: </td>
                    <td>{formatDateSM(get(orders, 'report_by_date'))}</td>
                  </tr>
                  <tr>
                    <td> New Duty Station: </td>
                    <td> {get(orders, 'new_duty_station.name')}</td>
                  </tr>
                  <tr>
                    <td> Dependents?: </td>
                    <td> {orders && yesNoMap[get(orders, 'has_dependents').toString()]}</td>
                  </tr>
                  {orders &&
                    get(orders, 'spouse_has_pro_gear') && (
                      <tr>
                        <td> Spouse Pro Gear?: </td>
                        <td>{orders && yesNoMap[get(orders, 'spouse_has_pro_gear').toString()]}</td>
                      </tr>
                    )}
                  <tr>
                    <td> Orders Uploaded: </td>
                    <td>{get(orders, 'uploaded_orders.uploads') && get(orders, 'uploaded_orders.uploads').length}</td>
                  </tr>
                </tbody>
              </table>
            </Fragment>
          )}
        </div>

        <div className="usa-width-one-half review-section">
          <p className="heading">
            Contact Info
            <span className="edit-section-link">
              <Link to={editContactInfoPath}>Edit</Link>
            </span>
          </p>
          <table>
            <tbody>
              <tr>
                <td> Best Contact Phone: </td>
                <td>{get(serviceMember, 'telephone')}</td>
              </tr>
              <tr>
                <td> Alt. Phone: </td>
                <td>{get(serviceMember, 'secondary_telephone')}</td>
              </tr>
              <tr>
                <td> Personal Email: </td>
                <td>{get(serviceMember, 'personal_email')}</td>
              </tr>
              <tr>
                <td> Preferred Contact Method: </td>
                <td>{getFullContactPreferences(serviceMember)}</td>
              </tr>
              <tr>
                <td> Current Mailing Address: </td>
                <td>
                  <Address address={get(serviceMember, 'residential_address')} />
                </td>
              </tr>
              <tr>
                <td> Backup Mailing Address: </td>
                <td>
                  <Address address={get(serviceMember, 'backup_mailing_address')} />
                </td>
              </tr>
            </tbody>
          </table>
          {backupContacts.map(contact => (
            <Fragment key={contact.id}>
              <p className="heading">
                Backup Contact Info
                <span className="edit-section-link">
                  <Link to={editBackupContactPath}>Edit</Link>
                </span>
              </p>
              <table key={contact.id}>
                <tbody>
                  <tr>
                    <td> Backup Contact: </td>
                    <td>
                      {contact.name} <br />
                      {/* getFullBackupPermission(contact) */}
                    </td>
                  </tr>
                  <tr>
                    <td> Email: </td>
                    <td> {contact.email} </td>
                  </tr>
                  <tr>
                    <td> Phone: </td>
                    <td> {contact.telephone}</td>
                  </tr>
                </tbody>
              </table>
            </Fragment>
          ))}
        </div>
      </div>
    </div>
  );
}

ServiceMemberSummary.propTypes = {
  backupContacts: PropTypes.array.isRequired,
  serviceMember: PropTypes.object.isRequired,
  schemaRank: PropTypes.object.isRequired,
  schemaAffiliation: PropTypes.object.isRequired,
  schemaOrdersType: PropTypes.object.isRequired,
  orders: PropTypes.object,
  moveIsApproved: PropTypes.bool,
  editOrdersPath: PropTypes.string,
};

export default withRouter(ServiceMemberSummary);
