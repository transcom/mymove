import React, { Component, Fragment } from 'react';
import { get, isEmpty } from 'lodash';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';

import { getInternalSwaggerDefinition } from 'shared/Swagger/selectors';
import { loadMove, selectMove } from 'shared/Entities/modules/moves';
import {
  fetchLatestOrders,
  selectActiveOrLatestOrders,
  selectUploadsForActiveOrders,
} from 'shared/Entities/modules/orders';
import { SHIPMENT_OPTIONS } from 'shared/constants';

import { moveIsApproved, lastMoveIsCanceled } from 'scenes/Moves/ducks';
import { loadEntitlementsFromState } from 'shared/entitlements';
import Alert from 'shared/Alert';
import { titleCase } from 'shared/constants.js';
import { selectedMoveType as selectMoveType } from 'scenes/Moves/ducks';

import { checkEntitlement } from './ducks';
import ServiceMemberSummary from './ServiceMemberSummary';
import PPMShipmentSummary from './PPMShipmentSummary';
import HHGShipmentSummary from './HHGShipmentSummary';

import './Review.css';
import { selectActivePPMForMove } from '../../shared/Entities/modules/ppms';
// import { showLoggedInUser as showLoggedInUserAction, selectLoggedInUser } from 'shared/Entities/modules/user';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
// import { selectMTOShipmentForMTO } from 'shared/Entities/modules/mtoShipments';

// const shipmentHardcoded = {};
const shipmentHardcoded = {
  agents: [
    {
      agentType: 'RELEASING_AGENT',
      createdAt: '0001-01-01T00:00:00.000Z',
      email: 'ra@example.com',
      firstName: 'ra firstname',
      id: '00000000-0000-0000-0000-000000000000',
      lastName: 'ra lastname',
      mtoShipmentID: '00000000-0000-0000-0000-000000000000',
      phone: '415-444-4444',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
    {
      agentType: 'RECEIVING_AGENT',
      createdAt: '0001-01-01T00:00:00.000Z',
      email: 'andrea@truss.works',
      firstName: 'receivingangetfi',
      id: '00000000-0000-0000-0000-000000000000',
      lastName: 'receiginachlast',
      mtoShipmentID: '00000000-0000-0000-0000-000000000000',
      phone: '415-555-5555',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
  ],
  createdAt: '2020-07-29T00:17:53.236Z',
  customerRemarks: 'lkjlkj',
  destinationAddress: {
    city: 'San Francisco',
    id: '0fda108d-6c6c-44c8-b5ae-485b779f7539',
    postal_code: '94611',
    state: 'CA',
    street_address_1: '666 no',
  },
  id: '3dc3c94f-8264-4dd6-85e0-9a0ec1af3433',
  moveTaskOrderID: 'b21536b7-22a3-43c1-a4a7-3a8c392c1ad5',
  pickupAddress: {
    city: 'San Francisco',
    id: '3ea70395-e15d-485b-8b5a-51549069b9f0',
    postal_code: '94611',
    state: 'CA',
    street_address_1: '666 no',
  },
  requestedDeliveryDate: '2020-07-31',
  requestedPickupDate: '2020-07-30',
  shipmentType: 'HHG',
  updatedAt: '2020-07-29T00:17:53.236Z',
};

export class Summary extends Component {
  componentDidMount() {
    if (this.props.onDidMount) {
      this.props.onDidMount(this.props.serviceMember.id);
      const { showLoggedInUser } = this.props;
      showLoggedInUser();
    }
  }
  componentDidUpdate(prevProps) {
    const { selectedMoveType } = this.props;
    // Only check entitlement for PPMs, not HHGs
    if (prevProps.currentPPM !== this.props.currentPPM && selectedMoveType === SHIPMENT_OPTIONS.PPM) {
      this.props.onCheckEntitlement(this.props.match.params.moveId);
    }
  }

  render() {
    const {
      currentMove,
      currentPPM,
      mtoShipment,
      currentBackupContacts,
      currentOrders,
      schemaRank,
      schemaAffiliation,
      schemaOrdersType,
      moveIsApproved,
      serviceMember,
      entitlement,
      match,
      uploads,
    } = this.props;
    const currentStation = get(serviceMember, 'current_station');
    const stationPhone = get(currentStation, 'transportation_office.phone_lines.0');

    const rootAddressWithMoveId = `/moves/${this.props.match.params.moveId}/review`;
    // isReviewPage being false is the same thing as being in the /edit route
    const isReviewPage = rootAddressWithMoveId === match.url;
    const editSuccessBlurb = this.props.reviewState.editSuccess ? 'Your changes have been saved. ' : '';
    const editOrdersPath = rootAddressWithMoveId + '/edit-orders';

    const showPPMShipmentSummary =
      (isReviewPage && currentPPM) || (!isReviewPage && currentPPM && currentPPM.status !== 'DRAFT');
    const showHHGShipmentSummary = !isEmpty(mtoShipment) || (!isEmpty(mtoShipment) && !isReviewPage);

    const showProfileAndOrders = isReviewPage || !isReviewPage;
    return (
      <Fragment>
        {get(this.props.reviewState.error, 'statusCode', false) === 409 && (
          <Alert type="warning" heading={editSuccessBlurb + 'Your estimated weight is above your entitlement.'}>
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
            <Alert type="info" heading={editSuccessBlurb + 'Note that the entitlement has also changed.'}>
              Your weight entitlement is now {entitlement.sum.toLocaleString()} lbs.
            </Alert>
          )}

        {showProfileAndOrders && (
          <ServiceMemberSummary
            orders={currentOrders}
            uploads={uploads}
            backupContacts={currentBackupContacts}
            serviceMember={serviceMember}
            schemaRank={schemaRank}
            schemaAffiliation={schemaAffiliation}
            schemaOrdersType={schemaOrdersType}
            moveIsApproved={moveIsApproved}
            editOrdersPath={editOrdersPath}
          />
        )}

        {showHHGShipmentSummary && (
          <HHGShipmentSummary mtoShipment={mtoShipment} movePath={rootAddressWithMoveId} entitlements={entitlement} />
        )}

        {showPPMShipmentSummary && (
          <PPMShipmentSummary ppm={currentPPM} movePath={rootAddressWithMoveId} orders={currentOrders} />
        )}
        {moveIsApproved && (
          <div className="approved-edit-warning">
            *To change these fields, contact your local PPPO office at {get(currentStation, 'name')}{' '}
            {stationPhone ? ` at ${stationPhone}` : ''}.
          </div>
        )}
      </Fragment>
    );
  }
}

Summary.propTypes = {
  currentBackupContacts: PropTypes.array,
  getCurrentMove: PropTypes.func,
  currentOrders: PropTypes.object,
  currentPPM: PropTypes.object,
  mtoShipment: PropTypes.object,
  schemaRank: PropTypes.object,
  schemaOrdersType: PropTypes.object,
  moveIsApproved: PropTypes.bool,
  lastMoveIsCanceled: PropTypes.bool,
  error: PropTypes.object,
  selectedMoveType: PropTypes.string.isRequired,
  showLoggedInUser: PropTypes.func.isRequired,
};

function mapStateToProps(state, ownProps) {
  const moveID = state.moves.currentMove.id;
  const currentOrders = selectActiveOrLatestOrders(state);
  // TODO: temporary workaround until moves is consolidated from move_task_orders - this should be the move id
  // const moveTaskOrderID = get(selectLoggedInUser(state), 'service_member.orders[0].move_task_order_id', '');

  return {
    currentPPM: selectActivePPMForMove(state, moveID),
    // mtoShipment: selectMTOShipmentForMTO(state, moveTaskOrderID),
    mtoShipment: shipmentHardcoded,
    serviceMember: state.serviceMember.currentServiceMember,
    currentMove: selectMove(state, ownProps.match.params.moveId),
    currentBackupContacts: state.serviceMember.currentBackupContacts,
    currentOrders: currentOrders,
    uploads: selectUploadsForActiveOrders(state),
    selectedMoveType: selectMoveType(state),
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
    onDidMount: function (smId) {
      const moveID = ownProps.match.params.moveId;
      dispatch(loadMove(moveID, 'Summary.getMove'));
      dispatch(fetchLatestOrders(smId));
    },
    onCheckEntitlement: (moveId) => {
      dispatch(checkEntitlement(moveId));
    },
    showLoggedInUser: showLoggedInUserAction,
  };
}
export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Summary));
