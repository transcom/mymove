import React from 'react';
import { Fieldset, FormGroup, Label, TextInput, Grid } from '@trussworks/react-uswds';
import { Field } from 'formik';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { AddressFields } from 'components/form/AddressFields/AddressFields';

const StorageFacilityAddress = () => {
  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset className={styles.Fieldset}>
        <h2 className={styles.SectionHeader}>Storage facility address</h2>
        <AddressFields
          name="storageFacility.address"
          render={(fields) => (
            <>
              {fields}
              <Grid row gap>
                <Grid col={6}>
                  <FormGroup>
                    <Label htmlFor="facilityLotNumber" className={styles.Label}>
                      Lot number
                      <span className="float-right">Optional</span>
                    </Label>
                    <Field as={TextInput} id="facilityLotNumber" name="storageFacility.lotNumber" />
                  </FormGroup>
                </Grid>
              </Grid>
            </>
          )}
        />
      </Fieldset>
    </SectionWrapper>
  );
};

export default StorageFacilityAddress;
