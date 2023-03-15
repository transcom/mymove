import React, { useState } from 'react';
import { Formik } from 'formik';
import { useQueryClient } from '@tanstack/react-query';
import * as Yup from 'yup';
import { PropTypes, func, number } from 'prop-types';
import { Button, Fieldset, Label, Textarea } from '@trussworks/react-uswds';

import styles from './EditPPMNetWeight.module.scss';

import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { ErrorMessage } from 'components/form/ErrorMessage';
import { formatWeight } from 'utils/formatters';
import { calculateWeightTicketWeightDifference, getWeightTicketNetWeight } from 'utils/shipmentWeights';
import { useCalculatedWeightRequested } from 'hooks/custom';
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

// Form Error Indicator
const ErrorIndicator = ({ children, hasErrors }) => {
  return (
    <div data-testid="errorIndicator" className={hasErrors ? 'usa-form-group--error' : ''}>
      {children}
    </div>
  );
};

const WeightCalculationHint = ({ type, firstValue, secondValue, thirdValue }) => {
  const { firstLabel, secondLabel, thirdLabel } = weightLabels[type];
  return (
    <>
      <FlexContainer className={styles.minus}>
        {thirdValue && <>-</>}
        <div className={styles.calculationWrapper}>
          <div className={styles.calculations}>
            <strong className={styles.value}>{formatWeight(firstValue)}</strong>
            <span className={styles.label}> {firstLabel}</span>
          </div>
          {secondValue && (
            <div className={styles.calculations}>
              <strong className={styles.value}>{formatWeight(secondValue)}</strong>
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

const validationSchema = Yup.object({
  adjustedNetWeight: Yup.number().min(0, 'Net weight must be 0 lbs or greater').required('Required'),
  netWeightRemarks: Yup.string().required('Required'),
});

const EditPPMNetWeightForm = ({ onSave, onCancel, initialValues }) => {
  const queryClient = useQueryClient();

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema}>
      {({ handleChange, values, isValid, errors, touched, setTouched }) => (
        <div>
          <Fieldset className={styles.fieldset}>
            <MaskedTextField
              data-testid="weightInput"
              defaultValue="0"
              id="adjustedNetWeight"
              name="adjustedNetWeight"
              mask={Number}
              lazy={false}
              scale={0}
              signed={false} // no negative numbers
              thousandsSeparator=","
              suffix="lbs"
              inputClassName={styles.weightInput}
              errorClassName={styles.errors}
              labelClassName={styles.weightLabel}
              value={values.adjustedNetWeight}
            />
            <Label htmlFor="remarks">Remarks</Label>
            <ErrorMessage className={styles.errors} display={!!touched.netWeightRemarks && !!errors.netWeightRemarks}>
              {errors.netWeightRemarks}
            </ErrorMessage>
            <ErrorIndicator hasErrors={!!touched.netWeightRemarks && !!errors.netWeightRemarks}>
              <Textarea
                id="netWeightRemarks"
                data-testid="formRemarks"
                maxLength={500}
                placeholder=""
                onChange={handleChange}
                onBlur={() => {
                  setTouched({ netWeightRemarks: true }, false);
                }}
                value={values.netWeightRemarks}
              />
            </ErrorIndicator>
            <FlexContainer className={styles.wrapper}>
              <Button
                onClick={() => {
                  onSave(
                    { ...initialValues, ...values },
                    {
                      onSuccess: (/* data, variables, context */) => {
                        queryClient.invalidateQueries(); // FIX: @rogeruiz - This invalidates everything for now.
                        onCancel();
                      },
                    },
                  );
                }}
                disabled={!isValid}
              >
                Save changes
              </Button>
              <Button unstyled onClick={onCancel}>
                Cancel
              </Button>
            </FlexContainer>
          </Fieldset>
        </div>
      )}
    </Formik>
  );
};

const EditPPMNetWeight = ({ weightTicket, weightAllowance, shipments, editNetWeight }) => {
  const [showEditForm, setShowEditForm] = useState(false);

  const toggleEditForm = () => {
    setShowEditForm(!showEditForm);
  };

  // Original weight is the full weight - empty weight
  const originalWeight = calculateWeightTicketWeightDifference(weightTicket);
  // moveWeightTotal = Sum of all ppm weights + sum of all non-ppm shipments
  const moveWeightTotal = useCalculatedWeightRequested(shipments);
  const excessWeight = moveWeightTotal - weightAllowance;
  const hasExcessWeight = Boolean(excessWeight > 0);
  const netWeight = getWeightTicketNetWeight(weightTicket);
  const toFitValue = hasExcessWeight ? -Math.min(excessWeight, netWeight) : null;
  const showWarning = Boolean(hasExcessWeight && !showEditForm);
  const showReduceWeight = Boolean(-originalWeight === toFitValue);
  return (
    <div className={styles.wrapper}>
      <div>
        <h4 className={styles.mainHeader}>Edit PPM net weight</h4>
        {Boolean(showEditForm && hasExcessWeight) && (
          <WeightCalculationHint
            firstValue={moveWeightTotal}
            secondValue={weightAllowance}
            thirdValue={excessWeight}
            type={CALCULATION_TYPE.EXCESS_WEIGHT}
          />
        )}
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
          {!showEditForm ? (
            <div className={styles.wrapper}>
              {formatWeight(netWeight)}
              {weightTicket.netWeightRemarks && (
                <>
                  <h5 className={styles.remarksHeader}>Remarks</h5>
                  <p className={styles.remarks}>{weightTicket.netWeightRemarks}</p>
                </>
              )}
              <Button onClick={toggleEditForm} className={styles.editButton}>
                Edit
              </Button>
            </div>
          ) : (
            <EditPPMNetWeightForm
              initialValues={{ adjustedNetWeight: String(netWeight), netWeightRemarks: weightTicket.netWeightRemarks }}
              onSave={editNetWeight}
              onCancel={toggleEditForm}
            />
          )}
        </div>
      </FlexContainer>
    </div>
  );
};

EditPPMNetWeight.propTypes = {
  weightTicket: WeightTicketShape.isRequired,
  weightAllowance: number.isRequired,
  shipments: PropTypes.arrayOf(ShipmentShape).isRequired,
  editNetWeight: func.isRequired,
};

export default EditPPMNetWeight;
