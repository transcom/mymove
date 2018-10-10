import React, { Component, Fragment } from 'react';
import { Link } from 'react-router-dom';
import { forEach, get } from 'lodash';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';

import { getInternalSwaggerDefinition } from 'shared/Swagger/selectors';
import { getShipment, selectShipment } from 'shared/Entities/modules/shipments';
import { getMove } from 'shared/Entities/modules/moves';
import { getCurrentShipmentID } from 'shared/UI/ducks';

import { moveIsApproved, lastMoveIsCanceled } from 'scenes/Moves/ducks';
import { loadEntitlementsFromState } from 'shared/entitlements';
import Alert from 'shared/Alert';
import { titleCase } from 'shared/constants.js';
import { formatDateSM } from 'shared/formatters';

import { checkEntitlement } from './ducks';
import PPMShipmentSummary from './PPMShipmentSummary';
import HHGShipmentSummary from './HHGShipmentSummary';
import Address from './Address';

import './Review.css';

export class Summary extends Component {
  componentDidMount() {
    if (this.props.onDidMount) {
      this.props.onDidMount();
    }
  }
  render() {
    const {
      currentMove,
      currentPpm,
      currentShipment,
      currentBackupContacts,
      currentOrders,
      schemaRank,
      schemaAffiliation,
      schemaOrdersType,
      moveIsApproved,
      lastMoveIsCanceled,
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
    function getFullContactPreferences() {
      if (!serviceMember) return;
      const prefs = {
        phone_is_preferred: 'Phone',
        text_message_is_preferred: 'Text',
        email_is_preferred: 'Email',
      };
      const preferredMethods = [];
      Object.keys(prefs).forEach(propertyName => {
        /* eslint-disable */
        if (serviceMember[propertyName]) {
          preferredMethods.push(prefs[propertyName]);
        }
        /* eslint-enable */
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
    const currentStation = get(serviceMember, 'current_station');
    const stationPhone = get(
      currentStation,
      'transportation_office.phone_lines.0',
    );

    const rootAddress = `/moves/review`;
    const rootAddressWithMoveId = `/moves/${
      this.props.match.params.moveId
    }/review`;
    const editProfileAddress = rootAddress + '/edit-profile';
    const editBackupContactAddress = rootAddress + '/edit-backup-contact';
    const editContactInfoAddress = rootAddress + '/edit-contact-info';
    const editOrdersAddress = rootAddressWithMoveId + '/edit-orders';
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
        {currentMove &&
          this.props.reviewState.entitlementChange &&
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
                    Profile
                    <span className="edit-section-link">
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
            {!lastMoveIsCanceled && (
              <table>
                <tbody>
                  <tr>
                    <th>
                      Orders
                      {moveIsApproved && '*'}
                      {!moveIsApproved && (
                        <span className="edit-section-link">
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
                    <td> {formatDateSM(get(currentOrders, 'issue_date'))}</td>
                  </tr>
                  <tr>
                    <td> Report-by Date: </td>
                    <td>
                      {formatDateSM(get(currentOrders, 'report_by_date'))}
                    </td>
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
                        yesNoMap[
                          get(currentOrders, 'has_dependents').toString()
                        ]}
                    </td>
                  </tr>
                  {currentOrders &&
                    get(currentOrders, 'spouse_has_pro_gear') && (
                      <tr>
                        <td> Spouse Pro Gear?: </td>
                        <td>
                          {currentOrders &&
                            yesNoMap[
                              get(
                                currentOrders,
                                'spouse_has_pro_gear',
                              ).toString()
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
            )}
          </div>

          <div className="usa-width-one-half review-section">
            <table>
              <tbody>
                <tr>
                  <th>
                    Contact Info
                    <span className="edit-section-link">
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
                    <Address
                      address={get(serviceMember, 'residential_address')}
                    />
                  </td>
                </tr>
                <tr>
                  <td> Backup Mailing Address: </td>
                  <td>
                    <Address
                      address={get(serviceMember, 'backup_mailing_address')}
                    />
                  </td>
                </tr>
              </tbody>
            </table>
            {currentBackupContacts.map(contact => (
              <table key={contact.id}>
                <tbody>
                  <tr>
                    <th>
                      Backup Contact Info
                      <span className="edit-section-link">
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
        {currentPpm &&
          !lastMoveIsCanceled && (
            <PPMShipmentSummary
              ppm={currentPpm}
              movePath={rootAddressWithMoveId}
            />
          )}

        {currentShipment &&
          !lastMoveIsCanceled && (
            <HHGShipmentSummary
              shipment={currentShipment}
              movePath={rootAddressWithMoveId}
              entitlements={entitlement}
            />
          )}

        {moveIsApproved && (
          <div className="approved-edit-warning">
            *To change these fields, contact your local PPPO office at{' '}
            {get(currentStation, 'name')}{' '}
            {stationPhone ? ` at ${stationPhone}` : ''}.
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
  lastMoveIsCanceled: PropTypes.bool,
  error: PropTypes.object,
};

function mapStateToProps(state, ownProps) {
  return {
    currentPpm: state.ppm.currentPpm,
    currentShipment: selectShipment(state, getCurrentShipmentID(state)),
    serviceMember: state.serviceMember.currentServiceMember,
    currentMove: getMove(state, ownProps.match.params.moveId),
    // latestMove: state.moves.latestMove,
    currentBackupContacts: state.serviceMember.currentBackupContacts,
    currentOrders: state.orders.currentOrders,
    schemaRank: getInternalSwaggerDefinition(state, 'ServiceMemberRank'),
    schemaOrdersType: getInternalSwaggerDefinition(state, 'OrdersType'),
    schemaAffiliation: getInternalSwaggerDefinition(state, 'Affiliation'),
    moveIsApproved: moveIsApproved(state),
    lastMoveIsCanceled: lastMoveIsCanceled(state),
    reviewState: state.review,
    entitlement: loadEntitlementsFromState(state),
  };
}
function mapDispatchToProps(dispatch, ownProps) {
  return {
    onDidMount: function() {
      const moveID = ownProps.match.params.moveId;
      dispatch(getMove('Summary.getMove', moveID)).then(function(action) {
        forEach(action.entities.shipments, function(shipment) {
          dispatch(getShipment('Summary.getShipment', shipment.id));
        });
      });
      dispatch(checkEntitlement(moveID));
    },
  };
}
export default withRouter(
  connect(mapStateToProps, mapDispatchToProps)(Summary),
);
