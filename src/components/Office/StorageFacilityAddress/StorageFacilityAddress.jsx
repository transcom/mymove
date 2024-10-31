import { React, useState } from 'react';
import { Fieldset, FormGroup, Label, TextInput, Grid } from '@trussworks/react-uswds';
import { Formik, Field } from 'formik';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { AddressFields } from 'components/form/AddressFields/AddressFields';

const StorageFacilityAddress = () => {
  const [isLookupErrorVisible, setIsLookupErrorVisible] = useState(false);
  return (
    <Formik validateOnChange={false} validateOnMount>
      {({ values, setValues }) => {
        const handleLocationChange = (value) => {
          setValues(
            {
              ...values,
              storageFacility: {
                ...values.storageFacility,
                address: {
                  ...values.storageFacility.address,
                  city: value.city,
                  state: value.state ? value.state : '',
                  county: value.county,
                  postalCode: value.postalCode,
                },
              },
            },
            { shouldValidate: true },
          );

          if (!value.city || !value.state || !value.county || !value.postalCode) {
            setIsLookupErrorVisible(true);
          } else {
            setIsLookupErrorVisible(false);
          }
        };
        return (
          <SectionWrapper className={formStyles.formSection}>
            <Fieldset className={styles.Fieldset}>
              <h2 className={styles.SectionHeader}>Storage facility address</h2>
              <AddressFields
                name="storageFacility.address"
                zipCityEnabled
                zipCityError={isLookupErrorVisible}
                handleLocationChange={handleLocationChange}
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
      }}
    </Formik>
  );
};

export default StorageFacilityAddress;
