import React, { useState } from 'react';
import * as Yup from 'yup';
import { func, number, string, bool } from 'prop-types';
import { Formik } from 'formik';
import { Button, Fieldset, Label, Textarea } from '@trussworks/react-uswds';

import styles from './EditBillableWeight.module.scss';

import { ErrorMessage } from 'components/form/ErrorMessage';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { formatWeight } from 'utils/formatters';
import RequiredAsterisk, { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

function BillableWeightHintText({
  billableWeight,
  estimatedWeight,
  maxBillableWeight,
  originalWeight,
  totalBillableWeight,
  isNTSRShipment,
}) {
  const estimatedWeightTimes110 = Math.round(estimatedWeight * 1.1);

  const showToFit = billableWeight > maxBillableWeight && billableWeight > estimatedWeightTimes110;
  // the to fit value is the max billable weight minus the total billable weight, excludes the shipment currently in view
  const toFitValue = maxBillableWeight - totalBillableWeight + billableWeight;

  return (
    <>
      <div className={styles.hintText}>
        <strong>{formatWeight(originalWeight)}</strong> <span>| original weight</span>
      </div>
      {!isNTSRShipment && estimatedWeightTimes110 !== 0 && (
        <div className={styles.hintText}>
          <strong>{formatWeight(estimatedWeightTimes110)}</strong> <span>| 110% of total estimated weight</span>
        </div>
      )}
      {showToFit && (
        <div className={styles.hintText}>
          <strong>{formatWeight(toFitValue)}</strong> <span>| to fit within max billable weight</span>
        </div>
      )}
    </>
  );
}

BillableWeightHintText.propTypes = {
  billableWeight: number,
  estimatedWeight: number,
  maxBillableWeight: number.isRequired,
  originalWeight: number,
  totalBillableWeight: number,
  isNTSRShipment: bool,
};

BillableWeightHintText.defaultProps = {
  billableWeight: null,
  estimatedWeight: null,
  originalWeight: null,
  totalBillableWeight: null,
  isNTSRShipment: false,
};

function MaxBillableWeightHintText({ estimatedWeight, isNTSRShipment }) {
  return (
    !Number.isNaN(estimatedWeight) &&
    estimatedWeight &&
    !isNTSRShipment && (
      <div className={styles.hintText}>
        <strong data-testid="maxWeight-estimatedWeight">{formatWeight(estimatedWeight * 1.1)}</strong>{' '}
        <span>| 110% of total estimated weight</span>
      </div>
    )
  );
}

MaxBillableWeightHintText.propTypes = {
  estimatedWeight: number,
  isNTSRShipment: bool,
};

MaxBillableWeightHintText.defaultProps = {
  estimatedWeight: null,
  isNTSRShipment: false,
};

const validationSchema = Yup.object({
  billableWeight: Yup.number().min(1, 'Authorized weight must be greater than or equal to 1').required('Required'),
  billableWeightJustification: Yup.string().required('Required'),
});
export default function EditBillableWeight({
  billableWeight,
  billableWeightJustification,
  editEntity,
  estimatedWeight,
  maxBillableWeight,
  originalWeight,
  showFieldsInitial,
  title,
  totalBillableWeight,
  isNTSRShipment,
}) {
  const [showFields, setShowFields] = useState(showFieldsInitial);
  const toggleEdit = () => {
    setShowFields((show) => !show);
  };

  const initialValues = {
    // Check for billable weight first since a maxBillableWeight will always exist, this avoids
    // a bug caused by short-circuiting where the value defaults to the maxBillableWeight
    billableWeight: (billableWeight && String(billableWeight)) || (maxBillableWeight && String(maxBillableWeight)), // Formik is expecting these weights as a string
    billableWeightJustification,
  };

  return (
    <div className={styles.wrapper} data-testid="maxBillableWeightEdit">
      <h4 className={styles.header}>{title}</h4>
      {!showFields ? (
        <>
          <span data-testid="billableWeightValue">
            {billableWeight ? formatWeight(billableWeight) : formatWeight(maxBillableWeight)}
          </span>
          {billableWeightJustification && (
            <>
              <h5 className={styles.remarksHeader}>Remarks</h5>
              <p data-testid="billableWeightRemarks" className={styles.remarks}>
                {billableWeightJustification}
              </p>
            </>
          )}
          <Button className={styles.editBtn} onClick={toggleEdit}>
            Edit
          </Button>
        </>
      ) : (
        <Formik enableReinitialize initialValues={initialValues} validationSchema={validationSchema}>
          {({ handleChange, values, isValid, errors, touched, setTouched }) => (
            <div className={styles.container}>
              {billableWeight ? (
                <BillableWeightHintText
                  billableWeight={billableWeight}
                  estimatedWeight={estimatedWeight}
                  maxBillableWeight={maxBillableWeight}
                  originalWeight={originalWeight}
                  totalBillableWeight={totalBillableWeight}
                  isNTSRShipment={isNTSRShipment}
                />
              ) : (
                <MaxBillableWeightHintText estimatedWeight={estimatedWeight} isNTSRShipment={isNTSRShipment} />
              )}
              {requiredAsteriskMessage}
              <Fieldset className={styles.fieldset}>
                <MaskedTextField
                  defaultValue="0"
                  inputClassName={styles.maxBillableWeight}
                  inputTestId="textInput"
                  errorClassName={styles.errorMessage}
                  labelClassName={styles.label}
                  id="billableWeight"
                  lazy={false} // immediate masking evaluation
                  mask={Number}
                  name="billableWeight"
                  scale={0} // digits after point, 0 for integers
                  signed={false} // disallow negative
                  thousandsSeparator=","
                  required
                  aria-label="Billable weight"
                >
                  {' '}
                  <span data-testid="billable-weight-label">
                    lbs <RequiredAsterisk />
                  </span>
                </MaskedTextField>
                <Label htmlFor="remarks">
                  <span data-testid="remarks-label">
                    Remarks <RequiredAsterisk />
                  </span>
                </Label>
                <ErrorMessage
                  className={styles.errorMessage}
                  display={!!touched.billableWeightJustification && !!errors.billableWeightJustification}
                >
                  {errors.billableWeightJustification}
                </ErrorMessage>
                <div
                  className={
                    !!touched.billableWeightJustification && !!errors.billableWeightJustification
                      ? 'usa-form-group--error'
                      : ''
                  }
                >
                  <Textarea
                    data-testid="remarks"
                    id="billableWeightJustification"
                    maxLength={500}
                    onChange={handleChange}
                    placeholder=""
                    onBlur={() => setTouched({ billableWeightJustification: true }, false)}
                    value={values.billableWeightJustification}
                    aria-label="Remarks"
                    required
                  />
                </div>
              </Fieldset>
              <div className={styles.btnContainer}>
                <Button type="button" onClick={toggleEdit} secondary>
                  Cancel
                </Button>
                <Button
                  disabled={!isValid}
                  onClick={() => {
                    editEntity({
                      ...initialValues,
                      ...values,
                    });
                    toggleEdit();
                  }}
                >
                  Save changes
                </Button>
              </div>
            </div>
          )}
        </Formik>
      )}
    </div>
  );
}

EditBillableWeight.propTypes = {
  billableWeight: number,
  billableWeightJustification: string,
  editEntity: func.isRequired,
  estimatedWeight: number,
  maxBillableWeight: number,
  originalWeight: number,
  showFieldsInitial: bool,
  title: string.isRequired,
  totalBillableWeight: number,
  isNTSRShipment: bool,
};

EditBillableWeight.defaultProps = {
  billableWeight: null,
  billableWeightJustification: '',
  estimatedWeight: null,
  maxBillableWeight: null,
  originalWeight: null,
  showFieldsInitial: false,
  totalBillableWeight: null,
  isNTSRShipment: false,
};
