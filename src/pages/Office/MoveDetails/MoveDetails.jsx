import React, { useEffect, useMemo, useState } from 'react';
import { Link, useHistory, useParams } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { func } from 'prop-types';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import 'styles/office.scss';
import hasRiskOfExcess from 'utils/hasRiskOfExcess';
import { MOVES, MTO_SERVICE_ITEMS, MTO_SHIPMENTS } from 'constants/queryKeys';
import SERVICE_ITEM_STATUSES from 'constants/serviceItems';
import { shipmentStatuses } from 'constants/shipments';
import AllowancesList from 'components/Office/DefinitionLists/AllowancesList';
import CustomerInfoList from 'components/Office/DefinitionLists/CustomerInfoList';
import OrdersList from 'components/Office/DefinitionLists/OrdersList';
import DetailsPanel from 'components/Office/DetailsPanel/DetailsPanel';
import FinancialReviewButton from 'components/Office/FinancialReviewButton/FinancialReviewButton';
import FinancialReviewModal from 'components/Office/FinancialReviewModal/FinancialReviewModal';
import ApprovedRequestedShipments from 'components/Office/RequestedShipments/ApprovedRequestedShipments';
import SubmittedRequestedShipments from 'components/Office/RequestedShipments/SubmittedRequestedShipments';
import { useMoveDetailsQueries } from 'hooks/queries';
import { updateMoveStatus, updateMTOShipmentStatus, updateFinancialFlag } from 'services/ghcApi';
import LeftNav from 'components/LeftNav/LeftNav';
import LeftNavTag from 'components/LeftNavTag/LeftNavTag';
import Restricted from 'components/Restricted/Restricted';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { SIT_EXTENSION_STATUS } from 'constants/sitExtensions';
import { ORDERS_TYPE } from 'constants/orders';
import { permissionTypes } from 'constants/permissions';
import { objectIsMissingFieldWithCondition } from 'utils/displayFlags';

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
};

const MoveDetails = ({
  setUnapprovedShipmentCount,
  setUnapprovedServiceItemCount,
  setExcessWeightRiskCount,
  setUnapprovedSITExtensionCount,
}) => {
  const { moveCode } = useParams();
  const [isFinancialModalVisible, setIsFinancialModalVisible] = useState(false);
  const [shipmentMissingRequiredInformation, setShipmentMissingRequiredInformation] = useState(false);
  const [alertMessage, setAlertMessage] = useState(null);
  const [alertType, setAlertType] = useState('success');
  const history = useHistory();

  const { move, order, mtoShipments, mtoServiceItems, isLoading, isError } = useMoveDetailsQueries(moveCode);

  // for now we are only showing dest type on retiree and separatee orders
  let isRetirementOrSeparation = false;

  isRetirementOrSeparation =
    order?.order_type === ORDERS_TYPE.RETIREMENT || order?.order_type === ORDERS_TYPE.SEPARATION;

  if (isRetirementOrSeparation) {
    // destination type must be set for for HHG, NTSR shipments only
    errorIfMissing.HHG = [{ fieldName: 'destinationType' }];
    errorIfMissing.HHG_OUTOF_NTS_DOMESTIC.push({ fieldName: 'destinationType' });
    errorIfMissing.HHG_SHORTHAUL_DOMESTIC = [{ fieldName: 'destinationType' }];
    errorIfMissing.HHG_LONGHAUL_DOMESTIC = [{ fieldName: 'destinationType' }];
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

  const shipmentsInfoNonPPM = mtoShipments?.filter((shipment) => shipment.shipmentType !== 'PPM');

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
    let estimatedWeightCalc = null;
    const riskOfExcessAcknowledged = !!move?.excess_weight_acknowledged_at;

    if (mtoShipments?.some((s) => s.primeEstimatedWeight)) {
      estimatedWeightCalc = mtoShipments
        ?.filter((s) => s.primeEstimatedWeight && s.status === shipmentStatuses.APPROVED)
        .reduce((prev, current) => {
          return prev + current.primeEstimatedWeight;
        }, 0);
    }

    if (hasRiskOfExcess(estimatedWeightCalc, order?.entitlement.totalWeight) && !riskOfExcessAcknowledged) {
      setExcessWeightRiskCount(1);
    } else {
      setExcessWeightRiskCount(0);
    }
  }, [mtoShipments, setExcessWeightRiskCount, order, move]);

  useEffect(() => {
    let unapprovedSITExtensionCount = 0;
    mtoShipments?.forEach((mtoShipment) => {
      if (mtoShipment.sitExtensions?.find((sitEx) => sitEx.status === SIT_EXTENSION_STATUS.PENDING)) {
        unapprovedSITExtensionCount += 1;
      }
    });
    setUnapprovedSITExtensionCount(unapprovedSITExtensionCount);
  }, [mtoShipments, setUnapprovedSITExtensionCount]);

  useEffect(() => {
    let shipmentIsMissingInformation = false;

    mtoShipments?.forEach((mtoShipment) => {
      const fieldsToCheckForShipment = errorIfMissing[mtoShipment.shipmentType];
      const existsMissingFieldsOnShipment = fieldsToCheckForShipment?.some((field) =>
        objectIsMissingFieldWithCondition(mtoShipment, field),
      );

      // If there were no fields to check, then nothing was required.
      if (fieldsToCheckForShipment && existsMissingFieldsOnShipment) {
        shipmentIsMissingInformation = true;
      }
    });
    setShipmentMissingRequiredInformation(shipmentIsMissingInformation);
  }, [mtoShipments]);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { customer, entitlement: allowances } = order;

  if (submittedShipments.length > 0 && approvedOrCanceledShipments.length > 0) {
    sections = ['requested-shipments', 'approved-shipments', ...sections];
  } else if (approvedOrCanceledShipments.length > 0) {
    sections = ['approved-shipments', ...sections];
  } else if (submittedShipments.length > 0) {
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
  };
  const allowancesInfo = {
    branch: customer.agency,
    rank: order.grade,
    weightAllowance: allowances.totalWeight,
    authorizedWeight: allowances.authorizedWeight,
    progear: allowances.proGearWeight,
    spouseProgear: allowances.proGearWeightSpouse,
    storageInTransit: allowances.storageInTransit,
    dependents: allowances.dependentsAuthorized,
    requiredMedicalEquipmentWeight: allowances.requiredMedicalEquipmentWeight,
    organizationalClothingAndIndividualEquipment: allowances.organizationalClothingAndIndividualEquipment,
  };
  const customerInfo = {
    name: `${customer.last_name}, ${customer.first_name}`,
    dodId: customer.dodID,
    phone: `+1 ${customer.phone}`,
    email: customer.email,
    currentAddress: customer.current_address,
    backupContact: customer.backup_contact,
  };

  const requiredOrdersInfo = {
    ordersNumber: order.order_number,
    ordersType: order.order_type,
    ordersTypeDetail: order.order_type_detail,
    tacMDC: order.tac,
  };

  const hasMissingOrdersRequiredInfo = Object.values(requiredOrdersInfo).some((value) => !value || value === '');
  const hasAmendedOrders = ordersInfo.uploadedAmendedOrderID && !ordersInfo.amendedOrdersAcknowledgedAt;

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <LeftNav sections={sections}>
          <LeftNavTag
            className="usa-tag usa-tag--alert"
            associatedSectionName="orders"
            showTag={hasMissingOrdersRequiredInfo}
            testID="tag"
          >
            <FontAwesomeIcon icon="exclamation" />
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
            showTag={!shipmentMissingRequiredInformation}
            testID="requestedShipmentsTag"
          >
            {submittedShipments?.length}
          </LeftNavTag>
          <LeftNavTag
            className="usa-tag usa-tag--alert"
            associatedSectionName="requested-shipments"
            showTag={shipmentMissingRequiredInformation}
            testID="shipment-missing-info-alert"
          >
            <FontAwesomeIcon icon="exclamation" />
          </LeftNavTag>
        </LeftNav>

        <GridContainer className={styles.gridContainer} data-testid="too-move-details">
          <div className={styles.tooMoveDetailsHeadingFlexbox}>
            <h1 className={styles.tooMoveDetailsH1}>Move details</h1>
            <Restricted to={permissionTypes.updateFinancialReviewFlag}>
              <div className={styles.tooFinancialReviewContainer}>
                <FinancialReviewButton
                  onClick={handleShowFinancialReviewModal}
                  reviewRequested={move.financialReviewFlag}
                />
              </div>
            </Restricted>
          </div>
          {isFinancialModalVisible && (
            <FinancialReviewModal
              onClose={handleCancelFinancialReviewModal}
              onSubmit={handleSubmitFinancialReviewModal}
              initialRemarks={move?.financialReviewRemarks}
              initialSelection={move?.financialReviewFlag}
            />
          )}
          <Grid row className={styles.pageHeader}>
            {alertMessage && (
              <Grid col={12} className={styles.alertContainer}>
                <Alert headingLevel="h4" slim type={alertType}>
                  {alertMessage}
                </Alert>
              </Grid>
            )}
          </Grid>
          {submittedShipments.length > 0 && (
            <div className={styles.section} id="requested-shipments">
              <SubmittedRequestedShipments
                mtoShipments={submittedShipments}
                ordersInfo={ordersInfo}
                allowancesInfo={allowancesInfo}
                customerInfo={customerInfo}
                approveMTO={mutateMoveStatus}
                approveMTOShipment={mutateMTOShipmentStatus}
                moveTaskOrder={move}
                missingRequiredOrdersInfo={hasMissingOrdersRequiredInfo}
                handleAfterSuccess={history.push}
                moveCode={moveCode}
                errorIfMissing={errorIfMissing}
                displayDestinationType={isRetirementOrSeparation}
                mtoServiceItems={mtoServiceItems}
              />
            </div>
          )}
          {approvedOrCanceledShipments.length > 0 && (
            <div className={styles.section} id="approved-shipments">
              <ApprovedRequestedShipments
                mtoShipments={approvedOrCanceledShipments}
                ordersInfo={ordersInfo}
                mtoServiceItems={mtoServiceItems}
                moveCode={moveCode}
                displayDestinationType={isRetirementOrSeparation}
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
                    <Link className="usa-button usa-button--secondary" data-testid="view-orders" to="orders">
                      View orders
                    </Link>
                  }
                >
                  <Link className="usa-button usa-button--secondary" data-testid="edit-orders" to="orders">
                    Edit orders
                  </Link>
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
                    <Link className="usa-button usa-button--secondary" data-testid="view-allowances" to="allowances">
                      View allowances
                    </Link>
                  }
                >
                  <Link className="usa-button usa-button--secondary" data-testid="edit-allowances" to="allowances">
                    Edit allowances
                  </Link>
                </Restricted>
              }
              shipmentsInfoNonPpm={shipmentsInfoNonPPM}
            >
              <AllowancesList info={allowancesInfo} />
            </DetailsPanel>
          </div>
          <div className={styles.section} id="customer-info">
            <DetailsPanel title="Customer info">
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
};

export default MoveDetails;
