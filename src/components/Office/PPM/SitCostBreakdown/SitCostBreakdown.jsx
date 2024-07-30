import React from 'react';
import classnames from 'classnames';

import styles from './SitCostBreakdown.module.scss';

import { formatCents, formatDate, formatWeight, formatWeightCWTFromLbs, toDollarString } from 'utils/formatters';

export default function SitCostBreakdown({
  ppmShipmentInfo,
  ppmSITLocation,
  sitStartDate,
  sitAdditionalStartDate,
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
                SIT Information:
              </small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>
                {estimatedCost.paramsFirstDaySIT.serviceAreaOrigin
                  ? `Origin service area: ${estimatedCost?.paramsFirstDaySIT.serviceAreaOrigin}`
                  : `Destination service area: ${estimatedCost?.paramsFirstDaySIT.serviceAreaDestination}`}
              </small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Actual move date: {formatDate(ppmShipmentInfo.actualMoveDate)}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>{estimatedCost.paramsFirstDaySIT.isPeak ? 'Domestic peak' : 'Domestic non-peak'}</small>
            </div>
          </div>
          <div data-testid="column" className={styles.col}>
            <div data-testid="row" className={styles.row}>
              <small data-testid="label" className={styles.descriptionTitle}>
                Billable weight:
              </small>
              <small data-testid="value" className={styles.value}>
                {formatWeightCWTFromLbs(weightStored)}
              </small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Adjusted weight: {formatWeight(weightStored)}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Estimated SIT weight: {formatWeight(ppmShipmentInfo.sitEstimatedWeight)}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Actual PPM weight: {formatWeight(actualWeight)}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Estimated PPM weight: {formatWeight(ppmShipmentInfo.estimatedWeight)}</small>
            </div>
          </div>
          <div data-testid="column" className={styles.col}>
            <div data-testid="row" className={styles.row}>
              <small data-testid="label" className={styles.descriptionTitle}>
                First day SIT price:
              </small>
              <small data-testid="value" className={styles.value}>
                {toDollarString(formatCents(estimatedCost?.priceFirstDaySIT || 0))}/cwt
              </small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>SIT start date: {formatDate(sitStartDate)}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Base price: {toDollarString(estimatedCost.paramsFirstDaySIT.priceRateOrFactor)}</small>
            </div>
          </div>
          <div data-testid="column" className={styles.col}>
            <div data-testid="row" className={styles.row}>
              <small data-testid="label" className={styles.descriptionTitle}>
                Additional Day SIT price:
              </small>
              <small data-testid="value" className={styles.value}>
                {toDollarString(formatCents(estimatedCost?.priceAdditionalDaySIT || 0))}
              </small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>SIT add&apos;l day start: {formatDate(sitAdditionalStartDate)} </small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>SIT end date: {formatDate(sitEndDate)}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Additional days used: {estimatedCost?.paramsAdditionalDaySIT.numberDaysSIT}</small>
            </div>
            <div data-testid="details" className={styles.row}>
              <small>Price per day: {toDollarString(estimatedCost.paramsAdditionalDaySIT.priceRateOrFactor)}/cwt</small>
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
