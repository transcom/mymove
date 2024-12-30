import React, { useEffect, useState } from 'react';
import { node, string } from 'prop-types';
import moment from 'moment';
import { connect } from 'react-redux';
import { Alert, Button } from '@trussworks/react-uswds';
import { generatePath, useNavigate, useParams, useLocation } from 'react-router-dom';

import styles from './Home.module.scss';
import {
  HelperAmendedOrders,
  HelperApprovedMove,
  HelperNeedsOrders,
  HelperNeedsShipment,
  HelperNeedsSubmitMove,
  HelperSubmittedMove,
  HelperPPMCloseoutSubmitted,
} from './HomeHelpers';

import CancelMoveConfirmationModal from 'components/ConfirmationModals/CancelMoveConfirmationModal';
import AsyncPacketDownloadLink from 'shared/AsyncPacketDownloadLink/AsyncPacketDownloadLink';
import ErrorModal from 'shared/ErrorModal/ErrorModal';
import ConnectedDestructiveShipmentConfirmationModal from 'components/ConfirmationModals/DestructiveShipmentConfirmationModal';
import Contact from 'components/Customer/Home/Contact';
import DocsUploaded from 'components/Customer/Home/DocsUploaded';
import PrintableLegalese from 'components/Customer/Home/PrintableLegalese';
import Step from 'components/Customer/Home/Step';
import SectionWrapper from 'components/Customer/SectionWrapper';
import PPMSummaryList from 'components/PPMSummaryList/PPMSummaryList';
import ShipmentList from 'components/ShipmentList/ShipmentList';
import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { profileStates } from 'constants/customerStates';
import MOVE_STATUSES from 'constants/moves';
import { customerRoutes } from 'constants/routes';
import { ppmShipmentStatuses, shipmentTypes } from 'constants/shipments';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import {
  deleteMTOShipment,
  getAllMoves,
  getMTOShipmentsForMove,
  downloadPPMAOAPacket,
  cancelMove,
} from 'services/internalApi';
import { withContext } from 'shared/AppContext';
import { SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';
import {
  getSignedCertification as getSignedCertificationAction,
  selectSignedCertification,
} from 'shared/Entities/modules/signed_certifications';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { updateMTOShipments, updateAllMoves as updateAllMovesAction } from 'store/entities/actions';
import {
  selectAllMoves,
  selectCurrentOrders,
  selectIsProfileComplete,
  selectMTOShipmentsForCurrentMove,
  selectServiceMemberFromLoggedInUser,
  selectUploadsForCurrentAmendedOrders,
  selectUploadsForCurrentOrders,
} from 'store/entities/selectors';
import { formatCustomerDate, formatUBAllowanceWeight, formatWeight } from 'utils/formatters';
import {
  isPPMAboutInfoComplete,
  isPPMShipmentComplete,
  isBoatShipmentComplete,
  isMobileHomeShipmentComplete,
  isWeightTicketComplete,
} from 'utils/shipments';
import withRouter from 'utils/routing';
import { ADVANCE_STATUSES } from 'constants/ppms';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import ToolTip from 'shared/ToolTip/ToolTip';

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

const MoveHome = ({ serviceMemberMoves, isProfileComplete, serviceMember, signedCertification, updateAllMoves }) => {
  // loading the moveId in params to select move details from serviceMemberMoves in state
  const { moveId } = useParams();
  const navigate = useNavigate();
  let { state } = useLocation();
  state = { ...state, moveId };
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [showCancelMoveModal, setShowCancelMoveModal] = useState(false);
  const [targetShipmentId, setTargetShipmentId] = useState(null);
  const [showCancelSuccessAlert, setShowCancelSuccessAlert] = useState(false);
  const [showDeleteSuccessAlert, setShowDeleteSuccessAlert] = useState(false);
  const [showDeleteErrorAlert, setShowDeleteErrorAlert] = useState(false);
  const [showErrorAlert, setShowErrorAlert] = useState(false);
  const [isManageSupportingDocsEnabled, setIsManageSupportingDocsEnabled] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      setIsManageSupportingDocsEnabled(await isBooleanFlagEnabled('manage_supporting_docs'));
    };
    fetchData();
  }, []);

  const handleCancelMove = () => {
    cancelMove(moveId)
      .then(() => {
        setShowCancelSuccessAlert(true);
      })
      .catch(() => {
        setShowDeleteErrorAlert(true);
        setShowCancelSuccessAlert(false);
      })
      .finally(() => {
        const path = generatePath('/');
        navigate(path);
      });
  };

  // fetching all move data on load since this component is dependent on that data
  // this will run each time the component is loaded/accessed
  useEffect(() => {
    getAllMoves(serviceMember.id).then((response) => {
      updateAllMoves(response);
    });
  }, [updateAllMoves, serviceMember]);

  // loading placeholder while data loads - this handles any async issues
  if (!serviceMemberMoves || !serviceMemberMoves.currentMove || !serviceMemberMoves.previousMoves) {
    return (
      <div className={styles.homeContainer}>
        <div className={`usa-prose grid-container ${styles['grid-container']}`}>
          <LoadingPlaceholder />
        </div>
      </div>
    );
  }

  // Find the move in the currentMove array
  const currentMove = serviceMemberMoves.currentMove.find((move) => move.id === moveId);
  // Find the move in the previousMoves array if not found in currentMove
  const previousMove = serviceMemberMoves.previousMoves.find((move) => move.id === moveId);
  // the move will either be in the currentMove or previousMove object
  const move = currentMove || previousMove;
  const { orders } = move;
  const uploadedOrderDocuments = orders?.uploaded_orders?.uploads || [];
  let mtoShipments;
  if (!move.mtoShipments) {
    mtoShipments = [];
  } else {
    mtoShipments = move.mtoShipments;
  }

  // checking to see if the orders object has a length
  const hasOrdersAndUpload = () => {
    return !!Object.keys(orders).length && !!uploadedOrderDocuments.length;
  };

  // checking if there are amended orders and if the move status is not approved
  const hasUnapprovedAmendedOrders = () => {
    const amendedOrders = orders?.uploaded_amended_orders || {};
    return !!Object.keys(amendedOrders).length && move.status !== 'APPROVED';
  };

  // checking if the user has order info, but no uploads
  const hasOrdersNoUpload = () => {
    return !!Object.keys(orders).length && !uploadedOrderDocuments.length;
  };

  // checking if there are any shipments in the move object
  const hasAnyShipments = () => {
    return !!Object.keys(orders).length && !!mtoShipments.length;
  };

  // checking status of the move is in any status other than DRAFT
  const hasSubmittedMove = () => {
    return !!Object.keys(move).length && move.status !== 'DRAFT';
  };

  // checking if a PPM shipment is waiting on payment approval
  const hasSubmittedPPMCloseout = () => {
    const finishedCloseout = mtoShipments.filter(
      (shipment) => shipment?.ppmShipment?.status === ppmShipmentStatuses.NEEDS_CLOSEOUT,
    );
    return !!finishedCloseout.length;
  };

  // checking if there are any incomplete shipment
  const hasIncompleteShipment = () => {
    if (!mtoShipments) return false;
    const shipmentValidators = {
      [SHIPMENT_TYPES.PPM]: isPPMShipmentComplete,
      [SHIPMENT_TYPES.BOAT_HAUL_AWAY]: isBoatShipmentComplete,
      [SHIPMENT_TYPES.BOAT_TOW_AWAY]: isBoatShipmentComplete,
      [SHIPMENT_TYPES.MOBILE_HOME]: isMobileHomeShipmentComplete,
    };

    return mtoShipments.some((shipment) => {
      const validateShipment = shipmentValidators[shipment.shipmentType];
      return validateShipment && !validateShipment(shipment);
    });
  };

  // determine if at least one advance was APPROVED (advance_status in ppm_shipments table is not nil)
  const hasAdvanceApproved = () => {
    const approvedAdvances = mtoShipments.filter(
      (shipment) =>
        shipment?.ppmShipment?.advanceStatus === ADVANCE_STATUSES.APPROVED.apiValue ||
        shipment?.ppmShipment?.advanceStatus === ADVANCE_STATUSES.EDITED.apiValue,
    );
    return !!approvedAdvances.length;
  };

  // checking if the customer has requested an advance
  const hasAdvanceRequested = () => {
    const requestedAdvances = mtoShipments.filter((shipment) => shipment?.ppmShipment?.hasRequestedAdvance);
    return !!requestedAdvances.length;
  };

  // check to see if all advance_status are REJECTED
  const hasAllAdvancesRejected = () => {
    const rejectedAdvances = mtoShipments.filter(
      (shipment) => shipment?.ppmShipment?.advanceStatus === ADVANCE_STATUSES.REJECTED.apiValue,
    );
    return !hasAdvanceApproved() && rejectedAdvances.length > 0;
  };

  // checking the move status, if approved, return true
  const isMoveApproved = () => {
    return move.status === MOVE_STATUSES.APPROVED;
  };

  // checking to see if prime is counseling this move, return true
  const isPrimeCounseled = () => {
    return !orders.providesServicesCounseling;
  };

  // checking to see if prime has completed counseling, return true
  const isPrimeCounselingComplete = () => {
    return move.primeCounselingCompletedAt?.indexOf('0001-01-01') < 0;
  };

  // check for FF and if move is submitted, can refactor once FF is removed
  // to just use hasSubmittedMove
  const isAdditionalDocumentsButtonAvailable = () => {
    return isManageSupportingDocsEnabled && hasSubmittedMove();
  };

  // logic that handles deleting a shipment
  // calls internal API and updates shipments
  const handleDeleteShipmentConfirmation = (shipmentId) => {
    deleteMTOShipment(shipmentId)
      .then(() => {
        getAllMoves(serviceMember.id).then((response) => {
          updateAllMoves(response);
        });
        getMTOShipmentsForMove(move.id).then((response) => {
          updateMTOShipments(response);
          setShowDeleteErrorAlert(false);
          setShowDeleteSuccessAlert(true);
        });
      })
      .catch(() => {
        setShowDeleteErrorAlert(true);
        setShowDeleteSuccessAlert(false);
      })
      .finally(() => {
        setShowDeleteModal(false);
      });
  };

  const shipmentActionBtnLabel = () => {
    if (hasSubmittedMove()) {
      return '';
    }
    if (hasAnyShipments()) {
      return 'Add another shipment';
    }
    return 'Set up your shipments';
  };

  const reportByLabel = () => {
    switch (orders.orders_type) {
      case 'RETIREMENT':
        return 'Retirement date';
      case 'SEPARATION':
        return 'Separation date';
      default:
        return 'Report by';
    }
  };

  const hideDeleteModal = () => {
    setShowDeleteModal(false);
  };

  const handleShipmentClick = (shipmentId, shipmentNumber, shipmentType) => {
    let queryString = '';
    if (shipmentNumber) {
      queryString = `?shipmentNumber=${shipmentNumber}`;
    }

    let destLink = '';
    if (
      shipmentType === shipmentTypes.HHG ||
      shipmentType === shipmentTypes.PPM ||
      shipmentType === shipmentTypes.BOAT ||
      shipmentType === shipmentTypes.MOBILE_HOME
    ) {
      destLink = `${generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: move.id,
        mtoShipmentId: shipmentId,
      })}${queryString}`;
    } else {
      // this will handle nts/ntsr shipments
      destLink = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
        moveId: move.id,
        mtoShipmentId: shipmentId,
      });
    }

    navigate(destLink);
  };

  const handleDeleteClick = (shipmentId) => {
    setShowDeleteModal(true);
    setTargetShipmentId(shipmentId);
  };

  const handlePPMUploadClick = (shipmentId) => {
    const shipment = mtoShipments.find((mtoShipment) => mtoShipment.id === shipmentId);

    const aboutInfoComplete = isPPMAboutInfoComplete(shipment.ppmShipment);

    let path = generatePath(customerRoutes.SHIPMENT_PPM_ABOUT_PATH, {
      moveId: move.id,
      mtoShipmentId: shipmentId,
    });

    if (aboutInfoComplete) {
      if (shipment.ppmShipment.weightTickets.length === 0) {
        path = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH, {
          moveId: move.id,
          mtoShipmentId: shipmentId,
        });
      } else if (!shipment.ppmShipment.weightTickets.some(isWeightTicketComplete)) {
        path = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
          moveId: move.id,
          mtoShipmentId: shipmentId,
          weightTicketId: shipment.ppmShipment.weightTickets[0].id,
        });
      } else {
        path = generatePath(customerRoutes.SHIPMENT_PPM_REVIEW_PATH, {
          moveId: move.id,
          mtoShipmentId: shipmentId,
        });
      }
    }

    navigate(path);
  };

  // eslint-disable-next-line class-methods-use-this
  const sortAllShipments = () => {
    const allShipments = JSON.parse(JSON.stringify(mtoShipments));
    allShipments.sort((a, b) => moment(a.createdAt) - moment(b.createdAt));

    return allShipments;
  };

  // eslint-disable-next-line class-methods-use-this
  const handlePrintLegalese = (e) => {
    e.preventDefault();
    window.print();
  };

  const handleNewPathClick = (path) => {
    navigate(path, { state });
  };

  const handlePPMFeedbackClick = (shipmentId) => {
    const path = generatePath(customerRoutes.SHIPMENT_PPM_FEEDBACK_PATH, {
      moveId: move.id,
      mtoShipmentId: shipmentId,
    });

    navigate(path);
  };

  // if the move has amended orders that aren't approved, it will display an info box at the top of the page
  const renderAlert = () => {
    if (hasUnapprovedAmendedOrders()) {
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

  // handles logic of which helper boxes to render
  const renderHelper = () => {
    if (!hasOrdersAndUpload()) return <HelperNeedsOrders />;
    if (!hasAnyShipments()) return <HelperNeedsShipment />;
    if (!hasSubmittedMove()) return <HelperNeedsSubmitMove />;
    if (hasSubmittedPPMCloseout()) return <HelperPPMCloseoutSubmitted />;
    if (hasUnapprovedAmendedOrders()) return <HelperAmendedOrders />;
    if (isMoveApproved()) return <HelperApprovedMove />;
    return <HelperSubmittedMove />;
  };

  const renderCustomerHeaderText = () => {
    return (
      <>
        <p>
          You’re moving to <strong>{orders.new_duty_location.name}</strong> from{' '}
          <strong>{orders.origin_duty_location?.name}</strong>
          <br />
          {` ${reportByLabel()} `}
          <strong>{moment(orders.report_by_date).format('DD MMM YYYY')}</strong>
        </p>

        <dl className={styles.subheaderContainer}>
          <div className={styles.subheaderSubsection}>
            <dt>Weight allowance</dt>
            <dd>{formatWeight(orders.authorizedWeight)}</dd>
          </div>
          {orders?.entitlement?.ub_allowance > 0 && (
            <div className={styles.subheaderSubsection}>
              <dt>UB allowance</dt>
              <dd>
                {formatUBAllowanceWeight(orders?.entitlement?.ub_allowance)}{' '}
                <ToolTip
                  color="#8cafea"
                  text="The weight of your UB shipment is also part of your overall authorized weight allowance."
                  data-testid="ubAllowanceToolTip"
                />
              </dd>
            </div>
          )}
          {move.moveCode && (
            <div className={styles.subheaderSubsection}>
              <dt>Move code</dt>
              <dd>#{move.moveCode}</dd>
            </div>
          )}
        </dl>
      </>
    );
  };

  const togglePPMPacketErrorModal = () => {
    setShowErrorAlert(!showErrorAlert);
  };

  const errorModalMessage =
    "Something went wrong downloading PPM paperwork. Please try again later. If that doesn't fix it, contact the ";

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

  const additionalDocumentsClick = () => {
    const uploadAdditionalDocumentsPath = generatePath(customerRoutes.UPLOAD_ADDITIONAL_DOCUMENTS_PATH, {
      moveId: move.id,
    });
    navigate(uploadAdditionalDocumentsPath, { state, moveId });
  };

  // eslint-disable-next-line camelcase
  const { current_location } = serviceMember;
  const ordersPath = hasOrdersNoUpload() ? `/orders/upload/${orders.id}` : `/orders/upload/${orders.id}`;

  const shipmentSelectionPath =
    move?.id &&
    (hasAnyShipments()
      ? generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, { moveId: move.id })
      : generatePath(customerRoutes.SHIPMENT_MOVING_INFO_PATH, { moveId: move.id }));

  const confirmationPath = move?.id && generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId: move.id });
  const profileEditPath = generatePath(customerRoutes.PROFILE_PATH);
  const ordersEditPath = `/move/${move.id}/review/edit-orders/${orders.id}`;
  const ordersAmendPath = `/orders/amend/${orders.id}`;
  const allSortedShipments = sortAllShipments(mtoShipments);
  const ppmShipments = allSortedShipments.filter((shipment) => shipment.shipmentType === SHIPMENT_OPTIONS.PPM);

  // eslint-disable-next-line camelcase
  const currentLocation = current_location;
  const shipmentNumbersByType = {};
  return (
    <>
      <CancelMoveConfirmationModal
        isOpen={showCancelMoveModal}
        moveID={moveId}
        onClose={() => setShowCancelMoveModal(false)}
        onSubmit={handleCancelMove}
      />
      <ConnectedDestructiveShipmentConfirmationModal
        isOpen={showDeleteModal}
        shipmentID={targetShipmentId}
        onClose={hideDeleteModal}
        onSubmit={handleDeleteShipmentConfirmation}
        title="Delete this?"
        content="Your information will be gone. You’ll need to start over if you want it back."
        submitText="Yes, Delete"
        closeText="No, Keep It"
      />
      <ErrorModal isOpen={showErrorAlert} closeModal={togglePPMPacketErrorModal} errorMessage={errorModalMessage} />
      <div className={styles.homeContainer}>
        <header data-testid="customer-header" className={styles['customer-header']}>
          <div className={`usa-prose grid-container ${styles['grid-container']}`}>
            <h2>
              {serviceMember.first_name} {serviceMember.last_name}
            </h2>
            {(hasOrdersNoUpload() || hasOrdersAndUpload()) && renderCustomerHeaderText()}
          </div>
        </header>
        <div className={`usa-prose grid-container ${styles['grid-container']}`}>
          {showCancelSuccessAlert && (
            <Alert headingLevel="h4" slim type="success">
              Your move was canceled.
            </Alert>
          )}
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
              {renderAlert()}
              {renderHelper()}
              {!hasSubmittedMove() && !showCancelSuccessAlert ? (
                <div className={styles.cancelMoveContainer}>
                  <Button
                    onClick={() => {
                      setShowCancelMoveModal(true);
                    }}
                    unstyled
                    data-testid="cancel-move-button"
                  >
                    Cancel move
                  </Button>
                </div>
              ) : null}
              <SectionWrapper>
                <Step
                  complete={serviceMember.is_profile_complete}
                  completedHeaderText="Profile complete"
                  editBtnLabel="Edit"
                  headerText="Profile complete"
                  step="1"
                  onEditBtnClick={() => handleNewPathClick(profileEditPath)}
                  actionBtnLabel={
                    isAdditionalDocumentsButtonAvailable() ? 'Upload/Manage Non-Orders Documentation' : null
                  }
                  onActionBtnClick={() => additionalDocumentsClick()}
                >
                  <Description>Make sure to keep your personal information up to date during your move.</Description>
                </Step>
                {!hasSubmittedMove() && (
                  <Step
                    complete={hasOrdersAndUpload()}
                    completedHeaderText="Orders uploaded"
                    editBtnLabel={hasOrdersAndUpload() ? 'Edit' : ''}
                    onEditBtnClick={() => handleNewPathClick(ordersEditPath)}
                    headerText="Upload orders"
                    actionBtnLabel={!hasOrdersAndUpload() ? 'Add orders' : ''}
                    onActionBtnClick={() => handleNewPathClick(ordersPath)}
                    step="2"
                  >
                    {hasOrdersAndUpload() && !hasSubmittedMove() ? (
                      <DocsUploaded files={uploadedOrderDocuments} />
                    ) : (
                      <Description>Upload photos of each page, or upload a PDF.</Description>
                    )}
                  </Step>
                )}
                {hasSubmittedMove() && hasOrdersAndUpload() && (
                  <Step
                    complete={hasOrdersAndUpload() && hasSubmittedMove()}
                    completedHeaderText="Orders"
                    editBtnLabel="Upload/Manage Orders Documentation"
                    onEditBtnClick={() => handleNewPathClick(ordersAmendPath)}
                    headerText="Orders"
                    step="2"
                    containerClassName="step-amended-orders"
                  >
                    <p>If you receive amended orders</p>
                    <ul>
                      <li>Upload the new document(s) here</li>
                      <li>If you have not had a counseling session talk to your local transportation office</li>
                      <li>If you have been assigned a Customer Care Representative, you can speak directly to them</li>
                      <li>They will update your move info to reflect the new orders</li>
                    </ul>
                  </Step>
                )}
                <Step
                  actionBtnLabel={shipmentActionBtnLabel()}
                  actionBtnDisabled={!hasOrdersAndUpload() || showCancelSuccessAlert}
                  actionBtnId="shipment-selection-btn"
                  onActionBtnClick={() => handleNewPathClick(shipmentSelectionPath)}
                  complete={!hasIncompleteShipment() && hasAnyShipments()}
                  completedHeaderText="Shipments"
                  headerText="Set up shipments"
                  secondaryBtn={hasAnyShipments()}
                  secondaryClassName="margin-top-2"
                  step="3"
                >
                  {hasAnyShipments() ? (
                    <div>
                      <ShipmentList
                        shipments={allSortedShipments}
                        onShipmentClick={handleShipmentClick}
                        onDeleteClick={handleDeleteClick}
                        moveSubmitted={hasSubmittedMove()}
                      />
                      {hasSubmittedMove() && (
                        <p className={styles.descriptionExtra}>
                          If you need to change, add, or cancel shipments, talk to your move counselor or Customer Care
                          Representative
                        </p>
                      )}
                    </div>
                  ) : (
                    <Description>
                      We will collect addresses, dates, and how you want to move your personal property.
                      <br /> Note: You can change these details later by talking to a move counselor or customer care
                      representative.
                    </Description>
                  )}
                </Step>
                <Step
                  actionBtnDisabled={hasIncompleteShipment() || !hasAnyShipments()}
                  actionBtnId="review-and-submit-btn"
                  actionBtnLabel={!hasSubmittedMove() ? 'Review and submit' : 'Review your request'}
                  complete={hasSubmittedMove()}
                  completedHeaderText="Move request confirmed"
                  containerClassName="margin-bottom-8"
                  headerText="Confirm move request"
                  onActionBtnClick={() => handleNewPathClick(confirmationPath)}
                  secondaryBtn={hasSubmittedMove()}
                  secondaryBtnClassName={styles.secondaryBtn}
                  step="4"
                >
                  {hasSubmittedMove() ? (
                    <Description className={styles.moveSubmittedDescription} dataTestId="move-submitted-description">
                      Move submitted {formatCustomerDate(move.submittedAt) || 'Not submitted yet'}.<br />
                      <Button unstyled onClick={handlePrintLegalese} className={styles.printBtn}>
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
                {!!ppmShipments.length && hasSubmittedMove() && hasAdvanceRequested() && (
                  <Step
                    complete={hasAdvanceApproved() || hasAllAdvancesRejected()}
                    completedHeaderText={
                      hasAllAdvancesRejected() ? 'Advance request denied' : 'Advance request reviewed'
                    }
                    headerText="Advance request submitted"
                    step="5"
                  >
                    <SectionWrapper className={styles['ppm-shipment']}>
                      {hasAdvanceApproved() && (
                        <>
                          <Description>
                            Your Advance Operating Allowance (AOA) request has been reviewed. Download the paperwork for
                            approved requests and submit it to your Finance Office to receive your advance.
                            <br />
                            <br /> The amount you receive will be deducted from your PPM incentive payment. If your
                            incentive ends up being less than your advance, you will be required to pay back the
                            difference.
                            <br />
                            <br />
                          </Description>
                          {ppmShipments.map((shipment) => {
                            const { shipmentType } = shipment;
                            if (shipmentNumbersByType[shipmentType]) {
                              shipmentNumbersByType[shipmentType] += 1;
                            } else {
                              shipmentNumbersByType[shipmentType] = 1;
                            }
                            const shipmentNumber = shipmentNumbersByType[shipmentType];
                            return (
                              <>
                                <strong>
                                  {shipmentTypes[shipment.shipmentType]}
                                  {` ${shipmentNumber} `}
                                </strong>
                                {(shipment?.ppmShipment?.advanceStatus === ADVANCE_STATUSES.APPROVED.apiValue ||
                                  shipment?.ppmShipment?.advanceStatus === ADVANCE_STATUSES.EDITED.apiValue) && (
                                  // TODO: B-18060 will add link to method that will create the AOA packet and return for download
                                  <p className={styles.downloadLink}>
                                    <AsyncPacketDownloadLink
                                      id={shipment?.ppmShipment?.id}
                                      label="Download AOA Paperwork (PDF)"
                                      asyncRetrieval={downloadPPMAOAPacket}
                                      onFailure={togglePPMPacketErrorModal}
                                    />
                                  </p>
                                )}
                                {shipment?.ppmShipment?.advanceStatus === ADVANCE_STATUSES.REJECTED.apiValue && (
                                  <Description>Advance request denied</Description>
                                )}
                                {shipment?.ppmShipment?.advanceStatus == null && (
                                  <Description>Advance request pending</Description>
                                )}
                              </>
                            );
                          })}
                        </>
                      )}
                      {hasAllAdvancesRejected() && (
                        <Description>
                          Your Advance Operating Allowance (AOA) request has been denied. You may be able to use your
                          Government Travel Charge Card (GTCC). Contact your local transportation office to verify GTCC
                          usage authorization or ask any questions.
                        </Description>
                      )}
                      {!isPrimeCounseled() && !hasAdvanceApproved() && !hasAllAdvancesRejected() && (
                        <Description>
                          Your service will review your request for an Advance Operating Allowance (AOA). If approved,
                          you will be able to download the paperwork for your request and submit it to your Finance
                          Office to receive your advance.
                          <br />
                          <br /> The amount you receive will be deducted from your PPM incentive payment. If your
                          incentive ends up being less than your advance, you will be required to pay back the
                          difference.
                        </Description>
                      )}
                      {isPrimeCounseled() && !hasAdvanceApproved() && !hasAllAdvancesRejected() && (
                        <Description>
                          Once you have received counseling for your PPM you will receive emailed instructions on how to
                          download your Advance Operating Allowance (AOA) packet. Please consult with your
                          Transportation Office for review of your AOA packet.
                          <br />
                          <br /> The amount you receive will be deducted from your PPM incentive payment. If your
                          incentive ends up being less than your advance, you will be required to pay back the
                          difference.
                          <br />
                          <br />
                        </Description>
                      )}
                      {isPrimeCounselingComplete() && (
                        <>
                          {ppmShipments.map((shipment) => {
                            const { shipmentType } = shipment;
                            if (shipmentNumbersByType[shipmentType]) {
                              shipmentNumbersByType[shipmentType] += 1;
                            } else {
                              shipmentNumbersByType[shipmentType] = 1;
                            }
                            const shipmentNumber = shipmentNumbersByType[shipmentType];
                            return (
                              <>
                                <strong>
                                  {shipmentTypes[shipment.shipmentType]}
                                  {` ${shipmentNumber} `}
                                </strong>
                                {shipment?.ppmShipment?.hasRequestedAdvance && (
                                  <p className={styles.downloadLink}>
                                    <AsyncPacketDownloadLink
                                      id={shipment?.ppmShipment?.id}
                                      label="Download AOA Paperwork (PDF)"
                                      asyncRetrieval={downloadPPMAOAPacket}
                                      onFailure={togglePPMPacketErrorModal}
                                    />
                                  </p>
                                )}
                                {!shipment?.ppmShipment?.hasRequestedAdvance && (
                                  <>
                                    <br />
                                    <br />
                                  </>
                                )}
                              </>
                            );
                          })}
                        </>
                      )}
                    </SectionWrapper>
                  </Step>
                )}
                {!!ppmShipments.length && hasSubmittedMove() && (
                  <Step
                    headerText="Manage your PPM"
                    completedHeaderText="Manage your PPM"
                    step={hasAdvanceRequested() ? '6' : '5'}
                  >
                    <PPMSummaryList
                      shipments={ppmShipments}
                      onUploadClick={handlePPMUploadClick}
                      onDownloadError={togglePPMPacketErrorModal}
                      onFeedbackClick={handlePPMFeedbackClick}
                    />
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
};

MoveHome.defaultProps = {
  orders: {},
  serviceMember: null,
  signedCertification: {},
  uploadedAmendedOrderDocuments: [],
  router: {},
};

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberMoves = selectAllMoves(state);

  return {
    isProfileComplete: selectIsProfileComplete(state),
    orders: selectCurrentOrders(state) || {},
    uploadedOrderDocuments: selectUploadsForCurrentOrders(state),
    uploadedAmendedOrderDocuments: selectUploadsForCurrentAmendedOrders(state),
    serviceMember,
    serviceMemberMoves,
    backupContacts: serviceMember?.backup_contacts || [],
    signedCertification: selectSignedCertification(state),
    mtoShipments: selectMTOShipmentsForCurrentMove(state),
  };
};

const mapDispatchToProps = {
  getSignedCertification: getSignedCertificationAction,
  updateShipmentList: updateMTOShipments,
  updateAllMoves: updateAllMovesAction,
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
    )(requireCustomerState(MoveHome, profileStates.BACKUP_CONTACTS_COMPLETE)),
  ),
);
