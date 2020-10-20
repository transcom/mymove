/* eslint-disable camelcase */
import React, { Component } from 'react';
import { arrayOf, bool, shape, string, node, oneOfType } from 'prop-types';
import moment from 'moment';
import { connect } from 'react-redux';
import { get, isEmpty } from 'lodash';

import styles from './Home.module.scss';
import {
  HelperNeedsOrders,
  HelperNeedsShipment,
  HelperNeedsSubmitMove,
  HelperSubmittedMove,
  HelperSubmittedNoPPM,
} from './HomeHelpers';

import { withContext } from 'shared/AppContext';
import { getNextIncompletePage as getNextIncompletePageInternal } from 'scenes/MyMove/getWorkflowRoutes';
import Alert from 'shared/Alert';
import PpmAlert from 'scenes/PpmLanding/PpmAlert';
import SignIn from 'shared/User/SignIn';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import Step from 'components/Customer/Home/Step';
import DocsUploaded from 'components/Customer/Home/DocsUploaded';
import ShipmentList from 'components/Customer/Home/ShipmentList';
import Contact from 'components/Customer/Home/Contact';
import { isProfileComplete as isProfileCompleteCheck } from 'scenes/ServiceMembers/ducks';
import { selectServiceMemberFromLoggedInUser } from 'shared/Entities/modules/serviceMembers';
import { selectUploadedOrders, selectActiveOrLatestOrdersFromEntities } from 'shared/Entities/modules/orders';
import { selectActiveOrLatestMove } from 'shared/Entities/modules/moves';
import { selectMTOShipmentsByMoveId, selectMTOShipmentForMTO } from 'shared/Entities/modules/mtoShipments';
import { SHIPMENT_OPTIONS, MOVE_TYPES } from 'shared/constants';
import { selectActivePPMForMove } from 'shared/Entities/modules/ppms';
import {
  selectCurrentUser,
  selectGetCurrentUserIsError,
  selectGetCurrentUserIsLoading,
  selectGetCurrentUserIsSuccess,
} from 'shared/Data/users';

const Description = ({ children }) => <p className={styles.description}>{children}</p>;

Description.propTypes = {
  children: node.isRequired,
};

class Home extends Component {
  componentDidUpdate(prevProps) {
    const { serviceMember, loggedInUserSuccess, isProfileComplete } = this.props;
    if (!prevProps.loggedInUserSuccess && loggedInUserSuccess) {
      if (!isEmpty(serviceMember) && !isProfileComplete) {
        // If the service member exists, but is not complete, redirect to next incomplete page.
        this.resumeMove();
      }
    }

    if (isEmpty(prevProps.serviceMember) && !isEmpty(serviceMember) && !isProfileComplete) {
      this.resumeMove();
    }

    if (!isEmpty(prevProps.serviceMember) && prevProps.serviceMember !== serviceMember && !isProfileComplete) {
      // if service member existed but was updated, redirect to next incomplete page.
      this.resumeMove();
    }
  }

  get hasOrders() {
    const { orders, uploadedOrderDocuments } = this.props;
    return !!Object.keys(orders).length && !!uploadedOrderDocuments.length;
  }

  get hasOrdersNoUpload() {
    const { orders, uploadedOrderDocuments } = this.props;
    return !!Object.keys(orders).length && !uploadedOrderDocuments.length;
  }

  get hasAnyShipments() {
    const { mtoShipments, currentPpm } = this.props;
    return (this.hasOrders && !!mtoShipments.length) || !!Object.keys(currentPpm).length;
  }

  get hasSubmittedMove() {
    const { move } = this.props;
    return !!Object.keys(move).length && move.status !== 'DRAFT';
  }

  get hasHHGShipment() {
    const { mtoShipments } = this.props;
    return mtoShipments.some((s) => s.shipmentType === SHIPMENT_OPTIONS.HHG);
  }

  get hasNTSShipment() {
    const { mtoShipments } = this.props;
    return mtoShipments.some((s) => s.shipmentType === MOVE_TYPES.NTS);
  }

  get hasPPMShipment() {
    const { currentPpm } = this.props;
    return !!Object.keys(currentPpm).length;
  }

  get shipmentActionBtnLabel() {
    if (this.hasSubmittedMove && this.hasPPMShipment) {
      return '';
    }
    if (this.hasAnyShipments) {
      return 'Add another shipment';
    }
    return 'Plan your shipments';
  }

  resumeMove = () => {
    const { history } = this.props;
    history.push(this.getNextIncompletePage());
  };

  getNextIncompletePage = () => {
    const {
      selectedMoveType,
      lastMoveIsCanceled,
      serviceMember,
      orders,
      uploadedOrderDocuments,
      move,
      currentPpm,
      mtoShipment,
      backupContacts,
      context,
    } = this.props;
    return getNextIncompletePageInternal({
      selectedMoveType,
      lastMoveIsCanceled,
      serviceMember,
      orders,
      uploads: uploadedOrderDocuments,
      move,
      currentPpm,
      mtoShipment,
      backupContacts,
      context,
    });
  };

  renderHelper = () => {
    if (!this.hasOrders) return <HelperNeedsOrders />;
    if (!this.hasAnyShipments) return <HelperNeedsShipment />;
    if (!this.hasSubmittedMove) return <HelperNeedsSubmitMove />;
    // TODO: support PPM shipments; see MB-4267
    if (this.hasPPMShipment) return <HelperSubmittedMove />;
    return <HelperSubmittedNoPPM />;
  };

  renderCustomerHeader = () => {
    const { serviceMember, orders } = this.props;
    if (!this.hasOrders) {
      return (
        <p>
          You&apos;re leaving <strong>{serviceMember?.current_station?.name}</strong>
        </p>
      );
    }
    return (
      <p>
        You&apos;re moving to <strong>{orders.new_duty_station.name}</strong> from{' '}
        <strong>{serviceMember.current_station.name}.</strong> Report by{' '}
        <strong>{moment(orders.report_by_date).format('DD MMM YYYY')}.</strong>
        <br />
        Weight allowance: <strong>{serviceMember.weight_allotment.total_weight_self} lbs</strong>
      </p>
    );
  };

  handleShipmentClick = (shipmentId, shipmentNumber, shipmentType) => {
    const { move, history } = this.props;
    let queryString = '';
    if (shipmentNumber) {
      queryString = `?shipmentNumber=${shipmentNumber}`;
    }

    let destLink = '';
    if (shipmentType === 'PPM') {
      destLink = `/moves/${move.id}/review/edit-date-and-location`;
    } else if (shipmentType === 'HHG') {
      destLink = `/moves/${move.id}/mto-shipments/${shipmentId}/edit-shipment${queryString}`;
    }

    history.push(destLink);
  };

  handleNewPathClick = (path) => {
    const { history } = this.props;
    history.push(path);
  };

  renderAlert = (loggedInUserError, createdServiceMemberError, moveSubmitSuccess, currentPpm) => {
    return (
      <div>
        {moveSubmitSuccess && !currentPpm && (
          <Alert type="success" heading="Success">
            You&apos;ve submitted your move
          </Alert>
        )}
        {currentPpm && moveSubmitSuccess && <PpmAlert heading="Congrats - your move is submitted!" />}
        {loggedInUserError && (
          <Alert type="error" heading="An error occurred">
            There was an error loading your user information.
          </Alert>
        )}
        {createdServiceMemberError && (
          <Alert type="error" heading="An error occurred">
            There was an error creating your profile information.
          </Alert>
        )}
      </div>
    );
  };

  sortAllShipments = (mtoShipments, currentPpm) => {
    const allShipments = JSON.parse(JSON.stringify(mtoShipments));
    if (Object.keys(currentPpm).length) {
      const ppm = JSON.parse(JSON.stringify(currentPpm));
      ppm.shipmentType = MOVE_TYPES.PPM;
      // workaround for differing cases between mtoShipments and ppms (bigger change needed on yaml)
      ppm.createdAt = ppm.created_at;
      delete ppm.created_at;

      allShipments.push(ppm);
    }
    allShipments.sort((a, b) => moment(a.createdAt) - moment(b.createdAt));

    return allShipments;
  };

  render() {
    const {
      isLoggedIn,
      loggedInUserIsLoading,
      loggedInUserSuccess,
      loggedInUserError,
      isProfileComplete,
      createdServiceMemberError,
      moveSubmitSuccess,
      serviceMember,
      move,
      uploadedOrderDocuments,
      mtoShipments,
      currentPpm,
      location,
    } = this.props;
    const ordersPath = this.hasOrdersNoUpload ? '/orders/upload' : '/orders';
    const shipmentSelectionPath = this.hasAnyShipments
      ? `/moves/${move.id}/select-type`
      : `/moves/${move.id}/moving-info`;
    const confirmationPath = `/moves/${move.id}/review`;
    const profileEditPath = '/moves/review/edit-profile';
    const ordersEditPath = `/moves/${move.id}/review/edit-orders`;
    const allSortedShipments = this.sortAllShipments(mtoShipments, currentPpm);
    return (
      <div className={`usa-prose grid-container ${styles['grid-container']}`}>
        {loggedInUserIsLoading && <LoadingPlaceholder />}
        {!isLoggedIn && !loggedInUserIsLoading && <SignIn location={location} />}
        {isLoggedIn && !isEmpty(serviceMember) && isProfileComplete && (
          <>
            <header data-testid="customer-header" className={styles['customer-header']}>
              <h2>
                {serviceMember?.first_name} {serviceMember?.last_name}
              </h2>
              {this.renderCustomerHeader()}
            </header>
            {loggedInUserSuccess && (
              <>
                {this.renderAlert(loggedInUserError, createdServiceMemberError, moveSubmitSuccess, currentPpm)}
                {this.renderHelper()}
                <Step
                  complete={serviceMember.is_profile_complete}
                  completedHeaderText="Profile complete"
                  editBtnLabel="Edit"
                  headerText="Profile complete"
                  step="1"
                  onEditBtnClick={() => this.handleNewPathClick(profileEditPath)}
                >
                  <Description>Make sure to keep your personal information up to date during your move</Description>
                </Step>
                <Step
                  complete={this.hasOrders}
                  completedHeaderText="Orders uploaded"
                  editBtnLabel={this.hasOrders ? 'Edit' : ''}
                  onEditBtnClick={() => this.handleNewPathClick(ordersEditPath)}
                  headerText="Upload orders"
                  actionBtnLabel={!this.hasOrders ? 'Add orders' : ''}
                  onActionBtnClick={() => this.handleNewPathClick(ordersPath)}
                  step="2"
                >
                  {this.hasOrders ? (
                    <DocsUploaded files={uploadedOrderDocuments} />
                  ) : (
                    <Description>Upload photos of each page, or upload a PDF.</Description>
                  )}
                </Step>
                <Step
                  actionBtnLabel={this.shipmentActionBtnLabel}
                  actionBtnDisabled={!this.hasOrders || (this.hasSubmittedMove && this.doesPpmAlreadyExist)}
                  onActionBtnClick={() => this.handleNewPathClick(shipmentSelectionPath)}
                  complete={this.hasAnyShipments}
                  completedHeaderText="Shipments"
                  headerText="Shipment selection"
                  secondaryBtn={this.hasAnyShipments}
                  secondaryClassName="margin-top-2"
                  step="3"
                >
                  {this.hasAnyShipments ? (
                    <div>
                      {this.hasSubmittedMove && !this.doesPpmAlreadyExist && (
                        <p className={styles.descriptionExtra}>If you need to add shipments, let your movers know.</p>
                      )}
                      <ShipmentList
                        shipments={allSortedShipments}
                        onShipmentClick={this.handleShipmentClick}
                        moveSubmitted={this.hasSubmittedMove}
                      />
                    </div>
                  ) : (
                    <Description>
                      Tell us where you&apos;re going and when you want to get there. We&apos;ll help you set up
                      shipments to make it work.
                    </Description>
                  )}
                </Step>
                <Step
                  complete={this.hasSubmittedMove}
                  actionBtnDisabled={!this.hasAnyShipments}
                  actionBtnLabel={!this.hasSubmittedMove ? 'Review and submit' : ''}
                  containerClassName="margin-bottom-8"
                  headerText="Confirm move request"
                  completedHeaderText="Move request confirmed"
                  onActionBtnClick={() => this.handleNewPathClick(confirmationPath)}
                  step="4"
                >
                  {this.hasSubmittedMove ? (
                    <Description>Move submitted.</Description>
                  ) : (
                    <Description>
                      Review your move details and sign the legal paperwork, then send the info on to your move
                      counselor.
                    </Description>
                  )}
                </Step>
                <Contact
                  header="Contacts"
                  dutyStationName="Seymour Johnson AFB"
                  officeType="Origin Transportation Office"
                  telephone="(919) 722-5458"
                />
              </>
            )}
          </>
        )}
      </div>
    );
  }
}

Home.propTypes = {
  orders: shape({}).isRequired,
  serviceMember: shape({
    first_name: string,
    last_name: string,
  }).isRequired,
  mtoShipments: arrayOf(
    shape({
      id: string,
      shipmentType: string,
    }),
  ).isRequired,
  currentPpm: shape({
    id: string,
    shipmentType: string,
  }).isRequired,
  uploadedOrderDocuments: arrayOf(
    shape({
      filename: string.isRequired,
    }),
  ).isRequired,
  history: shape({}).isRequired,
  move: shape({}).isRequired,
  isLoggedIn: bool.isRequired,
  loggedInUserIsLoading: bool.isRequired,
  loggedInUserSuccess: bool.isRequired,
  loggedInUserError: bool.isRequired,
  isProfileComplete: bool.isRequired,
  createdServiceMemberError: string,
  moveSubmitSuccess: bool.isRequired,
  location: shape({}).isRequired,
  selectedMoveType: string,
  lastMoveIsCanceled: bool,
  backupContacts: arrayOf(oneOfType([string, shape({})])),
  context: shape({
    flags: shape({
      hhgFlow: bool,
      ghcFlow: bool,
    }),
  }),
  mtoShipment: shape({}).isRequired,
};

Home.defaultProps = {
  createdServiceMemberError: '',
  selectedMoveType: '',
  lastMoveIsCanceled: false,
  backupContacts: [],
  context: {
    flags: {
      hhgFlow: false,
      ghcFlow: false,
    },
  },
};

const mapStateToProps = (state) => {
  const user = selectCurrentUser(state);
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const move = selectActiveOrLatestMove(state);

  return {
    currentPpm: selectActivePPMForMove(state, move.id),
    isLoggedIn: user.isLoggedIn,
    loggedInUserIsLoading: selectGetCurrentUserIsLoading(state),
    loggedInUserSuccess: selectGetCurrentUserIsSuccess(state),
    loggedInUserError: selectGetCurrentUserIsError(state),
    createdServiceMemberError: state.serviceMember.error,
    isProfileComplete: isProfileCompleteCheck(state),
    moveSubmitSuccess: state.signedCertification.moveSubmitSuccess,
    orders: selectActiveOrLatestOrdersFromEntities(state),
    uploadedOrderDocuments: selectUploadedOrders(state),
    serviceMember,
    backupContacts: serviceMember.backup_contacts || state.serviceMember.currentBackupContacts || [],
    // TODO: change when we support PPM shipments as well
    mtoShipments: selectMTOShipmentsByMoveId(state, move.id),
    // TODO: change when we support multiple moves
    move,
    mtoShipment: selectMTOShipmentForMTO(state, get(move, 'id', '')),
  };
};

// in order to avoid setting up proxy server only for storybook, pass in stub function so API requests don't fail
const mergeProps = (stateProps, dispatchProps, ownProps) => ({
  ...stateProps,
  ...dispatchProps,
  ...ownProps,
});

export default withContext(connect(mapStateToProps, mergeProps)(Home));
