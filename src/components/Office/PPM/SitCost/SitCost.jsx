import React from 'react';
import classnames from 'classnames';

import styles from './SitCost.module.scss';

import { useGetPPMSITEstimatedCostQuery } from 'hooks/queries';
import { formatCents, toDollarString } from 'utils/formatters';

export default function SitCost({ ppmShipmentInfo, ppmSITLocation, sitStartDate, sitEndDate, weightStored }) {
  const { estimatedCost, isLoading, isError } = useGetPPMSITEstimatedCostQuery(
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

  return (
    <div>
      <legend className={classnames('usa-label', styles.label)}>{costLabel}</legend>
      <div className={styles.displayValue} data-testid="costAmount">
        {toDollarString(formatCents(estimatedCost?.estimatedCost || 0))}
      </div>
    </div>
  );
}
