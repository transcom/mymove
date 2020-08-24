import React from 'react';
import { string, arrayOf, shape, func } from 'prop-types';

import styles from './ShipmentList.module.scss';

import { ReactComponent as EditIcon } from 'shared/icon/edit.svg';

const ShipmentListItem = ({ shipment, onShipmentClick }) => {
  function handleEnterOrSpace(event) {
    const key = event.which || event.keyCode; // Use either which or keyCode, depending on browser support
    // enter or space
    if (key === 13 || key === 32) {
      onShipmentClick(shipment);
    }
  }
  const shipmentClassName = styles[`shipment-list-item-${shipment.type}`];
  return (
    <div
      className={`${styles['shipment-list-item-container']} ${shipmentClassName} margin-bottom-1`}
      onClick={() => onShipmentClick(shipment)}
      onKeyDown={(event) => handleEnterOrSpace(event)}
      role="button"
      tabIndex="0"
    >
      <strong>{shipment.type}</strong> <span>{shipment.id}</span> <EditIcon className={styles.edit} />
    </div>
  );
};

ShipmentListItem.propTypes = {
  shipment: shape({ id: string.isRequired, type: string.isRequired }).isRequired,
  onShipmentClick: func.isRequired,
};

const ShipmentList = ({ shipments, onShipmentClick }) => (
  <div>
    {shipments.map((shipment, index) => (
      <ShipmentListItem
        key={shipment.id}
        onShipmentClick={onShipmentClick}
        shipment={shipment}
        isLastItem={shipments.length - 1 === index}
      />
    ))}
  </div>
);

ShipmentList.propTypes = {
  shipments: arrayOf(shape({ id: string.isRequired, type: string.isRequired })).isRequired,
  onShipmentClick: func.isRequired,
};

export default ShipmentList;
