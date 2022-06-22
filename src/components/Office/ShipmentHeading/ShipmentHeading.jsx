import React from 'react';
import classNames from 'classnames';
import { PropTypes } from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';

import { AddressShape } from '../../../types/address';

import styles from './shipmentHeading.module.scss';

import { shipmentStatuses } from 'constants/shipments';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';

function formatDestinationAddress(address) {
  if (address.city) {
    // eslint-disable-next-line camelcase
    return `${address.city}, ${address.state} ${address.postalCode}`;
  }
  // eslint-disable-next-line camelcase
  return `${address.postalCode}`;
}

function ShipmentHeading({ shipmentInfo, handleShowCancellationModal }) {
  const { shipmentStatus } = shipmentInfo;
  // cancelation modal is visible if shipment is not already canceled, AND if shipment cancellation hasn't already been requested
  const isCancelModalVisible = shipmentStatus !== shipmentStatuses.CANCELED || shipmentStatuses.CANCELLATION_REQUESTED;
  const isCancellationRequested = shipmentStatus === shipmentStatuses.CANCELLATION_REQUESTED;

  return (
    <div className={classNames(styles.shipmentHeading, 'shipment-heading')}>
      <div className={styles.shipmentHeadingType}>
        <h2>{shipmentInfo.shipmentType}</h2>
        {shipmentStatus === shipmentStatuses.CANCELED && <Tag className="usa-tag--red">cancelled</Tag>}
        {shipmentInfo.isDiversion && <Tag>diversion</Tag>}
        {!shipmentInfo.isDiversion && shipmentStatus === shipmentStatuses.DIVERSION_REQUESTED && (
          <Tag>diversion requested</Tag>
        )}
      </div>
      <div className={styles.row}>
        <small>
          {`${shipmentInfo.originCity}, ${shipmentInfo.originState} ${shipmentInfo.originPostalCode} to
        ${formatDestinationAddress(shipmentInfo.destinationAddress)} on ${shipmentInfo.scheduledPickupDate}`}
        </small>
        {isCancelModalVisible && (
          <Restricted to={permissionTypes.createShipmentCancellation}>
            <Button type="button" onClick={() => handleShowCancellationModal(shipmentInfo)} unstyled>
              Request Cancellation
            </Button>
          </Restricted>
        )}
        {isCancellationRequested && <Tag>Cancellation Requested</Tag>}
      </div>
    </div>
  );
}

ShipmentHeading.propTypes = {
  shipmentInfo: PropTypes.shape({
    shipmentID: PropTypes.string.isRequired,
    shipmentType: PropTypes.string.isRequired,
    isDiversion: PropTypes.bool,
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
