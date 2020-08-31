import React from 'react';
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
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      // TODO: Handle error some how
      // eslint-disable-next-line no-console
      console.log(errorMsg);
    },
  });

  const handleUpdateMTOServiceItemStatus = (mtoServiceItemID, status) => {
    const mtoServiceItemForRequest = mtoServiceItemsArr.find((s) => s.id === mtoServiceItemID);

    mutateMTOServiceItemStatus({
      moveTaskOrderId: moveTaskOrder.id,
      mtoServiceItemID,
      status,
      ifMatchEtag: mtoServiceItemForRequest.eTag,
    });
  };

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const serviceItems = mtoServiceItems.map((item) => {
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

  return (
    <div className={styles.tabContent}>
      <GridContainer className={styles.gridContainer} data-testid="too-shipment-container">
        <div className={styles.pageHeader}>
          <h1>Move task order</h1>
          <div className={styles.pageHeaderDetails}>
            <h6>MTO Reference ID #{moveTaskOrder?.referenceId}</h6>
            <h6>Contract #1234567890</h6> {/* TODO - need this value from the API */}
          </div>
        </div>

        {mtoShipments.map((mtoShipment) => {
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
                  destinationCity: get(mtoShipment.destinationAddress, 'city'),
                  destinationState: get(mtoShipment.destinationAddress, 'state'),
                  destinationPostalCode: get(mtoShipment.destinationAddress, 'postal_code'),
                  scheduledPickupDate: formatShipmentDate(mtoShipment.scheduledPickupDate),
                }}
              />
              <ImportantShipmentDates
                requestedPickupDate={formatShipmentDate(mtoShipment.requestedPickupDate)}
                scheduledPickupDate={formatShipmentDate(mtoShipment.scheduledPickupDate)}
              />
              <ShipmentAddresses
                pickupAddress={mtoShipment?.pickupAddress}
                destinationAddress={mtoShipment?.destinationAddress}
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
                  statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
                />
              )}
              {approvedServiceItems?.length > 0 && (
                <RequestedServiceItemsTable
                  serviceItems={approvedServiceItems}
                  handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
                  statusForTableType={SERVICE_ITEM_STATUS.APPROVED}
                />
              )}
              {rejectedServiceItems?.length > 0 && (
                <RequestedServiceItemsTable
                  serviceItems={rejectedServiceItems}
                  handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
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
