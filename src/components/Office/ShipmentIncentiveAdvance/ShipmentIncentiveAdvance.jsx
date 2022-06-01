import React from 'react';
import { Fieldset, FormGroup, Radio, Grid, Label } from '@trussworks/react-uswds';
import { useField } from 'formik';
import PropTypes from 'prop-types';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { calculateMaxAdvanceAndFormatAdvanceAndIncentive, getFormattedMaxAdvancePercentage } from 'utils/incentives';

const ShipmentIncentiveAdvance = ({ estimatedIncentive }) => {
  const [advanceInput, , advanceHelper] = useField('hasRequestedAdvance');
  const hasRequestedAdvance = !!advanceInput.value;
  const [amountInput] = useField('advanceAmountRequested');
  const advanceAmountRequested = Number(amountInput.value || '0');

  const { maxAdvance, formattedMaxAdvance, formattedIncentive } =
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
                name="hasRequestedAdvance"
                value="Yes"
                title="Yes"
                checked={hasRequestedAdvance}
                onChange={handleHasRequestedAdvanceChange}
              />
              <Radio
                id="hasRequestedAdvanceNo"
                label="No"
                name="hasRequestedAdvance"
                value="No"
                title="No"
                checked={!hasRequestedAdvance}
                onChange={handleHasRequestedAdvanceChange}
              />
            </FormGroup>

            {hasRequestedAdvance && (
              <>
                <FormGroup>
                  <MaskedTextField
                    defaultValue="0"
                    name="advanceAmountRequested"
                    label="Amount requested"
                    id="advanceAmountRequested"
                    mask={Number}
                    scale={0} // digits after point, 0 for integers
                    signed={false} // disallow negative
                    thousandsSeparator=","
                    lazy={false} // immediate masking evaluation
                    prefix="$"
                    error={advanceAmountRequested > maxAdvance}
                    errorMessage={
                      advanceAmountRequested > maxAdvance
                        ? `Enter an amount that is less than or equal to the maximum advance (${getFormattedMaxAdvancePercentage()} of estimated incentive)`
                        : ''
                    }
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
