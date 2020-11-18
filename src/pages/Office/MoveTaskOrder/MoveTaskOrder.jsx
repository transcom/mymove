import React, { useState } from 'react';
import { withRouter } from 'react-router-dom';
import { get } from 'lodash';
import { GridContainer } from '@trussworks/react-uswds';
import { queryCache, useMutation } from 'react-query';

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

function formatShipmentType(shipmentType) {
  if (shipmentType === 'HHG') {
    return 'Household Goods';
  }
  return shipmentType;
}

function formatShipmentDate(shipmentDateString) {
  const dateObj = new Date(shipmentDateString);
  const weekday = new Intl.DateTimeFormat('en', { weekday: 'long' }).format(dateObj);
  const year = new Intl.DateTimeFormat('en', { year: 'numeric' }).format(dateObj);
  const month = new Intl.DateTimeFormat('en', { month: 'short' }).format(dateObj);
  const day = new Intl.DateTimeFormat('en', { day: '2-digit' }).format(dateObj);
  return `${weekday}, ${day} ${month} ${year}`;
}

export const MoveTaskOrder = ({ match }) => {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [selectedServiceItem, setSelectedServiceItem] = useState(undefined);

  const { moveOrderId } = match.params;

  // TODO - Do something with moveOrder and moveTaskOrder?
  const {
    moveOrders = {},
    moveTaskOrders,
    mtoShipments,
    mtoServiceItems,
    isLoading,
    isError,
  } = useMoveTaskOrderQueries(moveOrderId);

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
      // RA: As indicated, this error needs to be handled and needs further investigation.
      // RA Developer Status: Known Issue
      // RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
      // RA Validator: jneuner@mitre.org
      // RA Modified Severity:
      // eslint-disable-next-line no-console
      console.log(errorMsg);
    },
  });

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
      <div>
        <p>This Move does not have any approved shipments yet.</p>
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
                  shipmentType: formatShipmentType(mtoShipment.shipmentType),
                  originCity: get(mtoShipment.pickupAddress, 'city'),
                  originState: get(mtoShipment.pickupAddress, 'state'),
                  originPostalCode: get(mtoShipment.pickupAddress, 'postal_code'),
                  destinationAddress: mtoShipment.destinationAddress || dutyStationPostal,
                  scheduledPickupDate: formatShipmentDate(mtoShipment.scheduledPickupDate),
                }}
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
};

export default withRouter(MoveTaskOrder);
