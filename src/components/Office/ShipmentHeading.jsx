import React, { useState } from 'react';
import classNames from 'classnames';
import { PropTypes } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { AddressShape } from '../../types/address';

import styles from './shipmentHeading.module.scss';

import { MTO_SHIPMENT_STATUSES } from 'shared/constants';

function formatDestinationAddress(address) {
  if (address.city) {
    // eslint-disable-next-line camelcase
    return `${address.city}, ${address.state} ${address.postal_code}`;
  }
  // eslint-disable-next-line camelcase
  return `${address.postal_code}`;
}

function ShipmentHeading({ shipmentInfo }) {
  // Using hooks to illustrate disabled button state
  // This will be modified once the modal is hooked up, as the button will only
  // be used to trigger the modal.
  const [shipmentStatus, setShipmentStatus] = useState(shipmentInfo.shipmentStatus);

  return (
    <div className={classNames(styles.shipmentHeading, 'shipment-heading')}>
      <h3 data-testid="office-shipment-heading-h3">{shipmentInfo.shipmentType}</h3>
      <div className={styles.row}>
        <small>
          {`${shipmentInfo.originCity}, ${shipmentInfo.originState} ${shipmentInfo.originPostalCode} to
        ${formatDestinationAddress(shipmentInfo.destinationAddress)} on ${shipmentInfo.scheduledPickupDate}`}
        </small>
        <Button
          type="button"
          onClick={() => setShipmentStatus(MTO_SHIPMENT_STATUSES.CANCELLATION_REQUESTED)}
          unstyled
          disabled={shipmentStatus === MTO_SHIPMENT_STATUSES.CANCELLATION_REQUESTED}
        >
          <small>
            {shipmentStatus === MTO_SHIPMENT_STATUSES.CANCELLATION_REQUESTED
              ? 'Cancellation Requested'
              : 'Request Cancellation'}
          </small>
        </Button>
      </div>
    </div>
  );
}

ShipmentHeading.propTypes = {
  shipmentInfo: PropTypes.shape({
    shipmentType: PropTypes.string.isRequired,
    originCity: PropTypes.string.isRequired,
    originState: PropTypes.string.isRequired,
    originPostalCode: PropTypes.string.isRequired,
    destinationAddress: AddressShape,
    scheduledPickupDate: PropTypes.string.isRequired,
    shipmentStatus: PropTypes.string.isRequired,
  }).isRequired,
};

export default ShipmentHeading;
