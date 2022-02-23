import React, { useEffect, useMemo, useState } from 'react';
import { Link, useHistory, useParams } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { queryCache, useMutation } from 'react-query';
import { func } from 'prop-types';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import 'styles/office.scss';
import hasRiskOfExcess from 'utils/hasRiskOfExcess';
import { MOVES, MTO_SERVICE_ITEMS, MTO_SHIPMENTS } from 'constants/queryKeys';
import SERVICE_ITEM_STATUSES from 'constants/serviceItems';
import { shipmentStatuses } from 'constants/shipments';
import LeftNav from 'components/LeftNav/LeftNav';
import AllowancesList from 'components/Office/DefinitionLists/AllowancesList';
import CustomerInfoList from 'components/Office/DefinitionLists/CustomerInfoList';
import OrdersList from 'components/Office/DefinitionLists/OrdersList';
import DetailsPanel from 'components/Office/DetailsPanel/DetailsPanel';
import FinancialReviewModal from 'components/Office/FinancialReviewModal/FinancialReviewModal';
import FinancialReviewButton from 'components/Office/FinancialReviewButton/FinancialReviewButton';
import RequestedShipments from 'components/Office/RequestedShipments/RequestedShipments';
import { useMoveDetailsQueries } from 'hooks/queries';
import { updateMoveStatus, updateMTOShipmentStatus, updateFinancialFlag } from 'services/ghcApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { SIT_EXTENSION_STATUS } from 'constants/sitExtensions';
import LeftNavSection from 'components/LeftNavSection/LeftNavSection';
import LeftNavTag from 'components/LeftNavTag/LeftNavTag';

const sectionLabels = {
  'requested-shipments': 'Requested shipments',
  'approved-shipments': 'Approved shipments',
  orders: 'Orders',
  allowances: 'Allowances',
  'customer-info': 'Customer info',
};

const errorIfMissing = {
  HHG_OUTOF_NTS_DOMESTIC: ['ntsRecordedWeight', 'serviceOrderNumber', 'tacType'],
  HHG_INTO_NTS_DOMESTIC: ['tacType'],
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

  const [activeSection, setActiveSection] = useState('');

  const { move, order, mtoShipments, mtoServiceItems, isLoading, isError } = useMoveDetailsQueries(moveCode);

  let sections = useMemo(() => {
    return ['orders', 'allowances', 'customer-info'];
  }, []);

  // use mutation calls
  const [mutateMoveStatus] = useMutation(updateMoveStatus, {
    onSuccess: (data) => {
      queryCache.setQueryData([MOVES, data.locator], data);
      queryCache.invalidateQueries([MOVES, data.locator]);
      queryCache.invalidateQueries([MTO_SERVICE_ITEMS, data.id]);
    },
  });

  const [mutateMTOShipmentStatus] = useMutation(updateMTOShipmentStatus, {
    onSuccess: (updatedMTOShipment) => {
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryCache.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryCache.invalidateQueries([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID]);
      queryCache.invalidateQueries([MTO_SERVICE_ITEMS, updatedMTOShipment.moveTaskOrderID]);
    },
  });

  const [mutateFinancialReview] = useMutation(updateFinancialFlag, {
    onSuccess: (data) => {
      queryCache.setQueryData([MOVES, data.locator], data);
      queryCache.invalidateQueries([MOVES, data.locator]);
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
      const existsMissingFieldsOnShipment = fieldsToCheckForShipment?.some(
        (field) => !mtoShipment[field] || mtoShipment[field] === '',
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
    newDutyStation: order.destinationDutyStation,
    currentDutyStation: order.originDutyStation,
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
        <LeftNav className={styles.sidebar}>
          {sections.map((s) => {
            return (
              <LeftNavSection
                key={`sidenav_${s}`}
                sectionName={s}
                isActive={s === activeSection}
                onClickHandler={() => setActiveSection(s)}
              >
                {sectionLabels[`${s}`]}
                <LeftNavTag
                  className="usa-tag usa-tag--alert"
                  showTag={s === 'orders' && hasMissingOrdersRequiredInfo}
                  testID="tag"
                >
                  <FontAwesomeIcon icon="exclamation" />
                </LeftNavTag>
                <LeftNavTag
                  className={styles.tag}
                  showTag={Boolean(s === 'orders' && !hasMissingOrdersRequiredInfo && hasAmendedOrders)}
                  testID="newOrdersNavTag"
                >
                  NEW
                </LeftNavTag>
                <LeftNavTag
                  className={styles.tag}
                  showTag={s === 'requested-shipments' && !shipmentMissingRequiredInformation}
                  testID="requestedShipmentsTag"
                >
                  {submittedShipments?.length}
                </LeftNavTag>
                <LeftNavTag
                  className="usa-tag usa-tag--alert"
                  showTag={s === 'requested-shipments' && shipmentMissingRequiredInformation}
                  testID="shipment-missing-info-alert"
                >
                  <FontAwesomeIcon icon="exclamation" />
                </LeftNavTag>
              </LeftNavSection>
            );
          })}
        </LeftNav>

        <GridContainer className={styles.gridContainer} data-testid="too-move-details">
          <div className={styles.tooMoveDetailsHeadingFlexbox}>
            <h1 className={styles.tooMoveDetailsH1}>Move details</h1>
            <div className={styles.tooFinancialReviewContainer}>
              <FinancialReviewButton
                onClick={handleShowFinancialReviewModal}
                reviewRequested={move.financialReviewFlag}
              />
            </div>
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
                <Alert slim type={alertType}>
                  {alertMessage}
                </Alert>
              </Grid>
            )}
          </Grid>
          {/* TODO - RequestedShipments could be simplified, if extra time we could tackle this or just write a story to track */}
          {submittedShipments.length > 0 && (
            <div className={styles.section} id="requested-shipments">
              <RequestedShipments
                mtoShipments={submittedShipments}
                ordersInfo={ordersInfo}
                allowancesInfo={allowancesInfo}
                customerInfo={customerInfo}
                mtoServiceItems={mtoServiceItems}
                shipmentsStatus={shipmentStatuses.SUBMITTED}
                approveMTO={mutateMoveStatus}
                approveMTOShipment={mutateMTOShipmentStatus}
                moveTaskOrder={move}
                missingRequiredOrdersInfo={hasMissingOrdersRequiredInfo}
                handleAfterSuccess={history.push}
                moveCode={moveCode}
              />
            </div>
          )}
          {approvedOrCanceledShipments.length > 0 && (
            <div className={styles.section} id="approved-shipments">
              <RequestedShipments
                moveTaskOrder={move}
                mtoShipments={approvedOrCanceledShipments}
                ordersInfo={ordersInfo}
                allowancesInfo={allowancesInfo}
                customerInfo={customerInfo}
                mtoServiceItems={mtoServiceItems}
                shipmentsStatus={shipmentStatuses.APPROVED}
                moveCode={moveCode}
              />
            </div>
          )}
          <div className={styles.section} id="orders">
            <DetailsPanel
              title="Orders"
              tag={hasAmendedOrders ? 'NEW' : ''}
              editButton={
                <Link className="usa-button usa-button--secondary" data-testid="edit-orders" to="orders">
                  Edit orders
                </Link>
              }
            >
              <OrdersList ordersInfo={ordersInfo} showMissingWarnings />
            </DetailsPanel>
          </div>
          <div className={styles.section} id="allowances">
            <DetailsPanel
              title="Allowances"
              editButton={
                <Link className="usa-button usa-button--secondary" data-testid="edit-allowances" to="allowances">
                  Edit allowances
                </Link>
              }
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
