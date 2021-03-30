import React from 'react';
import classNames from 'classnames';
import { PropTypes } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { AddressShape } from '../../types/address';

import styles from './shipmentHeading.module.scss';

import { shipmentStatuses } from 'constants/shipments';

function formatDestinationAddress(address) {
  if (address.city) {
    // eslint-disable-next-line camelcase
    return `${address.city}, ${address.state} ${address.postal_code}`;
  }
  // eslint-disable-next-line camelcase
  return `${address.postal_code}`;
}

function ShipmentHeading({ shipmentInfo, handleShowCancellationModal }) {
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
          onClick={() => handleShowCancellationModal(shipmentInfo)}
          unstyled
          disabled={shipmentInfo.shipmentStatus === shipmentStatuses.CANCELLATION_REQUESTED}
        >
          {shipmentInfo.shipmentStatus === shipmentStatuses.CANCELLATION_REQUESTED
            ? 'Cancellation Requested'
            : 'Request Cancellation'}
        </Button>
      </div>
    </div>
  );
}

ShipmentHeading.propTypes = {
  shipmentInfo: PropTypes.shape({
    shipmentID: PropTypes.string.isRequired,
    shipmentType: PropTypes.string.isRequired,
    originCity: PropTypes.string.isRequired,
    originState: PropTypes.string.isRequired,
    originPostalCode: PropTypes.string.isRequired,
    destinationAddress: AddressShape,
    scheduledPickupDate: PropTypes.string.isRequired,
    shipmentStatus: PropTypes.string.isRequired,
    ifMatchEtag: PropTypes.string.isRequired,
    moveTaskOrderID: PropTypes.string.isRequired,
  }).isRequired,
  handleShowCancellationModal: PropTypes.func.isRequired,
};

export default ShipmentHeading;
