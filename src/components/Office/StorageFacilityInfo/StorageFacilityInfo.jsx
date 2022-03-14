import React from 'react';
import { Fieldset, FormGroup, Grid } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import TextField from 'components/form/fields/TextField/TextField';
import { officeRoles, roleTypes } from 'constants/userRoles';

const StorageFacilityInfo = ({ userRole }) => {
  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset className={styles.Fieldset}>
        <h2 className={styles.SectionHeader}>Storage facility info</h2>
        <Grid row>
          <Grid col={12}>
            <TextField label="Facility name" id="facilityName" name="storageFacility.facilityName" />
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
            <TextField label="Email" id="facilityEmail" name="storageFacility.email" optional />
          </Grid>
        </Grid>

        <Grid row gap>
          <Grid col={6}>
            <FormGroup>
              <TextField
                label="Service order number"
                id="facilityServiceOrderNumber"
                name="serviceOrderNumber"
                optional={userRole !== roleTypes.TOO}
              />
            </FormGroup>
          </Grid>
        </Grid>
      </Fieldset>
    </SectionWrapper>
  );
};

StorageFacilityInfo.propTypes = {
  userRole: PropTypes.oneOf(officeRoles).isRequired,
};

export default StorageFacilityInfo;
