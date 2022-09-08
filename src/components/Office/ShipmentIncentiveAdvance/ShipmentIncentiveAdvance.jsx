import React from 'react';
import { Fieldset, FormGroup, Radio, Grid, Label } from '@trussworks/react-uswds';
import { useField } from 'formik';
import PropTypes from 'prop-types';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { calculateMaxAdvanceAndFormatAdvanceAndIncentive } from 'utils/incentives';

const ShipmentIncentiveAdvance = ({ estimatedIncentive }) => {
  const [advanceInput, , advanceHelper] = useField('advanceRequested');
  const advanceRequested = String(advanceInput.value) === 'true';

  const { formattedMaxAdvance, formattedIncentive } =
    calculateMaxAdvanceAndFormatAdvanceAndIncentive(estimatedIncentive);

  const handleHasRequestedAdvanceChange = (event) => {
    const selected = event.target.value;
    advanceHelper.setValue(selected === 'Yes');
  };

  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset className={styles.Fieldset}>
        <h2 className={styles.SectionHeader}>Incentive &amp; advance</h2>
        <h3 className={styles.NoSpacing}>Estimated incentive: ${formattedIncentive}</h3>

        <Grid row>
          <Grid col={12}>
            <FormGroup>
              <Label className={styles.Label}>Advance (AOA) requested?</Label>
              <Radio
                id="hasRequestedAdvanceYes"
                label="Yes"
                name="advanceRequested"
                value="Yes"
                title="Yes"
                checked={advanceRequested}
                onChange={handleHasRequestedAdvanceChange}
              />
              <Radio
                id="hasRequestedAdvanceNo"
                label="No"
                name="advanceRequested"
                value="No"
                title="No"
                checked={!advanceRequested}
                onChange={handleHasRequestedAdvanceChange}
              />
            </FormGroup>

            {advanceRequested && (
              <>
                <FormGroup>
                  <MaskedTextField
                    defaultValue="0"
                    name="advance"
                    label="Amount requested"
                    id="advanceAmountRequested"
                    mask={Number}
                    scale={0} // digits after point, 0 for integers
                    signed={false} // disallow negative
                    thousandsSeparator=","
                    lazy={false} // immediate masking evaluation
                    prefix="$"
                  />
                </FormGroup>

                <FormGroup>
                  <div className={styles.AdvanceText}>Maximum advance: ${formattedMaxAdvance}</div>
                </FormGroup>
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
