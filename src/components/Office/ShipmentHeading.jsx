import React from 'react';
import classNames from 'classnames/bind';
import PropTypes from 'prop-types';

import styles from './shipmentHeading.module.scss';

const cx = classNames.bind(styles);

function ShipmentHeading({ shipmentInfo }) {
  return (
    <div className={cx('shipment-heading')}>
      <h3 data-testid="office-shipment-heading-h3">{shipmentInfo.shipmentType}</h3>
      <small>
        {`${shipmentInfo.originCity} ${shipmentInfo.originState} ${shipmentInfo.originPostalCode} to
    ${shipmentInfo.destinationCity} ${shipmentInfo.destinationState} ${shipmentInfo.destinationPostalCode}
    on ${shipmentInfo.scheduledPickupDate}`}
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
    destinationCity: PropTypes.string.isRequired,
    destinationState: PropTypes.string.isRequired,
    destinationPostalCode: PropTypes.string.isRequired,
    scheduledPickupDate: PropTypes.string.isRequired,
  }).isRequired,
};

export default ShipmentHeading;
