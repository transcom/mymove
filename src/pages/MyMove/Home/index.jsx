import React, { Component } from 'react';
import { arrayOf, bool, func, node, shape, string } from 'prop-types';
import moment from 'moment';
import { connect } from 'react-redux';
import { Alert, Button } from '@trussworks/react-uswds';
import { generatePath } from 'react-router';
import { withRouter } from 'react-router-dom';

import styles from './Home.module.scss';
import {
  HelperAmendedOrders,
  HelperApprovedMove,
  HelperNeedsOrders,
  HelperNeedsShipment,
  HelperNeedsSubmitMove,
  HelperSubmittedMove,
} from './HomeHelpers';

import ConnectedDestructiveShipmentConfirmationModal from 'components/ConfirmationModals/DestructiveShipmentConfirmationModal';
import Contact from 'components/Customer/Home/Contact';
import DocsUploaded from 'components/Customer/Home/DocsUploaded';
import PrintableLegalese from 'components/Customer/Home/PrintableLegalese';
import Step from 'components/Customer/Home/Step';
import SectionWrapper from 'components/Customer/SectionWrapper';
import PPMSummaryList from 'components/PPMSummaryList/PPMSummaryList';
import ScrollToTop from 'components/ScrollToTop';
import ShipmentList from 'components/ShipmentList/ShipmentList';
import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { profileStates } from 'constants/customerStates';
import MOVE_STATUSES from 'constants/moves';
import { customerRoutes } from 'constants/routes';
import { shipmentTypes } from 'constants/shipments';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import { deleteMTOShipment, getMTOShipmentsForMove } from 'services/internalApi';
import { withContext } from 'shared/AppContext';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import {
  getSignedCertification as getSignedCertificationAction,
  selectSignedCertification,
} from 'shared/Entities/modules/signed_certifications';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { updateMTOShipments } from 'store/entities/actions';
import {
  selectCurrentMove,
  selectCurrentOrders,
  selectIsProfileComplete,
  selectMTOShipmentsForCurrentMove,
  selectServiceMemberFromLoggedInUser,
  selectUploadsForCurrentAmendedOrders,
  selectUploadsForCurrentOrders,
} from 'store/entities/selectors';
import { HistoryShape, MoveShape, OrdersShape, UploadShape } from 'types/customerShapes';
import { ShipmentShape } from 'types/shipment';
import { formatCustomerDate, formatWeight } from 'utils/formatters';
import { isPPMShipmentComplete } from 'utils/shipments';

const Description = ({ className, children, dataTestId }) => (
  <p className={`${styles.description} ${className}`} data-testid={dataTestId}>
    {children}
  </p>
);

Description.propTypes = {
  className: string,
  children: node.isRequired,
  dataTestId: string,
};

Description.defaultProps = {
  className: '',
  dataTestId: '',
};

export class Home extends Component {
  constructor(props) {
    super(props);
    this.state = {
      showDeleteModal: false,
      targetShipmentId: null,
      showDeleteSuccessAlert: false,
      showDeleteErrorAlert: false,
    };
  }

  componentDidMount() {
    const { move, getSignedCertification } = this.props;
    if (Object.entries(move).length && move.status === MOVE_STATUSES.SUBMITTED) {
      getSignedCertification(move.id);
    }
  }

  componentDidUpdate(prevProps) {
    const { move, getSignedCertification } = this.props;

    if (!Object.entries(prevProps.move).length && Object.entries(move).length) {
      getSignedCertification(move.id);
    }
  }

  get hasOrders() {
    const { orders, uploadedOrderDocuments } = this.props;
    return !!Object.keys(orders).length && !!uploadedOrderDocuments.length;
  }

  get hasUnapprovedAmendedOrders() {
    const { move, uploadedAmendedOrderDocuments } = this.props;
    return !!uploadedAmendedOrderDocuments?.length && move.status !== 'APPROVED';
  }

  get hasOrdersNoUpload() {
    const { orders, uploadedOrderDocuments } = this.props;
    return !!Object.keys(orders).length && !uploadedOrderDocuments.length;
  }

  get hasAnyShipments() {
    const { mtoShipments } = this.props;
    return this.hasOrders && !!mtoShipments.length;
  }

  get hasSubmittedMove() {
    const { move } = this.props;
    return !!Object.keys(move).length && move.status !== 'DRAFT';
  }

  get hasPPMShipments() {
    const { mtoShipments } = this.props;
    return mtoShipments?.some((shipment) => shipment.ppmShipment);
  }

  get hasAllCompletedPPMShipments() {
    const { mtoShipments } = this.props;
    return mtoShipments?.filter((s) => s.shipmentType === SHIPMENT_OPTIONS.PPM)?.every((s) => isPPMShipmentComplete(s));
  }

  get isMoveApproved() {
    const { move } = this.props;
    return move.status === MOVE_STATUSES.APPROVED;
  }

  get shipmentActionBtnLabel() {
    if (this.hasSubmittedMove) {
      return '';
    }
    if (this.hasAnyShipments) {
      return 'Add another shipment';
    }
    return 'Set up your shipments';
  }

  get reportByLabel() {
    const { orders } = this.props;
    switch (orders.orders_type) {
      case 'RETIREMENT':
        return 'Retirement date';
      case 'SEPARATION':
        return 'Separation date';
      default:
        return 'Report by';
    }
  }

  renderAlert = () => {
    if (this.hasUnapprovedAmendedOrders) {
      return (
        <Alert headingLevel="h4" type="success" slim data-testid="unapproved-amended-orders-alert">
          <span className={styles.alertMessageFirstLine}>
            The transportation office will review your new documents and update your move info. Contact your movers to
            coordinate any changes to your move.
          </span>
          <span className={styles.alertMessageSecondLine}>You don&apos;t need to do anything else in MilMove.</span>
        </Alert>
      );
    }
    return null;
  };

  renderHelper = () => {
    if (!this.hasOrders) return <HelperNeedsOrders />;
    if (!this.hasAnyShipments) return <HelperNeedsShipment />;
    if (!this.hasSubmittedMove) return <HelperNeedsSubmitMove />;
    if (this.hasUnapprovedAmendedOrders) return <HelperAmendedOrders />;
    if (this.isMoveApproved) return <HelperApprovedMove />;
    return <HelperSubmittedMove />;
  };

  renderCustomerHeaderText = () => {
    const { serviceMember, orders, move } = this.props;
    return (
      <>
        <p>
          You’re moving to <strong>{orders.new_duty_location.name}</strong> from{' '}
          <strong>{orders.origin_duty_location?.name}.</strong>
          {` ${this.reportByLabel} `}
          <strong>{moment(orders.report_by_date).format('DD MMM YYYY')}.</strong>
        </p>

        <dl className={styles.subheaderContainer}>
          <div className={styles.subheaderSubsection}>
            <dt>Weight allowance</dt>
            <dd>
              {orders.has_dependents
                ? formatWeight(serviceMember.weight_allotment.total_weight_self_plus_dependents)
                : formatWeight(serviceMember.weight_allotment.total_weight_self)}
              .
            </dd>
          </div>
          {move.locator && (
            <div className={styles.subheaderSubsection}>
              <dt>Move code</dt>
              <dd>#{move.locator}</dd>
            </div>
          )}
        </dl>
      </>
    );
  };

  hideDeleteModal = () => {
    this.setState({
      showDeleteModal: false,
    });
  };

  handleShipmentClick = (shipmentId, shipmentNumber, shipmentType) => {
    const { move, history } = this.props;
    let queryString = '';
    if (shipmentNumber) {
      queryString = `?shipmentNumber=${shipmentNumber}`;
    }

    let destLink = '';
    if (shipmentType === shipmentTypes.HHG || shipmentType === shipmentTypes.PPM) {
      destLink = `${generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: move.id,
        mtoShipmentId: shipmentId,
      })}${queryString}`;
    } else {
      // nts/ntsr shipment
      destLink = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: move.id,
        mtoShipmentId: shipmentId,
      });
    }

    history.push(destLink);
  };

  handleDeleteClick = (shipmentId) => {
    this.setState({
      showDeleteModal: true,
      targetShipmentId: shipmentId,
    });
  };

  handleDeleteShipmentConfirmation = (shipmentId) => {
    const { move, updateShipmentList } = this.props;
    deleteMTOShipment(shipmentId)
      .then(() => {
        getMTOShipmentsForMove(move.id).then((response) => {
          updateShipmentList(response);
          this.setState({
            showDeleteSuccessAlert: true,
            showDeleteErrorAlert: false,
          });
        });
      })
      .catch(() => {
        this.setState({
          showDeleteErrorAlert: true,
          showDeleteSuccessAlert: false,
        });
      })
      .finally(() => {
        this.setState({ showDeleteModal: false });
      });
  };

  handleNewPathClick = (path) => {
    const { history } = this.props;
    history.push(path);
  };

  handlePPMUploadClick = (shipmentId) => {
    const { move, history } = this.props;
    const path = generatePath(customerRoutes.SHIPMENT_PPM_ABOUT_PATH, {
      moveId: move.id,
      mtoShipmentId: shipmentId,
    });
    history.push(path);
  };

  // eslint-disable-next-line class-methods-use-this
  sortAllShipments = (mtoShipments) => {
    const allShipments = JSON.parse(JSON.stringify(mtoShipments));
    allShipments.sort((a, b) => moment(a.createdAt) - moment(b.createdAt));

    return allShipments;
  };

  // eslint-disable-next-line class-methods-use-this
  handlePrintLegalese = (e) => {
    e.preventDefault();
    window.print();
  };

  render() {
    const { isProfileComplete, move, mtoShipments, serviceMember, signedCertification, uploadedOrderDocuments } =
      this.props;

    const { showDeleteModal, targetShipmentId, showDeleteSuccessAlert, showDeleteErrorAlert } = this.state;

    // early return if loading user/service member
    if (!serviceMember) {
      return (
        <div className={styles.homeContainer}>
          <div className={`usa-prose grid-container ${styles['grid-container']}`}>
            <LoadingPlaceholder />
          </div>
        </div>
      );
    }

    // eslint-disable-next-line camelcase
    const { current_location } = serviceMember;
    const ordersPath = this.hasOrdersNoUpload ? customerRoutes.ORDERS_UPLOAD_PATH : customerRoutes.ORDERS_INFO_PATH;

    const shipmentSelectionPath =
      move?.id &&
      (this.hasAnyShipments
        ? generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, { moveId: move.id })
        : generatePath(customerRoutes.SHIPMENT_MOVING_INFO_PATH, { moveId: move.id }));

    const confirmationPath = move?.id && generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId: move.id });
    const profileEditPath = customerRoutes.PROFILE_PATH;
    const ordersEditPath = `/moves/${move.id}/review/edit-orders`;
    const ordersAmendPath = customerRoutes.ORDERS_AMEND_PATH;
    const allSortedShipments = this.sortAllShipments(mtoShipments);
    const ppmShipments = allSortedShipments.filter((shipment) => shipment.shipmentType === SHIPMENT_OPTIONS.PPM);

    // eslint-disable-next-line camelcase
    const currentLocation = current_location;

    return (
      <>
        <ScrollToTop />
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
        <div className={styles.homeContainer}>
          <header data-testid="customer-header" className={styles['customer-header']}>
            <div className={`usa-prose grid-container ${styles['grid-container']}`}>
              <h2>
                {serviceMember.first_name} {serviceMember.last_name}
              </h2>
              {(this.hasOrdersNoUpload || this.hasOrders) && this.renderCustomerHeaderText()}
            </div>
          </header>
          <div className={`usa-prose grid-container ${styles['grid-container']}`}>
            {showDeleteSuccessAlert && (
              <Alert headingLevel="h4" slim type="success">
                The shipment was deleted.
              </Alert>
            )}
            {showDeleteErrorAlert && (
              <Alert headingLevel="h4" slim type="error">
                Something went wrong, and your changes were not saved. Please try again later or contact your counselor.
              </Alert>
            )}
            <ConnectedFlashMessage />

            {isProfileComplete && (
              <>
                {this.renderAlert()}
                {this.renderHelper()}
                <SectionWrapper>
                  <Step
                    complete={serviceMember.is_profile_complete}
                    completedHeaderText="Profile complete"
                    editBtnLabel="Edit"
                    headerText="Profile complete"
                    step="1"
                    onEditBtnClick={() => this.handleNewPathClick(profileEditPath)}
                  >
                    <Description>Make sure to keep your personal information up to date during your move.</Description>
                  </Step>
                  {!this.hasSubmittedMove && (
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
                      {this.hasOrders && !this.hasSubmittedMove ? (
                        <DocsUploaded files={uploadedOrderDocuments} />
                      ) : (
                        <Description>Upload photos of each page, or upload a PDF.</Description>
                      )}
                    </Step>
                  )}
                  {this.hasSubmittedMove && this.hasOrders && (
                    <Step
                      complete={this.hasOrders && this.hasSubmittedMove}
                      completedHeaderText="Orders"
                      editBtnLabel="Upload documents"
                      onEditBtnClick={() => this.handleNewPathClick(ordersAmendPath)}
                      headerText="Orders"
                      step="2"
                      containerClassName="step-amended-orders"
                    >
                      <p>If you receive amended orders:</p>
                      <ul>
                        <li>Upload the new documents here</li>
                        <li>Talk directly with your movers about changes</li>
                        <li>The transportation office will update your move info to reflect the new orders</li>
                      </ul>
                    </Step>
                  )}
                  <Step
                    actionBtnLabel={this.shipmentActionBtnLabel}
                    actionBtnDisabled={!this.hasOrders}
                    actionBtnId="shipment-selection-btn"
                    onActionBtnClick={() => this.handleNewPathClick(shipmentSelectionPath)}
                    complete={this.hasPPMShipments ? this.hasAllCompletedPPMShipments : this.hasAnyShipments}
                    completedHeaderText="Shipments"
                    headerText="Set up shipments"
                    secondaryBtn={this.hasAnyShipments}
                    secondaryClassName="margin-top-2"
                    step="3"
                  >
                    {this.hasAnyShipments ? (
                      <div>
                        <ShipmentList
                          shipments={allSortedShipments}
                          onShipmentClick={this.handleShipmentClick}
                          onDeleteClick={this.handleDeleteClick}
                          moveSubmitted={this.hasSubmittedMove}
                        />
                        {this.hasSubmittedMove && (
                          <p className={styles.descriptionExtra}>
                            If you need to change, add, or get rid of shipments, talk to your move counselor or to your
                            movers.
                          </p>
                        )}
                      </div>
                    ) : (
                      <Description>
                        We&apos;ll collect addresses, dates, and how you want to move your things.
                        <br />
                        Note: You can change these details later by talking to a move counselor or your movers.
                      </Description>
                    )}
                  </Step>
                  <Step
                    actionBtnDisabled={this.hasPPMShipments ? !this.hasAllCompletedPPMShipments : !this.hasAnyShipments}
                    actionBtnId="review-and-submit-btn"
                    actionBtnLabel={!this.hasSubmittedMove ? 'Review and submit' : 'Review your request'}
                    complete={this.hasSubmittedMove}
                    completedHeaderText="Move request confirmed"
                    containerClassName="margin-bottom-8"
                    headerText="Confirm move request"
                    onActionBtnClick={() => this.handleNewPathClick(confirmationPath)}
                    secondaryBtn={this.hasSubmittedMove}
                    secondaryBtnClassName={styles.secondaryBtn}
                    step="4"
                  >
                    {this.hasSubmittedMove ? (
                      <Description className={styles.moveSubmittedDescription} dataTestId="move-submitted-description">
                        Move submitted {formatCustomerDate(move.submitted_at)}.<br />
                        <Button unstyled onClick={this.handlePrintLegalese} className={styles.printBtn}>
                          Print the legal agreement
                        </Button>
                      </Description>
                    ) : (
                      <Description>
                        Review your move details and sign the legal paperwork, then send the info on to your move
                        counselor.
                      </Description>
                    )}
                  </Step>
                  {!!ppmShipments.length && this.hasSubmittedMove && (
                    <Step headerText="Manage your PPM" completedHeaderText="Manage your PPM" step="5">
                      <PPMSummaryList shipments={ppmShipments} onUploadClick={this.handlePPMUploadClick} />
                    </Step>
                  )}
                </SectionWrapper>
                <Contact
                  header="Contacts"
                  dutyLocationName={currentLocation?.transportation_office?.name}
                  officeType="Origin Transportation Office"
                  telephone={currentLocation?.transportation_office?.phone_lines[0]}
                />
              </>
            )}
          </div>
        </div>
        <PrintableLegalese signature={signedCertification.signature} signatureDate={signedCertification.created_at} />
      </>
    );
  }
}

Home.propTypes = {
  orders: OrdersShape,
  serviceMember: shape({
    first_name: string,
    last_name: string,
  }),
  mtoShipments: arrayOf(ShipmentShape).isRequired,
  uploadedOrderDocuments: arrayOf(UploadShape).isRequired,
  uploadedAmendedOrderDocuments: arrayOf(UploadShape),
  history: HistoryShape.isRequired,
  move: MoveShape.isRequired,
  isProfileComplete: bool.isRequired,
  signedCertification: shape({
    signature: string,
    created_at: string,
  }),
  getSignedCertification: func.isRequired,
  updateShipmentList: func.isRequired,
};

Home.defaultProps = {
  orders: {},
  serviceMember: null,
  signedCertification: {},
  uploadedAmendedOrderDocuments: [],
};

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const move = selectCurrentMove(state) || {};

  return {
    isProfileComplete: selectIsProfileComplete(state),
    orders: selectCurrentOrders(state) || {},
    uploadedOrderDocuments: selectUploadsForCurrentOrders(state),
    uploadedAmendedOrderDocuments: selectUploadsForCurrentAmendedOrders(state),
    serviceMember,
    backupContacts: serviceMember?.backup_contacts || [],
    signedCertification: selectSignedCertification(state),
    mtoShipments: selectMTOShipmentsForCurrentMove(state),
    // TODO: change when we support multiple moves
    move,
  };
};

const mapDispatchToProps = {
  getSignedCertification: getSignedCertificationAction,
  updateShipmentList: updateMTOShipments,
};

// in order to avoid setting up proxy server only for storybook, pass in stub function so API requests don't fail
const mergeProps = (stateProps, dispatchProps, ownProps) => ({
  ...stateProps,
  ...dispatchProps,
  ...ownProps,
});

export default withContext(
  withRouter(
    connect(
      mapStateToProps,
      mapDispatchToProps,
      mergeProps,
    )(requireCustomerState(Home, profileStates.BACKUP_CONTACTS_COMPLETE)),
  ),
);
