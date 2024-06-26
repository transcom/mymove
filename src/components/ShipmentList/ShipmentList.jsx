import React from 'react';
import { arrayOf, bool, func, number, shape, string } from 'prop-types';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Tag, Button } from '@trussworks/react-uswds';

import styles from './ShipmentList.module.scss';

import { shipmentTypes } from 'constants/shipments';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { ShipmentShape } from 'types/shipment';
import { formatWeight } from 'utils/formatters';
import { isPPMShipmentComplete } from 'utils/shipments';
import { shipmentIsOverweight } from 'utils/shipmentWeights';
import ToolTip from 'shared/ToolTip/ToolTip';

export const ShipmentListItem = ({
  shipment,
  onShipmentClick,
  onDeleteClick,
  shipmentNumber,
  canEditOrDelete,
  showNumber,
  showIncomplete,
  showShipmentWeight,
  isOverweight,
  isMissingWeight,
}) => {
  const shipmentClassName = classnames({
    [styles[`shipment-list-item-NTS-release`]]: shipment.shipmentType === SHIPMENT_OPTIONS.NTSR,
    [styles[`shipment-list-item-NTS`]]: shipment.shipmentType === SHIPMENT_OPTIONS.NTS,
    [styles[`shipment-list-item-HHG`]]: shipment.shipmentType === SHIPMENT_OPTIONS.HHG,
    [styles[`shipment-list-item-PPM`]]: shipment.shipmentType === SHIPMENT_OPTIONS.PPM,
  });
  const isPPM = shipment.shipmentType === SHIPMENT_OPTIONS.PPM;
  const estimated = 'Estimated';
  const actual = 'Actual';
  let requestedWeightPPM = 0;
  if (shipment.shipmentType === SHIPMENT_OPTIONS.PPM) {
    if (shipment.ppmShipment?.weightTickets !== undefined) {
      const wt = shipment.ppmShipment.weightTickets;
      for (let i = 0; i < wt.length; i += 1) {
        requestedWeightPPM += wt[i].fullWeight - wt[i].emptyWeight;
      }
    }
  }
  return (
    <div
      className={`${styles['shipment-list-item-container']} ${shipmentClassName} ${
        showShipmentWeight && styles['shipment-display']
      }`}
      data-testid="shipment-list-item-container"
    >
      <div>
        <strong>
          {shipmentTypes[shipment.shipmentType]}
          {showNumber && ` ${shipmentNumber}`}
        </strong>{' '}
        <br />
        {(shipment.shipmentType === SHIPMENT_OPTIONS.HHG || shipment.shipmentType === SHIPMENT_OPTIONS.NTS) && (
          <>
            `${formatWeight(shipment.primeEstimatedWeight * 1.1)} `
            <ToolTip text="110% Prime Estimated Weight" icon="circle-question" closeOnLeave />
          </>
        )}
      </div>
      {/* use substring of the UUID until actual shipment code is available */}
      {!showShipmentWeight && !showIncomplete && (
        <span className={styles['shipment-code']}>#{shipment.shipmentLocator}</span>
      )}{' '}
      {showIncomplete && <Tag>Incomplete</Tag>}
      {showShipmentWeight && (
        <div className={styles.shipmentWeight}>
          {isPPM && (
            <div className={styles.spaceBetween}>
              <div className={styles.textAlignRight}>
                <h6>{estimated}</h6>
                <h6>{actual}</h6>
              </div>
              <div className={styles.textAlignLeft}>
                <h6>
                  <strong>{formatWeight(shipment.ppmShipment.estimatedWeight)}</strong>
                </h6>
                <h6>
                  <strong>{requestedWeightPPM > 0 ? formatWeight(requestedWeightPPM) : '-'}</strong>
                </h6>
              </div>
            </div>
          )}
          {!isPPM && formatWeight(shipment.calculatedBillableWeight)}
        </div>
      )}
      {(isOverweight || (isMissingWeight && !isPPM)) && (
        <div className={styles['warning-section']}>
          <FontAwesomeIcon icon="exclamation-triangle" className={styles.warning} />
          <span className={styles.warningText}>{isOverweight ? 'Over weight' : 'Missing weight'}</span>
        </div>
      )}
      {canEditOrDelete ? (
        <div className={styles['shipment-btns']}>
          <Button className={styles['edit-btn']} onClick={onDeleteClick} type="button">
            Delete
          </Button>
          |
          <Button
            className={styles['edit-btn']}
            onClick={onShipmentClick}
            type="button"
            data-testid="editShipmentButton"
          >
            Edit
          </Button>
        </div>
      ) : (
        <div className={styles.noEdit} />
      )}
    </div>
  );
};

ShipmentListItem.propTypes = {
  shipment: shape({ id: string.isRequired, shipmentLocator: string.isRequired, shipmentType: string.isRequired })
    .isRequired,
  onShipmentClick: func,
  onDeleteClick: func,
  shipmentNumber: number.isRequired,
  canEditOrDelete: bool.isRequired,
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
  onDeleteClick: null,
};

const ShipmentList = ({ shipments, onShipmentClick, onDeleteClick, moveSubmitted, showShipmentWeight }) => {
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
        let canEditOrDelete = !moveSubmitted;
        let isOverweight;
        let isMissingWeight;
        let showNumber = shipmentCountByType[shipmentType] > 1;
        const ppmIsIncomplete = shipmentType === SHIPMENT_OPTIONS.PPM && !isPPMShipmentComplete(shipment);

        if (showShipmentWeight) {
          canEditOrDelete = false;
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

        return (
          <ShipmentListItem
            key={shipment.id}
            shipmentNumber={shipmentNumber}
            showNumber={showNumber}
            showShipmentWeight={showShipmentWeight}
            canEditOrDelete={canEditOrDelete}
            isOverweight={isOverweight}
            showIncomplete={ppmIsIncomplete}
            isMissingWeight={isMissingWeight}
            onShipmentClick={() => onShipmentClick(shipment.id, shipmentNumber, shipmentType)}
            onDeleteClick={() => onDeleteClick(shipment.id)}
            shipment={shipment}
          />
        );
      })}
    </div>
  );
};

ShipmentList.propTypes = {
  shipments: arrayOf(ShipmentShape).isRequired,
  onShipmentClick: func,
  onDeleteClick: func,
  moveSubmitted: bool.isRequired,
  showShipmentWeight: bool,
};

ShipmentList.defaultProps = {
  showShipmentWeight: false,
  onShipmentClick: null,
  onDeleteClick: null,
};

export default ShipmentList;
