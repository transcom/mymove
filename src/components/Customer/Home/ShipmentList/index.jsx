import React from 'react';
import { string, arrayOf, shape, func, number } from 'prop-types';

import styles from './ShipmentList.module.scss';

import { ReactComponent as EditIcon } from 'shared/icon/edit.svg';

const ShipmentListItem = ({ shipment, onShipmentClick, shipmentNumber }) => {
  function handleEnterOrSpace(event) {
    const key = event.which || event.keyCode; // Use either which or keyCode, depending on browser support
    // enter or space
    if (key === 13 || key === 32) {
      onShipmentClick(shipment);
    }
  }
  const shipmentClassName = styles[`shipment-list-item-${shipment.shipmentType}`];
  return (
    <div
      className={`${styles['shipment-list-item-container']} ${shipmentClassName}`}
      data-testid="shipment-list-item-container"
      onClick={() => onShipmentClick(shipment, shipmentNumber)}
      onKeyDown={(event) => handleEnterOrSpace(event)}
      role="button"
      tabIndex="0"
    >
      <strong>
        {shipment.shipmentType}
        {` ${shipmentNumber}`}
      </strong>{' '}
      <span>{shipment.id}</span> <EditIcon className={styles.edit} />
    </div>
  );
};

ShipmentListItem.propTypes = {
  shipment: shape({ id: string.isRequired, shipmentType: string.isRequired }).isRequired,
  onShipmentClick: func.isRequired,
  shipmentNumber: number.isRequired,
};

const ShipmentList = ({ shipments, onShipmentClick }) => (
  <div>
    {shipments.map((shipment, index) => (
      <ShipmentListItem
        key={shipment.id}
        shipmentNumber={index + 1}
        onShipmentClick={() => onShipmentClick(shipment.id, index + 1)}
        shipment={shipment}
      />
    ))}
  </div>
);

ShipmentList.propTypes = {
  shipments: arrayOf(shape({ id: string.isRequired, shipmentType: string.isRequired })).isRequired,
  onShipmentClick: func.isRequired,
};

export default ShipmentList;
