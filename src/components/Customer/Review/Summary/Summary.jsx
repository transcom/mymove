import React, { Component } from 'react';
import { connect } from 'react-redux';
import { generatePath, Link, matchPath } from 'react-router-dom';
import { func, shape, bool, string } from 'prop-types';
import moment from 'moment';
import { Button, Grid } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { isBooleanFlagEnabled } from '../../../../utils/featureFlags';
import { FEATURE_FLAG_KEYS, MOVE_STATUSES, SHIPMENT_OPTIONS, SHIPMENT_TYPES } from '../../../../shared/constants';

import styles from './Summary.module.scss';

import ConnectedDestructiveShipmentConfirmationModal from 'components/ConfirmationModals/DestructiveShipmentConfirmationModal';
import ConnectedAddShipmentModal from 'components/Customer/Review/AddShipmentModal/AddShipmentModal';
import ConnectedIncompleteShipmentModal from 'components/Customer/Review/IncompleteShipmentModal/IncompleteShipmentModal';
import OrdersTable from 'components/Customer/Review/OrdersTable/OrdersTable';
import ProfileTable from 'components/Customer/Review/ProfileTable/ProfileTable';
import HHGShipmentCard from 'components/Customer/Review/ShipmentCard/HHGShipmentCard/HHGShipmentCard';
import NTSRShipmentCard from 'components/Customer/Review/ShipmentCard/NTSRShipmentCard/NTSRShipmentCard';
import NTSShipmentCard from 'components/Customer/Review/ShipmentCard/NTSShipmentCard/NTSShipmentCard';
import PPMShipmentCard from 'components/Customer/Review/ShipmentCard/PPMShipmentCard/PPMShipmentCard';
import BoatShipmentCard from 'components/Customer/Review/ShipmentCard/BoatShipmentCard/BoatShipmentCard';
import MobileHomeShipmentCard from 'components/Customer/Review/ShipmentCard/MobileHomeShipmentCard/MobileHomeShipmentCard';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { ORDERS_BRANCH_OPTIONS, ORDERS_PAY_GRADE_OPTIONS } from 'constants/orders';
import { customerRoutes } from 'constants/routes';
import { deleteMTOShipment, getAllMoves, getMTOShipmentsForMove } from 'services/internalApi';
import { loadEntitlementsFromState } from 'shared/entitlements';
import { updateMTOShipments, updateAllMoves as updateAllMovesAction } from 'store/entities/actions';
import {
  selectServiceMemberFromLoggedInUser,
  selectMoveIsApproved,
  selectHasCanceledMove,
  selectAllMoves,
  selectCurrentMoveFromAllMoves,
  selectShipmentsFromMove,
} from 'store/entities/selectors';
import { setFlashMessage } from 'store/flash/actions';
import withRouter from 'utils/routing';
import { RouterShape } from 'types';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

export class Summary extends Component {
  constructor(props) {
    super(props);

    this.state = {
      showModal: false,
      showIncompletePPMModal: false,
      showDeleteModal: false,
      targetShipmentId: null,
      targetShipmentLabel: null,
      targetShipmentMoveCode: null,
      targetShipmentType: null,
      enablePPM: true,
      enableNTS: true,
      enableNTSR: true,
      enableBoat: true,
      enableMobileHome: true,
      enableUB: true,
    };
  }

  componentDidMount() {
    const { onDidMount, serviceMember, updateAllMoves } = this.props;

    if (onDidMount) {
      onDidMount(serviceMember.id);
    }

    getAllMoves(serviceMember.id).then((response) => {
      updateAllMoves(response);
    });

    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.PPM).then((enabled) => {
      this.setState({
        enablePPM: enabled,
      });
    });
    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.NTS).then((enabled) => {
      this.setState({
        enableNTS: enabled,
      });
    });
    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.NTSR).then((enabled) => {
      this.setState({
        enableNTSR: enabled,
      });
    });
    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.BOAT).then((enabled) => {
      this.setState({
        enableBoat: enabled,
      });
    });
    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.MOBILE_HOME).then((enabled) => {
      this.setState({
        enableMobileHome: enabled,
      });
    });
    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.UNACCOMPANIED_BAGGAGE).then((enabled) => {
      this.setState({
        enableUB: enabled,
      });
    });
  }

  handleEditClick = (path) => {
    const { router } = this.props;
    const { state } = this;
    router.navigate(path, { state });
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
    const { serviceMember, updateAllMoves, moveId, updateShipmentList, setMsg } = this.props;
    deleteMTOShipment(shipmentId)
      .then(() => {
        getAllMoves(serviceMember.id).then((res) => {
          updateAllMoves(res);
        });
        getMTOShipmentsForMove(moveId).then((response) => {
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
    const { router, serviceMember, serviceMemberMoves } = this.props;
    const { moveId } = router.params;

    const currentMove = selectCurrentMoveFromAllMoves(serviceMemberMoves, moveId);
    const { mtoShipments } = currentMove ?? {};
    const { orders } = currentMove ?? {};
    const currentOrders = orders;
    this.state = { ...this.state, moveId };

    const sortedShipments = mtoShipments.sort((a, b) => moment(a.createdAt) - moment(b.createdAt));

    // loading placeholder while data loads - this handles any async issues
    if (!currentMove || !mtoShipments) {
      return (
        <div className={styles.homeContainer}>
          <div className={`usa-prose grid-container ${styles['grid-container']}`}>
            <LoadingPlaceholder />
          </div>
        </div>
      );
    }

    const showEditAndDeleteBtn = currentMove.status === MOVE_STATUSES.DRAFT;
    let hhgShipmentNumber = 0;
    let ppmShipmentNumber = 0;
    let boatShipmentNumber = 0;
    let mobileHomeShipmentNumber = 0;
    let ubShipmentNumber = 0;
    return sortedShipments.map((shipment) => {
      let receivingAgent;
      let releasingAgent;

      if (shipment.shipmentType === SHIPMENT_OPTIONS.PPM) {
        ppmShipmentNumber += 1;
        return (
          <PPMShipmentCard
            key={shipment.id}
            move={currentMove}
            affiliation={serviceMember.affiliation}
            shipment={shipment}
            shipmentNumber={ppmShipmentNumber}
            showEditAndDeleteBtn={showEditAndDeleteBtn}
            onEditClick={this.handleEditClick}
            onDeleteClick={this.handleDeleteClick}
            onIncompleteClick={this.toggleIncompleteShipmentModal}
            marketCode={shipment.marketCode}
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
            tertiaryPickupAddress={shipment?.tertiaryPickupAddress}
            releasingAgent={releasingAgent}
            remarks={shipment.customerRemarks}
            requestedPickupDate={shipment.requestedPickupDate}
            shipmentId={shipment.id}
            shipmentType={shipment.shipmentType}
            status={shipment.status}
            onIncompleteClick={this.toggleIncompleteShipmentModal}
            shipmentLocator={shipment.shipmentLocator}
            marketCode={shipment.marketCode}
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
            tertiaryDeliveryAddress={shipment?.tertiaryDeliveryAddress}
            showEditAndDeleteBtn={showEditAndDeleteBtn}
            moveId={moveId}
            onEditClick={this.handleEditClick}
            onDeleteClick={this.handleDeleteClick}
            receivingAgent={receivingAgent}
            remarks={shipment.customerRemarks}
            requestedDeliveryDate={shipment.requestedDeliveryDate}
            shipmentId={shipment.id}
            shipmentType={shipment.shipmentType}
            status={shipment.status}
            onIncompleteClick={this.toggleIncompleteShipmentModal}
            shipmentLocator={shipment.shipmentLocator}
            marketCode={shipment.marketCode}
          />
        );
      }
      if (
        shipment.shipmentType === SHIPMENT_TYPES.BOAT_TOW_AWAY ||
        shipment.shipmentType === SHIPMENT_TYPES.BOAT_HAUL_AWAY
      ) {
        boatShipmentNumber += 1;
        return (
          <BoatShipmentCard
            key={shipment.id}
            shipment={shipment}
            destinationZIP={currentOrders.new_duty_location.address.postalCode}
            secondaryDeliveryAddress={shipment?.secondaryDeliveryAddress}
            tertiaryDeliveryAddress={shipment?.tertiaryDeliveryAddress}
            secondaryPickupAddress={shipment?.secondaryPickupAddress}
            tertiaryPickupAddress={shipment?.tertiaryPickupAddress}
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
            shipmentNumber={boatShipmentNumber}
            showEditAndDeleteBtn={showEditAndDeleteBtn}
            status={shipment.status}
            onIncompleteClick={this.toggleIncompleteShipmentModal}
            marketCode={shipment.marketCode}
          />
        );
      }
      if (shipment.shipmentType === SHIPMENT_TYPES.MOBILE_HOME) {
        mobileHomeShipmentNumber += 1;
        return (
          <MobileHomeShipmentCard
            key={shipment.id}
            shipment={shipment}
            destinationZIP={currentOrders.new_duty_location.address.postalCode}
            secondaryDeliveryAddress={shipment?.secondaryDeliveryAddress}
            tertiaryDeliveryAddress={shipment?.tertiaryDeliveryAddress}
            secondaryPickupAddress={shipment?.secondaryPickupAddress}
            tertiaryPickupAddress={shipment?.tertiaryPickupAddress}
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
            shipmentNumber={mobileHomeShipmentNumber}
            showEditAndDeleteBtn={showEditAndDeleteBtn}
            status={shipment.status}
            onIncompleteClick={this.toggleIncompleteShipmentModal}
            marketCode={shipment.marketCode}
          />
        );
      }
      if (shipment.shipmentType === SHIPMENT_TYPES.UNACCOMPANIED_BAGGAGE) {
        ubShipmentNumber += 1;
        return (
          <HHGShipmentCard
            key={shipment.id}
            destinationZIP={currentOrders.new_duty_location.address.postalCode}
            secondaryDeliveryAddress={shipment?.secondaryDeliveryAddress}
            tertiaryDeliveryAddress={shipment?.tertiaryDeliveryAddress}
            secondaryPickupAddress={shipment?.secondaryPickupAddress}
            tertiaryPickupAddress={shipment?.tertiaryPickupAddress}
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
            shipmentLocator={shipment.shipmentLocator}
            shipmentNumber={ubShipmentNumber}
            shipmentType={shipment.shipmentType}
            showEditAndDeleteBtn={showEditAndDeleteBtn}
            status={shipment.status}
            onIncompleteClick={this.toggleIncompleteShipmentModal}
            marketCode={shipment.marketCode}
          />
        );
      }
      hhgShipmentNumber += 1;
      return (
        <HHGShipmentCard
          key={shipment.id}
          destinationZIP={currentOrders.new_duty_location.address.postalCode}
          secondaryDeliveryAddress={shipment?.secondaryDeliveryAddress}
          tertiaryDeliveryAddress={shipment?.tertiaryDeliveryAddress}
          secondaryPickupAddress={shipment?.secondaryPickupAddress}
          tertiaryPickupAddress={shipment?.tertiaryPickupAddress}
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
          shipmentLocator={shipment.shipmentLocator}
          shipmentNumber={hhgShipmentNumber}
          shipmentType={shipment.shipmentType}
          showEditAndDeleteBtn={showEditAndDeleteBtn}
          status={shipment.status}
          onIncompleteClick={this.toggleIncompleteShipmentModal}
          marketCode={shipment.marketCode}
        />
      );
    });
  };

  toggleModal = () => {
    this.setState((state) => ({
      showModal: !state.showModal,
    }));
  };

  toggleIncompleteShipmentModal = (ShipmentLabel, shipmentMoveCode, shipmentType) => {
    this.setState((state) => ({
      showIncompletePPMModal: !state.showIncompletePPMModal,
      targetShipmentLabel: ShipmentLabel,
      targetShipmentMoveCode: shipmentMoveCode,
      targetShipmentType: shipmentType,
    }));
  };

  render() {
    const { serviceMemberMoves, router, moveIsApproved, serviceMember } = this.props;
    const {
      showModal,
      showDeleteModal,
      targetShipmentId,
      showIncompletePPMModal,
      targetShipmentLabel,
      targetShipmentMoveCode,
      targetShipmentType,
      enablePPM,
      enableNTS,
      enableNTSR,
      enableBoat,
      enableMobileHome,
    } = this.state;

    const { pathname } = router.location;
    const { moveId } = router.params;

    const currentMove = selectCurrentMoveFromAllMoves(serviceMemberMoves, moveId);
    const mtoShipments = selectShipmentsFromMove(currentMove);
    const { orders } = currentMove ?? {};
    const currentOrders = orders;

    // loading placeholder while data loads - this handles any async issues
    if (!currentMove || !mtoShipments) {
      return (
        <div className={styles.homeContainer}>
          <div className={`usa-prose grid-container ${styles['grid-container']}`}>
            <LoadingPlaceholder />
          </div>
        </div>
      );
    }

    const currentDutyLocation = orders?.origin_duty_location?.transportation_office;
    const officePhone = currentDutyLocation?.transportation_office?.phone_lines?.[0];

    const rootReviewAddressWithMoveId = generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId });

    // isReviewPage being false is the same thing as being in the /edit route
    const isReviewPage = matchPath(
      {
        path: rootReviewAddressWithMoveId,
        end: true,
      },
      pathname,
    );

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
        <ConnectedIncompleteShipmentModal
          isOpen={showIncompletePPMModal}
          closeModal={this.toggleIncompleteShipmentModal}
          shipmentLabel={targetShipmentLabel}
          shipmentMoveCode={targetShipmentMoveCode}
          shipmentType={targetShipmentType}
        />
        <SectionWrapper className={styles.SummarySectionWrapper}>
          <ProfileTable
            affiliation={ORDERS_BRANCH_OPTIONS[serviceMember?.affiliation] || ''}
            city={serviceMember.residential_address.city}
            edipi={serviceMember.edipi}
            email={serviceMember.personal_email}
            firstName={serviceMember.first_name}
            onEditClick={this.handleEditClick}
            lastName={serviceMember.last_name}
            postalCode={serviceMember.residential_address.postalCode}
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
            payGrade={ORDERS_PAY_GRADE_OPTIONS[currentOrders?.grade] || ''}
            originDutyLocationName={currentOrders.origin_duty_location.name}
            orderId={currentOrders.id}
            counselingOfficeName={currentMove.counselingOffice?.name || ''}
            accompaniedTour={currentOrders.entitlement?.accompanied_tour}
            dependentsUnderTwelve={currentOrders.entitlement?.dependents_under_twelve}
            dependentsTwelveAndOver={currentOrders.entitlement?.dependents_twelve_and_over}
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
                <FontAwesomeIcon icon={['far', 'circle-question']} />
              </Button>
            </Grid>
          </Grid>
        ) : (
          <p>Talk with your movers directly if you want to add or change shipments.</p>
        )}
        {moveIsApproved && currentDutyLocation && (
          <p>
            *To change these fields, contact your local PPPO office at {currentDutyLocation?.name}
            {officePhone ? ` at ${officePhone}` : ''}.
          </p>
        )}
        {moveIsApproved && !currentDutyLocation && <p>*To change these fields, contact your local PPPO office.</p>}
        <ConnectedAddShipmentModal
          isOpen={showModal}
          closeModal={this.toggleModal}
          enablePPM={enablePPM}
          enableNTS={enableNTS}
          enableNTSR={enableNTSR}
          enableBoat={enableBoat}
          enableMobileHome={enableMobileHome}
        />
      </>
    );
  }
}

Summary.propTypes = {
  router: RouterShape,
  moveIsApproved: bool.isRequired,
  onDidMount: func.isRequired,
  serviceMember: shape({ id: string.isRequired }).isRequired,
  updateShipmentList: func.isRequired,
  setMsg: func.isRequired,
};

Summary.defaultProps = {
  router: {},
};

function mapStateToProps(state, ownProps) {
  const serviceMemberMoves = selectAllMoves(state);
  const {
    router: {
      params: { moveId },
    },
  } = ownProps;
  const currentMove = selectCurrentMoveFromAllMoves(state, moveId);
  const mtoShipments = currentMove?.mtoShipments ?? [];
  const currentOrders = currentMove?.orders ?? {};

  return {
    serviceMemberMoves,
    mtoShipments,
    serviceMember: selectServiceMemberFromLoggedInUser(state),
    currentMove,
    currentOrders,
    moveId,
    moveIsApproved: selectMoveIsApproved(state),
    lastMoveIsCanceled: selectHasCanceledMove(state),
    entitlement: loadEntitlementsFromState(state),
  };
}

const mapDispatchToProps = {
  updateShipmentList: updateMTOShipments,
  updateAllMoves: updateAllMovesAction,
  setMsg: setFlashMessage,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Summary));
