import React, { useState, useEffect } from 'react';
import { withRouter } from 'react-router-dom';
import { get } from 'lodash';
import { GridContainer, Alert } from '@trussworks/react-uswds';
import { queryCache, useMutation } from 'react-query';
import { func } from 'prop-types';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import { MTO_SERVICE_ITEMS } from 'constants/queryKeys';
import ShipmentContainer from 'components/Office/ShipmentContainer';
import ShipmentHeading from 'components/Office/ShipmentHeading';
import ImportantShipmentDates from 'components/Office/ImportantShipmentDates';
import RequestedServiceItemsTable from 'components/Office/RequestedServiceItemsTable/RequestedServiceItemsTable';
import { useMoveTaskOrderQueries } from 'hooks/queries';
import { MatchShape } from 'types/router';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import ShipmentAddresses from 'components/Office/ShipmentAddresses/ShipmentAddresses';
import RejectServiceItemModal from 'components/Office/RejectServiceItemModal/RejectServiceItemModal';
import { SERVICE_ITEM_STATUS } from 'shared/constants';
import { patchMTOServiceItemStatus } from 'services/ghcApi';
import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';
import dimensionTypes from 'constants/dimensionTypes';
import customerContactTypes from 'constants/customerContactTypes';
import { mtoShipmentTypes } from 'constants/shipments';

function formatShipmentDate(shipmentDateString) {
  const dateObj = new Date(shipmentDateString);
  const weekday = new Intl.DateTimeFormat('en', { weekday: 'long' }).format(dateObj);
  const year = new Intl.DateTimeFormat('en', { year: 'numeric' }).format(dateObj);
  const month = new Intl.DateTimeFormat('en', { month: 'short' }).format(dateObj);
  const day = new Intl.DateTimeFormat('en', { day: '2-digit' }).format(dateObj);
  return `${weekday}, ${day} ${month} ${year}`;
}

export const MoveTaskOrder = ({ match, ...props }) => {
  // Using hooks to illustrate disabled button state for shipment cancellation
  // This will be modified once the modal is hooked up, as the button will only
  // be used to trigger the modal.
  const [mockShipmentStatus, setMockShipmentStatus] = useState(undefined);
  const [currentAlert, setCurrentAlert] = useState(undefined);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [selectedServiceItem, setSelectedServiceItem] = useState(undefined);

  const { moveCode } = match.params;
  const { setUnapprovedShipmentCount } = props;

  // TODO - Do something with moveOrder and moveTaskOrder?
  const {
    moveOrders = {},
    moveTaskOrders,
    mtoShipments,
    mtoServiceItems,
    isLoading,
    isError,
  } = useMoveTaskOrderQueries(moveCode);

  let mtoServiceItemsArr;
  if (mtoServiceItems) {
    mtoServiceItemsArr = Object.values(mtoServiceItems);
  }

  const moveOrder = Object.values(moveOrders)?.[0];
  let moveTaskOrder;
  if (moveTaskOrders) {
    moveTaskOrder = Object.values(moveTaskOrders)?.[0];
  }

  const [mutateMTOServiceItemStatus] = useMutation(patchMTOServiceItemStatus, {
    onSuccess: (data, variables) => {
      const newMTOServiceItem = data.mtoServiceItems[variables.mtoServiceItemID];
      queryCache.setQueryData([MTO_SERVICE_ITEMS, variables.mtoServiceItemID], {
        mtoServiceItems: {
          ...mtoServiceItems,
          [`${variables.mtoServiceItemID}`]: newMTOServiceItem,
        },
      });
      queryCache.invalidateQueries(MTO_SERVICE_ITEMS);
      setSelectedServiceItem({});
      setIsModalVisible(false);
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

  const handleUpdateMTOShipmentStatus = (mtoShipmentID, status) => {
    setCurrentAlert({
      type: 'success',
      msg: 'The request to cancel that shipment has been sent to the movers.',
    });
    setMockShipmentStatus({
      id: mtoShipmentID,
      status,
    });
    // TODO mutateMTOShipmentStatus(); to implement updateMTOShipmentStatus endpoint
  };

  const handleUpdateMTOServiceItemStatus = (mtoServiceItemID, status, rejectionReason) => {
    const mtoServiceItemForRequest = mtoServiceItemsArr.find((s) => s.id === mtoServiceItemID);

    mutateMTOServiceItemStatus({
      moveTaskOrderId: moveTaskOrder.id,
      mtoServiceItemID,
      status,
      rejectionReason,
      ifMatchEtag: mtoServiceItemForRequest.eTag,
    });
  };

  useEffect(() => {
    const shipmentCount = mtoShipments
      ? Object.values(mtoShipments).filter((shipment) => shipment.status === 'SUBMITTED').length
      : 0;
    setUnapprovedShipmentCount(shipmentCount);
  }, [mtoShipments, setUnapprovedShipmentCount]);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const serviceItems = mtoServiceItemsArr?.map((item) => {
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
    return newItem;
  });

  const handleShowRejectionDialog = (mtoServiceItemID) => {
    const serviceItem = serviceItems?.find((item) => item.id === mtoServiceItemID);
    setSelectedServiceItem(serviceItem);
    setIsModalVisible(true);
  };

  const approved = (shipment) => shipment.status === 'APPROVED';
  const mtoShipmentsArr = Object.values(mtoShipments);

  if (!mtoShipmentsArr.some(approved)) {
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
      <GridContainer className={styles.gridContainer} data-testid="too-shipment-container">
        {isModalVisible && (
          <RejectServiceItemModal
            serviceItem={selectedServiceItem}
            onSubmit={handleUpdateMTOServiceItemStatus}
            onClose={setIsModalVisible}
          />
        )}
        {currentAlert && (
          <Alert slim type={currentAlert.type}>
            {currentAlert.msg}
          </Alert>
        )}

        <div className={styles.pageHeader}>
          <h1>Move task order</h1>
          <div className={styles.pageHeaderDetails}>
            <h6>MTO Reference ID #{moveTaskOrder?.referenceId}</h6>
            <h6>Contract #1234567890</h6> {/* TODO - need this value from the API */}
          </div>
        </div>

        {mtoShipmentsArr.map((mtoShipment) => {
          if (mtoShipment.status !== 'APPROVED') {
            return false;
          }
          // This code mocks a "CANCELLATION_REQUESTED" status change on a shipment so we can test that behavior
          const mockStatus =
            mockShipmentStatus && mockShipmentStatus.id === mtoShipment.id
              ? mockShipmentStatus.status
              : mtoShipment.status;
          const serviceItemsForShipment = serviceItems.filter((item) => item.mtoShipmentID === mtoShipment.id);
          const requestedServiceItems = serviceItemsForShipment.filter(
            (item) => item.status === SERVICE_ITEM_STATUS.SUBMITTED,
          );
          const approvedServiceItems = serviceItemsForShipment.filter(
            (item) => item.status === SERVICE_ITEM_STATUS.APPROVED,
          );
          const rejectedServiceItems = serviceItemsForShipment.filter(
            (item) => item.status === SERVICE_ITEM_STATUS.REJECTED,
          );
          // eslint-disable-next-line camelcase
          const dutyStationPostal = { postal_code: moveOrder.destinationDutyStation.address.postal_code };
          return (
            <ShipmentContainer
              key={mtoShipment.id}
              shipmentType={mtoShipment.shipmentType}
              className={styles.shipmentCard}
            >
              <ShipmentHeading
                key={mtoShipment.id}
                shipmentInfo={{
                  shipmentID: mtoShipment.id,
                  shipmentType: mtoShipmentTypes[mtoShipment.shipmentType],
                  originCity: get(mtoShipment.pickupAddress, 'city'),
                  originState: get(mtoShipment.pickupAddress, 'state'),
                  originPostalCode: get(mtoShipment.pickupAddress, 'postal_code'),
                  destinationAddress: mtoShipment.destinationAddress || dutyStationPostal,
                  scheduledPickupDate: formatShipmentDate(mtoShipment.scheduledPickupDate),
                  shipmentStatus: mockStatus,
                }}
                handleUpdateMTOShipmentStatus={handleUpdateMTOShipmentStatus}
              />
              <ImportantShipmentDates
                requestedPickupDate={formatShipmentDate(mtoShipment.requestedPickupDate)}
                scheduledPickupDate={formatShipmentDate(mtoShipment.scheduledPickupDate)}
              />
              <ShipmentAddresses
                pickupAddress={mtoShipment?.pickupAddress}
                destinationAddress={mtoShipment?.destinationAddress || dutyStationPostal}
                originDutyStation={moveOrder?.originDutyStation?.address}
                destinationDutyStation={moveOrder?.destinationDutyStation?.address}
              />
              <ShipmentWeightDetails
                estimatedWeight={mtoShipment?.primeEstimatedWeight}
                actualWeight={mtoShipment?.primeActualWeight}
              />
              {requestedServiceItems?.length > 0 && (
                <RequestedServiceItemsTable
                  serviceItems={requestedServiceItems}
                  handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                  handleShowRejectionDialog={handleShowRejectionDialog}
                  statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
                />
              )}
              {approvedServiceItems?.length > 0 && (
                <RequestedServiceItemsTable
                  serviceItems={approvedServiceItems}
                  handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                  handleShowRejectionDialog={handleShowRejectionDialog}
                  statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
                />
              )}
              {rejectedServiceItems?.length > 0 && (
                <RequestedServiceItemsTable
                  serviceItems={rejectedServiceItems}
                  handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                  handleShowRejectionDialog={handleShowRejectionDialog}
                  statusForTableType={SERVICE_ITEM_STATUS.REJECTED}
                />
              )}
            </ShipmentContainer>
          );
        })}
      </GridContainer>
    </div>
  );
};

MoveTaskOrder.propTypes = {
  match: MatchShape.isRequired,
  setUnapprovedShipmentCount: func.isRequired,
};

export default withRouter(MoveTaskOrder);
