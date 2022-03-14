import React from 'react';
import { Fieldset, FormGroup, Grid } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import TextField from 'components/form/fields/TextField/TextField';
import { officeRoles, roleTypes } from 'constants/userRoles';

const ShipmentWeightInput = ({ userRole }) => {
  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset className={styles.Fieldset}>
        <h2 className={styles.SectionHeader}>Weight</h2>
        <Grid row gap>
          <Grid col={6}>
            <FormGroup>
              <TextField
                label="Previously recorded weight (lbs)"
                name="ntsRecordedWeight"
                id="ntsRecordedWeight"
                optional={userRole !== roleTypes.TOO}
              />
            </FormGroup>
          </Grid>
        </Grid>
      </Fieldset>
    </SectionWrapper>
  );
};

ShipmentWeightInput.propTypes = {
  userRole: PropTypes.oneOf(officeRoles).isRequired,
};

export default ShipmentWeightInput;
