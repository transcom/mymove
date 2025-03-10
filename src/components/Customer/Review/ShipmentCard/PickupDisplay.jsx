import React, { useEffect, useState } from 'react';
import { string, shape } from 'prop-types';

import { AddressShape } from '../../../../types/address';

import styles from './ShipmentCard.module.scss';

import { formatCustomerDate } from 'utils/formatters';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { formatCustomerDestination } from 'utils/shipmentDisplay';

const PickupDisplay = ({
  pickupLocation,
  secondaryPickupAddress,
  tertiaryPickupAddress,
  releasingAgent,
  requestedPickupDate,
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
        <dt>Requested pickup date</dt>
        <dd>{formatCustomerDate(requestedPickupDate)}</dd>
      </div>
      {pickupLocation && (
        <div className={styles.row}>
          <dt>Pickup Address</dt>
          <dd>{formatCustomerDestination(pickupLocation)}</dd>
        </div>
      )}
      {secondaryPickupAddress && (
        <div className={styles.row}>
          <dt>Second Pickup Address</dt>
          <dd>{formatCustomerDestination(secondaryPickupAddress)}</dd>
        </div>
      )}
      {isTertiaryAddressEnabled && tertiaryPickupAddress && secondaryPickupAddress && (
        <div className={styles.row}>
          <dt>Third Pickup Address</dt>
          <dd>{formatCustomerDestination(tertiaryPickupAddress)}</dd>
        </div>
      )}
      {releasingAgent && (
        <div className={styles.row}>
          <dt>Releasing agent</dt>
          <dd>
            {(releasingAgent.firstName || releasingAgent.lastName) && (
              <>
                {releasingAgent.firstName} {releasingAgent.lastName} <br />
              </>
            )}
            {releasingAgent.phone && (
              <>
                {releasingAgent.phone} <br />
              </>
            )}
            {releasingAgent.email}
          </dd>
        </div>
      )}
    </>
  );
};

PickupDisplay.propTypes = {
  requestedPickupDate: string.isRequired,
  pickupLocation: AddressShape.isRequired,
  secondaryPickupAddress: AddressShape,
  releasingAgent: shape({
    firstName: string,
    lastName: string,
    phone: string,
    email: string,
  }),
};

PickupDisplay.defaultProps = {
  secondaryPickupAddress: null,
  releasingAgent: null,
};

export default PickupDisplay;
