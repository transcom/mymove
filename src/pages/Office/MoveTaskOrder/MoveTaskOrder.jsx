import React, { useEffect, useMemo, useState } from 'react';
import { Link, useParams, generatePath } from 'react-router-dom';
import { Alert, Button, Grid, GridContainer, Tag } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { connect } from 'react-redux';
import { func } from 'prop-types';
import classnames from 'classnames';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import moveTaskOrderStyles from './MoveTaskOrder.module.scss';

import ConnectedEditMaxBillableWeightModal from 'components/Office/EditMaxBillableWeightModal/EditMaxBillableWeightModal';
import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';
import { formatStorageFacilityForAPI, formatAddressForAPI, removeEtag } from 'utils/formatMtoShipment';
import hasRiskOfExcess from 'utils/hasRiskOfExcess';
import customerContactTypes from 'constants/customerContactTypes';
import dimensionTypes from 'constants/dimensionTypes';
import { MTO_SERVICE_ITEMS, MOVES, MTO_SHIPMENTS, ORDERS } from 'constants/queryKeys';
import SERVICE_ITEM_STATUSES from 'constants/serviceItems';
import { mtoShipmentTypes, shipmentStatuses } from 'constants/shipments';
import FlashGridContainer from 'containers/FlashGridContainer/FlashGridContainer';
import { shipmentSectionLabels } from 'content/shipments';
import RejectServiceItemModal from 'components/Office/RejectServiceItemModal/RejectServiceItemModal';
import RequestedServiceItemsTable from 'components/Office/RequestedServiceItemsTable/RequestedServiceItemsTable';
import RequestShipmentCancellationModal from 'components/Office/RequestShipmentCancellationModal/RequestShipmentCancellationModal';
import RequestReweighModal from 'components/Office/RequestReweighModal/RequestReweighModal';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import ShipmentHeading from 'components/Office/ShipmentHeading/ShipmentHeading';
import ShipmentDetails from 'components/Office/ShipmentDetails/ShipmentDetails';
import { useMoveTaskOrderQueries } from 'hooks/queries';
import {
  acknowledgeExcessWeightRisk,
  patchMTOServiceItemStatus,
  updateBillableWeight,
  updateMTOShipmentRequestReweigh,
  updateMTOShipmentStatus,
  updateMTOShipment,
  approveSITExtension,
  denySITExtension,
  submitSITExtension,
  updateFinancialFlag,
} from 'services/ghcApi';
import { MOVE_STATUSES } from 'shared/constants';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { setFlashMessage } from 'store/flash/actions';
import WeightDisplay from 'components/Office/WeightDisplay/WeightDisplay';
import { calculateEstimatedWeight, calculateWeightRequested } from 'hooks/custom';
import { SIT_EXTENSION_STATUS } from 'constants/sitExtensions';
import FinancialReviewButton from 'components/Office/FinancialReviewButton/FinancialReviewButton';
import FinancialReviewModal from 'components/Office/FinancialReviewModal/FinancialReviewModal';
import leftNavStyles from 'components/LeftNav/LeftNav.module.scss';
import LeftNavSection from 'components/LeftNavSection/LeftNavSection';
import LeftNavTag from 'components/LeftNavTag/LeftNavTag';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import { tooRoutes } from 'constants/routes';

const nonShipmentSectionLabels = {
  'move-weights': 'Move weights',
};

function formatShipmentDate(shipmentDateString) {
  if (shipmentDateString == null) {
    return '';
  }
  const dateObj = new Date(shipmentDateString);
  const weekday = new Intl.DateTimeFormat('en', { weekday: 'long' }).format(dateObj);
  const year = new Intl.DateTimeFormat('en', { year: 'numeric' }).format(dateObj);
  const month = new Intl.DateTimeFormat('en', { month: 'short' }).format(dateObj);
  const day = new Intl.DateTimeFormat('en', { day: '2-digit' }).format(dateObj);
  return `${weekday}, ${day} ${month} ${year}`;
}

function showShipmentFilter(shipment) {
  return (
    shipment.status === shipmentStatuses.APPROVED ||
    shipment.status === shipmentStatuses.CANCELLATION_REQUESTED ||
    shipment.status === shipmentStatuses.DIVERSION_REQUESTED ||
    shipment.status === shipmentStatuses.CANCELED
  );
}

export const MoveTaskOrder = (props) => {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [isCancelModalVisible, setIsCancelModalVisible] = useState(false);
  const [isReweighModalVisible, setIsReweighModalVisible] = useState(false);
  const [isWeightModalVisible, setIsWeightModalVisible] = useState(false);
  const [isWeightAlertVisible, setIsWeightAlertVisible] = useState(false);
  const [isSuccessAlertVisible, setIsSuccessAlertVisible] = useState(false);
  const [isFinancialModalVisible, setIsFinancialModalVisible] = useState(false);
  const [alertMessage, setAlertMessage] = useState(null);
  const [alertType, setAlertType] = useState('success');

  const [selectedShipment, setSelectedShipment] = useState(undefined);
  const [selectedServiceItem, setSelectedServiceItem] = useState(undefined);
  const [sections, setSections] = useState([]);
  const [activeSection, setActiveSection] = useState('');
  const [unapprovedServiceItemsForShipment, setUnapprovedServiceItemsForShipment] = useState({});
  const [unapprovedSITExtensionForShipment, setUnApprovedSITExtensionForShipment] = useState({});
  const [estimatedWeightTotal, setEstimatedWeightTotal] = useState(null);
  const [externalVendorShipmentCount, setExternalVendorShipmentCount] = useState(0);

  const nonShipmentSections = useMemo(() => {
    return ['move-weights'];
  }, []);

  const { moveCode } = useParams();
  const {
    setUnapprovedShipmentCount,
    setUnapprovedServiceItemCount,
    setExcessWeightRiskCount,
    setMessage,
    setUnapprovedSITExtensionCount,
  } = props;

  const { orders = {}, move, mtoShipments, mtoServiceItems, isLoading, isError } = useMoveTaskOrderQueries(moveCode);
  const order = Object.values(orders)?.[0];

  const shipmentServiceItems = useMemo(() => {
    const serviceItemsForShipment = {};
    mtoServiceItems?.forEach((item) => {
      // We're not interested in basic service items
      if (!item.mtoShipmentID) {
        return;
      }
      const newItem = { ...item };
      newItem.code = item.reServiceCode;
      newItem.serviceItem = item.reServiceName;
      newItem.details = {
        pickupPostalCode: item.pickupPostalCode,
        SITPostalCode: item.SITPostalCode,
        reason: item.reason,
        description: item.description,
        itemDimensions: item.dimensions?.find((dimension) => dimension?.type === dimensionTypes.ITEM),
        crateDimensions: item.dimensions?.find((dimension) => dimension?.type === dimensionTypes.CRATE),
        firstCustomerContact: item.customerContacts?.find((contact) => contact?.type === customerContactTypes.FIRST),
        secondCustomerContact: item.customerContacts?.find((contact) => contact?.type === customerContactTypes.SECOND),
        estimatedWeight: item.estimatedWeight,
        rejectionReason: item.rejectionReason,
      };

      if (serviceItemsForShipment[`${newItem.mtoShipmentID}`]) {
        serviceItemsForShipment[`${newItem.mtoShipmentID}`].push(newItem);
      } else {
        serviceItemsForShipment[`${newItem.mtoShipmentID}`] = [newItem];
      }
    });
    return serviceItemsForShipment;
  }, [mtoServiceItems]);

  const queryClient = useQueryClient();
  const { mutate: mutateMTOServiceItemStatus } = useMutation({
    mutationFn: patchMTOServiceItemStatus,
    onSuccess: (data, variables) => {
      const newMTOServiceItem = data.mtoServiceItems[variables.mtoServiceItemID];
      mtoServiceItems[mtoServiceItems.find((serviceItem) => serviceItem.id === newMTOServiceItem.id)] =
        newMTOServiceItem;
      queryClient.setQueryData([MTO_SERVICE_ITEMS, variables.moveId, false], mtoServiceItems);
      queryClient.invalidateQueries({ queryKey: [MTO_SERVICE_ITEMS, variables.moveId] });
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const { mutate: mutateMTOShipment } = useMutation({
    mutationFn: updateMTOShipment,
    onSuccess: (_, variables) => {
      queryClient.setQueryData([MTO_SHIPMENTS, variables.moveTaskOrderID, false], mtoShipments);
      queryClient.invalidateQueries({ queryKey: [MTO_SHIPMENTS, variables.moveTaskOrderID] });
    },
  });

  const { mutate: mutateMTOShipmentStatus } = useMutation({
    mutationFn: updateMTOShipmentStatus,
    onSuccess: (data, variables) => {
      const updatedMTOShipment = data.mtoShipments[variables.shipmentID];
      // Update mtoShipments with our updated status and set query data to match
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryClient.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      // InvalidateQuery tells other components using this data that they need to re-fetch
      // This allows the requestCancellation button to update immediately
      queryClient.invalidateQueries({ queryKey: [MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID] });
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const { mutate: mutateMTOShipmentRequestReweigh } = useMutation({
    mutationFn: updateMTOShipmentRequestReweigh,
    onSuccess: (data) => {
      // Update mtoShipments with our updated status and set query data to match
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === data.shipmentID)] = data;
      queryClient.setQueryData([MTO_SHIPMENTS, move.id, false], mtoShipments);
      // InvalidateQuery tells other components using this data that they need to re-fetch
      // This allows the requestReweigh button to update immediately
      queryClient.invalidateQueries({ queryKey: [MTO_SHIPMENTS, move.id] });
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const { mutate: mutateOrderBillableWeight } = useMutation({
    mutationFn: updateBillableWeight,
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: [MOVES, move.locator] });
      const updatedOrder = data.orders[variables.orderID];
      queryClient.setQueryData([ORDERS, variables.orderID], {
        orders: {
          [`${variables.orderID}`]: updatedOrder,
        },
      });
      queryClient.invalidateQueries({ queryKey: [ORDERS, variables.orderID] });
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const { mutate: mutateAcknowledgeExcessWeightRisk } = useMutation({
    mutationFn: acknowledgeExcessWeightRisk,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [MOVES, move.locator] });
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      // TODO: Handle error some how
      // RA Summary: eslint: no-console - System Information Leak: External
      // RA: The linter flags any use of console.
      // RA: This console displays an error message from unsuccessful mutation.
      // RA: TODO: As indicated, this error needs to be handled and needs further investigation and work.
      // RA: POAM story here: https://dp3.atlassian.net/browse/MB-5597
      // RA Developer Status: Known Issue
      // RA Validator Status: Known Issue
      // RA Modified Severity: CAT II
      // eslint-disable-next-line no-console
      console.log(errorMsg);
    },
  });

  const { mutate: mutateSITExtensionApproval } = useMutation({
    mutationFn: approveSITExtension,
    onSuccess: (data, variables) => {
      const updatedMTOShipment = data.mtoShipments[variables.shipmentID];
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryClient.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryClient.invalidateQueries({ queryKey: [MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID] });
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const { mutate: mutateSITExtensionDenial } = useMutation({
    mutationFn: denySITExtension,
    onSuccess: (data, variables) => {
      const updatedMTOShipment = data.mtoShipments[variables.shipmentID];
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryClient.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryClient.invalidateQueries({ queryKey: [MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID] });
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const { mutate: mutateSubmitSITExtension } = useMutation({
    mutationFn: submitSITExtension,
    onSuccess: (data, variables) => {
      const updatedMTOShipment = data.mtoShipments[variables.shipmentID];
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryClient.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryClient.invalidateQueries({ queryKey: [MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID] });
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const { mutate: mutateFinancialReview } = useMutation({
    mutationFn: updateFinancialFlag,
    onSuccess: (data) => {
      queryClient.setQueryData([MOVES, data.locator], data);
      queryClient.invalidateQueries({ queryKey: [MOVES, data.locator] });
    },
  });

  const handleSubmitFinancialReviewModal = (remarks, flagForReview) => {
    // if it's set to yes let's send a true to the backend. If not we'll send false.
    const flagForReviewBool = flagForReview === 'yes';
    mutateFinancialReview(
      {
        moveID: move.id,
        ifMatchETag: move.eTag,
        body: { remarks, flagForReview: flagForReviewBool },
      },
      {
        onSuccess: (data) => {
          if (data.financialReviewFlag) {
            setAlertMessage('Move flagged for financial review.');
          } else {
            setAlertMessage('Move unflagged for financial review.');
          }
          setAlertType('success');
          setIsFinancialModalVisible(false);
        },
        onError: () => {
          setAlertMessage('There was a problem flagging the move for financial review. Please try again later.');
          setAlertType('error');
        },
      },
    );
  };

  const handleCancelFinancialReviewModal = () => {
    setIsFinancialModalVisible(false);
  };

  const handleShowFinancialReviewModal = () => {
    setIsFinancialModalVisible(true);
  };

  const handleReviewSITExtension = (sitExtensionID, formValues, shipment) => {
    if (formValues.acceptExtension === 'yes') {
      mutateSITExtensionApproval({
        shipmentID: shipment.id,
        sitExtensionID,
        ifMatchETag: shipment.eTag,
        body: {
          requestReason: formValues.requestReason,
          officeRemarks: formValues.officeRemarks,
          approvedDays: parseInt(formValues.daysApproved, 10) - shipment.sitDaysAllowance,
        },
      });
    } else if (formValues.acceptExtension === 'no') {
      mutateSITExtensionDenial({
        shipmentID: shipment.id,
        sitExtensionID,
        ifMatchETag: shipment.eTag,
        body: { officeRemarks: formValues.officeRemarks },
      });
    }
  };

  const handleSubmitSITExtension = (formValues, shipment) => {
    mutateSubmitSITExtension(
      {
        shipmentID: shipment.id,
        ifMatchETag: shipment.eTag,
        body: {
          requestReason: formValues.requestReason,
          officeRemarks: formValues.officeRemarks,
          approvedDays: parseInt(formValues.daysApproved, 10) - shipment.sitDaysAllowance,
        },
      },
      {
        onSuccess: () => setIsSuccessAlertVisible(true),
      },
    );
  };

  const handleDivertShipment = (mtoShipmentID, eTag) => {
    mutateMTOShipmentStatus(
      {
        shipmentID: mtoShipmentID,
        operationPath: 'shipment.requestShipmentDiversion',
        ifMatchETag: eTag,
        onSuccessFlashMsg: `Diversion successfully requested for Shipment #${mtoShipmentID}`,
      },
      {
        onSuccess: (data, variables) => {
          setIsCancelModalVisible(false);
          // Must set FlashMesage after hiding the modal, since FlashMessage will disappear when focus changes
          setMessage(`MSG_CANCEL_SUCCESS_${variables.shipmentID}`, 'success', variables.onSuccessFlashMsg, '', true);
        },
      },
    );
  };

  const handleReweighShipment = (mtoShipmentID, eTag) => {
    mutateMTOShipmentRequestReweigh(
      {
        shipmentID: mtoShipmentID,
        ifMatchETag: eTag,
        onSuccessFlashMsg: `Reweigh successfully requested.`,
      },
      {
        onSuccess: (data, variables) => {
          setIsReweighModalVisible(false);
          // Must set FlashMesage after hiding the modal, since FlashMessage will disappear when focus changes
          setMessage(`MSG_REWEIGH_SUCCESS_${variables.shipmentID}`, 'success', variables.onSuccessFlashMsg, '', true);
        },
      },
    );
  };

  const handleEditAccountingCodes = (fields, shipment) => {
    const body = { tacType: null, sacType: null, ...fields };
    mutateMTOShipment({
      moveTaskOrderID: shipment.moveTaskOrderID,
      shipmentID: shipment.id,
      ifMatchETag: shipment.eTag,
      body,
    });
  };

  const handleUpdateMTOShipmentStatus = (moveTaskOrderID, mtoShipmentID, eTag) => {
    mutateMTOShipmentStatus(
      {
        shipmentID: mtoShipmentID,
        operationPath: 'shipment.requestShipmentCancellation',
        ifMatchETag: eTag,
        onSuccessFlashMsg: 'The request to cancel that shipment has been sent to the movers.',
      },
      {
        onSuccess: (data, variables) => {
          setIsCancelModalVisible(false);
          // Must set FlashMesage after hiding the modal, since FlashMessage will disappear when focus changes
          setMessage(`MSG_CANCEL_SUCCESS_${variables.shipmentID}`, 'success', variables.onSuccessFlashMsg, '', true);
        },
      },
    );
  };

  const handleEditFacilityInfo = (fields, shipment) => {
    const formattedStorageFacility = formatStorageFacilityForAPI(fields.storageFacility);
    const formattedStorageFacilityAddress = removeEtag(formatAddressForAPI(fields.storageFacility.address));
    const body = {
      storageFacility: { ...formattedStorageFacility, address: formattedStorageFacilityAddress },
      serviceOrderNumber: fields.serviceOrderNumber,
    };
    mutateMTOShipment({
      moveTaskOrderID: shipment.moveTaskOrderID,
      shipmentID: shipment.id,
      ifMatchETag: shipment.eTag,
      body,
    });
  };

  const handleEditServiceOrderNumber = (fields, shipment) => {
    mutateMTOShipment({
      moveTaskOrderID: shipment.moveTaskOrderID,
      shipmentID: shipment.id,
      ifMatchETag: shipment.eTag,
      body: { serviceOrderNumber: fields.serviceOrderNumber },
    });
  };

  const handleUpdateMTOServiceItemStatus = (mtoServiceItemID, mtoShipmentID, status, rejectionReason) => {
    const mtoServiceItemForRequest = shipmentServiceItems[`${mtoShipmentID}`]?.find((s) => s.id === mtoServiceItemID);

    mutateMTOServiceItemStatus(
      {
        moveId: move.id,
        mtoServiceItemID,
        status,
        rejectionReason,
        ifMatchEtag: mtoServiceItemForRequest.eTag,
      },
      {
        onSuccess: () => {
          setIsModalVisible(false);
          setSelectedServiceItem({});
        },
      },
    );
  };

  const handleUpdateBillableWeight = (maxBillableWeight) => {
    mutateOrderBillableWeight(
      {
        orderID: order.id,
        ifMatchETag: order.eTag,
        body: { authorizedWeight: maxBillableWeight },
      },
      {
        onSuccess: (data, variables) => {
          setIsWeightModalVisible(false);
          setMessage(
            `MSG_MAX_BILLABLE_WEIGHT_SUCCESS_${variables.orderID}`,
            'success',
            'The maximum billable weight has been updated.',
            '',
            true,
          );
        },
      },
    );
  };

  const handleAcknowledgeExcessWeightRisk = () => {
    mutateAcknowledgeExcessWeightRisk({ orderID: order.id, ifMatchETag: move.eTag });
  };

  useEffect(() => {
    let serviceItemCount = 0;
    const serviceItemsCountForShipment = {};
    mtoShipments?.forEach((mtoShipment) => {
      if (
        mtoShipment.status === shipmentStatuses.APPROVED ||
        mtoShipment.status === shipmentStatuses.DIVERSION_REQUESTED
      ) {
        const requestedServiceItemCount = shipmentServiceItems[`${mtoShipment.id}`]?.filter(
          (serviceItem) => serviceItem.status === SERVICE_ITEM_STATUSES.SUBMITTED,
        )?.length;
        serviceItemCount += requestedServiceItemCount || 0;
        serviceItemsCountForShipment[`${mtoShipment.id}`] = requestedServiceItemCount;
      }
    });
    setUnapprovedServiceItemCount(serviceItemCount);
    setUnapprovedServiceItemsForShipment(serviceItemsCountForShipment);
  }, [mtoShipments, shipmentServiceItems, setUnapprovedServiceItemCount]);

  useEffect(() => {
    if (mtoShipments) {
      const shipmentCount = mtoShipments?.length
        ? mtoShipments.filter((shipment) => shipment.status === shipmentStatuses.SUBMITTED).length
        : 0;
      setUnapprovedShipmentCount(shipmentCount);

      const externalVendorShipments = mtoShipments?.length
        ? mtoShipments.filter((shipment) => shipment.usesExternalVendor).length
        : 0;
      setExternalVendorShipmentCount(externalVendorShipments);
    }
  }, [mtoShipments, setUnapprovedShipmentCount]);

  useEffect(() => {
    const shipmentSections = mtoShipments?.reduce((previous, shipment) => {
      if (showShipmentFilter(shipment)) {
        previous.push({
          id: shipment.id,
          label: shipmentSectionLabels[`${shipment.shipmentType}`] || shipment.shipmentType,
        });
      }
      return previous;
    }, []);
    setSections(shipmentSections || []);
  }, [mtoShipments]);

  useEffect(() => {
    setEstimatedWeightTotal(calculateEstimatedWeight(mtoShipments));
    let excessBillableWeightCount = 0;
    const riskOfExcessAcknowledged = !!move?.excess_weight_acknowledged_at;

    if (hasRiskOfExcess(estimatedWeightTotal, order?.entitlement.totalWeight) && !riskOfExcessAcknowledged) {
      excessBillableWeightCount = 1;
      setExcessWeightRiskCount(1);
    } else {
      setExcessWeightRiskCount(0);
    }

    const showWeightAlert = !riskOfExcessAcknowledged && !!excessBillableWeightCount;

    setIsWeightAlertVisible(showWeightAlert);
  }, [
    estimatedWeightTotal,
    move?.excess_weight_acknowledged_at,
    mtoShipments,
    order?.entitlement.totalWeight,
    setEstimatedWeightTotal,
    setExcessWeightRiskCount,
  ]);

  // Edge case of diversion shipments being counted twice
  const moveWeightTotal = calculateWeightRequested(mtoShipments);

  useEffect(() => {
    let unapprovedSITExtensionCount = 0;
    mtoShipments?.forEach((mtoShipment) => {
      if (mtoShipment.sitExtensions?.find((sitEx) => sitEx.status === SIT_EXTENSION_STATUS.PENDING)) {
        unapprovedSITExtensionCount += 1;
        unapprovedSITExtensionForShipment[`${mtoShipment.id}`] = 1;
        setUnApprovedSITExtensionForShipment(unapprovedSITExtensionForShipment);
      }
    });
    setUnapprovedSITExtensionCount(unapprovedSITExtensionCount);
  }, [
    mtoShipments,
    setUnapprovedSITExtensionCount,
    setUnApprovedSITExtensionForShipment,
    unapprovedSITExtensionForShipment,
  ]);

  const handleShowRejectionDialog = (mtoServiceItemID, mtoShipmentID) => {
    const serviceItem = shipmentServiceItems[`${mtoShipmentID}`]?.find((item) => item.id === mtoServiceItemID);
    setSelectedServiceItem(serviceItem);
    setIsModalVisible(true);
  };

  const handleShowCancellationModal = (mtoShipment) => {
    setSelectedShipment(mtoShipment);
    setIsCancelModalVisible(true);
  };

  const handleRequestReweighModal = (mtoShipment) => {
    setSelectedShipment(mtoShipment);
    setIsReweighModalVisible(true);
  };

  const handleShowWeightModal = () => {
    setIsWeightModalVisible(true);
  };

  const handleHideWeightAlert = () => {
    handleAcknowledgeExcessWeightRisk();
    setIsWeightAlertVisible(false);
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  if (move.status === MOVE_STATUSES.SUBMITTED || !mtoShipments.some(showShipmentFilter)) {
    return (
      <div className={styles.tabContent}>
        <GridContainer className={styles.gridContainer} data-testid="too-shipment-container">
          <div className={styles.pageHeader}>
            <h1>Move task order</h1>
          </div>
          <div className={styles.emptyMessage}>
            <p>This move does not have any approved shipments yet.</p>
          </div>
        </GridContainer>
      </div>
    );
  }

  const excessWeightAlertControl = (
    <Button type="button" onClick={handleHideWeightAlert} unstyled>
      <FontAwesomeIcon icon="times" />
    </Button>
  );

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        {/* nav is being used here instead of LeftNav since there are two separate sections that need to be interated through */}
        <nav className={classnames(leftNavStyles.LeftNav)}>
          {nonShipmentSections.map((s) => {
            return (
              <LeftNavSection
                key={`sidenav_${s}`}
                sectionName={s}
                isActive={`${s}` === activeSection}
                onClickHandler={() => setActiveSection(`${s}`)}
              >
                {nonShipmentSectionLabels[`${s}`]}
              </LeftNavSection>
            );
          })}
          {sections.map((s) => {
            return (
              <LeftNavSection
                key={`sidenav_${s.id}`}
                sectionName={`s-${s.id}`}
                isActive={`s-${s.id}` === activeSection}
                onClickHandler={() => setActiveSection(`s-${s.id}`)}
              >
                {s.label}{' '}
                <LeftNavTag
                  showTag={Boolean(
                    unapprovedServiceItemsForShipment[`${s.id}`] || unapprovedSITExtensionForShipment[`${s.id}`],
                  )}
                >
                  {(unapprovedServiceItemsForShipment[`${s.id}`] || 0) +
                    (unapprovedSITExtensionForShipment[`${s.id}`] || 0)}
                </LeftNavTag>
              </LeftNavSection>
            );
          })}
        </nav>
        <FlashGridContainer className={styles.gridContainer} data-testid="too-shipment-container">
          <Grid row className={styles.pageHeader}>
            {alertMessage && (
              <Grid col={12} className={styles.alertContainer}>
                <Alert headingLevel="h4" slim type={alertType}>
                  {alertMessage}
                </Alert>
              </Grid>
            )}
          </Grid>
          {isWeightAlertVisible && (
            <Alert
              headingLevel="h4"
              slim
              type="warning"
              cta={excessWeightAlertControl}
              className={styles.alertWithButton}
            >
              <span>
                This move is at risk for excess weight.{' '}
                <Restricted to={permissionTypes.updateBillableWeight}>
                  <span className={styles.rightAlignButtonWrapper}>
                    <Button type="button" onClick={handleShowWeightModal} unstyled>
                      Review billable weight
                    </Button>
                  </span>
                </Restricted>
              </span>
            </Alert>
          )}
          {isSuccessAlertVisible && (
            <Alert headingLevel="h4" slim type="success">
              Your changes were saved
            </Alert>
          )}

          {isModalVisible && (
            <RejectServiceItemModal
              serviceItem={selectedServiceItem}
              onSubmit={handleUpdateMTOServiceItemStatus}
              onClose={setIsModalVisible}
            />
          )}
          {isCancelModalVisible && (
            <RequestShipmentCancellationModal
              shipmentInfo={selectedShipment}
              onClose={setIsCancelModalVisible}
              onSubmit={handleUpdateMTOShipmentStatus}
            />
          )}
          {isReweighModalVisible && (
            <RequestReweighModal
              shipmentInfo={selectedShipment}
              onClose={setIsReweighModalVisible}
              onSubmit={handleReweighShipment}
            />
          )}

          <ConnectedEditMaxBillableWeightModal
            isOpen={isWeightModalVisible}
            defaultWeight={order.entitlement.totalWeight}
            maxBillableWeight={order.entitlement.authorizedWeight}
            onSubmit={handleUpdateBillableWeight}
            onClose={setIsWeightModalVisible}
          />

          {isFinancialModalVisible && (
            <FinancialReviewModal
              onClose={handleCancelFinancialReviewModal}
              onSubmit={handleSubmitFinancialReviewModal}
              initialRemarks={move?.financialReviewRemarks}
              initialSelection={move?.financialReviewFlag}
            />
          )}
          <div className={styles.pageHeader}>
            <h1>Move task order</h1>
            <div className={styles.pageHeaderDetails}>
              <h6>MTO Reference ID #{move?.referenceId}</h6>
              <h6>Contract #1234567890</h6> {/* TODO - need this value from the API */}
              <Restricted to={permissionTypes.updateFinancialReviewFlag}>
                <div className={moveTaskOrderStyles.financialReviewContainer}>
                  <FinancialReviewButton
                    onClick={handleShowFinancialReviewModal}
                    reviewRequested={move.financialReviewFlag}
                  />
                </div>
              </Restricted>
            </div>
          </div>

          <div className={moveTaskOrderStyles.weightHeader} id="move-weights">
            <WeightDisplay heading="Weight allowance" weightValue={order.entitlement.totalWeight} />
            <WeightDisplay heading="Estimated weight (total)" weightValue={estimatedWeightTotal}>
              {hasRiskOfExcess(estimatedWeightTotal, order.entitlement.totalWeight) && <Tag>Risk of excess</Tag>}
              {hasRiskOfExcess(estimatedWeightTotal, order.entitlement.totalWeight) &&
                externalVendorShipmentCount > 0 && <br />}
              {externalVendorShipmentCount > 0 && (
                <small>
                  {externalVendorShipmentCount} shipment{externalVendorShipmentCount > 1 && 's'} not moved by GHC prime.{' '}
                  <Link className="usa-link" to={generatePath(tooRoutes.MOVE_VIEW_PATH, { moveCode })}>
                    View move details
                  </Link>
                </small>
              )}
            </WeightDisplay>
            <WeightDisplay
              heading="Max billable weight"
              weightValue={order.entitlement.authorizedWeight}
              onEdit={handleShowWeightModal}
            />
            <WeightDisplay heading="Move weight (total)" weightValue={moveWeightTotal} />
          </div>
          {mtoShipments.map((mtoShipment) => {
            if (
              mtoShipment.status !== shipmentStatuses.APPROVED &&
              mtoShipment.status !== shipmentStatuses.CANCELLATION_REQUESTED &&
              mtoShipment.status !== shipmentStatuses.DIVERSION_REQUESTED &&
              mtoShipment.status !== shipmentStatuses.CANCELED
            ) {
              return false;
            }

            const serviceItemsForShipment = shipmentServiceItems[`${mtoShipment.id}`];
            const requestedServiceItems = serviceItemsForShipment?.filter(
              (item) => item.status === SERVICE_ITEM_STATUSES.SUBMITTED,
            );
            const approvedServiceItems = serviceItemsForShipment?.filter(
              (item) => item.status === SERVICE_ITEM_STATUSES.APPROVED,
            );
            const rejectedServiceItems = serviceItemsForShipment?.filter(
              (item) => item.status === SERVICE_ITEM_STATUSES.REJECTED,
            );
            const dutyLocationPostal = { postalCode: order.destinationDutyLocation.address.postalCode };
            const { pickupAddress, destinationAddress } = mtoShipment;
            const formattedScheduledPickup = formatShipmentDate(mtoShipment.scheduledPickupDate);

            return (
              <ShipmentContainer
                id={`s-${mtoShipment.id}`}
                key={mtoShipment.id}
                shipmentType={mtoShipment.shipmentType}
                className={styles.shipmentCard}
              >
                <ShipmentHeading
                  shipmentInfo={{
                    shipmentID: mtoShipment.id,
                    shipmentType: mtoShipmentTypes[mtoShipment.shipmentType],
                    isDiversion: mtoShipment.diversion,
                    originCity: pickupAddress?.city || '',
                    originState: pickupAddress?.state || '',
                    originPostalCode: pickupAddress?.postalCode || '',
                    destinationAddress: destinationAddress || dutyLocationPostal,
                    scheduledPickupDate: formattedScheduledPickup,
                    shipmentStatus: mtoShipment.status,
                    ifMatchEtag: mtoShipment.eTag,
                    moveTaskOrderID: mtoShipment.moveTaskOrderID,
                  }}
                  handleShowCancellationModal={handleShowCancellationModal}
                />
                <ShipmentDetails
                  shipment={mtoShipment}
                  order={order}
                  handleDivertShipment={handleDivertShipment}
                  handleRequestReweighModal={handleRequestReweighModal}
                  handleReviewSITExtension={handleReviewSITExtension}
                  handleSubmitSITExtension={handleSubmitSITExtension}
                  handleEditFacilityInfo={handleEditFacilityInfo}
                  handleEditServiceOrderNumber={handleEditServiceOrderNumber}
                  handleEditAccountingCodes={handleEditAccountingCodes}
                />
                {requestedServiceItems?.length > 0 && (
                  <RequestedServiceItemsTable
                    serviceItems={requestedServiceItems}
                    handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                    handleShowRejectionDialog={handleShowRejectionDialog}
                    statusForTableType={SERVICE_ITEM_STATUSES.SUBMITTED}
                  />
                )}
                {approvedServiceItems?.length > 0 && (
                  <RequestedServiceItemsTable
                    serviceItems={approvedServiceItems}
                    handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                    handleShowRejectionDialog={handleShowRejectionDialog}
                    statusForTableType={SERVICE_ITEM_STATUSES.APPROVED}
                  />
                )}
                {rejectedServiceItems?.length > 0 && (
                  <RequestedServiceItemsTable
                    serviceItems={rejectedServiceItems}
                    handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                    handleShowRejectionDialog={handleShowRejectionDialog}
                    statusForTableType={SERVICE_ITEM_STATUSES.REJECTED}
                  />
                )}
              </ShipmentContainer>
            );
          })}
        </FlashGridContainer>
      </div>
    </div>
  );
};

MoveTaskOrder.propTypes = {
  setUnapprovedShipmentCount: func.isRequired,
  setUnapprovedServiceItemCount: func.isRequired,
  setExcessWeightRiskCount: func.isRequired,
  setMessage: func.isRequired,
  setUnapprovedSITExtensionCount: func.isRequired,
};

const mapDispatchToProps = {
  setMessage: setFlashMessage,
};

export default connect(() => ({}), mapDispatchToProps)(MoveTaskOrder);
