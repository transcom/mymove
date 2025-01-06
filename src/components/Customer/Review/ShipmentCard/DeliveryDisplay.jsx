import { React, useEffect, useState } from 'react';
import { string, shape } from 'prop-types';

import { AddressShape } from '../../../../types/address';

import styles from './ShipmentCard.module.scss';

import { formatCustomerDate } from 'utils/formatters';
import { formatCustomerDestination } from 'utils/shipmentDisplay';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const DeliveryDisplay = ({
  destinationLocation,
  destinationZIP,
  secondaryDeliveryAddress,
  tertiaryDeliveryAddress,
  receivingAgent,
  requestedDeliveryDate,
}) => {
  const [isTertiaryAddressEnabled, setIsTertiaryAddressEnabled] = useState(false);
  useEffect(() => {
    const fetchData = async () => {
      isBooleanFlagEnabled('third_address_available').then((enabled) => {
        setIsTertiaryAddressEnabled(enabled);
      });
    };
    fetchData();
  }, []);

  return (
    <>
      <div className={styles.row}>
        <dt>Requested delivery date</dt>
        <dd>{formatCustomerDate(requestedDeliveryDate)}</dd>
      </div>
      <div className={styles.row}>
        <dt>Delivery Address</dt>
        <dd>{formatCustomerDestination(destinationLocation, destinationZIP)}</dd>
      </div>
      {secondaryDeliveryAddress && (
        <div className={styles.row}>
          <dt>Second Delivery Address</dt>
          <dd>{formatCustomerDestination(secondaryDeliveryAddress)}</dd>
        </div>
      )}
      {isTertiaryAddressEnabled && secondaryDeliveryAddress && tertiaryDeliveryAddress && (
        <div className={styles.row}>
          <dt>Third Delivery Address</dt>
          <dd>{formatCustomerDestination(tertiaryDeliveryAddress)}</dd>
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
    </>
  );
};

DeliveryDisplay.propTypes = {
  destinationLocation: AddressShape,
  secondaryDeliveryAddress: AddressShape,
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
  secondaryDeliveryAddress: null,
  receivingAgent: null,
};

export default DeliveryDisplay;
