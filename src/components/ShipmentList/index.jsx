import React from 'react';
import { arrayOf, bool, func, number, shape, string } from 'prop-types';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './ShipmentList.module.scss';

import { formatWeight } from 'shared/formatters';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { getShipmentTypeLabel } from 'utils/shipmentDisplay';

export const ShipmentListItem = ({
  shipment,
  onShipmentClick,
  shipmentNumber,
  canEdit,
  showNumber,
  showShipmentWeight,
  isOverweight,
  isMissingWeight,
}) => {
  const shipmentClassName = classnames({
    [styles[`shipment-list-item-NTS-R`]]: shipment.shipmentType === SHIPMENT_OPTIONS.NTSR,
    [styles[`shipment-list-item-NTS`]]: shipment.shipmentType === SHIPMENT_OPTIONS.NTS,
    [styles[`shipment-list-item-HHG`]]: shipment.shipmentType === SHIPMENT_OPTIONS.HHG,
    [styles[`shipment-list-item-PPM`]]: shipment.shipmentType === SHIPMENT_OPTIONS.PPM,
  });

  return (
    <button
      type="button"
      className={`${styles['shipment-list-item-container']} ${shipmentClassName} ${
        showShipmentWeight && styles['shipment-display']
      }`}
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
      {!showShipmentWeight && (
        <span className={styles['shipment-code']}>#{shipment.id.substring(0, 8).toUpperCase()}</span>
      )}{' '}
      {showShipmentWeight && <div className={styles.shipmentWeight}>{formatWeight(shipment.billableWeightCap)}</div>}
      {(isOverweight || isMissingWeight) && (
        <div>
          <FontAwesomeIcon icon="exclamation-triangle" className={styles.warning} />
          <span className={styles.warningText}>{isOverweight ? 'Over weight' : 'Missing weight'}</span>
        </div>
      )}
      {canEdit ? <FontAwesomeIcon icon="pen" className={styles.edit} /> : <div className={styles.noEdit} />}
    </button>
  );
};

ShipmentListItem.propTypes = {
  shipment: shape({ id: string.isRequired, shipmentType: string.isRequired }).isRequired,
  onShipmentClick: func,
  shipmentNumber: number.isRequired,
  canEdit: bool.isRequired,
  showNumber: bool,
  showShipmentWeight: bool,
  isOverweight: bool,
  isMissingWeight: bool,
};

ShipmentListItem.defaultProps = {
  showNumber: true,
  showShipmentWeight: false,
  isOverweight: false,
  isMissingWeight: false,
  onShipmentClick: null,
};

const ShipmentList = ({ shipments, onShipmentClick, moveSubmitted, showShipmentWeight }) => {
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
        let canEdit = moveSubmitted ? shipmentType === 'PPM' : true;
        let isOverweight;
        let isMissingWeight;
        let showNumber = shipmentCountByType[shipmentType] > 1;
        if (showShipmentWeight) {
          canEdit = false;
          showNumber = false;
          if (parseInt(shipment.billableWeightCap, 10) > parseInt(shipment.primeEstimatedWeight, 10) * 1.1) {
            isOverweight = true;
          }
          if (shipment.reweigh?.id && !shipment.reweigh?.weight) {
            isMissingWeight = true;
          }
        }
        return (
          <ShipmentListItem
            key={shipment.id}
            shipmentNumber={shipmentNumber}
            showNumber={showNumber}
            showShipmentWeight={showShipmentWeight}
            canEdit={canEdit}
            isOverweight={isOverweight}
            isMissingWeight={isMissingWeight}
            onShipmentClick={() => onShipmentClick(shipment.id, shipmentNumber, shipmentType)}
            shipment={shipment}
          />
        );
      })}
    </div>
  );
};

ShipmentList.propTypes = {
  shipments: arrayOf(
    shape({
      id: string.isRequired,
      shipmentType: string.isRequired,
      reweigh: shape({ id: string.isRequired, weight: string }),
    }),
  ).isRequired,
  onShipmentClick: func,
  moveSubmitted: bool.isRequired,
  showShipmentWeight: bool,
};

ShipmentList.defaultProps = {
  showShipmentWeight: false,
  onShipmentClick: null,
};

export default ShipmentList;
