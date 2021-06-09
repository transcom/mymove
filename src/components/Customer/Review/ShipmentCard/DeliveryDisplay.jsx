import React from 'react';
import { string, shape } from 'prop-types';

import { AddressShape } from '../../../../types/address';

import styles from './ShipmentCard.module.scss';

import { formatCustomerDate } from 'utils/formatters';
import { formatCustomerDestination } from 'utils/shipmentDisplay';

const DeliveryDisplay = ({
  destinationLocation,
  destinationZIP,
  secondaryDestinationLocation,
  receivingAgent,
  requestedDeliveryDate,
}) => {
  return (
    <div>
      <div className={styles.row}>
        <dt>Requested delivery date</dt>
        <dd>{formatCustomerDate(requestedDeliveryDate)}</dd>
      </div>
      <div className={styles.row}>
        <dt>Destination</dt>
        <dd>{formatCustomerDestination(destinationLocation, destinationZIP)}</dd>
      </div>
      {secondaryDestinationLocation && (
        <div className={styles.row}>
          <dt>Second Destination</dt>
          <dd>{formatCustomerDestination(secondaryDestinationLocation)}</dd>
        </div>
      )}

      {receivingAgent && (
        <div className={styles.row}>
          <dt>Receiving agent</dt>
          <dd>
            {(receivingAgent.firstName || receivingAgent.lastName) && (
              <>
                {receivingAgent.firstName} {receivingAgent.lastName} <br />
              </>
            )}
            {receivingAgent.phone && (
              <>
                {receivingAgent.phone} <br />
              </>
            )}
            {receivingAgent.email}
          </dd>
        </div>
      )}
    </div>
  );
};

DeliveryDisplay.propTypes = {
  destinationLocation: AddressShape,
  secondaryDestinationLocation: AddressShape,
  requestedDeliveryDate: string.isRequired,
  destinationZIP: string.isRequired,
  receivingAgent: shape({
    firstName: string,
    lastName: string,
    phone: string,
    email: string,
  }),
};

DeliveryDisplay.defaultProps = {
  destinationLocation: null,
  secondaryDestinationLocation: null,
  receivingAgent: null,
};

export default DeliveryDisplay;
