import React, { useState } from 'react';
import { number } from 'prop-types';
import { Button, TextInput, Fieldset, Label, Textarea } from '@trussworks/react-uswds';

import styles from './EditBillableWeight.module.scss';

import { formatWeight } from 'shared/formatters';

export default function EditBillableWeight({ weightAllowance, estimatedWeight }) {
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
      <h5>Max billable weight</h5>
      <div>
        <strong>{formatWeight(weightAllowance)}</strong> <span>| weight allowance</span>
      </div>
      <div className={styles.estimatedWeight}>
        <strong>{formatWeight(estimatedWeight)}</strong> <span>| 110% of total estimated weight</span>
      </div>

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
  weightAllowance: number.isRequired,
  estimatedWeight: number.isRequired,
};
