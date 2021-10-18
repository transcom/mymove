import React from 'react';

import descriptionListStyles from '../../../styles/descriptionList.module.scss';
import { shipmentTypeLabels } from '../../../content/shipments';
import { formatDateFromIso } from '../../../shared/formatters';
import { ShipmentShape } from '../../../types/shipment';

const Shipment = ({ shipment }) => {
  return (
    <dl className={descriptionListStyles.descriptionList}>
      <h3>{`${shipmentTypeLabels[shipment.shipmentType]} shipment`}</h3>
      <div className={descriptionListStyles.row}>
        <dt>Status:</dt>
        <dd>{shipment.status}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Shipment ID:</dt>
        <dd>{shipment.id}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Shipment eTag:</dt>
        <dd>{shipment.eTag}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Requested Pickup Date:</dt>
        <dd>{shipment.requestedPickupDate}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Actual Pickup Date:</dt>
        <dd>{shipment.actualPickupDate}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Estimated Weight:</dt>
        <dd>{shipment.primeEstimatedWeight}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Actual Weight:</dt>
        <dd>{shipment.primeActualWeight}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Pickup Address:</dt>
        <dd>
          {shipment.pickupAddress.streetAddress1} {shipment.pickupAddress.streetAddress2} {shipment.pickupAddress.city}{' '}
          {shipment.pickupAddress.state} {shipment.pickupAddress.postalCode}
        </dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Destination Address:</dt>
        <dd>
          {shipment.destinationAddress.streetAddress1} {shipment.destinationAddress.streetAddress2}{' '}
          {shipment.destinationAddress.city} {shipment.destinationAddress.state}{' '}
          {shipment.destinationAddress.postalCode}
        </dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Created at:</dt>
        <dd>{formatDateFromIso(shipment.createdAt, 'YYYY-MM-DD')}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Approved at:</dt>
        <dd>{shipment.approvedDate}</dd>
      </div>
    </dl>
  );
};

Shipment.propTypes = {
  shipment: ShipmentShape.isRequired,
};

export default Shipment;
