import React from 'react';
import { Fieldset, FormGroup, Radio, Grid, Label } from '@trussworks/react-uswds';
import { useField } from 'formik';
import PropTypes from 'prop-types';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';

const ShipmentIncentiveAdvance = ({ estimatedIncentive }) => {
  const [inputProps, , helperProps] = useField('advanceRequested');
  const advanceRequested = !!inputProps.value;

  const formattedIncentive = ((estimatedIncentive || 0) / 100).toLocaleString('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0,
  });
  const maximumAdvance = (((estimatedIncentive || 0) * 0.6) / 100).toLocaleString('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 2,
    minimumFractionDigits: 2,
  });

  const handleAdvanceRequestedChange = (event) => {
    const selected = event.target.value;
    helperProps.setValue(selected === 'Yes');
  };

  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset className={styles.Fieldset}>
        <h2 className={styles.SectionHeader}>Incentive &amp; advance</h2>
        <h3>Estimated incentive: {formattedIncentive}</h3>

        <Grid row gap>
          <Grid col={12}>
            <FormGroup>
              <Label className={styles.Label}>Advance (AOA) requested?</Label>
              <Radio
                id="advanceRequestedYes"
                label="Yes"
                name="advanceRequested"
                value="Yes"
                title="Yes"
                checked={advanceRequested}
                onChange={handleAdvanceRequestedChange}
              />
              <Radio
                id="advanceRequestedNo"
                label="No"
                name="advanceRequested"
                value="No"
                title="No"
                checked={!advanceRequested}
                onChange={handleAdvanceRequestedChange}
              />
            </FormGroup>

            {advanceRequested && (
              <>
                <FormGroup>
                  <Label>Amount requested</Label>
                  <MaskedTextField
                    defaultValue="0"
                    name="amountRequested"
                    label="Amount requested"
                    id="amountRequested"
                    mask={Number}
                    scale={0} // digits after point, 0 for integers
                    signed={false} // disallow negative
                    thousandsSeparator=","
                    lazy={false} // immediate masking evaluation
                    prefix="$"
                  />
                </FormGroup>
                Maximum advance: {maximumAdvance}
              </>
            )}
          </Grid>
        </Grid>
      </Fieldset>
    </SectionWrapper>
  );
};

export default ShipmentIncentiveAdvance;

ShipmentIncentiveAdvance.propTypes = {
  estimatedIncentive: PropTypes.number,
};

ShipmentIncentiveAdvance.defaultProps = {
  estimatedIncentive: 0,
};
