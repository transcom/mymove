/* eslint-disable security/detect-object-injection */
import React from 'react';
import { string, arrayOf, shape, func, number, bool } from 'prop-types';

import styles from './ShipmentList.module.scss';

import { ReactComponent as EditIcon } from 'shared/icon/edit.svg';

const ShipmentListItem = ({ shipment, onShipmentClick, shipmentNumber, showNumber }) => {
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
      onClick={() => onShipmentClick(shipment, shipmentNumber, shipment.shipmentType)}
      onKeyDown={(event) => handleEnterOrSpace(event)}
      role="button"
      tabIndex="0"
    >
      <strong>
        {shipment.shipmentType}
        {showNumber && ` ${shipmentNumber}`}
      </strong>{' '}
      {/* use substring  of the UUID until actual shipment code is available */}
      <span className={styles['shipment-code']}>{shipment.id.substring(0, 10)}</span>{' '}
      <EditIcon className={styles.edit} />
    </div>
  );
};

ShipmentListItem.propTypes = {
  shipment: shape({ id: string.isRequired, shipmentType: string.isRequired }).isRequired,
  onShipmentClick: func.isRequired,
  shipmentNumber: number.isRequired,
  showNumber: bool,
};

ShipmentListItem.defaultProps = {
  showNumber: true,
};

const ShipmentList = ({ shipments, onShipmentClick }) => {
  const shipmentNumbersByType = {};
  const shipmentCountByType = {};
  shipments.map((shipment) => {
    const { shipmentType } = shipment;
    if (shipmentCountByType[shipmentType]) {
      shipmentCountByType[shipmentType] += 1;
    } else {
      shipmentCountByType[shipmentType] = 1;
    }
    return shipmentCountByType;
  });

  return (
    <div>
      {shipments.map((shipment) => {
        const { shipmentType } = shipment;
        if (shipmentNumbersByType[shipmentType]) {
          shipmentNumbersByType[shipmentType] += 1;
        } else {
          shipmentNumbersByType[shipmentType] = 1;
        }
        const shipmentNumber = shipmentNumbersByType[shipmentType];
        return (
          <ShipmentListItem
            key={shipment.id}
            shipmentNumber={shipmentNumber}
            showNumber={shipmentCountByType[shipmentType] > 1}
            onShipmentClick={() => onShipmentClick(shipment.id, shipmentNumber, shipmentType)}
            shipment={shipment}
          />
        );
      })}
    </div>
  );
};

ShipmentList.propTypes = {
  shipments: arrayOf(shape({ id: string.isRequired, shipmentType: string.isRequired })).isRequired,
  onShipmentClick: func.isRequired,
};

export default ShipmentList;
