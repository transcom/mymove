/* eslint-disable security/detect-object-injection */
import React from 'react';
import { string, arrayOf, shape, func, number, bool } from 'prop-types';

import styles from './ShipmentList.module.scss';

import { ReactComponent as EditIcon } from 'shared/icon/edit.svg';

const ShipmentListItem = ({ shipment, onShipmentClick, shipmentNumber, canEdit, showNumber }) => {
  function handleEnterOrSpace(event) {
    if (!canEdit) return;
    const key = event.which || event.keyCode; // Use either which or keyCode, depending on browser support
    // enter or space
    if (key === 13 || key === 32) {
      onShipmentClick(shipment, shipmentNumber, shipment.shipmentType);
    }
  }
  const shipmentClassName = styles[`shipment-list-item-${shipment.shipmentType}`];
  return (
    <div
      className={`${styles['shipment-list-item-container']} ${shipmentClassName}`}
      data-testid="shipment-list-item-container"
      onClick={() => {
        if (!canEdit) return;
        onShipmentClick(shipment, shipmentNumber, shipment.shipmentType);
      }}
      onKeyDown={(event) => handleEnterOrSpace(event)}
      role="button"
      tabIndex="0"
    >
      <strong>
        {shipment.shipmentType}
        {showNumber && ` ${shipmentNumber}`}
      </strong>{' '}
      {/* use substring of the UUID until actual shipment code is available */}
      <span className={styles['shipment-code']}>{shipment.id.substring(0, 10)}</span>{' '}
      {canEdit ? <EditIcon className={styles.edit} /> : <div className={styles.noEdit} />}
    </div>
  );
};

ShipmentListItem.propTypes = {
  shipment: shape({ id: string.isRequired, shipmentType: string.isRequired }).isRequired,
  onShipmentClick: func.isRequired,
  shipmentNumber: number.isRequired,
  canEdit: bool.isRequired,
  showNumber: bool,
};

ShipmentListItem.defaultProps = {
  showNumber: true,
};

const ShipmentList = ({ shipments, onShipmentClick, moveSubmitted }) => {
  const shipmentNumbersByType = {};
  const shipmentCountByType = {};
  shipments.forEach((shipment) => {
    const { shipmentType } = shipment;
    if (shipmentCountByType[shipmentType]) {
      shipmentCountByType[shipmentType] += 1;
    } else {
      shipmentCountByType[shipmentType] = 1;
    }
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
        const canEdit = moveSubmitted ? shipmentType === 'PPM' : true;
        return (
          <ShipmentListItem
            key={shipment.id}
            shipmentNumber={shipmentNumber}
            showNumber={shipmentCountByType[shipmentType] > 1}
            canEdit={canEdit}
            onShipmentClick={() => onShipmentClick(shipment.id, shipmentNumber, shipmentType, canEdit)}
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
  moveSubmitted: bool.isRequired,
};

export default ShipmentList;
