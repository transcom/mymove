import React from 'react';
import { Fieldset, FormGroup, Label, TextInput, Grid } from '@trussworks/react-uswds';
import { Field } from 'formik';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';

const ShipmentWeightInput = () => {
  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset className={styles.Fieldset}>
        <h2 className={styles.SectionHeader}>Weight</h2>
        <Grid row gap>
          <Grid col={6}>
            <FormGroup>
              <Label htmlFor="ntsRecordedWeight" className={styles.Label}>
                Previous Recorded Weight (lbs)
                <span className="float-right">Optional</span>
              </Label>
              <Field as={TextInput} id="ntsRecordedWeight" name="ntsRecordedWeight" />
            </FormGroup>
          </Grid>
        </Grid>
      </Fieldset>
    </SectionWrapper>
  );
};

export default ShipmentWeightInput;
