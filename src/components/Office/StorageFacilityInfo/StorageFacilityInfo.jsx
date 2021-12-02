import React from 'react';
import { Fieldset, FormGroup, Label, TextInput, Grid } from '@trussworks/react-uswds';
import { Field } from 'formik';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ServicesCounselingShipmentForm/ServicesCounselingShipmentForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';

const StorageFacilityInfo = () => {
  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset className={styles.Fieldset}>
        <h2>Storage facility info</h2>
        <Grid row>
          <Grid col={12}>
            <FormGroup>
              <Label htmlFor="facilityName">Facility name</Label>
              <Field as={TextInput} id="facilityName" name="storageFacility.facilityName" />
            </FormGroup>
          </Grid>
        </Grid>

        <Grid row gap>
          <Grid col={6}>
            <MaskedTextField
              label="Phone"
              id="facilityPhone"
              name="storageFacility.phone"
              type="tel"
              minimum="12"
              mask="000{-}000{-}0000"
              optional
            />
          </Grid>
        </Grid>

        <Grid row>
          <Grid col={12}>
            <FormGroup>
              <Label htmlFor="facilityEmail" className={styles.Label}>
                Email
                <span className="float-right">Optional</span>
              </Label>
              <Field as={TextInput} id="facilityEmail" name="storageFacility.email" />
            </FormGroup>
          </Grid>
        </Grid>

        <Grid row gap>
          <Grid col={6}>
            <FormGroup>
              <Label htmlFor="facilityServiceOrderNumber" className={styles.Label}>
                Service order number
                <span className="float-right">Optional</span>
              </Label>
              <Field as={TextInput} id="facilityServiceOrderNumber" name="serviceOrderNumber" />
            </FormGroup>
          </Grid>
        </Grid>
      </Fieldset>
    </SectionWrapper>
  );
};

export default StorageFacilityInfo;
