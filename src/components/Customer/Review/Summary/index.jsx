import React, { Component } from 'react';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import { arrayOf, func, shape, bool, string } from 'prop-types';
import moment from 'moment';
import { Button } from '@trussworks/react-uswds';

import styles from './Summary.module.scss';

import { checkEntitlement } from 'scenes/Review/ducks';
// eslint-disable-next-line import/no-named-as-default
import PPMShipmentSummary from 'scenes/Review/PPMShipmentSummary';
import { selectActivePPMForMove } from 'shared/Entities/modules/ppms';
import { getInternalSwaggerDefinition } from 'shared/Swagger/selectors';
import { loadMove, selectMove } from 'shared/Entities/modules/moves';
import { selectActiveOrLatestOrdersFromEntities, selectUploadsForActiveOrders } from 'shared/Entities/modules/orders';
import { SHIPMENT_OPTIONS, titleCase } from 'shared/constants';
import {
  moveIsApproved as selectMoveIsApproved,
  lastMoveIsCanceled,
  selectedMoveType as selectMoveType,
} from 'scenes/Moves/ducks';
import { loadEntitlementsFromState } from 'shared/entitlements';
import { formatOrderType } from 'utils/formatters';
import Alert from 'shared/Alert';
import ProfileTable from 'components/Customer/Review/ProfileTable';
import OrdersTable from 'components/Customer/Review/OrdersTable';
import PPMShipmentCard from 'components/Customer/Review/ShipmentCard/PPMShipmentCard';
import HHGShipmentCard from 'components/Customer/Review/ShipmentCard/HHGShipmentCard';
import NTSShipmentCard from 'components/Customer/Review/ShipmentCard/NTSShipmentCard';
import NTSRShipmentCard from 'components/Customer/Review/ShipmentCard/NTSRShipmentCard';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { selectMTOShipmentsByMoveId } from 'shared/Entities/modules/mtoShipments';

export class Summary extends Component {
  componentDidMount() {
    const { onDidMount, serviceMember } = this.props;
    if (onDidMount) {
      onDidMount(serviceMember.id);
    }
  }

  componentDidUpdate(prevProps) {
    const { currentPPM, selectedMoveType, onCheckEntitlement, match } = this.props;
    const hhgMove = !Object.keys(prevProps.currentPPM).length && !Object.keys(currentPPM).length;
    // Only check entitlement for PPMs, not HHGs
    if (prevProps.currentPPM !== currentPPM && !hhgMove && selectedMoveType === SHIPMENT_OPTIONS.PPM) {
      onCheckEntitlement(match.params.moveId);
    }
  }

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

  handleEditClick = (path) => {
    const { history } = this.props;
    history.push(path);
  };

  renderShipments = () => {
    const { currentOrders, match } = this.props;
    const { moveId } = match.params;
    let hhgShipmentNumber = 0;
    return this.getSortedShipments.map((shipment) => {
      let receivingAgent;
      let releasingAgent;
      if (shipment.shipmentType === SHIPMENT_OPTIONS.PPM) {
        return (
          <PPMShipmentCard
            key={shipment.id}
            destinationZIP={shipment.destination_postal_code}
            moveId={moveId}
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
        receivingAgent = shipment.agents.find((agent) => agent.agentType === 'RECEIVING_AGENT');
        releasingAgent = shipment.agents.find((agent) => agent.agentType === 'RELEASING_AGENT');
      }
      if (shipment.shipmentType === SHIPMENT_OPTIONS.NTS) {
        return (
          <NTSShipmentCard
            key={shipment.id}
            moveId={moveId}
            onEditClick={this.handleEditClick}
            pickupLocation={shipment.pickupAddress}
            releasingAgent={releasingAgent}
            remarks={shipment.customerRemarks}
            requestedPickupDate={shipment.requestedPickupDate}
            shipmentId={shipment.id}
            shipmentType={shipment.shipmentType}
          />
        );
      }
      if (shipment.shipmentType === SHIPMENT_OPTIONS.NTSR) {
        return (
          <NTSRShipmentCard
            key={shipment.id}
            destinationZIP={currentOrders.new_duty_station.address.postal_code}
            destinationLocation={shipment?.destinationAddress}
            receivingAgent={receivingAgent}
            remarks={shipment.customerRemarks}
            requestedDeliveryDate={shipment.requestedDeliveryDate}
            shipmentId={shipment.id}
            shipmentType={shipment.shipmentType}
          />
        );
      }
      hhgShipmentNumber += 1;
      return (
        <HHGShipmentCard
          key={shipment.id}
          destinationZIP={currentOrders.new_duty_station.address.postal_code}
          destinationLocation={shipment?.destinationAddress}
          moveId={moveId}
          onEditClick={this.handleEditClick}
          pickupLocation={shipment.pickupAddress}
          receivingAgent={receivingAgent}
          releasingAgent={releasingAgent}
          remarks={shipment.customerRemarks}
          requestedDeliveryDate={shipment.requestedDeliveryDate}
          requestedPickupDate={shipment.requestedPickupDate}
          shipmentId={shipment.id}
          shipmentNumber={hhgShipmentNumber}
          shipmentType={shipment.shipmentType}
        />
      );
    });
  };

  render() {
    const {
      currentMove,
      currentOrders,
      currentPPM,
      entitlement,
      history,
      match,
      moveIsApproved,
      mtoShipments,
      serviceMember,
      reviewState,
    } = this.props;
    const { moveId } = match.params;
    const currentStation = get(serviceMember, 'current_station');
    const stationPhone = get(currentStation, 'transportation_office.phone_lines.0');

    const rootAddressWithMoveId = `/moves/${moveId}`;
    const rootReviewAddressWithMoveId = `${rootAddressWithMoveId}/review`;

    // isReviewPage being false is the same thing as being in the /edit route
    const isReviewPage = rootReviewAddressWithMoveId === match.url;
    const editSuccessBlurb = reviewState.editSuccess ? 'Your changes have been saved. ' : '';

    const showPPMShipmentSummary = !isReviewPage && Object.keys(currentPPM).length && currentPPM.status !== 'DRAFT';
    const showHHGShipmentSummary = isReviewPage && !!mtoShipments.length;
    const hasPPMorHHG = (isReviewPage && Object.keys(currentPPM).length) || !!mtoShipments.length;

    const showMoveSetup = showPPMShipmentSummary || showHHGShipmentSummary;
    const shipmentSelectionPath = `/moves/${currentMove.id}/select-type`;
    return (
      <>
        {get(reviewState.error, 'statusCode', false) === 409 && (
          <Alert type="warning" heading={`${editSuccessBlurb}Your estimated weight is above your entitlement.`}>
            {titleCase(reviewState.error.response.body.message)}.
          </Alert>
        )}
        {reviewState.editSuccess &&
          !reviewState.entitlementChange &&
          get(reviewState.error, 'statusCode', false) === false && <Alert type="success" heading={editSuccessBlurb} />}
        {currentMove && reviewState.entitlementChange && get(reviewState.error, 'statusCode', false) === false && (
          <Alert type="info" heading={`${editSuccessBlurb}Note that the entitlement has also changed.`}>
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
          postalCode={serviceMember.residential_address.postal_code}
          rank={serviceMember.rank}
          state={serviceMember.residential_address.state}
          streetAddress1={serviceMember.residential_address.street_address_1}
          streetAddress2={serviceMember.residential_address.street_address_2}
          telephone={serviceMember.telephone}
        />

        <OrdersTable
          hasDependents={currentOrders.has_dependents}
          issueDate={currentOrders.issue_date}
          moveId={moveId}
          newDutyStationName={currentOrders.new_duty_station.name}
          onEditClick={this.handleEditClick}
          orderType={formatOrderType(currentOrders.orders_type)}
          reportByDate={currentOrders.report_by_date}
          uploads={currentOrders.uploaded_orders.uploads}
        />

        {showMoveSetup && <h2 className={styles.moveSetup}>Move setup</h2>}
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
      </>
    );
  }
}

Summary.propTypes = {
  currentMove: shape({ id: string.isRequired }).isRequired,
  currentOrders: shape({}).isRequired,
  currentPPM: shape({}).isRequired,
  history: shape({ push: func.isRequired }).isRequired,
  match: shape({ params: shape({ moveId: string.isRequired }) }).isRequired,
  moveIsApproved: bool.isRequired,
  mtoShipments: arrayOf(shape({})).isRequired,
  onCheckEntitlement: func.isRequired,
  onDidMount: func.isRequired,
  selectedMoveType: string.isRequired,
  serviceMember: shape({ id: string.isRequired }).isRequired,
  entitlement: shape({}).isRequired,
  reviewState: shape({
    editSuccess: bool,
    entitlementChange: bool,
    error: bool,
  }),
};

Summary.defaultProps = {
  reviewState: {},
};

function mapStateToProps(state, ownProps) {
  const moveID = ownProps.match.params.moveId;
  const currentOrders = selectActiveOrLatestOrdersFromEntities(state);

  return {
    currentPPM: selectActivePPMForMove(state, moveID),
    mtoShipments: selectMTOShipmentsByMoveId(state, moveID),
    serviceMember: state.serviceMember.currentServiceMember,
    currentMove: selectMove(state, moveID),
    currentOrders,
    uploads: selectUploadsForActiveOrders(state),
    selectedMoveType: selectMoveType(state),
    schemaRank: getInternalSwaggerDefinition(state, 'ServiceMemberRank'),
    schemaOrdersType: getInternalSwaggerDefinition(state, 'OrdersType'),
    schemaAffiliation: getInternalSwaggerDefinition(state, 'Affiliation'),
    moveIsApproved: selectMoveIsApproved(state),
    lastMoveIsCanceled: lastMoveIsCanceled(state),
    reviewState: state.review,
    entitlement: loadEntitlementsFromState(state),
  };
}
function mapDispatchToProps(dispatch, ownProps) {
  return {
    onDidMount() {
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
