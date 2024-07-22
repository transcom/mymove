import React from 'react';
import classnames from 'classnames';

import styles from './SitCost.module.scss';

import { formatCents, toDollarString } from 'utils/formatters';

export default function SitCost({
  ppmShipmentInfo,
  ppmSITLocation,
  sitStartDate,
  sitEndDate,
  weightStored,
  useQueries,
  setEstimatedCost,
}) {
  const { estimatedCost, isLoading, isError } = useQueries(
    ppmShipmentInfo.id,
    ppmSITLocation,
    sitStartDate,
    sitEndDate,
    weightStored,
  );

  const costLabel = 'Government SIT Cost';

  if (isLoading || isError) {
    return (
      <div data-testid="costAmount">
        <legend className={classnames('usa-label', styles.label)}>{costLabel}</legend>
        <div className={styles.displayValue}> {toDollarString(0)} </div>
      </div>
    );
  }

  setEstimatedCost(estimatedCost?.sitCost || 0);
  return (
    <div>
      <legend className={classnames('usa-label', styles.label)}>{costLabel}</legend>
      <div className={styles.displayValue} data-testid="costAmountSuccess">
        {toDollarString(formatCents(estimatedCost?.sitCost || 0))}
      </div>
    </div>
  );
}
