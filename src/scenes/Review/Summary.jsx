import React, { Component, Fragment } from 'react';
import { forEach, get } from 'lodash';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';

import { getInternalSwaggerDefinition } from 'shared/Swagger/selectors';
import { getShipment, selectShipment } from 'shared/Entities/modules/shipments';
import { getMove, selectMove } from 'shared/Entities/modules/moves';
import { getCurrentShipmentID } from 'shared/UI/ducks';

import { moveIsApproved, lastMoveIsCanceled } from 'scenes/Moves/ducks';
import { loadEntitlementsFromState } from 'shared/entitlements';
import Alert from 'shared/Alert';
import { titleCase } from 'shared/constants.js';

import { checkEntitlement } from './ducks';
import ServiceMemberSummary from './ServiceMemberSummary';
import PPMShipmentSummary from './PPMShipmentSummary';
import HHGShipmentSummary from './HHGShipmentSummary';

import './Review.css';

export class Summary extends Component {
  componentDidMount() {
    if (this.props.onDidMount) {
      this.props.onDidMount();
    }
  }
  componentDidUpdate(prevProps) {
    // Only check entitlement for PPMs, not HHGs
    if (prevProps.currentPpm !== this.props.currentPpm) {
      this.props.onCheckEntitlement(this.props.match.params.moveId);
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
      serviceMember,
      entitlement,
    } = this.props;

    const currentStation = get(serviceMember, 'current_station');
    const stationPhone = get(currentStation, 'transportation_office.phone_lines.0');

    const rootAddressWithMoveId = `/moves/${this.props.match.params.moveId}/review`;
    const editSuccessBlurb = this.props.reviewState.editSuccess ? 'Your changes have been saved. ' : '';
    const editOrdersPath = rootAddressWithMoveId + '/edit-orders';

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

        <ServiceMemberSummary
          orders={currentOrders}
          backupContacts={currentBackupContacts}
          serviceMember={serviceMember}
          schemaRank={schemaRank}
          schemaAffiliation={schemaAffiliation}
          schemaOrdersType={schemaOrdersType}
          moveIsApproved={moveIsApproved}
          editOrdersPath={editOrdersPath}
        />

        {currentPpm && <PPMShipmentSummary ppm={currentPpm} movePath={rootAddressWithMoveId} />}
        {currentShipment && (
          <HHGShipmentSummary shipment={currentShipment} movePath={rootAddressWithMoveId} entitlements={entitlement} />
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
  currentPpm: PropTypes.object,
  currentShipment: PropTypes.object,
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
    currentMove: selectMove(state, ownProps.match.params.moveId),
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
    },
    onCheckEntitlement: moveId => {
      dispatch(checkEntitlement(moveId));
    },
  };
}
export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Summary));
