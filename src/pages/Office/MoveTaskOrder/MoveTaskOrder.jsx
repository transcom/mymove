import React from 'react';
import { withRouter } from 'react-router-dom';
import { get, map } from 'lodash';
import { GridContainer } from '@trussworks/react-uswds';

import styles from '../MoveDetails/MoveDetails.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer';
import ShipmentHeading from 'components/Office/ShipmentHeading';
import ImportantShipmentDates from 'components/Office/ImportantShipmentDates';
import RequestedServiceItemsTable from 'components/Office/RequestedServiceItemsTable';
import { useMoveTaskOrderQueries } from 'hooks/queries';
import { MatchShape } from 'types/router';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

function formatShipmentType(shipmentType) {
  if (shipmentType === 'HHG') {
    return 'Household Goods';
  }
  return shipmentType;
}

function formatShipmentDate(shipmentDateString) {
  const dateObj = new Date(shipmentDateString);
  const year = new Intl.DateTimeFormat('en', { year: 'numeric' }).format(dateObj);
  const month = new Intl.DateTimeFormat('en', { month: 'short' }).format(dateObj);
  const day = new Intl.DateTimeFormat('en', { day: '2-digit' }).format(dateObj);
  return `${day} ${month} ${year}`;
}

export const MoveTaskOrder = ({ match }) => {
  const { moveOrderId } = match.params;

  // TODO - Do something with moveOrder and moveTaskOrder?
  const { mtoShipments, mtoServiceItems, isLoading, isError } = useMoveTaskOrderQueries(moveOrderId);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const serviceItems = map(mtoServiceItems, (item) => {
    const newItem = { ...item };
    newItem.serviceItem = item.reServiceName;
    newItem.details = { text: { ZIP: item.pickupPostalCode, Reason: item.reason }, imgURL: '' };

    return newItem;
  });

  return (
    <div className={styles.MoveDetails}>
      <GridContainer className={styles.gridContainer} data-testid="too-shipment-container">
        <h1>Move task order</h1>

        {map(mtoShipments, (mtoShipment) => {
          const serviceItemsForShipment = serviceItems.filter((item) => item.mtoShipmentID === mtoShipment.id);

          return (
            <ShipmentContainer>
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
              <RequestedServiceItemsTable serviceItems={serviceItemsForShipment} />
            </ShipmentContainer>
          );
        })}
      </GridContainer>
    </div>
  );
};

MoveTaskOrder.propTypes = {
  // history: HistoryShape.isRequired,
  match: MatchShape.isRequired,
};

export default withRouter(MoveTaskOrder);
