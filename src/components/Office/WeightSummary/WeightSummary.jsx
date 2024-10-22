import React from 'react';
import { number, arrayOf, shape } from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './WeightSummary.module.scss';

import { formatWeight } from 'utils/formatters';
import { shipmentIsOverweight, getDisplayWeight } from 'utils/shipmentWeights';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { shipmentTypes, WEIGHT_ADJUSTMENT } from 'constants/shipments';

const WeightSummary = ({ maxBillableWeight, weightRequested, weightAllowance, totalBillableWeight, shipments }) => {
  const formatShipments = shipments.slice();
  let countHHG = 0;
  if (formatShipments.filter((shipment) => shipment.shipmentType === 'HHG').length > 1) {
    for (let i = 0; i < formatShipments?.length; i += 1) {
      if (shipmentTypes.HHG) {
        countHHG += 1;
        formatShipments[i].count = countHHG;
      }
    }
  }

  const formatShipmentType = (shipment) => {
    if (shipment.shipmentType === SHIPMENT_OPTIONS.HHG && countHHG > 1) return `HHG ${shipment.count}`;
    if (shipment.shipmentType === SHIPMENT_OPTIONS.HHG && countHHG <= 1) return 'HHG';
    if (shipment.shipmentType === SHIPMENT_OPTIONS.NTS) return 'NTS';
    if (shipment.shipmentType === SHIPMENT_OPTIONS.NTSR) return 'NTSR';
    if (shipment.shipmentType === SHIPMENT_OPTIONS.PPM) return 'PPM';
    return '';
  };

  const displayShipments = formatShipments?.map((shipment) => {
    const displayWeight = getDisplayWeight(shipment, 1.1);
    return (
      <div className={styles.weight} key={shipment.id} data-testid="billableWeightCap">
        {shipmentIsOverweight(shipment.primeEstimatedWeight * WEIGHT_ADJUSTMENT, shipment.calculatedBillableWeight) ||
        !shipment.primeEstimatedWeight ||
        (shipment.reweigh?.dateReweighRequested && !shipment.reweigh?.weight) ? (
          <FontAwesomeIcon icon="exclamation-triangle" data-testid="shipmentHasFlag" className={styles.warningFlag} />
        ) : (
          <div className={styles.noEdit} />
        )}
        {formatWeight(displayWeight)} {formatShipmentType(shipment)}
      </div>
    );
  });

  return (
    <div className={styles.weightSummaryContainer}>
      <div>
        <h4 className={styles.weightSummaryHeading}>Max billable weight</h4>
        <div data-testid="maxBillableWeight" className={styles.marginBottom}>
          {formatWeight(maxBillableWeight)}
        </div>
        <h4 className={styles.weightSummaryHeading}>Actual weight</h4>
        <div data-testid="weightRequested" className={styles.marginBottom}>
          {formatWeight(weightRequested)}
        </div>
        <h4 className={styles.weightSummaryHeading}>Weight allowance</h4>
        <div data-testid="weightAllowance">{formatWeight(weightAllowance)}</div>
      </div>
      <div>
        <h4 className={styles.weightSummaryHeading}>Actual billable weight</h4>
        <div data-testid="totalBillableWeight" className={styles.weight}>
          {totalBillableWeight > maxBillableWeight ? (
            <FontAwesomeIcon
              icon="exclamation-circle"
              data-testid="totalBillableWeightFlag"
              className={styles.errorFlag}
            />
          ) : (
            <div className={styles.noEdit} />
          )}
          {formatWeight(totalBillableWeight)}
        </div>
        <hr />
        {displayShipments}
      </div>
    </div>
  );
};

WeightSummary.propTypes = {
  maxBillableWeight: number.isRequired,
  weightRequested: number.isRequired,
  weightAllowance: number.isRequired,
  totalBillableWeight: number.isRequired,
  shipments: arrayOf(
    shape({
      calculatedBillableWeight: number,
      primeEstimatedWeight: number,
    }),
  ).isRequired,
};

export default WeightSummary;
