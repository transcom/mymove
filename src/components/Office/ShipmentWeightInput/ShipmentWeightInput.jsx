import React from 'react';
import { Fieldset, FormGroup, Label, TextInput, Grid } from '@trussworks/react-uswds';
import { Field } from 'formik';

import styles from './ShipmentWeightInput.module.scss';

import formStyles from 'styles/form.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';

const ShipmentWeightInput = () => {
  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset>
        <h2>Weight</h2>
        <Grid col={5}>
          <FormGroup>
            <Label htmlFor="primeActualWeight" className={styles.Label}>
              Shipment weight (lbs)
              <span className="float-right">Optional</span>
            </Label>
            <Field as={TextInput} id="primeActualWeight" name="primeActualWeight" />
          </FormGroup>
        </Grid>
      </Fieldset>
    </SectionWrapper>
  );
};

export default ShipmentWeightInput;
