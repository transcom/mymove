import React, { useState, useEffect, useMemo } from 'react';
import { withRouter } from 'react-router-dom';
import { GridContainer, Tag } from '@trussworks/react-uswds';
import { queryCache, useMutation } from 'react-query';
import { connect } from 'react-redux';
import { func } from 'prop-types';
import classnames from 'classnames';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import { MTO_SERVICE_ITEMS, MTO_SHIPMENTS } from 'constants/queryKeys';
import SERVICE_ITEM_STATUSES from 'constants/serviceItems';
import { MOVE_STATUSES } from 'shared/constants';
import dimensionTypes from 'constants/dimensionTypes';
import customerContactTypes from 'constants/customerContactTypes';
import { mtoShipmentTypes, shipmentStatuses } from 'constants/shipments';
import { patchMTOServiceItemStatus, updateMTOShipmentStatus } from 'services/ghcApi';
import { useMoveTaskOrderQueries } from 'hooks/queries';
import LeftNav from 'components/LeftNav';
import ShipmentContainer from 'components/Office/ShipmentContainer';
import ShipmentHeading from 'components/Office/ShipmentHeading';
import ImportantShipmentDates from 'components/Office/ImportantShipmentDates';
import RequestedServiceItemsTable from 'components/Office/RequestedServiceItemsTable/RequestedServiceItemsTable';
import { MatchShape } from 'types/router';
import ShipmentAddresses from 'components/Office/ShipmentAddresses/ShipmentAddresses';
import RejectServiceItemModal from 'components/Office/RejectServiceItemModal/RejectServiceItemModal';
import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';
import { RequestShipmentCancellationModal } from 'components/Office/RequestShipmentCancellationModal/RequestShipmentCancellationModal';
import { shipmentSectionLabels } from 'content/shipments';
import { setFlashMessage } from 'store/flash/actions';
import FlashGridContainer from 'containers/FlashGridContainer/FlashGridContainer';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

function formatShipmentDate(shipmentDateString) {
  const dateObj = new Date(shipmentDateString);
  const weekday = new Intl.DateTimeFormat('en', { weekday: 'long' }).format(dateObj);
  const year = new Intl.DateTimeFormat('en', { year: 'numeric' }).format(dateObj);
  const month = new Intl.DateTimeFormat('en', { month: 'short' }).format(dateObj);
  const day = new Intl.DateTimeFormat('en', { day: '2-digit' }).format(dateObj);
  return `${weekday}, ${day} ${month} ${year}`;
}

function approvedFilter(shipment) {
  return shipment.status === shipmentStatuses.APPROVED || shipment.status === shipmentStatuses.CANCELLATION_REQUESTED;
}

export const MoveTaskOrder = ({ match, ...props }) => {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [isCancelModalVisible, setIsCancelModalVisible] = useState(false);
  const [selectedShipment, setSelectedShipment] = useState(undefined);
  const [selectedServiceItem, setSelectedServiceItem] = useState(undefined);
  const [sections, setSections] = useState([]);
  const [activeSection, setActiveSection] = useState('');
  const [unapprovedServiceItemsForShipment, setUnapprovedServiceItemsForShipment] = useState({});

  const { moveCode } = match.params;
  const { setUnapprovedShipmentCount, setUnapprovedServiceItemCount, setMessage } = props;

  const { orders = {}, moveTaskOrders, mtoShipments, mtoServiceItems, isLoading, isError } = useMoveTaskOrderQueries(
    moveCode,
  );

  const order = Object.values(orders)?.[0];
  const moveTaskOrder = Object.values(moveTaskOrders || {})?.[0];

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
        reason: item.reason,
        imgURL: '',
        description: item.description,
        itemDimensions: item.dimensions?.find((dimension) => dimension?.type === dimensionTypes.ITEM),
        crateDimensions: item.dimensions?.find((dimension) => dimension?.type === dimensionTypes.CRATE),
        firstCustomerContact: item.customerContacts?.find((contact) => contact?.type === customerContactTypes.FIRST),
        secondCustomerContact: item.customerContacts?.find((contact) => contact?.type === customerContactTypes.SECOND),
      };

      if (serviceItemsForShipment[`${newItem.mtoShipmentID}`]) {
        serviceItemsForShipment[`${newItem.mtoShipmentID}`].push(newItem);
      } else {
        serviceItemsForShipment[`${newItem.mtoShipmentID}`] = [newItem];
      }
    });
    return serviceItemsForShipment;
  }, [mtoServiceItems]);

  const [mutateMTOServiceItemStatus] = useMutation(patchMTOServiceItemStatus, {
    onSuccess: (data, variables) => {
      const newMTOServiceItem = data.mtoServiceItems[variables.mtoServiceItemID];
      mtoServiceItems[
        mtoServiceItems.find((serviceItem) => serviceItem.id === newMTOServiceItem.id)
      ] = newMTOServiceItem;
      queryCache.setQueryData([MTO_SERVICE_ITEMS, variables.moveTaskOrderId, false], mtoServiceItems);
      queryCache.invalidateQueries([MTO_SERVICE_ITEMS, variables.moveTaskOrderId]);
      setIsModalVisible(false);
      setSelectedServiceItem({});
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      // TODO: Handle error some how
      // RA Summary: eslint: no-console - System Information Leak: External
      // RA: The linter flags any use of console.
      // RA: This console displays an error message from unsuccessful mutation.
      // RA: TODO: As indicated, this error needs to be handled and needs further investigation.
      // RA: POAM story here: https://dp3.atlassian.net/browse/MB-5597
      // RA Developer Status: Known Issue
      // RA Validator Status: Known Issue
      // RA Modified Severity: CAT II
      // eslint-disable-next-line no-console
      console.log(errorMsg);
    },
  });

  const [mutateMTOShipmentStatus] = useMutation(updateMTOShipmentStatus, {
    onSuccess: (data, variables) => {
      const updatedMTOShipment = data.mtoShipments[variables.shipmentID];
      // Update mtoShipments with our updated status and set query data to match
      mtoShipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryCache.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      // InvalidateQuery tells other components using this data that they need to re-fetch
      // This allows the requestCancellation button to update immediately
      queryCache.invalidateQueries([MTO_SHIPMENTS, variables.moveTaskOrderID]);

      setIsCancelModalVisible(false);
      // Must set FlashMesage after hiding the modal, since FlashMessage will disappear when focus changes
      setMessage(
        `MSG_CANCEL_SUCCESS_${variables.shipmentID}`,
        'success',
        'The request to cancel that shipment has been sent to the movers.',
        '',
        true,
      );
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      // TODO: Handle error some how
      // RA Summary: eslint: no-console - System Information Leak: External
      // RA: The linter flags any use of console.
      // RA: This console displays an error message from unsuccessful mutation.
      // RA: TODO: As indicated, this error needs to be handled and needs further investigation.
      // RA: POAM story here: https://dp3.atlassian.net/browse/MB-5597
      // RA Developer Status: Known Issue
      // RA Validator Status: Known Issue
      // RA Modified Severity: CAT II
      // eslint-disable-next-line no-console
      console.log(errorMsg);
    },
  });

  const handleUpdateMTOShipmentStatus = (moveTaskOrderID, mtoShipmentID, eTag) => {
    mutateMTOShipmentStatus({
      moveTaskOrderID,
      shipmentID: mtoShipmentID,
      shipmentStatus: shipmentStatuses.CANCELLATION_REQUESTED,
      ifMatchETag: eTag,
    });
  };

  const handleUpdateMTOServiceItemStatus = (mtoServiceItemID, mtoShipmentID, status, rejectionReason) => {
    const mtoServiceItemForRequest = shipmentServiceItems[`${mtoShipmentID}`]?.find((s) => s.id === mtoServiceItemID);

    mutateMTOServiceItemStatus({
      moveTaskOrderId: moveTaskOrder.id,
      mtoServiceItemID,
      status,
      rejectionReason,
      ifMatchEtag: mtoServiceItemForRequest.eTag,
    });
  };

  useEffect(() => {
    let serviceItemCount = 0;
    const serviceItemsCountForShipment = {};
    mtoShipments?.forEach((mtoShipment) => {
      if (mtoShipment.status === shipmentStatuses.APPROVED) {
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
    }
  }, [mtoShipments, setUnapprovedShipmentCount]);

  useEffect(() => {
    const shipmentSections = [];
    mtoShipments?.forEach((shipment) => {
      if (
        shipment.status === shipmentStatuses.APPROVED ||
        shipment.status === shipmentStatuses.CANCELLATION_REQUESTED
      ) {
        shipmentSections.push({
          id: shipment.id,
          label: shipmentSectionLabels[`${shipment.shipmentType}`] || shipment.shipmentType,
        });
      }
    });
    setSections(shipmentSections);
  }, [mtoShipments]);

  const handleScroll = () => {
    const distanceFromTop = window.scrollY;
    let newActiveSection;

    sections.forEach((section) => {
      const sectionEl = document.querySelector(`#shipment-${section.id}`);
      if (sectionEl?.offsetTop <= distanceFromTop && sectionEl?.offsetTop + sectionEl?.offsetHeight > distanceFromTop) {
        newActiveSection = section.id;
      }
    });

    if (activeSection !== newActiveSection) {
      setActiveSection(newActiveSection);
    }
  };

  useEffect(() => {
    // attach scroll listener
    window.addEventListener('scroll', handleScroll);

    // remove scroll listener
    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  });

  const handleShowRejectionDialog = (mtoServiceItemID, mtoShipmentID) => {
    const serviceItem = shipmentServiceItems[`${mtoShipmentID}`]?.find((item) => item.id === mtoServiceItemID);
    setSelectedServiceItem(serviceItem);
    setIsModalVisible(true);
  };

  const handleShowCancellationModal = (mtoShipment) => {
    setSelectedShipment(mtoShipment);
    setIsCancelModalVisible(true);
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  if (moveTaskOrder.status === MOVE_STATUSES.SUBMITTED || !mtoShipments.some(approvedFilter)) {
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

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <LeftNav className={styles.sidebar}>
          {sections.map((s) => {
            const classes = classnames({ active: s.id === activeSection });
            return (
              <a key={`sidenav_${s.id}`} href={`#shipment-${s.id}`} className={classes}>
                {s.label}{' '}
                {unapprovedServiceItemsForShipment[`${s.id}`] > 0 && (
                  <Tag>{unapprovedServiceItemsForShipment[`${s.id}`]}</Tag>
                )}
              </a>
            );
          })}
        </LeftNav>
        <FlashGridContainer className={styles.gridContainer} data-testid="too-shipment-container">
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
          <div className={styles.pageHeader}>
            <h1>Move task order</h1>
            <div className={styles.pageHeaderDetails}>
              <h6>MTO Reference ID #{moveTaskOrder?.referenceId}</h6>
              <h6>Contract #1234567890</h6> {/* TODO - need this value from the API */}
            </div>
          </div>
          {mtoShipments.map((mtoShipment) => {
            if (
              mtoShipment.status !== shipmentStatuses.APPROVED &&
              mtoShipment.status !== shipmentStatuses.CANCELLATION_REQUESTED
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
            // eslint-disable-next-line camelcase
            const dutyStationPostal = { postal_code: order.destinationDutyStation.address.postal_code };
            const { pickupAddress, destinationAddress } = mtoShipment;
            const formattedScheduledPickup = formatShipmentDate(mtoShipment.scheduledPickupDate);

            return (
              <ShipmentContainer
                id={`shipment-${mtoShipment.id}`}
                key={mtoShipment.id}
                shipmentType={mtoShipment.shipmentType}
                className={styles.shipmentCard}
              >
                <ShipmentHeading
                  shipmentInfo={{
                    shipmentID: mtoShipment.id,
                    shipmentType: mtoShipmentTypes[mtoShipment.shipmentType],
                    originCity: pickupAddress?.city,
                    originState: pickupAddress?.state,
                    originPostalCode: pickupAddress?.postal_code,
                    destinationAddress: destinationAddress || dutyStationPostal,
                    scheduledPickupDate: formattedScheduledPickup,
                    shipmentStatus: mtoShipment.status,
                    ifMatchEtag: mtoShipment.eTag,
                    moveTaskOrderID: mtoShipment.moveTaskOrderID,
                  }}
                  handleShowCancellationModal={handleShowCancellationModal}
                />
                <ImportantShipmentDates
                  requestedPickupDate={formatShipmentDate(mtoShipment.requestedPickupDate)}
                  scheduledPickupDate={formattedScheduledPickup}
                />
                <ShipmentAddresses
                  pickupAddress={pickupAddress}
                  destinationAddress={destinationAddress || dutyStationPostal}
                  originDutyStation={order.originDutyStation?.address}
                  destinationDutyStation={order.destinationDutyStation?.address}
                />
                <ShipmentWeightDetails
                  estimatedWeight={mtoShipment.primeEstimatedWeight}
                  actualWeight={mtoShipment.primeActualWeight}
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
  match: MatchShape.isRequired,
  setUnapprovedShipmentCount: func.isRequired,
  setUnapprovedServiceItemCount: func.isRequired,
  setMessage: func.isRequired,
};

const mapDispatchToProps = {
  setMessage: setFlashMessage,
};

export default withRouter(connect(() => ({}), mapDispatchToProps)(MoveTaskOrder));
