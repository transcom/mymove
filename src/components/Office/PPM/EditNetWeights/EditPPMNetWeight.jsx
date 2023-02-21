import React, { useState } from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
// import { func, number, string, bool } from 'prop-types';
import { Button, Fieldset, Label, Textarea } from '@trussworks/react-uswds';

import styles from './EditPPMNetWeight.module.scss';

import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { ErrorMessage } from 'components/form/ErrorMessage';
import { formatWeight } from 'utils/formatters';

// Labels & constants

const CALCULATION_TYPE = {
  NET_WEIGHT: 'NET_WEIGHT',
  EXCESS_WEIGHT: 'EXCESS_WEIGHT',
};
const weightLabels = {
  [CALCULATION_TYPE.NET_WEIGHT]: {
    firstLabel: ' | original weight',
    secondLabel: ' | to fit within weight allowance',
  },
  [CALCULATION_TYPE.EXCESS_WEIGHT]: {
    firstLabel: 'Move Weight (total)',
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

const WeightCalculation = ({ type, firstValue, secondValue, thirdValue }) => {
  const { firstLabel, secondLabel, thirdLabel } = weightLabels[type];
  return (
    <div className={styles.wrapper}>
      <div className={styles.calculations}>
        <strong className={styles.value}>{formatWeight(firstValue)}</strong>
        <span className={styles.label}> {firstLabel}</span>
      </div>
      <div className={styles.calculations}>
        <strong className={styles.value}>{formatWeight(secondValue)}</strong>
        <span className={styles.label}> {secondLabel}</span>
      </div>
      {thirdValue && (
        <>
          <hr className={styles.divider} />
          <div className={styles.calculations}>
            <strong className={styles.value}>{formatWeight(thirdValue)}</strong>
            <span className={styles.label}> {thirdLabel}</span>
          </div>
        </>
      )}
    </div>
  );
};

const validationSchema = Yup.object({
  ppmNetWeightRemarks: Yup.string().required('Required'),
});

const EditPPMNetWeightForm = ({ onSave, onCancel, initialValues }) => (
  <Formik initialValues={initialValues} validationSchema={validationSchema}>
    {({ handleChange, values, isValid, errors, touched, setTouched }) => (
      <div>
        <Fieldset className={styles.fieldset}>
          <MaskedTextField
            defaultValue="0"
            id="ppmNetWeight"
            name="ppmNetWeight"
            mask={Number}
            lazy={false}
            scale={0}
            signed={false} // no negative numbers
            thousandsSeparator=","
            inputClassName={styles.weightInput}
            errorClassName={styles.errorMessage}
            labelClassName={styles.weightLabel}
          >
            {' '}
            lbs
          </MaskedTextField>
          <Label htmlFor="remarks">Remarks</Label>
          <ErrorMessage
            className={styles.errorMessage}
            display={!!touched.ppmNetWeightRemarks && !!errors.ppmNetWeightRemarks}
          >
            {errors.ppmNetWeightRemarks}
          </ErrorMessage>
          <Textarea
            id="ppmNetWeightRemarks"
            maxLength={500}
            placeholder=""
            onChange={handleChange}
            onBlur={() => {
              setTouched({ ppmNetWeightRemarks: true }, false);
            }}
            values={values.ppmNetWeightRemarks}
          />
          <FlexContainer className={styles.wrapper}>
            <Button onClick={onSave} disabled={!isValid}>
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

const EditPPMNetWeight = ({
  billableWeight,
  estimatedWeight,
  originalWeight,
  maxBillableWeight,
  totalBillableWeight,
  weightAllowance,
}) => {
  const [showEditForm, setShowEditForm] = useState(false);
  const toggleEditForm = () => {
    setShowEditForm(!showEditForm);
  };
  const toFitValue = maxBillableWeight - totalBillableWeight + billableWeight;
  const excessWeight = totalBillableWeight - weightAllowance;
  const hasExcessWeight = Boolean(excessWeight > 0);
  return (
    <div className={styles.wrapper}>
      <div>
        <h4 className={styles.header}>Edit PPM net weight</h4>
        {Boolean(showEditForm && hasExcessWeight) && (
          <WeightCalculation
            firstValue={totalBillableWeight}
            secondValue={weightAllowance}
            thirdValue={excessWeight}
            type={CALCULATION_TYPE.EXCESS_WEIGHT}
          />
        )}
      </div>
      <FlexContainer className={styles.netWeightContainer}>
        {hasExcessWeight && <div className={styles.warnings} />}
        <div>
          <h5 className={styles.header}>Net Weight</h5>
          <WeightCalculation firstValue={originalWeight} secondValue={toFitValue} type={CALCULATION_TYPE.NET_WEIGHT} />
          {!showEditForm ? (
            <div className={styles.wrapper}>
              {formatWeight(originalWeight)}
              <Button onClick={toggleEditForm} className={styles.editButton}>
                Edit
              </Button>
            </div>
          ) : (
            <EditPPMNetWeightForm
              initialValues={{ ppmNetWeight: String(originalWeight), ppmNetWeightRemarks: '' }}
              onCancel={toggleEditForm}
            />
          )}
        </div>
      </FlexContainer>
    </div>
  );
};

export default EditPPMNetWeight;
