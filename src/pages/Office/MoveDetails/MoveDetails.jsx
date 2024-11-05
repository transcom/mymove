import React, { useEffect, useMemo, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { Alert, Grid, GridContainer, Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { func } from 'prop-types';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import NotificationScrollToTop from 'components/NotificationScrollToTop';
import hasRiskOfExcess from 'utils/hasRiskOfExcess';
import { MOVES, MTO_SERVICE_ITEMS, MTO_SHIPMENTS } from 'constants/queryKeys';
import { tooRoutes } from 'constants/routes';
import SERVICE_ITEM_STATUSES from 'constants/serviceItems';
import { ADDRESS_UPDATE_STATUS, shipmentStatuses, ppmShipmentStatuses } from 'constants/shipments';
import AllowancesList from 'components/Office/DefinitionLists/AllowancesList';
import CustomerInfoList from 'components/Office/DefinitionLists/CustomerInfoList';
import OrdersList from 'components/Office/DefinitionLists/OrdersList';
import DetailsPanel from 'components/Office/DetailsPanel/DetailsPanel';
import FinancialReviewButton from 'components/Office/FinancialReviewButton/FinancialReviewButton';
import FinancialReviewModal from 'components/Office/FinancialReviewModal/FinancialReviewModal';
import CancelMoveConfirmationModal from 'components/ConfirmationModals/CancelMoveConfirmationModal';
import ApprovedRequestedShipments from 'components/Office/RequestedShipments/ApprovedRequestedShipments';
import SubmittedRequestedShipments from 'components/Office/RequestedShipments/SubmittedRequestedShipments';
import { useMoveDetailsQueries } from 'hooks/queries';
import { updateMoveStatus, updateMTOShipmentStatus, cancelMove, updateFinancialFlag } from 'services/ghcApi';
import LeftNav from 'components/LeftNav/LeftNav';
import LeftNavTag from 'components/LeftNavTag/LeftNavTag';
import Restricted from 'components/Restricted/Restricted';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { SIT_EXTENSION_STATUS } from 'constants/sitExtensions';
import { ORDERS_TYPE } from 'constants/orders';
import { permissionTypes } from 'constants/permissions';
import { objectIsMissingFieldWithCondition } from 'utils/displayFlags';
import formattedCustomerName from 'utils/formattedCustomerName';
import { calculateEstimatedWeight } from 'hooks/custom';
import { ADVANCE_STATUSES } from 'constants/ppms';
import { MOVE_STATUSES, technicalHelpDeskURL } from 'shared/constants';

const MoveDetails = ({
  setUnapprovedShipmentCount,
  setUnapprovedServiceItemCount,
  setExcessWeightRiskCount,
  setUnapprovedSITExtensionCount,
  setShipmentErrorConcernCount,
  shipmentErrorConcernCount,
  setShipmentsWithDeliveryAddressUpdateRequestedCount,
  missingOrdersInfoCount,
  setMissingOrdersInfoCount,
  isMoveLocked,
}) => {
  const { moveCode } = useParams();
  const [isFinancialModalVisible, setIsFinancialModalVisible] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);
  const [isCancelMoveModalVisible, setIsCancelMoveModalVisible] = useState(false);
  const [alertMessage, setAlertMessage] = useState(null);
  const [alertType, setAlertType] = useState('success');

  // RA Summary: eslint-disable-next-line react-hooks/exhaustive-deps
  // RA: This rule is used to enforce correct dependency arrays in hooks like useEffect, useCallback, and useMemo.
  // RA: We are disabling this rule here because adding useMemo causes undesired behavior in our case.
  // RA Developer Status: Known Issue - Intentional decision to prevent page refresh issues related to action counts.
  // RA Validator Status: CODEOWNER ACCEPTED
  // RA Modified Severity: N/A
  // eslint-disable-next-line react-hooks/exhaustive-deps
  const errorIfMissing = {
    HHG_INTO_NTS_DOMESTIC: [
      { fieldName: 'storageFacility' },
      { fieldName: 'serviceOrderNumber' },
      { fieldName: 'tacType' },
    ],
    HHG_OUTOF_NTS_DOMESTIC: [
      { fieldName: 'storageFacility' },
      { fieldName: 'ntsRecordedWeight' },
      { fieldName: 'serviceOrderNumber' },
      { fieldName: 'tacType' },
    ],
    PPM: [
      {
        fieldName: 'advanceStatus',
        condition: (mtoShipment) =>
          mtoShipment?.ppmShipment?.hasRequestedAdvance === true &&
          mtoShipment?.ppmShipment?.advanceStatus !== ADVANCE_STATUSES.APPROVED.apiValue,
      },
    ],
  };

  const navigate = useNavigate();

  const { move, customerData, order, closeoutOffice, mtoShipments, mtoServiceItems, isLoading, isError } =
    useMoveDetailsQueries(moveCode);

  // for now we are only showing dest type on retiree and separatee orders
  let isRetirementOrSeparation = false;
  let numberOfShipmentsNotAllowedForCancel = 0;

  isRetirementOrSeparation =
    order?.order_type === ORDERS_TYPE.RETIREMENT || order?.order_type === ORDERS_TYPE.SEPARATION;

  if (isRetirementOrSeparation) {
    // destination type must be set for for HHG, NTSR shipments only
    errorIfMissing.HHG = [{ fieldName: 'destinationType' }];
    errorIfMissing.HHG_OUTOF_NTS_DOMESTIC.push({ fieldName: 'destinationType' });
  }

  let sections = useMemo(() => {
    return ['orders', 'allowances', 'customer-info'];
  }, []);

  // use mutation calls
  const queryClient = useQueryClient();
  const { mutate: mutateMoveStatus } = useMutation(updateMoveStatus, {
    onSuccess: (data) => {
      queryClient.setQueryData([MOVES, data.locator], data);
      queryClient.invalidateQueries([MOVES, data.locator]);
      queryClient.invalidateQueries([MTO_SERVICE_ITEMS, data.id]);
    },
  });

  const { mutate: mutateMTOShipmentStatus } = useMutation(updateMTOShipmentStatus, {
    onSuccess: (updatedMTOShipment) => {
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryClient.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryClient.invalidateQueries([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID]);
      queryClient.invalidateQueries([MTO_SERVICE_ITEMS, updatedMTOShipment.moveTaskOrderID]);
    },
  });

  const { mutate: mutateCancelMove } = useMutation(cancelMove, {
    onSuccess: (data) => {
      queryClient.setQueryData([MOVES, data.locator], data);
      queryClient.invalidateQueries([MOVES, data.locator]);
      queryClient.invalidateQueries({ queryKey: [MTO_SHIPMENTS] });
      setAlertMessage('Move canceled.');
      setAlertType('success');
    },
    onError: () => {
      setAlertMessage('There was a problem cancelling the move. Please try again later.');
      setAlertType('error');
    },
  });

  const { mutate: mutateFinancialReview } = useMutation(updateFinancialFlag, {
    onSuccess: (data) => {
      queryClient.setQueryData([MOVES, data.locator], data);
      queryClient.invalidateQueries([MOVES, data.locator]);
      if (data.financialReviewFlag) {
        setAlertMessage('Move flagged for financial review.');
        setAlertType('success');
      } else {
        setAlertMessage('Move unflagged for financial review.');
        setAlertType('success');
      }
    },
    onError: () => {
      setAlertMessage('There was a problem flagging the move for financial review. Please try again later.');
      setAlertType('error');
    },
  });

  const handleShowFinancialReviewModal = () => {
    setIsFinancialModalVisible(true);
  };

  const handleSubmitFinancialReviewModal = (remarks, flagForReview) => {
    // if it's set to yes let's send a true to the backend. If not we'll send false.
    const flagForReviewBool = flagForReview === 'yes';
    mutateFinancialReview({
      moveID: move.id,
      ifMatchETag: move.eTag,
      body: { remarks, flagForReview: flagForReviewBool },
    });
    setIsFinancialModalVisible(false);
  };

  const handleCancelFinancialReviewModal = () => {
    setIsFinancialModalVisible(false);
  };

  const handleShowCancelMoveModal = () => {
    setIsCancelMoveModalVisible(true);
  };

  const handleCancelMove = () => {
    mutateCancelMove({
      moveID: move.id,
      ifMatchETag: move.eTag,
    });
    setIsCancelMoveModalVisible(false);
  };

  const handleCloseCancelMoveModal = () => {
    setIsCancelMoveModalVisible(false);
  };

  const submittedShipments = mtoShipments?.filter(
    (shipment) => shipment.status === shipmentStatuses.SUBMITTED && !shipment.deletedAt,
  );

  const approvedOrCanceledShipments = mtoShipments?.filter(
    (shipment) =>
      shipment.status === shipmentStatuses.APPROVED ||
      shipment.status === shipmentStatuses.DIVERSION_REQUESTED ||
      shipment.status === shipmentStatuses.CANCELLATION_REQUESTED ||
      shipment.status === shipmentStatuses.CANCELED,
  );

  const shipmentWithDestinationAddressChangeRequest = mtoShipments?.filter(
    (shipment) => shipment?.deliveryAddressUpdate?.status === ADDRESS_UPDATE_STATUS.REQUESTED && !shipment.deletedAt,
  );

  useEffect(() => {
    const shipmentCount = shipmentWithDestinationAddressChangeRequest?.length || 0;
    if (setShipmentsWithDeliveryAddressUpdateRequestedCount)
      setShipmentsWithDeliveryAddressUpdateRequestedCount(shipmentCount);
  }, [shipmentWithDestinationAddressChangeRequest?.length, setShipmentsWithDeliveryAddressUpdateRequestedCount]);

  const shipmentsInfoNonPPM = mtoShipments?.filter((shipment) => shipment.shipmentType !== 'PPM');
  if (mtoShipments) {
    const nonDeletedShipments = mtoShipments?.filter((shipment) => !shipment.deletedAt);
    const nonPpmShipments = nonDeletedShipments.filter((shipment) => shipment.shipmentType !== 'PPM');
    const nonPpmApprovedShipments = nonPpmShipments.filter(
      (shipment) => shipment?.status === shipmentStatuses.APPROVED,
    );
    const onlyPpmShipments = nonDeletedShipments.filter((shipment) => shipment.shipmentType === 'PPM');
    const ppmCloseoutCompleteShipments = onlyPpmShipments.filter(
      (shipment) => shipment.ppmShipment?.status === ppmShipmentStatuses.CLOSEOUT_COMPLETE,
    );
    numberOfShipmentsNotAllowedForCancel = nonPpmApprovedShipments.length + ppmCloseoutCompleteShipments.length;
  }

  useEffect(() => {
    const shipmentCount = submittedShipments?.length || 0;
    setUnapprovedShipmentCount(shipmentCount);
  }, [mtoShipments, submittedShipments, setUnapprovedShipmentCount]);

  useEffect(() => {
    let serviceItemCount = 0;
    mtoServiceItems?.forEach((serviceItem) => {
      if (
        serviceItem.status === SERVICE_ITEM_STATUSES.SUBMITTED &&
        serviceItem.mtoShipmentID &&
        approvedOrCanceledShipments?.find((shipment) => shipment.id === serviceItem.mtoShipmentID)
      ) {
        serviceItemCount += 1;
      }
    });
    setUnapprovedServiceItemCount(serviceItemCount);
  }, [approvedOrCanceledShipments, mtoServiceItems, setUnapprovedServiceItemCount]);

  useEffect(() => {
    const estimatedWeight = calculateEstimatedWeight(mtoShipments);
    const riskOfExcessAcknowledged = !!move?.excess_weight_acknowledged_at;

    if (hasRiskOfExcess(estimatedWeight, order?.entitlement.totalWeight) && !riskOfExcessAcknowledged) {
      setExcessWeightRiskCount(1);
    } else {
      setExcessWeightRiskCount(0);
    }
  }, [move?.excess_weight_acknowledged_at, mtoShipments, order?.entitlement.totalWeight, setExcessWeightRiskCount]);

  useEffect(() => {
    const checkShipmentsForUnapprovedSITExtensions = (shipmentsWithStatus) => {
      let unapprovedSITExtensionCount = 0;
      shipmentsWithStatus?.forEach((mtoShipment) => {
        const unapprovedSITExtItems =
          mtoShipment.sitExtensions?.filter((sitEx) => sitEx.status === SIT_EXTENSION_STATUS.PENDING) ?? [];
        const unapprovedSITCount = unapprovedSITExtItems.length;
        unapprovedSITExtensionCount += unapprovedSITCount; // Top bar Label
      });
      return { count: unapprovedSITExtensionCount };
    };
    const { count } = checkShipmentsForUnapprovedSITExtensions(mtoShipments);
    setUnapprovedSITExtensionCount(count);
  }, [mtoShipments, setUnapprovedSITExtensionCount]);

  useEffect(() => {
    // Reset the error count before running any logic
    let numberOfErrorIfMissingForAllShipments = 0;

    // Process each shipment to accumulate errors
    mtoShipments?.forEach((mtoShipment) => {
      const errorIfMissingList = errorIfMissing[mtoShipment.shipmentType];

      if (errorIfMissingList) {
        errorIfMissingList.forEach((fieldToCheck) => {
          if (objectIsMissingFieldWithCondition(mtoShipment, fieldToCheck)) {
            numberOfErrorIfMissingForAllShipments += 1;
          }
        });
      }
    });

    // Set the error concern count after processing
    setShipmentErrorConcernCount(numberOfErrorIfMissingForAllShipments);
  }, [errorIfMissing, mtoShipments, setShipmentErrorConcernCount]);

  // using useMemo here due to this being used in a useEffect
  // using useMemo prevents the useEffect from being rendered on ever render by memoizing the object
  // so that it only recognizes the change when the orders object changes
  const requiredOrdersInfo = useMemo(
    () => ({
      ordersNumber: order?.order_number || '',
      ordersType: order?.order_type || '',
      ordersTypeDetail: order?.order_type_detail || '',
      tacMDC: order?.tac || '',
      departmentIndicator: order?.department_indicator || '',
    }),
    [order],
  );

  // Keep num of missing orders info synced up
  useEffect(() => {
    const ordersInfoCount = Object.values(requiredOrdersInfo).reduce((count, value) => {
      return !value ? count + 1 : count;
    }, 0);
    setMissingOrdersInfoCount(ordersInfoCount);
  }, [order, requiredOrdersInfo, setMissingOrdersInfoCount]);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { customer, entitlement: allowances } = order;

  if (submittedShipments?.length > 0 && approvedOrCanceledShipments?.length > 0) {
    sections = ['requested-shipments', 'approved-shipments', ...sections];
  } else if (approvedOrCanceledShipments?.length > 0) {
    sections = ['approved-shipments', ...sections];
  } else if (submittedShipments?.length > 0) {
    sections = ['requested-shipments', ...sections];
  }

  const ordersInfo = {
    newDutyLocation: order.destinationDutyLocation,
    currentDutyLocation: order.originDutyLocation,
    issuedDate: order.date_issued,
    reportByDate: order.report_by_date,
    departmentIndicator: order.department_indicator,
    ordersNumber: order.order_number,
    ordersType: order.order_type,
    ordersTypeDetail: order.order_type_detail,
    uploadedAmendedOrderID: order.uploadedAmendedOrderID,
    amendedOrdersAcknowledgedAt: order.amendedOrdersAcknowledgedAt,
    tacMDC: order.tac,
    sacSDN: order.sac,
    NTStac: order.ntsTac,
    NTSsac: order.ntsSac,
    payGrade: order.grade,
  };
  const allowancesInfo = {
    branch: customer.agency,
    grade: order.grade,
    totalWeight: allowances.totalWeight,
    progear: allowances.proGearWeight,
    spouseProgear: allowances.proGearWeightSpouse,
    storageInTransit: allowances.storageInTransit,
    dependents: allowances.dependentsAuthorized,
    requiredMedicalEquipmentWeight: allowances.requiredMedicalEquipmentWeight,
    organizationalClothingAndIndividualEquipment: allowances.organizationalClothingAndIndividualEquipment,
    gunSafe: allowances.gunSafe,
  };

  const customerInfo = {
    name: formattedCustomerName(customer.last_name, customer.first_name, customer.suffix, customer.middle_name),
    agency: customer.agency,
    edipi: customer.edipi,
    emplid: customer.emplid,
    phone: customer.phone,
    altPhone: customer.secondaryTelephone,
    email: customer.email,
    currentAddress: customer.current_address,
    backupAddress: customerData.backupAddress,
    backupContact: customer.backup_contact,
  };

  const hasMissingOrdersRequiredInfo = Object.values(requiredOrdersInfo).some((value) => !value || value === '');
  const hasAmendedOrders = ordersInfo.uploadedAmendedOrderID && !ordersInfo.amendedOrdersAcknowledgedAt;
  const hasDestinationAddressUpdate =
    shipmentWithDestinationAddressChangeRequest && shipmentWithDestinationAddressChangeRequest.length > 0;
  const tooCanCancelMove = move.status !== MOVE_STATUSES.CANCELED && numberOfShipmentsNotAllowedForCancel === 0;

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <LeftNav sections={sections}>
          <LeftNavTag
            background="#e34b11"
            associatedSectionName="orders"
            showTag={missingOrdersInfoCount !== 0}
            testID="tag"
          >
            {missingOrdersInfoCount}
          </LeftNavTag>
          <LeftNavTag
            associatedSectionName="orders"
            showTag={Boolean(!hasMissingOrdersRequiredInfo && hasAmendedOrders)}
            testID="newOrdersNavTag"
          >
            NEW
          </LeftNavTag>
          <LeftNavTag
            associatedSectionName="requested-shipments"
            showTag={submittedShipments?.length > 0}
            testID="requestedShipmentsTag"
          >
            {submittedShipments?.length || 0}
          </LeftNavTag>
          <LeftNavTag
            background="#e34b11"
            associatedSectionName="requested-shipments"
            showTag={shipmentErrorConcernCount !== 0}
            testID="shipment-missing-info-alert"
          >
            {shipmentErrorConcernCount}
          </LeftNavTag>
          <LeftNavTag
            associatedSectionName="approved-shipments"
            className="usa-tag usa-tag--alert"
            showTag={hasDestinationAddressUpdate}
          >
            <FontAwesomeIcon icon="exclamation" />
          </LeftNavTag>
        </LeftNav>

        <GridContainer className={styles.gridContainer} data-testid="too-move-details">
          <div>
            <Grid row className={styles.pageHeader}>
              {alertMessage && (
                <Grid col={12} className={styles.alertContainer}>
                  <Alert headingLevel="h4" slim type={alertType}>
                    {alertMessage}
                  </Alert>
                </Grid>
              )}
            </Grid>
            <Grid col={12} className={styles.tooMoveDetailsHeadingFlexbox}>
              <h1 className={styles.tooMoveDetailsH1}>Move details</h1>
              <Restricted to={permissionTypes.updateFinancialReviewFlag}>
                <div>
                  <FinancialReviewButton
                    onClick={handleShowFinancialReviewModal}
                    reviewRequested={move.financialReviewFlag}
                    isMoveLocked={isMoveLocked}
                  />
                </div>
              </Restricted>
            </Grid>
          </div>
          {isFinancialModalVisible && (
            <FinancialReviewModal
              onClose={handleCancelFinancialReviewModal}
              onSubmit={handleSubmitFinancialReviewModal}
              initialRemarks={move?.financialReviewRemarks}
              initialSelection={move?.financialReviewFlag}
            />
          )}
          <NotificationScrollToTop dependency={errorMessage} />
          {errorMessage && (
            <Alert data-testid="errorMessage" type="error" headingLevel="h4" heading="An error occurred">
              <p>
                {errorMessage} Please try again later, or contact the&nbsp;
                <Link to={technicalHelpDeskURL} target="_blank" rel="noreferrer">
                  Technical Help Desk
                </Link>
                .
              </p>
            </Alert>
          )}
          <Grid row col={12}>
            <Restricted to={permissionTypes.cancelMoveFlag}>
              <div className={styles.tooCancelMoveContainer}>
                {tooCanCancelMove && !isMoveLocked && (
                  <Button type="button" unstyled onClick={handleShowCancelMoveModal}>
                    Cancel move
                  </Button>
                )}
              </div>
            </Restricted>
          </Grid>
          <CancelMoveConfirmationModal
            isOpen={isCancelMoveModalVisible}
            onClose={handleCloseCancelMoveModal}
            onSubmit={handleCancelMove}
          />
          {submittedShipments?.length > 0 && (
            <div className={styles.section} id="requested-shipments">
              <SubmittedRequestedShipments
                mtoShipments={submittedShipments}
                closeoutOffice={closeoutOffice}
                ordersInfo={ordersInfo}
                allowancesInfo={allowancesInfo}
                customerInfo={customerInfo}
                approveMTO={mutateMoveStatus}
                approveMTOShipment={mutateMTOShipmentStatus}
                moveTaskOrder={move}
                missingRequiredOrdersInfo={hasMissingOrdersRequiredInfo}
                handleAfterSuccess={navigate}
                moveCode={moveCode}
                errorIfMissing={errorIfMissing}
                displayDestinationType={isRetirementOrSeparation}
                mtoServiceItems={mtoServiceItems}
                isMoveLocked={isMoveLocked}
                setErrorMessage={setErrorMessage}
              />
            </div>
          )}
          {approvedOrCanceledShipments?.length > 0 && (
            <div className={styles.section} id="approved-shipments">
              <ApprovedRequestedShipments
                mtoShipments={approvedOrCanceledShipments}
                closeoutOffice={closeoutOffice}
                ordersInfo={ordersInfo}
                mtoServiceItems={mtoServiceItems}
                moveCode={moveCode}
                displayDestinationType={isRetirementOrSeparation}
                isMoveLocked={isMoveLocked}
              />
            </div>
          )}
          <div className={styles.section} id="orders">
            <DetailsPanel
              title="Orders"
              tag={hasAmendedOrders ? 'NEW' : ''}
              editButton={
                <Restricted
                  to={permissionTypes.updateOrders}
                  fallback={
                    <Link className="usa-button usa-button--secondary" data-testid="view-orders" to="../orders">
                      View orders
                    </Link>
                  }
                >
                  {!isMoveLocked && (
                    <Link className="usa-button usa-button--secondary" data-testid="edit-orders" to="../orders">
                      Edit orders
                    </Link>
                  )}
                </Restricted>
              }
              shipmentsInfoNonPpm={shipmentsInfoNonPPM}
            >
              <OrdersList ordersInfo={ordersInfo} />
            </DetailsPanel>
          </div>
          <div className={styles.section} id="allowances">
            <DetailsPanel
              title="Allowances"
              editButton={
                <Restricted
                  to={permissionTypes.updateAllowances}
                  fallback={
                    <Link className="usa-button usa-button--secondary" data-testid="view-allowances" to="../allowances">
                      View allowances
                    </Link>
                  }
                >
                  {!isMoveLocked && (
                    <Link className="usa-button usa-button--secondary" data-testid="edit-allowances" to="../allowances">
                      Edit allowances
                    </Link>
                  )}
                </Restricted>
              }
              shipmentsInfoNonPpm={shipmentsInfoNonPPM}
            >
              <AllowancesList info={allowancesInfo} />
            </DetailsPanel>
          </div>
          <div className={styles.section} id="customer-info">
            <DetailsPanel
              title="Customer info"
              editButton={
                <Restricted to={permissionTypes.updateCustomer}>
                  {!isMoveLocked && (
                    <Link
                      className="usa-button usa-button--secondary"
                      data-testid="edit-customer-info"
                      to={`../${tooRoutes.CUSTOMER_INFO_EDIT_PATH}`}
                    >
                      Edit customer info
                    </Link>
                  )}
                </Restricted>
              }
            >
              <CustomerInfoList customerInfo={customerInfo} />
            </DetailsPanel>
          </div>
        </GridContainer>
      </div>
    </div>
  );
};

MoveDetails.propTypes = {
  setUnapprovedShipmentCount: func.isRequired,
  setUnapprovedServiceItemCount: func.isRequired,
  setExcessWeightRiskCount: func.isRequired,
  setUnapprovedSITExtensionCount: func.isRequired,
  setShipmentErrorConcernCount: func.isRequired,
  setShipmentsWithDeliveryAddressUpdateRequestedCount: func,
};

MoveDetails.defaultProps = {
  setShipmentsWithDeliveryAddressUpdateRequestedCount: () => {},
};

export default MoveDetails;
