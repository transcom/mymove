import React, { useState } from 'react';
import { number, string } from 'prop-types';
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
        <strong>{formatWeight(weightAllowance)}</strong> <span>| weight allowance</span>
      </div>
      <div className={styles.hintText}>
        <strong>{formatWeight(estimatedWeight * 1.1)}</strong> <span>| 110% of total estimated weight</span>
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
  estimatedWeight,
  maxBillableWeight,
  originalWeight,
  title,
  totalBillableWeight,
  weightAllowance,
}) {
  const [showEditBtn, setShowEditBtn] = useState(true);

  function toggleEdit() {
    setShowEditBtn(!showEditBtn);
  }

  return showEditBtn ? (
    <Button className={styles.editBtn} onClick={toggleEdit}>
      Edit
    </Button>
  ) : (
    <div className={styles.container}>
      <h5>{title}</h5>
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
        <TextInput className={styles.maxBillableWeight} type="number" /> lbs
        <Label htmlFor="remarks">Remarks</Label>
        <Textarea data-testid="remarks" name="remarks" placeholder="" id="remarks" maxLength={500} />
      </Fieldset>
      <div className={styles.btnContainer}>
        <Button onClick={toggleEdit}>Save changes</Button>
        <Button onClick={toggleEdit} unstyled>
          Cancel
        </Button>
      </div>
    </div>
  );
}

EditBillableWeight.propTypes = {
  billableWeight: number,
  estimatedWeight: number,
  maxBillableWeight: number.isRequired,
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
};
