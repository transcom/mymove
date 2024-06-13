import React from 'react';
import classnames from 'classnames';

import styles from './SitCost.module.scss';

import { useGetPPMSITEstimatedCostQuery } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { formatCents, toDollarString } from 'utils/formatters';

export default function SitCost({ ppmShipmentInfo, ppmSITLocation, sitStartDate, sitEndDate, weightStored }) {
  const { estimatedCost, isLoading, isError } = useGetPPMSITEstimatedCostQuery(
    ppmShipmentInfo.id,
    ppmSITLocation,
    sitStartDate,
    sitEndDate,
    weightStored,
  );

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;
  return (
    <div>
      <legend className={classnames('usa-label', styles.label)}>Cost</legend>
      <div className={styles.displayValue}> {toDollarString(formatCents(estimatedCost?.estimatedCost || 0))} </div>
    </div>
  );
}
