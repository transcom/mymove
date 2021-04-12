import React, { Component } from 'react';
import { get } from 'lodash';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import { arrayOf, func, shape, bool, string } from 'prop-types';
import moment from 'moment';
import { Button } from '@trussworks/react-uswds';
import { generatePath } from 'react-router';

import styles from './Summary.module.scss';

import { customerRoutes } from 'constants/routes';
import { validateEntitlement } from 'services/internalApi';
import ConnectedPPMShipmentSummary from 'scenes/Review/PPMShipmentSummary';
import { getInternalSwaggerDefinition } from 'shared/Swagger/selectors';
import { loadMove } from 'shared/Entities/modules/moves';
import { MOVE_STATUSES, SHIPMENT_OPTIONS, titleCase } from 'shared/constants';
import { loadEntitlementsFromState } from 'shared/entitlements';
import { formatOrderType } from 'utils/formatters';
import Alert from 'shared/Alert';
import ProfileTable from 'components/Customer/Review/ProfileTable';
import OrdersTable from 'components/Customer/Review/OrdersTable';
import PPMShipmentCard from 'components/Customer/Review/ShipmentCard/PPMShipmentCard';
import HHGShipmentCard from 'components/Customer/Review/ShipmentCard/HHGShipmentCard';
import SectionWrapper from 'components/Customer/SectionWrapper';
import NTSShipmentCard from 'components/Customer/Review/ShipmentCard/NTSShipmentCard';
import NTSRShipmentCard from 'components/Customer/Review/ShipmentCard/NTSRShipmentCard';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentOrders,
  selectCurrentMove,
  selectMoveIsApproved,
  selectHasCanceledMove,
  selectMTOShipmentsForCurrentMove,
  selectCurrentPPM,
} from 'store/entities/selectors';
import { OrdersShape, MoveShape, MtoShipmentShape, HistoryShape, MatchShape } from 'types/customerShapes';

export class Summary extends Component {
  constructor(props) {
    super(props);

    this.state = {
      entitlementWarning: null,
    };
  }

  componentDidMount() {
    const { onDidMount, serviceMember, currentPPM } = this.props;

    if (currentPPM) {
      this.checkEntitlement();
    }

    if (onDidMount) {
      onDidMount(serviceMember.id);
    }
  }

  componentDidUpdate(prevProps) {
    const { currentPPM } = this.props;

    // Only check entitlement for PPMs, not HHGs
    if (!prevProps.currentPPM && currentPPM) {
      this.checkEntitlement();
    }
  }

  get getSortedShipments() {
    const { currentPPM, mtoShipments } = this.props;
    const sortedShipments = [...mtoShipments];
    if (currentPPM) {
      const ppm = { ...currentPPM };
      ppm.shipmentType = SHIPMENT_OPTIONS.PPM;
      // workaround for differing cases between mtoShipments and ppms (bigger change needed on yaml)
      ppm.createdAt = ppm.created_at;
      delete ppm.created_at;

      sortedShipments.push(ppm);
    }

    return sortedShipments.sort((a, b) => moment(a.createdAt) - moment(b.createdAt));
  }

  checkEntitlement = () => {
    const { match } = this.props;
    const { entitlementWarning } = this.state;

    // Reset state
    if (entitlementWarning) {
      this.setState({
        entitlementWarning: null,
      });
    }

    validateEntitlement(match.params.moveId).catch((error) => {
      const { status, body } = error.response;

      if (status === 409) {
        this.setState({
          entitlementWarning: body?.message,
        });
      }
    });
  };

  handleEditClick = (path) => {
    const { history } = this.props;
    history.push(path);
  };

  renderShipments = () => {
    const { currentMove, currentOrders, match } = this.props;
    const { moveId } = match.params;
    const showEditBtn = currentMove.status === MOVE_STATUSES.DRAFT;
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
            showEditBtn={showEditBtn}
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
            destinationLocation={shipment?.destinationAddress}
            destinationZIP={currentOrders.new_duty_station.address.postal_code}
            showEditBtn={showEditBtn}
            moveId={moveId}
            onEditClick={this.handleEditClick}
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
          showEditBtn={showEditBtn}
        />
      );
    });
  };

  render() {
    const {
      currentMove,
      currentOrders,
      currentPPM,
      history,
      match,
      moveIsApproved,
      mtoShipments,
      serviceMember,
    } = this.props;
    const { entitlementWarning } = this.state;

    const { moveId } = match.params;
    const currentStation = get(serviceMember, 'current_station');
    const stationPhone = get(currentStation, 'transportation_office.phone_lines.0');

    const rootReviewAddressWithMoveId = generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId });

    // isReviewPage being false is the same thing as being in the /edit route
    const isReviewPage = rootReviewAddressWithMoveId === match.url;

    const showPPMShipmentSummary = !isReviewPage && currentPPM?.status !== 'DRAFT';
    const showHHGShipmentSummary = isReviewPage && !!mtoShipments.length;
    const hasPPM = !!currentPPM;

    // customer can add another shipment IFF the move is still draft OR it's not a draft & they don't have a PPM yet
    // double not is to prevent js from converting false to 0 and displaying said 0 on the page
    const canAddAnotherShipment = isReviewPage && !!(currentMove.status === MOVE_STATUSES.DRAFT || !hasPPM);

    const showMoveSetup = showPPMShipmentSummary || showHHGShipmentSummary;
    const shipmentSelectionPath = generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, { moveId: currentMove.id });

    return (
      <>
        {entitlementWarning && (
          <Alert type="warning" heading="Your estimated weight is above your entitlement.">
            {titleCase(entitlementWarning)}.
          </Alert>
        )}

        <SectionWrapper>
          <ProfileTable
            affiliation={serviceMember.affiliation}
            city={serviceMember.residential_address.city}
            currentDutyStationName={currentOrders.origin_duty_station.name}
            edipi={serviceMember.edipi}
            email={serviceMember.personal_email}
            firstName={serviceMember.first_name}
            onEditClick={this.handleEditClick}
            lastName={serviceMember.last_name}
            postalCode={serviceMember.residential_address.postal_code}
            rank={currentOrders.grade}
            state={serviceMember.residential_address.state}
            streetAddress1={serviceMember.residential_address.street_address_1}
            streetAddress2={serviceMember.residential_address.street_address_2}
            telephone={serviceMember.telephone}
          />
        </SectionWrapper>
        <SectionWrapper>
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
        </SectionWrapper>
        <SectionWrapper>
          {showMoveSetup && <h2 className={styles.moveSetup}>Move setup</h2>}
          {isReviewPage && this.renderShipments()}
          {showPPMShipmentSummary && (
            <ConnectedPPMShipmentSummary
              ppm={currentPPM}
              movePath={rootReviewAddressWithMoveId}
              orders={currentOrders}
            />
          )}
          <div className="grid-col-row margin-top-5">
            <span className="float-right">Optional</span>
            <h3>Add another shipment</h3>
            {canAddAnotherShipment ? (
              <>
                <p>Do you have more to move, either to or from another location, or by another method?</p>
                <Button className="usa-button--secondary" onClick={() => history.push(shipmentSelectionPath)}>
                  Add another shipment
                </Button>
              </>
            ) : (
              <p>Talk with your movers directly if you want to add or change shipments.</p>
            )}
          </div>
          {moveIsApproved && (
            <div className="approved-edit-warning">
              *To change these fields, contact your local PPPO office at {get(currentStation, 'name')}{' '}
              {stationPhone ? ` at ${stationPhone}` : ''}.
            </div>
          )}
        </SectionWrapper>
      </>
    );
  }
}

Summary.propTypes = {
  currentMove: MoveShape.isRequired,
  currentOrders: OrdersShape.isRequired,
  currentPPM: shape({}),
  history: HistoryShape.isRequired,
  match: MatchShape.isRequired,
  moveIsApproved: bool.isRequired,
  mtoShipments: arrayOf(MtoShipmentShape).isRequired,
  onDidMount: func.isRequired,
  serviceMember: shape({ id: string.isRequired }).isRequired,
};

Summary.defaultProps = {
  currentPPM: null,
};

function mapStateToProps(state) {
  return {
    currentPPM: selectCurrentPPM(state),
    mtoShipments: selectMTOShipmentsForCurrentMove(state),
    serviceMember: selectServiceMemberFromLoggedInUser(state),
    currentMove: selectCurrentMove(state) || {},
    currentOrders: selectCurrentOrders(state) || {},
    schemaRank: getInternalSwaggerDefinition(state, 'ServiceMemberRank'),
    schemaOrdersType: getInternalSwaggerDefinition(state, 'OrdersType'),
    schemaAffiliation: getInternalSwaggerDefinition(state, 'Affiliation'),
    moveIsApproved: selectMoveIsApproved(state),
    lastMoveIsCanceled: selectHasCanceledMove(state),
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
    showLoggedInUser: showLoggedInUserAction,
  };
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Summary));
