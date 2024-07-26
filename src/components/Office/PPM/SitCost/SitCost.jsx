import React, { useState } from 'react';
import classnames from 'classnames';
import { Button } from '@trussworks/react-uswds';

import styles from './SitCost.module.scss';

import { formatCents, toDollarString } from 'utils/formatters';
import SitCostBreakdown from 'components/Office/PPM/SitCostBreakdown/SitCostBreakdown';

export default function SitCost({
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
  const [calculationsVisible, setCalulationsVisible] = useState(false);

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
        <Button
          className={styles.togglePPMCalculations}
          type="button"
          data-testid="togglePPMCalculations"
          aria-expanded={calculationsVisible}
          unstyled
          onClick={() => {
            setCalulationsVisible((isVisible) => {
              return !isVisible;
            });
          }}
        >
          {calculationsVisible ? 'Hide calculations' : 'Show calculations'}
        </Button>
      </div>
      {calculationsVisible && (
        <div className={styles.calculationsContainer}>
          <SitCostBreakdown
            ppmShipmentInfo={ppmShipmentInfo}
            ppmSITLocation={ppmSITLocation}
            sitStartDate={sitStartDate}
            sitEndDate={sitEndDate}
            weightStored={weightStored}
            actualWeight={actualWeight}
            useQueries={useQueries}
            setEstimatedCost={setEstimatedCost}
          />
        </div>
      )}
    </div>
  );
}
