import React from 'react';
import { withRouter } from 'react-router-dom';
import { get, map } from 'lodash';
import { GridContainer } from '@trussworks/react-uswds';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer';
import ShipmentHeading from 'components/Office/ShipmentHeading';
import ImportantShipmentDates from 'components/Office/ImportantShipmentDates';
import RequestedServiceItemsTable from 'components/Office/RequestedServiceItemsTable';
import { useMoveTaskOrderQueries } from 'hooks/queries';
import { MatchShape } from 'types/router';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import ShipmentAddresses from 'components/Office/ShipmentAddresses/ShipmentAddresses';

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

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const moveOrder = Object.values(moveOrders)?.[0];
  const moveTaskOrder = Object.values(moveTaskOrders)?.[0];

  const serviceItems = map(mtoServiceItems, (item) => {
    const newItem = { ...item };
    newItem.serviceItem = item.reServiceName;
    newItem.details = { text: { ZIP: item.pickupPostalCode, Reason: item.reason }, imgURL: '' };

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

        {map(mtoShipments, (mtoShipment) => {
          const serviceItemsForShipment = serviceItems.filter((item) => item.mtoShipmentID === mtoShipment.id);
          return (
            <ShipmentContainer shipmentType={mtoShipment.shipmentType} className={styles.shipmentCard}>
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
                // eslint-disable-next-line react/prop-types
                originDutyStation={moveOrder?.originDutyStation?.address}
                // eslint-disable-next-line react/prop-types
                destinationDutyStation={moveOrder?.destinationDutyStation?.address}
              />
              <RequestedServiceItemsTable serviceItems={serviceItemsForShipment} />
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
