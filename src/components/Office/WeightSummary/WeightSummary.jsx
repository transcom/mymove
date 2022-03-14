import React from 'react';
import { number, arrayOf, shape } from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './WeightSummary.module.scss';

import { formatWeight } from 'utils/formatters';
import { shipmentIsOverweight } from 'utils/shipmentWeights';

const WeightSummary = ({ maxBillableWeight, weightRequested, weightAllowance, totalBillableWeight, shipments }) => {
  return (
    <div className={styles.weightSummaryContainer}>
      <div>
        <h4 className={styles.weightSummaryHeading}>Max billable weight</h4>
        <div data-testid="maxBillableWeight" className={styles.marginBottom}>
          {formatWeight(maxBillableWeight)}
        </div>
        <h4 className={styles.weightSummaryHeading}>Weight requested</h4>
        <div data-testid="weightRequested" className={styles.marginBottom}>
          {formatWeight(weightRequested)}
        </div>
        <h4 className={styles.weightSummaryHeading}>Weight allowance</h4>
        <div data-testid="weightAllowance">{formatWeight(weightAllowance)}</div>
      </div>
      <div>
        <h4 className={styles.weightSummaryHeading}>Total billable weight</h4>
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
        {shipments?.map((shipment) => {
          return (
            <div className={styles.weight} key={shipment.id} data-testid="billableWeightCap">
              {shipmentIsOverweight(shipment.primeEstimatedWeight, shipment.calculatedBillableWeight) ||
              !shipment.primeEstimatedWeight ||
              (shipment.reweigh?.dateReweighRequested && !shipment.reweigh?.weight) ? (
                <FontAwesomeIcon
                  icon="exclamation-triangle"
                  data-testid="shipmentHasFlag"
                  className={styles.warningFlag}
                />
              ) : (
                <div className={styles.noEdit} />
              )}
              {formatWeight(shipment.calculatedBillableWeight)}
            </div>
          );
        })}
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
      calculatedBillableWeight: number.isRequired,
      primeEstimatedWeight: number,
    }),
  ).isRequired,
};

export default WeightSummary;
