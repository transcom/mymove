/* eslint-ignore */
import React, { Component, Fragment } from 'react';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';
import moment from 'moment';
import { Button } from '@trussworks/react-uswds';

import { getInternalSwaggerDefinition } from 'shared/Swagger/selectors';
import { loadMove, selectMove } from 'shared/Entities/modules/moves';
import { selectActiveOrLatestOrdersFromEntities, selectUploadsForActiveOrders } from 'shared/Entities/modules/orders';
import { SHIPMENT_OPTIONS } from 'shared/constants';

import { moveIsApproved, lastMoveIsCanceled } from 'scenes/Moves/ducks';
import { loadEntitlementsFromState } from 'shared/entitlements';
import { formatOrderType } from 'shared/utils';
import Alert from 'shared/Alert';
import { titleCase } from 'shared/constants.js';
import { selectedMoveType as selectMoveType } from 'scenes/Moves/ducks';

import { checkEntitlement } from './ducks';
import PPMShipmentSummary from './PPMShipmentSummary';
import ProfileTable from 'components/Customer/Review/ProfileTable';
import OrdersTable from 'components/Customer/Review/OrdersTable';
import PPMShipmentCard from 'components/Customer/Review/ShipmentCard/PPMShipmentCard';
import HHGShipmentCard from 'components/Customer/Review/ShipmentCard/HHGShipmentCard';

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

  handleEditClick = (path) => {
    const { history } = this.props;
    history.push(path);
  };

  get getSortedShipments() {
    const { currentPPM, mtoShipments } = this.props;
    const sortedShipments = [...mtoShipments];
    if (Object.keys(currentPPM).length) {
      const ppm = { ...currentPPM };
      ppm.shipmentType = SHIPMENT_OPTIONS.PPM;
      // workaround for differing cases between mtoShipments and ppms (bigger change needed on yaml)
      ppm.createdAt = ppm.created_at;
      delete ppm.created_at;

      sortedShipments.push(ppm);
    }

    return sortedShipments.sort((a, b) => moment(a.createdAt) - moment(b.createdAt));
  }

  renderShipments = () => {
    const { currentOrders, match } = this.props;
    let hhgShipmentNumber = 1;
    return this.getSortedShipments.map((shipment) => {
      let receivingAgent;
      let releasingAgent;
      if (shipment.shipmentType === SHIPMENT_OPTIONS.PPM) {
        return (
          <PPMShipmentCard
            destinationZIP={shipment.destination_postal_code}
            editPath={`/moves/${match.params.moveId}/review/edit-date-and-location`}
            estimatedWeight="5,000"
            expectedDepartureDate={shipment.original_move_date}
            onEditClick={this.handleEditClick}
            originZIP={shipment.pickup_postal_code}
            shipmentId={shipment.id}
            sitDays={shipment.has_sit ? shipment.days_in_storage : ''}
          />
        );
      }

      if (shipment.agents) {
        receivingAgent = shipment.agents.find((agent) => (agent.agentType = 'RECEIVING_AGENT'));
        releasingAgent = shipment.agents.find((agent) => (agent.agentType = 'RELEASING_AGENT'));
      }

      return (
        <HHGShipmentCard
          destinationZIP={currentOrders.new_duty_station.address.postal_code}
          destinationLocation={shipment?.destinationAddress}
          pickupLocation={shipment.pickupAddress}
          receivingAgent={receivingAgent}
          releasingAgent={releasingAgent}
          remarks={shipment.remarks}
          requestedDeliveryDate={shipment.requestedDeliveryDate}
          requestedPickupDate={shipment.requestedPickupDate}
          shipmentId={shipment.id}
          shipmentNumber={hhgShipmentNumber++}
        />
      );
    });
  };

  render() {
    const {
      currentMove,
      currentPPM,
      mtoShipments,
      currentOrders,
      moveIsApproved,
      serviceMember,
      entitlement,
      match,
      history,
    } = this.props;
    console.log('hey', this.getSortedShipments);
    const currentStation = get(serviceMember, 'current_station');
    const stationPhone = get(currentStation, 'transportation_office.phone_lines.0');

    const rootAddressWithMoveId = `/moves/${match.params.moveId}`;
    const rootReviewAddressWithMoveId = rootAddressWithMoveId + `/review`;

    // isReviewPage being false is the same thing as being in the /edit route
    const isReviewPage = rootReviewAddressWithMoveId === match.url;
    const editSuccessBlurb = this.props.reviewState.editSuccess ? 'Your changes have been saved. ' : '';
    const editOrdersPath = rootReviewAddressWithMoveId + '/edit-orders';

    const showPPMShipmentSummary = !isReviewPage && Object.keys(currentPPM).length && currentPPM.status !== 'DRAFT';
    const showHHGShipmentSummary = isReviewPage && !!mtoShipments.length;
    const hasPPMorHHG = (isReviewPage && Object.keys(currentPPM).length) || !!mtoShipments.length;

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
        <ProfileTable
          affiliation={serviceMember.affiliation}
          city={serviceMember.residential_address.city}
          currentDutyStationName={serviceMember.current_station.name}
          edipi={serviceMember.edipi}
          email={serviceMember.personal_email}
          firstName={serviceMember.first_name}
          onEditClick={this.handleEditClick}
          lastName={serviceMember.last_name}
          postalCode={serviceMember.postal_code}
          rank={serviceMember.rank}
          state={serviceMember.residential_address.state}
          streetAddress1={serviceMember.residential_address.street_address_1}
          streetAddress2={serviceMember.residential_address.street_address_2}
          telephone={serviceMember.telephone}
        />

        <OrdersTable
          editPath={editOrdersPath}
          hasDependents={currentOrders.has_dependents}
          issueDate={currentOrders.issue_date}
          newDutyStationName={currentOrders.new_duty_station.name}
          onEditClick={this.handleEditClick}
          orderType={formatOrderType(currentOrders.orders_type)}
          reportByDate={currentOrders.report_by_date}
          uploads={currentOrders.uploaded_orders.uploads}
        />

        {showMoveSetup && <h3>Move setup</h3>}
        {isReviewPage && this.renderShipments()}
        {showPPMShipmentSummary && (
          <PPMShipmentSummary ppm={currentPPM} movePath={rootReviewAddressWithMoveId} orders={currentOrders} />
        )}
        {hasPPMorHHG && (
          <div className="grid-col-row margin-top-5">
            <span className="float-right">Optional</span>
            <h3>Add another shipment</h3>
            <p>Will you move any belongings to or from another location?</p>
            <Button className="usa-button--secondary" onClick={() => history.push(shipmentSelectionPath)}>
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
