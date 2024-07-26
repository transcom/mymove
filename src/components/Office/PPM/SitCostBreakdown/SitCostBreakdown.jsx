import React from 'react';
import classnames from 'classnames';

import styles from './SitCostBreakdown.module.scss';

import { formatCents, formatDate, toDollarString } from 'utils/formatters';

export default function SitCostBreakdown({
  ppmShipmentInfo,
  ppmSITLocation,
  sitStartDate,
  sitEndDate,
  weightStored,
  actualWeight,
  useQueries,
  setEstimatedCost,
}) {
  const { estimatedCost, isLoading, isError } = useQueries(
    ppmShipmentInfo.id,
    ppmSITLocation,
    sitStartDate,
    sitEndDate,
    weightStored,
    actualWeight,
  );

  if (isLoading || isError) {
    return (
      <div data-testid="costAmount">
        <legend className={classnames('usa-label', styles.label)}>Calculation Values</legend>
        <div className={styles.displayValue}> {toDollarString(0)} </div>
      </div>
    );
  }

  setEstimatedCost(estimatedCost?.sitCost || 0);
  return (
    <div data-testid="SitCostBreakdown" className={styles.SitCostBreakdown}>
      <div>
        <h4 className={styles.title}>Calculations</h4>
        <div data-testid="flexGridSmall" className={styles.flexGridSmall}>
          <div data-testid="column" className={styles.col}>
            <div data-testid="row" className={styles.row}>
              <small data-testid="label" className={styles.descriptionTitle}>
                Billable weight:
              </small>
              <small data-testid="value" className={styles.value}>
                {weightStored}
              </small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Adjusted weight: {weightStored}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Estimated SIT weight: {ppmShipmentInfo.sitEstimatedWeight}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Actual PPM weight: {actualWeight}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Estimated PPM weight: {ppmShipmentInfo.estimatedWeight}</small>
            </div>
          </div>
          <div data-testid="column" className={styles.col}>
            <div data-testid="row" className={styles.row}>
              <small data-testid="label" className={styles.descriptionTitle}>
                First day SIT price:
              </small>
              <small data-testid="value" className={styles.value}>
                {toDollarString(formatCents(estimatedCost?.priceFirstDaySIT || 0))}
              </small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Origin service area: {estimatedCost?.paramsFirstDaySIT.serviceAreaOrigin}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Actual pickup date: {formatDate(ppmShipmentInfo.actualMoveDate)}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Domestic peak</small>
            </div>
          </div>
          <div data-testid="column" className={styles.col}>
            <div data-testid="row" className={styles.row}>
              <small data-testid="label" className={styles.descriptionTitle}>
                Additonal SIT price:
              </small>
              <small data-testid="value" className={styles.value}>
                {toDollarString(formatCents(estimatedCost?.priceAdditionalDaySIT || 0))}
              </small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Origin service area: {estimatedCost?.paramsAdditionalDaySIT.serviceAreaOrigin}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Actual pickup date: {formatDate(ppmShipmentInfo.actualMoveDate)}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Domestic peak</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Additional days used: {estimatedCost?.paramsAdditionalDaySIT.numberDaysSIT}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>
                Price per day: {toDollarString(formatCents(estimatedCost?.priceAdditionalDaySIT || 0))}/
                {estimatedCost?.paramsAdditionalDaySIT.numberDaysSIT}
              </small>
            </div>
          </div>
          <div data-testid="column" className={styles.col}>
            <div data-testid="row" className={styles.row}>
              <small data-testid="label" className={styles.descriptionTitle}>
                Price escalation factor:
              </small>
              <small data-testid="value" className={styles.value}>
                {estimatedCost?.paramsFirstDaySIT.escalationCompounded}
              </small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Base year: {estimatedCost?.paramsFirstDaySIT.contractYearName}</small>
            </div>
          </div>
          <div data-testid="column" className={styles.col}>
            <div data-testid="row" className={styles.row}>
              <small data-testid="label" className={styles.descriptionTitle}>
                Total:
              </small>
              <small data-testid="value" className={styles.value}>
                {toDollarString(formatCents(estimatedCost?.sitCost || 0))}
              </small>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
