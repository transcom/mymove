import React, { Component } from 'react';
import { connect } from 'react-redux';
import { withRouter, Link } from 'react-router-dom';
import { arrayOf, func, shape, bool, string } from 'prop-types';
import moment from 'moment';
import { Button, Grid } from '@trussworks/react-uswds';
import { generatePath } from 'react-router';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './Summary.module.scss';

import { customerRoutes } from 'constants/routes';
import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';
import { MOVE_STATUSES, SHIPMENT_OPTIONS } from 'shared/constants';
import { loadEntitlementsFromState } from 'shared/entitlements';
import ProfileTable from 'components/Customer/Review/ProfileTable/ProfileTable';
import OrdersTable from 'components/Customer/Review/OrdersTable/OrdersTable';
import PPMShipmentCard from 'components/Customer/Review/ShipmentCard/PPMShipmentCard/PPMShipmentCard';
import HHGShipmentCard from 'components/Customer/Review/ShipmentCard/HHGShipmentCard/HHGShipmentCard';
import SectionWrapper from 'components/Customer/SectionWrapper';
import NTSShipmentCard from 'components/Customer/Review/ShipmentCard/NTSShipmentCard/NTSShipmentCard';
import NTSRShipmentCard from 'components/Customer/Review/ShipmentCard/NTSRShipmentCard/NTSRShipmentCard';
import ConnectedAddShipmentModal from 'components/Customer/Review/AddShipmentModal/AddShipmentModal';
import ConnectedDestructiveShipmentConfirmationModal from 'components/ConfirmationModals/DestructiveShipmentConfirmationModal';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentOrders,
  selectCurrentMove,
  selectMoveIsApproved,
  selectHasCanceledMove,
  selectMTOShipmentsForCurrentMove,
} from 'store/entities/selectors';
import { deleteMTOShipment, getMTOShipmentsForMove } from 'services/internalApi';
import { updateMTOShipments } from 'store/entities/actions';
import { setFlashMessage } from 'store/flash/actions';
import { OrdersShape, MoveShape, MtoShipmentShape, HistoryShape, MatchShape } from 'types/customerShapes';

export class Summary extends Component {
  constructor(props) {
    super(props);

    this.state = {
      showModal: false,
      showDeleteModal: false,
      targetShipmentId: null,
    };
  }

  componentDidMount() {
    const { onDidMount, serviceMember } = this.props;

    if (onDidMount) {
      onDidMount(serviceMember.id);
    }
  }

  get getSortedShipments() {
    const { mtoShipments } = this.props;
    const sortedShipments = [...mtoShipments];

    return sortedShipments.sort((a, b) => moment(a.createdAt) - moment(b.createdAt));
  }

  handleEditClick = (path) => {
    const { history } = this.props;
    history.push(path);
  };

  handleDeleteClick = (shipmentId) => {
    this.setState({
      showDeleteModal: true,
      targetShipmentId: shipmentId,
    });
  };

  hideDeleteModal = () => {
    this.setState({
      showDeleteModal: false,
    });
  };

  handleDeleteShipmentConfirmation = (shipmentId) => {
    const { currentMove, updateShipmentList, setMsg } = this.props;
    deleteMTOShipment(shipmentId)
      .then(() => {
        getMTOShipmentsForMove(currentMove.id).then((response) => {
          updateShipmentList(response);
          setMsg('MTO_SHIPMENT_DELETE_SUCCESS', 'success', 'The shipment was deleted.', '', true);
        });
      })
      .catch(() => {
        setMsg(
          'MTO_SHIPMENT_DELETE_FAILURE',
          'error',
          'Something went wrong, and your changes were not saved. Please try again later or contact your counselor.',
          '',
          true,
        );
      })
      .finally(() => {
        this.setState({ showDeleteModal: false });
      });
  };

  renderShipments = () => {
    const { currentMove, currentOrders, match } = this.props;
    const { moveId } = match.params;
    const showEditAndDeleteBtn = currentMove.status === MOVE_STATUSES.DRAFT;
    let hhgShipmentNumber = 0;
    let ppmShipmentNumber = 0;
    return this.getSortedShipments.map((shipment) => {
      let receivingAgent;
      let releasingAgent;

      if (shipment.shipmentType === SHIPMENT_OPTIONS.PPM) {
        ppmShipmentNumber += 1;
        return (
          <PPMShipmentCard
            key={shipment.id}
            shipment={shipment}
            shipmentNumber={ppmShipmentNumber}
            showEditAndDeleteBtn={showEditAndDeleteBtn}
            onEditClick={this.handleEditClick}
            onDeleteClick={this.handleDeleteClick}
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
            showEditAndDeleteBtn={showEditAndDeleteBtn}
            moveId={moveId}
            onEditClick={this.handleEditClick}
            onDeleteClick={this.handleDeleteClick}
            pickupLocation={shipment.pickupAddress}
            secondaryPickupAddress={shipment?.secondaryPickupAddress}
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
            destinationZIP={currentOrders.new_duty_location.address.postalCode}
            secondaryDeliveryAddress={shipment?.secondaryDeliveryAddress}
            showEditAndDeleteBtn={showEditAndDeleteBtn}
            moveId={moveId}
            onEditClick={this.handleEditClick}
            onDeleteClick={this.handleDeleteClick}
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
          destinationZIP={currentOrders.new_duty_location.address.postalCode}
          secondaryDeliveryAddress={shipment?.secondaryDeliveryAddress}
          secondaryPickupAddress={shipment?.secondaryPickupAddress}
          destinationLocation={shipment?.destinationAddress}
          moveId={moveId}
          onEditClick={this.handleEditClick}
          onDeleteClick={this.handleDeleteClick}
          pickupLocation={shipment.pickupAddress}
          receivingAgent={receivingAgent}
          releasingAgent={releasingAgent}
          remarks={shipment.customerRemarks}
          requestedDeliveryDate={shipment.requestedDeliveryDate}
          requestedPickupDate={shipment.requestedPickupDate}
          shipmentId={shipment.id}
          shipmentNumber={hhgShipmentNumber}
          shipmentType={shipment.shipmentType}
          showEditAndDeleteBtn={showEditAndDeleteBtn}
        />
      );
    });
  };

  toggleModal = () => {
    this.setState((state) => ({
      showModal: !state.showModal,
    }));
  };

  render() {
    const { currentMove, currentOrders, match, moveIsApproved, mtoShipments, serviceMember } = this.props;
    const { showModal, showDeleteModal, targetShipmentId } = this.state;

    const { moveId } = match.params;
    const currentDutyLocation = serviceMember?.current_location;
    const officePhone = currentDutyLocation?.transportation_office?.phone_lines?.[0];

    const rootReviewAddressWithMoveId = generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId });

    // isReviewPage being false is the same thing as being in the /edit route
    const isReviewPage = rootReviewAddressWithMoveId === match.url;

    const showHHGShipmentSummary = isReviewPage && !!mtoShipments.length;

    // customer can add another shipment IFF the move is still draft
    const canAddAnotherShipment = isReviewPage && currentMove.status === MOVE_STATUSES.DRAFT;

    const showMoveSetup = showHHGShipmentSummary;
    const shipmentSelectionPath = generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, { moveId: currentMove.id });

    const thirdSectionHasContent = showMoveSetup || (isReviewPage && mtoShipments.length > 0);

    return (
      <>
        <ConnectedDestructiveShipmentConfirmationModal
          isOpen={showDeleteModal}
          shipmentID={targetShipmentId}
          onClose={this.hideDeleteModal}
          onSubmit={this.handleDeleteShipmentConfirmation}
          title="Delete this?"
          content="Your information will be gone. You’ll need to start over if you want it back."
          submitText="Yes, Delete"
          closeText="No, Keep It"
        />
        <SectionWrapper className={styles.SummarySectionWrapper}>
          <ProfileTable
            affiliation={ORDERS_BRANCH_OPTIONS[serviceMember?.affiliation] || ''}
            city={serviceMember.residential_address.city}
            currentDutyLocationName={currentOrders.origin_duty_location.name}
            edipi={serviceMember.edipi}
            email={serviceMember.personal_email}
            firstName={serviceMember.first_name}
            onEditClick={this.handleEditClick}
            lastName={serviceMember.last_name}
            postalCode={serviceMember.residential_address.postalCode}
            rank={ORDERS_RANK_OPTIONS[serviceMember?.rank] || ''}
            state={serviceMember.residential_address.state}
            streetAddress1={serviceMember.residential_address.streetAddress1}
            streetAddress2={serviceMember.residential_address.streetAddress2}
            telephone={serviceMember.telephone}
          />
        </SectionWrapper>
        <SectionWrapper className={styles.SummarySectionWrapper}>
          <OrdersTable
            hasDependents={currentOrders.has_dependents}
            issueDate={currentOrders.issue_date}
            moveId={moveId}
            newDutyLocationName={currentOrders.new_duty_location.name}
            onEditClick={this.handleEditClick}
            orderType={currentOrders.orders_type}
            reportByDate={currentOrders.report_by_date}
            uploads={currentOrders.uploaded_orders.uploads}
          />
        </SectionWrapper>
        {thirdSectionHasContent && (
          <SectionWrapper className={styles.SummarySectionWrapper}>
            {showMoveSetup && <h2 className={styles.moveSetup}>Move setup</h2>}
            {isReviewPage && this.renderShipments()}
          </SectionWrapper>
        )}
        {canAddAnotherShipment ? (
          <Grid row>
            <Grid col="fill" tablet={{ col: 'auto' }}>
              <Link to={shipmentSelectionPath} className="usa-link">
                Add another shipment
              </Link>
            </Grid>
            <Grid col="auto" className={styles.buttonContainer}>
              <Button
                title="Help with adding shipments"
                type="button"
                onClick={this.toggleModal}
                unstyled
                className={styles.buttonRight}
              >
                <FontAwesomeIcon icon={['far', 'question-circle']} />
              </Button>
            </Grid>
          </Grid>
        ) : (
          <p>Talk with your movers directly if you want to add or change shipments.</p>
        )}
        {moveIsApproved && (
          <div className="approved-edit-warning">
            *To change these fields, contact your local PPPO office at {currentDutyLocation.name}{' '}
            {officePhone ? ` at ${officePhone}` : ''}.
          </div>
        )}
        <ConnectedAddShipmentModal isOpen={showModal} closeModal={this.toggleModal} />
      </>
    );
  }
}

Summary.propTypes = {
  currentMove: MoveShape.isRequired,
  currentOrders: OrdersShape.isRequired,
  history: HistoryShape.isRequired,
  match: MatchShape.isRequired,
  moveIsApproved: bool.isRequired,
  mtoShipments: arrayOf(MtoShipmentShape).isRequired,
  onDidMount: func.isRequired,
  serviceMember: shape({ id: string.isRequired }).isRequired,
  updateShipmentList: func.isRequired,
  setMsg: func.isRequired,
};

function mapStateToProps(state) {
  return {
    mtoShipments: selectMTOShipmentsForCurrentMove(state),
    serviceMember: selectServiceMemberFromLoggedInUser(state),
    currentMove: selectCurrentMove(state) || {},
    currentOrders: selectCurrentOrders(state) || {},
    moveIsApproved: selectMoveIsApproved(state),
    lastMoveIsCanceled: selectHasCanceledMove(state),
    entitlement: loadEntitlementsFromState(state),
  };
}

const mapDispatchToProps = {
  updateShipmentList: updateMTOShipments,
  setMsg: setFlashMessage,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Summary));
