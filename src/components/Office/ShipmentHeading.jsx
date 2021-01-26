import React from 'react';
import classNames from 'classnames';
import PropTypes from 'prop-types';

import { AddressShape } from '../../types/address';

import styles from './shipmentHeading.module.scss';

import { mtoShipmentTypes } from 'constants/shipments';

function formatDestinationAddress(address) {
  if (address.city) {
    // eslint-disable-next-line camelcase
    return `${address.city}, ${address.state} ${address.postal_code}`;
  }
  // eslint-disable-next-line camelcase
  return `${address.postal_code}`;
}

function ShipmentHeading({ shipmentInfo }) {
  return (
    <div className={classNames(styles.shipmentHeading, 'shipment-heading')}>
      <h3 data-testid="office-shipment-heading-h3">{mtoShipmentTypes[shipmentInfo.shipmentType]}</h3>
      <small>
        {`${shipmentInfo.originCity}, ${shipmentInfo.originState} ${shipmentInfo.originPostalCode} to
        ${formatDestinationAddress(shipmentInfo.destinationAddress)} on ${shipmentInfo.scheduledPickupDate}`}
      </small>
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
  }).isRequired,
};

export default ShipmentHeading;
