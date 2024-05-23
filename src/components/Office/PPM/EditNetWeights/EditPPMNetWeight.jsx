import React, { useState } from 'react';
import { Formik } from 'formik';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import * as Yup from 'yup';
import { PropTypes, number } from 'prop-types';
import { Button, Fieldset, Label, Textarea } from '@trussworks/react-uswds';

import styles from './EditPPMNetWeight.module.scss';

import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { ErrorMessage } from 'components/form/ErrorMessage';
import { formatWeight } from 'utils/formatters';
import { calculateWeightTicketWeightDifference, getWeightTicketNetWeight } from 'utils/shipmentWeights';
import { calculateWeightRequested } from 'hooks/custom';
import { patchWeightTicket } from 'services/ghcApi';
import { ShipmentShape, WeightTicketShape } from 'types/shipment';
import { DOCUMENTS } from 'constants/queryKeys';

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

const EditPPMNetWeightForm = ({ onCancel, initialValues, weightTicket }) => {
  const validationSchema = Yup.object({
    adjustedNetWeight: Yup.number()
      .min(0, 'Net weight must be 0 lbs or greater')
      .lessThan(weightTicket.fullWeight, 'Net weight must be less than or equal to the full weight')
      .required('Required'),
    netWeightRemarks: Yup.string().nullable().required('Required'),
  });
  const queryClient = useQueryClient();

  const { mutate: patchWeightTicketMutation } = useMutation({
    mutationFn: patchWeightTicket,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: [DOCUMENTS],
      });
    },
  });

  /**
   * @const onSubmit
   * @description This function is used to submit the mini form represented by the EditPPMNetWeightForm.
   * @param {Object} formValues - The values that are returned from the EditPPMNetWeightForm component on click.
   * @param {string} formValues.adjustedNetWeight - The adjusted net weight as a string. This value needs to be parsed into an integer before mutation.
   * @param {string} formValues.netWeightRemarks - The net weight remarks.
   * */
  const onSubmit = (formValues /* , actions */) => {
    const payload = {
      adjustedNetWeight: parseInt(formValues.adjustedNetWeight, 10),
      netWeightRemarks: formValues.netWeightRemarks,
    };
    patchWeightTicketMutation(
      {
        ppmShipmentId: weightTicket.ppmShipmentId,
        weightTicketId: weightTicket.id,
        payload,
        eTag: weightTicket.eTag,
      },
      {
        onSuccess: () => {
          onCancel();
        },
      },
    );
  };

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ handleChange, handleSubmit, values, isValid, errors, touched, setTouched }) => (
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
              <Button onClick={handleSubmit} type="submit" disabled={!isValid}>
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

const EditPPMNetWeight = ({ weightTicket, weightAllowance, shipments }) => {
  const [showEditForm, setShowEditForm] = useState(false);

  const toggleEditForm = () => {
    setShowEditForm(!showEditForm);
  };

  // Original weight is the full weight - empty weight
  const originalWeight = calculateWeightTicketWeightDifference(weightTicket);
  // moveWeightTotal = Sum of all ppm weights + sum of all non-ppm shipments
  const moveWeightTotal = calculateWeightRequested(shipments);
  const excessWeight = moveWeightTotal - weightAllowance;
  const hasExcessWeight = Boolean(excessWeight > 0);
  const netWeight = getWeightTicketNetWeight(weightTicket);

  const toFitValue = hasExcessWeight ? -Math.min(excessWeight, netWeight) : null;
  const showWarning = Boolean(hasExcessWeight && !showEditForm);
  const showReduceWeight = Boolean(-originalWeight === toFitValue);
  return (
    <div className={styles.main_wrapper}>
      <div>
        <h4 className={styles.mainHeader}>PPM net weight</h4>
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
            <div data-testid="net-weight-display" className={styles.wrapper}>
              <div className={styles.netWeightDisplay}>{formatWeight(netWeight)}</div>
              {weightTicket.netWeightRemarks && (
                <>
                  <h5 className={styles.remarksHeader}>Remarks</h5>
                  <p className={styles.remarks}>{weightTicket.netWeightRemarks}</p>
                </>
              )}
            </div>
          ) : (
            <EditPPMNetWeightForm
              initialValues={{
                adjustedNetWeight: String(netWeight),
                netWeightRemarks: weightTicket.netWeightRemarks,
              }}
              weightTicket={weightTicket}
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
};

export default EditPPMNetWeight;
