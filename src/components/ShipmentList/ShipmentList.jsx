import React from 'react';
import { arrayOf, bool, func, number, shape, string } from 'prop-types';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Tag } from '@trussworks/react-uswds';

import styles from './ShipmentList.module.scss';

import { isPPMShipmentComplete } from 'utils/shipments';
import { formatWeight } from 'utils/formatters';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { shipmentTypes } from 'constants/shipments';
import { shipmentIsOverweight } from 'utils/shipmentWeights';
import { PPMShipmentShape } from 'types/customerShapes';

export const ShipmentListItem = ({
  shipment,
  onShipmentClick,
  shipmentNumber,
  canEdit,
  showNumber,
  showIncomplete,
  showShipmentWeight,
  isOverweight,
  isMissingWeight,
}) => {
  const shipmentClassName = classnames({
    [styles[`shipment-list-item-NTS-release`]]: shipment.shipmentType === SHIPMENT_OPTIONS.NTSR,
    [styles[`shipment-list-item-NTS`]]: shipment.shipmentType === SHIPMENT_OPTIONS.NTS,
    [styles[`shipment-list-item-HHG`]]:
      shipment.shipmentType === SHIPMENT_OPTIONS.HHG ||
      shipment.shipmentType === SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC ||
      shipment.shipmentType === SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
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
        {shipmentTypes[shipment.shipmentType]}
        {showNumber && ` ${shipmentNumber}`}
      </strong>{' '}
      {/* use substring of the UUID until actual shipment code is available */}
      {!showShipmentWeight && !showIncomplete && (
        <span className={styles['shipment-code']}>#{shipment.id.substring(0, 8).toUpperCase()}</span>
      )}{' '}
      {showIncomplete && <Tag>Incomplete</Tag>}
      {showShipmentWeight && (
        <div className={styles.shipmentWeight}>{formatWeight(shipment.calculatedBillableWeight)}</div>
      )}
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
  showIncomplete: bool,
  showShipmentWeight: bool,
  isOverweight: bool,
  isMissingWeight: bool,
};

ShipmentListItem.defaultProps = {
  showNumber: true,
  showIncomplete: false,
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
        let canEdit = !moveSubmitted;
        let isOverweight;
        let isMissingWeight;
        let isIncomplete;
        let showNumber = shipmentCountByType[shipmentType] > 1;
        if (showShipmentWeight) {
          canEdit = false;
          showNumber = false;
          switch (shipmentType) {
            case SHIPMENT_OPTIONS.NTSR:
              // don't want “Over weight” or “Missing weight” warnings for NTSR
              break;
            default:
              if (shipmentIsOverweight(shipment.primeEstimatedWeight, shipment.calculatedBillableWeight)) {
                isOverweight = true;
              }
              if ((shipment.reweigh?.id && !shipment.reweigh?.weight) || !shipment.primeEstimatedWeight) {
                isMissingWeight = true;
              }
          }
        }
        if (shipmentType === SHIPMENT_OPTIONS.PPM) {
          if (isPPMShipmentComplete(shipment.ppmShipment)) {
            isIncomplete = true;
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
            showIncomplete={isIncomplete}
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
      reweigh: shape({ id: string.isRequired, weight: number }),
      ppmShipment: PPMShipmentShape,
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
