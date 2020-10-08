/* eslint-ignore */
import React, { Component, Fragment } from 'react';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { getInternalSwaggerDefinition } from 'shared/Swagger/selectors';
import { loadMove, selectMove } from 'shared/Entities/modules/moves';
import { selectActiveOrLatestOrdersFromEntities, selectUploadsForActiveOrders } from 'shared/Entities/modules/orders';
import { SHIPMENT_OPTIONS } from 'shared/constants';

import { moveIsApproved, lastMoveIsCanceled } from 'scenes/Moves/ducks';
import { loadEntitlementsFromState } from 'shared/entitlements';
import Alert from 'shared/Alert';
import { titleCase } from 'shared/constants.js';
import { selectedMoveType as selectMoveType } from 'scenes/Moves/ducks';

import { checkEntitlement } from './ducks';
import ServiceMemberSummary from './ServiceMemberSummary';
import PPMShipmentSummary from './PPMShipmentSummary';
import HHGShipmentSummary from 'pages/MyMove/HHGShipmentSummary';

import './Review.css';
import { selectActivePPMForMove } from '../../shared/Entities/modules/ppms';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { selectMTOShipmentsByMoveId } from 'shared/Entities/modules/mtoShipments';

export class Summary extends Component {
  componentDidMount() {
    if (this.props.onDidMount) {
      this.props.onDidMount(this.props.serviceMember.id);
    }
  }
  componentDidUpdate(prevProps) {
    const { selectedMoveType } = this.props;
    const hhgMove = !Object.keys(prevProps.currentPPM).length && !Object.keys(this.props.currentPPM).length;
    // Only check entitlement for PPMs, not HHGs
    if (prevProps.currentPPM !== this.props.currentPPM && !hhgMove && selectedMoveType === SHIPMENT_OPTIONS.PPM) {
      this.props.onCheckEntitlement(this.props.match.params.moveId);
    }
  }

  render() {
    const {
      currentMove,
      currentPPM,
      mtoShipments,
      currentBackupContacts,
      currentOrders,
      schemaRank,
      schemaAffiliation,
      schemaOrdersType,
      moveIsApproved,
      serviceMember,
      entitlement,
      match,
      history,
      uploads,
    } = this.props;

    const currentStation = get(serviceMember, 'current_station');
    const stationPhone = get(currentStation, 'transportation_office.phone_lines.0');

    const rootAddressWithMoveId = `/moves/${this.props.match.params.moveId}`;
    const rootReviewAddressWithMoveId = rootAddressWithMoveId + `/review`;

    // isReviewPage being false is the same thing as being in the /edit route
    const isReviewPage = rootReviewAddressWithMoveId === match.url;
    const editSuccessBlurb = this.props.reviewState.editSuccess ? 'Your changes have been saved. ' : '';
    const editOrdersPath = rootReviewAddressWithMoveId + '/edit-orders';

    const showPPMShipmentSummary =
      (isReviewPage && Object.keys(currentPPM).length) ||
      (!isReviewPage && Object.keys(currentPPM).length && currentPPM.status !== 'DRAFT');
    const showHHGShipmentSummary = isReviewPage && !!mtoShipments.length;
    const hasPPMorHHG = (isReviewPage && Object.keys(currentPPM).length) || !!mtoShipments.length;

    const showProfileAndOrders = isReviewPage || !isReviewPage;
    const showMoveSetup = showPPMShipmentSummary || showHHGShipmentSummary;
    const shipmentSelectionPath = `/moves/${currentMove.id}/select-type`;

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
        {showMoveSetup && <h3>Move setup</h3>}
        {showPPMShipmentSummary && (
          <PPMShipmentSummary ppm={currentPPM} movePath={rootReviewAddressWithMoveId} orders={currentOrders} />
        )}
        {showHHGShipmentSummary &&
          mtoShipments.map((shipment, index) => {
            return (
              <HHGShipmentSummary
                key={shipment.id}
                mtoShipment={shipment}
                shipmentNumber={index + 1}
                movePath={rootAddressWithMoveId}
                newDutyStationPostalCode={currentOrders.new_duty_station.address.postal_code}
              />
            );
          })}
        {hasPPMorHHG && (
          <div className="grid-col-row margin-top-5">
            <span className="float-right">Optional</span>
            <h3>Add another shipment</h3>
            <p>Will you move any belongings to or from another location?</p>
            <Button data-testid="addAnotherShipmentBtn" secondary onClick={() => history.push(shipmentSelectionPath)}>
              Add another shipment
            </Button>
          </div>
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
  schemaRank: PropTypes.object,
  schemaOrdersType: PropTypes.object,
  moveIsApproved: PropTypes.bool,
  lastMoveIsCanceled: PropTypes.bool,
  error: PropTypes.object,
  selectedMoveType: PropTypes.string.isRequired,
  showLoggedInUser: PropTypes.func.isRequired,
};

function mapStateToProps(state, ownProps) {
  const moveID = ownProps.match.params.moveId;
  const currentOrders = selectActiveOrLatestOrdersFromEntities(state);

  return {
    currentPPM: selectActivePPMForMove(state, moveID),
    mtoShipments: selectMTOShipmentsByMoveId(state, moveID),
    serviceMember: state.serviceMember.currentServiceMember,
    currentMove: selectMove(state, moveID),
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
    onDidMount: function () {
      const moveID = ownProps.match.params.moveId;
      dispatch(loadMove(moveID, 'Summary.getMove'));
      dispatch(showLoggedInUserAction());
    },
    onCheckEntitlement: (moveId) => {
      dispatch(checkEntitlement(moveId));
    },
    showLoggedInUser: showLoggedInUserAction,
  };
}
export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Summary));
