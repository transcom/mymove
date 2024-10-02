import React, { useEffect, useMemo, useState } from 'react';
import { generatePath, Link, useParams } from 'react-router-dom';
import { Alert, Button, Grid, GridContainer, Tag } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { connect } from 'react-redux';
import { func } from 'prop-types';
import classnames from 'classnames';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import moveTaskOrderStyles from './MoveTaskOrder.module.scss';

import ConnectedEditMaxBillableWeightModal from 'components/Office/EditMaxBillableWeightModal/EditMaxBillableWeightModal';
import { milmoveLogger } from 'utils/milmoveLog';
import { formatAddressForAPI, formatStorageFacilityForAPI, removeEtag } from 'utils/formatMtoShipment';
import hasRiskOfExcess from 'utils/hasRiskOfExcess';
import dimensionTypes from 'constants/dimensionTypes';
import { MOVES, MTO_SERVICE_ITEMS, MTO_SHIPMENTS, ORDERS } from 'constants/queryKeys';
import SERVICE_ITEM_STATUSES from 'constants/serviceItems';
import { mtoShipmentTypes, shipmentStatuses } from 'constants/shipments';
import FlashGridContainer from 'containers/FlashGridContainer/FlashGridContainer';
import { shipmentSectionLabels } from 'content/shipments';
import RejectServiceItemModal from 'components/Office/RejectServiceItemModal/RejectServiceItemModal';
import RequestedServiceItemsTable from 'components/Office/RequestedServiceItemsTable/RequestedServiceItemsTable';
import RequestShipmentCancellationModal from 'components/Office/RequestShipmentCancellationModal/RequestShipmentCancellationModal';
import RequestShipmentDiversionModal from 'components/Office/RequestShipmentDiversionModal/RequestShipmentDiversionModal';
import RequestReweighModal from 'components/Office/RequestReweighModal/RequestReweighModal';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import ShipmentHeading from 'components/Office/ShipmentHeading/ShipmentHeading';
import ShipmentDetails from 'components/Office/ShipmentDetails/ShipmentDetails';
import ServiceItemContainer from 'components/Office/ServiceItemContainer/ServiceItemContainer';
import { useMoveTaskOrderQueries } from 'hooks/queries';
import {
  acknowledgeExcessWeightRisk,
  approveSITExtension,
  denySITExtension,
  patchMTOServiceItemStatus,
  submitSITExtension,
  updateBillableWeight,
  updateFinancialFlag,
  updateMTOShipment,
  updateMTOShipmentRequestReweigh,
  updateMTOShipmentStatus,
  updateServiceItemSITEntryDate,
  updateSITServiceItemCustomerExpense,
} from 'services/ghcApi';
import { MOVE_STATUSES, MTO_SERVICE_ITEM_STATUS, SHIPMENT_OPTIONS } from 'shared/constants';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { setFlashMessage } from 'store/flash/actions';
import WeightDisplay from 'components/Office/WeightDisplay/WeightDisplay';
import {
  calculateEstimatedWeight,
  calculateWeightRequested,
  includedStatusesForCalculatingWeights,
} from 'hooks/custom';
import { SIT_EXTENSION_STATUS } from 'constants/sitExtensions';
import FinancialReviewButton from 'components/Office/FinancialReviewButton/FinancialReviewButton';
import FinancialReviewModal from 'components/Office/FinancialReviewModal/FinancialReviewModal';
import leftNavStyles from 'components/LeftNav/LeftNav.module.scss';
import LeftNavSection from 'components/LeftNavSection/LeftNavSection';
import LeftNavTag from 'components/LeftNavTag/LeftNavTag';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import { tooRoutes } from 'constants/routes';
import { formatDateForSwagger } from 'shared/dates';
import EditSitEntryDateModal from 'components/Office/EditSitEntryDateModal/EditSitEntryDateModal';
import { formatWeight } from 'utils/formatters';
import { roleTypes } from 'constants/userRoles';

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
  /* ------------------ Modals ------------------------- */
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [isCancelModalVisible, setIsCancelModalVisible] = useState(false);
  // Diversion
  const [isDiversionModalVisible, setIsDiversionModalVisible] = useState(false);
  // Weights
  const [isReweighModalVisible, setIsReweighModalVisible] = useState(false);
  const [isWeightModalVisible, setIsWeightModalVisible] = useState(false);
  // SIT Address Updates
  const [isEditSitEntryDateModalVisible, setIsEditSitEntryDateModalVisible] = useState(false);
  /* ------------------ Alerts ------------------------- */
  const [alertMessage, setAlertMessage] = useState(null);
  const [alertType, setAlertType] = useState('success');
  const [isSuccessAlertVisible, setIsSuccessAlertVisible] = useState(false);
  const [isWeightAlertVisible, setIsWeightAlertVisible] = useState(false);
  const [isFinancialModalVisible, setIsFinancialModalVisible] = useState(false);
  /* ------------------ Selected / Active Item ------------------------- */
  const [selectedShipment, setSelectedShipment] = useState(undefined);
  const [selectedServiceItem, setSelectedServiceItem] = useState(undefined);
  const [activeSection, setActiveSection] = useState('');
  const [sections, setSections] = useState([]);
  /* ------------------ Unapproved requests / counts ------------------------- */
  const [unapprovedServiceItemsForShipment, setUnapprovedServiceItemsForShipment] = useState({});
  const [unapprovedSITExtensionForShipment, setUnApprovedSITExtensionForShipment] = useState({});
  const [externalVendorShipmentCount, setExternalVendorShipmentCount] = useState(0);
  /* ------------------ Miscellaneous ------------------------- */
  const [estimatedWeightTotal, setEstimatedWeightTotal] = useState(null);
  const [estimatedHHGWeightTotal, setEstimatedHHGWeightTotal] = useState(null);
  const [estimatedNTSWeightTotal, setEstimatedNTSWeightTotal] = useState(null);
  const [estimatedNTSReleaseWeightTotal, setEstimatedNTSReleaseWeightTotal] = useState(null);
  const [estimatedPPMWeightTotal, setEstimatedPPMWeightTotal] = useState(null);
  const [, setSubmittedChangeTime] = useState(Date.now());
  const [breakdownVisible, setBreakdownVisible] = useState(false);

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
    userRole,
    isMoveLocked,
  } = props;

  const { orders = {}, move, mtoShipments, mtoServiceItems, isLoading, isError } = useMoveTaskOrderQueries(moveCode);
  const order = Object.values(orders)?.[0];
  const nonPPMShipments = mtoShipments?.filter((shipment) => shipment.shipmentType !== 'PPM');
  const onlyPPMShipments = mtoShipments?.filter((shipment) => shipment.shipmentType === 'PPM');

  const shipmentServiceItems = useMemo(() => {
    const serviceItemsForShipment = {};
    mtoServiceItems?.forEach((item) => {
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
        customerContacts: item.customerContacts,
        estimatedWeight: item.estimatedWeight,
        rejectionReason: item.rejectionReason,
        sitDepartureDate: item.sitDepartureDate,
        sitEntryDate: item.sitEntryDate,
        sitOriginHHGOriginalAddress: item.sitOriginHHGOriginalAddress,
        sitOriginHHGActualAddress: item.sitOriginHHGActualAddress,
        sitDestinationFinalAddress: item.sitDestinationFinalAddress,
        sitDestinationOriginalAddress: item.sitDestinationOriginalAddress,
        sitCustomerContacted: item.sitCustomerContacted,
        sitRequestedDelivery: item.sitRequestedDelivery,
        sitDeliveryMiles: item.sitDeliveryMiles,
        status: item.status,
        estimatedPrice: item.estimatedPrice,
        standaloneCrate: item.standaloneCrate,
        lockedPricedCents: item.lockedPriceCents,
      };

      if (serviceItemsForShipment[`${newItem.mtoShipmentID}`]) {
        serviceItemsForShipment[`${newItem.mtoShipmentID}`].push(newItem);
      } else {
        serviceItemsForShipment[`${newItem.mtoShipmentID}`] = [newItem]; // Basic service items belong under shipmentServiceItems[`${undefined}`]
      }
    });
    return serviceItemsForShipment;
  }, [mtoServiceItems]);

  const serviceItemsForMove = shipmentServiceItems[`${undefined}`];
  const requestedMoveServiceItems = serviceItemsForMove?.filter(
    (item) => item.status === SERVICE_ITEM_STATUSES.SUBMITTED,
  );
  const approvedMoveServiceItems = serviceItemsForMove?.filter(
    (item) => item.status === SERVICE_ITEM_STATUSES.APPROVED,
  );
  const rejectedMoveServiceItems = serviceItemsForMove?.filter(
    (item) => item.status === SERVICE_ITEM_STATUSES.REJECTED,
  );

  /*
  *
  -------------------------  Mutation Funtions  -------------------------
  * using istanbul ignore next to omit from test coverage since these functions
  * cannot be exported
  */

  const queryClient = useQueryClient();
  /* istanbul ignore next */
  const { mutate: mutateMTOServiceItemStatus } = useMutation({
    mutationFn: patchMTOServiceItemStatus,
    onSuccess: (data, variables) => {
      const newMTOServiceItem = data.mtoServiceItems[variables.mtoServiceItemID];
      mtoServiceItems[mtoServiceItems.find((serviceItem) => serviceItem.id === newMTOServiceItem.id)] =
        newMTOServiceItem;
      queryClient.setQueryData([MTO_SERVICE_ITEMS, variables.moveId, false], mtoServiceItems);
      queryClient.invalidateQueries({ queryKey: [MTO_SERVICE_ITEMS, variables.moveId] });
      queryClient.invalidateQueries({ queryKey: [MTO_SHIPMENTS] });
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  /* istanbul ignore next */
  const { mutate: mutateSITServiceItemCustomerExpense } = useMutation({
    mutationFn: updateSITServiceItemCustomerExpense,
    onSuccess: (data, variables) => {
      const updatedMTOShipment = data.mtoShipments[variables.shipmentID];
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryClient.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryClient.invalidateQueries({ queryKey: [MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID] });
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  /* istanbul ignore next */
  const { mutate: mutateMTOShipment } = useMutation({
    mutationFn: updateMTOShipment,
    onSuccess: (_, variables) => {
      queryClient.setQueryData([MTO_SHIPMENTS, variables.moveTaskOrderID, false], mtoShipments);
      queryClient.invalidateQueries({ queryKey: [MTO_SHIPMENTS, variables.moveTaskOrderID] });
    },
  });

  /* istanbul ignore next */
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
      milmoveLogger.error(errorMsg);
    },
  });

  /* istanbul ignore next */
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
      milmoveLogger.error(errorMsg);
    },
  });

  /* istanbul ignore next */
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
      milmoveLogger.error(errorMsg);
    },
  });

  /* istanbul ignore next */
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

  /* istanbul ignore next */
  const { mutate: mutateSITExtensionApproval } = useMutation({
    mutationFn: approveSITExtension,
    onSuccess: (data, variables) => {
      const updatedMTOShipment = data.mtoShipments[variables.shipmentID];
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryClient.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryClient.invalidateQueries({ queryKey: [MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID] });
      setSubmittedChangeTime(Date.now());
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  /* istanbul ignore next */
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
      milmoveLogger.error(errorMsg);
    },
  });

  /* istanbul ignore next */
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
      milmoveLogger.error(errorMsg);
    },
  });

  /* istanbul ignore next */
  const { mutate: mutateFinancialReview } = useMutation({
    mutationFn: updateFinancialFlag,
    onSuccess: (data) => {
      queryClient.setQueryData([MOVES, data.locator], data);
      queryClient.invalidateQueries({ queryKey: [MOVES, data.locator] });
    },
  });

  /* istanbul ignore next */
  const { mutate: mutateServiceItemSitEntryDate } = useMutation({
    mutationFn: updateServiceItemSITEntryDate,
    onSuccess: (data) => {
      // here we are updating the service item
      const updatedServiceItems = [...mtoServiceItems];
      updatedServiceItems[updatedServiceItems.findIndex((serviceItem) => serviceItem.id === data.id)] = data;
      queryClient.setQueryData([MTO_SERVICE_ITEMS, move.id, false], updatedServiceItems);
      queryClient.invalidateQueries({ queryKey: [MTO_SERVICE_ITEMS, move.id, false] });

      // here we are updating the shipment (focusing on the currentSit object)
      const updatedMTOShipment = data;
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === data.mtoShipmentID)] = updatedMTOShipment;
      queryClient.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryClient.invalidateQueries({ queryKey: [MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID] });
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });
  /*
    *
    -------------------------  Toggle Modals  -------------------------
                  Functions to show and hide modals
    * using istanbul ignore next to omit from test coverage since these functions
    * cannot be exported
    *
    */

  /* istanbul ignore next */
  const handleCancelFinancialReviewModal = () => {
    setIsFinancialModalVisible(false);
  };

  /* istanbul ignore next */
  const handleShowFinancialReviewModal = () => {
    setIsFinancialModalVisible(true);
  };

  /* istanbul ignore next */
  const handleShowRejectionDialog = (mtoServiceItemID, mtoShipmentID) => {
    const serviceItem = shipmentServiceItems[`${mtoShipmentID}`]?.find((item) => item.id === mtoServiceItemID);
    setSelectedServiceItem(serviceItem);
    setIsModalVisible(true);
  };

  /* istanbul ignore next */
  const handleShowEditSitEntryDateModal = (mtoServiceItemID, mtoShipmentID) => {
    const serviceItem = shipmentServiceItems[`${mtoShipmentID}`]?.find((item) => item.id === mtoServiceItemID);
    setSelectedServiceItem(serviceItem);
    setIsEditSitEntryDateModalVisible(true);
  };

  /* istanbul ignore next */
  const handleCancelEditSitEntryDateModal = () => {
    setIsEditSitEntryDateModalVisible(false);
  };

  /* istanbul ignore next */
  const handleShowCancellationModal = (mtoShipment) => {
    setSelectedShipment(mtoShipment);
    setIsCancelModalVisible(true);
  };

  /* istanbul ignore next */
  const handleShowDiversionModal = (mtoShipment) => {
    setSelectedShipment(mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === mtoShipment.id)]);
    setIsDiversionModalVisible(true);
  };
  /* istanbul ignore next */
  const handleRequestReweighModal = (mtoShipment) => {
    setSelectedShipment(mtoShipment);
    setIsReweighModalVisible(true);
  };

  // To-do: Combine handle Acknowldge Weights and handle Weight alert into one one mutation function
  const handleAcknowledgeExcessWeightRisk = () => {
    mutateAcknowledgeExcessWeightRisk({ orderID: order.id, ifMatchETag: move.eTag });
  };
  const handleHideWeightAlert = () => {
    handleAcknowledgeExcessWeightRisk();
    setIsWeightAlertVisible(false);
  };

  const handleShowWeightModal = () => {
    handleHideWeightAlert();
    setIsWeightModalVisible(true);
  };
  /*
  *
  -------------------------  Submit Handlers  -------------------------
              Contain mutation functions to handle form submissions
              Using istanbul ignore next to omit from test coverage
              since they cannot be exported
  *
  */

  /* istanbul ignore next */
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
  /* istanbul ignore next */
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
        body: {
          officeRemarks: formValues.officeRemarks,
          convertToCustomerExpense: formValues.convertToCustomerExpense,
        },
      });
    }
    setSubmittedChangeTime(Date.now());
  };

  /* istanbul ignore next */
  const handleSubmitSITExtension = (formValues, shipment) => {
    mutateSubmitSITExtension(
      {
        shipmentID: shipment.id,
        ifMatchETag: shipment.eTag,
        body: {
          requestReason: formValues.requestReason,
          officeRemarks: formValues.officeRemarks,
          approvedDays: parseInt(formValues.daysApproved, 10) - shipment.sitDaysAllowance,
          sitEntryDate: formatDateForSwagger(formValues.sitEntryDate),
          moveID: shipment.moveTaskOrderID,
        },
      },
      {
        onSuccess: () => {
          setIsSuccessAlertVisible(true);
          setSubmittedChangeTime(Date.now());
        },
      },
    );
  };

  /* istanbul ignore next */
  const handleDivertShipment = (mtoShipmentID, eTag, shipmentLocator, diversionReason) => {
    mutateMTOShipmentStatus(
      {
        shipmentID: mtoShipmentID,
        operationPath: 'shipment.requestShipmentDiversion',
        ifMatchETag: eTag,
        onSuccessFlashMsg: `Diversion successfully requested for Shipment #${shipmentLocator}`,
        shipmentLocator,
        diversionReason,
      },
      {
        onSuccess: (data, variables) => {
          setIsDiversionModalVisible(false);
          // Must set FlashMesage after hiding the modal, since FlashMessage will disappear when focus changes
          setMessage(
            `MSG_CANCEL_SUCCESS_${variables.shipmentLocator}`,
            'success',
            variables.onSuccessFlashMsg,
            '',
            true,
          );
        },
        onError: () => {
          setIsDiversionModalVisible(false);
          setAlertMessage('There was a problem requesting a diversion on this shipment. Please try again later.');
          setAlertType('error');
        },
      },
    );
  };

  /* istanbul ignore next */
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

  /* istanbul ignore next */
  const handleEditAccountingCodes = (fields, shipment) => {
    const body = { tacType: null, sacType: null, ...fields };
    mutateMTOShipment({
      moveTaskOrderID: shipment.moveTaskOrderID,
      shipmentID: shipment.id,
      ifMatchETag: shipment.eTag,
      body,
    });
  };

  /* istanbul ignore next */
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
        onError: (data, error) => {
          const errorMsg = error?.response?.body;
          milmoveLogger.error(errorMsg);
          setIsCancelModalVisible(false);
          setAlertMessage(`${data.response.body.message}`);
          setAlertType('error');
        },
      },
    );
  };

  /* istanbul ignore next */
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

  /* istanbul ignore next */
  const handleEditServiceOrderNumber = (fields, shipment) => {
    mutateMTOShipment({
      moveTaskOrderID: shipment.moveTaskOrderID,
      shipmentID: shipment.id,
      ifMatchETag: shipment.eTag,
      body: { serviceOrderNumber: fields.serviceOrderNumber },
    });
  };

  /* istanbul ignore next */
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

  /* istanbul ignore next */
  const handleUpdateSITServiceItemCustomerExpense = (
    mtoShipmentID,
    convertToCustomerExpense,
    customerExpenseReason,
    eTag,
  ) => {
    mutateSITServiceItemCustomerExpense(
      {
        shipmentID: mtoShipmentID,
        convertToCustomerExpense,
        customerExpenseReason,
        ifMatchETag: eTag,
        onSuccessFlashMsg: `SIT successfully converted to customer expense`,
      },
      {
        onSuccess: (data, variables) => {
          setMessage(
            `MSG_CONVERT_TO_CUSTOMER_EXPENSE_SUCCESS_${variables.shipmentID}`,
            'success',
            variables.onSuccessFlashMsg,
            '',
            true,
          );
        },
      },
    );
  };

  /* istanbul ignore next */
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

  /**
   * @typedef AddressShape
   * @prop {string} city
   * @prop {string} state
   * @prop {string} postalCode
   * @prop {string} streetAddress1
   * @prop {string} streetAddress2
   * @prop {string} streetAddress3
   * @prop {string} country
   */

  /**
   * @function
   * @param {string} mtoServiceItemID
   * @param {Date} newSitEntryDate
   * @description Updates the selected SIT entry date
   * OnSuccess, it closes the modal and sets a success message.
   */
  /* istanbul ignore next */
  const handleSubmitSitEntryDateChange = (mtoServiceItemID, newSitEntryDate) => {
    mutateServiceItemSitEntryDate(
      {
        mtoServiceItemID,
        body: { ID: mtoServiceItemID, SitEntryDate: newSitEntryDate },
      },
      {
        onSuccess: () => {
          setSelectedServiceItem({});
          setIsEditSitEntryDateModalVisible(false);
          setAlertMessage('SIT entry date updated');
          setAlertType('success');
        },
      },
    );
  };

  /*
  *
  -------------------------  useEffect Handlers  -------------------------
  *
  */

  /* ------------------ Update Notification counts ------------------------- */
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

  /* ------------------ Update Shipment approvals ------------------------- */
  useEffect(() => {
    if (mtoShipments && userRole !== roleTypes.SERVICES_COUNSELOR) {
      const shipmentCount = mtoShipments?.length
        ? mtoShipments.filter((shipment) => shipment.status === shipmentStatuses.SUBMITTED).length
        : 0;
      setUnapprovedShipmentCount(shipmentCount);

      const externalVendorShipments = mtoShipments?.length
        ? mtoShipments.filter((shipment) => shipment.usesExternalVendor).length
        : 0;
      setExternalVendorShipmentCount(externalVendorShipments);
    }
  }, [mtoShipments, setUnapprovedShipmentCount, userRole]);

  /* ------------------ Update Weight related alerts and estimates ------------------------- */
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
    setEstimatedWeightTotal(calculateEstimatedWeight(nonPPMShipments));
    setEstimatedHHGWeightTotal(calculateEstimatedWeight(nonPPMShipments, SHIPMENT_OPTIONS.HHG));
    setEstimatedNTSWeightTotal(calculateEstimatedWeight(nonPPMShipments, SHIPMENT_OPTIONS.NTS));
    setEstimatedNTSReleaseWeightTotal(calculateEstimatedWeight(nonPPMShipments, SHIPMENT_OPTIONS.NTSR));
    setEstimatedPPMWeightTotal(calculateEstimatedWeight(onlyPPMShipments));
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
    nonPPMShipments,
    onlyPPMShipments,
    order?.entitlement.totalWeight,
    setEstimatedWeightTotal,
    setEstimatedHHGWeightTotal,
    setEstimatedNTSWeightTotal,
    setEstimatedNTSReleaseWeightTotal,
    setExcessWeightRiskCount,
  ]);

  /* ------------------ Update SIT extension counts ------------------------- */
  useEffect(() => {
    const copyItemsFromTempArrayToSourceArray = (temp, target) => {
      Object.keys(temp).forEach((item) => {
        const targetArray = target;
        targetArray[item] = temp[item];
      });
    };
    const checkShipmentsForUnapprovedSITExtensions = (shipmentsWithStatus) => {
      const unapprovedSITExtensionShipmentItems = [];
      let unapprovedSITExtensionCount = 0;
      shipmentsWithStatus?.forEach((mtoShipment) => {
        const unapprovedSITExtItems =
          mtoShipment.sitExtensions?.filter((sitEx) => sitEx.status === SIT_EXTENSION_STATUS.PENDING) ?? [];
        const unapprovedSITCount = unapprovedSITExtItems.length;
        unapprovedSITExtensionCount += unapprovedSITCount; // Top bar Label
        unapprovedSITExtensionShipmentItems[`${mtoShipment.id}`] = unapprovedSITCount; // Nav bar Label
      });
      return { count: unapprovedSITExtensionCount, items: unapprovedSITExtensionShipmentItems };
    };
    const { count, items } = checkShipmentsForUnapprovedSITExtensions(mtoShipments);
    setUnapprovedSITExtensionCount(count);
    copyItemsFromTempArrayToSourceArray(items, unapprovedSITExtensionForShipment);
    setUnApprovedSITExtensionForShipment(unapprovedSITExtensionForShipment);
  }, [
    mtoShipments,
    setUnapprovedSITExtensionCount,
    setUnApprovedSITExtensionForShipment,
    unapprovedSITExtensionForShipment,
  ]);

  /* ------------------ Utils ------------------------- */
  // determine if max billable weight should be displayed yet
  const displayMaxBillableWeight = (shipments) => {
    return shipments?.some(
      (shipment) =>
        includedStatusesForCalculatingWeights(shipment.status) &&
        (shipment.primeEstimatedWeight || shipment.ntsRecordedWeight),
    );
  };
  // Edge case of diversion shipments being counted twice
  const moveWeightTotal = calculateWeightRequested(nonPPMShipments);
  const ppmWeightTotal = calculateWeightRequested(onlyPPMShipments);
  const maxBillableWeight = displayMaxBillableWeight(nonPPMShipments) ? order?.entitlement?.authorizedWeight : '-';

  /*
  *
  -------------------------  UI -------------------------
  *
  */
  // this should always be 110% of estimated weight regardless of allowance
  // or max billable weight

  const estimateWeightBreakdown = (
    <div>
      <div>110% Estimated HHG</div>
      <div className={moveTaskOrderStyles.subValue}>
        {Number.isFinite(estimatedHHGWeightTotal) ? formatWeight(Math.round(estimatedHHGWeightTotal * 1.1)) : '—'}
      </div>
      <div>110% Estimated NTS</div>
      <div className={moveTaskOrderStyles.subValue}>
        {Number.isFinite(estimatedNTSWeightTotal) ? formatWeight(Math.round(estimatedNTSWeightTotal * 1.1)) : '—'}
      </div>
      <div>110% Estimated NTS-Release</div>
      <div className={moveTaskOrderStyles.subValue}>
        {Number.isFinite(estimatedNTSReleaseWeightTotal)
          ? formatWeight(Math.round(estimatedNTSReleaseWeightTotal * 1.1))
          : '—'}
      </div>
    </div>
  );

  const estimateWeight110 = (
    <div className={moveTaskOrderStyles.childHeader}>
      <div>110% of estimated weight (TOTAL)</div>
      <div className={moveTaskOrderStyles.value}>
        {Number.isFinite(estimatedWeightTotal) ? formatWeight(Math.round(estimatedWeightTotal * 1.1)) : '—'}
      </div>
      <Button
        className={styles.toggleBreakdown}
        type="button"
        data-testid="toggleBreakdown"
        aria-expanded={breakdownVisible}
        unstyled
        onClick={() => {
          setBreakdownVisible((isVisible) => {
            return !isVisible;
          });
        }}
      >
        {breakdownVisible ? 'Hide Breakdown' : 'Show Breakdown'}
      </Button>
      &nbsp;
      {breakdownVisible && estimateWeightBreakdown}
    </div>
  );

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  /* ------------------ No approved shipments ------------------------- */
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
    <Button
      data-testid="excessWeightAlertButton"
      type="button"
      onClick={handleHideWeightAlert}
      unstyled
      disabled={isMoveLocked}
    >
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
                  <Restricted to={permissionTypes.updateMTOPage}>
                    <span className={styles.rightAlignButtonWrapper}>
                      <Button
                        data-testid="reviewBillableWeightBtn"
                        type="button"
                        onClick={handleShowWeightModal}
                        unstyled
                        disabled={isMoveLocked}
                      >
                        Review billable weight
                      </Button>
                    </span>
                  </Restricted>
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
          {isDiversionModalVisible && (
            <RequestShipmentDiversionModal
              shipmentInfo={selectedShipment}
              onClose={setIsDiversionModalVisible}
              onSubmit={handleDivertShipment}
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
          {isEditSitEntryDateModalVisible && (
            <EditSitEntryDateModal
              onClose={handleCancelEditSitEntryDateModal}
              onSubmit={handleSubmitSitEntryDateChange}
              isOpen={isEditSitEntryDateModalVisible}
              serviceItem={selectedServiceItem}
              shipmentInfo={selectedShipment}
            />
          )}
          <div className={styles.pageHeader}>
            <h1>Move task order</h1>
            <div className={styles.pageHeaderDetails}>
              <h6>MTO Reference ID #{move?.referenceId}</h6>
              <h6>Contract #{move?.contractor?.contractNumber}</h6>
              <h6>NAICS: {order?.naics}</h6>
              <Restricted to={permissionTypes.updateFinancialReviewFlag}>
                <Restricted to={permissionTypes.updateMTOPage}>
                  <div className={moveTaskOrderStyles.financialReviewContainer}>
                    <FinancialReviewButton
                      onClick={handleShowFinancialReviewModal}
                      reviewRequested={move.financialReviewFlag}
                      isMoveLocked={isMoveLocked}
                    />
                  </div>
                </Restricted>
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
              {estimateWeight110}
            </WeightDisplay>
            <WeightDisplay
              heading="Max billable weight"
              weightValue={maxBillableWeight}
              onEdit={displayMaxBillableWeight(nonPPMShipments) ? handleShowWeightModal : null}
              isMoveLocked={isMoveLocked}
            />
            <WeightDisplay heading="Move weight (total)" weightValue={moveWeightTotal} />
          </div>
          {onlyPPMShipments.length > 0 && (
            <div className={moveTaskOrderStyles.secondRow} id="move-weights">
              <WeightDisplay heading="PPM estimated weight (total)" weightValue={estimatedPPMWeightTotal} />
              <WeightDisplay heading="Actual PPM weight (total)" weightValue={ppmWeightTotal} />
            </div>
          )}
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
                    shipmentLocator: mtoShipment.shipmentLocator,
                  }}
                  handleShowCancellationModal={handleShowCancellationModal}
                  isMoveLocked={isMoveLocked}
                />
                <ShipmentDetails
                  shipment={mtoShipment}
                  order={order}
                  handleRequestReweighModal={handleRequestReweighModal}
                  handleShowDiversionModal={handleShowDiversionModal}
                  handleReviewSITExtension={handleReviewSITExtension}
                  handleSubmitSITExtension={handleSubmitSITExtension}
                  handleUpdateSITServiceItemCustomerExpense={handleUpdateSITServiceItemCustomerExpense}
                  handleEditFacilityInfo={handleEditFacilityInfo}
                  handleEditServiceOrderNumber={handleEditServiceOrderNumber}
                  handleEditAccountingCodes={handleEditAccountingCodes}
                  isMoveLocked={isMoveLocked}
                />
                {requestedServiceItems?.length > 0 && (
                  <RequestedServiceItemsTable
                    serviceItems={requestedServiceItems}
                    handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                    handleShowRejectionDialog={handleShowRejectionDialog}
                    handleShowEditSitEntryDateModal={handleShowEditSitEntryDateModal}
                    statusForTableType={SERVICE_ITEM_STATUSES.SUBMITTED}
                    shipment={mtoShipment}
                    sitStatus={mtoShipment.sitStatus}
                    isMoveLocked={isMoveLocked}
                  />
                )}
                {approvedServiceItems?.length > 0 && (
                  <RequestedServiceItemsTable
                    serviceItems={approvedServiceItems}
                    handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                    handleShowRejectionDialog={handleShowRejectionDialog}
                    handleShowEditSitEntryDateModal={handleShowEditSitEntryDateModal}
                    statusForTableType={SERVICE_ITEM_STATUSES.APPROVED}
                    shipment={mtoShipment}
                    sitStatus={mtoShipment.sitStatus}
                    isMoveLocked={isMoveLocked}
                  />
                )}
                {rejectedServiceItems?.length > 0 && (
                  <RequestedServiceItemsTable
                    serviceItems={rejectedServiceItems}
                    handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                    handleShowRejectionDialog={handleShowRejectionDialog}
                    statusForTableType={SERVICE_ITEM_STATUSES.REJECTED}
                    shipment={mtoShipment}
                    sitStatus={mtoShipment.sitStatus}
                    isMoveLocked={isMoveLocked}
                  />
                )}
              </ShipmentContainer>
            );
          })}
          <ServiceItemContainer className={styles.shipmentCard}>
            {requestedMoveServiceItems?.length > 0 && (
              <RequestedServiceItemsTable
                serviceItems={requestedMoveServiceItems}
                handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                handleShowRejectionDialog={handleShowRejectionDialog}
                statusForTableType={MTO_SERVICE_ITEM_STATUS.SUBMITTED}
              />
            )}
            {approvedMoveServiceItems?.length > 0 && (
              <RequestedServiceItemsTable
                serviceItems={approvedMoveServiceItems}
                handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                handleShowRejectionDialog={handleShowRejectionDialog}
                statusForTableType={MTO_SERVICE_ITEM_STATUS.APPROVED}
              />
            )}
            {rejectedMoveServiceItems?.length > 0 && (
              <RequestedServiceItemsTable
                serviceItems={rejectedMoveServiceItems}
                handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                handleShowRejectionDialog={handleShowRejectionDialog}
                statusForTableType={MTO_SERVICE_ITEM_STATUS.REJECTED}
              />
            )}
          </ServiceItemContainer>
          <div className={styles.pageFooter}>
            <div className={styles.pageFooterDetails}>
              <h6>{order?.packingAndShippingInstructions}</h6>
              <h6>{order?.methodOfPayment}</h6>
            </div>
          </div>
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
