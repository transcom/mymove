import React, { Component, Fragment } from 'react';
import { Link } from 'react-router-dom';
import { get } from 'lodash';
import moment from 'moment';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';
import ppmBlack from 'shared/icon/ppm-black.svg';
import { moveIsApproved } from 'scenes/Moves/ducks';
import { formatCentsRange } from 'shared/formatters';
import { loadEntitlementsFromState } from 'shared/entitlements';
import { checkEntitlement } from './ducks';
import Alert from 'shared/Alert';
import { titleCase } from 'shared/constants.js';
import './Review.css';

export class Summary extends Component {
  componentDidMount() {
    this.props.checkEntitlement(this.props.match.params.moveId);
  }
  render() {
    const {
      currentPpm,
      currentBackupContacts,
      currentOrders,
      schemaRank,
      schemaAffiliation,
      schemaOrdersType,
      moveIsApproved,
      serviceMember,
      entitlement,
    } = this.props;
    const yesNoMap = { true: 'Yes', false: 'No' };
    function getFullName() {
      if (!serviceMember) return;
      return `${serviceMember.first_name} ${serviceMember.middle_name || ''} ${
        serviceMember.last_name
      } ${serviceMember.suffix || ''}`;
    }
    function getFullAddress(address) {
      if (address) {
        return (
          <Fragment>
            <div>{address.street_address_1}</div>
            {address.street_address_2 && <div>{address.street_address_2}</div>}
            <div>
              {address.city}, {address.state} {address.postal_code}
            </div>
          </Fragment>
        );
      }
    }
    function getFullContactPreferences() {
      if (!serviceMember) return;
      const prefs = {
        phone_is_preferred: 'Phone',
        text_message_is_preferred: 'Text',
        email_is_preferred: 'Email',
      };
      const preferredMethods = [];
      Object.keys(prefs).forEach(propertyName => {
        if (serviceMember[propertyName]) {
          preferredMethods.push(prefs[propertyName]);
        }
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
    function formatDate(date) {
      if (!date) return;
      return moment(date, 'YYYY-MM-DD').format('MM/DD/YYYY');
    }

    const rootAddress = `/moves/${this.props.match.params.moveId}/review`;
    const editProfileAddress = rootAddress + '/edit-profile';
    const editBackupContactAddress = rootAddress + '/edit-backup-contact';
    const editContactInfoAddress = rootAddress + '/edit-contact-info';
    const editOrdersAddress = rootAddress + '/edit-orders';
    const editDateAndLocationAddress = rootAddress + '/edit-date-and-location';
    const editWeightAddress = rootAddress + '/edit-weight';
    const privateStorageString = get(
      currentPpm,
      'estimated_storage_reimbursement',
    )
      ? `(spend up to ${currentPpm.estimated_storage_reimbursement.toLocaleString()} on private storage)`
      : '';
    const sitDisplay = get(currentPpm, 'has_sit', false)
      ? `${currentPpm.days_in_storage} days ${privateStorageString}`
      : 'Not requested';
    const editSuccessBlurb = this.props.reviewState.editSuccess
      ? 'Your changes have been saved. '
      : '';
    return (
      <Fragment>
        {get(this.props.reviewState.error, 'statusCode', false) === 409 && (
          <Alert
            type="warning"
            heading={
              editSuccessBlurb +
              'Your estimated weight is above your entitlement.'
            }
          >
            {titleCase(this.props.reviewState.error.response.body.message)}.
          </Alert>
        )}
        {this.props.reviewState.editSuccess &&
          !this.props.reviewState.entitlementChange &&
          get(this.props.reviewState.error, 'statusCode', false) === false && (
            <Alert type="success" heading={editSuccessBlurb} />
          )}
        {this.props.reviewState.entitlementChange &&
          get(this.props.reviewState.error, 'statusCode', false) === false && (
            <Alert
              type="info"
              heading={
                editSuccessBlurb + 'Note that the entitlement has also changed.'
              }
            >
              Your weight entitlement is now {entitlement.sum.toLocaleString()}{' '}
              lbs.
            </Alert>
          )}

        <h3>Profile and Orders</h3>
        <div className="usa-grid-full review-content">
          <div className="usa-width-one-half review-section">
            <table>
              <tbody>
                <tr>
                  <th>
                    Profile{' '}
                    <span className="align-right">
                      <Link to={editProfileAddress}>Edit</Link>
                    </span>
                  </th>
                </tr>
                <tr>
                  <td> Name: </td>
                  <td>{getFullName()}</td>
                </tr>
                <tr>
                  <td>Branch:</td>
                  <td>
                    {get(
                      schemaAffiliation['x-display-value'],
                      get(serviceMember, 'affiliation'),
                    )}
                  </td>
                </tr>
                <tr>
                  <td> Rank/Pay Grade: </td>
                  <td>
                    {get(
                      schemaRank['x-display-value'],
                      get(serviceMember, 'rank'),
                    )}
                  </td>
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

            <table>
              <tbody>
                <tr>
                  <th>
                    Orders{moveIsApproved && '*'}
                    {!moveIsApproved && (
                      <span className="align-right">
                        <Link to={editOrdersAddress}>Edit</Link>
                      </span>
                    )}
                  </th>
                </tr>
                <tr>
                  <td> Orders Type: </td>
                  <td>
                    {get(
                      schemaOrdersType['x-display-value'],
                      get(currentOrders, 'orders_type'),
                    )}
                  </td>
                </tr>
                <tr>
                  <td> Orders Date: </td>
                  <td> {formatDate(get(currentOrders, 'issue_date'))}</td>
                </tr>
                <tr>
                  <td> Report-by Date: </td>
                  <td>{formatDate(get(currentOrders, 'report_by_date'))}</td>
                </tr>
                <tr>
                  <td> New Duty Station: </td>
                  <td> {get(currentOrders, 'new_duty_station.name')}</td>
                </tr>
                <tr>
                  <td> Dependents?: </td>
                  <td>
                    {' '}
                    {currentOrders &&
                      yesNoMap[get(currentOrders, 'has_dependents').toString()]}
                  </td>
                </tr>
                {currentOrders &&
                  get(currentOrders, 'spouse_has_pro_gear') && (
                    <tr>
                      <td> Spouse Pro Gear?: </td>
                      <td>
                        {currentOrders &&
                          yesNoMap[
                            get(currentOrders, 'spouse_has_pro_gear').toString()
                          ]}
                      </td>
                    </tr>
                  )}
                <tr>
                  <td> Orders Uploaded: </td>
                  <td>
                    {get(currentOrders, 'uploaded_orders.uploads') &&
                      get(currentOrders, 'uploaded_orders.uploads').length}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div className="usa-width-one-half review-section">
            <table>
              <tbody>
                <tr>
                  <th>
                    Contact Info{' '}
                    <span className="align-right">
                      <Link to={editContactInfoAddress}>Edit</Link>
                    </span>
                  </th>
                </tr>
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
                  <td>{getFullContactPreferences()}</td>
                </tr>
                <tr>
                  <td> Current Mailing Address: </td>
                  <td>
                    {getFullAddress(get(serviceMember, 'residential_address'))}
                  </td>
                </tr>
                <tr>
                  <td> Backup Mailing Address: </td>
                  <td>
                    {getFullAddress(
                      get(serviceMember, 'backup_mailing_address'),
                    )}
                  </td>
                </tr>
              </tbody>
            </table>
            {currentBackupContacts.map(contact => (
              <table key={contact.id}>
                <tbody>
                  <tr>
                    <th>
                      Backup Contact{' '}
                      <span className="align-right">
                        <Link to={editBackupContactAddress}>Edit</Link>
                      </span>
                    </th>
                  </tr>
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
            ))}
          </div>
        </div>
        {currentPpm && (
          <div className="usa-grid-full ppm-container">
            <h3>
              <img src={ppmBlack} alt="PPM shipment" /> Shipment - You move your
              stuff (PPM)
            </h3>
            <div className="usa-width-one-half review-section ppm-review-section">
              <table>
                <tbody>
                  <tr>
                    <th>
                      Dates & Locations
                      <span className="align-right">
                        <Link to={editDateAndLocationAddress}>Edit</Link>
                      </span>
                    </th>
                  </tr>
                  <tr>
                    <td> Move Date: </td>
                    <td>{formatDate(get(currentPpm, 'planned_move_date'))}</td>
                  </tr>
                  <tr>
                    <td> Pickup ZIP Code: </td>
                    <td> {currentPpm && currentPpm.pickup_postal_code}</td>
                  </tr>
                  {currentPpm.has_additional_postal_code && (
                    <tr>
                      <td> Additional Pickup: </td>
                      <td> {currentPpm.additional_pickup_postal_code}</td>
                    </tr>
                  )}
                  <tr>
                    <td> Delivery ZIP Code: </td>
                    <td> {currentPpm && currentPpm.destination_postal_code}</td>
                  </tr>
                  <tr>
                    <td> Storage: </td>
                    <td>{sitDisplay}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div className="usa-width-one-half review-section ppm-review-section">
              <table>
                <tbody>
                  <tr>
                    <th>
                      Weight
                      <span className="align-right">
                        <Link to={editWeightAddress}>Edit</Link>
                      </span>
                    </th>
                  </tr>
                  <tr>
                    <td> Estimated Weight: </td>
                    <td>
                      {' '}
                      {currentPpm &&
                        currentPpm.weight_estimate.toLocaleString()}{' '}
                      lbs
                    </td>
                  </tr>
                  <tr>
                    <td> Estimated PPM Incentive: </td>
                    <td>
                      {' '}
                      {currentPpm &&
                        formatCentsRange(
                          currentPpm.incentive_estimate_min,
                          currentPpm.incentive_estimate_max,
                        )}
                    </td>
                  </tr>
                  {currentPpm.has_requested_advance && (
                    <tr>
                      <td> Advance: </td>
                      <td>
                        {' '}
                        ${(
                          currentPpm.advance.requested_amount / 100
                        ).toLocaleString()}
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>
        )}
        {moveIsApproved && (
          <div className="approved-edit-warning">
            *To change these fields, contact your local PPPO office.
          </div>
        )}
      </Fragment>
    );
  }
}

Summary.propTypes = {
  currentPpm: PropTypes.object,
  currentBackupContacts: PropTypes.array,
  currentOrders: PropTypes.object,
  schemaRank: PropTypes.object,
  schemaOrdersType: PropTypes.object,
  moveIsApproved: PropTypes.bool,
  checkEntitlement: PropTypes.func.isRequired,
  error: PropTypes.object,
};

function mapStateToProps(state) {
  return {
    currentPpm: state.ppm.currentPpm,
    serviceMember: state.serviceMember.currentServiceMember,
    currentMove: state.moves.currentMove,
    currentBackupContacts: state.serviceMember.currentBackupContacts,
    currentOrders: state.orders.currentOrders,
    schemaRank: get(state, 'swagger.spec.definitions.ServiceMemberRank', {}),
    schemaOrdersType: get(state, 'swagger.spec.definitions.OrdersType', {}),
    schemaAffiliation: get(state, 'swagger.spec.definitions.Affiliation', {}),
    moveIsApproved: moveIsApproved(state),
    reviewState: state.review,
    entitlement: loadEntitlementsFromState(state),
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      checkEntitlement,
      loadEntitlementsFromState,
    },
    dispatch,
  );
}
export default withRouter(
  connect(mapStateToProps, mapDispatchToProps)(Summary),
);
