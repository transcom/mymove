import React from 'react';
import { PropTypes, number } from 'prop-types';

import styles from './EditPPMNetWeight.module.scss';

import { formatWeight } from 'utils/formatters';
import { calculateWeightTicketWeightDifference, getWeightTicketNetWeight } from 'utils/shipmentWeights';
import { calculateWeightRequested } from 'hooks/custom';
import { ShipmentShape, WeightTicketShape } from 'types/shipment';

// Labels & constants
const CALCULATION_TYPE = {
  NET_WEIGHT: 'NET_WEIGHT',
  EXCESS_WEIGHT: 'EXCESS_WEIGHT',
  REDUCE_WEIGHT: 'REDUCE_WEIGHT',
};
const weightLabels = {
  [CALCULATION_TYPE.NET_WEIGHT]: {
    firstLabel: ' | original weight',
    secondLabel: ' | to fit within weight allowance',
  },
  [CALCULATION_TYPE.REDUCE_WEIGHT]: {
    firstLabel: ' | original weight',
    secondLabel: ' | to reduce excess weight',
  },
  [CALCULATION_TYPE.EXCESS_WEIGHT]: {
    firstLabel: 'Move weight (total)',
    secondLabel: 'Weight allowance',
    thirdLabel: 'Excess weight (total)',
  },
};

// Flexbox wrapper
const FlexContainer = ({ children, className }) => {
  return (
    <div className={className} style={{ display: 'flex' }}>
      {children}
    </div>
  );
};

const WeightCalculationHint = ({ type, firstValue, secondValue, thirdValue }) => {
  const { firstLabel, secondLabel, thirdLabel } = weightLabels[type];
  return (
    <>
      <FlexContainer className={styles.minus}>
        <div className={(!thirdValue && styles.calculationWrapperDisplay) || styles.calculationWrapper}>
          <div className={styles.calculations}>
            <strong className={styles.value}>{formatWeight(firstValue)}</strong>
            <span className={styles.label}> {firstLabel}</span>
          </div>
          {secondValue && (
            <div className={styles.calculations}>
              <strong className={styles.value}>
                {thirdValue && <>– </>}
                {formatWeight(secondValue)}
              </strong>
              <span className={styles.label}> {secondLabel}</span>
            </div>
          )}
        </div>
      </FlexContainer>
      {thirdValue && (
        <>
          <hr className={styles.divider} />
          <div className={styles.calculations}>
            <strong className={styles.value}>{formatWeight(thirdValue)}</strong>
            <span className={styles.label}> {thirdLabel}</span>
          </div>
        </>
      )}
    </>
  );
};

const EditPPMNetWeight = ({ weightTicket, weightAllowance, shipments }) => {
  // Original weight is the full weight - empty weight
  const originalWeight = calculateWeightTicketWeightDifference(weightTicket);
  // moveWeightTotal = Sum of all ppm weights + sum of all non-ppm shipments
  const moveWeightTotal = calculateWeightRequested(shipments);
  const excessWeight = moveWeightTotal - weightAllowance;
  const hasExcessWeight = Boolean(excessWeight > 0);
  const netWeight = getWeightTicketNetWeight(weightTicket);

  const toFitValue = hasExcessWeight ? -Math.min(excessWeight, netWeight) : null;
  const showWarning = Boolean(hasExcessWeight);
  const showReduceWeight = Boolean(-originalWeight === toFitValue);
  return (
    <div className={styles.main_wrapper}>
      <div>
        <h4 className={styles.mainHeader}>PPM Trip Weight</h4>
      </div>
      <FlexContainer className={styles.netWeightContainer}>
        {showWarning && <div className={styles.warnings} data-testid="warning" />}
        <div>
          <h5 className={styles.header}>Net weight</h5>
          <WeightCalculationHint
            firstValue={originalWeight}
            secondValue={toFitValue}
            type={showReduceWeight ? CALCULATION_TYPE.REDUCE_WEIGHT : CALCULATION_TYPE.NET_WEIGHT}
          />
          <div data-testid="net-weight-display" className={styles.wrapper}>
            <div className={styles.netWeightDisplay}>{formatWeight(netWeight)}</div>
            {weightTicket.netWeightRemarks && (
              <>
                <h5 className={styles.remarksHeader}>Remarks</h5>
                <p className={styles.remarks}>{weightTicket.netWeightRemarks}</p>
              </>
            )}
          </div>
        </div>
      </FlexContainer>
    </div>
  );
};

EditPPMNetWeight.propTypes = {
  weightTicket: WeightTicketShape.isRequired,
  weightAllowance: number.isRequired,
  shipments: PropTypes.arrayOf(ShipmentShape).isRequired,
};

export default EditPPMNetWeight;
