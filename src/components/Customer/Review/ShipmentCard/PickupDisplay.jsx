import React from 'react';
import { string, shape } from 'prop-types';

import { AddressShape } from '../../../../types/address';

import styles from './ShipmentCard.module.scss';

import { formatCustomerDate } from 'utils/formatters';

const PickupDisplay = ({ pickupLocation, secondaryPickupAddress, releasingAgent, requestedPickupDate }) => {
  return (
    <div>
      <div className={styles.row}>
        <dt>Requested pickup date</dt>
        <dd>{formatCustomerDate(requestedPickupDate)}</dd>
      </div>
      {pickupLocation && (
        <div className={styles.row}>
          <dt>Pickup location</dt>
          <dd>
            {pickupLocation.street_address_1} {pickupLocation.street_address_2}
            <br />
            {pickupLocation.city}, {pickupLocation.state} {pickupLocation.postal_code}
          </dd>
        </div>
      )}
      {secondaryPickupAddress && (
        <div className={styles.row}>
          <dt>Second pickup location</dt>
          <dd>
            {secondaryPickupAddress.street_address_1} {secondaryPickupAddress.street_address_2}
            <br />
            {secondaryPickupAddress.city}, {secondaryPickupAddress.state} {secondaryPickupAddress.postal_code}
          </dd>
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
    </div>
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
