/* eslint-disable security/detect-object-injection */
import React from 'react';
import { string, arrayOf, shape, func, number, bool } from 'prop-types';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './ShipmentList.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { getShipmentTypeLabel } from 'utils/shipmentDisplay';

const ShipmentListItem = ({ shipment, onShipmentClick, shipmentNumber, canEdit, showNumber }) => {
  const shipmentClassName = classnames({
    [styles[`shipment-list-item-NTS-R`]]: shipment.shipmentType === SHIPMENT_OPTIONS.NTSR,
    [styles[`shipment-list-item-NTS`]]: shipment.shipmentType === SHIPMENT_OPTIONS.NTS,
    [styles[`shipment-list-item-HHG`]]: shipment.shipmentType === SHIPMENT_OPTIONS.HHG,
    [styles[`shipment-list-item-PPM`]]: shipment.shipmentType === SHIPMENT_OPTIONS.PPM,
  });

  return (
    <button
      type="button"
      className={`${styles['shipment-list-item-container']} ${shipmentClassName}`}
      data-testid="shipment-list-item-container"
      onClick={() => {
        if (!canEdit) return;
        onShipmentClick();
      }}
      tabIndex="0"
    >
      <strong>
        {getShipmentTypeLabel(shipment.shipmentType)}
        {showNumber && ` ${shipmentNumber}`}
      </strong>{' '}
      {/* use substring of the UUID until actual shipment code is available */}
      <span className={styles['shipment-code']}>#{shipment.id.substring(0, 8).toUpperCase()}</span>{' '}
      {canEdit ? <FontAwesomeIcon icon="pen" className={styles.edit} /> : <div className={styles.noEdit} />}
    </button>
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
  moveSubmitted: bool.isRequired,
};

export default ShipmentList;
