import React, { useState } from 'react';
import { func, number, string } from 'prop-types';
import { Formik } from 'formik';
import { Button, TextInput, Fieldset, Label, Textarea } from '@trussworks/react-uswds';

import styles from './EditBillableWeight.module.scss';

import { formatWeight } from 'shared/formatters';

function BillableWeightHintText({
  billableWeight,
  estimatedWeight,
  maxBillableWeight,
  originalWeight,
  totalBillableWeight,
}) {
  const showToFit = billableWeight > maxBillableWeight && billableWeight < estimatedWeight * 1.1;

  return (
    <>
      <div>
        <strong>{formatWeight(originalWeight)}</strong> <span>| original weight</span>
      </div>
      <div className={styles.hintText}>
        <strong>{formatWeight(estimatedWeight * 1.1)}</strong> <span>| 110% of total estimated weight</span>
      </div>
      {showToFit && (
        <div className={styles.hintText}>
          <strong>{formatWeight(totalBillableWeight - billableWeight)}</strong>{' '}
          <span>| to fit within max billable weight</span>
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
};

BillableWeightHintText.defaultProps = {
  billableWeight: null,
  estimatedWeight: null,
  originalWeight: null,
  totalBillableWeight: null,
};

function MaxBillableWeightHintText({ weightAllowance, estimatedWeight }) {
  return (
    <>
      <div>
        <strong data-testid="maxWeight-weightAllowance">{formatWeight(weightAllowance)}</strong>{' '}
        <span>| weight allowance</span>
      </div>
      <div className={styles.hintText}>
        <strong data-testid="maxWeight-estimatedWeight">{formatWeight(estimatedWeight * 1.1)}</strong>{' '}
        <span>| 110% of total estimated weight</span>
      </div>
    </>
  );
}

MaxBillableWeightHintText.propTypes = {
  estimatedWeight: number,
  weightAllowance: number,
};

MaxBillableWeightHintText.defaultProps = {
  estimatedWeight: null,
  weightAllowance: null,
};
export default function EditBillableWeight({
  billableWeight,
  billableWeightJustification,
  editMTOShipment,
  estimatedWeight,
  maxBillableWeight,
  originalWeight,
  title,
  totalBillableWeight,
  weightAllowance,
}) {
  const [showEditBtn, setShowEditBtn] = useState(true);

  const toggleEdit = () => {
    setShowEditBtn(!showEditBtn);
  };

  const initialValues = { billableWeight, billableWeightJustification };

  return (
    <div className={styles.wrapper}>
      <h4 className={styles.header}>{title}</h4>
      {showEditBtn ? (
        <>
          <span>{billableWeight ? formatWeight(billableWeight) : formatWeight(maxBillableWeight)}</span>
          {billableWeightJustification && (
            <>
              <h5 className={styles.remarksHeader}>Remarks</h5>
              <p className={styles.remarks}>{billableWeightJustification}</p>
            </>
          )}
          <Button className={styles.editBtn} onClick={toggleEdit}>
            Edit
          </Button>
        </>
      ) : (
        <Formik initialValues={initialValues}>
          {({ handleChange, values }) => (
            <div className={styles.container}>
              {billableWeight ? (
                <BillableWeightHintText
                  billableWeight={billableWeight}
                  estimatedWeight={estimatedWeight}
                  maxBillableWeight={maxBillableWeight}
                  originalWeight={originalWeight}
                  totalBillableWeight={totalBillableWeight}
                />
              ) : (
                <MaxBillableWeightHintText weightAllowance={weightAllowance} estimatedWeight={estimatedWeight} />
              )}
              <Fieldset className={styles.fieldset}>
                <TextInput
                  className={styles.maxBillableWeight}
                  defaultValue={maxBillableWeight}
                  id="billableWeight"
                  onChange={handleChange}
                  type="number"
                  value={values.billableWeight}
                />{' '}
                lbs
                <Label htmlFor="remarks">Remarks</Label>
                <Textarea
                  data-testid="remarks"
                  id="billableWeightJustification"
                  maxLength={500}
                  onChange={handleChange}
                  placeholder=""
                  value={values.billableWeightJustification}
                />
              </Fieldset>
              <div className={styles.btnContainer}>
                <Button
                  onClick={() => {
                    editMTOShipment({
                      ...initialValues,
                      ...values,
                    });
                    toggleEdit();
                  }}
                >
                  Save changes
                </Button>
                <Button onClick={toggleEdit} unstyled>
                  Cancel
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
  editMTOShipment: func.isRequired,
  estimatedWeight: number,
  maxBillableWeight: number,
  originalWeight: number,
  title: string.isRequired,
  totalBillableWeight: number,
  weightAllowance: number,
};

EditBillableWeight.defaultProps = {
  billableWeight: null,
  estimatedWeight: null,
  originalWeight: null,
  totalBillableWeight: null,
  weightAllowance: null,
  maxBillableWeight: null,
  billableWeightJustification: '',
};
